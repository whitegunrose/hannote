package vocab

import "testing"

func TestStatusRoundTrip(t *testing.T) {
	for _, s := range []Status{Unknown, Seen, Known, Hidden} {
		got, err := ParseStatus(s.String())
		if err != nil || got != s {
			t.Errorf("ParseStatus(%q) = %v, %v; want %v", s.String(), got, err, s)
		}
	}
	if _, err := ParseStatus("famous"); err == nil {
		t.Error("ParseStatus accepted an invalid name")
	}
}
