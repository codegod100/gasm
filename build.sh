#!/bin/bash

# Exit immediately if any command fails
set -e

# Function to handle errors
handle_error() {
    echo "❌ Build failed at line $1"
    exit 1
}

# Set up error handling
trap 'handle_error $LINENO' ERR

echo "Building Go WASM project with Templ and Tailwind..."

# Install npm dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "Installing npm dependencies..."
    npm install || {
        echo "❌ Failed to install npm dependencies"
        exit 1
    }
fi

# Generate Tailwind CSS
echo "Generating Tailwind CSS..."
npx tailwindcss build -i ./input.css -o ./output.css || {
    echo "❌ Failed to generate Tailwind CSS"
    exit 1
}

# Generate templ templates - FAIL HARD if this doesn't work
echo "Generating templ templates..."

# Check if templ command exists
if ! command -v templ >/dev/null 2>&1; then
    echo "❌ templ command not found. Install with: go install github.com/a-h/templ/cmd/templ@latest"
    exit 1
fi

# Generate templ templates from .templ files
echo "Generating page templates..."
templ generate -f page_templates.templ || {
    echo "❌ Failed to generate page templates"
    echo "Check page_templates.templ for syntax errors"
    exit 1
}

echo "Generating component templates..."
templ generate -f component_templates.templ || {
    echo "❌ Failed to generate component templates"
    echo "Check component_templates.templ for syntax errors"
    exit 1
}

echo "Generating styles templates..."
templ generate -f styles_templates.templ || {
    echo "❌ Failed to generate styles templates"
    echo "Check styles_templates.templ for syntax errors"
    exit 1
}

# Verify generated files exist
if ! ls *_templ.go >/dev/null 2>&1; then
    echo "❌ No *_templ.go files were generated"
    exit 1
fi

# Generate HTML from templates
echo "Generating HTML from templates..."
go run generate_html.go *_templ.go || {
    echo "❌ Failed to generate HTML from templates"
    echo "Check generate_html.go and template files for errors"
    exit 1
}

# Build the WASM file with templates
echo "Building WASM with templates..."
env GOOS=js GOARCH=wasm go build -o main.wasm main.go *_templ.go || {
    echo "❌ Failed to build WASM with templates"
    echo "Check main.go and template files for compilation errors"
    exit 1
}

# Copy the wasm_exec.js file from Go (detect correct path for different Go versions)
echo "Copying wasm_exec.js..."
if [ -f "$(go env GOROOT)/misc/wasm/wasm_exec.js" ]; then
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" . || {
        echo "❌ Failed to copy wasm_exec.js"
        exit 1
    }
elif [ -f "$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" . || {
        echo "❌ Failed to copy wasm_exec.js"
        exit 1
    }
else
    echo "❌ Could not find wasm_exec.js in Go installation"
    echo "Searching for wasm_exec.js..."
    find "$(go env GOROOT)" -name "wasm_exec.js" -type f || true
    exit 1
fi

echo "✅ Build complete! Open index.html in a web browser to test."
echo "Note: You may need to serve the files over HTTP due to CORS restrictions."
echo "You can use: python3 -m http.server 8000"