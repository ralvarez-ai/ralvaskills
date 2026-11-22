package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"go.ralva.dev/ralvaskills/v3/internal/digest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
)

// digestPreview is the text a pin would add to the always-on .rsk/SKILLS.md.
type digestPreview struct {
	Name   string
	Origin string
	Text   string
}

// digestPreviews returns a preview for each skill from a source other than the
// default registry or the local clone: remote text that would enter every
// session. A skill without a DIGEST.md contributes its description, as it does
// in the generated file.
func digestPreviews(skills []skill.Skill) []digestPreview {
	var out []digestPreview
	for i := range skills {
		s := &skills[i]
		if !s.Source.UsesGit() {
			continue
		}
		out = append(out, digestPreview{
			Name: s.Name, Origin: originOf(s.Source), Text: digest.Section(digest.Load(s.Path, s.Name)),
		})
	}
	return out
}

func originOf(src skill.Source) string {
	if id := src.Identity(); id != "" {
		return id
	}
	return string(src)
}

// installedDigestPreview previews an already installed skill, for `rsk pin`.
func installedDigestPreview(projectDir, name, origin string) digestPreview {
	return digestPreview{
		Name: name, Origin: origin, Text: digest.Section(digest.LoadInstalled(projectDir, name)),
	}
}

// confirmDigests shows the text each preview would put in the always-on file and
// gets agreement to load it (trust rule 7). A person at a terminal is asked;
// --yes accepts it after showing it; anything else is refused, since nobody read
// it. It reports false when the person declined.
func confirmDigests(cmd *cobra.Command, previews []digestPreview, yes bool) (bool, error) {
	if len(previews) == 0 {
		return true, nil
	}
	out := cmd.OutOrStdout()
	fmt.Fprintln(out)
	for i := range previews {
		p := &previews[i]
		fmt.Fprintf(out, "Pinning %s from %s adds this to .rsk/SKILLS.md, which your agent loads every session:\n\n",
			p.Name, p.Origin)
		fmt.Fprintln(out, p.Text)
		fmt.Fprintln(out)
	}
	if yes {
		return true, nil
	}
	answer, asked := askHuman(cmd, "Load this text into every session?")
	if !asked {
		return false, errors.New(
			"the text above becomes always-on agent instructions and nobody has confirmed it: " +
				"read it, then run again with --yes, or run in a terminal to be asked",
		)
	}
	return answer, nil
}
