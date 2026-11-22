package main

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/term"
	"go.ralva.dev/lib/term/output"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/resolve"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

// newResolveCmd builds `rsk resolve`, the read-only twin of the resolution
// `rsk install` performs. It is a factory so a test can run it on a fresh tree.
func newResolveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve <reference>",
		Short: "Explain what a skill reference resolves to.",
		Long: `Show the source, path, version and commit a reference gives, without installing
anything, and for a bare name say why it is ambiguous or unavailable.

A reference is a bare name (go-grpc), a qualified one (@owner/repo/name) or either
with a version (name@1.2.0). A bare name is looked up in every source; it resolves
only when exactly one has it, and fails when a source cannot be consulted, so an
unreachable private source never lets a public one answer in its place.

The command exits non-zero unless the reference resolves (or names a bundle). With
--output json the result is printed either way, so a script reads its "status".

Examples:
  rsk resolve go-grpc
  rsk resolve @acme/skills/api-design
  rsk resolve both -o json`,
		Args: cobra.ExactArgs(1),
		RunE: runResolve,
	}
	addOutputFlag(cmd)
	return cmd
}

var resolveCmd = newResolveCmd()

func runResolve(cmd *cobra.Command, args []string) error {
	if err := checkOutput(cmd); err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("%w\n  Run 'rsk init' to set up rsk on this machine", err)
	}
	chain, err := sourceChain(cmd, cfg)
	if err != nil {
		return err
	}
	catalog, catalogWarn := config.LoadCatalog("")
	if catalogWarn != nil {
		ui.Warn(cmd.ErrOrStderr(), fmt.Sprintf("user catalog: %v", catalogWarn))
	}

	res, err := resolve.Reference(
		cmd.Context(),
		args[0],
		chain,
		resolve.Options{Bundles: catalog, ImplicitHint: sources.ImplicitHint},
	)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if printer.Format(cmd) == output.FormatJSON {
		if err = writeJSON(out, res); err != nil {
			return err
		}
	} else {
		printResolution(out, &res)
	}
	if res.Status != resolve.Resolved && res.Status != resolve.Bundle {
		return fmt.Errorf("%s is %s", args[0], res.Status)
	}
	return nil
}

func printResolution(out io.Writer, res *resolve.Result) {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	row := func(key, value string) {
		if value != "" {
			fmt.Fprintf(w, "%s\t%s\n", key, value)
		}
	}
	row("reference", res.Reference)
	row("status", string(res.Status))
	if s := res.Skill; s != nil {
		row("source", s.Source+" ("+s.Type+")")
		row("url", s.URL)
		row("path", s.Path)
		row("version", s.Version)
		row("requested", res.Requested)
		row("rev", s.Rev)
		row("commit", s.Commit)
		row("tree", s.Tree)
		row("key", s.Key)
	}
	for _, m := range res.Members {
		row("member", m)
	}
	row("why", res.Message)
	_ = w.Flush()

	if len(res.Consulted) > 0 {
		fmt.Fprintln(out)
		t := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		fmt.Fprintln(t, "SOURCE\tTYPE\tHAS\tVERSION\tNOTE")
		for _, a := range res.Consulted {
			has := "no"
			if a.Has {
				has = "yes"
			}
			fmt.Fprintf(t, "%s\t%s\t%s\t%s\t%s\n", a.Source, a.Type, has, a.Version, a.Error)
		}
		_ = t.Flush()
	}
}

func registerResolveCommand() {
	rootCmd.AddCommand(resolveCmd)
	term.MarkTableView(resolveCmd)
}
