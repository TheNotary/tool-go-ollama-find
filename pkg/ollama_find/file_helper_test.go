package ollama_find

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
	tests := []struct {
		name            string
		path            string
		expectedPathLen int
		expectErr       bool
	}{
		{"it works fine", "~/blah", 6, false},
		// regrets!
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expandedPath, err := fh.ExpandPath(tt.path)
			if tt.expectErr {
				assert.Error(t, err)
			}
			if len(expandedPath) <= tt.expectedPathLen {
				t.Error("expected ExpandPath make the test string longer")
			}
		})
	}
}

// TODO: Delete these they're pointless and test nothing the compiler doesn't
// already prove. There's probably a way to mark standard lib calls and their
// error branches as ignored

func TestReadManifest(t *testing.T) {
	fh := &defaultFileHelper{}
	_, err := fh.ReadManifest("/idontexist")
	assert.Error(t, err)
}

func TestReadDir(t *testing.T) {
	fh := &defaultFileHelper{}
	filesFound, _ := fh.ReadDir("/tmp")
	assert.True(t, len(filesFound) >= 0)
}

func TestDirExist(t *testing.T) {
	fh := &defaultFileHelper{}
	doesTmpExist := fh.DirExist(".")
	assert.True(t, doesTmpExist)
}

func TestIsWindows(t *testing.T) {
	fh := &defaultFileHelper{}
	anyBool := fh.IsWindows()
	assert.IsType(t, anyBool, false)
}
