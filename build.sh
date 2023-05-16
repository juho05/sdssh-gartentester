#!/bin/bash

OUT_DIR=./bin

export CGO_ENABLED=0
echo "Compiling for Windows x64..." &&
GOOS=windows GOARCH=amd64 go build -o "$OUT_DIR/gartentester-windows-x64.exe" &&
echo "Compiling for Windows ARM64..." &&
GOOS=windows GOARCH=arm64 go build -o "$OUT_DIR/gartentester-windows-arm64.exe" &&
echo "Compiling for macOS x64..." &&
GOOS=darwin GOARCH=amd64 go build -o "$OUT_DIR/gartentester-macos-x64" &&
echo "Compiling for macOS ARM64..." &&
GOOS=darwin GOARCH=arm64 go build -o "$OUT_DIR/gartentester-macos-arm64" &&
echo "Compiling for Linux x64..." &&
GOOS=linux GOARCH=amd64 go build -o "$OUT_DIR/gartentester-linux-x64" &&
echo "Compiling for Linux ARM64..." &&
GOOS=linux GOARCH=arm64 go build -o "$OUT_DIR/gartentester-linux-arm64" &&

echo "Done. The output files can be found in $OUT_DIR/"
