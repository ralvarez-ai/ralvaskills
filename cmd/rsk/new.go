package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/term"

	"go.ralva.dev/ralvaskills/v3/internal/cmdx"
	"go.ralva.dev/ralvaskills/v3/internal/project"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Initialize an rsk project in the current directory.",
	Long: fmt.Sprintf(`Create .rsk/rsk.mod and tool-specific config in the current directory so
this folder becomes an rsk project. Skill directories are created
lazily the first time a skill is installed. Commit .rsk/rsk.mod,
.rsk/rsk.lock and .rsk/CLAUDE.md: the lock resolves the same versions on
every machine.

Claude Code also gets an agent guide: AGENTS.md is created from a short
skeleton when the project has none, and CLAUDE.md imports it (@AGENTS.md) above
the generated pin file. An existing AGENTS.md is never modified.

Use --for to select which tools to configure, from %s, or all:

  claude-code  AGENTS.md + CLAUDE.md (imports it and .rsk/CLAUDE.md)
  opencode     Writes pinned skills to opencode.json instructions
  others       Discover skills in .agents/skills/, the shared Agent Skills path

Examples:
  rsk new
  rsk new --for opencode
  rsk new --for all`, strings.Join(tool.Names(), ", ")),

	RunE: runNew,
}

func toolsFromFlag(flag string) ([]tool.ID, error) {
	scope, err := cmdx.ParseTargetScope(flag)
	if err != nil {
		return nil, fieldError(cmdx.FlagFor, "%v", err)
	}
	if scope == cmdx.ScopeAll {
		tools := tool.All()
		ids := make([]tool.ID, len(tools))
		for i, t := range tools {
			ids[i] = t.ID()
		}
		return ids, nil
	}
	id, parseErr := tool.ParseID(string(scope))
	if parseErr != nil {
		return nil, fieldError(cmdx.FlagFor, "--for must be a supported tool or %s; got %q", cmdx.ForAll, flag)
	}
	return []tool.ID{id}, nil
}

func runNew(cmd *cobra.Command, _ []string) error {
	out := cmd.OutOrStdout()

	forFlag, err := resolveForFlag(cmd, out)
	if err != nil {
		return err
	}

	tools, err := toolsFromFlag(forFlag)
	if err != nil {
		return err
	}
	registerMCP := resolveMCPFlag(cmd)

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	res, err := project.Setup(cmd.Context(), cwd, project.SetupOptions{
		Tools:        tools,
		ForLabel:     forFlag,
		RegisterMCP:  registerMCP,
		InstallGuide: resolveGuideFlag(cmd),
	})
	printNotices(out, res.Notices)
	if err != nil {
		return err
	}

	fmt.Fprintln(out)
	if res.Created {
		ui.Success(out, "initialized .rsk/ for "+forFlag)
	} else {
		ui.Success(out, "updated .rsk/ for "+forFlag)
	}
	ui.Indent(out, "Run 'rsk install <name>' to add skills to this project.")
	fmt.Fprintln(out)
	return nil
}

// printNotices shows the messages a project operation produced, in order.
func printNotices(out io.Writer, notices []project.Notice) {
	for _, n := range notices {
		switch n.Level {
		case project.Warn:
			ui.Warn(out, n.Text)
		case project.Success:
			ui.Success(out, n.Text)
		default:
			ui.Info(out, n.Text)
		}
	}
}

// resolveForFlag returns --for's value, prompting interactively when it
// wasn't set on the command line.
func resolveForFlag(cmd *cobra.Command, out io.Writer) (string, error) {
	forFlag := cmdx.String(cmd, cmdx.FlagFor)
	if cmd.Flags().Changed(cmdx.FlagFor) {
		return forFlag, nil
	}
	if ui.InSession(cmd.Context()) {
		return "", fieldError(cmdx.FlagFor, "choose which tools to configure")
	}

	choices := append(tool.Names(), cmdx.ForAll)
	idx, err := ui.Select(out, "Tools to configure", choices, 0)
	if err != nil {
		return "", err
	}
	return choices[idx], nil
}

// flagMCP is the boolean flag that registers the rsk MCP server without a
// prompt, for non-interactive use.
const flagMCP = "mcp"

// resolveMCPFlag returns whether to register the rsk MCP server: the --mcp flag
// when set, an interactive prompt otherwise.
func resolveMCPFlag(cmd *cobra.Command) bool {
	if cmd.Flags().Changed(flagMCP) {
		return cmdx.Bool(cmd, flagMCP)
	}
	if ui.InSession(cmd.Context()) {
		return false
	}
	return term.Confirm(cmd, false, term.StdioIsTerminal(), "Register the rsk MCP server in this project?")
}

// flagGuide is the boolean flag that installs rsk-guide without a prompt.
const flagGuide = "guide"

// resolveGuideFlag returns whether to install rsk-guide: the --guide flag when
// set, an interactive prompt otherwise.
func resolveGuideFlag(cmd *cobra.Command) bool {
	if cmd.Flags().Changed(flagGuide) {
		return cmdx.Bool(cmd, flagGuide)
	}
	if ui.InSession(cmd.Context()) {
		return false
	}
	return term.Confirm(cmd, false, term.StdioIsTerminal(), "Install the rsk-guide skill in this project?")
}
