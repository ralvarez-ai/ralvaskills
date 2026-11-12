package tool

import (
	"go.ralva.dev/ralvaskills/v3/internal"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

const (
	// AntigravityID identifies the Google Antigravity target shared by the
	// Antigravity App, IDE and CLI.
	AntigravityID ID = "antigravity"

	// AntigravityMCPFileName is the MCP config Antigravity reads inside the
	// shared .agents/ workspace folder.
	AntigravityMCPFileName = "mcp_config.json"

	// antigravityHomeFolder is home-relative on purpose: Antigravity ignores
	// XDG_CONFIG_HOME, so this must not go through internal.ConfigHome().
	antigravityHomeFolder = ".gemini"
	antigravityConfigDir  = "config"
)

// newAntigravityTool builds the adapter. The global skills directory is shared
// by the App, IDE and CLI; the project directory is the .agents/skills
// convention, which Codex and the other Agent Skills clients read too.
func newAntigravityTool() *jsonAdapter {
	return &jsonAdapter{
		id:         AntigravityID,
		skillsDir:  func() string { return homePath(antigravityHomeFolder, antigravityConfigDir, skill.SkillsFolderName) },
		projectDir: agentsSkillsDir,
		mcp: jsonMCPConfig{
			file:   internal.AgentsFolderName + "/" + AntigravityMCPFileName,
			key:    mcpServersKey,
			server: readOnlyMCPServer,
		},
	}
}
