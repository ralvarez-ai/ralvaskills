package trust_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.ralva.dev/ralvaskills/v3/internal/schema"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
)

var now = time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC)

const (
	acme    = "@acme/skills"
	acmeURL = "git@github.com:acme/skills.git"
)

func TestGrantIsKeyedByIdentityAndURL(t *testing.T) {
	t.Parallel()
	f := trust.Grant(trust.File{}, acme, acmeURL, now)

	if !f.Granted(acme, acmeURL) || !f.Granted("@ACME/Skills", acmeURL) {
		t.Error("a granted identity and URL is not trusted")
	}
	if f.Granted(acme, "git@evil.example:acme/skills.git") {
		t.Error("changing the URL kept the grant: a project could redirect a trusted identity")
	}
	if f.Granted("@acme/other", acmeURL) {
		t.Error("a grant covered a different identity")
	}
	if (trust.File{}).Granted(acme, acmeURL) {
		t.Error("nothing is trusted by default")
	}
}

func TestGrantReplacesAndRevokeRemoves(t *testing.T) {
	t.Parallel()
	f := trust.Grant(trust.File{}, acme, acmeURL, now)
	moved := "ssh://git@github.com/acme/skills.git"
	f = trust.Grant(f, acme, moved, now.Add(time.Hour))

	if len(f.Grants) != 1 || f.Granted(acme, acmeURL) || !f.Granted(acme, moved) {
		t.Errorf("a new grant for an identity must replace the old one: %+v", f.Grants)
	}
	f, revoked := trust.Revoke(f, "@Acme/Skills")
	if !revoked || len(f.Grants) != 0 {
		t.Errorf("Revoke = %+v, %v", f, revoked)
	}
	if _, revoked = trust.Revoke(f, acme); revoked {
		t.Error("revoked a grant that was not there")
	}
}

func TestRequireNamesTheCommandToRun(t *testing.T) {
	t.Parallel()
	err := trust.Require(trust.File{}, acme, acmeURL)
	var untrusted *trust.UntrustedError
	if !errors.As(err, &untrusted) || !errors.Is(err, trust.ErrUntrusted) {
		t.Fatalf("err = %v, want an UntrustedError", err)
	}
	for _, want := range []string{acme, acmeURL, "rsk source trust " + acme} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q lacks %q", err, want)
		}
	}
	if err = trust.Require(trust.Grant(trust.File{}, acme, acmeURL, now), acme, acmeURL); err != nil {
		t.Errorf("a trusted source was refused: %v", err)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "rsk", "trust.json")
	if f, err := trust.Load(path); err != nil || len(f.Grants) != 0 {
		t.Fatalf("a missing file should load empty: %+v, %v", f, err)
	}
	f := trust.Grant(trust.File{}, acme, acmeURL, now)
	if err := trust.Save(path, f); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(path); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("trust.json mode = %v, %v; want 0600", info.Mode().Perm(), err)
	}
	got, err := trust.Load(path)
	if err != nil || !got.Granted(acme, acmeURL) || !got.Grants[0].GrantedAt.Equal(now) {
		t.Fatalf("round trip = %+v, %v", got, err)
	}
	if got.SchemaVersion != schema.Trust {
		t.Errorf("version = %d", got.SchemaVersion)
	}
}

func TestLoadRefusesNewerOrCorruptFiles(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	write := func(name, body string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	if _, err := trust.Load(write("n.json", `{"version":9,"grants":[]}`)); !errors.Is(err, schema.ErrUnsupported) {
		t.Errorf("newer file err = %v", err)
	}
	if _, err := trust.Load(write("c.json", `nope`)); err == nil {
		t.Error("corrupt file accepted")
	}
}
