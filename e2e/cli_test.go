// Package e2e exercises the built rsk binary. The unit tests in cmd/rsk call
// the command functions in process, so they cannot see exit codes, the
// stdout/stderr split, or the files an invocation leaves on disk; these few
// journeys are the only place those contracts are proven.
//
// The tier runs from its own command (`task test:e2e`) and its own CI job. It
// needs no network and no Docker: each test works in a temporary directory with
// RSK_CONFIG_HOME and HOME pointed inside it, so nothing touches the machine.
package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// binary is the rsk executable built once by TestMain. Empty under -short,
// where every test skips instead of paying for the build.
var binary string

// repoRoot walks up from the test's working directory to the module root.
func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("go.mod not found above the test directory")
		}
		dir = parent
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}
	root, err := repoRoot()
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	dir, err := os.MkdirTemp("", "rsk-e2e-bin-")
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	binary = filepath.Join(dir, "rsk")
	build := exec.CommandContext(context.Background(), "go", "build", "-o", binary, "./cmd/rsk")
	build.Dir = root
	if out, buildErr := build.CombinedOutput(); buildErr != nil {
		os.RemoveAll(dir)
		os.Stderr.WriteString("building rsk: " + buildErr.Error() + "\n" + string(out))
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// result is one finished invocation: the two streams and the exit code.
type result struct {
	stdout string
	stderr string
	code   int
}

// project is a temporary directory with an isolated rsk config and home.
type project struct {
	dir string
}

func newProject(t *testing.T) *project {
	t.Helper()
	dir := t.TempDir()
	for _, sub := range []string{"config", "home"} {
		if err := os.MkdirAll(filepath.Join(dir, sub), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", sub, err)
		}
	}

	return &project{dir: dir}
}

// run invokes the binary inside the project with its streams captured.
func (p *project) run(t *testing.T, args ...string) result {
	t.Helper()
	if binary == "" {
		t.Skip("the e2e tier does not run under -short")
	}
	cmd := exec.CommandContext(t.Context(), binary, args...)
	cmd.Dir = p.dir
	cmd.Env = append(os.Environ(),
		"RSK_CONFIG_HOME="+filepath.Join(p.dir, "config"),
		"HOME="+filepath.Join(p.dir, "home"),
		"NO_COLOR=1",
		"TERM=dumb",
	)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		code = exitErr.ExitCode()
	} else if err != nil {
		t.Fatalf("running rsk %s: %v", strings.Join(args, " "), err)
	}

	return result{stdout: stdout.String(), stderr: stderr.String(), code: code}
}

// read returns a file's contents, failing the test when it is missing.
func (p *project) read(t *testing.T, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(p.dir, rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}

	return string(data)
}

func (p *project) exists(rel string) bool {
	_, err := os.Stat(filepath.Join(p.dir, rel))

	return err == nil
}

func TestVersionGoesToStdout(t *testing.T) {
	t.Parallel()
	res := newProject(t).run(t, "--version")
	if res.code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr: %s)", res.code, res.stderr)
	}
	if !strings.HasPrefix(res.stdout, "rsk ") || strings.TrimSpace(res.stdout) == "rsk" {
		t.Errorf("stdout = %q, want a version after 'rsk '", res.stdout)
	}
	if res.stderr != "" {
		t.Errorf("stderr = %q, want it empty", res.stderr)
	}
}

func TestHelpGoesToStdout(t *testing.T) {
	t.Parallel()
	res := newProject(t).run(t, "--help")
	if res.code != 0 {
		t.Fatalf("exit = %d, want 0", res.code)
	}
	if !strings.Contains(res.stdout, "rsk manages your local and official AI skills") {
		t.Errorf("stdout does not look like the usage text: %q", res.stdout)
	}
}

func TestUnknownCommandExitsNonZero(t *testing.T) {
	t.Parallel()
	res := newProject(t).run(t, "frobnicate")
	if res.code == 0 {
		t.Fatalf("exit = 0, want a failure for an unknown command")
	}
	if !strings.Contains(res.stderr, "unknown command") {
		t.Errorf("stderr = %q, want it to name the unknown command", res.stderr)
	}
}

func TestNewOpencodeIsIdempotent(t *testing.T) {
	t.Parallel()
	p := newProject(t)
	first := p.run(t, "new", "--for", "opencode")
	if first.code != 0 {
		t.Fatalf("first new: exit = %d (stderr: %s)", first.code, first.stderr)
	}
	mod := p.read(t, ".rsk/rsk.mod")
	if !strings.Contains(mod, `tools = ["opencode"]`) {
		t.Errorf("rsk.mod does not scope the project to opencode:\n%s", mod)
	}
	if !strings.Contains(p.read(t, "opencode.json"), `"rsk"`) {
		t.Errorf("opencode.json does not register the rsk MCP server")
	}

	second := p.run(t, "new", "--for", "opencode")
	if second.code != 0 {
		t.Fatalf("second new: exit = %d, want a no-op (stderr: %s)", second.code, second.stderr)
	}
	if again := p.read(t, ".rsk/rsk.mod"); again != mod {
		t.Errorf("a second new rewrote rsk.mod:\nfirst:\n%s\nsecond:\n%s", mod, again)
	}
}

func TestNewClaudeCodeImportsAgents(t *testing.T) {
	t.Parallel()
	p := newProject(t)
	if res := p.run(t, "new", "--for", "claude-code"); res.code != 0 {
		t.Fatalf("exit = %d (stderr: %s)", res.code, res.stderr)
	}
	if !p.exists("AGENTS.md") {
		t.Error("new --for claude-code did not seed AGENTS.md")
	}
	claude := p.read(t, "CLAUDE.md")
	for _, want := range []string{"@AGENTS.md", "@.rsk/CLAUDE.md"} {
		if !strings.Contains(claude, want) {
			t.Errorf("CLAUDE.md does not import %s:\n%s", want, claude)
		}
	}
}

func TestListJSONIsJSONEvenWhenEmpty(t *testing.T) {
	t.Parallel()
	p := newProject(t)
	if res := p.run(t, "new", "--for", "opencode"); res.code != 0 {
		t.Fatalf("new: exit = %d (stderr: %s)", res.code, res.stderr)
	}
	res := p.run(t, "list", "-o", "json")
	if res.code != 0 {
		t.Fatalf("list: exit = %d (stderr: %s)", res.code, res.stderr)
	}
	var rows []map[string]any
	if err := json.Unmarshal([]byte(res.stdout), &rows); err != nil {
		t.Fatalf("stdout is not a JSON array: %v\nstdout: %q\nstderr: %q", err, res.stdout, res.stderr)
	}
	if len(rows) != 0 {
		t.Errorf("a fresh project listed %d skills, want 0", len(rows))
	}
}

func TestDestroyKeepsAgentsMd(t *testing.T) {
	t.Parallel()
	p := newProject(t)
	if res := p.run(t, "new", "--for", "claude-code"); res.code != 0 {
		t.Fatalf("new: exit = %d (stderr: %s)", res.code, res.stderr)
	}
	if res := p.run(t, "destroy", "--yes"); res.code != 0 {
		t.Fatalf("destroy: exit = %d (stderr: %s)", res.code, res.stderr)
	}
	if p.exists(".rsk") {
		t.Error("destroy left .rsk/ behind")
	}
	if !p.exists("AGENTS.md") {
		t.Error("destroy removed the hand-written AGENTS.md")
	}
	if claude := p.read(t, "CLAUDE.md"); !strings.Contains(claude, "@AGENTS.md") {
		t.Errorf("destroy dropped the @AGENTS.md import:\n%s", claude)
	}
}
