# TinyGo WASM Chat

A basic chat application built with TinyGo that compiles to WebAssembly.

## Prerequisites

- [TinyGo](https://tinygo.org/getting-started/install/) installed
- A web browser
- A local web server (for CORS compliance)

## Building

Run the build script:

```bash
./build.sh
```

This will:
1. Compile the Go code to WebAssembly using TinyGo
2. Copy the required `wasm_exec.js` file

## Running

Since browsers enforce CORS policies, you need to serve the files over HTTP:

```bash
python3 -m http.server 8000
```

Then open http://localhost:8000 in your browser.

## Features

- Send messages with username
- View message history with timestamps
- Clear all messages
- Responsive chat interface
- Real-time message display powered by WASM

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
- Select "Deploy TinyGo WASM Chat to GitHub Pages"
- Click "Run workflow"

## Files

- `main.go` - TinyGo source code with chat logic
- `index.html` - HTML page with chat UI and WASM loader
- `build.sh` - Build script
- `main.wasm` - Compiled WebAssembly (generated)
- `wasm_exec.js` - TinyGo WASM runtime (copied during build)
- `.github/workflows/deploy.yml` - GitHub Actions workflow for automatic deployment