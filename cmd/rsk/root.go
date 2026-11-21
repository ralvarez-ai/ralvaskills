// Command rsk manages local and official AI skills for Claude Code and OpenCode.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/lifecycle"
	"go.ralva.dev/lib/term/cli"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

// Cobra group IDs for the root command's sections, in help/picker order.
const (
	groupMachine   = "machine"
	groupProject   = "project"
	groupInstall   = "install"
	groupPinning   = "pinning"
	groupToolsConf = "tools"
	groupViews     = "views"
	groupOther     = "other"
)

// rootGroups are the root command's help sections, in display order. The TUI
// picker orders its top-level rows by the same list.
var rootGroups = []*cobra.Group{
	{ID: groupMachine, Title: "Machine setup:"},
	{ID: groupProject, Title: "Project lifecycle:"},
	{ID: groupInstall, Title: "Install / uninstall / update (use --global for system-wide):"},
	{ID: groupPinning, Title: "Project pinning (auto-import into each tool's project config):"},
	{ID: groupToolsConf, Title: "Tool configuration:"},
	{ID: groupViews, Title: "Read-only views:"},
	{ID: groupOther, Title: "Other:"},
}

var rootCmd = cli.NewRoot(cli.RootConfig{
	Use:   "rsk",
	Short: "Manage ralvaskills — install, update, and check AI skill bundles.",
	Long:  "rsk manages your local and official AI skills for Claude Code and OpenCode.",
})

// rootRunE is a bare `rsk`: the interactive session on a TTY, the banner and
// help otherwise.
func rootRunE(cmd *cobra.Command, _ []string) error {
	notice := updateNotice()
	maybeSpawnUpdateCheck()

	if !ui.IsTTY() {
		return printHome(cmd, notice)
	}
	return ui.RunSession(cmd.Context(), ui.SessionOptions{
		Root:     cmd.Root(),
		Runnable: sessionRunnable,
		Notice:   notice,
	})
}

// sessionRunnable offers a command in the interactive picker unless it is
// cobra's own help/completion machinery or `init`, whose setup wizard prompts
// on the terminal, which the session owns while it runs.
func sessionRunnable(cmd *cobra.Command) bool {
	top := cmd
	for top.HasParent() && top.Parent().HasParent() {
		top = top.Parent()
	}
	switch top.Name() {
	case "help", "completion", mcpCmd.Name(), initCmd.Name():
		return false
	}
	return true
}

// printHome prints the plain banner + cobra help shown for a bare `rsk`
// invocation outside a TTY (pipes, CI, scripts). notice, when non-empty, is
// an updates-available flag printed above the banner.
func printHome(cmd *cobra.Command, notice string) error {
	if notice != "" {
		ui.Warn(cmd.OutOrStdout(), notice)
	}
	if _, err := fmt.Fprintln(cmd.OutOrStdout(), ui.Banner()); err != nil {
		return fmt.Errorf("write banner: %w", err)
	}
	return cmd.Help()
}

//nolint:gochecknoinits // init() used for root command setup
func init() {
	rootCmd.RunE = rootRunE

	pf := rootCmd.PersistentFlags()
	pf.Bool("no-color", false, "disable color output")
	pf.String("config", "", "rsk config directory (default: $RSK_CONFIG_HOME or ~/.config/rsk)")
	pf.BoolP("verbose", "v", false, "print diagnostic detail")
	pf.BoolP("quiet", "q", false, "suppress informational output")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		noColor, err := cmd.Flags().GetBool("no-color")
		if err != nil {
			return err
		}
		ui.SetNoColor(noColor)

		if dir, dirErr := cmd.Flags().GetString("config"); dirErr != nil {
			return dirErr
		} else if dir != "" {
			config.SetConfigHome(dir)
		}

		quiet, err := cmd.Flags().GetBool("quiet")
		if err != nil {
			return err
		}
		ui.SetQuiet(quiet)

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}
		ui.SetVerbose(verbose)
		ui.Debug(cmd.ErrOrStderr(), "config dir "+config.DefaultConfigFolderPath())

		return nil
	}
	cli.SetupVersion(rootCmd, "rsk")
	setupCommands()
}

// Execute runs the root command and exits on error.
// A user cancellation (Ctrl+C during a prompt) or a caught SIGINT/SIGTERM exits 130 silently per SIGINT convention.
func Execute() {
	os.Exit(execute())
}

func execute() int {
	ctx, stop := lifecycle.NotifyContext(context.Background())
	defer stop()

	err := rootCmd.ExecuteContext(ctx)
	if errors.Is(context.Cause(ctx), lifecycle.ErrInterrupted) || errors.Is(err, ui.ErrAborted) {
		return cli.ExitInterrupt
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "rsk: "+err.Error())
		return cli.ExitCode(err)
	}
	return 0
}
