#!/bin/bash

echo "Building Go WASM project..."

# Build the WASM file using standard Go
GOOS=js GOARCH=wasm go build -o main.wasm main.go

# Copy the wasm_exec.js file from Go
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

echo "Build complete! Open index.html in a web browser to test."
echo "Note: You may need to serve the files over HTTP due to CORS restrictions."
echo "You can use: python3 -m http.server 8000"