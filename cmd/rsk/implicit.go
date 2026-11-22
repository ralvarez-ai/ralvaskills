package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"go.ralva.dev/lib/term"

	"go.ralva.dev/ralvaskills/v3/internal/git"
	"go.ralva.dev/ralvaskills/v3/internal/manifest"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/sources"
	"go.ralva.dev/ralvaskills/v3/internal/trust"
	"go.ralva.dev/ralvaskills/v3/internal/ui"
)

// defaultImplicitProbe is the real transport probe: git ls-remote through the
// hardened runner.
var defaultImplicitProbe sources.Probe = git.LsRemote

var (
	// implicitProbe picks the transport for an approved implicit source. Tests
	// replace it so none reaches a network.
	implicitProbe = defaultImplicitProbe

	// isInteractive reports whether a person can answer a prompt. Tests replace it.
	isInteractive = term.StdioIsTerminal
)

const (
	flagTrust = "trust"
	flagSave  = "save"
	flagFrom  = "from"
)

// errImplicitDeclined is returned when the person at the terminal says no.
var errImplicitDeclined = errors.New("source not approved")

// askHuman asks a yes/no question that only a person at a terminal may answer.
// Unlike term.Confirm it never proceeds on its own: not with --yes, not when
// there is no terminal, not inside the interactive session. asked is false when
// no question could be put.
func askHuman(cmd *cobra.Command, prompt string) (yes, asked bool) {
	if !isInteractive() || term.InCapture(cmd.Context()) {
		return false, false
	}
	fmt.Fprintf(cmd.ErrOrStderr(), "%s [y/N] ", prompt)
	answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil && answer == "" {
		return false, true
	}
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "y", "yes":
		return true, true
	default:
		return false, true
	}
}

// implicitCandidates returns the sources the references name that this machine
// has not declared and rsk can derive: each is a source rsk would have to clone
// for the first time. Built-in and declared sources resolve through the chain.
func implicitCandidates(refs []string, declared sources.File) []sources.Candidate {
	var out []sources.Candidate
	for _, raw := range refs {
		ref, _, err := skill.ParseRef(raw)
		if err != nil || ref.Source == "" {
			continue
		}
		if _, ok := sources.Find(declared, ref.Source); ok {
			continue
		}
		c, err := sources.Derive(ref.Source)
		if err != nil || slices.ContainsFunc(out, func(have sources.Candidate) bool { return have.ID == c.ID }) {
			continue
		}
		out = append(out, c)
	}
	return out
}

// approveImplicit settles each candidate for this run: approved by --trust, or
// by the person at the terminal after seeing the URLs, and otherwise refused
// with the ways forward. Only then does it contact the host, to pick SSH or
// HTTPS, and every URL goes through the trust gate first. Nothing is written.
func approveImplicit(
	cmd *cobra.Command, cands []sources.Candidate, declared sources.File, trustIDs []string,
) ([]sources.Entry, error) {
	approvedIDs, err := checkTrustFlag(trustIDs, cands, declared)
	if err != nil {
		return nil, err
	}
	if len(cands) == 0 {
		return nil, nil
	}
	if offline, offlineErr := offlineRequested(cmd); offlineErr != nil {
		return nil, offlineErr
	} else if offline {
		return nil, errors.New("--offline cannot contact " + cands[0].ID + ", which this machine has not cached")
	}

	entries := make([]sources.Entry, 0, len(cands))
	for i := range cands {
		c := &cands[i]
		if !slices.Contains(approvedIDs, c.ID) {
			if err = askImplicit(cmd, c); err != nil {
				return nil, err
			}
		}
		gate := func(id, url string) error {
			if id == c.ID && slices.Contains(c.URLs(), url) {
				return nil
			}
			return &trust.UntrustedError{ID: id, URL: url}
		}
		entry, chooseErr := sources.Choose(cmd.Context(), *c, "", gate, implicitProbe)
		if chooseErr != nil {
			return nil, chooseErr
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// checkTrustFlag validates --trust and returns the identities it approves. The
// flag only covers a source this command would otherwise open implicitly: a
// declared one is trusted with `rsk source trust`, which needs a terminal.
func checkTrustFlag(ids []string, cands []sources.Candidate, declared sources.File) ([]string, error) {
	var approved []string
	for _, raw := range ids {
		id, err := skill.ParseSourceID(raw)
		if err != nil {
			return nil, fieldError(flagTrust, "%v", err)
		}
		switch {
		case slices.ContainsFunc(cands, func(c sources.Candidate) bool { return c.ID == id }):
			approved = append(approved, id)
		case declaredHas(declared, id):
			return nil, fieldError(flagTrust, "%s is declared on this machine: --trust does not cover it; "+
				"trust it yourself in a terminal with: rsk source trust %s", id, id)
		default:
			return nil, fieldError(flagTrust, "%s is not a source this command opens", id)
		}
	}
	return approved, nil
}

func declaredHas(declared sources.File, id string) bool {
	_, ok := sources.Find(declared, id)
	return ok
}

// askImplicit shows what rsk would contact and asks the person to approve it for
// this run. With nobody to ask it refuses and says how to proceed.
func askImplicit(cmd *cobra.Command, c *sources.Candidate) error {
	out := cmd.ErrOrStderr()
	if !isInteractive() || term.InCapture(cmd.Context()) {
		return fmt.Errorf(
			"source %s is not registered on this machine, and rsk will not contact it without your approval "+
				"(%s over SSH, then %s over HTTPS). --yes does not cover this: run it in a terminal to be asked, "+
				"pass --trust %s to approve it for this one run, or register it yourself: "+
				"rsk source add %s, then rsk source trust %s",
			c.ID, c.SSH, c.HTTPS, c.ID, c.ID, c.ID,
		)
	}
	fmt.Fprintf(out, "\nSource %s is not registered on this machine.\n", c.ID)
	fmt.Fprintln(out, "rsk would let git contact it with your own SSH key or credential helper:")
	fmt.Fprintf(out, "  %s  (SSH, tried first)\n", c.SSH)
	fmt.Fprintf(out, "  %s  (HTTPS, if SSH cannot be used)\n", c.HTTPS)
	fmt.Fprintln(
		out,
		"Its skills are text your agents will follow. Nothing is saved: this approval is for this run only.",
	)
	yes, _ := askHuman(cmd, "Contact "+c.ID+" for this run?")
	if !yes {
		return errImplicitDeclined
	}
	return nil
}

// saveImplicit registers the sources an install opened for one run, untrusted.
// It never writes trust.json: a grant is made only by `rsk source trust` at a
// terminal, however the run was approved.
func saveImplicit(out io.Writer, entries []sources.Entry) error {
	declared, err := sources.Load(sources.DefaultPath())
	if err != nil {
		return err
	}
	for i := range entries {
		if declared, err = sources.Add(declared, entries[i]); err != nil {
			return err
		}
	}
	if err = sources.Save(sources.DefaultPath(), declared); err != nil {
		return err
	}
	for i := range entries {
		ui.Success(out, "registered "+entries[i].ID+" at "+entries[i].URL)
		ui.Info(out, trustHint(entries[i].ID))
	}
	return nil
}

// lockImplicitCandidates returns the git sources a lock names that this machine
// has not declared, as candidates to approve for the run. A source rsk cannot
// derive is left out: restoring it fails with the command to declare it.
func lockImplicitCandidates(lock manifest.Lock, declared sources.File) []sources.Candidate {
	refs := make([]string, 0, len(lock.Skills))
	for i := range lock.Skills {
		e := &lock.Skills[i]
		if e.Git != nil && e.Source.IsDeclared() {
			refs = append(refs, string(e.Source)+"/"+e.Name)
		}
	}
	return implicitCandidates(refs, declared)
}

// approveLockSources approves, for one `sync`, the undeclared sources a lock
// names. A lock is untrusted input: it never opens a source on its own.
func approveLockSources(cmd *cobra.Command, lock manifest.Lock, trustIDs []string) ([]sources.Entry, error) {
	declared, err := sources.Load(sources.DefaultPath())
	if err != nil {
		return nil, err
	}
	return approveImplicit(cmd, lockImplicitCandidates(lock, declared), declared, trustIDs)
}

// prepareImplicit settles the unregistered sources the references name and, with
// --save, arranges for them to be registered once the install has happened.
func prepareImplicit(cmd *cobra.Command, opts *installOpts, args []string) ([]sources.Entry, error) {
	declared, err := sources.Load(sources.DefaultPath())
	if err != nil {
		return nil, err
	}
	implicit, err := approveImplicit(cmd, implicitCandidates(args, declared), declared, opts.trust)
	if err != nil {
		return nil, err
	}
	if !opts.save {
		return implicit, nil
	}
	if len(implicit) == 0 {
		return nil, fieldError(flagSave, "--save has nothing to register: no reference names an unregistered source")
	}
	out := cmd.OutOrStdout()
	opts.onApplied = func() error { return saveImplicit(out, implicit) }
	return implicit, nil
}
