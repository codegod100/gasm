#!/bin/bash

echo "Building Go WASM project with Templ and Tailwind..."

# Install npm dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "Installing npm dependencies..."
    npm install
fi

# Generate Tailwind CSS
echo "Generating Tailwind CSS..."
npx tailwindcss build -i ./input.css -o ./output.css

# Generate templ templates from separate files
templ generate page_templates.go
templ generate component_templates.go  
templ generate styles_templates.go

# Generate HTML from templates
go run generate_html.go *_templ.go

# Build the WASM file using standard Go
env GOOS=js GOARCH=wasm go build -o main.wasm main.go *_templ.go

# Copy the wasm_exec.js file from Go (detect correct path for different Go versions)
if [ -f "$(go env GOROOT)/misc/wasm/wasm_exec.js" ]; then
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
elif [ -f "$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .
else
    echo "Could not find wasm_exec.js in Go installation"
    find "$(go env GOROOT)" -name "wasm_exec.js" -type f
fi

echo "Build complete! Open index.html in a web browser to test."
echo "Note: You may need to serve the files over HTTP due to CORS restrictions."
echo "You can use: python3 -m http.server 8000"