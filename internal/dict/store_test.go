package dict

import (
	"path/filepath"
	"testing"
)

func TestWriteAndLookup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "krdict.db")
	w, err := Create(path)
	if err != nil {
		t.Fatal(err)
	}
	entries := []Entry{
		{KrdictID: 92065, Headword: "방금", Homonym: 2, Unit: "단어", POS: "부사", Level: "초급", Hanja: "方今",
			Senses: []Sense{{Order: 1, DefKo: "말하고 있는 시점보다 바로 조금 전에.", GlossEn: "a moment ago"}}},
		{KrdictID: 57833, Headword: "방금", Homonym: 1, Unit: "단어", POS: "명사", Level: "초급", Hanja: "方今",
			Senses: []Sense{{Order: 1, GlossEn: "moment ago"}}},
		{KrdictID: 71696, Headword: "일어나다", Unit: "단어", POS: "동사", Pronunciation: "이러나다",
			Forms:  []string{"일어나", "일어나니"},
			Senses: []Sense{{Order: 2, GlossEn: "get up"}, {Order: 1, GlossEn: "stand up"}}},
	}
	for _, e := range entries {
		if err := w.Insert(e); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if w.Entries != 3 || w.Senses != 4 || w.Forms != 2 {
		t.Errorf("counts = %d entries, %d senses, %d forms", w.Entries, w.Senses, w.Forms)
	}

	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	got, err := s.Lookup("방금")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Homonym != 1 || got[1].Homonym != 2 {
		t.Fatalf("방금 = %+v, want both homographs in order", got)
	}

	got, err = s.Lookup("일어나")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Headword != "일어나다" {
		t.Fatalf("일어나 = %+v, want 일어나다 via its form", got)
	}
	if got[0].Senses[0].GlossEn != "stand up" {
		t.Errorf("senses not ordered: %+v", got[0].Senses)
	}

	if got, _ := s.Lookup("없는말"); len(got) != 0 {
		t.Errorf("unknown word returned %+v", got)
	}
}

func TestOpenRejectsNonDictionary(t *testing.T) {
	if _, err := Open(filepath.Join(t.TempDir(), "missing.db")); err == nil {
		t.Error("Open succeeded on a missing file")
	}
}
