package paths

import (
	"path/filepath"
	"testing"
)

func TestDataDirPrecedence(t *testing.T) {
	t.Setenv("HANNOTE_DATA_DIR", "/tmp/override")
	t.Setenv("XDG_DATA_HOME", "/tmp/xdg")
	if d, _ := DataDir(); d != "/tmp/override" {
		t.Errorf("with HANNOTE_DATA_DIR: got %q", d)
	}

	t.Setenv("HANNOTE_DATA_DIR", "")
	if d, _ := DataDir(); d != filepath.Join("/tmp/xdg", "hannote") {
		t.Errorf("with XDG_DATA_HOME: got %q", d)
	}
}
