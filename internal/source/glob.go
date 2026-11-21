package source

import (
	"regexp"
	"strings"
)

// compileGlob turns an exclude pattern into a regular expression over
// slash-separated paths: "**" crosses directories, "*" and "?" stay inside one,
// "**/" also matches no directory at all, and every other character is literal.
func compileGlob(pattern string) (*regexp.Regexp, error) {
	var b strings.Builder
	b.WriteString("^")
	for i := 0; i < len(pattern); {
		switch {
		case strings.HasPrefix(pattern[i:], "**/"):
			b.WriteString("(?:.*/)?")
			i += 3
		case strings.HasPrefix(pattern[i:], "**"):
			b.WriteString(".*")
			i += 2
		case pattern[i] == '*':
			b.WriteString("[^/]*")
			i++
		case pattern[i] == '?':
			b.WriteString("[^/]")
			i++
		default:
			b.WriteString(regexp.QuoteMeta(string(pattern[i])))
			i++
		}
	}
	b.WriteString("$")
	return regexp.Compile(b.String())
}

// matchGlob reports whether re matches dir, or dir as a directory prefix, so
// "**/node_modules/**" excludes a skill directly under node_modules too.
func matchGlob(re *regexp.Regexp, dir string) bool {
	return re.MatchString(dir) || re.MatchString(dir+"/")
}

// globMatch compiles pattern and matches it against dir.
func globMatch(pattern, dir string) bool {
	re, err := compileGlob(pattern)
	return err == nil && matchGlob(re, dir)
}
