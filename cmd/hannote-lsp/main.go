// Command hannote-lsp is Hannote's language server. Editors start it and
// talk to it over stdin/stdout.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "hannote-lsp: not implemented yet")
	os.Exit(1)
}
