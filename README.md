# Go WASM Chat

A basic chat application built with Go that compiles to WebAssembly, featuring Chart.js integration for message statistics.

## Prerequisites

- [Go](https://golang.org/dl/) 1.16+ installed
- [Node.js](https://nodejs.org/) and npm for Chart.js dependency
- A web browser
- A local web server (for CORS compliance)

## Building

First install npm dependencies:

```bash
npm install
```

Then run the build script:

```bash
./build.sh
```

This will:
1. Compile the Go code to WebAssembly using standard Go
2. Copy the required `wasm_exec.js` file from Go's standard library

## Running

Since browsers enforce CORS policies, you need to serve the files over HTTP:

```bash
python3 -m http.server 8000
```

Then open http://localhost:8000 in your browser.

## Features

- Send messages with username and timestamps
- View message history with localStorage persistence
- Clear all messages
- Interactive message statistics with Chart.js doughnut chart
- Toggle statistics view
- Responsive chat interface
- Real-time message display powered by Go WASM
- Full Chart.js integration demonstrating JS library wrapping

## GitHub Pages Deployment

This project includes a GitHub Action that automatically builds and deploys to GitHub Pages using the `acifani/setup-tinygo` action.

### Setup Instructions:

1. Push this repository to GitHub
2. Go to your repository Settings → Pages
3. Under "Source", select "GitHub Actions"
4. The workflow will automatically trigger on pushes to `main` or `master` branch
5. Your chat app will be available at `https://yourusername.github.io/yourrepo`

### Manual Deployment:

You can also trigger the deployment manually:
- Go to Actions tab in your GitHub repository
- Select "Deploy Go WASM Chat to GitHub Pages"
- Click "Run workflow"

## Files

- `main.go` - Go source code with chat logic and Chart.js integration
- `index.html` - HTML page with chat UI and Chart.js
- `build.sh` - Build script for Go WASM
- `package.json` - npm configuration with Chart.js dependency
- `main.wasm` - Compiled WebAssembly (generated)
- `wasm_exec.js` - Go WASM runtime (copied during build)
- `.github/workflows/deploy.yml` - GitHub Actions workflow for automatic deployment