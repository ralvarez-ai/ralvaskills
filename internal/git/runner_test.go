package git

import (
	"slices"
	"strings"
	"testing"
)

func TestScrubRedactsCredentials(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"fatal: unable to access 'https://user:hunter2@github.com/acme/x.git/': 403": "fatal: unable to access 'https://***@github.com/acme/x.git/': 403",
		"remote: token ghp_abcdefghijklmnopqrstuvwxyz0123456789 rejected":            "remote: token *** rejected",
		"auth github_pat_11AAAAAAA0abcdefghijklmnopqrstuvwxyz failed":                "auth *** failed",
		"glpat-abcdefghij1234567890 bad":                                             "*** bad",
		"ssh://git:secret@host/x":                                                    "ssh://***@host/x",
		"Permission denied (publickey).":                                             "Permission denied (publickey).",
	}
	for in, want := range tests {
		if got := scrub(in); got != want {
			t.Errorf("scrub(%q)\n got %q\nwant %q", in, got, want)
		}
	}
}

func TestErrorCarriesScrubbedStderrOnly(t *testing.T) {
	t.Parallel()
	err := &Error{Op: "fetch", Stderr: scrub("fatal: https://u:p@h/x: denied\nmore"), Err: errExit}
	msg := err.Error()
	if strings.Contains(msg, "u:p@") {
		t.Errorf("error leaks credentials: %q", msg)
	}
	if !strings.Contains(msg, "git fetch") || !strings.Contains(msg, "denied") {
		t.Errorf("error lacks op or reason: %q", msg)
	}
}

var errExit = &exitStatusError{}

type exitStatusError struct{}

func (*exitStatusError) Error() string {
	return "exit status 128"
}

func TestCommandIsHardened(t *testing.T) {
	t.Parallel()
	r := runner{}
	cmd := r.command(t.Context(), "/tmp/repo.git", "base-ssh", "fetch", "--", "git@github.com:a/b.git", "main")

	args := strings.Join(cmd.Args, " ")
	for _, want := range []string{
		"-C /tmp/repo.git",
		"protocol.ext.allow=never",
		"protocol.file.allow=never",
		"protocol.allow=never",
		"protocol.ssh.allow=always",
		"protocol.https.allow=always",
		"submodule.recurse=false",
		"transfer.fsckObjects=true",
		"core.hooksPath=/dev/null",
		"fetch -- git@github.com:a/b.git main",
	} {
		if !strings.Contains(args, want) {
			t.Errorf("args lack %q:\n%s", want, args)
		}
	}
	if cmd.Stdin != nil {
		t.Error("stdin must be closed (nil means /dev/null)")
	}

	env := strings.Join(cmd.Env, "\n")
	for _, want := range []string{"GIT_TERMINAL_PROMPT=0", "GIT_SSH_COMMAND=base-ssh -o BatchMode=yes"} {
		if !strings.Contains(env, want) {
			t.Errorf("env lacks %q", want)
		}
	}
	for _, banned := range []string{"GIT_ASKPASS=", "SSH_ASKPASS="} {
		for _, e := range cmd.Env {
			if strings.HasPrefix(e, banned) && e != banned {
				t.Errorf("env keeps %q: an askpass program could prompt", e)
			}
		}
	}
}

func TestCommandAllowsFileOnlyForTheTestHook(t *testing.T) {
	t.Parallel()
	hooked := runner{allowFile: true}.command(t.Context(), "/r", "ssh", "fetch")
	if !slices.Contains(hooked.Args, "protocol.file.allow=always") {
		t.Error("the in-package test hook should allow the file protocol")
	}
	if slices.Contains(runner{}.command(t.Context(), "/r", "ssh", "fetch").Args, "protocol.file.allow=always") {
		t.Error("the default runner must never allow the file protocol")
	}
}

func TestSSHCommandKeepsTheUsersChoice(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, env, configured, want string
	}{
		{"nothing set", "", "", "ssh"},
		{"configured", "", "ssh -i /k", "ssh -i /k"},
		{"env wins", "ssh -F /c", "ssh -i /k", "ssh -F /c"},
	}
	for _, tc := range tests {
		if got := baseSSHCommand(tc.env, tc.configured); got != tc.want {
			t.Errorf("%s: baseSSHCommand = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestParseVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		want    [3]int
		wantErr bool
	}{
		{"git version 2.56.1", [3]int{2, 56, 1}, false},
		{"git version 2.39.2 (Apple Git-143)", [3]int{2, 39, 2}, false},
		{"git version 2.30.0.windows.1", [3]int{2, 30, 0}, false},
		{"git version 2.5", [3]int{2, 5, 0}, false},
		{"not git", [3]int{}, true},
	}
	for _, tc := range tests {
		got, err := parseVersion(tc.in)
		if (err != nil) != tc.wantErr || got != tc.want {
			t.Errorf("parseVersion(%q) = %v, %v; want %v, err=%v", tc.in, got, err, tc.want, tc.wantErr)
		}
	}
	if !versionAtLeast([3]int{2, 30, 0}, MinVersion) || versionAtLeast([3]int{2, 29, 9}, MinVersion) {
		t.Errorf("versionAtLeast disagrees with MinVersion %v", MinVersion)
	}
}
