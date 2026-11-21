package manifest

import (
	"errors"
	"os"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/schema"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestRegistryOnlyFilesKeepWritingV1(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	m := Mod{
		Tools:  []tool.ID{tool.ClaudeID},
		Pinned: []string{"go-architect"},
		Skills: map[string]string{"go-architect": "*", "tdd": "1.2.0"},
	}
	if err := WriteMod(dir, m); err != nil {
		t.Fatal(err)
	}
	wantMod := "version = 1\ntools = [\"claude-code\"]\npinned = [\"go-architect\"]\n\n[skills]\n  go-architect = \"*\"\n  tdd = \"1.2.0\"\n"
	if got := readFile(t, ModPath(dir)); got != wantMod {
		t.Errorf("rsk.mod =\n%s\nwant\n%s", got, wantMod)
	}

	lock := Lock{Skills: []LockEntry{{Name: "tdd", Version: "1.2.0", Source: skill.SourceRegistry}}}
	if err := WriteLock(dir, lock); err != nil {
		t.Fatal(err)
	}
	wantLock := "version = 1\n\n[[skills]]\n  name = \"tdd\"\n  version = \"1.2.0\"\n  source = \"registry\"\n"
	if got := readFile(t, LockPath(dir)); got != wantLock {
		t.Errorf("rsk.lock =\n%s\nwant\n%s", got, wantLock)
	}
}

func TestOfficialSkillsMakeBothFilesV2WithQualifiedKeys(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	m := Mod{
		Tools:  []tool.ID{tool.ClaudeID},
		Pinned: []string{"xlsx", "tdd"},
		Skills: map[string]string{"xlsx": "*", "tdd": "*"},
		Origin: map[string]string{"xlsx": skill.OfficialID},
	}
	if err := WriteMod(dir, m); err != nil {
		t.Fatal(err)
	}
	text := readFile(t, ModPath(dir))
	for _, want := range []string{"version = 2", `"@anthropics/skills/xlsx"`, `tdd = "*"`} {
		if !strings.Contains(text, want) {
			t.Errorf("rsk.mod lacks %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, `"xlsx"`) {
		t.Errorf("rsk.mod still carries the bare key:\n%s", text)
	}

	got, err := ReadMod(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.Skills["xlsx"] != "*" || got.Origin["xlsx"] != skill.OfficialID || got.Origin["tdd"] != "" {
		t.Errorf("read back skills=%v origin=%v", got.Skills, got.Origin)
	}
	if len(got.Pinned) != 2 || got.Pinned[0] != "xlsx" {
		t.Errorf("pinned = %v, want bare names in memory", got.Pinned)
	}

	lock := Lock{Skills: []LockEntry{
		{Name: "xlsx", Source: skill.SourceOfficial},
		{Name: "tdd", Version: "1.0.0", Source: skill.SourceRegistry},
	}}
	if err = WriteLock(dir, lock); err != nil {
		t.Fatal(err)
	}
	text = readFile(t, LockPath(dir))
	for _, want := range []string{
		"version = 2", `name = "@anthropics/skills/xlsx"`, `source = "@anthropics/skills"`, `transport = "git"`,
		`name = "tdd"`, `source = "@ralvarez-ai/ralvaskills"`, `transport = "registry"`,
	} {
		if !strings.Contains(text, want) {
			t.Errorf("rsk.lock lacks %q:\n%s", want, text)
		}
	}
	back, err := ReadLock(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(back.Skills) != 2 || back.Skills[0].Name != "xlsx" || back.Skills[0].Source != skill.SourceOfficial ||
		back.Skills[1].Source != skill.SourceRegistry {
		t.Errorf("lock read back = %+v", back.Skills)
	}
}

func TestOriginOfAnAbsentSkillDoesNotForceV2(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	m := Mod{Skills: map[string]string{"tdd": "*"}, Origin: map[string]string{"xlsx": skill.OfficialID}}
	if err := WriteMod(dir, m); err != nil {
		t.Fatal(err)
	}
	if text := readFile(t, ModPath(dir)); !strings.HasPrefix(text, "version = 1\n") {
		t.Errorf("an uninstalled skill's origin forced v2:\n%s", text)
	}
}

func TestV1FilesStillReadAndUpgradeOnWrite(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	oldMod := "version = 1\ntools = [\"claude-code\"]\npinned = [\"xlsx\"]\n\n[skills]\n  xlsx = \"*\"\n"
	oldLock := "version = 1\n\n[[skills]]\n  name = \"xlsx\"\n  version = \"\"\n  source = \"official\"\n"
	if err := os.WriteFile(ModPath(dir), []byte(oldMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(LockPath(dir), []byte(oldLock), 0o644); err != nil {
		t.Fatal(err)
	}

	m, err := ReadMod(dir)
	if err != nil || m.Skills["xlsx"] != "*" || len(m.Origin) != 0 {
		t.Fatalf("ReadMod = %+v, %v", m, err)
	}
	l, err := ReadLock(dir)
	if err != nil || len(l.Skills) != 1 || l.Skills[0].Name != "xlsx" || l.Skills[0].Source != skill.SourceOfficial {
		t.Fatalf("ReadLock = %+v, %v", l, err)
	}

	if err = WriteLock(dir, l); err != nil {
		t.Fatal(err)
	}
	if text := readFile(t, LockPath(dir)); !strings.Contains(text, `name = "@anthropics/skills/xlsx"`) {
		t.Errorf("lock was not upgraded:\n%s", text)
	}
	if err = WriteMod(dir, m); err != nil {
		t.Fatal(err)
	}
	if text := readFile(t, ModPath(dir)); !strings.HasPrefix(text, "version = 1\n") {
		t.Errorf("rewriting an unchanged v1 mod must keep it v1:\n%s", text)
	}
}

func TestNewerFilesAreRefusedWithTheUpgradeMessage(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.WriteFile(ModPath(dir), []byte("version = 3\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(LockPath(dir), []byte("version = 3\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadMod(dir); !errors.Is(err, schema.ErrUnsupported) {
		t.Errorf("ReadMod err = %v, want ErrUnsupported", err)
	}
	if _, err := ReadLock(dir); !errors.Is(err, schema.ErrUnsupported) ||
		!strings.Contains(err.Error(), "upgrade rsk") {
		t.Errorf("ReadLock err = %v", err)
	}
}

func TestReadRejectsWhatItCannotTrust(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, mod, lock string
	}{
		{
			name: "two skills with one name in a scope",
			mod:  "version = 2\n[skills]\n  xlsx = \"*\"\n  \"@anthropics/skills/xlsx\" = \"*\"\n",
		},
		{
			name: "key and source disagree",
			lock: "version = 2\n[[skills]]\n  name = \"@anthropics/skills/xlsx\"\n  source = \"@ralvarez-ai/ralvaskills\"\n",
		},
		{
			name: "a source this rsk cannot read",
			lock: "version = 2\n[[skills]]\n  name = \"x\"\n  source = \"bogus\"\n",
		},
		{
			name: "a malformed key",
			mod:  "version = 2\n[skills]\n  \"@acme/x\" = \"*\"\n",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			var err error
			if tc.mod != "" {
				_ = os.WriteFile(ModPath(dir), []byte(tc.mod), 0o644)
				_, err = ReadMod(dir)
			} else {
				_ = os.WriteFile(LockPath(dir), []byte(tc.lock), 0o644)
				_, err = ReadLock(dir)
			}
			if err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func gitRef() skill.GitRef {
	return skill.GitRef{
		URL: "git@github.com:acme/skills.git", Rev: "main", Path: "skills/alpha",
		Commit: "1111111111111111111111111111111111111111", Tree: "2222222222222222222222222222222222222222",
	}
}

func gitEntry() LockEntry {
	ref := gitRef()
	return LockEntry{Name: "alpha", Version: "1.0.0", Source: skill.DeclaredSource("@acme/skills"), Git: &ref}
}

func TestGitEntriesRoundTripThroughAV2Lock(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	want := gitEntry()
	if err := WriteLock(
		dir,
		Lock{Skills: []LockEntry{want, {Name: "tdd", Version: "1.0.0", Source: skill.SourceRegistry}}},
	); err != nil {
		t.Fatal(err)
	}
	text := readFile(t, LockPath(dir))
	for _, line := range []string{
		"version = 2", `name = "@acme/skills/alpha"`, `source = "@acme/skills"`, `transport = "git"`,
		`url = "git@github.com:acme/skills.git"`, `rev = "main"`, `path = "skills/alpha"`,
		`commit = "1111111111111111111111111111111111111111"`, `tree = "2222222222222222222222222222222222222222"`,
	} {
		if !strings.Contains(text, line) {
			t.Errorf("rsk.lock lacks %q:\n%s", line, text)
		}
	}
	got, err := ReadLock(dir)
	if err != nil || len(got.Skills) != 2 {
		t.Fatalf("ReadLock = %+v, %v", got, err)
	}
	if g := got.Skills[0]; g.Name != want.Name || g.Source != want.Source || g.Version != want.Version ||
		g.Git == nil || *g.Git != *want.Git {
		t.Errorf("entry = %+v, want %+v", g, want)
	}
	if got.Skills[1].Git != nil {
		t.Errorf("a registry entry grew git fields: %+v", got.Skills[1])
	}
}

func TestReadLockRefusesAMalformedGitEntry(t *testing.T) {
	t.Parallel()
	tests := map[string]func(*skill.GitRef){
		"no commit":        func(e *skill.GitRef) { e.Commit = "" },
		"short commit":     func(e *skill.GitRef) { e.Commit = "abc123" },
		"branch as commit": func(e *skill.GitRef) { e.Commit = "main" },
		"no tree":          func(e *skill.GitRef) { e.Tree = "" },
		"traversing path":  func(e *skill.GitRef) { e.Path = "../alpha" },
		"path vs name":     func(e *skill.GitRef) { e.Path = "skills/beta" },
		"hostile url":      func(e *skill.GitRef) { e.URL = "ext::sh -c id" },
		"credential url": func(e *skill.GitRef) {
			e.URL = "https://u:ghp_abcdefghijklmnopqrstuvwxyz0123456789@github.com/a/b.git"
		},
		"option as rev": func(e *skill.GitRef) { e.Rev = "--upload-pack=x" },
		"no url":        func(e *skill.GitRef) { e.URL = "" },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			e := gitRef()
			mutate(&e)
			// Encode through WriteLock's own shape but bypass its checks by writing
			// the text directly: a lock file is untrusted input.
			text := "version = 2\n[[skills]]\n  name = \"@acme/skills/alpha\"\n  version = \"1.0.0\"\n  source = \"@acme/skills\"\n" +
				"  transport = \"git\"\n  url = " + quote(
				e.URL,
			) + "\n  rev = " + quote(
				e.Rev,
			) + "\n  commit = " + quote(
				e.Commit,
			) +
				"\n  tree = " + quote(
				e.Tree,
			) + "\n  path = " + quote(
				e.Path,
			) + "\n"
			if err := os.WriteFile(LockPath(dir), []byte(text), 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := ReadLock(dir); err == nil {
				t.Error("a malformed git entry was accepted")
			}
		})
	}
}

func quote(s string) string {
	return "\"" + strings.ReplaceAll(s, "\"", "\\\"") + "\""
}

func TestReadLockRefusesGitFieldsOnOtherSourcesAndAMissingTransport(t *testing.T) {
	t.Parallel()
	for name, text := range map[string]string{
		"git fields on a registry entry":             "version = 2\n[[skills]]\n  name = \"tdd\"\n  version = \"1\"\n  source = \"@ralvarez-ai/ralvaskills\"\n  transport = \"registry\"\n  commit = \"" + strings.Repeat("a", 40) + "\"\n",
		"a declared source without transport":        "version = 2\n[[skills]]\n  name = \"@acme/skills/alpha\"\n  version = \"1\"\n  source = \"@acme/skills\"\n",
		"a declared source over the wrong transport": "version = 2\n[[skills]]\n  name = \"@acme/skills/alpha\"\n  version = \"1\"\n  source = \"@acme/skills\"\n  transport = \"registry\"\n",
	} {
		dir := t.TempDir()
		if err := os.WriteFile(LockPath(dir), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := ReadLock(dir); err == nil {
			t.Errorf("%s was accepted", name)
		}
	}
}

func TestModQualifiesAGitSkillsKey(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	m := Mod{Skills: map[string]string{"alpha": "1.2.0"}, Pinned: []string{"alpha"}}
	m.SetOrigin("alpha", skill.DeclaredSource("@acme/skills"))
	if err := WriteMod(dir, m); err != nil {
		t.Fatal(err)
	}
	text := readFile(t, ModPath(dir))
	if !strings.Contains(text, "version = 2") || !strings.Contains(text, `"@acme/skills/alpha" = "1.2.0"`) {
		t.Errorf("rsk.mod =\n%s", text)
	}
	got, err := ReadMod(dir)
	if err != nil || got.Origin["alpha"] != "@acme/skills" || got.Skills["alpha"] != "1.2.0" {
		t.Errorf("read back %+v, %v", got, err)
	}
}
