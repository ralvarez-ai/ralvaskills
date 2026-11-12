package tool

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"go.ralva.dev/ralvaskills/v3/internal/fsperm"
)

const (
	// rskMCPServerName is the key rsk registers its server under.
	rskMCPServerName = "rsk"

	// mcpSubcommand is the argv every client passes to rsk.
	mcpSubcommand = "mcp"
	mcpCommandKey = "command"
	mcpTypeKey    = "type"
	mcpStdioType  = "stdio"

	// mcpServersKey is the JSON key most clients use for their server map.
	mcpServersKey = "mcpServers"

	// mcpReadOnlyFlag keeps a registration to the read-only server, which omits
	// install_skills until the client's per-install approval behavior is verified.
	mcpReadOnlyFlag = "--read-only"
)

// mcpRecognizer is implemented by a client whose MCP layout needs custom
// detection or preflight naming. JSON clients use jsonMCPProvider instead, and
// each client keeps its own recognizer beside its adapter.
type mcpRecognizer interface {
	hasMCP(projectDir string) (bool, error)
	mcpConfigName(projectDir string) (string, error)
}

// ProjectHasMCP reports whether a selected client already has an rsk server.
func ProjectHasMCP(projectDir string, ids []ID) (bool, error) {
	for _, id := range ids {
		found, err := clientHasMCP(projectDir, id)
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}

func clientHasMCP(projectDir string, id ID) (bool, error) {
	if provider, ok := jsonMCPProviderFor(id); ok {
		cfg, err := readJSONConfig(filepath.Join(projectDir, filepath.FromSlash(provider.JSONMCP().file)))
		if err != nil {
			return false, err
		}
		servers, valid := cfg[provider.JSONMCP().key].(map[string]any)
		if !valid {
			return false, nil
		}
		_, exists := servers[rskMCPServerName]
		return exists, nil
	}
	if recognizer, ok := mcpRecognizerFor(id); ok {
		return recognizer.hasMCP(projectDir)
	}
	return false, fmt.Errorf("unknown tool %q", id)
}

// jsonMCPProviderFor returns the JSON MCP provider for id.
//
//nolint:ireturn // returns the interface the caller checks
func jsonMCPProviderFor(id ID) (jsonMCPProvider, bool) {
	t, ok := Get(id)
	if !ok {
		return nil, false
	}
	provider, ok := t.(jsonMCPProvider)
	return provider, ok
}

// mcpRecognizerFor returns the custom MCP recognizer for id, when it has one.
//
//nolint:ireturn // returns the interface the caller checks
func mcpRecognizerFor(id ID) (mcpRecognizer, bool) {
	t, ok := Get(id)
	if !ok {
		return nil, false
	}
	recognizer, ok := t.(mcpRecognizer)
	return recognizer, ok
}

// PreflightMCP checks all selected client configurations before any registration.
// It stages their existing files in a temporary project and invokes the same
// registration logic used for the real project.
func PreflightMCP(projectDir string, ids []ID) error {
	stage, err := os.MkdirTemp("", "rsk-mcp-preflight-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	for _, id := range ids {
		if err = preflightClientMCP(projectDir, stage, id); err != nil {
			return err
		}
	}
	return nil
}

func preflightClientMCP(projectDir, stage string, id ID) error {
	client, ok := Get(id)
	if !ok {
		return fmt.Errorf("unknown tool %q", id)
	}
	name, err := mcpConfigName(projectDir, id)
	if err != nil {
		return err
	}
	source := filepath.Join(projectDir, name)
	if err = refuseIfLink(source); err != nil {
		return err
	}
	data, err := os.ReadFile(source)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		target := filepath.Join(stage, name)
		if err = os.MkdirAll(filepath.Dir(target), fsperm.Dir); err != nil {
			return err
		}
		if err = os.WriteFile(target, data, fsperm.File); err != nil {
			return err
		}
	}
	if _, err = client.RegisterMCP(stage); err != nil {
		return fmt.Errorf("%s MCP preflight: %w", id, err)
	}
	return nil
}

// mcpConfigName returns the project-relative config file a client reads, so
// preflight can stage it.
func mcpConfigName(projectDir string, id ID) (string, error) {
	if provider, ok := jsonMCPProviderFor(id); ok {
		cfg := provider.JSONMCP()
		if err := cfg.checkAlternates(projectDir); err != nil {
			return "", err
		}
		return cfg.file, nil
	}
	if recognizer, ok := mcpRecognizerFor(id); ok {
		return recognizer.mcpConfigName(projectDir)
	}
	return "", fmt.Errorf("MCP preflight is not implemented for %q", id)
}

// matchesAny reports whether v equals one of the accepted legacy shapes. A
// client that changed its server entry passes the old shape here, so the
// generic writer stays ignorant of which client needed it.
func matchesAny(v any, candidates []map[string]any) bool {
	for _, candidate := range candidates {
		if reflect.DeepEqual(v, candidate) {
			return true
		}
	}
	return false
}

func registerMCPServer(
	path, key string,
	server map[string]any,
	legacy []map[string]any,
	read func(string) (map[string]any, error),
	write func(string, map[string]any) error,
) (string, error) {
	if err := refuseIfLink(path); err != nil {
		return "", err
	}
	cfg, err := read(path)
	if err != nil {
		return "", err
	}

	servers, ok := cfg[key].(map[string]any)
	if !ok {
		if _, exists := cfg[key]; exists {
			return "", fmt.Errorf("%s: %s is not an object", path, key)
		}
		servers = make(map[string]any)
	}
	if reflect.DeepEqual(servers[rskMCPServerName], server) {
		return path, nil
	}
	if existing, exists := servers[rskMCPServerName]; exists && !matchesAny(existing, legacy) {
		return "", fmt.Errorf("%s: MCP server rsk already exists with different settings", path)
	}

	servers[rskMCPServerName] = server
	cfg[key] = servers
	if writeErr := write(path, cfg); writeErr != nil {
		return "", writeErr
	}
	return path, nil
}

func removeMCPServer(
	path, key string,
	managed map[string]any,
	legacy []map[string]any,
	read func(string) (map[string]any, error),
	write func(string, map[string]any) error,
) error {
	if err := refuseIfLink(path); err != nil {
		return err
	}
	cfg, err := read(path)
	if err != nil {
		return err
	}

	servers, ok := cfg[key].(map[string]any)
	if !ok {
		return nil
	}
	existing, present := servers[rskMCPServerName]
	if !present {
		return nil
	}
	if !reflect.DeepEqual(existing, managed) && !matchesAny(existing, legacy) {
		return fmt.Errorf("%s: refusing to remove a modified rsk MCP server", path)
	}

	delete(servers, rskMCPServerName)
	if len(servers) == 0 {
		delete(cfg, key)
	} else {
		cfg[key] = servers
	}
	return write(path, cfg)
}
