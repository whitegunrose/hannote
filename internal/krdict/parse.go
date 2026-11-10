// Package krdict reads the 한국어기초사전 (KRDICT) XML dump.
//
// The dump is LMF: LexicalResource > Lexicon > LexicalEntry, where nearly
// every value is a <feat att="name" val="value"/> child. Parse streams one
// LexicalEntry at a time, so a 35 MB file never sits in memory as a tree.
package krdict

import (
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"strconv"
	"strings"

	"github.com/whitegunrose/hannote/internal/dict"
)

const english = "영어"

type feat struct {
	Att string `xml:"att,attr"`
	Val string `xml:"val,attr"`
}

type feats []feat

// get returns the first value for att, or "".
func (fs feats) get(att string) string {
	for _, f := range fs {
		if f.Att == att {
			return f.Val
		}
	}
	return ""
}

// text returns the value for att as clean display text. Some KRDICT
// strings were escaped twice at the source ("&amp;quot;" in the file), so
// after XML decoding they still read "&quot;"; html.UnescapeString fixes
// that. Surrounding whitespace is trimmed.
func (fs feats) text(att string) string {
	return strings.TrimSpace(html.UnescapeString(fs.get(att)))
}

// all returns every value for att, in order.
func (fs feats) all(att string) []string {
	var out []string
	for _, f := range fs {
		if f.Att == att {
			out = append(out, f.Val)
		}
	}
	return out
}

type xEntry struct {
	ID        string      `xml:"val,attr"`
	Feats     feats       `xml:"feat"`
	Lemmas    []xFeatOnly `xml:"Lemma"`
	WordForms []xWordForm `xml:"WordForm"`
	Senses    []xSense    `xml:"Sense"`
}

type xFeatOnly struct {
	Feats feats `xml:"feat"`
}

type xWordForm struct {
	Feats feats       `xml:"feat"`
	Reps  []xFeatOnly `xml:"FormRepresentation"`
}

type xSense struct {
	ID          string      `xml:"val,attr"`
	Feats       feats       `xml:"feat"`
	Equivalents []xFeatOnly `xml:"Equivalent"`
}

// Parse reads one KRDICT XML file and calls fn for every entry, in file order.
// It stops at the first error from the XML or from fn.
func Parse(r io.Reader, fn func(dict.Entry) error) error {
	dec := xml.NewDecoder(newCleaner(r))
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		start, ok := tok.(xml.StartElement)
		if !ok || start.Name.Local != "LexicalEntry" {
			continue
		}
		var x xEntry
		if err := dec.DecodeElement(&x, &start); err != nil {
			return err
		}
		e, err := convert(x)
		if err != nil {
			return err
		}
		if err := fn(e); err != nil {
			return err
		}
	}
}

func convert(x xEntry) (dict.Entry, error) {
	id, err := strconv.ParseInt(x.ID, 10, 64)
	if err != nil {
		return dict.Entry{}, fmt.Errorf("krdict: bad entry id %q: %w", x.ID, err)
	}
	e := dict.Entry{
		KrdictID: id,
		Unit:     x.Feats.get("lexicalUnit"),
		POS:      x.Feats.get("partOfSpeech"),
		Level:    x.Feats.get("vocabularyLevel"),
		Hanja:    x.Feats.get("origin"),
	}
	if h := x.Feats.get("homonym_number"); h != "" {
		e.Homonym, _ = strconv.Atoi(h)
	}
	for _, l := range x.Lemmas {
		if w := l.Feats.get("writtenForm"); w != "" {
			e.Headword = strings.TrimSpace(html.UnescapeString(w))
			break
		}
	}
	if e.Headword == "" {
		return dict.Entry{}, fmt.Errorf("krdict: entry %d has no headword", id)
	}

	var prons []string
	seen := map[string]bool{}
	addForm := func(w string) {
		w = strings.TrimSpace(w)
		if w != "" && w != e.Headword && !seen[w] {
			seen[w] = true
			e.Forms = append(e.Forms, w)
		}
	}
	for _, wf := range x.WordForms {
		switch wf.Feats.get("type") {
		case "발음":
			prons = append(prons, wf.Feats.all("pronunciation")...)
		case "활용":
			addForm(wf.Feats.get("writtenForm"))
		}
		for _, rep := range wf.Reps { // contracted forms (준말)
			addForm(rep.Feats.get("writtenForm"))
		}
	}
	e.Pronunciation = strings.Join(prons, ", ")

	for i, s := range x.Senses {
		sense := dict.Sense{Order: i + 1, DefKo: s.Feats.text("definition")}
		if n, err := strconv.Atoi(s.ID); err == nil {
			sense.Order = n
		}
		for _, eq := range s.Equivalents {
			if eq.Feats.get("language") == english {
				sense.GlossEn = eq.Feats.text("lemma")
				sense.DefEn = eq.Feats.text("definition")
				break
			}
		}
		e.Senses = append(e.Senses, sense)
	}
	return e, nil
}
