#!/bin/bash

# Build script for FlashAppleEDGE
# Builds for Linux x86_64, Linux aarch64, and Windows x86_64

VERSION=0.0.5

echo "Building FlashAppleEDGE..."

rm -rf release/*


# Linux x86_64
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'main.Version=${VERSION}'" -o release/flashAppleEDGE.x86_64 cmd/simple/main.go
echo "Built: flashAppleEDGE.x86_64"

# Linux aarch64
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X 'main.Version=${VERSION}'" -o release/flashAppleEDGE.aarch64 cmd/simple/main.go
echo "Built: flashAppleEDGE.aarch64"

# Windows x86_64
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X 'main.Version=${VERSION}'" -o release/flashAppleEDGE.exe cmd/simple/main.go
echo "Built: flashAppleEDGE.exe"

echo "Build complete!"

touch release/v${VERSION}.info
echo "Version: ${VERSION}" > release/v${VERSION}.info
echo "Erstellt am $(date)" >> release/v${VERSION}.info
