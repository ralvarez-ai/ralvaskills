package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/cmdx"
	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

func TestPlanSyncUsesLockedVersionAndConfiguredDirectories(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	cache := t.TempDir()
	cfg := config.Config{OfficialCache: filepath.Join(cache, "anthropic")}
	skillDir := filepath.Join(cfg.RegistryCache(), "example", "1.0.0")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(skillDir, "SKILL.md"),
		[]byte("---\nname: example\nversion: 1.0.0\ndescription: Example\n---\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}
	m := manifest.Mod{Tools: []tool.ID{tool.ClaudeID, tool.OpenCodeID}, Skills: map[string]string{"example": "^1.0.0"}}
	lock := manifest.Lock{
		Skills: []manifest.LockEntry{{Name: "example", Version: "1.0.0", Source: skill.SourceRegistry}},
	}
	links, err := planSyncLinks(t.Context(), project, m, lock, cfg)
	if err != nil || len(links) != 2 {
		t.Fatalf("links=%+v err=%v", links, err)
	}
	for _, link := range links {
		if link.skill.Path != skillDir {
			t.Fatalf("source = %s, want %s", link.skill.Path, skillDir)
		}
		if err = skill.Link(link.skill, link.target); err != nil {
			t.Fatal(err)
		}
	}
	links, err = planSyncLinks(t.Context(), project, m, lock, cfg)
	if err != nil || len(links) != 0 {
		t.Fatalf("second plan: links=%+v err=%v", links, err)
	}
	if m.Skills["example"] != "^1.0.0" || lock.Skills[0].Version != "1.0.0" {
		t.Fatal("manifest or lock was changed")
	}
}

func TestPlanSyncRefusesUserOwnedEntry(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	cfg := config.Config{OfficialCache: filepath.Join(t.TempDir(), "anthropic")}
	skillDir := filepath.Join(cfg.RegistryCache(), "example", "1.0.0")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(skillDir, "SKILL.md"),
		[]byte("---\nversion: 1.0.0\n---\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(project, ".opencode", "skills", "example")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	m := manifest.Mod{Tools: []tool.ID{tool.OpenCodeID}, Skills: map[string]string{"example": "*"}}
	lock := manifest.Lock{
		Skills: []manifest.LockEntry{{Name: "example", Version: "1.0.0", Source: skill.SourceRegistry}},
	}
	if _, err := planSyncLinks(t.Context(), project, m, lock, cfg); err == nil {
		t.Fatal("expected conflict for existing directory")
	}
}

// TestPlanSyncRefusesStaleLink keeps the lock authoritative: an existing rsk
// link that points at another cached version must be reported, never silently
// re-pointed, so a half-synced project cannot hide a stale install.
func TestPlanSyncRefusesStaleLink(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	cfg := config.Config{OfficialCache: filepath.Join(t.TempDir(), "anthropic")}
	for _, version := range []string{"0.9.0", "1.0.0"} {
		dir := filepath.Join(cfg.RegistryCache(), "example", version)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		body := "---\nname: example\nversion: " + version + "\ndescription: Example\n---\n"
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	old := skill.Skill{
		Name:    "example",
		Version: "0.9.0",
		Source:  skill.SourceRegistry,
		Path:    filepath.Join(cfg.RegistryCache(), "example", "0.9.0"),
	}
	if err := skill.Link(old, filepath.Join(project, ".claude", "skills")); err != nil {
		t.Fatal(err)
	}
	m := manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"example": "*"}}
	lock := manifest.Lock{
		Skills: []manifest.LockEntry{{Name: "example", Version: "1.0.0", Source: skill.SourceRegistry}},
	}
	_, err := planSyncLinks(t.Context(), project, m, lock, cfg)
	if err == nil || !strings.Contains(err.Error(), "points to") {
		t.Fatalf("expected a conflicting-link error, got %v", err)
	}
}

// TestPlanSyncSupportsLocalSource covers the local-clone configuration, where
// skills resolve from repo_path instead of the registry cache.
func TestPlanSyncSupportsLocalSource(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	repo := t.TempDir()
	skillDir := filepath.Join(repo, "skills", "example")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: example\nversion: 1.0.0\ndescription: Example\n---\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Config{RepoPath: repo, OfficialCache: filepath.Join(t.TempDir(), "anthropic")}
	m := manifest.Mod{Tools: []tool.ID{tool.CodexID}, Skills: map[string]string{"example": "*"}}
	lock := manifest.Lock{
		Skills: []manifest.LockEntry{{Name: "example", Version: "1.0.0", Source: skill.SourceLocal}},
	}
	links, err := planSyncLinks(t.Context(), project, m, lock, cfg)
	if err != nil || len(links) != 1 {
		t.Fatalf("links=%+v err=%v", links, err)
	}
	if links[0].target != filepath.Join(project, ".agents", "skills") {
		t.Fatalf("target = %s, want the Codex project directory", links[0].target)
	}
	if links[0].skill.Path != skillDir {
		t.Fatalf("source = %s, want %s", links[0].skill.Path, skillDir)
	}
}

// TestPlanSyncSupportsVersionlessOfficialSkills covers the official-source
// lock entry: those skills carry no version in their frontmatter, so the lock
// records an empty one and sync must not reject it.
func TestPlanSyncSupportsVersionlessOfficialSkills(t *testing.T) {
	t.Parallel()
	project := t.TempDir()
	cfg := config.Config{OfficialCache: filepath.Join(t.TempDir(), "anthropic")}
	skillDir := filepath.Join(cfg.OfficialCache, "skills", "frontend-design")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: frontend-design\ndescription: Frontend design\n---\n\nBody.\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m := manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"frontend-design": "*"}}
	lock := manifest.Lock{
		Skills: []manifest.LockEntry{{Name: "frontend-design", Source: skill.SourceOfficial}},
	}
	links, err := planSyncLinks(t.Context(), project, m, lock, cfg)
	if err != nil || len(links) != 1 {
		t.Fatalf("links=%+v err=%v", links, err)
	}
	if links[0].skill.Path != skillDir {
		t.Fatalf("source = %s, want %s", links[0].skill.Path, skillDir)
	}
}

func TestPlanSyncRejectsBrokenLocks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		mod     manifest.Mod
		lock    manifest.Lock
		wantErr string
	}{
		{
			name:    "missing lock entry",
			mod:     manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"example": "*"}},
			wantErr: "no lock entry",
		},
		{
			name: "duplicate lock entry",
			mod:  manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"example": "*"}},
			lock: manifest.Lock{Skills: []manifest.LockEntry{
				{Name: "example", Version: "1.0.0", Source: skill.SourceRegistry},
				{Name: "example", Version: "1.0.0", Source: skill.SourceRegistry},
			}},
			wantErr: "duplicate lock entry",
		},
		{
			name: "invalid name",
			mod:  manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"Example_Bad": "*"}},
			lock: manifest.Lock{Skills: []manifest.LockEntry{
				{Name: "Example_Bad", Version: "1.0.0", Source: skill.SourceRegistry},
			}},
			wantErr: "not kebab-case",
		},
		{
			name: "invalid version",
			mod:  manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"example": "*"}},
			lock: manifest.Lock{Skills: []manifest.LockEntry{
				{Name: "example", Version: "../1.0.0", Source: skill.SourceRegistry},
			}},
			wantErr: "invalid locked version",
		},
		{
			name: "uncached registry version",
			mod:  manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"example": "*"}},
			lock: manifest.Lock{Skills: []manifest.LockEntry{
				{Name: "example", Version: "1.0.0", Source: skill.SourceRegistry},
			}},
			wantErr: "is unavailable",
		},
		{
			name: "missing official cache",
			mod:  manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"example": "*"}},
			lock: manifest.Lock{Skills: []manifest.LockEntry{
				{Name: "example", Version: "1.0.0", Source: skill.SourceOfficial},
			}},
			wantErr: "official cache unavailable",
		},
		{
			name: "unsupported source",
			mod:  manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"example": "*"}},
			lock: manifest.Lock{Skills: []manifest.LockEntry{
				{Name: "example", Version: "1.0.0", Source: skill.Source("pypi")},
			}},
			wantErr: "unsupported source",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cfg := config.Config{OfficialCache: filepath.Join(t.TempDir(), "anthropic")}
			_, err := planSyncLinks(t.Context(), t.TempDir(), tc.mod, tc.lock, cfg)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error = %v, want %q", err, tc.wantErr)
			}
		})
	}
}

// TestRunSyncRestoresLegacyProject keeps a project written by an earlier
// release working: the legacy manifests are accepted unchanged, missing links
// are restored from the lock, and hand-written client files stay intact.
//
//nolint:paralleltest // t.Chdir and t.Setenv change process-wide state
func TestRunSyncRestoresLegacyProject(t *testing.T) {
	project, skillDir := seedLegacySyncProject(t, `["claude-code", "opencode"]`)

	if err := os.WriteFile(
		filepath.Join(project, "AGENTS.md"), []byte("User guide\n"), 0o644,
	); err != nil {
		t.Fatal(err)
	}
	if out, err := runSyncInTest(t, false); err != nil {
		t.Fatalf("runSync: %v\n%s", err, out)
	}

	for _, target := range []string{".claude/skills/example", ".opencode/skills/example"} {
		assertLinkTarget(t, filepath.Join(project, target), skillDir)
	}
	if got := readTestFile(t, filepath.Join(project, ".rsk", "SKILLS.md")); !strings.Contains(got, "example") {
		t.Fatalf(".rsk/SKILLS.md does not mention the pinned skill: %s", got)
	}
	claude := readTestFile(t, filepath.Join(project, "CLAUDE.md"))
	if !strings.Contains(claude, "Keep me.") || !strings.Contains(claude, "@.rsk/CLAUDE.md") {
		t.Fatalf("CLAUDE.md lost user content or the rsk import: %s", claude)
	}
	if got := readTestFile(t, filepath.Join(project, "AGENTS.md")); got != "User guide\n" {
		t.Fatalf("AGENTS.md was rewritten: %q", got)
	}
	opencode := readTestFile(t, filepath.Join(project, "opencode.json"))
	if !strings.Contains(opencode, "docs/keep.md") || !strings.Contains(opencode, ".rsk/SKILLS.md") {
		t.Fatalf("opencode.json lost instructions: %s", opencode)
	}
	if !strings.Contains(opencode, `"other"`) {
		t.Fatalf("opencode.json lost the unrelated MCP server: %s", opencode)
	}

	mod := readTestFile(t, filepath.Join(project, ".rsk", "rsk.mod"))
	lock := readTestFile(t, filepath.Join(project, ".rsk", "rsk.lock"))
	skills := readTestFile(t, filepath.Join(project, ".rsk", "SKILLS.md"))
	out, err := runSyncInTest(t, false)
	if err != nil {
		t.Fatalf("second runSync: %v\n%s", err, out)
	}
	if strings.Contains(out, "link ") {
		t.Fatalf("second run re-created links: %s", out)
	}
	for path, want := range map[string]string{
		filepath.Join(project, ".rsk", "rsk.mod"):   mod,
		filepath.Join(project, ".rsk", "rsk.lock"):  lock,
		filepath.Join(project, ".rsk", "SKILLS.md"): skills,
		filepath.Join(project, "opencode.json"):     opencode,
		filepath.Join(project, ".rsk", "CLAUDE.md"): readTestFile(t, filepath.Join(project, ".rsk", "CLAUDE.md")),
	} {
		if got := readTestFile(t, path); got != want {
			t.Fatalf("%s changed on the second sync", path)
		}
	}
}

// TestRunSyncUpgradesLegacyMCPRegistration covers the MCP half of a sync for a
// project that still has the pre-read-only OpenCode entry.
//
//nolint:paralleltest // t.Chdir and t.Setenv change process-wide state
func TestRunSyncUpgradesLegacyMCPRegistration(t *testing.T) {
	project, _ := seedLegacySyncProject(t, `["opencode"]`)
	path := filepath.Join(project, "opencode.json")
	writeTestFile(t, path, `{"instructions":["docs/keep.md"],"mcp":{`+
		`"rsk":{"type":"local","command":["rsk","mcp"],"enabled":true},`+
		`"other":{"type":"local","command":["other"]}}}`)
	installFakeRsk(t, "#!/bin/sh\necho ' --read-only'\n")

	out, err := runSyncInTest(t, true)
	if err != nil {
		t.Fatalf("runSync --mcp: %v\n%s", err, out)
	}
	cfg := readTestFile(t, path)
	if !strings.Contains(cfg, `"--read-only"`) {
		t.Fatalf("legacy MCP entry was not upgraded: %s", cfg)
	}
	if !strings.Contains(cfg, "docs/keep.md") || !strings.Contains(cfg, `"other"`) {
		t.Fatalf("MCP sync dropped existing settings: %s", cfg)
	}
}

// TestRunSyncRefusesMCPWithoutReadOnlySupport proves the guard runs before any
// file is written: an old rsk on PATH must leave both links and MCP config alone.
//
//nolint:paralleltest // t.Chdir and t.Setenv change process-wide state
func TestRunSyncRefusesMCPWithoutReadOnlySupport(t *testing.T) {
	project, _ := seedLegacySyncProject(t, `["opencode"]`)
	path := filepath.Join(project, "opencode.json")
	original := `{"mcp":{"rsk":{"type":"local","command":["rsk","mcp"],"enabled":true}}}`
	writeTestFile(t, path, original)
	installFakeRsk(t, "#!/bin/sh\necho 'Usage: rsk mcp'\n")

	out, err := runSyncInTest(t, true)
	if err == nil || !strings.Contains(err.Error(), "read-only") {
		t.Fatalf("error = %v, want the read-only guard\n%s", err, out)
	}
	if got := readTestFile(t, path); got != original {
		t.Fatalf("MCP config changed despite the guard: %s", got)
	}
	if _, statErr := os.Lstat(filepath.Join(project, ".opencode", "skills", "example")); !os.IsNotExist(statErr) {
		t.Fatalf("sync wrote links before the MCP guard failed: %v", statErr)
	}
}

func TestCheckSyncMCPExecutable(t *testing.T) {
	tests := []struct {
		name    string
		script  string
		ids     []tool.ID
		wantErr bool
	}{
		{name: "new binary for codex", script: "#!/bin/sh\necho ' --read-only'\n", ids: []tool.ID{tool.CodexID}},
		{name: "new binary for opencode", script: "#!/bin/sh\necho ' --read-only'\n", ids: []tool.ID{tool.OpenCodeID}},
		{
			name:    "old binary for opencode",
			script:  "#!/bin/sh\necho 'Usage: rsk mcp'\n",
			ids:     []tool.ID{tool.OpenCodeID},
			wantErr: true,
		},
		{name: "old binary for claude", script: "#!/bin/sh\necho 'Usage: rsk mcp'\n", ids: []tool.ID{tool.ClaudeID}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "rsk")
			writeTestFile(t, path, tc.script)
			if err := os.Chmod(path, 0o755); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PATH", dir)
			err := checkSyncMCPExecutable(t.Context(), tc.ids)
			if tc.wantErr != (err != nil) {
				t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
	t.Run("missing binary", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		if err := checkSyncMCPExecutable(t.Context(), []tool.ID{tool.ClaudeID}); err == nil {
			t.Fatal("expected an error when rsk is not on PATH")
		}
	})
}

func seedLegacySyncProject(t *testing.T, tools string) (project, skillDir string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv(config.EnvConfigHome, home)
	cfg := config.Config{
		RegistryURL:        "https://example.invalid",
		GlobalTargets:      map[string]string{"claude-code": filepath.Join(home, "global-claude")},
		DefaultTargetScope: "claude-code",
		OfficialCache:      filepath.Join(home, "official"),
		VersionsCache:      filepath.Join(home, "versions"),
	}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	skillDir = filepath.Join(cfg.RegistryCache(), "example", "1.0.0")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: example\nversion: 1.0.0\ndescription: Example skill\n---\n\n# Example\n\nBody.\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	project = t.TempDir()
	rskDir := filepath.Join(project, ".rsk")
	if err := os.MkdirAll(rskDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Deliberately no version key: this is the pre-versioning format.
	mod := "tools = " + tools + "\npinned = [\"example\"]\n\n[skills]\nexample = \"^1.0.0\"\n"
	lock := "[[skills]]\nname = \"example\"\nversion = \"1.0.0\"\nsource = \"registry\"\n"
	writeTestFile(t, filepath.Join(rskDir, "rsk.mod"), mod)
	writeTestFile(t, filepath.Join(rskDir, "rsk.lock"), lock)
	writeTestFile(t, filepath.Join(project, "CLAUDE.md"), "# Project\n\nKeep me.\n")
	writeTestFile(t, filepath.Join(project, "opencode.json"),
		`{"instructions":["docs/keep.md"],"mcp":{"servers":{"other":{"type":"local","command":["other"]}}}}`)
	t.Chdir(project)
	return project, skillDir
}

func runSyncInTest(t *testing.T, mcp bool) (string, error) {
	t.Helper()
	flags := syncCmd.Flags()
	if err := flags.Set(cmdx.FlagYes, "true"); err != nil {
		t.Fatal(err)
	}
	if err := flags.Set(flagMCP, strconv.FormatBool(mcp)); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		for _, name := range []string{cmdx.FlagYes, flagMCP} {
			flag := flags.Lookup(name)
			_ = flag.Value.Set(flag.DefValue)
			flag.Changed = false
		}
	})
	var buf bytes.Buffer
	syncCmd.SetOut(&buf)
	syncCmd.SetContext(t.Context())
	t.Cleanup(func() {
		syncCmd.SetOut(nil)
		syncCmd.SetContext(context.Background())
	})
	return buf.String(), runSync(syncCmd, nil)
}

func installFakeRsk(t *testing.T, script string) {
	t.Helper()
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "rsk"), script)
	if err := os.Chmod(filepath.Join(dir, "rsk"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func assertLinkTarget(t *testing.T, link, want string) {
	t.Helper()
	got, err := filepath.EvalSymlinks(link)
	if err != nil {
		t.Fatalf("resolve %s: %v", link, err)
	}
	expected, err := filepath.EvalSymlinks(want)
	if err != nil {
		t.Fatal(err)
	}
	if got != expected {
		t.Fatalf("%s -> %s, want %s", link, got, expected)
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
