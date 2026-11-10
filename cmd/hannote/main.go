// Command hannote is the custom terminal editor. It comes after the
// language server works in Neovim.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "hannote: the editor is not implemented yet; use Neovim with nvim/ for now")
	os.Exit(1)
}
