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

// HandleHelpCommand peaks at the args and if appropriate will print the help
// message for the CLI tool.
// Returns true to indicate the CLI operation was handled or false if not.
func HandleHelpCommand(args []string) bool {
	helpCommands := map[string]bool{
		"help":   true,
		"halp":   true,
		"--help": true,
		"-h":     true,
	}

	if len(args) == 0 || helpCommands[args[0]] {
		fmt.Println(Help())
		return true
	}
	return false
}

// HandleFindCommand peaks at the args and if appropriate will perform the
// FindCommand CLI operation.
// Returns true to indicate the CLI operation was handled or false if not.
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

// ParseModelNameAndTag peaks at the args and returns the modelName and
// modelTag that the user requested to find.
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
