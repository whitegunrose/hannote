// Command hannote-lookup prints the dictionary entries for one or more words.
//
// Usage:
//
//	hannote-lookup [-db path/to/krdict.db] [-ko] <word>...
//
// Each argument is looked up separately; quote a phrase to look it up whole
// ('눈이 높다'). Punctuation around a word is ignored, so a word copied out of
// a sentence with its comma or full stop still works.
//
// Exit status is 0 when every word was found, 1 when any word had no entry,
// and 2 for usage errors or a missing dictionary.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/whitegunrose/hannote/internal/dict"
	"github.com/whitegunrose/hannote/internal/lookup"
	"github.com/whitegunrose/hannote/internal/paths"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run is main without the process around it, so tests can call it.
func run(args []string, stdout, stderr io.Writer) int {
	fl := flag.NewFlagSet("hannote-lookup", flag.ContinueOnError)
	fl.SetOutput(stderr)
	dbPath := fl.String("db", "", "dictionary database (default: data dir/krdict.db)")
	showKo := fl.Bool("ko", false, "also show the Korean definitions")
	fl.Usage = func() {
		fmt.Fprintln(stderr, "usage: hannote-lookup [-db krdict.db] [-ko] <word>...")
		fl.PrintDefaults()
	}
	if err := fl.Parse(args); err != nil {
		return 2
	}
	if fl.NArg() == 0 {
		fl.Usage()
		return 2
	}

	store, err := openDict(*dbPath)
	if err != nil {
		fmt.Fprintln(stderr, "hannote-lookup:", err)
		return 2
	}
	defer store.Close()

	// No Analyzer yet: words are looked up as written. Kiwi goes here next.
	l := &lookup.Looker{Dict: store}

	status := 0
	for i, word := range fl.Args() {
		if i > 0 {
			fmt.Fprintln(stdout)
		}
		w, matches, err := l.Lookup(word)
		if err != nil {
			fmt.Fprintln(stderr, "hannote-lookup:", err)
			return 2
		}
		if len(matches) == 0 {
			if w == "" {
				w = word
			}
			hint := ""
			if l.Analyzer == nil {
				hint = " (conjugated forms aren't analyzed yet; try the dictionary form)"
			}
			fmt.Fprintf(stdout, "%s: no entry%s\n", w, hint)
			status = 1
			continue
		}
		for j, m := range matches {
			if j > 0 {
				fmt.Fprintln(stdout)
			}
			printMatch(stdout, w, m, *showKo)
		}
	}
	return status
}

// openDict opens krdict.db at path, or at the default location when path is
// empty, with a hint to run the importer when it isn't there.
func openDict(path string) (*dict.Store, error) {
	if path == "" {
		var err error
		if path, err = paths.DictDB(); err != nil {
			return nil, err
		}
	}
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("no dictionary at %s (run hannote-import first, or pass -db)", path)
	}
	return dict.Open(path)
}

// printMatch writes one entry:
//
//	통해 → 통하다 (通하다) [통하다] 동사 · 중급
//	  1. go through; pass
//	     To go through a place.
func printMatch(w io.Writer, query string, m lookup.Match, showKo bool) {
	e := m.Entry
	head := e.Headword + homonymMark(e.Homonym)
	if e.Headword != query {
		head = query + " → " + head
	}
	if e.Hanja != "" {
		head += " (" + e.Hanja + ")"
	}
	if e.Pronunciation != "" {
		head += " [" + e.Pronunciation + "]"
	}
	var tags []string
	if e.POS != "" {
		tags = append(tags, e.POS)
	} else if e.Unit != "" && e.Unit != "단어" {
		tags = append(tags, e.Unit) // idioms and proverbs have no POS
	}
	if e.Level != "" && e.Level != "없음" {
		tags = append(tags, e.Level)
	}
	if len(tags) > 0 {
		head += " " + strings.Join(tags, " · ")
	}
	fmt.Fprintln(w, head)

	for _, s := range e.Senses {
		first := s.GlossEn
		if first == "" {
			first = s.DefEn
		}
		if first == "" {
			first = s.DefKo // no English for this sense: Korean beats nothing
		}
		fmt.Fprintf(w, "  %d. %s\n", s.Order, first)
		if s.DefEn != "" && s.DefEn != first {
			fmt.Fprintf(w, "     %s\n", s.DefEn)
		}
		if showKo && s.DefKo != "" && s.DefKo != first {
			fmt.Fprintf(w, "     %s\n", s.DefKo)
		}
	}
}

// homonymMark turns homonym number 2 into "²". 0 means the headword has no
// homographs and gets no mark.
func homonymMark(n int) string {
	if n <= 0 {
		return ""
	}
	const digits = "⁰¹²³⁴⁵⁶⁷⁸⁹"
	sup := []rune(digits)
	var b strings.Builder
	for _, c := range fmt.Sprint(n) {
		b.WriteRune(sup[c-'0'])
	}
	return b.String()
}
