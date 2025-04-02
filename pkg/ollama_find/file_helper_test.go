package ollama_find

import (
	"path/filepath"
	"testing"
)

/////////////////////
// Tests           //
/////////////////////

func TestFileMissing(t *testing.T) {
	fh := &defaultFileHelper{}
	if fh.FileMissing(filepath.Join("/nonexistent", "file")) != true {
		t.Error("expected file to be missing")
	}
}

func TestExpandPath(t *testing.T) {
	fh := &defaultFileHelper{}
	expandedPath, _ := fh.ExpandPath("~/blah")
	if len(expandedPath) <= 6 {
		t.Error("expected ExpandPath make the test string longer")
	}
}
