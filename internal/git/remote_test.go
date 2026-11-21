package git

import "testing"

func TestParseRemoteURL(t *testing.T) {
	t.Parallel()
	accept := []string{
		"git@github.com:acme/skills.git",
		"git@github.com:acme/skills",
		"deploy@git.example.com:team/repo.git",
		"ssh://git@github.com/acme/skills.git",
		"ssh://git@git.example.com:2222/team/repo.git",
		"ssh://git.example.com/team/repo.git",
		"https://github.com/acme/skills.git",
		"https://github.com/acme/skills",
		"https://git.example.com:8443/team/repo.git",
	}
	for _, in := range accept {
		got, err := ParseRemoteURL(in)
		if err != nil {
			t.Errorf("ParseRemoteURL(%q): %v", in, err)
			continue
		}
		if got.String() != in {
			t.Errorf("String() = %q, want %q", got.String(), in)
		}
	}

	reject := map[string]string{
		"empty":                    "",
		"ext transport":            "ext::sh -c touch% /tmp/pwned",
		"file scheme":              "file:///etc/passwd",
		"bare path":                "/srv/git/skills.git",
		"relative path":            "../skills.git",
		"plain http":               "http://github.com/acme/skills.git",
		"git protocol":             "git://github.com/acme/skills.git",
		"leading dash":             "--upload-pack=touch /tmp/pwned",
		"dash host in ssh url":     "ssh://-oProxyCommand=touch/tmp/pwned/x",
		"dash host in scp form":    "git@-oProxyCommand=x:path",
		"dash user":                "-oProxyCommand=x@host:path",
		"https credentials":        "https://user:ghp_token123@github.com/acme/skills.git",
		"https username only":      "https://user@github.com/acme/skills.git",
		"ssh password":             "ssh://git:secret@github.com/acme/skills.git",
		"no host":                  "https:///acme/skills.git",
		"no path":                  "https://github.com",
		"scp no path":              "git@github.com:",
		"newline":                  "https://github.com/acme/skills.git\nfoo",
		"space":                    "https://github.com/acme/ skills.git",
		"nul":                      "https://github.com/acme/\x00skills.git",
		"scp with scheme-looking":  "git@github.com:ext::sh",
		"unknown scheme":           "ftp://example.com/x.git",
		"transport helper":         "foo::https://example.com/x.git",
		"fragment in host":         "https://exa mple.com/x.git",
		"host with backslash":      "https://example.com\\evil/x.git",
		"userinfo trick in https":  "https://github.com@evil.example/acme/skills.git",
		"percent encoded dash ssh": "ssh://%2DoProxyCommand/x",
	}
	for name, in := range reject {
		if got, err := ParseRemoteURL(in); err == nil {
			t.Errorf("%s: ParseRemoteURL(%q) accepted as %q", name, in, got)
		}
	}
}
