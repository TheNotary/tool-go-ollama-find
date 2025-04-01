package main

import (
	"testing"
)

// TODO: finish the test cases up
// func TestHandleFindCommand(t *testing.T) {

// }

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
