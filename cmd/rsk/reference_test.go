package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/resolve"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
)

// resolveMachine saves a config with a local clone holding "mine" and "both",
// the official source holding "theirs" and "both", and a declared path source
// holding "scratch-skill". It returns the config folder.
func resolveMachine(t *testing.T) string {
	t.Helper()
	dir := isolateConfig(t)
	withOfficialSkills(t, map[string]string{"theirs": "", "both": ""})
	clone, scratch := filepath.Join(dir, "clone"), filepath.Join(dir, "scratch")
	for _, name := range []string{"mine", "both"} {
		writeCatalogSkill(t, filepath.Join(clone, "skills"), name)
	}
	writeCatalogSkill(t, scratch, "scratch-skill")
	cfg := config.Config{
		RepoPath: clone, OfficialCache: filepath.Join(dir, "official"), VersionsCache: filepath.Join(dir, "versions"),
		GlobalTargets: map[string]string{"claude-code": filepath.Join(dir, "claude")}, DefaultTargetScope: "all",
	}
	if err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}
	entry, err := sources.NewEntry(scratch, sources.Options{ID: "@me/scratch"})
	if err != nil {
		t.Fatal(err)
	}
	file, err := sources.Add(sources.File{}, entry)
	if err != nil {
		t.Fatal(err)
	}
	if err = sources.Save(sources.DefaultPath(), file); err != nil {
		t.Fatal(err)
	}
	return dir
}

func runResolveCmd(t *testing.T, args ...string) (stdout string, err error) {
	t.Helper()
	cmd := newResolveCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(io.Discard)
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	cmd.SetArgs(args)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = cmd.ExecuteContext(ctx)
	return out.String(), err
}

//nolint:paralleltest // the config folder override and the clone seam are process-wide
func TestResolveExplainsAResolvedReferenceAsText(t *testing.T) {
	resolveMachine(t)
	out, err := runResolveCmd(t, "theirs")
	if err != nil {
		t.Fatalf("err = %v\n%s", err, out)
	}
	for _, want := range []string{"resolved", "@anthropics/skills", "skills/theirs", strings.Repeat("c", 40)} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

//nolint:paralleltest // the config folder override and the clone seam are process-wide
func TestResolveSaysWhyABareNameIsAmbiguousAndFailsTheCommand(t *testing.T) {
	resolveMachine(t)
	out, err := runResolveCmd(t, "both")
	if err == nil {
		t.Errorf("an ambiguous name exits 0:\n%s", out)
	}
	for _, want := range []string{"ambiguous", "@anthropics/skills", "@anthropics/skills/both"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

//nolint:paralleltest // the config folder override and the clone seam are process-wide
func TestResolveJSONCarriesTheResultEvenWhenItIsNotResolved(t *testing.T) {
	resolveMachine(t)
	for ref, want := range map[string]resolve.Status{
		"@me/scratch/scratch-skill": resolve.Resolved,
		"both":                      resolve.Ambiguous,
		"nope":                      resolve.NotFound,
		"@nobody/repo/x":            resolve.UnknownSource,
	} {
		out, err := runResolveCmd(t, ref, "-o", "json")
		if (want == resolve.Resolved) != (err == nil) {
			t.Errorf("%s: err = %v", ref, err)
		}
		var res resolve.Result
		if jsonErr := json.Unmarshal([]byte(out), &res); jsonErr != nil || res.Status != want {
			t.Errorf("%s: status %q, want %q (%v)\n%s", ref, res.Status, want, jsonErr, out)
		}
	}
}

//nolint:paralleltest // the config folder override and the clone seam are process-wide
func TestResolveSaysWhatInstallWouldDoForAnUnregisteredSourceWithoutContactingIt(t *testing.T) {
	resolveMachine(t)
	contacted := false
	implicitProbe = func(context.Context, git.RemoteURL, string) (string, error) { contacted = true; return "", nil }
	t.Cleanup(func() { implicitProbe = defaultImplicitProbe })

	out, err := runResolveCmd(t, "@nobody/repo/x")
	if err == nil || contacted {
		t.Errorf("err = %v, contacted = %v", err, contacted)
	}
	for _, want := range []string{"unknown-source", "git@github.com:nobody/repo.git", "https://github.com/nobody/repo.git", "this run only"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}
}

//nolint:paralleltest // the config folder override and the clone seam are process-wide
func TestResolveChangesNoDeclarationAndGrantsNoTrust(t *testing.T) {
	dir := resolveMachine(t)
	read := func(name string) string {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	configBefore, sourcesBefore := read(config.ConfigFileName), read(sources.FileName)
	for _, ref := range []string{"theirs", "both", "@me/scratch/scratch-skill", "@nobody/repo/x"} {
		_, _ = runResolveCmd(t, ref)
	}
	if read(config.ConfigFileName) != configBefore || read(sources.FileName) != sourcesBefore {
		t.Error("resolve rewrote the machine's configuration")
	}
	if _, err := os.Stat(filepath.Join(dir, "trust.json")); err == nil {
		t.Error("resolve wrote trust.json")
	}
}

//nolint:paralleltest // the config folder override and the clone seam are process-wide
func TestResolveRejectsAMalformedReference(t *testing.T) {
	resolveMachine(t)
	if _, err := runResolveCmd(t, "@bad"); err == nil {
		t.Error("a malformed reference exited 0")
	}
}
