package api_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/thenotary/tool-go-ollama-find/pkg/api"
)

type MockOllamaFind struct {
	mock.Mock
}

func (m *MockOllamaFind) FileMissing(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

func (m *MockOllamaFind) ReadManifest(path string) ([]byte, error) {
	args := m.Called(path)
	return args.Get(0).([]byte), args.Error(1)
}

var testManifestPath = "./testdata/ollama_manifest.json"
var manifestData, _ = os.ReadFile(filepath.Clean(testManifestPath))

func TestLookupGGUFPath(t *testing.T) {
	assert := assert.New(t)

	modelName := "deepseek-r1"

	mockFind := new(MockOllamaFind)
	mockFind.On("FileMissing", mock.Anything).Return(false)
	mockFind.On("ReadManifest", mock.Anything).Return(manifestData, nil)

	path, err := api.LookupGGUFPath(modelName, "", mockFind)
	assert.NoError(err)
	assert.Equal("~/.ollama/models/blobs/sha256-96c415656d377afbff962f6cdb2394ab092ccbcbaab4b82525bc4ca800fe8a49", path)
}

func TestLookupGGUFPathUgly(t *testing.T) {
	cases := []struct {
		name             string
		modelURI         string
		modelTag         string
		fileMissingValue bool
		expectErr        bool
	}{
		{"Valid model with tag", "mymodel", "v1.0", false, false},
		{"Valid model without tag", "mymodel", "", false, false},
		{"Invalid model", "unknownmodel", "", true, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockFind := new(MockOllamaFind)
			mockFind.On("FileMissing", mock.Anything).Return(tc.fileMissingValue)
			mockFind.On("ReadManifest", mock.Anything).Return(manifestData, nil)

			_, err := api.LookupGGUFPath(tc.modelURI, tc.modelTag, mockFind)

			if (err != nil) != tc.expectErr {
				t.Errorf("expected error: %v, got: %v", tc.expectErr, err)
			}
		})
	}
}

func TestSplitRegistryFromModelName(t *testing.T) {
	cases := []struct {
		input       string
		expRegistry string
		expModel    string
	}{
		{"registry.com/user/mymodel", "registry.com/user", "mymodel"},
		{"subcatalog/mymodel", "registry.ollama.ai/subcatalog", "mymodel"},
		{"mymodel", "registry.ollama.ai/library", "mymodel"},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			reg, model := api.SplitRegistryFromModelName(tc.input)
			if reg != tc.expRegistry || model != tc.expModel {
				t.Errorf("expected (%s, %s), got (%s, %s)", tc.expRegistry, tc.expModel, reg, model)
			}
		})
	}
}

func TestExtractModelDigest(t *testing.T) {
	manifest := api.Manifest{
		Layers: []struct {
			MediaType string `json:"mediaType"`
			Digest    string `json:"digest"`
		}{
			{MediaType: "application/vnd.ollama.image.model", Digest: "sha256:abc123"},
		},
	}
	digest, err := api.ExtractModelDigest(manifest)
	if err != nil || digest != "sha256-abc123" {
		t.Errorf("expected sha256-abc123, got %s, err: %v", digest, err)
	}
}

func TestFileMissing(t *testing.T) {
	if api.FileMissing(filepath.Join("/nonexistent", "file")) != true {
		t.Error("expected file to be missing")
	}
}

func TestExpandPath(t *testing.T) {
	expandedPath, _ := api.ExpandPath("~/blah")
	if len(expandedPath) <= 6 {
		t.Error("expected ExpandPath make the test string longer")
	}
}
