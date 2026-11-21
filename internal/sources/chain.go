package sources

import (
	"context"
	"fmt"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
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

// NotDeclaredError says a lock names a git source this machine has not
// declared. The fix is a command, so the message gives it.
type NotDeclaredError struct {
	ID string
}

// Error implements error.
func (e *NotDeclaredError) Error() string {
	return fmt.Sprintf(
		"source %s is in the lock but not declared on this machine: add it with 'rsk source add %s'",
		e.ID,
		e.ID,
	)
}

// locker is implemented by a source that can restore a skill at a locked commit.
type locker interface {
	Locked(ctx context.Context, l source.LockedSkill) (skill.Skill, error)
}

// Restore resolves a git skill at the commit a lock recorded. A lock names its
// source, so this never searches the chain; it still goes through the trust
// gate, and a lock whose URL differs from the one declared here is refused
// instead of followed.
func Restore(
	ctx context.Context, cfg config.Config, o ChainOptions, id, url string, l source.LockedSkill,
) (skill.Skill, error) {
	return RestoreAt(ctx, cfg, DefaultPath(), trust.DefaultPath(), o, id, url, l)
}

// RestoreAt is Restore reading the declarations and grants from the given files.
func RestoreAt(
	ctx context.Context, cfg config.Config, sourcesPath, trustPath string, o ChainOptions,
	id, url string, l source.LockedSkill,
) (skill.Skill, error) {
	declared, err := Load(sourcesPath)
	if err != nil {
		return skill.Skill{}, err
	}
	entry, ok := Find(declared, id)
	if !ok || entry.Type != source.KindGit {
		return skill.Skill{}, &NotDeclaredError{ID: skill.CanonicalSourceID(id)}
	}
	if entry.URL != url {
		return skill.Skill{}, fmt.Errorf(
			"the lock fetches %s from %s but this machine declares %s: refusing to follow it", entry.ID, url, entry.URL,
		)
	}

	chain, err := ChainAt(cfg, sourcesPath, trustPath, o)
	if err != nil {
		return skill.Skill{}, err
	}
	for _, r := range chain.Select(source.KindGit) {
		if r.Spec().ID != entry.ID {
			continue
		}
		lk, isLocker := r.(locker)
		if !isLocker {
			break
		}
		return lk.Locked(ctx, l)
	}
	return skill.Skill{}, &NotDeclaredError{ID: entry.ID}
}
