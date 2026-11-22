package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"go.ralva.dev/ralvaskills/v3/internal/cmdx"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

func runFrom(t *testing.T, opts installOpts, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	installCmd.SetOut(&buf)
	installCmd.SetErr(&buf)
	installCmd.SetContext(context.Background())
	t.Cleanup(func() { installCmd.SetOut(nil); installCmd.SetErr(nil) })
	err := runInstall(installCmd, opts, args)
	return buf.String(), err
}

func TestFromLinksOneDirectoryIntoTheProjectAndRecordsNothing(t *testing.T) {
	implicitMachine(t)
	scratch := t.TempDir()
	writeCatalogSkill(t, scratch, "scratch-skill")
	project := t.TempDir()
	rsk := filepath.Join(project, ".rsk")
	if err := manifest.WriteMod(rsk, manifest.Mod{Tools: []tool.ID{tool.ClaudeID}}); err != nil {
		t.Fatal(err)
	}
	t.Chdir(project)

	out, err := runFrom(t, installOpts{from: filepath.Join(scratch, "scratch-skill"), yes: true})
	if err != nil {
		t.Fatalf("err = %v\n%s", err, out)
	}
	link := filepath.Join(project, ".claude", "skills", "scratch-skill")
	if info, statErr := os.Lstat(link); statErr != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("the skill is not linked into the project: %v", statErr)
	}
	m, err := manifest.ReadMod(rsk)
	if err != nil || len(m.Skills) != 0 || len(m.Pinned) != 0 {
		t.Errorf("manifest = %+v, %v: --from must record nothing", m, err)
	}
	lock, err := manifest.ReadLock(rsk)
	if err != nil || len(lock.Skills) != 0 {
		t.Errorf("lock = %+v, %v: --from must leave no lock record of a machine path", lock, err)
	}
	for _, want := range []string{"rsk.mod", "rsk.lock", "another machine"} {
		if !strings.Contains(out, want) {
			t.Errorf("output does not say the install is not reproducible (%q):\n%s", want, out)
		}
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestFromGlobalLinksIntoTheGlobalDirectory(t *testing.T) {
	w := implicitMachine(t)
	scratch := t.TempDir()
	writeCatalogSkill(t, scratch, "scratch-skill")

	out, err := runFrom(t, installOpts{from: filepath.Join(scratch, "scratch-skill"), global: true, yes: true})
	if err != nil {
		t.Fatalf("err = %v\n%s", err, out)
	}
	if _, statErr := os.Lstat(filepath.Join(w.global, "scratch-skill")); statErr != nil {
		t.Errorf("not linked globally: %v", statErr)
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestFromDryRunLinksNothing(t *testing.T) {
	w := implicitMachine(t)
	scratch := t.TempDir()
	writeCatalogSkill(t, scratch, "scratch-skill")
	out, err := runFrom(t, installOpts{from: filepath.Join(scratch, "scratch-skill"), global: true, dryRun: true})
	if err != nil || !strings.Contains(out, "scratch-skill") {
		t.Fatalf("err = %v\n%s", err, out)
	}
	if _, statErr := os.Lstat(filepath.Join(w.global, "scratch-skill")); statErr == nil {
		t.Error("a dry run linked the skill")
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestFromRefusesWhatItCannotRecordOrTrust(t *testing.T) {
	w := implicitMachine(t)
	scratch := t.TempDir()
	writeCatalogSkill(t, scratch, "scratch-skill")
	good := filepath.Join(scratch, "scratch-skill")
	if err := os.MkdirAll(filepath.Join(scratch, "pile"), 0o750); err != nil {
		t.Fatal(err)
	}
	writeCatalogSkill(t, filepath.Join(scratch, "pile"), "inner")

	for name, tc := range map[string]struct {
		opts installOpts
		args []string
	}{
		"with names":            {installOpts{from: good, global: true}, []string{"go-grpc"}},
		"with --pin":            {installOpts{from: good, pin: true}, nil},
		"with --save":           {installOpts{from: good, global: true, save: true}, nil},
		"with --trust":          {installOpts{from: good, global: true, trust: []string{"@a/b"}}, nil},
		"a directory of skills": {installOpts{from: filepath.Join(scratch, "pile"), global: true}, nil},
		"a missing directory":   {installOpts{from: filepath.Join(scratch, "absent"), global: true}, nil},
	} {
		if out, err := runFrom(t, tc.opts, tc.args...); err == nil {
			t.Errorf("%s: accepted\n%s", name, out)
		}
	}
	if _, err := os.Lstat(filepath.Join(w.global, "scratch-skill")); err == nil {
		t.Error("a refused --from linked the skill")
	}
}

// RunE reads these flags through helpers that panic on one nobody registered,
// and the tests above call runInstall directly, so this is what keeps a flag
// from being read but never declared.
func TestEveryFlagTheCommandsReadIsRegistered(t *testing.T) {
	t.Parallel()
	for cmd, names := range map[*cobra.Command][]string{
		installCmd: {flagTrust, flagSave, flagFrom, cmdx.FlagYes, cmdx.FlagPin, cmdx.FlagGlobal, cmdx.FlagVersion},
		syncCmd:    {flagTrust, cmdx.FlagYes},
		pinCmd:     {cmdx.FlagYes, cmdx.FlagRemove},
	} {
		for _, name := range names {
			if cmd.Flags().Lookup(name) == nil {
				t.Errorf("%s: flag --%s is read but not registered", cmd.Name(), name)
			}
		}
	}
}

// Running the whole command, not runInstall, reads every flag the way a user's
// invocation does.
//
//nolint:paralleltest // the config folder override and the seams are process-wide
func TestInstallCommandReadsItsFlagsEndToEnd(t *testing.T) {
	w := implicitMachine(t)
	scratch := t.TempDir()
	writeCatalogSkill(t, scratch, "scratch-skill")

	root := rootCmd
	root.SetArgs([]string{"install", "--global", "--from", filepath.Join(scratch, "scratch-skill"), "--yes"})
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	t.Cleanup(func() {
		root.SetOut(nil)
		root.SetErr(nil)
		root.SetArgs(nil)
		for name, def := range map[string]string{flagFrom: "", cmdx.FlagGlobal: "false", cmdx.FlagYes: "false"} {
			_ = installCmd.Flags().Set(name, def)
		}
	})
	if err := root.ExecuteContext(t.Context()); err != nil {
		t.Fatalf("err = %v\n%s", err, buf.String())
	}
	if _, err := os.Lstat(filepath.Join(w.global, "scratch-skill")); err != nil {
		t.Errorf("not linked: %v\n%s", err, buf.String())
	}
}
