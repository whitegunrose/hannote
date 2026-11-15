package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/whitegunrose/hannote/internal/dict"
)

// testDB builds a tiny krdict.db with the entries the tests look up.
func testDB(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "krdict.db")
	w, err := dict.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	entries := []dict.Entry{
		{KrdictID: 1, Headword: "통하다", Unit: "단어", POS: "동사", Level: "중급", Hanja: "通하다",
			Pronunciation: "통하다", Forms: []string{"통해", "통하니"},
			Senses: []dict.Sense{
				{Order: 1, DefKo: "어떤 곳을 거쳐 지나가다.", GlossEn: "go through; pass",
					DefEn: "To go through a certain place."},
				{Order: 2, DefKo: "막힘이 없이 흐르다."}, // no English
			}},
		{KrdictID: 2, Headword: "방금", Homonym: 1, Unit: "단어", POS: "명사", Level: "초급",
			Senses: []dict.Sense{{Order: 1, GlossEn: "moment ago"}}},
		{KrdictID: 3, Headword: "방금", Homonym: 2, Unit: "단어", POS: "부사", Level: "초급",
			Senses: []dict.Sense{{Order: 1, GlossEn: "a moment ago"}}},
		{KrdictID: 4, Headword: "눈이 높다", Unit: "관용구", Level: "없음",
			Senses: []dict.Sense{{Order: 1, GlossEn: "have high standards"}}},
	}
	for _, e := range entries {
		if err := w.Insert(e); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}

func runLookup(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run(args, &out, &errOut)
	return out.String(), errOut.String(), code
}

func TestHeadword(t *testing.T) {
	db := testDB(t)
	out, _, code := runLookup(t, "-db", db, "통하다")
	want := "통하다 (通하다) [통하다] 동사 · 중급\n" +
		"  1. go through; pass\n" +
		"     To go through a certain place.\n" +
		"  2. 막힘이 없이 흐르다.\n"
	if code != 0 || out != want {
		t.Errorf("exit %d, got\n%s\nwant\n%s", code, out, want)
	}
}

func TestListedFormAndBrackets(t *testing.T) {
	out, _, code := runLookup(t, "-db", testDB(t), "[통해].")
	if code != 0 || !strings.HasPrefix(out, "통해 → 통하다 (通하다)") {
		t.Errorf("exit %d, got\n%s", code, out)
	}
}

func TestHomographs(t *testing.T) {
	out, _, _ := runLookup(t, "-db", testDB(t), "방금")
	if !strings.Contains(out, "방금¹ 명사 · 초급") || !strings.Contains(out, "방금² 부사 · 초급") {
		t.Errorf("got\n%s", out)
	}
}

func TestPhraseShowsUnit(t *testing.T) {
	out, _, _ := runLookup(t, "-db", testDB(t), "눈이 높다")
	if !strings.HasPrefix(out, "눈이 높다 관용구\n") {
		t.Errorf("got\n%s", out)
	}
}

func TestKoreanFlag(t *testing.T) {
	out, _, _ := runLookup(t, "-db", testDB(t), "-ko", "통하다")
	if !strings.Contains(out, "     어떤 곳을 거쳐 지나가다.\n") {
		t.Errorf("-ko did not add the Korean definition:\n%s", out)
	}
}

func TestNoEntry(t *testing.T) {
	out, _, code := runLookup(t, "-db", testDB(t), "방금", "통해야")
	if code != 1 {
		t.Errorf("exit %d, want 1 when a word is missing", code)
	}
	if !strings.Contains(out, "방금¹") || !strings.Contains(out, "통해야: no entry") {
		t.Errorf("got\n%s", out)
	}
}

func TestMissingDictionary(t *testing.T) {
	_, errOut, code := runLookup(t, "-db", filepath.Join(t.TempDir(), "nope.db"), "말")
	if code != 2 || !strings.Contains(errOut, "run hannote-import") {
		t.Errorf("exit %d, stderr %q", code, errOut)
	}
}

func TestUsage(t *testing.T) {
	if _, _, code := runLookup(t); code != 2 {
		t.Errorf("no words: exit %d, want 2", code)
	}
}

func TestHomonymMark(t *testing.T) {
	for n, want := range map[int]string{0: "", 1: "¹", 2: "²", 12: "¹²"} {
		if got := homonymMark(n); got != want {
			t.Errorf("homonymMark(%d) = %q, want %q", n, got, want)
		}
	}
}
