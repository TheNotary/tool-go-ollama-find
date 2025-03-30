# OllamaFind
![Tests](https://github.com/thenotary/tool-go-ollama-find/actions/workflows/build.yml/badge.svg)

This is a CLI tool that allows you to quickly generate a path to a gguf file that's been pulled via Ollama.


## Install From Github

    go install github.com/thenotary/tool-go-ollama-find@v0.1.1
    go install github.com/thenotary/tool-go-ollama-find/cmd/ollama-find@v0.1.2


## Build/ Run from Source

    go run main.go
    go build main.go
    ./ollama-find


## Run Tests

    go test ./...


# Ecosystem Notes
## Delete These for Public Projects

## Distribution Troublshooting
Versioned releases are cut exclusively by their git tag in Go eliminated the need to create "bumps version to x.y.z" commits!

Go releases are cached by "the go module proxy" which strictly adheres to immutable releases.  Even the `@latest` tag can render the wrong package when doing pre-release git repo cleanup.

Also, local caches will exist after installing a package, so you may want to clean those up with:

    go clean -cache -modcache
    rm "${GOPATH}/bin/tool-go-ollama-find"
