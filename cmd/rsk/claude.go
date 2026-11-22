package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/term"
	"go.ralva.dev/lib/term/output"

	"go.ralva.dev/ralvaskills/v3/internal/cmdx"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

// toolsForHint is the --for flag's help text, shared by every tools subcommand.
const toolsForHint = "Tool to configure (only claude-code is supported for now)"

type (
	// toolsCmds is one registered copy of the tools command tree: the parent
	// and its four subcommands.
	toolsCmds struct {
		root, list, allow, deny, remove *cobra.Command
	}
)

// canonicalTools is `rsk tools`.
var canonicalTools = newToolsCmds("rsk tools")

// newToolsCmds builds a tools command tree whose help examples are written
// with prefix (e.g. "rsk tools"). Every subcommand carries its own --for flag
// so the TUI parameter form and row actions see it.
func newToolsCmds(prefix string) toolsCmds {
	root := &cobra.Command{
		Use:   "tools",
		Short: "Manage tool permissions (allow/deny/list/remove).",
		Long: fmt.Sprintf(`Manage which tools (Bash, Read, Write, etc.) are allowed or denied in
this project. For Claude Code, changes are written to .claude/settings.json
and override global tool permissions.

Use --for to pick the AI tool (default claude-code; other tools are not
supported yet).

Examples:
  %[1]s list
  %[1]s allow Bash(npm run *)
  %[1]s deny Write(**)
  %[1]s remove Bash(npm run *)`, prefix),
	}

	list := &cobra.Command{
		Use:   "list",
		Short: "List current tool permissions for this project.",
		Long: fmt.Sprintf(`List current tool permissions for this project.

Examples:
  %[1]s list
  %[1]s list -o json`, prefix),
		RunE: runClaudeToolsList,
	}
	allow := &cobra.Command{
		Use:   "allow [rule]",
		Short: "Allow a tool in this project.",
		Long: `Add a tool rule to the permissions.allow list in .claude/settings.json.
The rule format is Tool(specifier), for example:
  Bash(npm run *)
  Read(~/docs/**)
  WebFetch(domain:example.com)`,
		RunE: runClaudeToolsAllow,
	}
	deny := &cobra.Command{
		Use:   "deny [rule]",
		Short: "Deny a tool in this project.",
		Long: `Add a tool rule to the permissions.deny list in .claude/settings.json.
The rule format is Tool(specifier), for example:
  Write(**)
  Bash
  WebFetch`,
		RunE: runClaudeToolsDeny,
	}
	remove := &cobra.Command{
		Use:   "remove [rule]",
		Short: "Remove a tool rule from permissions.",
		Long:  `Remove a tool rule from either the allow or deny list.`,
		RunE:  runClaudeToolsRemove,
	}

	root.AddCommand(list, allow, deny, remove)
	for _, c := range []*cobra.Command{list, allow, deny, remove} {
		c.Flags().String(cmdx.FlagFor, string(tool.ClaudeID), toolsForHint)
	}
	addOutputFlag(list)
	return toolsCmds{root: root, list: list, allow: allow, deny: deny, remove: remove}
}

// requireClaudeTarget validates the --for flag of a tools subcommand. Only
// Claude Code is supported today: other known tools (or "all") get a "not
// supported yet" error, anything else an unknown-tool error.
func requireClaudeTarget(cmd *cobra.Command) error {
	target := cmdx.String(cmd, cmdx.FlagFor)
	if target == "" {
		return nil
	}
	scope, err := cmdx.ParseTargetScope(target)
	if err != nil {
		return fmt.Errorf("unknown tool %q for --for (supported: %s)", target, tool.ClaudeID)
	}
	if scope != cmdx.TargetScope(tool.ClaudeID) {
		return fmt.Errorf("--for %s is not supported yet: tools subcommands only support %s", target, tool.ClaudeID)
	}
	return nil
}

func claudeToolGet() (*tool.ClaudeTool, error) {
	t, ok := tool.Get(tool.ClaudeID)
	if !ok {
		return nil, fmt.Errorf("tool %s not registered", tool.ClaudeID)
	}
	ct, ok := t.(*tool.ClaudeTool)
	if !ok {
		return nil, fmt.Errorf("tool %s has unexpected type %T", tool.ClaudeID, t)
	}
	return ct, nil
}

// claudeToolsPermissions is the JSON shape for `rsk claude tools list -o json`
// — the raw allow/deny rule lists as configured in .claude/settings.json.
type claudeToolsPermissions struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

func runClaudeToolsList(cmd *cobra.Command, args []string) error {
	if err := checkOutput(cmd); err != nil {
		return err
	}
	if err := requireClaudeTarget(cmd); err != nil {
		return err
	}
	format := printer.Format(cmd)

	cwd, err := manifest.ProjectFolderPath()
	if err != nil {
		return err
	}
	projectDir := cwd[:len(cwd)-len("/.rsk")]

	claudeTool, err := claudeToolGet()
	if err != nil {
		return err
	}
	allow, deny, err := claudeTool.ReadPermissions(projectDir)
	if err != nil {
		return err
	}

	if format == output.FormatJSON {
		return newFormatter(cmd, format).JSON(claudeToolsPermissions{Allow: allow, Deny: deny})
	}

	rows := tool.ClaudeToolStates(allow, deny)

	header := []string{headerName, "Status"}
	tableRows := term.Rows(rows, claudeToolsTable)

	return renderList(cmd, format, "Claude Code tools available:", rows, header, tableRows, tool.AvailableClaudeTools())
}

// claudeToolsTable projects one tool.State, marking each tool's
// permission state the same way install/uninstall mark their own rows.
var claudeToolsTable = term.Table[tool.State]{
	Row: func(r tool.State) []any {
		switch r.Permission {
		case tool.Allowed:
			return []any{r.Name, ui.SuccessMark + " explicitly allowed"}
		case tool.Denied:
			return []any{r.Name, ui.ErrorMark + " explicitly denied"}
		default:
			return []any{r.Name, ui.MutedStyle.Render("• not configured (uses global settings)")}
		}
	},
}

// claudeRuleTarget resolves what the allow, deny and remove subcommands need:
// the validated --for target, the rule (from the argument or a prompt) and the
// project directory with the Claude tool that owns its settings.
func claudeRuleTarget(
	cmd *cobra.Command, args []string, prompt string,
) (claudeTool *tool.ClaudeTool, projectDir, rule string, err error) {
	if err = requireClaudeTarget(cmd); err != nil {
		return nil, "", "", err
	}
	rule, err = nameFromArgsOrPrompt(cmd, args, fieldArgRule, prompt)
	if err != nil {
		return nil, "", "", err
	}
	projectDir, err = claudeProjectDir()
	if err != nil {
		return nil, "", "", err
	}
	claudeTool, err = claudeToolGet()
	if err != nil {
		return nil, "", "", err
	}
	return claudeTool, projectDir, rule, nil
}

// claudeProjectDir is the project root: the parent of the .rsk folder.
func claudeProjectDir() (string, error) {
	cwd, err := manifest.ProjectFolderPath()
	if err != nil {
		return "", err
	}
	return cwd[:len(cwd)-len("/.rsk")], nil
}

func runClaudeToolsAllow(cmd *cobra.Command, args []string) error {
	claudeTool, projectDir, rule, err := claudeRuleTarget(cmd, args, "Tool rule (e.g. Bash(npm run *))")
	if err != nil {
		return err
	}
	return reportClaudeRule(cmd, rule, " is already allowed", "allowed ")(claudeTool.AllowRule(projectDir, rule))
}

func runClaudeToolsDeny(cmd *cobra.Command, args []string) error {
	claudeTool, projectDir, rule, err := claudeRuleTarget(cmd, args, "Tool rule (e.g. Write(**) or Bash)")
	if err != nil {
		return err
	}
	return reportClaudeRule(cmd, rule, " is already denied", "denied ")(claudeTool.DenyRule(projectDir, rule))
}

func runClaudeToolsRemove(cmd *cobra.Command, args []string) error {
	claudeTool, projectDir, rule, err := claudeRuleTarget(cmd, args, "Tool rule to remove")
	if err != nil {
		return err
	}
	return reportClaudeRule(cmd, rule, " not found in permissions", "removed ")(claudeTool.RemoveRule(projectDir, rule))
}

// reportClaudeRule prints the outcome of a rule change: the unchanged message
// when nothing was written, the success line framed by blank lines otherwise.
func reportClaudeRule(cmd *cobra.Command, rule, unchanged, done string) func(bool, error) error {
	out := cmd.OutOrStdout()
	return func(changed bool, err error) error {
		if err != nil {
			return err
		}
		if !changed {
			ui.Info(out, rule+unchanged)
			return nil
		}
		fmt.Fprintln(out)
		ui.Success(out, done+rule)
		fmt.Fprintln(out)
		return nil
	}
}
