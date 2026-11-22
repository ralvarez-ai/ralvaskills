package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.ralva.dev/lib/term"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/sourceops"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
)

// runSource runs `rsk source <args>` against a fresh command tree and an
// isolated config folder. The config folder override is process-wide, so these
// tests are sequential.
func runSource(t *testing.T, args ...string) (stdout string, err error) {
	t.Helper()
	cmds := newSourceCmds()
	var out bytes.Buffer
	cmds.root.SetOut(&out)
	cmds.root.SetErr(io.Discard)
	cmds.root.SilenceUsage, cmds.root.SilenceErrors = true, true
	cmds.root.SetArgs(args)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = cmds.root.ExecuteContext(ctx)
	return out.String(), err
}

func isolateConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	config.SetConfigHome(dir)
	t.Cleanup(func() { config.SetConfigHome("") })
	return dir
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceAddDryRunWritesNothingAndPreviewsJSON(t *testing.T) {
	dir := isolateConfig(t)
	out, err := runSource(t, "add", "@acme/skills-private", "--roots", "skills", "--dry-run", "-o", "json")
	if err != nil {
		t.Fatal(err)
	}
	var change sourceChange
	if err = json.Unmarshal([]byte(out), &change); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out)
	}
	if change.Action != "add" || !change.DryRun || change.Applied || change.Source.ID != "@acme/skills-private" ||
		change.Source.Trust != trustNone || !strings.Contains(change.Next, "rsk source trust @acme/skills-private") {
		t.Errorf("change = %+v", change)
	}
	if _, statErr := os.Stat(filepath.Join(dir, sources.FileName)); statErr == nil {
		t.Error("a dry run wrote sources.json")
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceAddRegistersUntrustedAndListShowsIt(t *testing.T) {
	dir := isolateConfig(t)
	if _, err := runSource(t, "add", "@acme/skills-private", "--yes", "--rev", "main"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, trust.FileName)); err == nil {
		t.Error("adding a source granted trust")
	}

	out, err := runSource(t, "list", "-o", "json")
	if err != nil {
		t.Fatal(err)
	}
	var views []sourceView
	if err = json.Unmarshal([]byte(out), &views); err != nil {
		t.Fatalf("list output: %v\n%s", err, out)
	}
	if len(views) != 1 || views[0].ID != "@acme/skills-private" || views[0].Trust != trustNone ||
		views[0].URL != "git@github.com:acme/skills-private.git" || views[0].Rev != "main" {
		t.Errorf("views = %+v", views)
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceAddRefusesWhatShouldNotBeDeclared(t *testing.T) {
	isolateConfig(t)
	for name, args := range map[string][]string{
		"reserved identity":  {"add", "@anthropics/skills", "--yes"},
		"hostile transport":  {"add", "ext::sh -c id", "--id", "@a/b", "--yes"},
		"host needs an id":   {"add", "git@git.example.com:t/s.git", "--yes"},
		"credentials in url": {"add", "https://u:ghp_abcdefghijklmnopqrstuvwxyz0123456789@github.com/a/b.git", "--yes"},
		"traversing root":    {"add", "@acme/x", "--roots", "../etc", "--yes"},
	} {
		if out, err := runSource(t, args...); err == nil {
			t.Errorf("%s was accepted:\n%s", name, out)
		}
	}
	if st, err := loadSourceState(); err != nil || len(st.file.Sources) != 0 {
		t.Errorf("a refused add left state behind: %+v, %v", st.file, err)
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceAddPathSource(t *testing.T) {
	isolateConfig(t)
	skills := t.TempDir()
	if _, err := runSource(t, "add", skills, "--id", "@me/scratch", "--yes"); err != nil {
		t.Fatal(err)
	}
	out, err := runSource(t, "list", "-o", "json")
	if err != nil {
		t.Fatal(err)
	}
	var views []sourceView
	if err = json.Unmarshal([]byte(out), &views); err != nil || len(views) != 1 || views[0].Trust != trustNA {
		t.Errorf("views = %+v, %v", views, err)
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceTrustNeedsATerminalAndSaysWhatToRun(t *testing.T) {
	dir := isolateConfig(t)
	if _, err := runSource(t, "add", "@acme/skills-private", "--yes"); err != nil {
		t.Fatal(err)
	}
	if term.StdioIsTerminal() {
		t.Skip("the tests are running on a terminal")
	}
	_, err := runSource(t, "trust", "@acme/skills-private")
	if err == nil {
		t.Fatal("trust succeeded without a terminal")
	}
	for _, want := range []string{"terminal", "rsk source trust @acme/skills-private"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q lacks %q", err, want)
		}
	}
	if _, statErr := os.Stat(filepath.Join(dir, trust.FileName)); statErr == nil {
		t.Error("a refused trust still wrote trust.json")
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceTestFailsClosedWhenUntrustedAndNeverContactsTheHost(t *testing.T) {
	isolateConfig(t)
	if _, err := runSource(t, "add", "@acme/skills-private", "--yes"); err != nil {
		t.Fatal(err)
	}
	out, err := runSource(t, "test", "@acme/skills-private", "-o", "json")
	if err == nil || !strings.Contains(err.Error(), "rsk source trust @acme/skills-private") {
		t.Fatalf("err = %v, want the trust hint", err)
	}
	var res sourceTest
	if jsonErr := json.Unmarshal([]byte(out), &res); jsonErr != nil || res.Reachable || res.Error == "" {
		t.Errorf("result = %+v, %v", res, jsonErr)
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceTestReportsAnUnreachableTrustedSource(t *testing.T) {
	dir := isolateConfig(t)
	url := "https://127.0.0.1:1/team/skills.git"
	if _, err := runSource(t, "add", url, "--id", "@team/skills", "--yes"); err != nil {
		t.Fatal(err)
	}
	grants := trust.Grant(trust.File{}, "@team/skills", url, time.Now())
	if err := trust.Save(filepath.Join(dir, trust.FileName), grants); err != nil {
		t.Fatal(err)
	}

	out, err := runSource(t, "test", "@team/skills", "-o", "json")
	if err == nil {
		t.Fatal("an unreachable source tested fine")
	}
	var res sourceTest
	if jsonErr := json.Unmarshal([]byte(out), &res); jsonErr != nil || res.Reachable || res.Trust != trustTrusted {
		t.Errorf("result = %+v, %v", res, jsonErr)
	}
	if strings.Contains(out, "ghp_") {
		t.Error("output leaked a credential")
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceTestSaysWhenGitRewritesTheURL(t *testing.T) {
	dir := isolateConfig(t)
	url := "https://127.0.0.1:1/team/skills.git"
	if _, err := runSource(t, "add", url, "--id", "@team/skills", "--yes"); err != nil {
		t.Fatal(err)
	}
	grants := trust.Grant(trust.File{}, "@team/skills", url, time.Now())
	if err := trust.Save(filepath.Join(dir, trust.FileName), grants); err != nil {
		t.Fatal(err)
	}
	effectiveURL = func(context.Context, git.RemoteURL) (string, error) { return "git@alias:team/skills.git", nil }
	t.Cleanup(func() { effectiveURL = git.EffectiveURL })

	out, _ := runSource(t, "test", "@team/skills", "-o", "json")
	var res sourceTest
	if err := json.Unmarshal([]byte(out), &res); err != nil || res.Via != "git@alias:team/skills.git" {
		t.Errorf("result = %+v, %v", res, err)
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceRemoveDropsTheEntryAndItsGrant(t *testing.T) {
	dir := isolateConfig(t)
	if _, err := runSource(t, "add", "@acme/skills-private", "--yes"); err != nil {
		t.Fatal(err)
	}
	url := "git@github.com:acme/skills-private.git"
	grants := trust.Grant(trust.File{}, "@acme/skills-private", url, time.Now())
	if err := trust.Save(filepath.Join(dir, trust.FileName), grants); err != nil {
		t.Fatal(err)
	}

	if _, err := runSource(t, "remove", "@acme/skills-private", "--dry-run"); err != nil {
		t.Fatal(err)
	}
	if st, _ := loadSourceState(); len(st.file.Sources) != 1 || !st.grants.Granted("@acme/skills-private", url) {
		t.Fatal("a dry run removed something")
	}

	if _, err := runSource(t, "remove", "@acme/skills-private", "--yes"); err != nil {
		t.Fatal(err)
	}
	st, err := loadSourceState()
	if err != nil || len(st.file.Sources) != 0 || st.grants.Granted("@acme/skills-private", url) {
		t.Errorf("state after remove = %+v, %v", st, err)
	}
	if _, err = runSource(t, "remove", "@acme/skills-private", "--yes"); err == nil {
		t.Error("removing a source that is not declared succeeded")
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceUntrust(t *testing.T) {
	dir := isolateConfig(t)
	grants := trust.Grant(trust.File{}, "@acme/x", "git@github.com:acme/x.git", time.Now())
	if err := trust.Save(filepath.Join(dir, trust.FileName), grants); err != nil {
		t.Fatal(err)
	}
	if _, err := runSource(t, "untrust", "@acme/x"); err != nil {
		t.Fatal(err)
	}
	if st, _ := loadSourceState(); st.grants.Granted("@acme/x", "git@github.com:acme/x.git") {
		t.Error("untrust left the grant")
	}
	if _, err := runSource(t, "untrust", "@acme/x"); err == nil {
		t.Error("untrusting twice succeeded")
	}
}

//nolint:paralleltest // shared cobra tree mutates on read; see root_test.go
func TestSourceCommandIsRegisteredGroupedAndGuarded(t *testing.T) {
	found, _, err := rootCmd.Find([]string{"source", "add"})
	if err != nil || found.Name() != "add" {
		t.Fatalf("rsk source add not reachable: %v", err)
	}
	if canonicalSource.root.GroupID != groupSources || canonicalSource.root.Hidden {
		t.Errorf("source: group=%q hidden=%v", canonicalSource.root.GroupID, canonicalSource.root.Hidden)
	}
	if !term.IsDestructive(canonicalSource.add) || !term.IsDestructive(canonicalSource.remove) {
		t.Error("add and remove must be marked destructive so the session asks first")
	}
	if sessionRunnable(canonicalSource.trust) {
		t.Error("the interactive session must never offer 'source trust'")
	}
	if !sessionRunnable(canonicalSource.list) {
		t.Error("source list should be offered by the session")
	}
	if canonicalSource.trust.Flags().Lookup("yes") != nil {
		t.Error("'source trust' must not accept --yes: trust is a human decision")
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceScanShowsAClashAndRegistersNothing(t *testing.T) {
	dir := isolateConfig(t)
	root := t.TempDir()
	writeCatalogSkill(t, filepath.Join(root, "a"), "same")
	writeCatalogSkill(t, filepath.Join(root, "b"), "same")
	writeCatalogSkill(t, root, "fine")

	out, err := runSource(t, "scan", root, "-o", "json")
	if err != nil {
		t.Fatalf("err = %v\n%s", err, out)
	}
	var res sourceops.ScanOutcome
	if err = json.Unmarshal([]byte(out), &res); err != nil || len(res.Clashes) != 1 || res.Clashes[0].Name != "same" ||
		len(res.Skills) != 1 || res.Skills[0].Name != "fine" {
		t.Errorf("scan = %+v, %v\n%s", res, err, out)
	}
	text, err := runSource(t, "scan", root)
	if err != nil || !strings.Contains(text, "same") || !strings.Contains(text, "a/same") ||
		!strings.Contains(text, "b/same") {
		t.Errorf("text scan = %q, %v", text, err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, sources.FileName)); statErr == nil {
		t.Error("scan registered the directory")
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceScanRefusesAnUnregisteredRepositoryAndAnUntrustedOne(t *testing.T) {
	isolateConfig(t)
	if _, err := runSource(
		t,
		"scan",
		"@acme/skills",
	); err == nil ||
		!strings.Contains(err.Error(), "rsk source add @acme/skills") {
		t.Errorf("an unregistered repository: err = %v", err)
	}
	if _, err := runSource(t, "add", "@acme/skills", "--yes"); err != nil {
		t.Fatal(err)
	}
	if _, err := runSource(
		t,
		"scan",
		"@acme/skills",
	); err == nil ||
		!strings.Contains(err.Error(), "rsk source trust @acme/skills") {
		t.Errorf("an untrusted repository: err = %v", err)
	}
}

//nolint:paralleltest // the config folder override is process-wide
func TestSourceAddPrivateIsRegisteredListedAndKeptOutOfTheCatalog(t *testing.T) {
	isolateConfig(t)
	if _, err := runSource(t, "add", "@acme/hidden", "--private", "--yes"); err != nil {
		t.Fatal(err)
	}
	out, err := runSource(t, "list", "-o", "json")
	if err != nil {
		t.Fatal(err)
	}
	var views []sourceView
	if err = json.Unmarshal([]byte(out), &views); err != nil || len(views) != 1 || !views[0].Private {
		t.Fatalf("views = %+v, %v", views, err)
	}
	if _, err = runSource(t, "add", "@acme/x", "--private", "--dry-run", "-o", "json"); err != nil {
		t.Errorf("a dry run of a private add failed: %v", err)
	}
}
