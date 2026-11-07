package config

import "testing"

func TestEffectiveRegistryURL(t *testing.T) {
	cfg := Config{RegistryURL: "https://file.example"}
	t.Setenv(EnvRegistryURL, "")
	if got := cfg.EffectiveRegistryURL(); got != "https://file.example" {
		t.Fatalf("file value: got %q", got)
	}
	t.Setenv(EnvRegistryURL, " https://env.example/registry/ ")
	if got := cfg.EffectiveRegistryURL(); got != "https://env.example/registry" {
		t.Fatalf("env override: got %q", got)
	}
	if cfg.RegistryURL != "https://file.example" {
		t.Fatalf("the stored value must not change: %q", cfg.RegistryURL)
	}
}
