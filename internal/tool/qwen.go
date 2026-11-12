package tool

import "go.ralva.dev/ralvaskills/v3/internal/skill"

const (
	// QwenID identifies the Qwen Code target.
	QwenID ID = "qwen"

	qwenHomeFolder = ".qwen"
	qwenMCPFile    = ".qwen/settings.json"
)

// newQwenTool builds the adapter for Qwen Code. It is a Gemini CLI fork, so the
// layout matches: mcpServers inside settings.json, and .agents/skills as the
// shared skills directory. Qwen renames a settings file it cannot parse to
// .corrupted, which is why rsk only ever re-encodes a file it just parsed.
func newQwenTool() *jsonAdapter {
	return &jsonAdapter{
		id:         QwenID,
		skillsDir:  func() string { return homePath(qwenHomeFolder, skill.SkillsFolderName) },
		projectDir: agentsSkillsDir,
		mcp: jsonMCPConfig{
			file:   qwenMCPFile,
			key:    mcpServersKey,
			server: readOnlyMCPServer,
		},
	}
}
