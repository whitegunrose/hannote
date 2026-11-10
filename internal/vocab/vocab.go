// Package vocab stores the user's own data in vocab.db: word statuses,
// notes on entries, and custom entries for words KRDICT doesn't have.
package vocab

import "fmt"

// Status is how well the user knows a word.
type Status int

const (
	Unknown Status = iota
	Seen
	Known
	Hidden
)

var statusNames = [...]string{"unknown", "seen", "known", "hidden"}

func (s Status) String() string {
	if s < Unknown || s > Hidden {
		return fmt.Sprintf("Status(%d)", int(s))
	}
	return statusNames[s]
}

// ParseStatus parses the name stored in vocab.db.
func ParseStatus(name string) (Status, error) {
	for i, n := range statusNames {
		if n == name {
			return Status(i), nil
		}
	}
	return Unknown, fmt.Errorf("vocab: unknown status %q", name)
}

// CustomEntry is a user-made entry for a word with no KRDICT match.
type CustomEntry struct {
	ID       int64
	Headword string
	Gloss    string
	Notes    string
}
