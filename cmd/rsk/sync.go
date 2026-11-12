package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"go.ralva.dev/termkit"

	"go.ralva.dev/ralvaskills/v3/internal/cmdx"
	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/install"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
	"go.ralva.dev/ralvaskills/v3/internal/tool"
)

type syncLink struct {
	skill  skill.Skill
	target string
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Reconcile the project's installed skills and tool configuration.",
	Long: `Recreate missing skill links for the tools selected in rsk.mod, using
the exact versions in rsk.lock. Refresh pinned instructions. Use --mcp to also
register the rsk MCP server with each configured client.

Existing links to another target and user-owned files are reported as conflicts.
Locked registry skills must already be cached; sync does not update skills.`,
	Args: cobra.NoArgs,
	RunE: runSync,
}

func runSync(cmd *cobra.Command, _ []string) error {
	rskDir, err := manifest.ProjectFolderPath()
	if err != nil {
		return err
	}
	m, err := manifest.ReadMod(rskDir)
	if err != nil {
		return err
	}
	lock, err := manifest.ReadLock(rskDir)
	if err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load rsk config: %w", err)
	}
	project := filepath.Dir(rskDir)
	links, err := planSyncLinks(cmd.Context(), project, m, lock, cfg)
	if err != nil {
		return err
	}
	mcp, err := prepareSyncMCP(cmd, project, m.Tools)
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	for _, link := range links {
		fmt.Fprintf(out, "link %s -> %s\n", filepath.Join(link.target, link.skill.Name), link.skill.Path)
	}
	printSyncConfigPlan(out, m.Tools, mcp)
	if cmdx.Bool(cmd, cmdx.FlagDryRun) {
		return nil
	}
	if !termkit.Confirm(cmd, cmdx.Bool(cmd, cmdx.FlagYes), termkit.StdioIsTerminal(), "Proceed?") {
		fmt.Fprintln(out, messageAborted)
		return nil
	}
	for _, link := range links {
		if err = skill.Link(link.skill, link.target); err != nil {
			return err
		}
	}
	if err = install.SyncPinnedAllTools(rskDir, m); err != nil {
		return fmt.Errorf("sync pinned instructions: %w", err)
	}
	if mcp {
		if err = registerProjectMCP(out, project, m.Tools); err != nil {
			return fmt.Errorf("sync MCP registration: %w", err)
		}
	}
	fmt.Fprintln(out, "Project synchronization complete.")
	return nil
}

func printSyncConfigPlan(out io.Writer, ids []tool.ID, mcp bool) {
	fmt.Fprintln(out, "reconcile .rsk/SKILLS.md and pinned instructions for configured clients")
	for _, id := range ids {
		switch id {
		case tool.ClaudeID:
			fmt.Fprintln(out, "check AGENTS.md, CLAUDE.md, and .rsk/CLAUDE.md")
		case tool.OpenCodeID:
			fmt.Fprintln(out, "check opencode.json instructions")
		case tool.CodexID:
			fmt.Fprintln(out, "Codex discovers skills on demand")
		case tool.AntigravityID:
			fmt.Fprintln(out, "Antigravity discovers skills on demand")
		}
	}
	if mcp {
		for _, id := range ids {
			fmt.Fprintf(out, "check %s MCP registration", id)
			switch id {
			case tool.ClaudeID:
				fmt.Fprintln(out, " (.mcp.json and CLAUDE.md)")
			case tool.CodexID:
				fmt.Fprintln(out, " (.codex/config.toml)")
			case tool.OpenCodeID:
				fmt.Fprintln(out, " (opencode.json)")
			case tool.AntigravityID:
				fmt.Fprintln(out, " (.agents/mcp_config.json)")
			}
		}
	}
}

func prepareSyncMCP(cmd *cobra.Command, project string, ids []tool.ID) (bool, error) {
	mcp := cmdx.Bool(cmd, flagMCP)
	if !cmd.Flags().Changed(flagMCP) {
		var err error
		mcp, err = tool.ProjectHasMCP(project, ids)
		if err != nil {
			return false, err
		}
	}
	if !mcp {
		return false, nil
	}
	if err := checkSyncMCPExecutable(cmd.Context(), ids); err != nil {
		return false, err
	}
	if err := tool.PreflightMCP(project, ids); err != nil {
		return false, err
	}
	return true, nil
}

func checkSyncMCPExecutable(ctx context.Context, ids []tool.ID) error {
	path, err := exec.LookPath("rsk")
	if err != nil {
		return fmt.Errorf("rsk must be on PATH to sync MCP registration: %w", err)
	}
	for _, id := range ids {
		if id == tool.CodexID || id == tool.OpenCodeID {
			output, runErr := exec.CommandContext(ctx, path, "mcp", "--help").Output()
			if runErr != nil || !strings.Contains(string(output), "--read-only") {
				return fmt.Errorf("%s does not support read-only MCP mode; upgrade rsk before syncing MCP", path)
			}
			break
		}
	}
	return nil
}

func planSyncLinks(
	ctx context.Context,
	project string,
	m manifest.Mod,
	lock manifest.Lock,
	cfg config.Config,
) ([]syncLink, error) {
	locked := make(map[string]manifest.LockEntry, len(lock.Skills))
	for _, entry := range lock.Skills {
		if _, exists := locked[entry.Name]; exists {
			return nil, fmt.Errorf("duplicate lock entry for %q", entry.Name)
		}
		locked[entry.Name] = entry
	}
	names := make([]string, 0, len(m.Skills))
	for name := range m.Skills {
		names = append(names, name)
	}
	sort.Strings(names)
	var links []syncLink
	for _, name := range names {
		entry, ok := locked[name]
		if !ok {
			return nil, fmt.Errorf("%s has no lock entry; restore rsk.lock before syncing", name)
		}
		s, err := lockedSyncSkill(ctx, entry, cfg)
		if err != nil {
			return nil, err
		}
		for _, target := range install.ProjectSkillsDirs(project, m) {
			missing, checkErr := syncLinkMissing(filepath.Join(target, name), s.Path)
			if checkErr != nil {
				return nil, checkErr
			}
			if missing {
				links = append(links, syncLink{skill: s, target: target})
			}
		}
	}
	return links, nil
}

func lockedSyncSkill(ctx context.Context, entry manifest.LockEntry, cfg config.Config) (skill.Skill, error) {
	if err := skill.ValidateName(entry.Name); err != nil {
		return skill.Skill{}, err
	}
	if entry.Version == "" || filepath.Base(entry.Version) != entry.Version || entry.Version == "." ||
		entry.Version == ".." {
		return skill.Skill{}, fmt.Errorf("invalid locked version for %s", entry.Name)
	}
	var s skill.Skill
	var err error
	switch entry.Source {
	case skill.SourceRegistry:
		s = skill.Skill{Name: entry.Name, Version: entry.Version, Source: entry.Source,
			Path: filepath.Join(cfg.RegistryCache(), entry.Name, entry.Version)}
	case skill.SourceLocal:
		s, err = source.NewLocal(cfg.RepoPath).Find(ctx, entry.Name)
	case skill.SourceOfficial:
		if _, statErr := os.Stat(filepath.Join(cfg.OfficialCache, skill.SkillsFolderName)); statErr != nil {
			return skill.Skill{}, fmt.Errorf("official cache unavailable for %s: %w", entry.Name, statErr)
		}
		s, err = source.NewOfficial(cfg.OfficialCache).Find(ctx, entry.Name)
	default:
		return skill.Skill{}, fmt.Errorf("unsupported source for %s: %s", entry.Name, entry.Source)
	}
	if err != nil {
		return skill.Skill{}, err
	}
	actualVersion, err := skill.ReadVersion(s.Path)
	if err != nil || actualVersion != entry.Version {
		return skill.Skill{}, fmt.Errorf("locked skill %s@%s is unavailable at %s", entry.Name, entry.Version, s.Path)
	}
	return s, nil
}

func syncLinkMissing(path, want string) (bool, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	if !skill.IsLink(info) {
		return false, fmt.Errorf("%s exists and is not an rsk link", path)
	}
	got, err := filepath.EvalSymlinks(path)
	if err != nil {
		return false, fmt.Errorf("resolve %s: %w", path, err)
	}
	expected, err := filepath.EvalSymlinks(want)
	if err != nil {
		return false, err
	}
	if got != expected {
		return false, fmt.Errorf("%s points to %s, not locked skill %s", path, got, expected)
	}
	return false, nil
}
