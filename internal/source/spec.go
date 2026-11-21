package source

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

	// KindOfficial is the cached anthropics/skills clone.
	KindOfficial Kind = "official"
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

// Options carries what a source needs besides its Spec: where it caches.
type Options struct {
	// RegistryCache is where downloaded registry tarballs are extracted.
	RegistryCache string

	// OfficialCache is where the anthropics/skills clone lives.
	OfficialCache string
}
