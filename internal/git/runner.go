package git

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const maxStderr = 16 << 10

var (
	urlUserinfo = regexp.MustCompile(`(?i)\b(https?|ssh|git)://[^/@\s'"]+@`)
	tokenLike   = regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{20,}|github_pat_[A-Za-z0-9_]{20,}|glpat-[A-Za-z0-9_-]{10,}`)

	// hardening is every setting that does not depend on the runner's mode.
	hardening = []string{
		"protocol.allow=never",
		"protocol.ssh.allow=always",
		"protocol.https.allow=always",
		"protocol.ext.allow=never",
		"submodule.recurse=false",
		"fetch.recurseSubmodules=false",
		"transfer.fsckObjects=true",
		"core.hooksPath=/dev/null",
	}

	// droppedEnv are the variables the runner sets itself or refuses to inherit:
	// an askpass program could open a prompt, and the ssh command is rebuilt.
	droppedEnv = []string{
		"GIT_TERMINAL_PROMPT=", "GIT_SSH_COMMAND=", "GIT_ASKPASS=", "SSH_ASKPASS=",
		"SSH_ASKPASS_REQUIRE=", "LC_ALL=",
	}
)

// Error is a failed git command. Stderr is scrubbed of credentials before it is
// kept, so the error is safe to print and to log.
type Error struct {
	Op     string
	Stderr string
	Err    error
}

// Error implements error.
func (e *Error) Error() string {
	detail := strings.Join(strings.Fields(strings.ReplaceAll(strings.TrimSpace(e.Stderr), "\n", "; ")), " ")
	if detail == "" {
		return fmt.Sprintf("git %s: %v", e.Op, e.Err)
	}
	return fmt.Sprintf("git %s: %s: %v", e.Op, detail, e.Err)
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.Err
}

// scrub removes what could be a credential from text git printed: userinfo in a
// URL and the shapes of common hosting tokens.
func scrub(s string) string {
	s = urlUserinfo.ReplaceAllString(s, "${1}://***@")
	return tokenLike.ReplaceAllString(s, "***")
}

// runner builds and runs git commands with the hardening every command that
// talks to a remote needs. The zero value is the production configuration;
// allowFile exists only so in-package tests can fetch from a local repository.
type runner struct {
	allowFile bool
}

// command builds the git invocation. It refuses every transport but ssh and
// https, never recurses into submodules, validates the objects it receives,
// runs no hooks, never prompts, and closes stdin.
func (r runner) command(ctx context.Context, dir, sshBase string, args ...string) *exec.Cmd {
	full := make([]string, 0, len(args)+2*(len(hardening)+1)+2)
	if dir != "" {
		full = append(full, "-C", dir)
	}
	fileProtocol := "protocol.file.allow=never"
	if r.allowFile {
		fileProtocol = "protocol.file.allow=always"
	}
	for _, setting := range append([]string{fileProtocol}, hardening...) {
		full = append(full, "-c", setting)
	}
	full = append(full, args...)

	cmd := exec.CommandContext(ctx, "git", full...)
	cmd.Env = append(filteredEnv(os.Environ()),
		"GIT_TERMINAL_PROMPT=0",
		"GIT_SSH_COMMAND="+sshBase+" -o BatchMode=yes",
		"GIT_ASKPASS=",
		"SSH_ASKPASS_REQUIRE=never",
		"LC_ALL=C",
	)
	return cmd
}

func filteredEnv(environ []string) []string {
	out := make([]string, 0, len(environ))
	for _, kv := range environ {
		dropped := false
		for _, prefix := range droppedEnv {
			if strings.HasPrefix(kv, prefix) {
				dropped = true
				break
			}
		}
		if !dropped {
			out = append(out, kv)
		}
	}
	return out
}

// baseSSHCommand keeps the user's own ssh command when they set one: the
// GIT_SSH_COMMAND environment variable first, then core.sshCommand, then ssh.
func baseSSHCommand(env, configured string) string {
	switch {
	case env != "":
		return env
	case configured != "":
		return configured
	default:
		return "ssh"
	}
}

func configuredSSHCommand(ctx context.Context) string {
	out, err := exec.CommandContext(ctx, "git", "config", "--get", "core.sshCommand").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (r runner) prepare(ctx context.Context, dir string, args ...string) *exec.Cmd {
	return r.command(ctx, dir, baseSSHCommand(os.Getenv("GIT_SSH_COMMAND"), configuredSSHCommand(ctx)), args...)
}

// run executes a command and returns its stdout. stderr is captured, capped and
// scrubbed instead of being forwarded to the terminal.
func (r runner) run(ctx context.Context, dir, op string, args ...string) ([]byte, error) {
	cmd := r.prepare(ctx, dir, args...)
	var stdout bytes.Buffer
	stderr := &cappedBuffer{limit: maxStderr}
	cmd.Stdout = &stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return nil, &Error{Op: op, Stderr: scrub(stderr.String()), Err: err}
	}
	return stdout.Bytes(), nil
}

// stream runs a command and hands its stdout to read. If read fails the
// process is stopped and read's error is returned.
func (r runner) stream(ctx context.Context, dir, op string, args []string, read func(io.Reader) error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd := r.prepare(ctx, dir, args...)
	stderr := &cappedBuffer{limit: maxStderr}
	cmd.Stderr = stderr
	out, err := cmd.StdoutPipe()
	if err != nil {
		return &Error{Op: op, Err: err}
	}
	if err = cmd.Start(); err != nil {
		return &Error{Op: op, Err: err}
	}

	readErr := read(out)
	if readErr == nil {
		_, readErr = io.Copy(io.Discard, out)
	}
	if readErr != nil {
		cancel()
	}
	waitErr := cmd.Wait()
	if readErr != nil {
		return readErr
	}
	if waitErr != nil {
		return &Error{Op: op, Stderr: scrub(stderr.String()), Err: waitErr}
	}
	return nil
}

// cappedBuffer keeps the first limit bytes written and discards the rest, so a
// chatty failure cannot grow memory.
type cappedBuffer struct {
	buf   bytes.Buffer
	limit int
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	if room := b.limit - b.buf.Len(); room > 0 {
		b.buf.Write(p[:min(len(p), room)])
	}
	return len(p), nil
}

func (b *cappedBuffer) String() string {
	return b.buf.String()
}
