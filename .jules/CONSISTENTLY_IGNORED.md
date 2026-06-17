## IGNORE: Out-of-Scope Dependency Updates

**- Pattern:** Restrictive personas (e.g., Docs, Janitor) modifying `go.mod` or `go.sum` files.
**- Justification:** Updating or downgrading dependencies like `go` version when instructed to only fix formatting or documentation leads to CI and code review failures, as it changes the executable logic or dependency graph beyond the persona's scope.
**- Files Affected:** `go.mod`, `go.sum`

## IGNORE: Out-of-Scope Formatting and Cleanups

**- Pattern:** Non-Janitor personas (e.g., Docs, Sentinel, Refactor) performing formatting fixes, removing ghost comments, or replacing deprecated functions like `ioutil.ReadAll`.
**- Justification:** These personas have specific operational contracts (e.g., 'Docs' strictly modifies documentation). Bundling general codebase cleanups with specific fixes violates the scope of their contracts, causing PR rejections.
**- Files Affected:** `api/selichoje.go`, `api/selichoje_test.go`

## IGNORE: Janitor Extracting Abstractions and Changing Behavior

**- Pattern:** Janitor persona extracting new abstractions like central `ReportError` functions, modifying error responses (e.g., hiding raw errors), or replacing HTTP requests with mock payloads in tests.
**- Justification:** The Janitor persona is strictly responsible for fixing lints, formatting, and small quality improvements (< 50 lines) without changing behavior or extracting new abstractions (which belong to the Refactor persona).
**- Files Affected:** `api/selichoje.go`, `api/selichoje_test.go`

## IGNORE: Invalid CI Workflows for Vercel Projects

**- Pattern:** The Arrumador persona adding Release and Artifacts steps to `.github/workflows/autorelease.yml`.
**- Justification:** The project is a Vercel/Cloudflare Workers project. The Release and Artifact steps must be skipped for these projects according to project guidelines, as they lead to unnecessary or failing CI steps.
**- Files Affected:** `.github/workflows/autorelease.yml`

## IGNORE: Restrictive Personas Modifying Tooling Configs

**- Pattern:** Non-Arrumador personas (e.g., Docs, Sentinel, Janitor, Refactor) creating or modifying configuration files like `mise.toml`, `shell.nix`, `.gitignore`, or `.github/workflows/autorelease.yml` to fix lints or add tasks.
**- Justification:** Modifying CI pipelines or tooling configurations is the strict domain of the Arrumador persona. When other personas attempt to 'fix' the environment to pass tests instead of fixing the code, the PR gets rejected.
**- Files Affected:** `mise.toml`, `shell.nix`, `.gitignore`, `.github/workflows/autorelease.yml`
