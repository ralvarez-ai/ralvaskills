package main

import (
	"context"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/term/output"

	"go.ralva.dev/ralvaskills/v3/internal/sourceops"
)

func runSourceScan(cmd *cobra.Command, args []string) error {
	if err := checkOutput(cmd); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(cmd.Context(), testTimeout)
	defer cancel()
	res, err := sourceops.Scan(ctx, sourceEnv(), args[0])
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	if printer.Format(cmd) == output.FormatJSON {
		return writeJSON(out, res)
	}
	printScan(out, &res)
	return nil
}

func printScan(out io.Writer, res *sourceops.ScanOutcome) {
	fmt.Fprintf(out, "%s (%s, %s): %d skills", res.Source.ID, res.Source.Type, res.Source.URL, len(res.Skills))
	fmt.Fprintf(out, ", %d clashes, %d left out\n", len(res.Clashes), len(res.Skipped))

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	for i := range res.Skills {
		s := &res.Skills[i]
		fmt.Fprintf(w, "skill\t%s\t%s\t%s\n", s.Name, s.Version, s.Path)
		for _, issue := range s.Issues {
			fmt.Fprintf(w, "\t\tissue\t%s\n", issue)
		}
	}
	for _, c := range res.Clashes {
		for _, p := range c.Paths {
			fmt.Fprintf(w, "clash\t%s\t\t%s\n", c.Name, p)
		}
	}
	for _, sk := range res.Skipped {
		fmt.Fprintf(w, "left out\t%s\t\t%s\n", sk.Path, sk.Reason)
	}
	_ = w.Flush()
	if len(res.Clashes) > 0 {
		fmt.Fprintln(out, "\nA source with a clash cannot be resolved: exclude one directory (--exclude) or rename it.")
	}
}
