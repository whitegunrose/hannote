// Package dict holds the dictionary model and queries against krdict.db,
// the read-only SQLite database built by hannote-import from the KRDICT XML.
package dict

// Entry is one dictionary headword (a KRDICT LexicalEntry).
type Entry struct {
	ID            int64  // Hannote's own key; KRDICT IDs are not unique
	KrdictID      int64  // LexicalEntry@val
	Headword      string // Lemma/writtenForm, e.g. 일어나다
	Homonym       int    // homonym_number: 방금 has 1 (noun) and 2 (adverb)
	Unit          string // lexicalUnit: 단어, 관용구, 속담, 문법‧표현, 구
	POS           string // partOfSpeech: 명사, 동사, ... (empty for idioms)
	Level         string // vocabularyLevel: 초급, 중급, 고급, 없음
	Hanja         string // origin, stored as-is, e.g. "豫約 하다"
	Pronunciation string // in Hangul; ː marks a long vowel
	Senses        []Sense
	Forms         []string // listed conjugated/contracted forms, e.g. 일어나, 일어나니
}

// Sense is one numbered meaning of an Entry.
type Sense struct {
	Order   int
	DefKo   string // Korean definition
	GlossEn string // short English equivalent, e.g. "make a reservation; book"
	DefEn   string // English definition
}

// Key identifies an entry for vocab.db. The KRDICT ID alone is not unique
// (idioms share their parent word's ID), so the headword is part of the key.
type Key struct {
	KrdictID int64
	Headword string
}

// Key returns the vocab key for e.
func (e Entry) Key() Key { return Key{KrdictID: e.KrdictID, Headword: e.Headword} }
