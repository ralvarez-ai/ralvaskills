package mcp_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/mcp"
	"go.ralva.dev/ralvaskills/v3/internal/resolve"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/source/sourcetest"
	"go.ralva.dev/ralvaskills/v3/internal/sourceops"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
)

const sourcesCommit = "cccccccccccccccccccccccccccccccccccccccc"

type sourcesWorld struct {
	dir     string
	global  string
	session *sdk.ClientSession
	repos   map[string]*sourcetest.Repo
	probed  []string
}

// newSourcesWorld is a server whose machine files live in a temp folder and
// whose git sources read an in-memory clone holding "alpha".
func newSourcesWorld(t *testing.T, readOnly bool) *sourcesWorld {
	t.Helper()
	w := &sourcesWorld{dir: t.TempDir(), repos: map[string]*sourcetest.Repo{}}
	w.global = filepath.Join(w.dir, "global")
	repo := filepath.Join(w.dir, "repo")
	writeSkill(t, filepath.Join(repo, "skills", "local-one"), "local-one", "1.0.0", "A local skill.")
	w.session = newSession(t, mcp.Deps{
		Cfg: config.Config{
			RepoPath:           repo,
			GlobalTargets:      map[string]string{"claude-code": w.global},
			DefaultTargetScope: "claude-code",
			OfficialCache:      filepath.Join(w.dir, "cache", "anthropic"),
		},
		SourcesPath: filepath.Join(w.dir, "sources.json"),
		TrustPath:   filepath.Join(w.dir, "trust.json"),
		ReadOnly:    readOnly,
		RepoFor: func(spec source.Spec) source.RepoOps {
			repo := &sourcetest.Repo{Commit: sourcesCommit, Skills: map[string]string{"alpha": "1.0.0"}}
			w.repos[spec.ID] = repo
			return repo
		},
		LsRemote: func(_ context.Context, remote git.RemoteURL, _ string) (string, error) {
			w.probed = append(w.probed, remote.String())
			return sourcesCommit, nil
		},
		Logger: slog.New(slog.DiscardHandler), Version: "test",
	})
	return w
}

func (w *sourcesWorld) call(t *testing.T, name string, args map[string]any, out any) *sdk.CallToolResult {
	t.Helper()
	res := callTool(t, w.session, name, args)
	if out != nil && !res.IsError {
		decode(t, res.StructuredContent, out)
	}
	return res
}

// fetched lists the sources whose remote was actually contacted.
func (w *sourcesWorld) fetched() []string {
	var ids []string
	for id, repo := range w.repos {
		if len(repo.Fetches) > 0 {
			ids = append(ids, id)
		}
	}
	return ids
}

func (w *sourcesWorld) exists(name string) bool {
	_, err := os.Stat(filepath.Join(w.dir, name))
	return err == nil
}

func (w *sourcesWorld) grant(t *testing.T) {
	t.Helper()
	grants := trust.Grant(trust.File{}, "@acme/skills", "git@github.com:acme/skills.git", time.Now())
	if err := trust.Save(filepath.Join(w.dir, "trust.json"), grants); err != nil {
		t.Fatal(err)
	}
}

func errorText(res *sdk.CallToolResult) string {
	var b strings.Builder
	for _, c := range res.Content {
		if tc, ok := c.(*sdk.TextContent); ok {
			b.WriteString(tc.Text)
		}
	}
	return b.String()
}

func TestAddSourceRegistersUntrustedAndGrantsNothing(t *testing.T) {
	t.Parallel()
	w := newSourcesWorld(t, false)

	var preview sourceops.Change
	res := w.call(t, "add_source", map[string]any{"target": "@acme/skills", "dry_run": true}, &preview)
	if res.IsError || preview.Applied || !preview.DryRun || preview.Source.Trust != sourceops.TrustNone ||
		!strings.Contains(preview.Next, "rsk source trust @acme/skills") {
		t.Fatalf("preview = %+v (error %v)", preview, errorText(res))
	}
	if w.exists("sources.json") {
		t.Fatal("a dry run wrote sources.json")
	}

	var done sourceops.Change
	if res = w.call(
		t,
		"add_source",
		map[string]any{"target": "@acme/skills", "roots": []string{"skills"}},
		&done,
	); res.IsError {
		t.Fatal(errorText(res))
	}
	if !done.Applied || !w.exists("sources.json") || w.exists("trust.json") {
		t.Errorf(
			"done = %+v, sources.json = %v, trust.json = %v",
			done,
			w.exists("sources.json"),
			w.exists("trust.json"),
		)
	}
	if res = w.call(t, "add_source", map[string]any{"target": "@acme/skills"}, nil); !res.IsError {
		t.Error("a source was added twice")
	}
	for _, hostile := range []string{"ext::sh -c id", "-oProxyCommand=x", "file:///etc"} {
		if res = w.call(t, "add_source", map[string]any{"target": hostile, "id": "@a/b"}, nil); !res.IsError {
			t.Errorf("add_source accepted %q", hostile)
		}
	}
	if len(w.fetched())+len(w.probed) != 0 {
		t.Errorf("adding a source contacted %v %v", w.fetched(), w.probed)
	}
}

func TestNoToolTakesATrustFlagASecretOrAConfirmation(t *testing.T) {
	t.Parallel()
	w := newSourcesWorld(t, false)
	tools, err := w.session.ListTools(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	forbidden := []string{
		"trust",
		"yes",
		"token",
		"secret",
		"password",
		"credential",
		"auth",
		"key",
		"private",
		"confirm",
		"force",
	}
	for _, tool := range tools.Tools {
		data, marshalErr := json.Marshal(tool.InputSchema)
		if marshalErr != nil {
			t.Fatal(marshalErr)
		}
		var schema struct {
			Properties map[string]any `json:"properties"`
		}
		if err = json.Unmarshal(data, &schema); err != nil {
			t.Fatal(err)
		}
		for prop := range schema.Properties {
			for _, bad := range forbidden {
				if strings.Contains(strings.ToLower(prop), bad) {
					t.Errorf("%s takes a %q parameter", tool.Name, prop)
				}
			}
		}
	}
}

func TestMutatingSourceToolsForceApprovalAndReadOnlyOnesDoNot(t *testing.T) {
	t.Parallel()
	w := newSourcesWorld(t, false)
	tools, err := w.session.ListTools(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"add_source", "remove_source", "install_skills"} {
		tool := findTool(t, tools.Tools, name)
		if v, ok := tool.Meta["anthropic/requiresUserInteraction"].(bool); !ok || !v || tool.Annotations.ReadOnlyHint {
			t.Errorf("%s does not force approval: %+v", name, tool.Meta)
		}
	}
	for _, name := range []string{"list_sources", "test_source", "scan_source", "resolve_reference", "project_profile", "search_skills"} {
		tool := findTool(t, tools.Tools, name)
		if !tool.Annotations.ReadOnlyHint || tool.Meta["anthropic/requiresUserInteraction"] != nil {
			t.Errorf("%s is not plainly read-only: %+v", name, tool)
		}
	}

	ro := newSourcesWorld(t, true)
	tools, err = ro.session.ListTools(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	got := make([]string, len(tools.Tools))
	for i, tool := range tools.Tools {
		got[i] = tool.Name
	}
	slices.Sort(got)
	want := []string{
		"list_sources",
		"project_profile",
		"resolve_reference",
		"scan_source",
		"search_skills",
		"test_source",
	}
	if !slices.Equal(got, want) {
		t.Errorf("read-only tools = %v, want %v", got, want)
	}
}

func TestRegisteredToolsAreTheOnesTheParityTestKnows(t *testing.T) {
	t.Parallel()
	w := newSourcesWorld(t, false)
	tools, err := w.session.ListTools(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	got := make([]string, len(tools.Tools))
	for i, tool := range tools.Tools {
		got[i] = tool.Name
	}
	slices.Sort(got)
	want := mcp.ToolNames()
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("registered %v, ToolNames %v", got, want)
	}
}

func TestListAndRemoveSource(t *testing.T) {
	t.Parallel()
	w := newSourcesWorld(t, false)
	w.call(t, "add_source", map[string]any{"target": "@acme/skills"}, nil)
	w.grant(t)

	var list struct {
		Sources []sourceops.View `json:"sources"`
	}
	w.call(t, "list_sources", nil, &list)
	trustOf := map[string]string{}
	for _, v := range list.Sources {
		trustOf[v.ID] = v.Trust
	}
	if trustOf["@acme/skills"] != sourceops.TrustTrusted || trustOf["@anthropics/skills"] != sourceops.TrustBuiltIn {
		t.Errorf("sources = %v", trustOf)
	}

	var preview sourceops.Change
	w.call(t, "remove_source", map[string]any{"id": "@acme/skills", "dry_run": true}, &preview)
	if preview.Applied || !w.exists("trust.json") {
		t.Fatalf("a dry run removed something: %+v", preview)
	}
	if res := w.call(t, "remove_source", map[string]any{"id": "@acme/skills"}, nil); res.IsError {
		t.Fatal(errorText(res))
	}
	w.call(t, "list_sources", nil, &list)
	for _, v := range list.Sources {
		if v.ID == "@acme/skills" {
			t.Error("the source is still listed")
		}
	}
	if res := w.call(t, "remove_source", map[string]any{"id": "@acme/skills"}, nil); !res.IsError {
		t.Error("removing an undeclared source succeeded")
	}
}

func TestTestSourceNamesTheTrustCommandAndNeverContactsAnUntrustedHost(t *testing.T) {
	t.Parallel()
	w := newSourcesWorld(t, false)
	w.call(t, "add_source", map[string]any{"target": "@acme/skills"}, nil)

	res := w.call(t, "test_source", map[string]any{"id": "@acme/skills"}, nil)
	if !res.IsError || !strings.Contains(errorText(res), "rsk source trust @acme/skills") || len(w.probed) != 0 {
		t.Fatalf("untrusted: isError=%v text=%q probed=%v", res.IsError, errorText(res), w.probed)
	}

	w.grant(t)
	var out sourceops.TestResult
	if res = w.call(t, "test_source", map[string]any{"id": "@acme/skills"}, &out); res.IsError || !out.Reachable ||
		out.Commit != sourcesCommit || len(w.probed) != 1 {
		t.Errorf("trusted: %+v %v", out, errorText(res))
	}
}

func TestScanSourceShowsClashesAndRefusesWhatItWouldHaveToCloneBlind(t *testing.T) {
	t.Parallel()
	w := newSourcesWorld(t, false)
	root := t.TempDir()
	writeSkill(t, filepath.Join(root, "a", "same"), "same", "1.0.0", "x")
	writeSkill(t, filepath.Join(root, "b", "same"), "same", "1.0.0", "x")
	writeSkill(t, filepath.Join(root, "fine"), "fine", "1.0.0", "x")

	var out sourceops.ScanOutcome
	if res := w.call(t, "scan_source", map[string]any{"target": root}, &out); res.IsError || len(out.Clashes) != 1 ||
		len(out.Skills) != 1 {
		t.Fatalf("scan = %+v %v", out, errorText(res))
	}
	if res := w.call(t, "scan_source", map[string]any{"target": "@acme/skills"}, nil); !res.IsError ||
		!strings.Contains(errorText(res), "rsk source add @acme/skills") || len(w.fetched()) != 0 {
		t.Errorf("an unregistered repository was scanned: %q reached %v", errorText(res), w.fetched())
	}
	w.call(t, "add_source", map[string]any{"target": "@acme/skills"}, nil)
	if res := w.call(
		t,
		"scan_source",
		map[string]any{"target": "@acme/skills"},
		nil,
	); !res.IsError ||
		len(w.fetched()) != 0 {
		t.Errorf("an untrusted repository was scanned: %q reached %v", errorText(res), w.fetched())
	}
}

func TestResolveReferenceExplainsWithoutOpeningAnUnregisteredSource(t *testing.T) {
	t.Parallel()
	w := newSourcesWorld(t, false)

	var out resolve.Result
	res := w.call(t, "resolve_reference", map[string]any{"reference": "@acme/skills/alpha"}, &out)
	if res.IsError || out.Status != resolve.UnknownSource ||
		!strings.Contains(out.Message, "git@github.com:acme/skills.git") {
		t.Fatalf("out = %+v %v", out, errorText(res))
	}
	if len(w.fetched())+len(w.probed) != 0 {
		t.Errorf("resolving an unregistered reference contacted %v %v", w.fetched(), w.probed)
	}

	w.call(t, "add_source", map[string]any{"target": "@acme/skills"}, nil)
	w.grant(t)
	if res = w.call(t, "resolve_reference", map[string]any{"reference": "@acme/skills/alpha"}, &out); res.IsError ||
		out.Status != resolve.Resolved || out.Skill == nil || out.Skill.Commit != sourcesCommit {
		t.Errorf("registered: %+v %v", out, errorText(res))
	}
	if res = w.call(t, "resolve_reference", map[string]any{"reference": "@bad"}, nil); !res.IsError {
		t.Error("a malformed reference was accepted")
	}
}

func TestInstallSkillsUsesATrustedDeclaredSourceAndStillRefusesTheRest(t *testing.T) {
	t.Parallel()
	w := newSourcesWorld(t, false)
	w.call(t, "add_source", map[string]any{"target": "@acme/skills"}, nil)

	res := w.call(t, "install_skills", map[string]any{"names": []string{"@acme/skills/alpha"}, "scope": "global"}, nil)
	if !res.IsError || !strings.Contains(errorText(res), "rsk source trust @acme/skills") {
		t.Fatalf("an untrusted source installed: %q", errorText(res))
	}

	w.grant(t)
	if res = w.call(
		t,
		"install_skills",
		map[string]any{"names": []string{"@acme/skills/alpha"}, "scope": "global"},
		nil,
	); res.IsError {
		t.Fatalf("a trusted declared source was refused: %q", errorText(res))
	}
	if _, err := os.Lstat(filepath.Join(w.global, "alpha")); err != nil {
		t.Errorf("not linked: %v", err)
	}
	if res = w.call(
		t,
		"install_skills",
		map[string]any{"names": []string{"@other/repo/alpha"}, "scope": "global"},
		nil,
	); !res.IsError {
		t.Error("an unregistered source installed")
	}
}

func TestProjectProfileListsTheConfiguredSources(t *testing.T) {
	t.Parallel()
	w := newSourcesWorld(t, false)
	w.call(t, "add_source", map[string]any{"target": "@acme/skills"}, nil)
	var out mcp.ProjectProfileOut
	if res := w.call(t, "project_profile", map[string]any{"path": t.TempDir()}, &out); res.IsError {
		t.Fatal(errorText(res))
	}
	if len(out.Sources) != 1 || out.Sources[0].ID != "@acme/skills" || out.Sources[0].Trust != sourceops.TrustNone {
		t.Errorf("sources = %+v, want the declared source and its trust state", out.Sources)
	}
}
