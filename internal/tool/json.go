package tool

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal"
	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
	"go.ralva.dev/ralvaskills/v3/internal/fsx"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

// jsonMCPConfig describes where a client keeps its MCP servers inside a JSON
// config file. It is the only per-client MCP data the JSON adapters carry.
type jsonMCPConfig struct {
	// file is the project-relative path of the config file.
	file string

	// key is the top-level key holding the server map. It is matched literally,
	// so a dotted key such as Amp's "amp.mcpServers" works unchanged.
	key string

	// server returns the entry rsk registers under rskMCPServerName.
	server func() map[string]any

	// alternate lists project-relative files the same client also reads. When
	// one exists, registration refuses instead of creating a second source of
	// truth the client resolves by precedence.
	alternate []string
}

// jsonMCPProvider is implemented by every tool whose MCP servers live in a JSON
// object at a fixed key. Detection and preflight use it instead of a per-tool
// switch, so a new JSON adapter needs no change in mcp.go.
type jsonMCPProvider interface {
	JSONMCP() jsonMCPConfig
}

// Register writes the rsk server into the client's config, preserving unrelated
// servers and keys, and returns the path written.
func (c jsonMCPConfig) Register(projectDir string) (string, error) {
	path, err := c.locate(projectDir)
	if err != nil {
		return "", err
	}
	if err = os.MkdirAll(filepath.Dir(path), fsperm.Dir); err != nil {
		return "", err
	}
	return registerMCPServer(path, c.key, c.server(), nil, readJSONConfig, writeJSONConfig)
}

// Unregister removes the rsk server from the client's config, leaving other
// servers intact. It returns nil when nothing was registered.
func (c jsonMCPConfig) Unregister(projectDir string) error {
	path, err := c.locate(projectDir)
	if err != nil {
		return err
	}
	return removeMCPServer(path, c.key, c.server(), nil, readJSONConfig, writeJSONConfig)
}

// locate resolves the config path and refuses the shapes rsk cannot write
// safely: an alternate config the client would prefer, or a symlinked directory
// that would put the file outside the project.
func (c jsonMCPConfig) locate(projectDir string) (string, error) {
	if err := c.checkAlternates(projectDir); err != nil {
		return "", err
	}
	path := filepath.Join(projectDir, filepath.FromSlash(c.file))
	if dir := filepath.Dir(path); dir != filepath.Clean(projectDir) {
		if err := refuseIfLinkDir(dir); err != nil {
			return "", err
		}
	}
	return path, nil
}

// checkAlternates rejects a config file in an alternate location, which the
// client may resolve ahead of the one rsk writes.
func (c jsonMCPConfig) checkAlternates(projectDir string) error {
	return rejectAlternates(projectDir, c.file, c.alternate)
}

// rejectAlternates fails when any project-relative alternate file exists. A
// client that reads two files would otherwise get two sources of truth whose
// precedence rsk cannot see.
func rejectAlternates(projectDir, preferred string, alternates []string) error {
	for _, alternate := range alternates {
		path := filepath.Join(projectDir, filepath.FromSlash(alternate))
		if _, err := os.Lstat(path); err == nil {
			return fmt.Errorf(
				"%s: this client also reads it; edit it or remove it, then retry (%s stays rsk-managed)",
				path, preferred,
			)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

// jsonAdapter implements Tool for a client that discovers skills natively and
// keeps MCP servers in a JSON object. Such a client needs no generated import,
// so the whole adapter is the three fields below; only the constants and the
// MCP config differ between clients.
type jsonAdapter struct {
	id         ID
	skillsDir  func() string
	projectDir func(string) string
	mcp        jsonMCPConfig
}

func (a *jsonAdapter) ID() ID {
	return a.id
}

func (a *jsonAdapter) SkillsDir() string {
	return a.skillsDir()
}

func (a *jsonAdapter) ProjectSkillsDir(projectRoot string) string {
	return a.projectDir(projectRoot)
}

// SyncPinned is a no-op: the client scans skills/<name>/SKILL.md itself and
// reads AGENTS.md, so there is no generated import for rsk to write.
func (*jsonAdapter) SyncPinned(string, []string) error {
	return nil
}

// RemovePinned is a no-op: rsk never writes a skill import for this client.
func (*jsonAdapter) RemovePinned(string) error {
	return nil
}

func (a *jsonAdapter) RegisterMCP(projectDir string) (string, error) {
	return a.mcp.Register(projectDir)
}

func (a *jsonAdapter) UnregisterMCP(projectDir string) error {
	return a.mcp.Unregister(projectDir)
}

func (a *jsonAdapter) JSONMCP() jsonMCPConfig {
	return a.mcp
}

// agentsSkillsDir is the project skills directory most clients share. The
// Agent Skills convention made .agents/skills the interoperable path, so one
// link per skill serves every tool that reads it.
func agentsSkillsDir(projectRoot string) string {
	return filepath.Join(projectRoot, internal.AgentsFolderName, skill.SkillsFolderName)
}

// homePath joins parts to the user's home directory. It returns "" when the home
// directory cannot be resolved, matching the other adapters.
func homePath(parts ...string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(append([]string{home}, parts...)...)
}

// readOnlyMCPServer is the entry rsk registers for every client whose
// per-install approval behavior is unverified. Omitting install_skills keeps a
// client from offering a destructive tool it cannot gate.
func readOnlyMCPServer() map[string]any {
	return map[string]any{
		mcpCommandKey: rskMCPServerName,
		"args":        []any{"mcp", mcpReadOnlyFlag},
	}
}

// readJSONConfig reads a JSON object config, returning an empty map when the
// file does not exist.
func readJSONConfig(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return make(map[string]any), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	cfg := make(map[string]any)
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return cfg, nil
}

// writeJSONConfig writes a JSON object config atomically, preserving the file's
// mode when it already exists.
func writeJSONConfig(path string, cfg map[string]any) error {
	return fsx.WriteAtomic(path, rskTempPattern, func(w io.Writer) error {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(cfg)
	})
}

// refuseIfLink refuses to write through a symlink, or over a non-regular file.
//
// Failing loudly is deliberate: someone who symlinked the file did it on
// purpose, and rsk cannot know whether the target wants an rsk-managed entry,
// so it must not guess in either direction.
func refuseIfLink(path string) error {
	return refuseLink(path, false)
}

// refuseIfLinkDir rejects a symlinked directory so rsk does not write a config
// outside the project.
func refuseIfLinkDir(dir string) error {
	return refuseLink(dir, true)
}

func refuseLink(path string, wantDir bool) error {
	info, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if skill.IsLink(info) {
		return fmt.Errorf(
			"%s is a symlink (→ %s): rsk will not write through it, because that would modify "+
				"the target or replace the link. Remove the link, then retry",
			path, linkTarget(path),
		)
	}
	if !wantDir && !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file (mode %s): refusing to write",
			path, info.Mode())
	}
	return nil
}

// linkTarget returns a symlink's target for an error message, or "unknown" when
// the mode bit came from a Windows junction, which has no readable target.
func linkTarget(path string) string {
	target, err := os.Readlink(path)
	if err != nil {
		return "unknown"
	}
	return target
}

// ProjectMCPFile returns the project-relative MCP config file a client reads,
// for the commands that print what a sync will check.
func ProjectMCPFile(id ID) (string, bool) {
	t, ok := Get(id)
	if !ok {
		return "", false
	}
	provider, ok := t.(jsonMCPProvider)
	if !ok || provider.JSONMCP().file == "" {
		return "", false
	}
	return provider.JSONMCP().file, true
}

// trimSlash is used by adapters that accept a directory from the environment.
func trimSlash(dir string) string {
	return strings.TrimRight(strings.TrimSpace(dir), `/\`)
}
