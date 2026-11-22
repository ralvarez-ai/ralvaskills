// Package project sets up an rsk project directory: the manifest, the
// per-tool pin files, the MCP registration and the rsk-guide skill. It never
// writes to the terminal; it returns the notices a caller shows, in order.
package project

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal"
	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/install"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

// Level is how a Notice is shown.
type Level int

const (
	// Info is a neutral message.
	Info Level = iota
	// Warn is a skipped or degraded step that does not fail the command.
	Warn
	// Success reports a step that completed.
	Success
)

// Notice is one message for the user, in the order it was produced.
type Notice struct {
	Level Level
	Text  string
}

// SetupOptions describe what Setup configures.
type SetupOptions struct {
	// Tools are the AI tools to configure.
	Tools []tool.ID
	// ForLabel is how the --for choice is named in messages.
	ForLabel string
	// RegisterMCP also writes the rsk MCP server into each tool's project config.
	RegisterMCP bool
	// InstallGuide also installs and pins the rsk-guide skill.
	InstallGuide bool
}

// SetupResult is what Setup did. Notices is filled even when Setup fails, so
// the caller can show what happened before the failure.
type SetupResult struct {
	Notices []Notice
	// Created is true when no project existed and Setup made one.
	Created bool
}

// Setup creates or updates the rsk project in cwd for opts.Tools: it reads or
// builds the manifest, writes it, syncs the pinned instructions, and optionally
// registers the MCP server and installs rsk-guide.
func Setup(ctx context.Context, cwd string, opts SetupOptions) (SetupResult, error) {
	var res SetupResult
	if opts.RegisterMCP {
		if err := tool.PreflightMCP(cwd, opts.Tools); err != nil {
			return res, err
		}
	}

	m, rskDir, created, notices, err := loadOrCreateManifest(cwd, opts.Tools, opts.ForLabel)
	res.Notices = append(res.Notices, notices...)
	res.Created = created
	if err != nil {
		return res, err
	}

	if err = manifest.WriteMod(rskDir, m); err != nil {
		return res, err
	}
	if err = install.SyncPinnedAllTools(rskDir, m); err != nil {
		return res, err
	}

	if opts.RegisterMCP {
		notes, mcpErr := RegisterMCP(cwd, opts.Tools)
		res.Notices = append(res.Notices, notes...)
		if mcpErr != nil {
			return res, fmt.Errorf(
				"MCP registration incomplete; previously written project configuration was retained: %w",
				mcpErr,
			)
		}
	}

	if opts.InstallGuide {
		notes, guideErr := installRskGuide(ctx, cwd, rskDir)
		res.Notices = append(res.Notices, notes...)
		if guideErr != nil {
			return res, guideErr
		}
	}
	return res, nil
}

// loadOrCreateManifest reads the existing project manifest and adds any of
// tools it's missing, or builds a fresh one if no project exists yet.
func loadOrCreateManifest(
	cwd string, tools []tool.ID, forLabel string,
) (m manifest.Mod, rskDir string, isNew bool, notices []Notice, err error) {
	existing, existErr := manifest.ProjectFolderPath()
	if existErr == nil {
		loaded, readErr := manifest.ReadMod(existing)
		if readErr != nil {
			return manifest.Mod{}, "", false, nil, readErr
		}

		var added []tool.ID
		for _, t := range tools {
			if !slices.Contains(loaded.Tools, t) {
				added = append(added, t)
				loaded.Tools = append(loaded.Tools, t)
			}
		}
		if len(added) == 0 {
			return loaded, existing, false,
				[]Notice{{Info, "rsk project already configured for " + forLabel}}, nil
		}
		return loaded, existing, false,
			[]Notice{{Info, fmt.Sprintf("Adding tools to existing project: %v", added)}}, nil
	}

	// No project yet — existErr just signals "not found", nothing to surface.
	return manifest.Mod{
		Tools:  tools,
		Skills: make(map[string]string),
		Pinned: []string{},
	}, filepath.Join(cwd, internal.ProjectFolderName), true, nil, nil
}

// RegisterMCP writes the rsk server into each configured tool's project config
// and, for Claude Code, the workflow pointer in ./CLAUDE.md. It skips, with a
// warning notice, when rsk is not on PATH, since the client would fail to
// launch the server.
func RegisterMCP(cwd string, tools []tool.ID) ([]Notice, error) {
	if _, err := exec.LookPath("rsk"); err != nil {
		//nolint:nilerr // not being on PATH is a skip with a warning, not a failure
		return []Notice{{Warn, "'rsk' is not on PATH — skipping MCP registration. Put rsk on PATH " +
			"first, or the client will fail to launch the server."}}, nil
	}

	paths := make([]string, 0, len(tools))
	for _, id := range tools {
		t, ok := tool.Get(id)
		if !ok {
			continue
		}
		path, err := t.RegisterMCP(cwd)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	if slices.Contains(tools, tool.ClaudeID) {
		if err := tool.WriteMCPPointer(cwd); err != nil {
			return nil, err
		}
	}

	return []Notice{{Success, "registered the rsk MCP server (" + strings.Join(paths, ", ") + ")"}}, nil
}

// installRskGuide installs rsk-guide into the project and pins it, so the agent
// loads the rsk policy automatically. It is best-effort: a missing config or an
// unresolvable skill warns instead of failing the whole command.
func installRskGuide(ctx context.Context, cwd, rskDir string) ([]Notice, error) {
	cfg, err := config.Load()
	if err != nil {
		//nolint:nilerr // rsk not initialized is a skip with a warning, not a failure
		return []Notice{{Warn, "skipping rsk-guide install: " + err.Error() + "\n  Run 'rsk init' first."}}, nil
	}

	var notices []Notice
	bundles, catalogErr := config.LoadCatalog("")
	if catalogErr != nil {
		notices = append(notices, Notice{Warn, "user catalog: " + catalogErr.Error()})
	}

	chain, resolveErr := cfg.Sources()
	var skills []skill.Skill
	if resolveErr == nil {
		skills, _, resolveErr = install.Resolve(ctx, []string{"rsk-guide"}, bundles, chain)
	}
	if resolveErr != nil {
		//nolint:nilerr // best-effort convenience install; the project is still set up
		return append(notices, Notice{Warn, "could not resolve rsk-guide: " + resolveErr.Error()}), nil
	}

	m, err := manifest.ReadMod(rskDir)
	if err != nil {
		return notices, err
	}
	results, err := install.Apply(skills, install.Options{
		Scope:   install.ScopeProject,
		Targets: install.ProjectSkillsDirs(cwd, m),
		RskDir:  rskDir,
		Pin:     true,
	})
	if err != nil {
		return notices, err
	}
	for _, r := range results {
		if r.Err != nil {
			notices = append(notices, Notice{Warn, "rsk-guide: " + r.Err.Error()})
			continue
		}
		notices = append(notices, Notice{Success, "installed and pinned rsk-guide"})
	}
	return notices, nil
}
