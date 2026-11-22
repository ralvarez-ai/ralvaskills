package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/term"
	"go.ralva.dev/lib/term/output"

	"go.ralva.dev/ralvaskills/v3/internal/cmdx"
	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/sourceops"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

const (
	actionAdd    = sourceops.ActionAdd
	actionRemove = sourceops.ActionRemove
	flagID       = "id"
	flagRev      = "rev"
	flagRoots    = "roots"
	flagExclude  = "exclude"
	flagPrivate  = "private"
	testTimeout  = 30 * time.Second
	trustTrusted = sourceops.TrustTrusted
	trustNone    = sourceops.TrustNone
	trustNA      = sourceops.TrustNA
	trustBuiltIn = sourceops.TrustBuiltIn
)

type (
	// sourceCmds is the `rsk source` command tree.
	sourceCmds struct {
		root    *cobra.Command
		add     *cobra.Command
		list    *cobra.Command
		remove  *cobra.Command
		test    *cobra.Command
		scan    *cobra.Command
		trust   *cobra.Command
		untrust *cobra.Command
	}
)

// The result types are the core's, so the CLI and the MCP tools print the same
// JSON.
type (
	sourceView   = sourceops.View
	sourceChange = sourceops.Change
	sourceTest   = sourceops.TestResult
)

// canonicalSource is the registered `rsk source` tree.
var canonicalSource = newSourceCmds()

func newSourceCmds() sourceCmds {
	root := &cobra.Command{
		Use:   "source",
		Short: "Manage the skill sources on this machine.",
		Long: `Manage where skills come from besides the hosted registry: git repositories
(including private ones) and local directories.

Sources are declared per machine in ~/.config/rsk/sources.json. rsk stores no
credential: a private repository is reached with your own SSH key or git
credential helper. A newly added git source is untrusted; rsk contacts nothing
until you run 'rsk source trust <id>' yourself in a terminal.

Examples:
  rsk source add @acme/skills-private --roots skills
  rsk source trust @acme/skills-private
  rsk source test @acme/skills-private
  rsk source list`,
	}

	add := &cobra.Command{
		Use:   "add <@owner/repo | git-url | path>",
		Short: "Declare a skill source (registered untrusted).",
		Long: `Declare a git repository or a local directory as a skill source.

@owner/repo means the GitHub repository over SSH (git@github.com:owner/repo.git).
A full ssh://, https:// or user@host:path URL works for any host; a host other
than GitHub needs --id. A directory becomes a path source, named @local/<dir>
unless --id says otherwise.

The source is registered as untrusted and nothing is contacted. Preview with
--dry-run --output json.

Examples:
  rsk source add @acme/skills-private --roots skills
  rsk source add git@git.example.com:team/skills.git --id @team/skills --rev main
  rsk source add ./my-skills --id @me/scratch
  rsk source add @acme/skills-private --dry-run -o json`,
		Args: cobra.ExactArgs(1),
		RunE: runSourceAdd,
	}
	f := add.Flags()
	f.String(flagID, "", "Identity @owner/repo (derived for GitHub; required for other hosts)")
	f.String(flagRev, "", "Git branch or tag to read (default: the repository's default branch)")
	f.StringSlice(flagRoots, nil, "Directories inside the source to search for skills")
	f.StringSlice(flagExclude, nil, "Glob patterns to skip while searching")
	f.Bool(flagPrivate, false, "Keep the source out of listings (not supported yet)")
	f.Bool(cmdx.FlagDryRun, false, "Show what would be added without doing it")
	f.BoolP(cmdx.FlagYes, "y", false, "Skip the confirmation prompt")
	addOutputFlag(add)

	list := &cobra.Command{
		Use:   "list",
		Short: "List the sources rsk knows about and whether each is trusted.",
		Args:  cobra.NoArgs,
		RunE:  runSourceList,
	}
	addOutputFlag(list)

	remove := &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove a declared source and its trust grant.",
		Args:  cobra.ExactArgs(1),
		RunE:  runSourceRemove,
	}
	remove.Flags().Bool(cmdx.FlagDryRun, false, "Show what would be removed without doing it")
	remove.Flags().BoolP(cmdx.FlagYes, "y", false, "Skip the confirmation prompt")
	addOutputFlag(remove)

	test := &cobra.Command{
		Use:   "test <id>",
		Short: "Check that a source is trusted and reachable.",
		Long: `Check a source without installing anything. A git source must be trusted first;
the check then asks the remote for its branch with 'git ls-remote'.`,
		Args: cobra.ExactArgs(1),
		RunE: runSourceTest,
	}
	addOutputFlag(test)

	scan := &cobra.Command{
		Use:   "scan <directory | @owner/repo>",
		Short: "Show what a source would provide, and what is wrong with it.",
		Long: `Dry-run discovery: the skills a source would resolve, the directories left out
and why, the names two directories claim (which make resolving the source fail),
and frontmatter problems such as a name that differs from its directory.

A directory needs no registration and contacts nothing. A repository must be
registered and trusted first: scanning it fetches it. Nothing is installed or
registered.

Examples:
  rsk source scan ./my-skills
  rsk source scan @acme/skills-private`,
		Args: cobra.ExactArgs(1),
		RunE: runSourceScan,
	}
	addOutputFlag(scan)

	trustCmd := &cobra.Command{
		Use:   "trust <id>",
		Short: "Allow rsk to contact a declared git source (needs a terminal).",
		Long: `Record that you trust a declared git source. From then on rsk may let git
contact its host with your SSH key or credential helper, and may load its skills
into your agents' context.

This works only in an interactive terminal, and only you can run it: an agent
that runs shell commands is refused and told to hand the command back to you.
The grant covers the identity at its current URL; changing the URL revokes it.`,
		Args: cobra.ExactArgs(1),
		RunE: runSourceTrust,
	}

	untrust := &cobra.Command{
		Use:   "untrust <id>",
		Short: "Withdraw the trust granted to a source.",
		Args:  cobra.ExactArgs(1),
		RunE:  runSourceUntrust,
	}

	root.AddCommand(add, list, remove, test, scan, trustCmd, untrust)
	return sourceCmds{
		root: root, add: add, list: list, remove: remove, test: test, scan: scan, trust: trustCmd, untrust: untrust,
	}
}

type sourceState struct {
	file   sources.File
	grants trust.File
}

func loadSourceState() (sourceState, error) {
	file, err := sources.Load(sources.DefaultPath())
	if err != nil {
		return sourceState{}, err
	}
	grants, err := trust.Load(trust.DefaultPath())
	if err != nil {
		return sourceState{}, err
	}
	return sourceState{file: file, grants: grants}, nil
}

func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("write json: %w", err)
	}
	return nil
}

// sourceEnv is where the source commands read and write: the machine's own
// files, plus the built-in sources when the configuration loads.
func sourceEnv() sourceops.Env {
	env := sourceops.Env{RepoFor: sourceRepoFor, EffectiveURL: effectiveURL}
	if cfg, err := config.Load(); err == nil {
		env.Cfg = &cfg
	}
	return env
}

func runSourceAdd(cmd *cobra.Command, args []string) error {
	if err := checkOutput(cmd); err != nil {
		return err
	}
	if cmdx.Bool(cmd, flagPrivate) {
		return fieldError(
			flagPrivate,
			"--private is not supported yet: private pins need the local overlay, which has not shipped",
		)
	}
	roots, err := cmd.Flags().GetStringSlice(flagRoots)
	if err != nil {
		return err
	}
	excludes, err := cmd.Flags().GetStringSlice(flagExclude)
	if err != nil {
		return err
	}
	change, apply, err := sourceops.PlanAdd(sourceEnv(), args[0], sources.Options{
		ID: cmdx.String(cmd, flagID), Rev: cmdx.String(cmd, flagRev), Roots: roots, Excludes: excludes,
	})
	if err != nil {
		return err
	}
	change.DryRun = cmdx.Bool(cmd, cmdx.FlagDryRun)
	return applySourceChange(cmd, change, apply)
}

func runSourceRemove(cmd *cobra.Command, args []string) error {
	if err := checkOutput(cmd); err != nil {
		return err
	}
	change, apply, err := sourceops.PlanRemove(sourceEnv(), args[0])
	if err != nil {
		return err
	}
	change.DryRun = cmdx.Bool(cmd, cmdx.FlagDryRun)
	return applySourceChange(cmd, change, apply)
}

// applySourceChange previews a change, asks, applies it, and reports the result
// as text or as the same JSON the preview used.
func applySourceChange(cmd *cobra.Command, change sourceChange, apply func() error) error {
	out := cmd.OutOrStdout()
	jsonOut := printer.Format(cmd) == output.FormatJSON

	if !jsonOut {
		verb := map[string]string{actionAdd: "Add", actionRemove: "Remove"}[change.Action]
		fmt.Fprintf(
			out,
			"%s source %s (%s, %s), trust: %s\n",
			verb,
			change.Source.ID,
			change.Source.Type,
			change.Source.URL,
			change.Source.Trust,
		)
	}
	if change.DryRun {
		if jsonOut {
			return writeJSON(out, change)
		}
		ui.Info(out, "dry run: nothing changed")
		return nil
	}
	if !term.Confirm(cmd, cmdx.Bool(cmd, cmdx.FlagYes), term.StdioIsTerminal(), "Proceed?") {
		fmt.Fprintln(out, messageAborted)
		return nil
	}
	if err := apply(); err != nil {
		return err
	}
	change.Applied = true
	if jsonOut {
		return writeJSON(out, change)
	}
	ui.Success(out, map[string]string{actionAdd: "added ", actionRemove: "removed "}[change.Action]+change.Source.ID)
	if change.Next != "" {
		ui.Info(out, change.Next)
	}
	return nil
}

func runSourceList(cmd *cobra.Command, _ []string) error {
	if err := checkOutput(cmd); err != nil {
		return err
	}
	views, err := sourceops.List(sourceEnv())
	if err != nil {
		return err
	}
	if len(views) == 0 {
		emptyHint(cmd, false, "No sources. Add one with 'rsk source add'.")
		return nil
	}

	header := []string{"ID", "TYPE", "TARGET", "REV", "TRUST"}
	rows := make([][]any, len(views))
	ids := make([]string, len(views))
	for i := range views {
		v := &views[i]
		rows[i] = []any{v.ID, v.Type, v.URL, v.Rev, v.Trust}
		ids[i] = v.ID
	}
	return renderList(cmd, printer.Format(cmd), "", views, header, rows, ids)
}

func runSourceTest(cmd *cobra.Command, args []string) error {
	if err := checkOutput(cmd); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(cmd.Context(), testTimeout)
	defer cancel()
	result, testErr := sourceops.Test(ctx, sourceEnv(), args[0])
	if result.ID == "" {
		return testErr // the source is not declared, or its files are unreadable
	}

	out := cmd.OutOrStdout()
	if printer.Format(cmd) == output.FormatJSON {
		if err := writeJSON(out, result); err != nil {
			return err
		}
	} else if testErr == nil {
		ui.Success(out, fmt.Sprintf("%s is reachable (%s %s)", result.ID, result.Type, result.Commit))
		if result.Via != "" {
			ui.Info(out, "your git configuration rewrites "+result.URL+" to "+result.Via)
		}
	}
	return testErr
}

func runSourceTrust(cmd *cobra.Command, args []string) error {
	state, err := loadSourceState()
	if err != nil {
		return err
	}
	entry, ok := sources.Find(state.file, args[0])
	if !ok {
		return fmt.Errorf("source %s is not declared: add it first with 'rsk source add'", args[0])
	}
	if entry.Type == source.KindPath {
		return fmt.Errorf("source %s is a local directory and needs no trust", entry.ID)
	}
	if !term.StdioIsTerminal() {
		return fmt.Errorf(
			"trusting a source needs an interactive terminal and cannot be done from a script or an agent: run this yourself: rsk source trust %s",
			entry.ID,
		)
	}

	out := cmd.OutOrStdout()
	if state.grants.Granted(entry.ID, entry.URL) {
		ui.Info(out, fmt.Sprintf("%s is already trusted at %s", entry.ID, entry.URL))
		return nil
	}
	fmt.Fprintf(out, "Trust %s at %s?\n", entry.ID, entry.URL)
	fmt.Fprintln(out, "rsk will let git contact this host with your own SSH key or credential helper, and may")
	fmt.Fprintln(out, "load the skills it provides into your agents' context. The grant is for this URL only.")
	if !term.Confirm(cmd, false, true, "Trust this source?") {
		fmt.Fprintln(out, messageAborted)
		return nil
	}
	grants := trust.Grant(state.grants, entry.ID, entry.URL, time.Now())
	if err = trust.Save(trust.DefaultPath(), grants); err != nil {
		return err
	}
	ui.Success(out, "trusted "+entry.ID)
	return nil
}

func runSourceUntrust(cmd *cobra.Command, args []string) error {
	state, err := loadSourceState()
	if err != nil {
		return err
	}
	canonical := skill.CanonicalSourceID(args[0])
	grants, removed := trust.Revoke(state.grants, canonical)
	if !removed {
		return errors.New("source " + canonical + " has no trust grant")
	}
	if err = trust.Save(trust.DefaultPath(), grants); err != nil {
		return err
	}
	ui.Success(cmd.OutOrStdout(), "no longer trusting "+canonical)
	return nil
}

// registerSourceCommands adds the `rsk source` tree to the root command.
func registerSourceCommands() {
	rootCmd.AddCommand(canonicalSource.root)
	term.MarkTableView(canonicalSource.list)
	for _, c := range []*cobra.Command{canonicalSource.add, canonicalSource.remove} {
		term.MarkDestructive(c)
	}
	canonicalSource.root.GroupID = groupSources
}
