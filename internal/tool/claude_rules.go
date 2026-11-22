package tool

import (
	"slices"
)

// availableClaudeTools lists the Claude Code tools whose permission state
// `rsk tools list` reports.
var availableClaudeTools = []string{
	"Agent",
	"AskUserQuestion",
	"Bash",
	"CronCreate",
	"CronDelete",
	"CronList",
	"Edit",
	"EnterPlanMode",
	"EnterWorktree",
	"ExitPlanMode",
	"ExitWorktree",
	"Glob",
	"Grep",
	"ListMcpResourcesTool",
	"LSP",
	"Monitor",
	"NotebookEdit",
	"PowerShell",
	"PushNotification",
	"Read",
	"ReadMcpResourceTool",
	"RemoteTrigger",
	"ScheduleWakeup",
	"SendMessage",
	"ShareOnboardingGuide",
	"Skill",
	"TaskCreate",
	"TaskGet",
	"TaskList",
	"TaskOutput",
	"TaskStop",
	"TaskUpdate",
	"TeamCreate",
	"TeamDelete",
	"TodoWrite",
	"ToolSearch",
	"WaitForMcpServers",
	"WebFetch",
	"WebSearch",
	"Write",
}

// Permission is a Claude tool's permission state relative to a project's
// settings.
type Permission int

const (
	// Allowed means the bare tool name is in the allow list.
	Allowed Permission = iota
	// Denied means the bare tool name is in the deny list.
	Denied
	// Unconfigured means neither list names the tool, so global settings apply.
	Unconfigured
)

// State is one Claude tool and its permission state.
type State struct {
	Name       string
	Permission Permission
}

// AvailableClaudeTools returns the Claude tool names `rsk tools list` reports,
// in display order. The slice is a copy.
func AvailableClaudeTools() []string {
	return slices.Clone(availableClaudeTools)
}

// ClaudeToolStates reports the permission state of every available Claude tool
// given the project's allow and deny rules, in display order. An allow rule
// wins over a deny rule for the same name.
func ClaudeToolStates(allow, deny []string) []State {
	states := make([]State, 0, len(availableClaudeTools))
	for _, name := range availableClaudeTools {
		switch {
		case slices.Contains(allow, name):
			states = append(states, State{Name: name, Permission: Allowed})
		case slices.Contains(deny, name):
			states = append(states, State{Name: name, Permission: Denied})
		default:
			states = append(states, State{Name: name, Permission: Unconfigured})
		}
	}
	return states
}

// AllowRule adds rule to the project's allow list and drops it from the deny
// list. It reports false, without writing, when rule is already allowed.
func (ct *ClaudeTool) AllowRule(projectDir, rule string) (bool, error) {
	allow, deny, err := ct.ReadPermissions(projectDir)
	if err != nil {
		return false, err
	}
	if slices.Contains(allow, rule) {
		return false, nil
	}
	deny = removeRule(deny, rule)
	allow = append(allow, rule)
	if err = ct.WritePermissions(projectDir, allow, deny); err != nil {
		return false, err
	}
	return true, nil
}

// DenyRule adds rule to the project's deny list and drops it from the allow
// list. It reports false, without writing, when rule is already denied.
func (ct *ClaudeTool) DenyRule(projectDir, rule string) (bool, error) {
	allow, deny, err := ct.ReadPermissions(projectDir)
	if err != nil {
		return false, err
	}
	if slices.Contains(deny, rule) {
		return false, nil
	}
	allow = removeRule(allow, rule)
	deny = append(deny, rule)
	if err = ct.WritePermissions(projectDir, allow, deny); err != nil {
		return false, err
	}
	return true, nil
}

// RemoveRule drops rule from both lists. It reports false, without writing,
// when neither list held it.
func (ct *ClaudeTool) RemoveRule(projectDir, rule string) (bool, error) {
	allow, deny, err := ct.ReadPermissions(projectDir)
	if err != nil {
		return false, err
	}
	initialAllow, initialDeny := len(allow), len(deny)
	allow = removeRule(allow, rule)
	deny = removeRule(deny, rule)
	if len(allow) == initialAllow && len(deny) == initialDeny {
		return false, nil
	}
	if err = ct.WritePermissions(projectDir, allow, deny); err != nil {
		return false, err
	}
	return true, nil
}

func removeRule(rules []string, rule string) []string {
	return slices.DeleteFunc(rules, func(v string) bool { return v == rule })
}
