// count-tokens scans every SKILL.md under skills/, estimates token counts
// (body + description-only + side files), and writes docs/TOKEN_COUNTS.md.
//
// Usage:
//
//	task tokens
package main

import (
	"flag"
	"fmt"
	"os"

	"go.ralva.dev/ralvaskills/v3/internal/tokensreport"
)

func main() {
	format := flag.String("format", "markdown", "output format: markdown or json")
	sortBy := flag.String("sort", "body", "sort column: body, desc, name, category")
	flag.Parse()

	path, err := tokensreport.Run(*format, *sortBy)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("wrote", path)
}
