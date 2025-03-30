# OllamaFind
![Tests](https://github.com/thenotary/tool-go-ollama-find/actions/workflows/build.yml/badge.svg)

This is a CLI tool that allows you to quickly generate a path to a gguf file that's been pulled via Ollama.


## Install From Github

    go install github.com/thenotary/tool-go-ollama-find@v0.1.0


## Install Dependencies and Run

    go mod tidy
    go build cmd/ollama-find.go
    go run cmd/ollama-find.go


## Run Tests

    go test ./...
