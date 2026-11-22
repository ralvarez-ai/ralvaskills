// Package sourceops is the one implementation of what `rsk source` and the MCP
// source tools do: list, add, remove, test and scan sources. It takes and
// returns plain structs, so the two surfaces only parse input and render
// output and cannot disagree. It never grants trust: that is `rsk source trust`
// in a terminal, and nothing here writes trust.json except to withdraw a grant
// when a source is removed.
package sourceops

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
)

const (
	// TrustTrusted says a grant covers the source at its URL.
	TrustTrusted = "trusted"

	// TrustNone says a git source has no grant.
	TrustNone = "untrusted"

	// TrustNA says the source needs no grant: a directory contacts nothing.
	TrustNA = "not required"

	// TrustBuiltIn says rsk ships the source and trusts its own URL.
	TrustBuiltIn = "built-in"

	// ActionAdd names the operation that registers a source.
	ActionAdd = "add"

	// ActionRemove names the operation that drops one.
	ActionRemove = "remove"
)

type (
	// Env is where the operations read and write. Zero paths mean the machine's
	// own files; the function fields exist so tests reach no network.
	Env struct {
		// Cfg, when set, adds the built-in sources to List and supplies the cache
		// locations a scan needs.
		Cfg *config.Config

		SourcesPath string
		TrustPath   string

		// RepoFor supplies the clone a git source reads. Nil means the real one.
		RepoFor func(source.Spec) source.RepoOps

		// LsRemote asks a remote for its branch head. Nil means git.LsRemote.
		LsRemote sources.Probe

		// EffectiveURL says where git really goes for a URL. Nil means
		// git.EffectiveURL.
		EffectiveURL func(context.Context, git.RemoteURL) (string, error)
	}

	// View is one source as it is listed, and as JSON.
	View struct {
		ID       string   `json:"id"`
		Type     string   `json:"type"`
		URL      string   `json:"url"`
		Rev      string   `json:"rev,omitempty"`
		Roots    []string `json:"roots,omitempty"`
		Excludes []string `json:"excludes,omitempty"`
		Trust    string   `json:"trust"`
	}

	// Change is the preview and the result of a mutating operation: the same JSON
	// either way, so an agent shows it and then applies it.
	Change struct {
		Action  string `json:"action"`
		DryRun  bool   `json:"dry_run"`
		Applied bool   `json:"applied"`
		Source  View   `json:"source"`
		Next    string `json:"next,omitempty"`
	}

	// TestResult is the outcome of testing one source.
	TestResult struct {
		View
		Reachable bool   `json:"reachable"`
		Via       string `json:"via,omitempty"`
		Commit    string `json:"commit,omitempty"`
		Error     string `json:"error,omitempty"`
	}

	// ScanOutcome is a scan together with the source it read.
	ScanOutcome struct {
		Source View `json:"source"`
		source.ScanResult
	}

	state struct {
		file   sources.File
		grants trust.File
	}
)

func (e Env) sourcesPath() string {
	if e.SourcesPath != "" {
		return e.SourcesPath
	}
	return sources.DefaultPath()
}

func (e Env) trustPath() string {
	if e.TrustPath != "" {
		return e.TrustPath
	}
	return trust.DefaultPath()
}

func (e Env) load() (state, error) {
	file, err := sources.Load(e.sourcesPath())
	if err != nil {
		return state{}, err
	}
	grants, err := trust.Load(e.trustPath())
	if err != nil {
		return state{}, err
	}
	return state{file: file, grants: grants}, nil
}

// TrustHint is the command that trusts a source: only a person at a terminal
// can run it.
func TrustHint(id string) string {
	return "Trust it yourself in a terminal before use: rsk source trust " + id
}

// ViewOf describes a declaration and whether it is trusted.
func ViewOf(e sources.Entry, grants trust.File) View {
	state := TrustNone
	switch {
	case e.Type == source.KindPath:
		state = TrustNA
	case grants.Granted(e.ID, e.URL):
		state = TrustTrusted
	}
	return View{
		ID: e.ID, Type: string(e.Type), URL: e.URL, Rev: e.Rev, Roots: e.Roots, Excludes: e.Excludes, Trust: state,
	}
}

// List returns the built-in sources (when Env.Cfg is set) and then the declared
// ones, each with its trust state.
func List(env Env) ([]View, error) {
	st, err := env.load()
	if err != nil {
		return nil, err
	}
	var views []View
	if env.Cfg != nil {
		views = append(
			views,
			View{ID: skill.DefaultRegistryID, Type: string(source.KindRegistry), URL: env.Cfg.EffectiveRegistryURL(),
				Trust: TrustBuiltIn},
			View{
				ID:    skill.OfficialID,
				Type:  string(source.KindGit),
				URL:   source.OfficialSkillsURL,
				Trust: TrustBuiltIn,
			},
		)
	}
	for i := range st.file.Sources {
		views = append(views, ViewOf(st.file.Sources[i], st.grants))
	}
	return views, nil
}

// PlanAdd validates a declaration and returns its preview and the function that
// writes it. The source is registered untrusted; nothing is contacted.
func PlanAdd(env Env, target string, o sources.Options) (Change, func() error, error) {
	entry, err := sources.NewEntry(target, o)
	if err != nil {
		return Change{}, nil, err
	}
	st, err := env.load()
	if err != nil {
		return Change{}, nil, err
	}
	next, err := sources.Add(st.file, entry)
	if err != nil {
		return Change{}, nil, err
	}
	change := Change{Action: ActionAdd, Source: ViewOf(entry, st.grants)}
	if entry.Type == source.KindGit {
		change.Next = TrustHint(entry.ID)
	}
	return change, func() error { return sources.Save(env.sourcesPath(), next) }, nil
}

// PlanRemove previews dropping a declared source and returns the function that
// does it, withdrawing the source's trust grant with it.
func PlanRemove(env Env, id string) (Change, func() error, error) {
	st, err := env.load()
	if err != nil {
		return Change{}, nil, err
	}
	entry, ok := sources.Find(st.file, id)
	if !ok {
		return Change{}, nil, fmt.Errorf("source %s is not declared", id)
	}
	next, _ := sources.Remove(st.file, entry.ID)
	change := Change{Action: ActionRemove, Source: ViewOf(entry, st.grants)}
	return change, func() error {
		if saveErr := sources.Save(env.sourcesPath(), next); saveErr != nil {
			return saveErr
		}
		grants, revoked := trust.Revoke(st.grants, entry.ID)
		if !revoked {
			return nil
		}
		return trust.Save(env.trustPath(), grants)
	}, nil
}

// Test checks that a declared source can be used. A git source must be trusted
// before anything contacts its host, so an untrusted one fails here with the
// command to run and nothing is sent. The result is filled in even on failure.
func Test(ctx context.Context, env Env, id string) (TestResult, error) {
	st, err := env.load()
	if err != nil {
		return TestResult{}, err
	}
	entry, ok := sources.Find(st.file, id)
	if !ok {
		return TestResult{}, fmt.Errorf("source %s is not declared", id)
	}
	res := TestResult{View: ViewOf(entry, st.grants)}
	if err = env.probe(ctx, entry, st.grants, &res); err != nil {
		res.Error = err.Error()
		return res, err
	}
	return res, nil
}

func (e Env) probe(ctx context.Context, entry sources.Entry, grants trust.File, res *TestResult) error {
	if entry.Type == source.KindPath {
		res.Reachable = true
		return nil
	}
	if err := trust.Require(grants, entry.ID, entry.URL); err != nil {
		return err
	}
	remote, err := git.ParseRemoteURL(entry.URL)
	if err != nil {
		return err
	}
	effective := e.EffectiveURL
	if effective == nil {
		effective = git.EffectiveURL
	}
	if actual, urlErr := effective(ctx, remote); urlErr == nil && actual != "" && actual != entry.URL {
		res.Via = actual
	}
	lsRemote := e.LsRemote
	if lsRemote == nil {
		lsRemote = git.LsRemote
		have, versionErr := git.Version(ctx)
		if versionErr != nil {
			return versionErr
		}
		if !git.MeetsMinimum(have) {
			return fmt.Errorf("git %d.%d.%d is older than the %d.%d.%d rsk needs", have[0], have[1], have[2],
				git.MinVersion[0], git.MinVersion[1], git.MinVersion[2])
		}
	}
	commit, err := lsRemote(ctx, remote, entry.Rev)
	if err != nil {
		return err
	}
	res.Reachable, res.Commit = true, commit
	return nil
}

// scanner is implemented by a source that can scan itself.
type scanner interface {
	Scan(ctx context.Context) (source.ScanResult, error)
}

// Scan discovers the skills in a directory or in a declared source without
// installing or registering anything. A directory needs no trust; a declared git
// source goes through the trust gate like any other read; an @owner/repo that is
// not declared is refused, since scanning it would clone code nobody approved.
func Scan(ctx context.Context, env Env, target string) (ScanOutcome, error) {
	st, err := env.load()
	if err != nil {
		return ScanOutcome{}, err
	}
	if strings.HasPrefix(target, "@") {
		entry, ok := sources.Find(st.file, target)
		if !ok {
			return ScanOutcome{}, fmt.Errorf("source %s is not declared: register it first with 'rsk source add %s'",
				target, target)
		}
		return env.scanEntry(ctx, entry, st)
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return ScanOutcome{}, fmt.Errorf("resolve %s: %w", target, err)
	}
	if info, statErr := os.Stat(abs); statErr != nil || !info.IsDir() {
		return ScanOutcome{}, fmt.Errorf(
			"%s is not a directory (a repository is scanned by its identity once registered)",
			target,
		)
	}
	entry := sources.Entry{ID: skill.LocalScope + "/scan", Type: source.KindPath, URL: abs}
	return env.scanEntry(ctx, entry, st)
}

func (e Env) scanEntry(ctx context.Context, entry sources.Entry, st state) (ScanOutcome, error) {
	chain, err := e.chainFor(entry, st)
	if err != nil {
		return ScanOutcome{}, err
	}
	if len(chain) == 0 {
		return ScanOutcome{}, errors.New("source " + entry.ID + " could not be built")
	}
	sc, ok := chain[0].(scanner)
	if !ok {
		return ScanOutcome{}, fmt.Errorf("source %s cannot be scanned", entry.ID)
	}
	res, err := sc.Scan(ctx)
	if err != nil {
		return ScanOutcome{}, err
	}
	return ScanOutcome{Source: ViewOf(entry, st.grants), ScanResult: res}, nil
}

// chainFor builds a chain of the one source entry, gated by the grants on file.
func (e Env) chainFor(entry sources.Entry, st state) (source.Chain, error) {
	cfg := config.Config{}
	if e.Cfg != nil {
		cfg = *e.Cfg
	}
	opts := cfg.SourceOptions()
	opts.RepoFor = e.RepoFor
	opts.Gate = func(id, url string) error {
		return trust.Require(st.grants, id, url)
	}
	spec := entry.Spec()
	r, err := source.New(spec, opts)
	if err != nil {
		return nil, err
	}
	return source.Chain{r}, nil
}
