package catalog

import (
	"context"
	"fmt"
	"sort"

	"go.ralva.dev/ralvaskills/v3/internal/config"
	"go.ralva.dev/ralvaskills/v3/internal/skill"
	"go.ralva.dev/ralvaskills/v3/internal/source"
)

// ListOptions selects the skills List returns.
type ListOptions struct {
	// Source limits the listing to one source kind; "" lists every source.
	Source skill.Source
	// Private includes the declared sources marked private.
	Private bool
	// Bundle keeps only the skills named by this bundle; "" keeps all.
	Bundle string
	// Personal keeps the personal skills; false drops them.
	Personal bool
}

// ListResult is what List found: the skills, ordered by source then name, and
// the warnings about sources that could not be walked. A warning never stops
// the listing, so the rest of the sources are still returned.
type ListResult struct {
	Skills   []skill.Skill
	Warnings []string
}

// List walks the local, official and declared sources of chain per opts. An
// unreachable or untrusted git source becomes a warning rather than an error.
// A bundle that does not exist is an error, and the warnings gathered before it
// are still returned.
func List(ctx context.Context, chain source.Chain, opts ListOptions) (ListResult, error) {
	res := ListResult{}
	res.Skills, res.Warnings = loadBySource(ctx, chain, opts.Source, opts.Private)

	if opts.Bundle != "" {
		kept, warnings, err := filterByBundle(res.Skills, opts.Bundle)
		res.Warnings = append(res.Warnings, warnings...)
		if err != nil {
			return res, err
		}
		res.Skills = kept
	}

	if !opts.Personal {
		res.Skills = keep(res.Skills, func(s skill.Skill) bool { return !s.IsPersonal })
	}

	sort.Slice(res.Skills, func(i, j int) bool {
		if res.Skills[i].Source != res.Skills[j].Source {
			return res.Skills[i].Source < res.Skills[j].Source
		}
		return res.Skills[i].Name < res.Skills[j].Name
	})
	return res, nil
}

// loadBySource walks the local, official and declared sources per sourceFilter
// ("" means all), collecting a warning for a source that can't be walked.
func loadBySource(
	ctx context.Context, chain source.Chain, sourceFilter skill.Source, includePrivate bool,
) (all []skill.Skill, warnings []string) {
	walk := func(want skill.Source, label func(source.Resolver, error) string, srcs source.Chain) {
		if sourceFilter != "" && sourceFilter != want {
			return
		}
		for _, src := range srcs {
			listed, err := src.All(ctx)
			if err != nil {
				warnings = append(warnings, label(src, err))
				continue
			}
			all = append(all, listed...)
		}
	}

	walk(skill.SourceLocal, func(_ source.Resolver, err error) string {
		return fmt.Sprintf("walk local skills: %v", err)
	}, chain.Primary())
	walk(skill.SourceOfficial, func(_ source.Resolver, err error) string { return err.Error() },
		chain.ByID(skill.OfficialID))
	walk(skill.SourceGit, func(src source.Resolver, err error) string {
		return fmt.Sprintf("source %s: %v", src.Spec().ID, err)
	}, declaredSources(chain, includePrivate))
	return all, warnings
}

// declaredSources is the chain's declared sources: git ones other than the
// built-in official one, which the catalog lists as its own source, and path
// ones.
func declaredSources(chain source.Chain, includePrivate bool) source.Chain {
	var out source.Chain
	for _, r := range chain {
		spec := r.Spec()
		if spec.Private && !includePrivate {
			continue
		}
		if (spec.Type == source.KindGit && spec.ID != skill.OfficialID) ||
			(spec.Type == source.KindPath && spec.Declared) {
			out = append(out, r)
		}
	}
	return out
}

// filterByBundle keeps only the skills named in bundleName's catalog entry.
func filterByBundle(all []skill.Skill, bundleName string) ([]skill.Skill, []string, error) {
	var warnings []string
	bundles, catalogWarn := config.LoadCatalog("")
	if catalogWarn != nil {
		warnings = append(warnings, fmt.Sprintf("user catalog: %v", catalogWarn))
	}
	bundle, ok := config.FindBundle(bundles, bundleName)
	if !ok {
		return nil, warnings, fmt.Errorf(
			"bundle %q not found — run 'rsk catalog bundles' to see all bundles", bundleName)
	}
	want := make(map[string]bool, len(bundle.Skills))
	for _, ref := range bundle.Skills {
		want[ref.Name] = true
	}
	return keep(all, func(s skill.Skill) bool { return want[s.Name] }), warnings, nil
}

// keep returns the skills for which pred holds.
func keep(in []skill.Skill, pred func(skill.Skill) bool) []skill.Skill {
	out := make([]skill.Skill, 0, len(in))
	for _, s := range in {
		if pred(s) {
			out = append(out, s)
		}
	}
	return out
}

// BundleRow is one bundle with how many of its skills are installed.
type BundleRow struct {
	config.Bundle
	Linked int
	Total  int
}

// BundleRows lists the catalog's bundles ordered by name, counting a skill as
// linked when isLinked reports it. The user catalog failing to load is a
// warning: the built-in bundles are still listed.
func BundleRows(isLinked func(name string) bool) (rows []BundleRow, warnings []string) {
	bundles, catalogWarn := config.LoadCatalog("")
	if catalogWarn != nil {
		warnings = append(warnings, fmt.Sprintf("user catalog: %v", catalogWarn))
	}
	rows = make([]BundleRow, 0, len(bundles))
	for _, b := range bundles {
		linked := 0
		for _, ref := range b.Skills {
			if isLinked(ref.Name) {
				linked++
			}
		}
		rows = append(rows, BundleRow{Bundle: b, Linked: linked, Total: len(b.Skills)})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })
	return rows, warnings
}
