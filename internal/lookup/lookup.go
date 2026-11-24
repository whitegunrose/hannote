// Package lookup turns a word as written in a note into dictionary entries.
//
// It is the piece hannote-lookup and the language server share: both hand it
// the word under the cursor (or on the command line) and print whatever it
// returns. The steps are:
//
//  1. Clean the word: drop the punctuation around it.
//  2. Look the word up as-is (headword, then KRDICT's listed forms).
//  3. If an Analyzer is set, look up every dictionary form it suggests.
//
// Step 3 is how conjugated forms like 통해야 will reach 통하다 once Kiwi is
// plugged in. Until then Analyzer is nil and only steps 1 and 2 run.
package lookup

import (
	"strings"
	"unicode"

	"github.com/whitegunrose/hannote/internal/analyzer"
	"github.com/whitegunrose/hannote/internal/dict"
)

// Dict is the part of *dict.Store that lookup needs. Taking an interface
// keeps tests free of SQLite and lets a second dictionary slot in later.
type Dict interface {
	Lookup(word string) ([]dict.Entry, error)
}

// Looker finds every entry a word could be.
type Looker struct {
	Dict     Dict
	Analyzer analyzer.Analyzer // nil: no analysis, the word is looked up as written
}

// Match is one entry found for a word, with the reason it was found.
type Match struct {
	Entry dict.Entry
	// Candidate is the analyzer reading that led here, or nil when the
	// word itself matched a headword or a listed form.
	Candidate *analyzer.Candidate
}

// Lookup cleans word and returns every matching entry, the direct matches
// first. Each entry appears once, even when several readings reach it. The
// cleaned word is returned so callers can show what was actually searched.
func (l *Looker) Lookup(word string) (string, []Match, error) {
	w := Clean(word)
	if w == "" {
		return w, nil, nil
	}

	seen := map[int64]bool{}
	var out []Match
	add := func(entries []dict.Entry, c *analyzer.Candidate) {
		for _, e := range entries {
			if !seen[e.ID] {
				seen[e.ID] = true
				out = append(out, Match{Entry: e, Candidate: c})
			}
		}
	}

	direct, err := l.Dict.Lookup(w)
	if err != nil {
		return w, nil, err
	}
	add(direct, nil)

	if l.Analyzer == nil {
		return w, out, nil
	}
	cands, err := l.Analyzer.Analyze(w)
	if err != nil {
		return w, nil, err
	}
	for i := range cands {
		c := &cands[i]
		if c.Lemma == "" || c.Lemma == w {
			continue // already looked up above
		}
		entries, err := l.Dict.Lookup(c.Lemma)
		if err != nil {
			return w, nil, err
		}
		add(entries, c)
	}
	return w, out, nil
}

// Clean trims whitespace and any punctuation or symbols around a word, so
// "[주로]", "통해야." and “말” look up as 주로, 통해야 and 말. Characters
// inside the word are kept: "눈이 높다" stays a phrase.
func Clean(word string) string {
	return strings.TrimFunc(word, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}
