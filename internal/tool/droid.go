package tool

import "go.ralva.dev/ralvaskills/v3/internal/skill"

const (
	// DroidID identifies the Factory Droid CLI target.
	DroidID ID = "droid"

	droidHomeFolder = ".factory"
	droidMCPFile    = ".factory/mcp.json"
)

// newDroidTool builds the adapter for Factory Droid. Droid reads
// .agents/skills in both scopes and .factory/skills of its own; rsk uses the
// shared directory so one link serves every compatible client.
func newDroidTool() *jsonAdapter {
	return &jsonAdapter{
		id:         DroidID,
		skillsDir:  func() string { return homePath(droidHomeFolder, skill.SkillsFolderName) },
		projectDir: agentsSkillsDir,
		mcp: jsonMCPConfig{
			file:   droidMCPFile,
			key:    mcpServersKey,
			server: readOnlyMCPServer,
		},
	}
}
