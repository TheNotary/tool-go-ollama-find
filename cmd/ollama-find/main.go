package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/thenotary/tool-go-ollama-find/pkg/api"
)

func main() {
	args := os.Args[1:]
	var handled bool

	handled = HandleHelpCommand(args)
	if handled {
		return
	}

	handled = HandleFindCommand(args)
	if handled {
		return
	}

	fmt.Println("error: unable to handle operation.")
}

func HandleHelpCommand(args []string) bool {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Println(Help())
		return true
	}
	return false
}

func HandleFindCommand(args []string) bool {
	var arg1, arg2 string
	if len(args) > 0 {
		arg1 = args[0]
	}
	if len(args) > 1 {
		arg2 = args[1]
	}

	modelName, modelTag := ParseModelNameAndTag(arg1, arg2)

	gguf_path, err := api.LookupGGUF(modelName, modelTag)
	if err != nil {
		fmt.Println("error: something went wrong calling LookupGGUF", err)
		return false
	}

	fmt.Println(gguf_path)
	return true
}

func ParseModelNameAndTag(arg1 string, arg2 string) (string, string) {
	modelName := strings.TrimSpace(arg1)
	modelTag := ""

	if len(arg2) > 1 {
		modelTag = strings.TrimSpace(arg2)
	}

	// Handle case where modelTag is specified in first argument via ':' symbol
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
