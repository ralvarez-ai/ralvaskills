package sources

import (
	"context"
	"fmt"
	"time"

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

	// Implicit are GitHub sources this machine has not declared, offered for this
	// run only: each was approved by the user at a terminal or with --trust, and
	// is granted in memory so it passes the same gate as a declared source. Nothing
	// is written to sources.json or trust.json, and the MCP server never sets it.
	Implicit []Entry

	// RepoFor supplies the clone a git source reads, for tests above this
	// package that must not reach a network. Nil means the real one.
	RepoFor func(source.Spec) source.RepoOps
}

// Chain builds the machine's source chain: the built-in sources first (the
// local clone or the hosted registry, then the official cache), then every
// declared source in file order. Each git source is gated by trust.json, so an
// untrusted one fails when it is consulted instead of being skipped. A declared
// path source needs no grant: it reads a directory and contacts nothing.
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
		specs = append(specs, declared.Sources[i].Spec())
	}
	seen := map[string]bool{}
	for i := range o.Implicit {
		if err = checkImplicit(o.Implicit[i], declared, seen); err != nil {
			return nil, err
		}
		specs = append(specs, o.Implicit[i].Spec())
		grants = trust.Grant(grants, o.Implicit[i].ID, o.Implicit[i].URL, time.Now())
	}

	opts := cfg.SourceOptions()
	opts.Offline = o.Offline
	opts.RepoFor = o.RepoFor
	opts.Gate = func(id, url string) error {
		if config.BuiltinTrusted(id, url) {
			return nil
		}
		return trust.Require(grants, id, url)
	}
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
		"source %s is in the lock but not declared on this machine: add it with 'rsk source add %s', "+
			"or approve it for this run only with --trust %s",
		e.ID, e.ID, e.ID,
	)
}

func findEntry(entries []Entry, id string) (Entry, bool) {
	return Find(File{Sources: entries}, id)
}

// locker is implemented by a source that can restore a skill at a locked commit.
type locker interface {
	Locked(ctx context.Context, l source.LockedSkill) (skill.Skill, error)
}

// Restore resolves a git skill at the commit a lock recorded. A lock names its
// source, so this never searches the chain; it still goes through the trust
// gate. The machine's own declaration decides the URL: a lock written over SSH
// restores over HTTPS on a machine that declares it so, while a lock naming
// another host or repository is refused instead of followed.
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
	entry, err := restoreEntry(sourcesPath, o, id, url)
	if err != nil {
		return skill.Skill{}, err
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

// restoreEntry finds the declaration a lock entry restores through: the built-in
// official source at its own URL, a declared source, or one approved for this run.
func restoreEntry(sourcesPath string, o ChainOptions, id, url string) (Entry, error) {
	canonical := skill.CanonicalSourceID(id)
	if canonical == skill.OfficialID {
		// The official source is built in: it is not declared, and the only URL it
		// may be restored from is its own.
		if url != source.OfficialSkillsURL {
			return Entry{}, fmt.Errorf(
				"the lock fetches %s from %s but it is built in at %s: refusing to follow it",
				canonical, url, source.OfficialSkillsURL,
			)
		}
		return Entry{ID: canonical, URL: url}, nil
	}
	declared, err := Load(sourcesPath)
	if err != nil {
		return Entry{}, err
	}
	found, ok := Find(declared, id)
	if !ok {
		found, ok = findEntry(o.Implicit, id)
	}
	if !ok || found.Type != source.KindGit {
		return Entry{}, &NotDeclaredError{ID: canonical}
	}
	if !sameRepository(found.ID, found.URL, url) {
		return Entry{}, fmt.Errorf(
			"the lock fetches %s from %s but this machine declares %s: refusing to follow it",
			found.ID, url, found.URL,
		)
	}
	return found, nil
}
