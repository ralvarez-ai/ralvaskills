package sources_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/source/sourcetest"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
)

func TestDeriveGivesTheSSHURLFirstAndHTTPSAsTheFallback(t *testing.T) {
	t.Parallel()
	c, err := sources.Derive("@Acme/Skills")
	if err != nil {
		t.Fatal(err)
	}
	if c.ID != "@acme/skills" || c.SSH != "git@github.com:acme/skills.git" ||
		c.HTTPS != "https://github.com/acme/skills.git" {
		t.Errorf("candidate = %+v", c)
	}
	for _, u := range c.URLs() {
		if _, parseErr := git.ParseRemoteURL(u); parseErr != nil {
			t.Errorf("derived URL %q is not one the runner accepts: %v", u, parseErr)
		}
	}
}

func TestDeriveRefusesWhatIsNotAnImplicitGitHubSource(t *testing.T) {
	t.Parallel()
	for _, id := range []string{
		skill.OfficialID, skill.DefaultRegistryID, "@ralvarezdev/ralvaskills", // built in, by name and by alias
		"@local/scratch", "@bad", "acme/skills", "@acme/skills/extra", "@acme/--upload-pack=x",
	} {
		if c, err := sources.Derive(id); err == nil {
			t.Errorf("Derive(%q) = %+v, want an error", id, c)
		}
	}
}

type probeLog struct {
	urls   []string
	failOn map[string]error
	commit string
}

func (p *probeLog) probe(_ context.Context, remote git.RemoteURL, _ string) (string, error) {
	p.urls = append(p.urls, remote.String())
	if err := p.failOn[remote.String()]; err != nil {
		return "", err
	}
	return p.commit, nil
}

func allow(string, string) error {
	return nil
}

func TestChooseKeepsSSHWhenItWorks(t *testing.T) {
	t.Parallel()
	c, _ := sources.Derive("@acme/skills")
	p := &probeLog{commit: strings.Repeat("a", 40)}

	e, err := sources.Choose(t.Context(), c, "", allow, p.probe)
	if err != nil || e.URL != c.SSH || e.ID != c.ID || e.Type != source.KindGit {
		t.Fatalf("entry = %+v, %v", e, err)
	}
	if len(p.urls) != 1 || p.urls[0] != c.SSH {
		t.Errorf("probed %v, want only SSH", p.urls)
	}
}

func TestChooseFallsBackToHTTPSWhenSSHCannotBeUsed(t *testing.T) {
	t.Parallel()
	c, _ := sources.Derive("@acme/skills")
	p := &probeLog{
		commit: strings.Repeat("a", 40),
		failOn: map[string]error{c.SSH: errors.New("Permission denied (publickey)")},
	}

	e, err := sources.Choose(t.Context(), c, "main", allow, p.probe)
	if err != nil || e.URL != c.HTTPS {
		t.Fatalf("entry = %+v, %v", e, err)
	}
	if strings.Join(p.urls, " ") != c.SSH+" "+c.HTTPS {
		t.Errorf("probed %v, want SSH then HTTPS", p.urls)
	}
}

func TestChooseFailsNamingBothAttemptsWhenNeitherWorks(t *testing.T) {
	t.Parallel()
	c, _ := sources.Derive("@acme/private")
	p := &probeLog{failOn: map[string]error{c.SSH: errors.New("no key"), c.HTTPS: errors.New("not found")}}

	_, err := sources.Choose(t.Context(), c, "", allow, p.probe)
	if err == nil || !strings.Contains(err.Error(), "no key") || !strings.Contains(err.Error(), "not found") {
		t.Errorf("err = %v, want both failures", err)
	}
}

func TestChooseAsksTheGateBeforeContactingEitherURL(t *testing.T) {
	t.Parallel()
	c, _ := sources.Derive("@acme/skills")
	var gated []string
	gate := func(id, url string) error {
		gated = append(gated, url)
		if url == c.HTTPS {
			return &trust.UntrustedError{ID: id, URL: url}
		}
		return nil
	}
	p := &probeLog{failOn: map[string]error{c.SSH: errors.New("no key")}}

	_, err := sources.Choose(t.Context(), c, "", gate, p.probe)
	if !errors.Is(err, trust.ErrUntrusted) {
		t.Errorf("err = %v, want the gate's refusal", err)
	}
	if strings.Join(p.urls, " ") != c.SSH {
		t.Errorf("probed %v: the HTTPS URL was contacted although the gate refused it", p.urls)
	}

	refuseAll := func(id, url string) error { return &trust.UntrustedError{ID: id, URL: url} }
	p = &probeLog{}
	if _, err = sources.Choose(t.Context(), c, "", refuseAll, p.probe); err == nil || len(p.urls) != 0 {
		t.Errorf("err = %v, probed %v: nothing may be contacted past a refusing gate", err, p.urls)
	}
}

// implicitChain builds the machine's chain with one source offered for this run.
func implicitChain(t *testing.T, m machine, implicit ...sources.Entry) (source.Chain, *sourcetest.Repo) {
	t.Helper()
	repo := &sourcetest.Repo{Commit: strings.Repeat("d", 40), Skills: map[string]string{"alpha": "1.0.0"}}
	official := &sourcetest.Repo{Commit: strings.Repeat("c", 40)}
	chain, err := sources.ChainAt(m.cfg, m.sources, m.trustFile, sources.ChainOptions{
		Implicit: implicit,
		RepoFor: func(spec source.Spec) source.RepoOps {
			if spec.ID == skill.OfficialID {
				return official
			}
			return repo
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return chain, repo
}

func TestAnImplicitSourceResolvesForOneRunThroughTheSameGateAndLeavesNoGrant(t *testing.T) {
	t.Parallel()
	m := newMachine(t)
	c, _ := sources.Derive("@acme/skills")
	entry := c.Entry(c.SSH)

	// Without the run's approval the source is not in the chain at all.
	plain := m.chain(t, false)
	if _, err := plain.Find(
		t.Context(),
		skill.Ref{Source: c.ID, Name: "alpha"},
	); !errors.Is(
		err,
		source.ErrUnknownSource,
	) {
		t.Fatalf("an unregistered source resolved without approval: %v", err)
	}

	chain, repo := implicitChain(t, m, entry)
	got, err := chain.Find(t.Context(), skill.Ref{Source: c.ID, Name: "alpha"})
	if err != nil || got.Git == nil || got.Git.URL != c.SSH || got.Git.Commit != strings.Repeat("d", 40) {
		t.Fatalf("Find = %+v, %v", got, err)
	}
	if len(repo.Fetches) == 0 {
		t.Error("the source was never fetched")
	}

	for _, name := range []string{"trust.json", "sources.json"} {
		if _, statErr := os.Stat(filepath.Join(filepath.Dir(m.trustFile), name)); statErr == nil {
			t.Errorf("a one-run source wrote %s", name)
		}
	}
	if _, err = m.chain(t, false).
		Find(t.Context(), skill.Ref{Source: c.ID, Name: "alpha"}); !errors.Is(
		err,
		source.ErrUnknownSource,
	) {
		t.Errorf("the next run still knows the source: %v", err)
	}
}

func TestAnImplicitSourceIsHeldToTheSameRulesAsADeclaredOne(t *testing.T) {
	t.Parallel()
	m := newMachine(t)
	declared := m.declare(t, "@acme/skills", sources.Options{})
	hostile, err := git.ParseRemoteURL("https://github.com/acme/other.git")
	if err != nil {
		t.Fatal(err)
	}
	for name, e := range map[string]sources.Entry{
		"already declared": {ID: "@acme/skills", Type: source.KindGit, URL: declared.URL},
		"built in":         {ID: skill.OfficialID, Type: source.KindGit, URL: source.OfficialSkillsURL},
		"not github":       {ID: "@acme/other", Type: source.KindGit, URL: "git@git.example.com:acme/other.git"},
		"identity vs url":  {ID: "@acme/skills2", Type: source.KindGit, URL: hostile.String()},
		"a path":           {ID: "@me/scratch", Type: source.KindPath, URL: t.TempDir()},
		"hostile url":      {ID: "@acme/x", Type: source.KindGit, URL: "ext::sh -c id"},
	} {
		_, buildErr := sources.ChainAt(
			m.cfg,
			m.sources,
			m.trustFile,
			sources.ChainOptions{Implicit: []sources.Entry{e}},
		)
		if buildErr == nil {
			t.Errorf("%s: an implicit source was accepted", name)
		}
	}
}

func TestRestoreUsesTheMachinesOwnTransportForTheSameRepository(t *testing.T) {
	t.Parallel()
	m := newMachine(t)
	l := source.LockedSkill{
		Name:   "alpha",
		Commit: strings.Repeat("d", 40),
		Tree:   sourcetest.TreeOf("skills/alpha"),
		Path:   "skills/alpha",
	}
	e := m.declare(t, "@acme/skills", sources.Options{}) // this machine uses SSH
	m.grant(t, e)
	repo := &sourcetest.Repo{Commit: l.Commit, Skills: map[string]string{"alpha": ""}}
	opts := sources.ChainOptions{RepoFor: func(source.Spec) source.RepoOps { return repo }}

	// The lock was written on a machine that used HTTPS for the same repository.
	got, err := sources.RestoreAt(t.Context(), m.cfg, m.sources, m.trustFile, opts,
		"@acme/skills", "https://github.com/acme/skills.git", l)
	if err != nil || got.Git == nil || got.Git.URL != e.URL {
		t.Fatalf("restore = %+v, %v: want the machine's own URL %s", got, err, e.URL)
	}

	// A different repository under the same identity is still refused.
	if _, err = sources.RestoreAt(t.Context(), m.cfg, m.sources, m.trustFile, opts,
		"@acme/skills", "https://github.com/evil/skills.git", l); err == nil {
		t.Error("a lock naming another repository was followed")
	}
}

func TestRestoreOfAnUndeclaredSourceNeedsTheRunsApproval(t *testing.T) {
	t.Parallel()
	m := newMachine(t)
	l := source.LockedSkill{
		Name:   "alpha",
		Commit: strings.Repeat("d", 40),
		Tree:   sourcetest.TreeOf("skills/alpha"),
		Path:   "skills/alpha",
	}
	repo := &sourcetest.Repo{Commit: l.Commit, Skills: map[string]string{"alpha": ""}}
	c, _ := sources.Derive("@acme/skills")

	opts := sources.ChainOptions{RepoFor: func(source.Spec) source.RepoOps { return repo }}
	var notDeclared *sources.NotDeclaredError
	if _, err := sources.RestoreAt(t.Context(), m.cfg, m.sources, m.trustFile, opts,
		c.ID, "https://github.com/acme/skills.git", l); !errors.As(err, &notDeclared) {
		t.Fatalf("err = %v, want NotDeclaredError without the run's approval", err)
	}
	if !strings.Contains(notDeclared.Error(), "--trust @acme/skills") {
		t.Errorf("the refusal does not name the one-run flag: %v", notDeclared)
	}

	opts.Implicit = []sources.Entry{c.Entry(c.SSH)}
	got, err := sources.RestoreAt(t.Context(), m.cfg, m.sources, m.trustFile, opts,
		c.ID, "https://github.com/acme/skills.git", l)
	if err != nil || got.Git == nil || got.Git.URL != c.SSH {
		t.Fatalf("restore = %+v, %v", got, err)
	}
}
