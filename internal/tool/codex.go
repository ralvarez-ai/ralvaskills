package tool

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"

	"go.ralva.dev/ralvaskills/v3/internal"
	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/fsx"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

// CodexID identifies the Codex skill target.
const CodexID ID = "codex"

type codexTool struct{}

func (*codexTool) ID() ID {
	return CodexID
}

func (*codexTool) SkillsDir() string {
	return homePath(internal.AgentsFolderName, skill.SkillsFolderName)
}

func (*codexTool) ProjectSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, internal.AgentsFolderName, skill.SkillsFolderName)
}

// Codex discovers installed skills on demand. It has no automatic import of
// the generated digest, and rsk does not edit the human-owned AGENTS.md.
func (*codexTool) SyncPinned(string, []string) error {
	return nil
}

func (*codexTool) RemovePinned(string) error {
	return nil
}

const (
	codexFolder   = ".codex"
	codexConfig   = "config.toml"
	codexMCPStart = "# rsk-mcp:start"
	codexMCPEnd   = "# rsk-mcp:end"
	codexMCPBlock = codexMCPStart + "\n[mcp_servers.rsk]\ncommand = \"rsk\"\nargs = [\"mcp\", \"--read-only\"]\n" + codexMCPEnd + "\n"
)

func (*codexTool) RegisterMCP(projectDir string) (string, error) {
	dir := filepath.Join(projectDir, codexFolder)
	path := filepath.Join(dir, codexConfig)
	data, err := readCodexConfig(path)
	if err != nil {
		return "", err
	}
	content := string(data)
	var cfg map[string]any
	if _, err = toml.Decode(content, &cfg); err != nil {
		return "", fmt.Errorf("parse %s: %w", path, err)
	}
	if strings.Contains(content, codexMCPStart) {
		if strings.Count(content, codexMCPBlock) != 1 || strings.Count(content, codexMCPStart) != 1 ||
			strings.Count(content, codexMCPEnd) != 1 {
			return "", fmt.Errorf("%s: modified rsk MCP block; resolve manually", path)
		}
		return path, nil
	}
	if strings.Contains(content, codexMCPEnd) {
		return "", fmt.Errorf("%s: incomplete rsk MCP block", path)
	}
	if servers, ok := cfg["mcp_servers"].(map[string]any); ok {
		if _, exists := servers[rskMCPServerName]; exists {
			return "", fmt.Errorf("%s: MCP server rsk already exists and is not managed by rsk", path)
		}
	} else if _, exists := cfg["mcp_servers"]; exists {
		return "", fmt.Errorf("%s: mcp_servers is not a table", path)
	}
	if err = os.MkdirAll(dir, fsperm.Dir); err != nil {
		return "", err
	}
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	if err = writeCodexConfig(path, content+codexMCPBlock); err != nil {
		return "", err
	}
	return path, nil
}

func (*codexTool) UnregisterMCP(projectDir string) error {
	path := filepath.Join(projectDir, codexFolder, codexConfig)
	data, err := readCodexConfig(path)
	if err != nil {
		return err
	}
	content := string(data)
	if !strings.Contains(content, codexMCPStart) && !strings.Contains(content, codexMCPEnd) {
		return nil
	}
	if strings.Count(content, codexMCPBlock) != 1 || strings.Count(content, codexMCPStart) != 1 ||
		strings.Count(content, codexMCPEnd) != 1 {
		return fmt.Errorf("%s: modified rsk MCP block; resolve manually", path)
	}
	return writeCodexConfig(path, strings.Replace(content, codexMCPBlock, "", 1))
}

func readCodexConfig(path string) ([]byte, error) {
	if err := refuseIfLinkDir(filepath.Dir(path)); err != nil {
		return nil, err
	}
	if err := refuseIfLink(path); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return data, err
}

func writeCodexConfig(path, content string) error {
	return fsx.WriteAtomic(path, rskTempPattern, func(w io.Writer) error {
		_, err := io.WriteString(w, content)
		return err
	})
}
