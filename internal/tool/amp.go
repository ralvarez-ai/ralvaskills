package tool

import "go.ralva.dev/ralvaskills/v3/internal/skill"

const (
	// AmpID identifies the Sourcegraph Amp CLI target.
	AmpID ID = "amp"

	// ampHomeFolder is where Amp keeps settings on every platform; it does not
	// follow XDG_CONFIG_HOME.
	ampHomeFolder = ".config/amp"
	ampMCPFile    = ".amp/settings.json"
	ampMCPKey     = "amp.mcpServers"

	// ampSettingsJSONC is a settings file Amp also reads. Registration refuses it
	// rather than creating a second file whose precedence is unclear.
	ampSettingsJSONC = ".amp/settings.jsonc"

	// ampGlobalSkillsFolder is the user-level skills root Amp prefers.
	ampGlobalSkillsFolder = ".config/agents"
)

// newAmpTool builds the adapter for Amp. Two details differ from the other JSON
// clients: the MCP map sits under the literal dotted key "amp.mcpServers", and
// the user-level skills root is ~/.config/agents/skills rather than
// ~/.config/amp/skills.
func newAmpTool() *jsonAdapter {
	return &jsonAdapter{
		id:         AmpID,
		skillsDir:  func() string { return homePath(ampGlobalSkillsFolder, skill.SkillsFolderName) },
		projectDir: agentsSkillsDir,
		mcp: jsonMCPConfig{
			file:      ampMCPFile,
			key:       ampMCPKey,
			server:    readOnlyMCPServer,
			alternate: []string{ampSettingsJSONC},
		},
	}
}
