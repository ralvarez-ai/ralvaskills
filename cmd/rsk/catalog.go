package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/term"
	"go.ralva.dev/lib/term/output"

	"go.ralva.dev/ralvaskills/v3/internal/catalog"
	"go.ralva.dev/ralvaskills/v3/internal/cmdx"
	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Browse the catalog of available skills and bundles.",
	Args:  cobra.NoArgs,
	Long: `Browse what rsk knows about. By default lists every skill.

The bundles subcommand lists the bundles themselves; use
catalog bundle <name> to list the skills inside one bundle.

Examples:
  rsk catalog                       # all skills
  rsk catalog bundles               # all bundles
  rsk catalog bundle go-grpc        # skills inside the go-grpc bundle
  rsk catalog --source local        # filter by source
  rsk catalog --include personal    # include personal/ skills
  rsk catalog -o json`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		opts, err := catalogSkillsOpts(cmd)
		if err != nil {
			return err
		}
		return runCatalogSkills(cmd, opts)
	},
}

// catalogBundlesCmd lists the catalog's bundles and how much of each is
// installed.
var catalogBundlesCmd = &cobra.Command{
	Use:   "bundles",
	Short: "List the catalog's bundles and how many of their skills are installed.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return runCatalogBundles(cmd, catalogOpts{format: printer.Format(cmd)})
	},
}

// catalogBundleCmd lists the skills that belong to one bundle.
var catalogBundleCmd = &cobra.Command{
	Use:   "bundle <name>",
	Short: "List the skills inside one bundle.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		opts, err := catalogSkillsOpts(cmd)
		if err != nil {
			return err
		}
		opts.bundle = args[0]
		return runCatalogSkills(cmd, opts)
	},
}

// catalogSkillsOpts reads the flags shared by the skill-list views: the
// --include/--personal scope, the --source filter, and the --output format.
func catalogSkillsOpts(cmd *cobra.Command) (catalogOpts, error) {
	inc, err := cmdx.ReadIncludes(cmd, cmdx.IncludePersonal, cmdx.IncludePrivate)
	if err != nil {
		return catalogOpts{}, err
	}
	sourceFilter, err := parseSourceFilter(cmdx.String(cmd, cmdx.FlagSource))
	if err != nil {
		return catalogOpts{}, err
	}
	return catalogOpts{
		personal: inc.Personal,
		private:  inc.Private,
		source:   sourceFilter,
		format:   printer.Format(cmd),
	}, nil
}

type (
	catalogOpts struct {
		personal bool
		private  bool
		bundle   string
		source   skill.Source
		format   output.OutputFormat
	}

	skillEntry struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		Source   string `json:"source"`
		SourceID string `json:"source_id,omitempty"`
		Personal bool   `json:"personal,omitempty"`
		Path     string `json:"path"`
	}

	bundleEntry struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Skills      int    `json:"skills"`
		Linked      int    `json:"linked"`
	}
)

func runCatalogSkills(cmd *cobra.Command, opts catalogOpts) error {
	if err := checkOutput(cmd); err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	ctx := cmd.Context()
	capture := term.CaptureFromContext(ctx)

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("%w\n  Run 'rsk init' to set up rsk on this machine", err)
	}

	chain, err := sourceChain(cmd, cfg)
	if err != nil {
		return err
	}
	res, err := catalog.List(ctx, chain, catalog.ListOptions{
		Source:   opts.source,
		Private:  opts.private,
		Bundle:   opts.bundle,
		Personal: opts.personal,
	})
	for _, w := range res.Warnings {
		warnOrCapture(out, capture, w)
	}
	if err != nil {
		return err
	}
	all := res.Skills

	if len(all) == 0 {
		if capture != nil {
			capture.AddMessage("No skills found.")
			return nil
		}
		ui.Info(out, "No skills found.")
		return nil
	}

	header := []string{headerSource, headerName, headerVersion}
	tableRows := term.Rows(all, catalogSkillTable)

	return renderList(cmd, opts.format, "", skillsToEntries(all), header, tableRows, skillNames(all))
}

// skillNames extracts each skill's name, parallel to catalogSkillTable's
// rows, so a captured table can identify the row a RowAction fires on
// independently of how it's displayed.
func skillNames(skills []skill.Skill) []string {
	names := make([]string, len(skills))
	for i, s := range skills {
		names[i] = s.Name
	}
	return names
}

// catalogSkillTable projects one catalog skill row.
var catalogSkillTable = term.Table[skill.Skill]{
	Row: func(s skill.Skill) []any {
		return []any{ui.SourceLabel(s.Source), ui.SkillName(s.Name), ui.SkillVersion(s.Version)}
	},
}

// declaredID returns the identity of a declared source, and "" for the others.
func declaredID(src skill.Source) string {
	if src.IsDeclared() {
		return string(src)
	}
	return ""
}

func skillsToEntries(skills []skill.Skill) []skillEntry {
	out := make([]skillEntry, len(skills))
	for i, s := range skills {
		out[i] = skillEntry{
			Name:     s.Name,
			Version:  s.Version,
			Source:   s.Source.Kind().String(),
			SourceID: declaredID(s.Source),
			Personal: s.IsPersonal,
			Path:     s.Path,
		}
	}
	return out
}

func runCatalogBundles(cmd *cobra.Command, opts catalogOpts) error {
	if err := checkOutput(cmd); err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	capture := term.CaptureFromContext(cmd.Context())

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("%w\n  Run 'rsk init' to set up rsk on this machine", err)
	}

	targets := allTargetDirs(cfg)
	rows, warnings := catalog.BundleRows(func(name string) bool { return isLinkedAnywhere(name, targets) })
	for _, w := range warnings {
		warnOrCapture(out, capture, w)
	}

	if len(rows) == 0 {
		if capture != nil {
			capture.AddMessage("No bundles found.")
			return nil
		}
		ui.Info(out, "No bundles found.")
		return nil
	}

	header := []string{"", "Bundle", "Linked", "Description"}
	tableRows := term.Rows(rows, catalogBundleTable)

	return renderList(cmd, opts.format, "", bundleRowsToEntries(rows), header, tableRows, bundleNames(rows))
}

// bundleNames extracts each bundle's name, parallel to catalogBundleTable's
// rows, so a captured table can identify the row a RowAction fires on
// independently of how it's displayed.
func bundleNames(rows []catalog.BundleRow) []string {
	names := make([]string, len(rows))
	for i, r := range rows {
		names[i] = r.Name
	}
	return names
}

// catalogBundleTable projects one bundle row.
var catalogBundleTable = term.Table[catalog.BundleRow]{
	Row: func(r catalog.BundleRow) []any {
		mark := ui.ErrorMark
		switch {
		case r.Linked == r.Total && r.Total > 0:
			mark = ui.SuccessMark
		case r.Linked > 0:
			mark = ui.WarnMark
		}
		return []any{
			mark,
			ui.BundleTag(r.Name),
			ui.MutedStyle.Render(fmt.Sprintf("%d/%d", r.Linked, r.Total)),
			ui.MutedStyle.Render(r.Description),
		}
	},
}

func bundleRowsToEntries(rows []catalog.BundleRow) []bundleEntry {
	out := make([]bundleEntry, len(rows))
	for i, r := range rows {
		out[i] = bundleEntry{
			Name:        r.Name,
			Description: r.Description,
			Skills:      r.Total,
			Linked:      r.Linked,
		}
	}
	return out
}
