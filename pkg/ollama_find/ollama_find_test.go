package ollama_find

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

///////////
// Mocks //
///////////

type MockOllamaFind struct {
	mock.Mock
}

func (m *MockOllamaFind) FileMissing(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

func (m *MockOllamaFind) DirExist(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

func (m *MockOllamaFind) ReadManifest(path string) ([]byte, error) {
	args := m.Called(path)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockOllamaFind) IsWindows() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockOllamaFind) ExpandPath(path string) (string, error) {
	return defaultFileHelper{}.ExpandPath(path)
}

func (m *MockOllamaFind) ReadDir(path string) ([]os.DirEntry, error) {
	if path == filepath.Join("/", "empty") {
		return []os.DirEntry{}, nil
	}
	return os.ReadDir(".")
}

//////////
// Vars //
//////////

var homeDir, _ = os.UserHomeDir()

var testManifestPath = "./testdata/ollama_manifest.json"
var manifestData, _ = os.ReadFile(filepath.Clean(testManifestPath))
var expectedWindowsBlobPath = filepath.Join(homeDir, ".ollama", "models", "blobs",
	"sha256-96c415656d377afbff962f6cdb2394ab092ccbcbaab4b82525bc4ca800fe8a49")
var expectedNormalBlobPath = filepath.Join("~", ".ollama", "models", "blobs",
	"sha256-96c415656d377afbff962f6cdb2394ab092ccbcbaab4b82525bc4ca800fe8a49")

///////////
// Tests //
///////////

func TestLookupGGUF(t *testing.T) {
	assert := assert.New(t)

	output, err := LookupGGUF("missing-model", "latest")

	assert.Contains(err.Error(), "error: Manifest for")
	assert.Equal("", output)
}

func TestLookupGGUFPath(t *testing.T) {
	assert := assert.New(t)

	modelName := "deepseek-r1"

	mockFind := new(MockOllamaFind)
	mockFind.On("IsWindows", mock.Anything).Return(false)
	mockFind.On("FileMissing", mock.Anything).Return(false)
	mockFind.On("ReadManifest", mock.Anything).Return(manifestData, nil)

	path, err := LookupGGUFPath(modelName, "", mockFind)

	assert.NoError(err)
	assert.Equal(expectedNormalBlobPath, path)
}

func TestOutputsForWindows(t *testing.T) {
	assert := assert.New(t)

	modelName := "deepseek-r1"

	mockFind := new(MockOllamaFind)
	mockFind.On("IsWindows", mock.Anything).Return(true)
	mockFind.On("FileMissing", mock.Anything).Return(false)
	mockFind.On("ReadManifest", mock.Anything).Return(manifestData, nil)

	path, err := LookupGGUFPath(modelName, "", mockFind)

	assert.NoError(err)
	assert.Equal(expectedWindowsBlobPath, path)
}

func TestLookupGGUFPathUgly(t *testing.T) {
	cases := []struct {
		name     string
		modelURI string
		modelTag string

		fileMissing    bool
		unreadableFile bool
		unparsableFile bool
		unusableFile   bool

		expectErr bool
	}{
		{"Valid model with tag", "mymodel", "v1.0",
			false, false, false, false,
			false},
		{"Valid model without tag", "mymodel", "",
			false, false, false, false,
			false},

		{"fileMissing for model's manifest", "unknownmodel", "",
			true, false, false, false,
			true},
		{"Unreadable manifest", "mymodel", "",
			false, true, false, false,
			true},
		{"unparsableFile", "mymodel", "",
			false, false, true, false,
			true},
		{"manifest missing keys", "mymodel", "",
			false, false, false, true,
			true},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockFind := new(MockOllamaFind)
			mockFind.On("IsWindows", mock.Anything).Return(false)
			mockFind.On("FileMissing", mock.Anything).Return(tt.fileMissing)

			if tt.unreadableFile {
				mockFind.On("ReadManifest", mock.Anything).Return([]byte{}, errors.New("error"))
			} else {
				if tt.unparsableFile {
					mockFind.On("ReadManifest", mock.Anything).Return([]byte("..."), nil)
				} else if tt.unusableFile {
					mockFind.On("ReadManifest", mock.Anything).Return([]byte("{}"), nil)
				} else {
					mockFind.On("ReadManifest", mock.Anything).Return(manifestData, nil)
				}
			}

			mockFind.On("DirExist", mock.Anything).Return(true)

			_, err := LookupGGUFPath(tt.modelURI, tt.modelTag, mockFind)

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

func TestCommandName(t *testing.T) {
	name := commandName()
	assert.IsType(t, name, "")
}

func TestGetTagNameSuggestion(t *testing.T) {
	fh := new(MockOllamaFind)

	suggestion, err := getTagNameSuggestion(fh, filepath.Join("/empty", ".keep"))

	assert.Empty(t, suggestion)
	assert.NoError(t, err)
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
			reg, model := splitRegistryFromModelName(tc.input)
			if reg != tc.expRegistry || model != tc.expModel {
				t.Errorf("expected (%s, %s), got (%s, %s)", tc.expRegistry, tc.expModel, reg, model)
			}
		})
	}
}

func TestExtractModelDigest(t *testing.T) {
	manifest := manifest{
		Layers: []struct {
			MediaType string `json:"mediaType"`
			Digest    string `json:"digest"`
		}{
			{MediaType: "application/vnd.ollama.image.model", Digest: "sha256:abc123"},
		},
	}
	digest, err := extractModelDigest(manifest)
	if err != nil || digest != "sha256-abc123" {
		t.Errorf("expected sha256-abc123, got %s, err: %v", digest, err)
	}
}
