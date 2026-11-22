package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.ralva.dev/lib/term"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source/sourcetest"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

const implicitCommit = "dddddddddddddddddddddddddddddddddddddddd"

type implicitWorld struct {
	dir    string
	global string
	repo   *sourcetest.Repo
	probed []string
	sshErr error // when set, the SSH probe fails with it
}

// implicitMachine isolates the config folder and points every seam at in-memory
// stand-ins: the clone every git source reads, the probe that picks a transport,
// and the terminal check. The seams are process-wide, so a test using it is
// sequential.
func implicitMachine(t *testing.T) *implicitWorld {
	t.Helper()
	dir := isolateConfig(t)
	w := &implicitWorld{dir: dir, global: filepath.Join(dir, "global-claude")}
	w.repo = withOfficialSkills(t, map[string]string{"alpha": "1.0.0"})
	w.repo.Commit = implicitCommit

	clone := filepath.Join(dir, "clone")
	writeCatalogSkill(t, filepath.Join(clone, "skills"), "local-one")
	cfg := config.Config{
		RepoPath: clone, OfficialCache: filepath.Join(dir, "official"), VersionsCache: filepath.Join(dir, "versions"),
		GlobalTargets: map[string]string{"claude-code": w.global}, DefaultTargetScope: "all",
	}
	if err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}

	implicitProbe = func(_ context.Context, remote git.RemoteURL, _ string) (string, error) {
		w.probed = append(w.probed, remote.String())
		if w.sshErr != nil && strings.HasPrefix(remote.String(), "git@") {
			return "", w.sshErr
		}
		return implicitCommit, nil
	}
	t.Cleanup(func() { implicitProbe, isInteractive = defaultImplicitProbe, term.StdioIsTerminal })
	isInteractive = func() bool { return false }
	return w
}

func (w *implicitWorld) install(t *testing.T, stdin string, opts installOpts, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	installCmd.SetOut(&out)
	installCmd.SetErr(&out)
	installCmd.SetIn(strings.NewReader(stdin))
	installCmd.SetContext(context.Background())
	t.Cleanup(func() { installCmd.SetOut(nil); installCmd.SetErr(nil); installCmd.SetIn(nil) })
	opts.global = true
	err := runInstall(installCmd, opts, args)
	return out.String(), err
}

func (w *implicitWorld) wrote(name string) bool {
	_, err := os.Stat(filepath.Join(w.dir, name))
	return err == nil
}

func (w *implicitWorld) linked() bool {
	_, err := os.Lstat(filepath.Join(w.global, "alpha"))
	return err == nil
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestAnImplicitSourceIsRefusedWithoutATerminalEvenWithYes(t *testing.T) {
	w := implicitMachine(t)
	out, err := w.install(t, "y\n", installOpts{yes: true}, "@acme/skills/alpha")
	if err == nil || !strings.Contains(err.Error(), "--trust @acme/skills") {
		t.Fatalf("err = %v\n%s", err, out)
	}
	if !strings.Contains(err.Error(), "--yes does not") {
		t.Errorf("the refusal does not say --yes is not enough: %v", err)
	}
	if len(w.probed) != 0 || len(w.repo.Fetches) != 0 {
		t.Errorf("probed %v, fetched %v: nothing may be contacted before approval", w.probed, w.repo.Fetches)
	}
	if w.linked() || w.wrote(sources.FileName) || w.wrote("trust.json") {
		t.Error("a refused implicit source left state behind")
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestAnImplicitSourceIsRefusedInsideTheSession(t *testing.T) {
	w := implicitMachine(t)
	isInteractive = func() bool { return true }
	var out bytes.Buffer
	installCmd.SetOut(&out)
	installCmd.SetIn(strings.NewReader("y\n"))
	installCmd.SetContext(term.WithCapture(context.Background(), &term.Capture{}))
	t.Cleanup(func() { installCmd.SetOut(nil); installCmd.SetIn(nil); installCmd.SetContext(context.Background()) })

	err := runInstall(installCmd, installOpts{global: true, yes: true}, []string{"@acme/skills/alpha"})
	if err == nil || len(w.probed) != 0 {
		t.Errorf("err = %v, probed %v: the session must not approve an implicit source", err, w.probed)
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestTrustFlagApprovesOneRunAndWritesNothing(t *testing.T) {
	w := implicitMachine(t)
	out, err := w.install(t, "", installOpts{trust: []string{"@acme/skills"}}, "@acme/skills/alpha")
	if err != nil {
		t.Fatalf("err = %v\n%s", err, out)
	}
	if !w.linked() {
		t.Error("the skill was not installed")
	}
	if len(w.probed) != 1 || w.probed[0] != "git@github.com:acme/skills.git" {
		t.Errorf("probed %v, want SSH only", w.probed)
	}
	if w.wrote(sources.FileName) || w.wrote("trust.json") {
		t.Error("a one-run trust wrote sources.json or trust.json")
	}

	// The next run, without the flag, is refused again.
	if _, err = w.install(t, "", installOpts{}, "@acme/skills/alpha"); err == nil {
		t.Error("the source stayed trusted after the run")
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestTrustFlagMustNameASourceTheCommandUsesAndNotADeclaredOne(t *testing.T) {
	w := implicitMachine(t)
	if _, err := w.install(t, "", installOpts{trust: []string{"@other/repo"}}, "@acme/skills/alpha"); err == nil ||
		len(w.probed) != 0 {
		t.Errorf("err = %v, probed %v: --trust for a source nothing uses was accepted", err, w.probed)
	}

	entry, err := sources.NewEntry("@acme/skills", sources.Options{})
	if err != nil {
		t.Fatal(err)
	}
	file, _ := sources.Add(sources.File{}, entry)
	if err = sources.Save(sources.DefaultPath(), file); err != nil {
		t.Fatal(err)
	}
	_, err = w.install(t, "", installOpts{trust: []string{"@acme/skills"}}, "@acme/skills/alpha")
	if err == nil || !strings.Contains(err.Error(), "rsk source trust @acme/skills") || len(w.probed) != 0 {
		t.Errorf("err = %v, probed %v: --trust must not stand in for trusting a declared source", err, w.probed)
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestAnImplicitSourceAsksAtATerminalAndYesDoesNotAnswerForTheUser(t *testing.T) {
	w := implicitMachine(t)
	isInteractive = func() bool { return true }

	out, err := w.install(t, "n\n", installOpts{yes: true}, "@acme/skills/alpha")
	if err != nil || w.linked() || len(w.probed) != 0 {
		t.Fatalf(
			"a refusal at the prompt: err = %v, linked = %v, probed = %v\n%s",
			err,
			w.linked(),
			w.probed,
			out,
		)
	}
	for _, want := range []string{"@acme/skills", "git@github.com:acme/skills.git", "https://github.com/acme/skills.git"} {
		if !strings.Contains(out, want) {
			t.Errorf("the prompt does not show %q:\n%s", want, out)
		}
	}

	if out, err = w.install(t, "y\n", installOpts{yes: true}, "@acme/skills/alpha"); err != nil || !w.linked() {
		t.Fatalf("approving: err = %v, linked = %v\n%s", err, w.linked(), out)
	}
	if w.wrote(sources.FileName) || w.wrote("trust.json") {
		t.Error("approving for one run wrote state")
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestSaveRegistersTheSourceUntrustedAndNeverGrantsTrust(t *testing.T) {
	w := implicitMachine(t)
	out, err := w.install(t, "", installOpts{trust: []string{"@acme/skills"}, save: true}, "@acme/skills/alpha")
	if err != nil {
		t.Fatalf("err = %v\n%s", err, out)
	}
	file, err := sources.Load(sources.DefaultPath())
	if err != nil || len(file.Sources) != 1 || file.Sources[0].ID != "@acme/skills" ||
		file.Sources[0].URL != "git@github.com:acme/skills.git" {
		t.Fatalf("sources = %+v, %v", file, err)
	}
	if w.wrote("trust.json") {
		t.Error("--save with --trust granted trust: only a person at a terminal may")
	}
	for _, want := range []string{"rsk source trust @acme/skills", "bare skill name fails", "rsk source remove @acme/skills"} {
		if !strings.Contains(out, want) {
			t.Errorf("--save does not say %q:\n%s", want, out)
		}
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestHTTPSIsTheFallbackAndSaveKeepsWhatWorked(t *testing.T) {
	w := implicitMachine(t)
	w.sshErr = errors.New("Permission denied (publickey)")
	if out, err := w.install(
		t,
		"",
		installOpts{trust: []string{"@acme/skills"}, save: true},
		"@acme/skills/alpha",
	); err != nil {
		t.Fatalf("err = %v\n%s", err, out)
	}
	if len(w.probed) != 2 || !strings.HasPrefix(w.probed[1], "https://") {
		t.Errorf("probed %v, want SSH then HTTPS", w.probed)
	}
	file, _ := sources.Load(sources.DefaultPath())
	if len(file.Sources) != 1 || file.Sources[0].URL != "https://github.com/acme/skills.git" {
		t.Errorf("sources = %+v, want the URL that worked", file.Sources)
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestAnImplicitSourceNeedsTheNetworkSoOfflineRefusesIt(t *testing.T) {
	w := implicitMachine(t)
	if installCmd.Flags().Lookup(flagOffline) == nil {
		installCmd.Flags().Bool(flagOffline, false, "")
	}
	_ = installCmd.Flags().Set(flagOffline, "true")
	t.Cleanup(func() { _ = installCmd.Flags().Set(flagOffline, "false") })

	if _, err := w.install(t, "", installOpts{trust: []string{"@acme/skills"}}, "@acme/skills/alpha"); err == nil ||
		len(w.probed) != 0 {
		t.Errorf("err = %v, probed %v: --offline must not contact an implicit source", err, w.probed)
	}
}

func TestAProjectInstallLocksTheIdentityAndCommitAndKeepsNoTrust(t *testing.T) {
	w := implicitMachine(t)
	project := t.TempDir()
	if err := manifest.WriteMod(
		filepath.Join(project, ".rsk"),
		manifest.Mod{Tools: []tool.ID{tool.ClaudeID}},
	); err != nil {
		t.Fatal(err)
	}
	t.Chdir(project)

	var out bytes.Buffer
	installCmd.SetOut(&out)
	installCmd.SetErr(&out)
	installCmd.SetContext(context.Background())
	t.Cleanup(func() { installCmd.SetOut(nil); installCmd.SetErr(nil) })
	err := runInstall(
		installCmd,
		installOpts{yes: true, trust: []string{"@acme/skills"}},
		[]string{"@acme/skills/alpha"},
	)
	if err != nil {
		t.Fatalf("err = %v\n%s", err, out.String())
	}
	lock, err := manifest.ReadLock(filepath.Join(project, ".rsk"))
	if err != nil || len(lock.Skills) != 1 {
		t.Fatalf("lock = %+v, %v", lock, err)
	}
	e := lock.Skills[0]
	if e.Source != skill.DeclaredSource("@acme/skills") || e.Git == nil || e.Git.Commit != implicitCommit {
		t.Errorf("entry = %+v", e)
	}
	if w.wrote("trust.json") {
		t.Error("the project install granted trust")
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestSyncRestoresAnUndeclaredSourceOnlyWithTheRunsApproval(t *testing.T) {
	w := implicitMachine(t)
	project := t.TempDir()
	m := manifest.Mod{Tools: []tool.ID{tool.ClaudeID}, Skills: map[string]string{"alpha": "*"}}
	lock := manifest.Lock{Skills: []manifest.LockEntry{{
		Name: "alpha", Version: "1.0.0", Source: skill.DeclaredSource("@acme/skills"),
		Git: &skill.GitRef{
			URL: "https://github.com/acme/skills.git", Path: "skills/alpha",
			Commit: implicitCommit, Tree: sourcetest.TreeOf("skills/alpha"),
		},
	}}}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	syncCmd.SetOut(&out)
	syncCmd.SetErr(&out)
	syncCmd.SetContext(context.Background())
	t.Cleanup(func() { syncCmd.SetOut(nil); syncCmd.SetErr(nil) })

	if _, err = approveLockSources(syncCmd, lock, nil); err == nil || len(w.probed) != 0 {
		t.Fatalf("err = %v, probed %v: a lock must not open an undeclared source on its own", err, w.probed)
	}
	approved, err := approveLockSources(syncCmd, lock, []string{"@acme/skills"})
	if err != nil || len(approved) != 1 {
		t.Fatalf("approved = %+v, %v", approved, err)
	}
	links, err := planSyncLinks(t.Context(), project, m, lock, cfg, approved...)
	if err != nil || len(links) != 1 || links[0].skill.Git.URL != "git@github.com:acme/skills.git" {
		t.Fatalf("links = %+v, %v: the machine's own transport must be used", links, err)
	}
	if len(w.repo.Fetches) != 0 {
		t.Errorf("restore fetched %v: it must not consult the branch", w.repo.Fetches)
	}
}
