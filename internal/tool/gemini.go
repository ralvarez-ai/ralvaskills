package tool

import "go.ralva.dev/ralvaskills/v3/internal/skill"

const (
	// GeminiID identifies the Gemini CLI target.
	GeminiID ID = "gemini"

	geminiHomeFolder = ".gemini"
	geminiMCPFile    = ".gemini/settings.json"
)

// newGeminiTool builds the adapter for Gemini CLI. MCP servers share
// settings.json with hooks, permissions and model settings, so registration
// merges the mcpServers key and leaves every other key byte-identical after
// re-encoding. Gemini's personal use moved to Antigravity; this target remains
// for Gemini Code Assist enterprise and API-key users.
func newGeminiTool() *jsonAdapter {
	return &jsonAdapter{
		id:         GeminiID,
		skillsDir:  func() string { return homePath(geminiHomeFolder, skill.SkillsFolderName) },
		projectDir: agentsSkillsDir,
		mcp: jsonMCPConfig{
			file:   geminiMCPFile,
			key:    mcpServersKey,
			server: readOnlyMCPServer,
		},
	}
}
