package resolve_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/resolve"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/source/sourcetest"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
)

const commit = "cccccccccccccccccccccccccccccccccccccccc"

func mkSkill(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, name), 0o750); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: " + name + "\nversion: 1.0.0\ndescription: x\n---\n"
	if err := os.WriteFile(filepath.Join(dir, name, "SKILL.md"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

type world struct {
	down   bool // a source that cannot be reached
	locked bool // a source that is not trusted
}

// chain holds: a local clone (only-local, both), the official source
// (only-official, both), @acme/skills (acme-only), and optionally a source that
// is down and one that is untrusted.
func (w world) chain(t *testing.T) source.Chain {
	t.Helper()
	clone, cache := t.TempDir(), t.TempDir()
	mkSkill(t, filepath.Join(clone, "skills"), "only-local")
	mkSkill(t, filepath.Join(clone, "skills"), "both")
	cfg := config.Config{RepoPath: clone, OfficialCache: filepath.Join(cache, "anthropic")}
	specs := append(cfg.SourceSpecs(),
		source.Spec{ID: "@acme/skills", Type: source.KindGit, URL: "git@github.com:acme/skills.git"})
	if w.down {
		specs = append(
			specs,
			source.Spec{ID: "@down/skills", Type: source.KindGit, URL: "git@github.com:down/skills.git"},
		)
	}
	if w.locked {
		specs = append(
			specs,
			source.Spec{ID: "@locked/skills", Type: source.KindGit, URL: "git@github.com:locked/skills.git"},
		)
	}
	opts := cfg.SourceOptions()
	opts.RepoFor = func(spec source.Spec) source.RepoOps {
		switch spec.ID {
		case skill.OfficialID:
			return &sourcetest.Repo{Commit: commit, Skills: map[string]string{"only-official": "", "both": ""}}
		case "@down/skills":
			return &sourcetest.Repo{FetchErr: errors.New("connection refused")}
		default:
			return &sourcetest.Repo{Commit: commit, Skills: map[string]string{"acme-only": "2.0.0"}}
		}
	}
	opts.Gate = func(id, url string) error {
		if id == "@locked/skills" {
			return &trust.UntrustedError{ID: id, URL: url}
		}
		return nil
	}
	chain, err := source.Build(specs, opts)
	if err != nil {
		t.Fatal(err)
	}
	return chain
}

func TestAQualifiedReferenceNamesItsSourcePathVersionAndCommit(t *testing.T) {
	t.Parallel()
	got, err := resolve.Reference(context.Background(), "@acme/skills/acme-only", world{}.chain(t), resolve.Options{})
	if err != nil {
		t.Fatal(err)
	}
	s := got.Skill
	if got.Status != resolve.Resolved || s == nil || s.Source != "@acme/skills" || s.Type != "git" ||
		s.Version != "2.0.0" || s.Path != "skills/acme-only" || s.Commit != commit || s.Tree == "" ||
		s.URL != "git@github.com:acme/skills.git" || s.Key != "@acme/skills/acme-only" {
		t.Errorf("result = %+v, skill = %+v", got, s)
	}
	if len(got.Consulted) != 0 {
		t.Errorf("a qualified reference consults one source, listed %+v", got.Consulted)
	}
}

func TestABareNameResolvesWhenOneSourceHasIt(t *testing.T) {
	t.Parallel()
	got, err := resolve.Reference(context.Background(), "only-official", world{}.chain(t), resolve.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != resolve.Resolved || got.Skill.Source != skill.OfficialID || got.Skill.Commit != commit {
		t.Errorf("result = %+v, skill = %+v", got, got.Skill)
	}
	var has []string
	for _, a := range got.Consulted {
		if a.Has {
			has = append(has, a.Source)
		}
	}
	if len(got.Consulted) != 3 || len(has) != 1 {
		t.Errorf("consulted = %+v", got.Consulted)
	}
}

func TestABareNameInTwoSourcesIsAmbiguousAndSaysHowToChoose(t *testing.T) {
	t.Parallel()
	got, err := resolve.Reference(context.Background(), "both", world{}.chain(t), resolve.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != resolve.Ambiguous || got.Skill != nil {
		t.Fatalf("result = %+v", got)
	}
	for _, want := range []string{skill.OfficialID, "@local/", "@anthropics/skills/both"} {
		if !strings.Contains(got.Message, want) {
			t.Errorf("message %q does not mention %q", got.Message, want)
		}
	}
}

func TestAnUnknownNameAndAnUnknownSourceAreNotFallbacks(t *testing.T) {
	t.Parallel()
	chain := world{}.chain(t)
	hint := func(id string) string { return "would derive git@github.com:" + strings.TrimPrefix(id, "@") + ".git" }

	got, err := resolve.Reference(context.Background(), "nope", chain, resolve.Options{})
	if err != nil || got.Status != resolve.NotFound {
		t.Errorf("bare unknown: %+v, %v", got, err)
	}
	got, err = resolve.Reference(context.Background(), "@acme/skills/absent", chain, resolve.Options{})
	if err != nil || got.Status != resolve.NotFound {
		t.Errorf("qualified absent: %+v, %v", got, err)
	}
	got, err = resolve.Reference(context.Background(), "@nobody/repo/x", chain, resolve.Options{ImplicitHint: hint})
	if err != nil || got.Status != resolve.UnknownSource ||
		!strings.Contains(got.Message, "git@github.com:nobody/repo.git") {
		t.Errorf("unknown source: %+v, %v", got, err)
	}
}

func TestASourceThatCannotBeConsultedFailsABareNameLikeInstallDoes(t *testing.T) {
	t.Parallel()
	for name, w := range map[string]world{"down": {down: true}, "untrusted": {locked: true}} {
		chain := w.chain(t)
		got, err := resolve.Reference(context.Background(), "only-local", chain, resolve.Options{})
		if err != nil || got.Status != resolve.Unavailable || got.Skill != nil {
			t.Errorf("%s: result = %+v, %v", name, got, err)
		}
		if _, findErr := chain.Find(context.Background(), skill.Ref{Name: "only-local"}); findErr == nil {
			t.Errorf("%s: install would resolve what resolve refuses", name)
		}
	}
	got, _ := resolve.Reference(context.Background(), "only-local", world{locked: true}.chain(t), resolve.Options{})
	if !strings.Contains(got.Message, "rsk source trust @locked/skills") {
		t.Errorf("an untrusted source does not say how to trust it: %q", got.Message)
	}

	// A qualified reference never needs the others.
	got, err := resolve.Reference(
		context.Background(),
		"@acme/skills/acme-only",
		world{down: true}.chain(t),
		resolve.Options{},
	)
	if err != nil || got.Status != resolve.Resolved {
		t.Errorf("qualified with another source down: %+v, %v", got, err)
	}
}

func TestResolveAgreesWithChainFindOnEveryOutcome(t *testing.T) {
	t.Parallel()
	chain := world{down: true, locked: true}.chain(t)
	for _, ref := range []string{"only-local", "both", "nope", "@acme/skills/acme-only", "@nobody/x/y", "@down/skills/x"} {
		got, err := resolve.Reference(context.Background(), ref, chain, resolve.Options{})
		if err != nil {
			t.Fatal(err)
		}
		parsed, _, _ := skill.ParseRef(ref)
		_, findErr := chain.Find(context.Background(), parsed)
		if (got.Status == resolve.Resolved) != (findErr == nil) {
			t.Errorf("%s: status %s but Find err = %v", ref, got.Status, findErr)
		}
		if got.Status == resolve.Ambiguous && !errors.Is(findErr, source.ErrAmbiguous) {
			t.Errorf("%s: ambiguous but Find err = %v", ref, findErr)
		}
	}
}

func TestABundleNameIsReportedAsABundle(t *testing.T) {
	t.Parallel()
	bundle := config.Bundle{Name: "starter", Skills: []config.SkillRef{{Name: "only-local"}, {Name: "both"}}}
	got, err := resolve.Reference(
		context.Background(),
		"starter",
		world{}.chain(t),
		resolve.Options{Bundles: []config.Bundle{bundle}},
	)
	if err != nil || got.Status != resolve.Bundle || strings.Join(got.Members, ",") != "only-local,both" {
		t.Errorf("result = %+v, %v", got, err)
	}
}

func TestResultIsPlainDataThatEncodesAsJSON(t *testing.T) {
	t.Parallel()
	got, err := resolve.Reference(
		context.Background(),
		"@acme/skills/acme-only@2.0.0",
		world{}.chain(t),
		resolve.Options{},
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	var back map[string]any
	if err = json.Unmarshal(
		data,
		&back,
	); err != nil || back["status"] != "resolved" ||
		back["requested_version"] != "2.0.0" {
		t.Errorf("json = %s, %v", data, err)
	}
}

func TestAMalformedReferenceIsAnError(t *testing.T) {
	t.Parallel()
	if _, err := resolve.Reference(context.Background(), "@bad", world{}.chain(t), resolve.Options{}); err == nil {
		t.Error("a malformed reference was accepted")
	}
}
