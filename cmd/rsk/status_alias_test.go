package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/status"
)

func writeStatusFile(t *testing.T, path, data string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestOfflineAliasesPrefersTheLocalCloneThenTheCache(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	cfg := config.Config{OfficialCache: filepath.Join(root, "official")}

	if got := offlineAliases(cfg); got != nil {
		t.Fatalf("no clone and no cache: got %v, want nil", got)
	}

	writeStatusFile(t, filepath.Join(cfg.RegistryCache(), "index.json"),
		`{"fetched_at":"2020-01-01T00:00:00Z","aliases":{"cached-old":"cached-new"},"skills":{"x":{}}}`)
	if got := offlineAliases(cfg); got["cached-old"] != "cached-new" {
		t.Fatalf("cache fallback: got %v", got)
	}

	cfg.RepoPath = repo
	writeStatusFile(t, filepath.Join(repo, "registry", "aliases.toml"), "[aliases]\nlocal-old = \"local-new\"\n")
	got := offlineAliases(cfg)
	if got["local-old"] != "local-new" || got["cached-old"] != "" {
		t.Fatalf("local clone must win over the cache: got %v", got)
	}
}

func TestStatusMarksARetiredName(t *testing.T) {
	t.Parallel()

	renamed := map[string]string{"hexagonal-arch": "hexagonal-architecture"}

	tags := statusRowTags(status.Linked{Name: "hexagonal-arch"}, nil, renamed)
	if !strings.Contains(tags, "[renamed → hexagonal-architecture]") {
		t.Errorf("tags = %q, want a renamed mark", tags)
	}
	if current := statusRowTags(status.Linked{Name: "go-architect"}, nil, renamed); current != "" {
		t.Errorf("a current name got tags %q", current)
	}

	sections := []status.Section{{Skills: []status.Linked{{Name: "hexagonal-arch"}, {Name: "go-architect"}}}}
	entries := statusSectionsToEntries(sections, nil, renamed)
	if got := entries[0].Skills[0].RenamedTo; got != "hexagonal-architecture" {
		t.Errorf("RenamedTo = %q", got)
	}
	if got := entries[0].Skills[1].RenamedTo; got != "" {
		t.Errorf("a current name has RenamedTo %q", got)
	}
}
