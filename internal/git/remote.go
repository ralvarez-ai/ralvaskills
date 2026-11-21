package git

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ErrInvalidRemote is returned for a remote URL the runner will not accept.
var ErrInvalidRemote = errors.New("invalid git remote")

var (
	// scpPattern matches the scp-like form user@host:path. The user and host
	// cannot start with "-" (option injection into ssh) and the path cannot
	// hold another ":" (transport helpers such as ext::).
	scpPattern = regexp.MustCompile(`^[A-Za-z0-9._][A-Za-z0-9._-]*@[A-Za-z0-9][A-Za-z0-9.-]*:[^\s:\\]+$`)

	hostPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9.-]*$`)
	userPattern = regexp.MustCompile(`^[A-Za-z0-9._][A-Za-z0-9._-]*$`)
)

// RemoteURL is a git remote that passed ParseRemoteURL. The hardened runner
// takes only this type, so no unvalidated string reaches a git command line.
type RemoteURL struct {
	raw string
}

// String returns the URL as given. It never holds a credential: a URL that
// embeds one is refused at parse time.
func (u RemoteURL) String() string {
	return u.raw
}

// ParseRemoteURL accepts ssh://, https:// and scp-like user@host:path remotes
// and nothing else. It refuses ext:: and other transport helpers, file://, git://,
// plain http://, a value that starts with "-", whitespace or control
// characters, an https URL with any userinfo, an ssh URL with a password, and a
// host that is not a plain hostname.
func ParseRemoteURL(raw string) (RemoteURL, error) {
	if err := screenRemote(raw); err != nil {
		return RemoteURL{}, err
	}
	if !strings.Contains(raw, "://") {
		if !scpPattern.MatchString(raw) {
			return RemoteURL{}, fmt.Errorf("%w: %q is not ssh://, https:// or user@host:path", ErrInvalidRemote, raw)
		}
		return RemoteURL{raw: raw}, nil
	}

	u, err := url.Parse(raw)
	if err != nil {
		return RemoteURL{}, fmt.Errorf("%w: %w", ErrInvalidRemote, err)
	}
	if err = checkURLParts(u); err != nil {
		return RemoteURL{}, err
	}
	return RemoteURL{raw: raw}, nil
}

func screenRemote(raw string) error {
	switch {
	case raw == "":
		return fmt.Errorf("%w: empty", ErrInvalidRemote)
	case strings.HasPrefix(raw, "-"):
		return fmt.Errorf("%w: must not start with \"-\"", ErrInvalidRemote)
	case strings.Contains(raw, "::"):
		return fmt.Errorf("%w: transport helpers are not allowed", ErrInvalidRemote)
	}
	for _, r := range raw {
		if r <= ' ' || r == 0x7f || r == '\\' {
			return fmt.Errorf("%w: whitespace, control or backslash character", ErrInvalidRemote)
		}
	}
	return nil
}

func checkURLParts(u *url.URL) error {
	scheme := strings.ToLower(u.Scheme)
	if scheme != "ssh" && scheme != "https" {
		return fmt.Errorf("%w: scheme %q is not allowed", ErrInvalidRemote, u.Scheme)
	}
	if !hostPattern.MatchString(u.Hostname()) {
		return fmt.Errorf("%w: host %q is not a plain hostname", ErrInvalidRemote, u.Hostname())
	}
	if port := u.Port(); port != "" && strings.Trim(port, "0123456789") != "" {
		return fmt.Errorf("%w: port %q is not a number", ErrInvalidRemote, port)
	}
	if len(u.Path) < 2 || u.Opaque != "" || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("%w: needs a repository path and no query or fragment", ErrInvalidRemote)
	}
	if scheme == "https" && u.User != nil {
		return fmt.Errorf("%w: https URLs must not carry credentials", ErrInvalidRemote)
	}
	if u.User != nil {
		if _, hasPassword := u.User.Password(); hasPassword || !userPattern.MatchString(u.User.Username()) {
			return fmt.Errorf("%w: ssh URLs may carry a user name only", ErrInvalidRemote)
		}
	}
	return nil
}
