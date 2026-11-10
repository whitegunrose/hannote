package krdict

import (
	"bufio"
	"io"
)

// cleaner repairs two known defects in the KRDICT dump before the XML
// parser sees them:
//
//  1. Raw '<' or '>' inside an attribute value, e.g. a French gloss
//     "(gudeul <système ...>)" in the 2019 dump. A raw '<' there is not
//     well-formed XML, so it is escaped to &lt; ('>' to &gt; for symmetry).
//  2. ASCII control characters that XML 1.0 forbids everywhere, e.g. the
//     backspaces (U+0008) inside Arabic definitions in the 2026-10-19 dump.
//     They are replaced with a space. Tab, LF and CR are legal and kept.
//
// Everything else passes through unchanged. Both checks are byte-level and
// safe on UTF-8: every byte of a multi-byte character is >= 0x80.
type cleaner struct {
	r       *bufio.Reader
	inTag   bool
	quote   byte   // the open quote character, or 0 when not in a value
	pending []byte // escaped bytes not yet returned
}

func newCleaner(r io.Reader) io.Reader {
	return &cleaner{r: bufio.NewReaderSize(r, 64*1024)}
}

func (c *cleaner) Read(p []byte) (int, error) {
	n := 0
	for n < len(p) {
		if len(c.pending) > 0 {
			k := copy(p[n:], c.pending)
			c.pending = c.pending[k:]
			n += k
			continue
		}
		b, err := c.r.ReadByte()
		if err != nil {
			if n > 0 {
				return n, nil // return what we have; the error comes next call
			}
			return 0, err
		}
		if b < 0x20 && b != '\t' && b != '\n' && b != '\r' {
			b = ' '
		}
		switch {
		case c.quote != 0: // inside an attribute value
			switch b {
			case c.quote:
				c.quote = 0
			case '<':
				c.pending = []byte("&lt;")
				continue
			case '>':
				c.pending = []byte("&gt;")
				continue
			}
		case c.inTag:
			switch b {
			case '"', '\'':
				c.quote = b
			case '>':
				c.inTag = false
			}
		case b == '<':
			c.inTag = true
		}
		p[n] = b
		n++
	}
	return n, nil
}
