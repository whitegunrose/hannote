package main

import (
	"strings"
	"testing"
)

func TestSortNatural(t *testing.T) {
	files := []string{
		"x/10_5000_20261019.xml", "x/2_5000_20261019.xml", "x/notes.xml",
		"x/1_5000_20261019.xml", "x/11_3671_20261019.xml", "x/5000.xml",
	}
	sortNatural(files)
	got := strings.Join(files, " ")
	want := "x/1_5000_20261019.xml x/2_5000_20261019.xml x/10_5000_20261019.xml " +
		"x/11_3671_20261019.xml x/5000.xml x/notes.xml"
	if got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}
