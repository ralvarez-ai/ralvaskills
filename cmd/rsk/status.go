package main

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/term"
	"go.ralva.dev/lib/term/output"

	"go.ralva.dev/ralvaskills/v3/internal/cmdx"
	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/status"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

var statusCmd = &cobra.Command{
	Use:   "status [flags]",
	Short: "Show installed skills and version drift (see also: list).",
	Long: `List installed skills with their versions and drift. Unlike 'rsk list'
(the offline manifest/installed view), status compares versions across
project + global.

Without --stack, no network calls are made.
With --stack, fetches latest versions from proxy.golang.org and pypi.org
and highlights skills whose STACK.md may be outdated (results cached 24h).

Examples:
  rsk status
  rsk status --global
  rsk status --project
  rsk status --stack
  rsk status --stack --refresh
  rsk status -o json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inc, err := cmdx.ReadIncludes(cmd, cmdx.IncludePersonal)
		if err != nil {
			return err
		}
		forTool, err := forToolFlag(cmd)
		if err != nil {
			return err
		}
		return runStatus(cmd, statusOpts{
			global:   cmdx.Bool(cmd, cmdx.FlagGlobal),
			project:  cmdx.Bool(cmd, "project"),
			stack:    cmdx.Bool(cmd, cmdx.FlagStack),
			refresh:  cmdx.Bool(cmd, cmdx.FlagRefresh),
			personal: inc.Personal,
			forTool:  forTool,
			format:   printer.Format(cmd),
		})
	},
}

type (
	// statusOpts holds the options for the status command.
	statusOpts struct {
		global, project, stack, refresh, personal bool
		forTool                                   tool.ID
		format                                    output.OutputFormat
	}

	// statusSkillEntry is one skill row inside a statusSectionEntry, serialized
	// as JSON when -o json is set.
	statusSkillEntry struct {
		Name    string   `json:"name"`
		Version string   `json:"version,omitempty"`
		Source  string   `json:"source,omitempty"`
		Pinned  bool     `json:"pinned,omitempty"`
		Bundles []string `json:"bundles,omitempty"`
		// RenamedTo is set when the linked name is a retired alias; it holds the
		// name that replaced it.
		RenamedTo string `json:"renamed_to,omitempty"`
	}

	// statusSectionEntry is one section (a scanned target directory) in the
	// JSON status output, mirroring status.Section's text rendering.
	statusSectionEntry struct {
		Title    string             `json:"title"`
		Subtitle string             `json:"subtitle,omitempty"`
		Dir      string             `json:"dir"`
		Skills   []statusSkillEntry `json:"skills"`
	}
)

func runStatus(cmd *cobra.Command, opts statusOpts) error {
	if err := checkOutput(cmd); err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	capture := term.CaptureFromContext(cmd.Context())

	if err := validateStatusOpts(opts); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("%w\n  Run 'rsk init' to set up rsk on this machine", err)
	}

	// --global means "report on the global install only" — no project
	// (rsk.mod) needs to exist in the current directory, so skip resolving
	// it entirely, matching how resolveTargetDirs branches for list/install.
	var rskDir string
	if !opts.global {
		rskDir, err = manifest.ProjectFolderPath()
		if err != nil {
			return err
		}
	}

	catalog, catalogWarn := config.LoadCatalog("")
	if catalogWarn != nil {
		warnOrCapture(out, capture, fmt.Sprintf("user catalog: %v", catalogWarn))
	}
	membership := status.BundleMembership(catalog)

	m, modErr := manifest.ReadMod(rskDir)
	var projectDirs []string
	if modErr == nil {
		projectDirs = projectSkillsDirs(filepath.Dir(rskDir), m)
	}
	sections := status.BuildSections(cfg, opts.global, opts.project, opts.forTool, projectDirs)

	pinnedSet := make(map[string]bool)
	if modErr == nil {
		for _, p := range m.Pinned {
			pinnedSet[p] = true
		}
	}

	for _, w := range status.Scan(sections, cfg, membership, opts.personal) {
		warnOrCapture(out, capture, w)
	}
	renamed := offlineAliases(cfg)

	if opts.format == output.FormatJSON {
		return newFormatter(cmd, opts.format).JSON(statusSectionsToEntries(sections, pinnedSet, renamed))
	}
	if err = printStatusText(out, capture, cmd, opts.format, sections, pinnedSet, renamed); err != nil {
		return err
	}
	// The captured (picker) run must not write behind the TUI, and --global
	// has no project index to size.
	if capture == nil && !opts.global && modErr == nil {
		printSkillsIndexStatus(out, filepath.Dir(rskDir), m.Pinned)
		printGitLockDrift(out, rskDir, projectDirs, cfg)
	}
	return nil
}

// printGitLockDrift warns about git skills whose links are not what the project's
// lock records. It reads only local files.
func printGitLockDrift(out io.Writer, rskDir string, projectDirs []string, cfg config.Config) {
	lock, err := manifest.ReadLock(rskDir)
	if err != nil {
		return
	}
	for _, line := range gitLockDrift(lock, projectDirs, cfg.SourceOptions().CacheRoot) {
		ui.Warn(out, line)
	}
}

// warnOrCapture reports msg on the capture (when a captured in-process run is
// in progress) instead of writing it to out, since nothing written to out is
// seen while the TUI owns the terminal.
func warnOrCapture(out io.Writer, capture *term.Capture, msg string) {
	if capture != nil {
		capture.AddMessage(msg)
		return
	}
	ui.Warn(out, msg)
}

// validateStatusOpts checks runStatus's flag combinations before doing any
// scanning work.
func validateStatusOpts(opts statusOpts) error {
	if opts.stack {
		return errors.New("--stack is not yet implemented")
	}
	if !opts.stack && opts.refresh {
		return errors.New("--refresh requires --stack")
	}
	if opts.global && opts.project {
		return errors.New("--global and --project are mutually exclusive")
	}
	if !opts.global && opts.forTool != "" {
		return errors.New("--for requires --global")
	}
	return nil
}

// printStatusText renders every non-empty section as a table, or a "no
// skills installed" notice when all sections are empty.
func printStatusText(
	out io.Writer, capture *term.Capture, cmd *cobra.Command,
	format output.OutputFormat, sections []status.Section, pinnedSet map[string]bool, renamed map[string]string,
) error {
	hasAny := false
	for _, sec := range sections {
		if len(sec.Skills) == 0 {
			continue
		}
		hasAny = true
		if capture == nil {
			fmt.Fprintln(out)
			ui.SectionHeader(out, sec.Title, sec.Subtitle)
		}
		if err := printStatusSectionRows(cmd, format, sec.Title, sec.Skills, pinnedSet, renamed); err != nil {
			return err
		}
	}

	if capture != nil {
		if !hasAny {
			capture.AddMessage("No skills installed. Run 'rsk install <bundle>' to get started.")
		}
		return nil
	}

	fmt.Fprintln(out)
	if !hasAny {
		ui.Warn(out, "No skills installed.")
		ui.Indent(out, "Run 'rsk install <bundle>' to get started.")
	}
	return nil
}

// printStatusSectionRows renders one row per linked skill, tagged with
// [pinned] and/or its bundle memberships. title becomes the captured table's
// Title, so the result view's tab line identifies which section is active
// when a captured status run has more than one.
func printStatusSectionRows(
	cmd *cobra.Command, format output.OutputFormat, title string,
	skills []status.Linked, pinnedSet map[string]bool, renamed map[string]string,
) error {
	header := []string{headerSource, headerName, headerVersion, "", ""}
	rows := term.Rows(skills, statusRowTable(pinnedSet, renamed))

	return renderTitledTable(cmd, format, title, rows, header, rows, linkedSkillNames(skills))
}

// linkedSkillNames extracts each row's skill name, parallel to statusRowTable's
// rows, so a captured table can identify the row a RowAction fires on
// independently of how it's displayed.
func linkedSkillNames(skills []status.Linked) []string {
	names := make([]string, len(skills))
	for i, e := range skills {
		names[i] = e.Name
	}
	return names
}

// statusRowActionScope resolves the captured status table a row action fired
// on to the uninstall scope its rows belong to: a "Global — <tool>" section
// runs uninstall with --global --for <tool>, while a "Project" section keeps
// uninstall's project default. This is what lets a bare `rsk status`, which
// renders one table per scanned directory, uninstall from the right one.
func statusRowActionScope(table term.Data) map[string]string {
	if toolID, ok := strings.CutPrefix(table.Title, "Global — "); ok {
		return map[string]string{cmdx.FlagGlobal: "true", cmdx.FlagFor: toolID}
	}
	return nil
}

// statusRowTable projects one linked skill into its status row. The pinned set
// is captured because it is per-section, not per-entry.
func statusRowTable(pinnedSet map[string]bool, renamed map[string]string) term.Table[status.Linked] {
	return term.Table[status.Linked]{
		Row: func(e status.Linked) []any {
			return []any{
				ui.SourceLabel(e.Source),
				ui.SkillName(e.Name),
				ui.SkillVersion(e.Version),
				ui.SuccessMark,
				statusRowTags(e, pinnedSet, renamed),
			}
		},
	}
}

// statusRowTags builds the trailing "  [pinned]  bundle-tag ..." suffix for
// one status row.
func statusRowTags(e status.Linked, pinnedSet map[string]bool, renamed map[string]string) string {
	tags := ""
	if next, ok := renamed[e.Name]; ok {
		tags += "  [renamed \u2192 " + next + "]"
	}
	if pinnedSet[e.Name] {
		tags += "  [pinned]"
	}
	if len(e.Bundles) > 0 {
		bundleTags := make([]string, len(e.Bundles))
		for i, b := range e.Bundles {
			bundleTags[i] = ui.BundleTag(b)
		}
		tags += "  " + strings.Join(bundleTags, " ")
	}
	return tags
}

// statusSectionsToEntries converts the scanned statusSections into their
// JSON-serializable shape, carrying over the pinned marker (which is tracked
// separately from status.Section/status.Linked) onto each skill row.
func statusSectionsToEntries(
	sections []status.Section, pinnedSet map[string]bool, renamed map[string]string,
) []statusSectionEntry {
	entries := make([]statusSectionEntry, len(sections))
	for i, sec := range sections {
		skills := make([]statusSkillEntry, len(sec.Skills))
		for j, e := range sec.Skills {
			skills[j] = statusSkillEntry{
				Name:      e.Name,
				Version:   e.Version,
				Source:    e.Source.String(),
				Pinned:    pinnedSet[e.Name],
				Bundles:   e.Bundles,
				RenamedTo: renamed[e.Name],
			}
		}
		entries[i] = statusSectionEntry{
			Title:    sec.Title,
			Subtitle: sec.Subtitle,
			Dir:      sec.Dir,
			Skills:   skills,
		}
	}
	return entries
}
