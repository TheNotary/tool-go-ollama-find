# OllamaFind
[![Go Reference](https://pkg.go.dev/badge/github.com/thenotary/tool-go-ollama-find.svg)](https://pkg.go.dev/github.com/thenotary/tool-go-ollama-find)
[![Tests](https://github.com/thenotary/tool-go-ollama-find/actions/workflows/build.yml/badge.svg)](https://github.com/thenotary/tool-go-ollama-find/actions/workflows/build.yml)

This is a CLI tool that allows you to quickly generate a path to a gguf file that's been pulled via Ollama.  Also see [tool-ruby-ollama-find](https://github.com/TheNotary/tool-ruby-ollama-find) for the Ruby implementation.

## Install From Github

    go install github.com/thenotary/tool-go-ollama-find/cmd/ollama-find@latest

## Build/ Run from Source

    go run cmd/ollama-find/main.go
    go build -o ollama-find cmd/ollama-find/main.go
    ./ollama-find

## Run Tests

    go install github.com/mfridman/tparse@latest
    alias gotest="set -o pipefail && go test ./... -json | tparse -all"
    gotest

