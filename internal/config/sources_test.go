package config_test

import (
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

//nolint:paralleltest // t.Setenv cannot run in a parallel test
func TestSourceSpecsFollowTheMode(t *testing.T) {
	t.Setenv(config.EnvRegistryURL, "")
	base := config.Config{OfficialCache: "/cache/anthropic", RegistryURL: "https://skills.example"}

	t.Run("registry mode", func(t *testing.T) {
		specs := base.SourceSpecs()
		if len(specs) != 2 ||
			specs[0].Type != source.KindRegistry || specs[0].URL != "https://skills.example" ||
			specs[1].Type != source.KindOfficial {
			t.Fatalf("specs = %+v, want registry then official", specs)
		}
	})

	t.Run("local mode", func(t *testing.T) {
		cfg := base
		cfg.RepoPath = "/home/u/ralvaskills"
		specs := cfg.SourceSpecs()
		if len(specs) != 2 ||
			specs[0].Type != source.KindPath || specs[0].URL != "/home/u/ralvaskills" ||
			specs[0].ID != "@local/ralvaskills" || specs[1].Type != source.KindOfficial {
			t.Fatalf("specs = %+v, want path then official", specs)
		}
	})

	t.Run("the registry URL override is honoured", func(t *testing.T) {
		t.Setenv(config.EnvRegistryURL, "https://mirror.example/")
		if got := base.SourceSpecs()[0].URL; got != "https://mirror.example" {
			t.Errorf("registry URL = %q, want the RSK_REGISTRY_URL override", got)
		}
	})

	t.Run("a chain builds in both modes", func(t *testing.T) {
		for _, cfg := range []config.Config{base, {RepoPath: "/r", OfficialCache: "/o"}} {
			if _, err := cfg.Sources(); err != nil {
				t.Errorf("Sources: %v", err)
			}
		}
	})
}
