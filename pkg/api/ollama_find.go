package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	CleanModelDir = "~/.ollama/models"
)

var ModelDir, _ = filepath.Abs(CleanModelDir)

type Manifest struct {
	Layers []struct {
		MediaType string `json:"mediaType"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}

////////////////////////////
// FileHelper for testing //
////////////////////////////

// FileHelper defines an interface for file operations.
type FileHelper interface {
	FileMissing(path string) bool
	ReadManifest(path string) ([]byte, error)
}

// DefaultFileHelper provides real implementations of file operations.
type DefaultFileHelper struct{}

func (DefaultFileHelper) FileMissing(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func (DefaultFileHelper) ReadManifest(path string) ([]byte, error) {
	return os.ReadFile(filepath.Clean(path))
}

////////////////
// Main Logic //
////////////////

// wrapper for LookupGGUFPath
func LookupGGUF(modelURI, modelTag string) (string, error) {
	return LookupGGUFPath(modelURI, modelTag, &DefaultFileHelper{})
}

func LookupGGUFPath(modelURI, modelTag string, fh FileHelper) (string, error) {
	if modelTag == "" {
		modelTag = "latest"
	}

	registryPath, modelName := SplitRegistryFromModelName(modelURI)
	pathToManifest := filepath.Join(ModelDir, "manifests", registryPath, modelName, modelTag)

	if fh.FileMissing(pathToManifest) {
		msg := fmt.Sprintf("error: Manifest for %s could not be found. Checked %s", modelName, CleanPath(pathToManifest))
		if LooksLikeTagNameNeeded(pathToManifest) {
			if taggedFile, err := GetExampleTagName(pathToManifest); err == nil {
				msg += fmt.Sprintf(".\n\nIf you meant to specify a version, try:\n  $  %s %s %s", CommandName(), modelURI, taggedFile)
			}
		}
		return "", errors.New(msg)
	}

	manifestData, err := fh.ReadManifest(pathToManifest)
	if err != nil {
		return "", fmt.Errorf("error: Unable to parse manifest at %s", pathToManifest)
	}

	var manifest Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return "", fmt.Errorf("error: Unable to parse manifest at %s", pathToManifest)
	}

	digest, err := ExtractModelDigest(manifest)
	if err != nil {
		return "", fmt.Errorf("error: Unable to extract digest from manifest at %s for model %s", pathToManifest, modelName)
	}

	if isWindows() {
		absolute_path, err := ExpandPath(filepath.Join(CleanModelDir, "blobs", digest))
		if err != nil {
			fmt.Println("error: Unable to ExpandPath")
		}
		return absolute_path, nil
	}

	return filepath.Join(CleanModelDir, "blobs", digest), nil
}

func CommandName() string {
	return filepath.Base(os.Args[0]) + " find"
}

func GetExampleTagName(path string) (string, error) {
	dirpath := filepath.Dir(path)
	files, err := os.ReadDir(dirpath)
	if err != nil || len(files) == 0 {
		return "", err
	}
	return files[0].Name(), nil
}

func LooksLikeTagNameNeeded(path string) bool {
	return DirExists(filepath.Dir(path))
}

func CleanPath(absolutePath string) string {
	if isWindows() {
		return absolutePath
	}

	return strings.Replace(absolutePath, ModelDir, CleanModelDir, 1)
}

func ExtractModelDigest(manifest Manifest) (string, error) {
	for _, layer := range manifest.Layers {
		if layer.MediaType == "application/vnd.ollama.image.model" {
			return strings.Replace(layer.Digest, ":", "-", 1), nil
		}
	}
	return "", errors.New("model digest not found")
}

func SplitRegistryFromModelName(modelName string) (string, string) {
	parts := strings.Split(modelName, "/")
	if strings.Contains(parts[0], ".") {
		return GetPrivateRegistryModelNameAndRegistry(modelName)
	}
	if len(parts) > 1 {
		subcatalog := parts[0]
		return fmt.Sprintf("registry.ollama.ai/%s", subcatalog), strings.TrimPrefix(modelName, subcatalog+"/")
	}
	return "registry.ollama.ai/library", modelName
}

func GetPrivateRegistryModelNameAndRegistry(modelName string) (string, string) {
	parts := strings.Split(modelName, "/")
	modelName = parts[len(parts)-1]
	registryPath := strings.Join(parts[:len(parts)-1], "/")
	return registryPath, modelName
}

func ReadManifest(pathToManifest string) ([]byte, error) {
	return os.ReadFile(filepath.Clean(pathToManifest))
}

func FileMissing(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

func ExpandPath(path string) (string, error) {
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
