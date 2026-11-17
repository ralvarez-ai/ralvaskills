package permission

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

// seedRules writes opencode.json with exactly rules allowed, replacing any
// previous content.
func seedRules(t *testing.T, root string, rules ...tool.Rule) {
	t.Helper()
	entries := make([]map[string]any, 0, len(rules))
	for _, rule := range rules {
		entries = append(entries, map[string]any{
			"action":   rule.Action,
			"resource": rule.Resource,
			"effect":   "allow",
		})
	}
	body, err := json.Marshal(map[string]any{tool.OpenCodePermissionsKey: entries})
	if err != nil {
		t.Fatalf("marshal seed: %v", err)
	}
	if err = os.WriteFile(filepath.Join(root, tool.OpenCodeFileName), body, 0o600); err != nil {
		t.Fatalf("seed opencode.json: %v", err)
	}
}

// rulesFromJSON reads the allow rules of an opencode.json payload, in order.
func rulesFromJSON(t *testing.T, data []byte) []tool.Rule {
	t.Helper()
	cfg := make(map[string]any)
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse opencode.json: %v", err)
	}
	entries, _ := cfg[tool.OpenCodePermissionsKey].([]any)
	rules := make([]tool.Rule, 0, len(entries))
	for _, entry := range entries {
		obj, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		effect, _ := obj["effect"].(string)
		if effect != "allow" {
			continue
		}
		action, _ := obj["action"].(string)
		resource, _ := obj["resource"].(string)
		rules = append(rules, tool.Rule{Action: action, Resource: resource})
	}
	return rules
}
