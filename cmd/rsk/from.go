package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"go.ralva.dev/ralvaskills/v3/internal/cmdx"
	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/install"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

// oneOffNote is shown with every --from install, and says plainly what is not
// kept: a path on this machine is not something another machine can restore.
const oneOffNote = "--from is a one-off: it is not recorded in rsk.mod or rsk.lock, so rsk sync, rsk update and " +
	"rsk status do not know it and another machine cannot restore it. For a directory you keep using, " +
	"register it: rsk source add <dir>"

// validateFromOpts checks the flag combinations --from allows: it installs one
// directory and nothing else, since everything else it could combine with needs
// a record it does not make.
func validateFromOpts(opts installOpts, args []string) error {
	switch {
	case len(args) > 0:
		return fieldError(flagFrom, "--from installs one directory and takes no names")
	case opts.pin:
		return fieldError(cmdx.FlagPin, "--pin needs a manifest entry, and --from records none")
	case opts.save || len(opts.trust) > 0:
		return fieldError(flagFrom, "--save and --trust apply to unregistered @owner/repo sources, not to --from")
	case opts.version != "":
		return fieldError(cmdx.FlagVersion, "--version does not apply to --from")
	}
	return nil
}

// runInstallFrom links the skill in one directory into the project (or, with
// --global, the global skills directories) without recording anything.
func runInstallFrom(cmd *cobra.Command, opts installOpts, args []string) error {
	if err := validateFromOpts(opts, args); err != nil {
		return err
	}
	s, err := source.ScanSkillDir(opts.from)
	if err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("%w\n  Run 'rsk init' to set up rsk on this machine", err)
	}
	targets, err := oneOffTargets(cfg, opts)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if !previewAndConfirm(cmd, out, []skill.Skill{s}, targets, []string{oneOffNote}, opts.dryRun, opts.yes) {
		return nil
	}
	results, err := install.Apply([]skill.Skill{s}, install.Options{Scope: install.ScopeOneOff, Targets: targets})
	if err != nil {
		return err
	}
	return reportInstall(out, cmd.ErrOrStderr(), results)
}

func oneOffTargets(cfg config.Config, opts installOpts) ([]string, error) {
	if opts.global {
		return resolveTargetDirs(cfg, true, opts.forTool)
	}
	rskDir, err := manifest.ProjectFolderPath()
	if err != nil {
		return nil, err
	}
	m, err := manifest.ReadMod(rskDir)
	if err != nil {
		return nil, err
	}
	return projectSkillsDirs(filepath.Dir(rskDir), m), nil
}
