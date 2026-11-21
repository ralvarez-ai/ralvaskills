// Package skill defines core types and operations for discovered skills on disk.
package skill

// Skill represents a skill discovered on disk.
type Skill struct {
	Name         string
	Version      string
	Description  string
	Path         string // absolute path to the skill folder
	Source       Source
	IsPersonal   bool     // true when path contains a "personal/" segment
	AllowedTools []string // permissions the skill declares it needs, in client rule syntax
	Git          *GitRef  // where a git skill came from; nil for any other source
}

// GitRef records where in a git repository a skill was read, enough to fetch
// the same content again.
type GitRef struct {
	URL    string // the remote it was fetched from
	Rev    string // the branch or tag asked for; empty means the remote's HEAD
	Commit string // the commit the content was read at
	Tree   string // the object id of the skill's own directory: its content identity
	Path   string // the skill's directory inside the repository
}
