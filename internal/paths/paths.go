// Package paths says where Hannote keeps its data.
package paths

import (
	"os"
	"path/filepath"
)

// DataDir returns $HANNOTE_DATA_DIR if set, else $XDG_DATA_HOME/hannote,
// else ~/.local/share/hannote.
func DataDir() (string, error) {
	if d := os.Getenv("HANNOTE_DATA_DIR"); d != "" {
		return d, nil
	}
	if d := os.Getenv("XDG_DATA_HOME"); d != "" {
		return filepath.Join(d, "hannote"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "hannote"), nil
}

// DictDB returns the path of krdict.db.
func DictDB() (string, error) { return inData("krdict.db") }

// VocabDB returns the path of vocab.db.
func VocabDB() (string, error) { return inData("vocab.db") }

func inData(name string) (string, error) {
	d, err := DataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, name), nil
}
