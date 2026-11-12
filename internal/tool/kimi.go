package tool

import (
	"os"
	"path/filepath"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

const (
	// KimiID identifies the Kimi Code CLI target.
	KimiID ID = "kimi"

	kimiHomeFolder = ".kimi-code"
	kimiHomeEnv    = "KIMI_CODE_HOME"
	kimiMCPFile    = ".kimi-code/mcp.json"
)

// newKimiTool builds the adapter for Kimi Code CLI. Kimi reads the shared
// .agents/skills directory in both scopes, so only the global directory and the
// MCP file are Kimi-specific. The archived Python kimi-cli is not supported.
func newKimiTool() *jsonAdapter {
	return &jsonAdapter{
		id:         KimiID,
		skillsDir:  kimiSkillsDir,
		projectDir: agentsSkillsDir,
		mcp: jsonMCPConfig{
			file:   kimiMCPFile,
			key:    mcpServersKey,
			server: readOnlyMCPServer,
		},
	}
}

// kimiSkillsDir returns $KIMI_CODE_HOME/skills when the variable is set, so an
// isolated data root keeps its own skills, and ~/.kimi-code/skills otherwise.
func kimiSkillsDir() string {
	if home := trimSlash(os.Getenv(kimiHomeEnv)); home != "" {
		return filepath.Join(home, skill.SkillsFolderName)
	}
	return homePath(kimiHomeFolder, skill.SkillsFolderName)
}
