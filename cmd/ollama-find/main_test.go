package main

import (
	"testing"
)

func TestParseModelNameAndTag(t *testing.T) {
	tests := []struct {
		arg1      string
		arg2      string
		expected1 string
		expected2 string
	}{
		{"model", "tag", "model", "tag"},
		{"model", "", "model", "latest"},
		{"model:7b", "", "model", "7b"},
		{" model ", " tag ", "model", "tag"},
		{"model", " ", "model", "latest"},
	}

	for _, tt := range tests {
		got1, got2 := ParseModelNameAndTag(tt.arg1, tt.arg2)
		if got1 != tt.expected1 || got2 != tt.expected2 {
			t.Errorf("ParseModelNameAndTag(%q, %q) = (%q, %q); want (%q, %q)", tt.arg1, tt.arg2, got1, got2, tt.expected1, tt.expected2)
		}
	}
}
