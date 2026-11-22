package mcp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"go.ralva.dev/lib/mcpserver"

	"go.ralva.dev/ralvaskills/v3/internal/catalog"
	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

// Deps is the environment the MCP server needs, injected so the server is
// testable without the CLI.
type Deps struct {
	Cfg      config.Config
	Version  string
	Logger   *slog.Logger
	ReadOnly bool

	// RepoFor supplies the clone a git source reads, for tests that must not
	// reach a network. Nil means the real one.
	RepoFor func(source.Spec) source.RepoOps
}

// chain builds the built-in source chain, which is all the MCP server resolves
// through: declared sources are the user's to add and trust from a terminal.
func (d Deps) chain() (source.Chain, error) {
	opts := d.Cfg.SourceOptions()
	opts.RepoFor = d.RepoFor
	return source.Build(d.Cfg.SourceSpecs(), opts)
}

// NewServer builds the MCP server with every tool and resource registered. The
// result is transport-independent: the caller decides whether to run it over
// stdio or in memory.
func NewServer(deps Deps) *sdk.Server {
	srv := sdk.NewServer(
		&sdk.Implementation{Name: "rsk", Title: "rsk", Version: deps.Version},
		&sdk.ServerOptions{Logger: deps.Logger},
	)
	srv.AddReceivingMiddleware(mcpserver.RecoverWith(mcpserver.RecoverConfig{Logger: deps.Logger}))
	registerTools(srv, deps)
	registerCatalogResource(srv, deps)
	return srv
}

// ServeStdio runs srv over stdin/stdout until the client closes the connection.
// Protocol frames go to stdout; logs must stay on stderr.
func ServeStdio(ctx context.Context, srv *sdk.Server) error {
	if err := srv.Run(ctx, &sdk.StdioTransport{}); err != nil {
		return fmt.Errorf("run mcp server: %w", err)
	}
	return nil
}

// catalog lists the skills available to this machine: from the local clone in
// local mode, or from the cached registry index otherwise.
func (d Deps) catalog(ctx context.Context) ([]catalog.Entry, error) {
	chain, err := d.chain()
	if err != nil {
		return nil, err
	}
	if d.Cfg.LocalMode() {
		local := chain.Primary()
		if len(local) == 0 {
			return nil, errors.New("no local source is configured")
		}
		return catalog.Build(ctx, local[0])
	}

	reg, ok := chain.Registry()
	if !ok {
		return nil, errors.New("no hosted registry is configured")
	}
	snap, err := NewIndexCache(d.Cfg.RegistryCache(), DefaultIndexTTL, reg.Index).Load(ctx)
	if err != nil {
		return nil, err
	}
	return catalog.FromIndex(snap.Skills), nil
}
