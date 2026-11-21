package sources_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
)

// machine is a config folder with a local clone holding "mine", an official
// cache holding "theirs", and whatever declared sources a test adds.
type machine struct {
	cfg       config.Config
	sources   string
	trustFile string
}

func newMachine(t *testing.T) machine {
	t.Helper()
	dir := t.TempDir()
	repo, cache := filepath.Join(dir, "repo"), filepath.Join(dir, "cache", "anthropic")
	for root, name := range map[string]string{repo: "mine", cache: "theirs"} {
		skillDir := filepath.Join(root, "skills", name)
		if err := os.MkdirAll(skillDir, 0o755); err != nil {
			t.Fatal(err)
		}
		body := "---\nname: " + name + "\nversion: 1.0.0\ndescription: x\n---\n"
		if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return machine{
		cfg:       config.Config{RepoPath: repo, OfficialCache: cache},
		sources:   filepath.Join(dir, "sources.json"),
		trustFile: filepath.Join(dir, "trust.json"),
	}
}

func (m machine) declare(t *testing.T, target string, o sources.Options) sources.Entry {
	t.Helper()
	e, err := sources.NewEntry(target, o)
	if err != nil {
		t.Fatal(err)
	}
	f, err := sources.Load(m.sources)
	if err != nil {
		t.Fatal(err)
	}
	if f, err = sources.Add(f, e); err != nil {
		t.Fatal(err)
	}
	if err = sources.Save(m.sources, f); err != nil {
		t.Fatal(err)
	}
	return e
}

func (m machine) grant(t *testing.T, e sources.Entry) {
	t.Helper()
	f, err := trust.Load(m.trustFile)
	if err != nil {
		t.Fatal(err)
	}
	if err = trust.Save(m.trustFile, trust.Grant(f, e.ID, e.URL, time.Now())); err != nil {
		t.Fatal(err)
	}
}

func (m machine) chain(t *testing.T, offline bool) source.Chain {
	t.Helper()
	chain, err := sources.ChainAt(m.cfg, m.sources, m.trustFile, sources.ChainOptions{Offline: offline})
	if err != nil {
		t.Fatal(err)
	}
	return chain
}

func ids(chain source.Chain) []string {
	out := make([]string, len(chain))
	for i, r := range chain {
		out[i] = r.Spec().ID
	}
	return out
}

func TestChainAddsDeclaredGitSourcesAfterTheBuiltInOnes(t *testing.T) {
	t.Parallel()
	m := newMachine(t)
	m.declare(t, "@acme/skills", sources.Options{})
	m.declare(t, "@acme/other", sources.Options{})
	m.declare(t, t.TempDir(), sources.Options{ID: "@me/scratch"}) // a path source: not resolvable yet

	got := ids(m.chain(t, false))
	want := []string{"@local/repo", skill.OfficialID, "@acme/skills", "@acme/other"}
	if len(got) != len(want) {
		t.Fatalf("chain = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("chain = %v, want %v", got, want)
		}
	}
}

func TestChainWithoutDeclaredSourcesIsTheBuiltInChain(t *testing.T) {
	t.Parallel()
	m := newMachine(t)
	if got := ids(m.chain(t, false)); len(got) != 2 {
		t.Errorf("chain = %v, want the two built-in sources", got)
	}
}

func TestAnUntrustedSourceFailsABareLookupAndNothingElse(t *testing.T) {
	t.Parallel()
	m := newMachine(t)
	e := m.declare(t, "@acme/skills", sources.Options{})
	chain := m.chain(t, false)

	// A bare name consults every source, so the untrusted one stops it: an
	// unconsulted source must never be skipped quietly.
	_, err := chain.Find(t.Context(), skill.Ref{Name: "mine"})
	var untrusted *trust.UntrustedError
	if !errors.As(err, &untrusted) || untrusted.ID != e.ID {
		t.Fatalf("bare lookup err = %v, want an UntrustedError for %s", err, e.ID)
	}

	// A qualified reference consults only the source it names.
	got, err := chain.Find(t.Context(), skill.Ref{Source: skill.OfficialID, Name: "theirs"})
	if err != nil || got.Name != "theirs" {
		t.Errorf("a qualified lookup in a built-in source = %+v, %v", got, err)
	}
	if _, err = chain.Find(t.Context(), skill.Ref{Source: e.ID, Name: "x"}); !errors.As(err, &untrusted) {
		t.Errorf("a qualified lookup in the untrusted source err = %v", err)
	}
}

func TestATrustedSourceThatWasNeverFetchedFailsOffline(t *testing.T) {
	t.Parallel()
	m := newMachine(t)
	e := m.declare(t, "@acme/skills", sources.Options{})
	m.grant(t, e)

	_, err := m.chain(t, true).Find(t.Context(), skill.Ref{Name: "mine"})
	if !errors.Is(err, source.ErrNeverFetched) {
		t.Errorf("err = %v, want ErrNeverFetched: offline must not reach the network or skip the source", err)
	}
}

func TestAGrantForAnotherURLDoesNotTrustTheSource(t *testing.T) {
	t.Parallel()
	m := newMachine(t)
	e := m.declare(t, "@acme/skills", sources.Options{})
	f, _ := trust.Load(m.trustFile)
	if err := trust.Save(
		m.trustFile,
		trust.Grant(f, e.ID, "git@evil.example:acme/skills.git", time.Now()),
	); err != nil {
		t.Fatal(err)
	}
	if _, err := m.chain(t, true).Find(t.Context(), skill.Ref{Name: "mine"}); !errors.Is(err, trust.ErrUntrusted) {
		t.Errorf("err = %v, want ErrUntrusted for a grant made at a different URL", err)
	}
}
