package main

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/thenotary/tool-go-ollama-find/pkg/ollama_find"
)

//////////////////
// Test Helpers //
//////////////////

var originalArgs = os.Args
var old = os.Stdout
var originalLookupGGUF = ollama_find.LookupGGUF

func stubbedLookupGGUF(modelName, modelTag string) (string, error) {
	if modelName == "valid-model" {
		return "/path/to/model.gguf", nil
	}
	return "", errors.New("lookup error")
}

func SetupStdoutForTesting() (*os.File, *os.File) {
	stdReadEnd, stdWriteEnd, _ := os.Pipe()
	os.Stdout = stdWriteEnd
	return stdReadEnd, stdWriteEnd
}

func CollectStringAndRestoreStdout(stdReadEnd, stdWriteEnd *os.File) string {
	stdWriteEnd.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(stdReadEnd)
	return buf.String()
}

///////////
// Tests //
///////////

func TestMainFunction(t *testing.T) {
	// Mock sticky things
	os.Args = []string{"_", "valid-model"}
	defer func() { os.Args = originalArgs }()
	ollama_find.LookupGGUF = stubbedLookupGGUF
	defer func() { ollama_find.LookupGGUF = originalLookupGGUF }()

	// Capture stdout
	stdReadEnd, stdWriteEnd := SetupStdoutForTesting()

	// Run main
	main()

	// Fix stdout
	output := CollectStringAndRestoreStdout(stdReadEnd, stdWriteEnd)

	// Expected output depends on HandleFindCommand
	expected := "/path/to/model.gguf\n" // Adjust based on expected behavior

	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestHandleFindCommand(t *testing.T) {
	tests := []struct {
		name            string
		model           string
		expectedSuccess bool
		expectedString  string
	}{
		{
			"it prints the expected thing for a valid model",
			"valid-model", true, "/path/to/model.gguf\n"},
		{
			"it prints an error message when an invalid model is passed in",
			"invalid-model", false, "error: something went wrong calling LookupGGUF lookup error\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ollama_find.LookupGGUF = stubbedLookupGGUF
			defer func() { ollama_find.LookupGGUF = originalLookupGGUF }()

			// Capture stdout
			stdReadEnd, stdWriteEnd := SetupStdoutForTesting()

			// Exercise code
			success := HandleFindCommand([]string{tt.model})

			// Put stdout back to normal
			output := CollectStringAndRestoreStdout(stdReadEnd, stdWriteEnd)

			if success != tt.expectedSuccess {
				t.Errorf("Expected success to be %t", tt.expectedSuccess)
			}

			// Validate output
			if output != tt.expectedString {
				t.Errorf("Expected %q, got %q", tt.expectedString, output)
			}
		})
	}
}

func TestParseModelNameAndTag(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expected1 string
		expected2 string
	}{
		{
			"model and tag are specified in separate arguments",
			[]string{"model", "7b"}, "model", "7b"},
		{
			"only model is specified",
			[]string{"model"}, "model", "latest"},
		{
			"model and tag are specified in the same argument",
			[]string{"model:7b"}, "model", "7b"},
		{
			"model and tag are specified in separate arguments and have spaces",
			[]string{" model ", " 7b "}, "model", "7b"},
		{
			"tag is specified as an empty space",
			[]string{"model", " "}, "model", "latest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := ParseModelNameAndTag(tt.args)
			if got1 != tt.expected1 || got2 != tt.expected2 {
				t.Errorf("ParseModelNameAndTag(%q) => (%q, %q); want (%q, %q)", tt.args, got1, got2, tt.expected1, tt.expected2)
			}
		})

	}
}

func TestHandleHelpCommand(t *testing.T) {
	tests := []struct {
		args              []string
		expectedHandled   bool
		expectedStdoutLen int
	}{
		{[]string{"--help"}, true, 246},
		{[]string{"help"}, true, 246},
		{[]string{"halp"}, true, 246},
		{[]string{"HELP"}, false, 0}, // plz don't shout at my computer programs
		{[]string{"-h"}, true, 246},
	}

	for _, tt := range tests {
		stdReadEnd, stdWriteEnd := SetupStdoutForTesting()

		// Excercise code
		handled := HandleHelpCommand(tt.args)

		output := CollectStringAndRestoreStdout(stdReadEnd, stdWriteEnd)

		if handled != tt.expectedHandled {
			t.Errorf("HandleHelpCommand(%q) => (%t); want (%t)", tt.args, handled, tt.expectedHandled)
		}

		if len(output) != tt.expectedStdoutLen {
			t.Errorf("HandleHelpCommand(%q) did not produce output of the expected length (%d).  want (%d)", tt.args, tt.expectedStdoutLen, len(output))
		}
	}

}
