// Package analyzer turns a word as written (일어났어요) into candidate
// dictionary forms (일어나다). Implementations live in subpackages:
// Kiwi first, a hand-written analyzer later.
package analyzer

// Candidate is one possible dictionary form of a word.
type Candidate struct {
	Lemma string // dictionary form, e.g. 일어나다
	POS   string // part of speech, mapped to KRDICT's names (명사, 동사, ...)
	Parts []Part // the breakdown (일어나 + 았 + 어요); may be empty
}

// Part is one piece of a word's breakdown.
type Part struct {
	Form string // surface text, e.g. 았
	Tag  string // analyzer tag, e.g. EP (pre-final ending)
}

// Analyzer returns every plausible dictionary form for a single
// space-separated word (어절). It must not filter by context: the editor
// shows all candidates.
type Analyzer interface {
	Analyze(word string) ([]Candidate, error)
}
