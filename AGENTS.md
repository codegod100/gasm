# Agent Memory

## Project Notes
- This is a TinyGo WASM chat application
- Do NOT auto-build the project - the build process works correctly but may appear to hang
- Use `./build.sh` only when explicitly requested by the user
- The syscall/js import errors in the editor are expected - TinyGo handles this during compilation
- localStorage functionality has been implemented to persist chat messages
- All JavaScript logic has been moved to WASM Go code for minimal JS footprint

## Build Process
- Build script: `./build.sh`
- Serves on: `python3 -m http.server 8000`
- Access at: http://localhost:8000

## Architecture
- Pure WASM implementation with minimal JavaScript
- localStorage persistence for chat history
- Event handlers managed entirely in Go/WASM