package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/thenotary/tool-go-ollama-find/pkg/ollama_find"
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
	modelName, modelTag := ParseModelNameAndTag(args)

	gguf_path, err := ollama_find.LookupGGUF(modelName, modelTag)
	if err != nil {
		fmt.Println("error: something went wrong calling LookupGGUF", err)
		return false
	}

	fmt.Println(gguf_path)
	return true
}

func ParseModelNameAndTag(args []string) (modelName, modelTag string) {
	if len(args) > 0 {
		modelName = strings.TrimSpace(args[0])
	}

	// Handle case where modelTag is specified in first argument via ':' symbol
	if strings.Contains(modelName, ":") {
		parts := strings.SplitN(modelName, ":", 2)
		modelName = parts[0]
		modelTag = parts[1]
		return
	}

	// Handle case where modelTag is specified in second argument
	if len(args) > 1 && strings.TrimSpace(args[1]) != "" {
		modelTag = strings.TrimSpace(args[1])
		return
	}

	// Handle case where modelTag is ommited, defaulting to "latest"
	modelTag = "latest"
	return
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
