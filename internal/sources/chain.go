package sources

import (
	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
)

// ChainOptions are the per-command choices about how declared sources are read.
type ChainOptions struct {
	// Offline makes a git source answer from its existing index and never
	// contact its remote.
	Offline bool
}

// Chain builds the machine's source chain: the built-in sources first (the
// local clone or the hosted registry, then the official cache), then every
// declared git source in file order. Each git source is gated by trust.json, so
// an untrusted one fails when it is consulted instead of being skipped. Declared
// path sources are not resolvable yet and are left out.
func Chain(cfg config.Config, o ChainOptions) (source.Chain, error) {
	return ChainAt(cfg, DefaultPath(), trust.DefaultPath(), o)
}

// ChainAt is Chain reading the declarations and grants from the given files.
func ChainAt(cfg config.Config, sourcesPath, trustPath string, o ChainOptions) (source.Chain, error) {
	declared, err := Load(sourcesPath)
	if err != nil {
		return nil, err
	}
	grants, err := trust.Load(trustPath)
	if err != nil {
		return nil, err
	}

	specs := cfg.SourceSpecs()
	for i := range declared.Sources {
		if declared.Sources[i].Type == source.KindGit {
			specs = append(specs, declared.Sources[i].Spec())
		}
	}

	opts := cfg.SourceOptions()
	opts.Offline = o.Offline
	opts.Gate = func(id, url string) error { return trust.Require(grants, id, url) }
	return source.Build(specs, opts)
}
