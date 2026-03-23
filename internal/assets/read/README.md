# assets/read/

Domain accessor packages for the embedded `assets.FS`.

Each sub-package owns a slice of the embedded filesystem and exposes
it through clean, domain-named functions (e.g., `entry.List()`,
`claude.Md()`). This avoids a god-object `assets` package where every
function is prefixed with its domain (`assets.EntryList`,
`assets.ClaudeMd`).

## Why this directory exists

- `assets.FS` must live in `internal/assets/` (Go `//go:embed` constraint).
- Domain packages here import `assets.FS` - the dependency flows one way.
- `assets` must **never** import packages under `read/` (that creates a cycle).
- Tests in `package assets` cannot import `read/` packages for the same reason;
  they use `assets.FS` directly.

## Adding a new domain package

1. Create `read/<domain>/<domain>.go`
2. Import `assets.FS` and `config/asset` for path constants
3. Export functions named from the domain's perspective: `List()`, `ForName()`,
   `Content()` - not `DomainList()`, `DomainContent()`
4. Callers read as `entry.List()`, `claude.Md()`, `agent.CopilotInstructions()`

## Do not

- Flatten these packages into `assets` - that's the opposite direction
- Delete this directory - the structure is intentional
- Import `read/` packages from `internal/assets/*.go` - cycle
