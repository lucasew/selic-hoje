# Repository Conventions and Guidelines

## Project Structure & Operational Memory
- `index.html` -> Frontend entrypoint (Vanilla JS, single page).
- `api/selichoje.go` -> Backend API handler (Go serverless function for Vercel).
- `api/selichoje_test.go` -> Unit tests for the API handler.
- `shell.nix` -> Development environment dependencies.

## Directives
- Go 1.22+ is used. Avoid deprecated packages like `io/ioutil`.
- HTTP handlers should use standard library `http.Error` for error responses instead of manually writing headers and response bodies, and raw system errors must never be exposed directly to clients.
- To prevent DoS and memory exhaustion, HTTP clients must explicitly define a `Timeout`, and response bodies must be read with a size limit (e.g., `io.LimitReader`).
- The `api/selichoje.go` implementation contains a known logic limitation where `requestData` sends a hardcoded date to the BCB API.
- All code paths that handle unexpected errors MUST funnel through a single, centralized error-reporting function.
- Do not extract helper functions into external directories like `src/`. Merge them directly within the API handler files to prevent Vercel module resolution failures.
- No ghost comments (commented-out code).
- Documentation should focus on explaining the 'why', non-obvious nuances, flow, side effects, and edge cases. Avoid obvious or redundant comments. Use standard `//` docstrings for both exported and internal functions.
- The `.gitignore` file must explicitly exclude `.vercel`, `mise`, and `mise.local.toml`.
