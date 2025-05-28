#!/bin/bash
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/markdown-refactor-linux-amd64 markdown_refactor.go

echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/markdown-refactor-windows-amd64.exe markdown_refactor.go

echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/markdown-refactor-macos-amd64 markdown_refactor.go

echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/markdown-refactor-macos-arm64 markdown_refactor.go

echo "Build complete. Binaries are in the 'dist' folder."