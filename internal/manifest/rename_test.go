package manifest

import "testing"

func TestRenameMod(t *testing.T) {
	t.Parallel()

	aliases := map[string]string{
		"hexagonal-arch":      "hexagonal-architecture",
		"narration-architect": "narration-writer",
	}
	m := Mod{
		Skills: map[string]string{
			"hexagonal-arch":      "*",
			"go-architect":        "^1.2.0",
			"narration-architect": "*",
		},
		Pinned: []string{"hexagonal-arch", "docs-architect"},
	}

	got, renamed := RenameMod(m, aliases)

	if len(renamed) != 2 {
		t.Fatalf("renamed = %v, want two entries", renamed)
	}
	if _, stale := got.Skills["hexagonal-arch"]; stale {
		t.Error("old skill key survived the rename")
	}
	if got.Skills["hexagonal-architecture"] != "*" {
		t.Errorf("new key = %q, want *", got.Skills["hexagonal-architecture"])
	}
	if got.Skills["go-architect"] != "^1.2.0" {
		t.Errorf("unrelated constraint changed: %q", got.Skills["go-architect"])
	}
	if got.Pinned[0] != "hexagonal-architecture" {
		t.Errorf("pinned[0] = %q, want the new name", got.Pinned[0])
	}
	if got.Pinned[1] != "docs-architect" {
		t.Errorf("unrelated pin changed: %q", got.Pinned[1])
	}
}

func TestRenameModKeepsExistingTarget(t *testing.T) {
	t.Parallel()

	m := Mod{
		Skills: map[string]string{"old-name": "*", "new-name": "^2.0.0"},
		Pinned: []string{},
	}
	got, renamed := RenameMod(m, map[string]string{"old-name": "new-name"})

	if len(renamed) != 1 {
		t.Fatalf("renamed = %v, want one entry", renamed)
	}
	if got.Skills["new-name"] != "^2.0.0" {
		t.Errorf("existing constraint overridden: %q", got.Skills["new-name"])
	}
}

func TestRenameLock(t *testing.T) {
	t.Parallel()

	l := Lock{Skills: []LockEntry{
		{Name: "hexagonal-arch", Version: "1.2.1"},
		{Name: "go-architect", Version: "1.20.0"},
	}}

	got, renamed := RenameLock(l, map[string]string{"hexagonal-arch": "hexagonal-architecture"})

	if len(renamed) != 1 || renamed[0].New != "hexagonal-architecture" {
		t.Fatalf("renamed = %v, want one hexagonal entry", renamed)
	}
	if got.Skills[0].Name != "hexagonal-architecture" || got.Skills[0].Version != "1.2.1" {
		t.Errorf("entry = %+v, want the new name and the old version", got.Skills[0])
	}
	if got.Skills[1].Name != "go-architect" {
		t.Errorf("unrelated entry changed: %+v", got.Skills[1])
	}
}

func TestRenameLockDropsDuplicate(t *testing.T) {
	t.Parallel()

	l := Lock{Skills: []LockEntry{
		{Name: "old-name", Version: "1.0.0"},
		{Name: "new-name", Version: "2.0.0"},
	}}

	got, _ := RenameLock(l, map[string]string{"old-name": "new-name"})

	if len(got.Skills) != 1 {
		t.Fatalf("entries = %+v, want the duplicate dropped", got.Skills)
	}
	if got.Skills[0].Name != "new-name" || got.Skills[0].Version != "2.0.0" {
		t.Errorf("kept = %+v, want the existing new-name entry", got.Skills[0])
	}
}

func TestRenameModNoAliases(t *testing.T) {
	t.Parallel()

	m := Mod{Skills: map[string]string{"go-architect": "*"}, Pinned: []string{"go-architect"}}
	got, renamed := RenameMod(m, map[string]string{"other": "thing"})

	if len(renamed) != 0 {
		t.Errorf("renamed = %v, want none", renamed)
	}
	if got.Skills["go-architect"] != "*" {
		t.Errorf("unrelated manifest changed: %+v", got.Skills)
	}
}
