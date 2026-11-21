package source

import "time"

// Kind names the transport a Spec is read through.
type Kind string

const (
	// KindRegistry is a hosted registry served as an index plus tarballs.
	KindRegistry Kind = "registry"

	// KindGit is a git repository read through the hardened git runner. The
	// factory does not build it yet; declarations can already name it.
	KindGit Kind = "git"

	// KindPath is a directory of skills on the local filesystem, such as a
	// clone of ralvaskills.
	KindPath Kind = "path"
)

// Spec describes one source: what it is called and where it is read from.
// Only the factory turns a Spec into a Resolver.
type Spec struct {
	// ID is the canonical identity, "@owner/repo".
	ID string

	// Type is the transport.
	Type Kind

	// URL is the registry base URL, or the filesystem path of a path source.
	URL string

	// Rev is a git ref. Not honored yet: the factory rejects it.
	Rev string

	// Roots are the directories to walk inside the source. Not honored yet:
	// the factory rejects it.
	Roots []string

	// Excludes are glob patterns pruned before walking. Not honored yet: the
	// factory rejects it.
	Excludes []string

	// Private keeps a source out of surfaces that can leave the machine
	// without a deliberate act. Not enforced yet: the factory rejects it, so a
	// source can never be believed private when nothing hides it.
	Private bool
}

// OfficialSkillsURL is the GitHub URL of the anthropics/skills repository, the
// built-in official source. It is the one git source that needs no trust grant.
const OfficialSkillsURL = "https://github.com/anthropics/skills"

// Options carries what a source needs besides its Spec: where it caches, and
// how a git source is gated.
type Options struct {
	// RegistryCache is where downloaded registry tarballs are extracted.
	RegistryCache string

	// CacheRoot is the folder git sources keep their bare clones, indexes and
	// extracted skills under. A git source cannot be built without it.
	CacheRoot string

	// Gate is called with a git source's identity and URL before anything
	// reads or contacts it, and returns an error to refuse. A git source
	// cannot be built without one, so no path can reach a remote ungated.
	Gate func(id, url string) error

	// RepoFor returns the bare clone a git source reads, for tests that supply an
	// in-memory one. Nil means the real clone under CacheRoot.
	RepoFor func(Spec) RepoOps

	// Offline makes a git source answer from its existing index however old,
	// and never contact its remote.
	Offline bool

	// IndexTTL is how long a git source's name index counts as fresh. Zero
	// means DefaultIndexTTL.
	IndexTTL time.Duration

	// Now returns the current time. Zero means time.Now; tests inject a clock.
	Now func() time.Time
}
