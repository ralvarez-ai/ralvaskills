package ui

import (
	"context"

	"go.ralva.dev/termkit"
)

// InSession reports whether ctx belongs to an interactive session run: the
// session always runs a command with a capture in its context, and there is no
// terminal to prompt on then.
func InSession(ctx context.Context) bool {
	return termkit.InCapture(ctx)
}
