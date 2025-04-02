package ollama_find

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// On mac and linux, to ensure program outputs don't contain sensitive
	// information, the absolute paths to files will be truncated with this
	// prefix
	cleanModelDir = "~/.ollama/models"
)

var fileHelper = &defaultFileHelper{}
var modelDir, _ = fileHelper.ExpandPath(cleanModelDir)

// A manifest allows us to unmarshall the important parts of Ollama's manifest
// json
type manifest struct {
	Layers []struct {
		MediaType string `json:"mediaType"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}

////////////////
// Main Logic //
////////////////

// Call this with a modelURI in mind and the path to that model will be returned
// if ollama has cached that model on the local file system
var LookupGGUF = func(modelURI, modelTag string) (string, error) {
	return LookupGGUFPath(modelURI, modelTag, &defaultFileHelper{})
}

// A lower level method than LookupGGUF which allows you to inject your own
// FileHelper (to do it all without a real filesystem)
func LookupGGUFPath(modelURI, modelTag string, fh FileHelper) (string, error) {
	if modelTag == "" {
		modelTag = "latest"
	}

	registryPath, modelName := splitRegistryFromModelName(modelURI)
	pathToManifest := filepath.Join(modelDir, "manifests", registryPath, modelName, modelTag)

	if fh.FileMissing(pathToManifest) {
		msg := fmt.Sprintf("error: Manifest for %s could not be found. Checked %s", modelName, pathToManifest)
		if supplyingATagNameFixesIt(fh, pathToManifest) {
			if taggedFile, err := getExampleTagName(pathToManifest); err == nil {
				msg += fmt.Sprintf(".\n\nIf you meant to specify a version, try:\n  $  %s %s %s", commandName(), modelURI, taggedFile)
			}
		}
		return "", errors.New(msg)
	}

	manifestData, err := fh.ReadManifest(pathToManifest)
	if err != nil {
		return "", fmt.Errorf("error: Unable to parse manifest at %s", pathToManifest)
	}

	var manifest manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return "", fmt.Errorf("error: Unable to parse manifest at %s", pathToManifest)
	}

	digest, err := extractModelDigest(manifest)
	if err != nil {
		return "", fmt.Errorf("error: Unable to extract digest from manifest at %s for model %s", pathToManifest, modelName)
	}

	if fh.IsWindows() {
		absolute_path, err := fh.ExpandPath(filepath.Join(cleanModelDir, "blobs", digest))
		if err != nil {
			fmt.Println("error: Unable to ExpandPath")
		}
		return absolute_path, nil
	}

	return filepath.Join(cleanModelDir, "blobs", digest), nil
}

func commandName() string {
	return filepath.Base(os.Args[0]) + " find"
}

func getExampleTagName(path string) (string, error) {
	dirpath := filepath.Dir(path)
	files, err := os.ReadDir(dirpath)
	if err != nil || len(files) == 0 {
		return "", err
	}
	return files[0].Name(), nil
}

func supplyingATagNameFixesIt(fh FileHelper, path string) bool {
	return fh.DirExist(filepath.Dir(path))
}

func extractModelDigest(manifest manifest) (string, error) {
	for _, layer := range manifest.Layers {
		if layer.MediaType == "application/vnd.ollama.image.model" {
			return strings.Replace(layer.Digest, ":", "-", 1), nil
		}
	}
	return "", errors.New("model digest not found")
}

func splitRegistryFromModelName(modelName string) (string, string) {
	parts := strings.Split(modelName, "/")
	if strings.Contains(parts[0], ".") {
		return getPrivateRegistryModelNameAndRegistry(modelName)
	}
	if len(parts) > 1 {
		subcatalog := parts[0]
		return fmt.Sprintf("registry.ollama.ai/%s", subcatalog), strings.TrimPrefix(modelName, subcatalog+"/")
	}
	return "registry.ollama.ai/library", modelName
}

func getPrivateRegistryModelNameAndRegistry(modelName string) (string, string) {
	parts := strings.Split(modelName, "/")
	modelName = parts[len(parts)-1]
	registryPath := strings.Join(parts[:len(parts)-1], "/")
	return registryPath, modelName
}

func expandPath(path string) (string, error) {
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
