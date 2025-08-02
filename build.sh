#!/bin/bash

echo "Building TinyGo WASM project..."

# Build the WASM file
tinygo build -o main.wasm -target wasm main.go

# Copy the wasm_exec.js file from TinyGo
cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js .

echo "Build complete! Open index.html in a web browser to test."
echo "Note: You may need to serve the files over HTTP due to CORS restrictions."
echo "You can use: python3 -m http.server 8000"