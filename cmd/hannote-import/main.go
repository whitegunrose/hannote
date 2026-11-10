// Command hannote-import builds krdict.db from the KRDICT XML download.
//
// Usage:
//
//	hannote-import [-o path/to/krdict.db] <dir-with-xml-files>
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/whitegunrose/hannote/internal/paths"
)

func main() {
	out := flag.String("o", "", "output database (default: data dir/krdict.db)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: hannote-import [-o krdict.db] <xml-dir>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	if *out == "" {
		p, err := paths.DictDB()
		if err != nil {
			fmt.Fprintln(os.Stderr, "hannote-import:", err)
			os.Exit(1)
		}
		*out = p
	}
	fmt.Fprintf(os.Stderr, "hannote-import: not implemented yet (would read %s, write %s)\n", flag.Arg(0), *out)
	os.Exit(1)
}
