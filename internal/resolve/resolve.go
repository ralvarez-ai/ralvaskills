// Package resolve explains what a skill reference resolves to, without
// installing anything. It returns plain structs, so the `rsk resolve` command
// and the MCP resolve_reference tool render the same answer.
package resolve

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

// Status says how a reference resolved.
type Status string

const (
	// Resolved means exactly one skill answers the reference.
	Resolved Status = "resolved"

	// Ambiguous means a bare name is provided by more than one source.
	Ambiguous Status = "ambiguous"

	// NotFound means no source that was asked has the skill.
	NotFound Status = "not-found"

	// Unavailable means a source could not be consulted, so a bare name is not
	// resolved from the others: failing closed keeps an unreachable private
	// source from letting a public one answer in its place.
	Unavailable Status = "unavailable"

	// UnknownSource means a qualified reference names a source this machine does
	// not have.
	UnknownSource Status = "unknown-source"

	// Bundle means the name is a bundle, which expands to its members.
	Bundle Status = "bundle"

	// Unsupported means the reference asks for something no source provides yet.
	Unsupported Status = "unsupported"
)

type (
	// Options carries what a lookup needs besides the chain.
	Options struct {
		// Bundles are the known bundles: a bare name matching one is a bundle.
		Bundles []config.Bundle

		// ImplicitHint says what `rsk install` would do for a qualified reference
		// to a source this machine does not have. Nil gives no hint.
		ImplicitHint func(id string) string
	}

	// Result is the answer for one reference.
	Result struct {
		Reference string   `json:"reference"`
		Name      string   `json:"name"`
		Source    string   `json:"source,omitempty"`
		Requested string   `json:"requested_version,omitempty"`
		Status    Status   `json:"status"`
		Message   string   `json:"message,omitempty"`
		Skill     *Skill   `json:"skill,omitempty"`
		Consulted []Answer `json:"consulted,omitempty"`
		Members   []string `json:"members,omitempty"`
	}

	// Skill is the skill a reference resolved to.
	Skill struct {
		Source  string `json:"source"`
		Type    string `json:"type"`
		URL     string `json:"url,omitempty"`
		Name    string `json:"name"`
		Version string `json:"version,omitempty"`
		// Key is how the skill is written in rsk.mod and rsk.lock: bare for the
		// registry and the local clone, qualified for any other source.
		Key    string `json:"key"`
		Path   string `json:"path,omitempty"`
		Rev    string `json:"rev,omitempty"`
		Commit string `json:"commit,omitempty"`
		Tree   string `json:"tree,omitempty"`
	}

	// Answer is what one source said when a bare name was looked up in all of
	// them.
	Answer struct {
		Source  string `json:"source"`
		Type    string `json:"type"`
		Has     bool   `json:"has"`
		Version string `json:"version,omitempty"`
		Error   string `json:"error,omitempty"`
	}
)

// Reference resolves raw the way `rsk install` would, and says why when it does
// not. A malformed reference is an error; every other outcome is a Result.
func Reference(ctx context.Context, raw string, chain source.Chain, opts Options) (Result, error) {
	ref, version, err := skill.ParseRef(raw)
	if err != nil {
		return Result{}, err
	}
	res := Result{Reference: raw, Name: ref.Name, Source: ref.Source, Requested: version}

	if ref.Source == "" {
		if bundle, ok := config.FindBundle(opts.Bundles, ref.Name); ok {
			res.Status = Bundle
			for _, m := range bundle.Skills {
				res.Members = append(res.Members, m.Name)
			}
			return res, nil
		}
		return bare(ctx, res, ref, chain, version), nil
	}
	return qualified(ctx, res, ref, chain, version, opts), nil
}

func qualified(
	ctx context.Context, res Result, ref skill.Ref, chain source.Chain, version string, opts Options,
) Result {
	found, err := chain.Find(ctx, ref)
	if err == nil && version != "" && found.Source.IsDeclared() {
		found, err = chain.FindVersion(ctx, string(found.Source), found.Name, version)
	}
	if err != nil {
		return fail(res, ref, err, opts)
	}
	res.Status = Resolved
	res.Skill = describe(ref.Source, chain, &found)
	return res
}

func fail(res Result, ref skill.Ref, err error, opts Options) Result {
	switch {
	case errors.Is(err, source.ErrUnknownSource):
		res.Status = UnknownSource
		res.Message = fmt.Sprintf("source %s is not declared on this machine", ref.Source)
		if opts.ImplicitHint != nil {
			if hint := opts.ImplicitHint(ref.Source); hint != "" {
				res.Message += ": " + hint
			}
		}
	case errors.Is(err, source.ErrNotFound):
		res.Status, res.Message = NotFound, err.Error()
	case errors.Is(err, source.ErrUnsupported):
		res.Status, res.Message = Unsupported, err.Error()
	default:
		res.Status, res.Message = Unavailable, err.Error()
	}
	return res
}

func bare(ctx context.Context, res Result, ref skill.Ref, chain source.Chain, version string) Result {
	answers := chain.Survey(ctx, ref.Name)
	var hits, down []int
	res.Consulted = make([]Answer, len(answers))
	for i := range answers {
		a := &answers[i]
		out := Answer{Source: a.Spec.ID, Type: string(a.Spec.Type)}
		switch {
		case a.Err == nil:
			out.Has, out.Version = true, a.Skill.Version
			hits = append(hits, i)
		case errors.Is(a.Err, source.ErrNotFound):
		default:
			out.Error = a.Err.Error()
			down = append(down, i)
		}
		res.Consulted[i] = out
	}

	switch {
	case len(down) > 0:
		res.Status = Unavailable
		res.Message = fmt.Sprintf(
			"cannot resolve %q: source %s could not be consulted (%s); a bare name is never resolved from the others when one is down",
			ref.Name,
			answers[down[0]].Spec.ID,
			res.Consulted[down[0]].Error,
		)
	case len(hits) == 0:
		res.Status = NotFound
		res.Message = fmt.Sprintf("%q is in no source", ref.Name)
	case len(hits) > 1:
		res.Status = Ambiguous
		res.Message = ambiguity(ref.Name, answers, hits)
	default:
		a := &answers[hits[0]]
		found := a.Skill
		if version != "" && found.Source.IsDeclared() {
			var err error
			if found, err = chain.FindVersion(ctx, string(found.Source), found.Name, version); err != nil {
				return fail(res, skill.Ref{Source: a.Spec.ID, Name: ref.Name}, err, Options{})
			}
		}
		res.Status = Resolved
		res.Skill = describe(a.Spec.ID, chain, &found)
	}
	return res
}

func ambiguity(name string, answers []source.Answer, hits []int) string {
	ids := make([]string, len(hits))
	forms := make([]string, len(hits))
	for i, h := range hits {
		id := answers[h].Spec.ID
		ids[i] = id
		forms[i] = id + "/" + name
	}
	return fmt.Sprintf(
		"%q is provided by %s: choose one with %s",
		name,
		strings.Join(ids, " and "),
		strings.Join(forms, " or "),
	)
}

func describe(id string, chain source.Chain, s *skill.Skill) *Skill {
	out := &Skill{
		Source:  id,
		Name:    s.Name,
		Version: s.Version,
		Path:    s.Path,
		Key:     skill.Ref{Source: id, Name: s.Name}.Key(),
	}
	if r := chain.ByID(id); len(r) > 0 {
		spec := r[0].Spec()
		out.Type, out.URL = string(spec.Type), spec.URL
	}
	if s.Git != nil {
		out.URL, out.Rev, out.Commit, out.Tree, out.Path = s.Git.URL, s.Git.Rev, s.Git.Commit, s.Git.Tree, s.Git.Path
	}
	return out
}
