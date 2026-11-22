package main

import (
	"slices"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/mcp"
)

// mcpTwins maps each MCP tool to the CLI command that does the same thing, so
// an agent and a person can do the same set-up from either surface.
var mcpTwins = map[string][]string{
	"search_skills":     {"catalog"},
	"install_skills":    {"install"},
	"list_sources":      {"source", "list"},
	"test_source":       {"source", "test"},
	"scan_source":       {"source", "scan"},
	"resolve_reference": {"resolve"},
	"add_source":        {"source", "add"},
	"remove_source":     {"source", "remove"},
}

// agentOnly are the tools with no CLI twin. Adding one here is a decision: the
// design wants every operation on both surfaces.
var agentOnly = []string{"project_profile"}

// TestEveryMCPToolHasACLICounterpart fails when a tool is registered without
// one, and when a mutating tool's twin cannot preview with --dry-run.
//
//nolint:paralleltest // the shared cobra tree mutates on read
func TestEveryMCPToolHasACLICounterpart(t *testing.T) {
	for _, name := range mcp.ToolNames() {
		if slices.Contains(agentOnly, name) {
			continue
		}
		path, ok := mcpTwins[name]
		if !ok {
			t.Errorf("MCP tool %s has no CLI counterpart: add it to mcpTwins, or to agentOnly with a reason", name)
			continue
		}
		cmd, _, err := rootCmd.Find(path)
		if err != nil || cmd.Name() != path[len(path)-1] || cmd == rootCmd {
			t.Errorf("MCP tool %s names CLI command %v, which does not exist (%v)", name, path, err)
			continue
		}
		if (name == "add_source" || name == "remove_source") && cmd.Flags().Lookup("dry-run") == nil {
			t.Errorf("%s: the CLI twin of a mutating MCP tool must take --dry-run", name)
		}
	}
	for name := range mcpTwins {
		if !slices.Contains(mcp.ToolNames(), name) {
			t.Errorf("mcpTwins names %s, which is not an MCP tool", name)
		}
	}
}
