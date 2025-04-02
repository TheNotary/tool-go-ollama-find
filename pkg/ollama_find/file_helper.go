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
	ReadManifest(path string) ([]byte, error)
	IsWindows() bool
	ExpandPath(path string) (string, error)
}

// DefaultFileHelper provides implementations of file operations many of which
// are handy to mock in tests.
type DefaultFileHelper struct{}

// Returns true if the path supplied does not exist
func (DefaultFileHelper) FileMissing(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func (DefaultFileHelper) DirExist(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func (DefaultFileHelper) ReadManifest(path string) ([]byte, error) {
	return os.ReadFile(filepath.Clean(path))
}

func (DefaultFileHelper) IsWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

// Replaces the leading "~" in a string with the users HOME path (even on window)
// as well as calls filepath.Abs on the path supplied
func (DefaultFileHelper) ExpandPath(path string) (string, error) {
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
