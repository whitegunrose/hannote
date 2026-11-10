// Command hannote-import builds krdict.db from the KRDICT XML download.
//
// Usage:
//
//	hannote-import [-o path/to/krdict.db] [-verify] <dir-with-xml-files>
//
// It writes to a temporary file next to the output and renames it into place
// only when the whole import succeeds, so a failed run never leaves a broken
// krdict.db behind.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/whitegunrose/hannote/internal/dict"
	"github.com/whitegunrose/hannote/internal/krdict"
	"github.com/whitegunrose/hannote/internal/paths"
)

// verifyWords are looked up after an import with -verify.
var verifyWords = []string{"주로", "방금", "예약하다", "칭찬하다", "일어나"}

func main() {
	out := flag.String("o", "", "output database (default: data dir/krdict.db)")
	verify := flag.Bool("verify", false, "look up a few known words after importing")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: hannote-import [-o krdict.db] [-verify] <xml-dir>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(flag.Arg(0), *out, *verify); err != nil {
		fmt.Fprintln(os.Stderr, "hannote-import:", err)
		os.Exit(1)
	}
}

func run(dir, out string, verify bool) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.xml"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no .xml files in %s", dir)
	}
	sortNatural(files)

	if out == "" {
		if out, err = paths.DictDB(); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	tmp := out + ".tmp"
	os.Remove(tmp)

	start := time.Now()
	w, err := dict.Create(tmp)
	if err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok {
			w.Abort()
			os.Remove(tmp)
		}
	}()

	for i, path := range files {
		before := w.Entries
		if err := importFile(w, path); err != nil {
			return fmt.Errorf("%s: %w", filepath.Base(path), err)
		}
		fmt.Fprintf(os.Stderr, "[%d/%d] %-14s %6d entries\n", i+1, len(files), filepath.Base(path), w.Entries-before)
	}
	if err := w.SetMeta("imported_at", time.Now().UTC().Format(time.RFC3339)); err != nil {
		return err
	}
	if err := w.SetMeta("source_files", fmt.Sprint(len(files))); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	ok = true
	if err := os.Rename(tmp, out); err != nil {
		return err
	}

	pct := 0.0
	if w.Senses > 0 {
		pct = 100 * float64(w.Eng) / float64(w.Senses)
	}
	fmt.Printf("Imported %d entries, %d senses (%d with English, %.0f%%), %d forms in %s\n",
		w.Entries, w.Senses, w.Eng, pct, w.Forms, time.Since(start).Round(time.Second/10))
	fmt.Printf("Wrote %s\n", out)

	if verify {
		return verifyLookups(out)
	}
	return nil
}

// sortNatural orders files by their leading number when they have one
// (KRDICT names them 1_5000_….xml … 11_3671_….xml), so 2 comes before 10.
// Names without a leading number sort after, alphabetically.
func sortNatural(files []string) {
	lead := func(path string) (int, bool) {
		base := filepath.Base(path)
		end := strings.IndexFunc(base, func(r rune) bool { return r < '0' || r > '9' })
		if end <= 0 {
			return 0, false
		}
		n, err := strconv.Atoi(base[:end])
		return n, err == nil
	}
	sort.SliceStable(files, func(i, j int) bool {
		a, aok := lead(files[i])
		b, bok := lead(files[j])
		if aok != bok {
			return aok
		}
		if aok && a != b {
			return a < b
		}
		return files[i] < files[j]
	})
}

func importFile(w *dict.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return krdict.Parse(f, w.Insert)
}

func verifyLookups(path string) error {
	s, err := dict.Open(path)
	if err != nil {
		return err
	}
	defer s.Close()
	fmt.Println()
	for _, word := range verifyWords {
		entries, err := s.Lookup(word)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			fmt.Printf("%s: no entry\n", word)
			continue
		}
		for _, e := range entries {
			var glosses []string
			for _, sn := range e.Senses {
				if sn.GlossEn != "" {
					glosses = append(glosses, sn.GlossEn)
				}
			}
			label := e.Headword
			if e.Hanja != "" {
				label += " (" + e.Hanja + ")"
			}
			fmt.Printf("%s → %s [%s] %s: %s\n", word, label, e.Pronunciation, e.POS, strings.Join(glosses, " | "))
		}
	}
	return nil
}
