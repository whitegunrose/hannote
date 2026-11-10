package krdict

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/whitegunrose/hannote/internal/dict"
)

func parseFile(t *testing.T, path string) []dict.Entry {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var got []dict.Entry
	if err := Parse(f, func(e dict.Entry) error { got = append(got, e); return nil }); err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return got
}

func TestParseSample(t *testing.T) {
	got := parseFile(t, "../../testdata/krdict_sample.xml")
	if len(got) != 3 {
		t.Fatalf("got %d entries, want 3", len(got))
	}

	ireonada := got[0]
	if ireonada.Headword != "일어나다" || ireonada.KrdictID != 71696 || ireonada.POS != "동사" || ireonada.Level != "초급" {
		t.Errorf("일어나다 header = %+v", ireonada)
	}
	if ireonada.Pronunciation != "이러나다" {
		t.Errorf("pronunciation = %q", ireonada.Pronunciation)
	}
	if strings.Join(ireonada.Forms, ",") != "일어나,일어나니" {
		t.Errorf("forms = %v", ireonada.Forms)
	}
	if len(ireonada.Senses) != 2 {
		t.Fatalf("senses = %d, want 2", len(ireonada.Senses))
	}
	// English is picked even when another language comes first, and
	// double-escaped quotes (&amp;quot; in the file) come out as plain quotes.
	if s := ireonada.Senses[1]; s.Order != 2 || s.GlossEn != "get up" || s.DefEn != `To "wake" up.` {
		t.Errorf("sense 2 = %+v", s)
	}

	yeyak := got[1]
	if yeyak.Hanja != "豫約 하다" || yeyak.Senses[0].GlossEn != "make a reservation; book" {
		t.Errorf("예약하다 = %+v", yeyak)
	}

	// Idioms reuse their parent's ID and have no part of speech or English.
	idiom := got[2]
	if idiom.Unit != "관용구" || idiom.POS != "" || idiom.Senses[0].GlossEn != "" {
		t.Errorf("idiom = %+v", idiom)
	}
}

func TestCleanerEscapesAttributes(t *testing.T) {
	in := `<a b="x <y> z" c='<q>'>text > <i/></a>`
	want := `<a b="x &lt;y&gt; z" c='&lt;q&gt;'>text > <i/></a>`
	got, err := io.ReadAll(newCleaner(strings.NewReader(in)))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}

func TestCleanerSmallReads(t *testing.T) {
	// Escapes must survive being split across tiny Read calls.
	r := newCleaner(strings.NewReader(`<a b="<>">`))
	var out []byte
	buf := make([]byte, 1)
	for {
		n, err := r.Read(buf)
		out = append(out, buf[:n]...)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	if string(out) != `<a b="&lt;&gt;">` {
		t.Errorf("got %s", out)
	}
}

func TestCleanerReplacesControlChars(t *testing.T) {
	in := "<a b=\"x\by\">\tok\r\n\x00</a>"
	want := "<a b=\"x y\">\tok\r\n </a>"
	got, err := io.ReadAll(newCleaner(strings.NewReader(in)))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
