package tool

import "go.ralva.dev/ralvaskills/v3/internal/skill"

const (
	// CursorID identifies the Cursor CLI (agent) target.
	CursorID ID = "cursor"

	cursorHomeFolder = ".cursor"
	cursorMCPFile    = ".cursor/mcp.json"
)

// newCursorTool builds the adapter for Cursor CLI. Cursor needs an explicit
// "type": "stdio" on a local server, which the shared read-only entry does not
// carry, so this adapter supplies its own. The CLI has been reported not to read
// the project file in some releases; the editor does, so rsk still writes there.
func newCursorTool() *jsonAdapter {
	return &jsonAdapter{
		id:         CursorID,
		skillsDir:  func() string { return homePath(cursorHomeFolder, skill.SkillsFolderName) },
		projectDir: agentsSkillsDir,
		mcp: jsonMCPConfig{
			file:   cursorMCPFile,
			key:    mcpServersKey,
			server: cursorMCPServer,
		},
	}
}

// cursorMCPServer is the read-only entry plus the transport type Cursor
// requires on a local server.
func cursorMCPServer() map[string]any {
	server := readOnlyMCPServer()
	server[mcpTypeKey] = mcpStdioType
	return server
}
