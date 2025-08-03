# Go WASM Chat with Templ

A modern chat application built with Go, compiled to WebAssembly, using Templ for templating and Tailwind CSS for styling.

## Features

- Real-time chat interface
- Persistent message storage using localStorage
- Message statistics with Chart.js visualization
- Modern UI with Catppuccin color scheme
- Fully built in Go/WASM with minimal JavaScript
- Refactored to use templ components for better maintainability

## Prerequisites

- [Go](https://golang.org/dl/) 1.20+ installed
- [Node.js](https://nodejs.org/) and npm for Tailwind CSS and Chart.js
- [Templ](https://templ.guide/) for template generation
- A web browser
- A local web server (for CORS compliance)

## Running Locally

1. Install dependencies:
   ```bash
   npm install
   ```

2. Build the project:
   ```bash
   ./build.sh
   ```

3. Serve the files:
   ```bash
   python3 -m http.server 8000
   ```

4. Open http://localhost:8000 in your browser

## Architecture

- **Go WASM**: Core application logic
- **Templ**: Template system for HTML generation  
- **Tailwind CSS**: Utility-first styling with Catppuccin theme
- **Chart.js**: Statistics visualization
- **localStorage**: Client-side data persistence

## GitHub Pages Deployment

This project includes a GitHub Action that automatically builds and deploys to GitHub Pages.

### Setup Instructions:

1. Push this repository to GitHub
2. Go to your repository Settings → Pages
3. Under "Source", select "GitHub Actions"
4. The workflow will automatically trigger on pushes to `main` branch
5. Your chat app will be available at `https://yourusername.github.io/yourrepo`

## Files

- `main.go` - Go source code with chat logic
- `*_templates.templ` - Templ template files
- `input.css` - Tailwind CSS source
- `build.sh` - Build script for Go WASM with Templ and Tailwind
- `package.json` - npm configuration with dependencies
- `tailwind.config.js` - Tailwind configuration with Catppuccin colors
- Generated files: `main.wasm`, `wasm_exec.js`, `index.html`, `output.css`

## Live Demo

Visit: https://vamshiaruru.github.io/gasm/