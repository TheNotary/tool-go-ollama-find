package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/thenotary/tool-go-ollama-find/pkg/api"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Println(Help())
		return
	}

	var arg1, arg2 string
	if len(args) > 0 {
		arg1 = args[0]
	}
	if len(args) > 1 {
		arg2 = args[1]
	}

	modelName, modelTag := ParseModelNameAndTag(arg1, arg2)

	fmt.Println(api.LookupGGUF(modelName, modelTag))
}

func ParseModelNameAndTag(arg1 string, arg2 string) (string, string) {
	// args := os.Args[1:]
	modelName := strings.TrimSpace(arg1)
	var modelTag string

	if len(arg2) > 1 {
		modelTag = strings.TrimSpace(arg2)
	}

	// Check if modelTag is specified in first argument via ':' symbol
	if strings.Contains(modelName, ":") {
		parts := strings.SplitN(modelName, ":", 2)
		modelName = parts[0]
		modelTag = parts[1]
	}

	// If no modelTag, default to "latest"
	if modelTag == "" {
		modelTag = "latest"
	}

	return modelName, modelTag
}

func Help() string {
	return strings.TrimSpace(`
ollama::find
  A CLI tool that allows you to quickly generate a path to a gguf file that's
  been pulled via Ollama.

Usage:
  $  ollama-find llama3
  ~/.ollama/models/blobs/sha256-6a0746a1ec1aef3e7ec53868f220ff6e389f6f8ef87a01d77c96807de94ca2aa
	`)
}
