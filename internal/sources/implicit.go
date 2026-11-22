package sources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

type (
	// Candidate is a GitHub repository named by a qualified reference that this
	// machine has not declared: the URLs rsk would contact if the user approves
	// it for one run. SSH comes first so a private repository works with the
	// user's own key; HTTPS is the fallback for a machine with none.
	Candidate struct {
		ID    string
		SSH   string
		HTTPS string
	}

	// Probe asks a remote for the commit rev points to, without fetching
	// anything. git.LsRemote is the production one; tests supply their own.
	Probe func(ctx context.Context, remote git.RemoteURL, rev string) (string, error)
)

// Derive builds the candidate for a source identity. It refuses what cannot be
// derived: a built-in identity (which resolves through its own source, never
// through an implicit clone), the reserved @local scope, and anything that is
// not @owner/repo.
func Derive(id string) (Candidate, error) {
	canonical, err := skill.ParseSourceID(id)
	if err != nil {
		return Candidate{}, err
	}
	switch {
	case reserved(canonical):
		return Candidate{}, fmt.Errorf("source %s is built in and has no implicit form", canonical)
	case strings.HasPrefix(canonical, skill.LocalScope+"/"):
		return Candidate{}, fmt.Errorf("the %s scope is reserved for path sources", skill.LocalScope)
	}
	path := strings.TrimPrefix(canonical, "@")
	return Candidate{
		ID:    canonical,
		SSH:   "git@" + githubHost + ":" + path + ".git",
		HTTPS: "https://" + githubHost + "/" + path + ".git",
	}, nil
}

// ImplicitHint says what `rsk install` would do with a qualified reference to
// the source id when this machine has not registered it, or "" when nothing can
// be derived. It contacts nothing.
func ImplicitHint(id string) string {
	c, err := Derive(id)
	if err != nil {
		return ""
	}
	return fmt.Sprintf(
		"rsk install would derive %s (%s if SSH cannot be used) and ask you to approve it for this run only",
		c.SSH, c.HTTPS,
	)
}

// URLs returns the candidate's URLs in the order they are tried.
func (c Candidate) URLs() []string {
	return []string{c.SSH, c.HTTPS}
}

// Entry returns the declaration for the candidate at url, which must be one of
// its URLs.
func (c Candidate) Entry(url string) Entry {
	return Entry{ID: c.ID, Type: source.KindGit, URL: url}
}

// Choose picks the transport for an approved candidate: SSH when git can use
// it, HTTPS when it cannot. Both URLs are tried through gate before they are
// contacted, so approving the candidate for a run goes through the same trust
// check as any other source, never around it, and a refusal stops the search.
// Any SSH failure counts as "cannot be used": a missing key, an unknown host
// key and a refused connection all look the same without a prompt, and HTTPS
// then fails by itself on a private repository.
func Choose(
	ctx context.Context, c Candidate, rev string, gate func(id, url string) error, probe Probe,
) (Entry, error) {
	var failures []string
	for _, url := range c.URLs() {
		if err := gate(c.ID, url); err != nil {
			return Entry{}, err
		}
		remote, err := git.ParseRemoteURL(url)
		if err != nil {
			return Entry{}, err
		}
		if _, err = probe(ctx, remote, rev); err == nil {
			e := c.Entry(url)
			e.Rev = rev
			return e, nil
		}
		if ctx.Err() != nil {
			return Entry{}, ctx.Err()
		}
		failures = append(failures, url+": "+err.Error())
	}
	return Entry{}, errors.New("cannot reach " + c.ID + " over SSH or HTTPS: " + strings.Join(failures, "; "))
}

// checkImplicit applies to a source offered for one run every rule a declared
// source must satisfy, and refuses one that shadows a declared or built-in
// source, so a one-run grant can never change what an existing identity means.
func checkImplicit(e Entry, declared File, seen map[string]bool) error {
	if err := Validate(e); err != nil {
		return err
	}
	if e.Type != source.KindGit {
		return fmt.Errorf("source %s: only a git repository can be offered for one run", e.ID)
	}
	if derived, onGitHub, err := githubIdentity(e.URL); err != nil || !onGitHub || derived != e.ID {
		return fmt.Errorf("source %s: only a GitHub repository can be offered for one run", e.ID)
	}
	if _, ok := Find(declared, e.ID); ok {
		return fmt.Errorf("source %s is already declared on this machine", e.ID)
	}
	if seen[e.ID] {
		return fmt.Errorf("source %s is offered twice", e.ID)
	}
	seen[e.ID] = true
	return nil
}

// sameRepository reports whether a lock's URL names the same repository as the
// one this machine reaches id at: equal, or both a GitHub URL for id, so one
// machine's SSH choice never breaks another's. A different host or repository
// is not followed.
func sameRepository(id, machineURL, lockURL string) bool {
	if machineURL == lockURL {
		return true
	}
	for _, u := range []string{machineURL, lockURL} {
		derived, onGitHub, err := githubIdentity(u)
		if err != nil || !onGitHub || derived != id {
			return false
		}
	}
	return true
}
