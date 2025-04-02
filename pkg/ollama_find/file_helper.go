package ollama_find

import (
	"os"
	"path/filepath"
	"strings"
)

////////////////////////////
// FileHelper for testing //
////////////////////////////

// FileHelper defines an interface for file operations.
type FileHelper interface {
	FileMissing(path string) bool
	DirExist(path string) bool
	ReadDir(path string) ([]os.DirEntry, error)
	ReadManifest(path string) ([]byte, error)
	IsWindows() bool
	ExpandPath(path string) (string, error)
}

// defaultFileHelper provides implementations of file operations many of which
// are handy to mock in tests.
type defaultFileHelper struct{}

// Returns true if the path supplied does not exist
func (defaultFileHelper) FileMissing(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func (defaultFileHelper) DirExist(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func (defaultFileHelper) ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func (defaultFileHelper) ReadManifest(path string) ([]byte, error) {
	return os.ReadFile(filepath.Clean(path))
}

func (defaultFileHelper) IsWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

// Replaces the leading "~" in a string with the users HOME path (even on window)
// as well as calls filepath.Abs on the path supplied
func (defaultFileHelper) ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(homeDir, path[1:])
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
