package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
)

// privateProject is a git project on a machine that has registered, and
// trusted, a private source holding "alpha" (a skill with a unique description).
func privateProject(t *testing.T) (project string) {
	t.Helper()
	w := implicitMachine(t)
	w.repo.Skills = map[string]string{"alpha": "1.0.0"}

	entry := sources.Entry{ID: "@acme/hidden", Type: "git", URL: "git@github.com:acme/hidden.git", Private: true}
	file, err := sources.Add(sources.File{}, entry)
	if err != nil || sources.Save(sources.DefaultPath(), file) != nil {
		t.Fatal(err)
	}
	grants := trust.Grant(trust.File{}, entry.ID, entry.URL, time.Now())
	if err = trust.Save(trust.DefaultPath(), grants); err != nil {
		t.Fatal(err)
	}

	project = t.TempDir()
	if out, gitErr := exec.CommandContext(t.Context(), "git", "-C", project, "init", "-q").
		CombinedOutput(); gitErr != nil {
		t.Skipf("git: %v %s", gitErr, out)
	}
	mod := manifest.Mod{Tools: []tool.ID{tool.ClaudeID, tool.OpenCodeID}}
	if err = manifest.WriteMod(filepath.Join(project, ".rsk"), mod); err != nil {
		t.Fatal(err)
	}
	t.Chdir(project)
	return project
}

func installPrivate(t *testing.T, opts installOpts) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	installCmd.SetOut(&buf)
	installCmd.SetErr(&buf)
	installCmd.SetContext(context.Background())
	t.Cleanup(func() { installCmd.SetOut(nil); installCmd.SetErr(nil) })
	opts.yes = true
	err := runInstall(installCmd, opts, []string{"@acme/hidden/alpha"})
	return buf.String(), err
}

func readOrEmpty(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestAPrivatePinNeverReachesACommittedFile(t *testing.T) {
	project := privateProject(t)
	out, err := installPrivate(t, installOpts{pin: true})
	if err != nil {
		t.Fatalf("err = %v\n%s", err, out)
	}

	for _, rel := range []string{".rsk/rsk.mod", ".rsk/rsk.lock", ".rsk/SKILLS.md", ".rsk/CLAUDE.md", "CLAUDE.md", "AGENTS.md", "opencode.json"} {
		text := readOrEmpty(filepath.Join(project, rel))
		for _, leak := range []string{"hidden", "alpha", "git@github.com"} {
			if strings.Contains(text, leak) {
				t.Errorf("%s carries %q:\n%s", rel, leak, text)
			}
		}
	}
	for rel, want := range map[string]string{
		".rsk/rsk.local.mod": "alpha", ".rsk/rsk.local.lock": "git@github.com:acme/hidden.git",
		".rsk/SKILLS.local.md": "alpha", "CLAUDE.local.md": "@.rsk/SKILLS.local.md", "opencode.json": ".rsk/SKILLS.local.md",
	} {
		if !strings.Contains(readOrEmpty(filepath.Join(project, rel)), want) {
			t.Errorf("%s lacks %q", rel, want)
		}
	}
	status, _ := exec.CommandContext(t.Context(), "git", "-C", project, "status", "--porcelain", "-uall").Output()
	for line := range strings.SplitSeq(string(status), "\n") {
		if strings.Contains(line, "local") {
			t.Errorf("git would commit a local file: %q", line)
		}
	}

	// A teammate's checkout, which has only what git carries, is clean.
	merged, err := manifest.ReadMod(filepath.Join(project, ".rsk"))
	if err != nil || !merged.Local["alpha"] {
		t.Fatalf("the owner's merged read lost the private skill: %+v, %v", merged, err)
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestSharedPutsAPrivateSkillInTheCommittedFilesOnPurpose(t *testing.T) {
	project := privateProject(t)
	if out, err := installPrivate(t, installOpts{pin: true, shared: true}); err != nil {
		t.Fatalf("err = %v\n%s", err, out)
	}
	if !strings.Contains(readOrEmpty(filepath.Join(project, ".rsk", "rsk.mod")), "alpha") ||
		!strings.Contains(readOrEmpty(filepath.Join(project, ".rsk", "SKILLS.md")), "alpha") {
		t.Error("--shared did not write the committed files")
	}
	if _, err := os.Stat(filepath.Join(project, ".rsk", "rsk.local.mod")); err == nil {
		t.Error("--shared left a local layer behind")
	}
}

//nolint:paralleltest // the config folder override and the seams are process-wide
func TestPinningAnInstalledPrivateSkillKeepsItLocal(t *testing.T) {
	project := privateProject(t)
	if out, err := installPrivate(t, installOpts{}); err != nil {
		t.Fatalf("err = %v\n%s", err, out)
	}
	var buf bytes.Buffer
	pinCmd.SetOut(&buf)
	pinCmd.SetErr(&buf)
	pinCmd.SetContext(context.Background())
	t.Cleanup(func() { pinCmd.SetOut(nil); pinCmd.SetErr(nil); _ = pinCmd.Flags().Set("yes", "false") })
	if err := pinCmd.Flags().Set("yes", "true"); err != nil {
		t.Fatal(err)
	}
	if err := runPin(pinCmd, []string{"alpha"}); err != nil {
		t.Fatalf("pin: %v\n%s", err, buf.String())
	}
	for _, rel := range []string{".rsk/rsk.mod", ".rsk/SKILLS.md"} {
		if strings.Contains(readOrEmpty(filepath.Join(project, rel)), "alpha") {
			t.Errorf("rsk pin moved the private skill into %s", rel)
		}
	}
	if !strings.Contains(readOrEmpty(filepath.Join(project, ".rsk", "SKILLS.local.md")), "alpha") {
		t.Error("the pin did not reach the local digest")
	}
}
