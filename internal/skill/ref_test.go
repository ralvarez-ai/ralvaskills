package skill_test

import (
	"strings"
	"testing"

	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

func TestParseRef(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in          string
		want        skill.Ref
		wantVersion string
		wantErr     bool
	}{
		{in: "go-grpc", want: skill.Ref{Name: "go-grpc"}},
		{in: "go-grpc@1.2.0", want: skill.Ref{Name: "go-grpc"}, wantVersion: "1.2.0"},
		{in: "foo@", want: skill.Ref{Name: "foo"}},
		{
			in:   "@ralvarez-ai/ralvaskills/go-grpc",
			want: skill.Ref{Source: "@ralvarez-ai/ralvaskills", Name: "go-grpc"},
		},
		{
			in:          "@anthropics/skills/xlsx@1.0.0",
			want:        skill.Ref{Source: "@anthropics/skills", Name: "xlsx"},
			wantVersion: "1.0.0",
		},
		{
			in:   "@ralvarezdev/ralvaskills/go-grpc",
			want: skill.Ref{Source: "@ralvarez-ai/ralvaskills", Name: "go-grpc"},
		},
		{in: "@Anthropics/Skills/xlsx", want: skill.Ref{Source: "@anthropics/skills", Name: "xlsx"}},
		{
			in:          "@acme/monorepo//packages/api/skills/api@1.2.0",
			want:        skill.Ref{Source: "@acme/monorepo", Name: "api", Path: "packages/api/skills/api"},
			wantVersion: "1.2.0",
		},
		{in: "", wantErr: true},
		{in: "@", wantErr: true},
		{in: "@1.0.0", wantErr: true},
		{in: "@acme", wantErr: true},
		{in: "@acme/skills", wantErr: true},
		{in: "@acme/skills/", wantErr: true},
		{in: "@acme/skills/Bad_Name", wantErr: true},
		{in: "@/skills/x", wantErr: true},
		{in: "@acme//x", wantErr: true},
		{in: "acme/skills/x", wantErr: true},
		{in: "@acme/skills/a/b", wantErr: true},
		{in: "@acme/skills//../x", wantErr: true},
		{in: "@acme/skills///x", wantErr: true},
		{in: "@acme/skills//x/", wantErr: true},
		{in: "@acme/skills/x@1@2", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			got, version, err := skill.ParseRef(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ParseRef(%q) = %+v, want an error", tc.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseRef(%q): %v", tc.in, err)
			}
			if got != tc.want || version != tc.wantVersion {
				t.Errorf("ParseRef(%q) = %+v, %q; want %+v, %q", tc.in, got, version, tc.want, tc.wantVersion)
			}
		})
	}
}

func TestRefKeyIsQualifiedOnlyOffTheDefaultRegistry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		ref  skill.Ref
		want string
	}{
		{skill.Ref{Name: "go-grpc"}, "go-grpc"},
		{skill.Ref{Source: skill.DefaultRegistryID, Name: "go-grpc"}, "go-grpc"},
		{skill.Ref{Source: skill.OfficialID, Name: "xlsx"}, "@anthropics/skills/xlsx"},
		{skill.Ref{Source: "@acme/monorepo", Name: "api", Path: "packages/api"}, "@acme/monorepo/api"},
	}
	for _, tc := range tests {
		if got := tc.ref.Key(); got != tc.want {
			t.Errorf("Key(%+v) = %q, want %q", tc.ref, got, tc.want)
		}
	}
}

func TestParseKeyRoundTrips(t *testing.T) {
	t.Parallel()
	for _, key := range []string{"go-grpc", "@anthropics/skills/xlsx", "@acme/skills-private/go-architect"} {
		ref, err := skill.ParseKey(key)
		if err != nil {
			t.Fatalf("ParseKey(%q): %v", key, err)
		}
		if got := ref.Key(); got != key {
			t.Errorf("ParseKey(%q).Key() = %q", key, got)
		}
	}
	if _, err := skill.ParseKey("go-grpc@1.0.0"); err == nil {
		t.Error("a key never carries a version; ParseKey accepted one")
	}
}

func TestSourceIdentity(t *testing.T) {
	t.Parallel()

	if got := skill.SourceRegistry.Identity(); got != skill.DefaultRegistryID {
		t.Errorf("registry identity = %q", got)
	}
	if got := skill.SourceOfficial.Identity(); got != skill.OfficialID {
		t.Errorf("official identity = %q", got)
	}
	if got := skill.SourceLocal.Identity(); got != "" {
		t.Errorf("local identity = %q, want none: a clone's identity is machine-specific", got)
	}

	for id, want := range map[string]skill.Source{
		skill.DefaultRegistryID:    skill.SourceRegistry,
		"@ralvarezdev/ralvaskills": skill.SourceRegistry,
		"@Anthropics/Skills":       skill.SourceOfficial,
		"@local/ralvaskills":       skill.SourceLocal,
		"registry":                 skill.SourceRegistry,
		"official":                 skill.SourceOfficial,
		"local":                    skill.SourceLocal,
	} {
		got, err := skill.SourceFromIdentity(id)
		if err != nil || got != want {
			t.Errorf("SourceFromIdentity(%q) = %q, %v; want %q", id, got, err, want)
		}
	}
	declared, declErr := skill.SourceFromIdentity("@Acme/Skills-Private")
	if declErr != nil || declared != skill.DeclaredSource("@acme/skills-private") {
		t.Errorf("a declared identity = %q, %v; want it kept, case-folded", declared, declErr)
	}
	for _, id := range []string{"", "bogus", "@acme", "@acme/x/y"} {
		if _, err := skill.SourceFromIdentity(id); err == nil {
			t.Errorf("SourceFromIdentity(%q) accepted a source this rsk cannot read", id)
		}
	}
}

func TestParseSourceID(t *testing.T) {
	t.Parallel()
	for in, want := range map[string]string{
		"@acme/skills":             "@acme/skills",
		"@Acme/Skills-Private":     "@acme/skills-private",
		"@ralvarezdev/ralvaskills": skill.DefaultRegistryID,
		"@local/my-dir":            "@local/my-dir",
		"@acme/.github":            "@acme/.github",
		"@acme/a.b":                "@acme/a.b",
	} {
		got, err := skill.ParseSourceID(in)
		if err != nil || got != want {
			t.Errorf("ParseSourceID(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{
		"", "acme/skills", "@acme", "@acme/", "@/skills", "@acme/skills/x", "@acme/sk ills", "@-a/b",
		"@acme/.", "@acme/..", "@acme/...",
	} {
		if got, err := skill.ParseSourceID(in); err == nil {
			t.Errorf("ParseSourceID(%q) accepted as %q", in, got)
		}
	}
}

func TestParseFrontmatter(t *testing.T) {
	t.Parallel()
	body := "---\nname: alpha\nversion: 1.2.0\ndescription: \"Does a thing\"\nallowed-tools:\n  - Read\n  - Bash(make lint)\n---\n# Body\nversion: 9\n"
	fm := skill.ParseFrontmatter(strings.NewReader(body))
	if fm.Name != "alpha" || fm.Version != "1.2.0" || fm.Description != "Does a thing" ||
		len(fm.AllowedTools) != 2 || fm.AllowedTools[1] != "Bash(make lint)" {
		t.Errorf("frontmatter = %+v", fm)
	}

	versionless := skill.ParseFrontmatter(strings.NewReader("---\nname: xlsx\ndescription: sheets\n---\n"))
	if versionless.Name != "xlsx" || versionless.Version != "" {
		t.Errorf("a versionless skill must parse with an empty version: %+v", versionless)
	}
	if empty := skill.ParseFrontmatter(strings.NewReader("no frontmatter here")); empty.Name != "" {
		t.Errorf("text without frontmatter parsed as %+v", empty)
	}
}

func TestDeclaredSourceCarriesItsIdentity(t *testing.T) {
	t.Parallel()
	src := skill.DeclaredSource("@acme/skills")
	if !src.IsDeclared() || src.Kind() != skill.SourceGit || src.Identity() != "@acme/skills" ||
		src.Label() != "git  " {
		t.Errorf("declared source = %q: declared=%v kind=%q identity=%q label=%q",
			src, src.IsDeclared(), src.Kind(), src.Identity(), src.Label())
	}
	for _, builtin := range []skill.Source{skill.SourceRegistry, skill.SourceOfficial, skill.SourceLocal} {
		if builtin.IsDeclared() || builtin.Kind() != builtin {
			t.Errorf("%q is built in but reports declared=%v kind=%q", builtin, builtin.IsDeclared(), builtin.Kind())
		}
	}
	if _, err := skill.ParseSource("@acme/skills"); err == nil {
		t.Error(
			"ParseSource accepted a declared identity: only the lock reader may map one, through SourceFromIdentity",
		)
	}
}
