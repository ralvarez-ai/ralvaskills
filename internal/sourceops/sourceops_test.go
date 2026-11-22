package sourceops_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/source/sourcetest"
	"go.ralva.dev/ralvaskills/v3/internal/sourceops"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
)

const commit = "cccccccccccccccccccccccccccccccccccccccc"

func newEnv(t *testing.T) (env sourceops.Env, dir string) {
	t.Helper()
	dir = t.TempDir()
	cfg := config.Config{
		RegistryURL: "https://skills.example", OfficialCache: filepath.Join(dir, "cache", "anthropic"),
		GlobalTargets: map[string]string{"claude-code": filepath.Join(dir, "g")},
	}
	env = sourceops.Env{
		Cfg: &cfg, SourcesPath: filepath.Join(dir, "sources.json"), TrustPath: filepath.Join(dir, "trust.json"),
		RepoFor: func(source.Spec) source.RepoOps {
			return &sourcetest.Repo{Commit: commit, Skills: map[string]string{"alpha": "1.0.0"}}
		},
		LsRemote: func(context.Context, git.RemoteURL, string) (string, error) { return commit, nil },
	}
	return env, dir
}

func TestListShowsTheBuiltInsThenTheDeclaredSourcesWithTheirTrust(t *testing.T) {
	t.Parallel()
	env, dir := newEnv(t)
	scratch := t.TempDir()
	for _, target := range []string{"@acme/skills", "@acme/other"} {
		plan, apply, err := sourceops.PlanAdd(env, target, sources.Options{})
		if err != nil || plan.Applied {
			t.Fatalf("PlanAdd(%s) = %+v, %v", target, plan, err)
		}
		if err = apply(); err != nil {
			t.Fatal(err)
		}
	}
	_, apply, err := sourceops.PlanAdd(env, scratch, sources.Options{ID: "@me/scratch"})
	if err != nil || apply() != nil {
		t.Fatal(err)
	}
	grants := trust.Grant(trust.File{}, "@acme/skills", "git@github.com:acme/skills.git", time.Now())
	if err = trust.Save(filepath.Join(dir, "trust.json"), grants); err != nil {
		t.Fatal(err)
	}

	views, err := sourceops.List(env)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	order := make([]string, 0, len(views))
	for _, v := range views {
		got[v.ID], order = v.Trust, append(order, v.ID)
	}
	want := map[string]string{
		skill.DefaultRegistryID: sourceops.TrustBuiltIn, skill.OfficialID: sourceops.TrustBuiltIn,
		"@acme/skills": sourceops.TrustTrusted, "@acme/other": sourceops.TrustNone, "@me/scratch": sourceops.TrustNA,
	}
	for id, trustState := range want {
		if got[id] != trustState {
			t.Errorf("%s: trust %q, want %q", id, got[id], trustState)
		}
	}
	if order[0] != skill.DefaultRegistryID || order[1] != skill.OfficialID {
		t.Errorf("order = %v, want the built-ins first", order)
	}
}

func TestPlanAddPreviewsWithoutWritingAndRegistersUntrusted(t *testing.T) {
	t.Parallel()
	env, dir := newEnv(t)
	plan, apply, err := sourceops.PlanAdd(env, "@acme/skills", sources.Options{Roots: []string{"skills"}})
	if err != nil {
		t.Fatal(err)
	}
	if plan.Action != "add" || plan.Source.Trust != sourceops.TrustNone ||
		!strings.Contains(plan.Next, "rsk source trust @acme/skills") {
		t.Errorf("plan = %+v", plan)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "sources.json")); statErr == nil {
		t.Error("planning wrote sources.json")
	}
	if err = apply(); err != nil {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "trust.json")); statErr == nil {
		t.Error("adding a source granted trust")
	}
	if _, _, err = sourceops.PlanAdd(env, "@acme/skills", sources.Options{}); err == nil {
		t.Error("a source was added twice")
	}
	if _, _, err = sourceops.PlanAdd(env, "ext::sh -c id", sources.Options{ID: "@a/b"}); err == nil {
		t.Error("a hostile transport was accepted")
	}
}

func TestPlanRemoveDropsTheEntryAndItsGrant(t *testing.T) {
	t.Parallel()
	env, dir := newEnv(t)
	_, apply, _ := sourceops.PlanAdd(env, "@acme/skills", sources.Options{})
	if err := apply(); err != nil {
		t.Fatal(err)
	}
	grants := trust.Grant(trust.File{}, "@acme/skills", "git@github.com:acme/skills.git", time.Now())
	if err := trust.Save(filepath.Join(dir, "trust.json"), grants); err != nil {
		t.Fatal(err)
	}
	if _, _, err := sourceops.PlanRemove(env, "@nobody/x"); err == nil {
		t.Error("removing an undeclared source succeeded")
	}
	plan, apply, err := sourceops.PlanRemove(env, "@acme/skills")
	if err != nil || plan.Action != "remove" {
		t.Fatalf("plan = %+v, %v", plan, err)
	}
	if err = apply(); err != nil {
		t.Fatal(err)
	}
	f, _ := sources.Load(filepath.Join(dir, "sources.json"))
	g, _ := trust.Load(filepath.Join(dir, "trust.json"))
	if len(f.Sources) != 0 || len(g.Grants) != 0 {
		t.Errorf("sources = %+v, grants = %+v", f.Sources, g.Grants)
	}
}

func TestTestRefusesAnUntrustedSourceBeforeAnyContact(t *testing.T) {
	t.Parallel()
	env, _ := newEnv(t)
	contacted := false
	env.LsRemote = func(context.Context, git.RemoteURL, string) (string, error) { contacted = true; return commit, nil }
	_, apply, _ := sourceops.PlanAdd(env, "@acme/skills", sources.Options{})
	if err := apply(); err != nil {
		t.Fatal(err)
	}

	res, err := sourceops.Test(t.Context(), env, "@acme/skills")
	if !errors.Is(err, trust.ErrUntrusted) || contacted || res.Reachable || res.Error == "" {
		t.Errorf("res = %+v, err = %v, contacted = %v", res, err, contacted)
	}
	if _, err = sourceops.Test(t.Context(), env, "@nobody/x"); err == nil {
		t.Error("testing an undeclared source succeeded")
	}
}

func TestTestReportsAReachableTrustedSourceAndTheRewrittenURL(t *testing.T) {
	t.Parallel()
	env, dir := newEnv(t)
	_, apply, _ := sourceops.PlanAdd(env, "@acme/skills", sources.Options{})
	if err := apply(); err != nil {
		t.Fatal(err)
	}
	grants := trust.Grant(trust.File{}, "@acme/skills", "git@github.com:acme/skills.git", time.Now())
	if err := trust.Save(filepath.Join(dir, "trust.json"), grants); err != nil {
		t.Fatal(err)
	}
	env.EffectiveURL = func(context.Context, git.RemoteURL) (string, error) { return "git@alias:acme/skills.git", nil }

	res, err := sourceops.Test(t.Context(), env, "@acme/skills")
	if err != nil || !res.Reachable || res.Commit != commit || res.Via != "git@alias:acme/skills.git" {
		t.Errorf("res = %+v, err = %v", res, err)
	}
}

func writeSkill(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, name), 0o750); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: " + name + "\nversion: 1.0.0\ndescription: x\n---\n"
	if err := os.WriteFile(filepath.Join(dir, name, "SKILL.md"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestScanADirectoryNeedsNoRegistrationAndShowsTheClash(t *testing.T) {
	t.Parallel()
	env, _ := newEnv(t)
	root := t.TempDir()
	writeSkill(t, filepath.Join(root, "a"), "same")
	writeSkill(t, filepath.Join(root, "b"), "same")
	writeSkill(t, root, "fine")

	out, err := sourceops.Scan(t.Context(), env, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Skills) != 1 || out.Skills[0].Name != "fine" || len(out.Clashes) != 1 ||
		out.Source.Trust != sourceops.TrustNA {
		t.Errorf("scan = %+v", out)
	}
	if _, statErr := os.Stat(env.SourcesPath); statErr == nil {
		t.Error("scanning a directory registered it")
	}
}

func TestScanADeclaredGitSourceNeedsTrustAndAnUnregisteredOneIsRefused(t *testing.T) {
	t.Parallel()
	env, dir := newEnv(t)
	_, apply, _ := sourceops.PlanAdd(env, "@acme/skills", sources.Options{})
	if err := apply(); err != nil {
		t.Fatal(err)
	}
	if _, err := sourceops.Scan(t.Context(), env, "@acme/skills"); !errors.Is(err, trust.ErrUntrusted) {
		t.Errorf("an untrusted source was scanned: %v", err)
	}
	grants := trust.Grant(trust.File{}, "@acme/skills", "git@github.com:acme/skills.git", time.Now())
	if err := trust.Save(filepath.Join(dir, "trust.json"), grants); err != nil {
		t.Fatal(err)
	}
	out, err := sourceops.Scan(t.Context(), env, "@acme/skills")
	if err != nil || len(out.Skills) != 1 || out.Skills[0].Name != "alpha" {
		t.Errorf("scan = %+v, %v", out, err)
	}
	if _, err = sourceops.Scan(t.Context(), env, "@nobody/repo"); err == nil ||
		!strings.Contains(err.Error(), "rsk source add @nobody/repo") {
		t.Errorf("an unregistered source: err = %v, want the command to register it", err)
	}
}
