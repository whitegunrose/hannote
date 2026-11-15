package lookup

import (
	"errors"
	"testing"

	"github.com/whitegunrose/hannote/internal/analyzer"
	"github.com/whitegunrose/hannote/internal/dict"
)

// fakeDict answers from a map of headword or form -> entries.
type fakeDict map[string][]dict.Entry

func (f fakeDict) Lookup(word string) ([]dict.Entry, error) { return f[word], nil }

// fakeAnalyzer returns fixed candidates per word.
type fakeAnalyzer map[string][]analyzer.Candidate

func (f fakeAnalyzer) Analyze(word string) ([]analyzer.Candidate, error) { return f[word], nil }

var (
	tonghada = dict.Entry{ID: 1, Headword: "통하다"}
	tong     = dict.Entry{ID: 2, Headword: "통", Homonym: 1}
	juro     = dict.Entry{ID: 3, Headword: "주로"}
)

func TestClean(t *testing.T) {
	cases := map[string]string{
		"[주로]":     "주로",
		"통해야.":     "통해야",
		"“말”":      "말",
		"  눈이 높다 ": "눈이 높다",
		"...":      "",
		"3월에":      "3월에",
	}
	for in, want := range cases {
		if got := Clean(in); got != want {
			t.Errorf("Clean(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLookupWithoutAnalyzer(t *testing.T) {
	l := &Looker{Dict: fakeDict{"주로": {juro}}}

	w, got, err := l.Lookup("[주로],")
	if err != nil {
		t.Fatal(err)
	}
	if w != "주로" || len(got) != 1 || got[0].Entry.ID != juro.ID || got[0].Candidate != nil {
		t.Fatalf("Lookup([주로],) = %q %+v", w, got)
	}

	// Without an analyzer a conjugated form only matches if KRDICT lists it.
	if _, got, _ := l.Lookup("통해야"); len(got) != 0 {
		t.Errorf("통해야 without analyzer = %+v, want nothing", got)
	}
}

func TestLookupWithAnalyzer(t *testing.T) {
	l := &Looker{
		Dict: fakeDict{"통하다": {tonghada}, "통": {tong}},
		Analyzer: fakeAnalyzer{"통해야": {
			{Lemma: "통하다", POS: "동사"},
			{Lemma: "통하다", POS: "동사"}, // a second reading of the same lemma
			{Lemma: "통", POS: "명사"},
			{Lemma: "통해야"}, // the word itself: already looked up
		}},
	}
	_, got, err := l.Lookup("통해야")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d matches, want 2 (each entry once): %+v", len(got), got)
	}
	if got[0].Entry.Headword != "통하다" || got[0].Candidate == nil || got[0].Candidate.POS != "동사" {
		t.Errorf("first match = %+v", got[0])
	}
	if got[1].Entry.Headword != "통" {
		t.Errorf("second match = %+v", got[1])
	}
}

func TestDirectMatchesComeFirst(t *testing.T) {
	l := &Looker{
		Dict:     fakeDict{"통": {tong}, "통하다": {tonghada}},
		Analyzer: fakeAnalyzer{"통": {{Lemma: "통하다"}}},
	}
	_, got, _ := l.Lookup("통")
	if len(got) != 2 || got[0].Candidate != nil || got[1].Candidate == nil {
		t.Errorf("got %+v, want the direct match before the analyzer's", got)
	}
}

type errDict struct{}

func (errDict) Lookup(string) ([]dict.Entry, error) { return nil, errors.New("disk on fire") }

func TestLookupError(t *testing.T) {
	l := &Looker{Dict: errDict{}}
	if _, _, err := l.Lookup("말"); err == nil {
		t.Error("error from the dictionary was swallowed")
	}
}

func TestEmptyWord(t *testing.T) {
	l := &Looker{Dict: errDict{}} // must not even be called
	if w, got, err := l.Lookup(" ?! "); w != "" || got != nil || err != nil {
		t.Errorf("Lookup of punctuation = %q %v %v", w, got, err)
	}
}
