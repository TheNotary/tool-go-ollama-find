# OllamaFind
![Tests](https://github.com/thenotary/tool-go-ollama-find/actions/workflows/build.yml/badge.svg)

This is a CLI tool that allows you to quickly generate a path to a gguf file that's been pulled via Ollama.

## Install From Github

    go install github.com/thenotary/tool-go-ollama-find/cmd/ollama-find@v0.1.3

## Build/ Run from Source

    go run cmd/ollama-find/main.go
    go build -o ollama-find cmd/ollama-find/main.go
    ./ollama-find

## Run Tests

    alias gotest="set -o pipefail && go test ./... -json | tparse -all"
    gotest

