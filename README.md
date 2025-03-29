# Ollama::Find

This is a CLI tool that allows you to quickly generate a path to a gguf file that's been pulled via Ollama.


## Install Dependencies and Run

Requirements are tracked in requirements.txt by convention.

```bash
go mod init github.com/TheNotary/ollama-find
go mod tidy
go get ./...
go build cmd/ollama-find.go
go run cmd/ollama-find.go
```


## Run Tests

```
go test ./...
```
