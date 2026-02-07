---
name: qa
description: "Run QA checks before committing. Use after writing code, before commits, or when CI might fail."
---

Run the project's QA pipeline locally to catch issues before
they hit CI.

## Before Running

1. **Check what changed**: run `git diff --name-only` to see
   which files were modified; this determines which checks
   to run (see "When to Run What" below)
2. **Build first**: `CGO_ENABLED=0 go build -o /dev/null ./cmd/ctx`
   to catch compile errors before running the full pipeline

## When to Use

- After writing or modifying Go code
- Before committing (catch issues locally, not in CI)
- When CI failed and you need to reproduce locally
- After dependency changes or refactors

## When NOT to Use

- When only docs, markdown, or config files changed (no Go
  code touched)
- When the user explicitly says "skip QA" or "commit without
  checks"
- Mid-development when code is intentionally incomplete; wait
  until a logical stopping point

## Usage Examples

```text
/qa
/qa (after refactoring the recall command)
```

## What to Run

Run these checks **in order**; each depends on the previous
passing:

### 1. Format

```bash
gofmt -l .
```

If files are listed, they need formatting. Fix with
`gofmt -w .` and include the formatted files in the commit.

### 2. Vet

```bash
CGO_ENABLED=0 go vet ./...
```

### 3. Lint

```bash
golangci-lint run --timeout=5m
```

### 4. Test

```bash
CGO_ENABLED=0 CTX_SKIP_PATH_CHECK=1 go test ./...
```

### 5. Smoke (if CLI behavior changed)

```bash
make smoke
```

Builds the binary and exercises `ctx init`, `ctx status`,
`ctx agent`, `ctx drift`, `ctx add task`, and
`ctx session save` in a temp directory.

## Shortcut

`make audit` runs steps 1-4 in sequence. Use it when you
want a single pass/fail answer.

## When to Run What

| Changed                | Minimum Check                |
|------------------------|------------------------------|
| Any `.go` file         | `make audit`                 |
| CLI command behavior   | `make audit` + `make smoke`  |
| Only docs/config       | Nothing                      |
| Template files changed | `go build` (embed must work) |

## Common Failures

| Failure                                  | Fix                                         |
|------------------------------------------|---------------------------------------------|
| `gofmt -l` lists files                   | `gofmt -w .`                                |
| `fmt.Printf` in CLI code                 | Use `cmd.Printf` (enforced by AST test)     |
| golangci-lint unused variable            | Remove it; do not rename to `_`             |
| Test needs `CTX_SKIP_PATH_CHECK`         | Already set in `make test` and `make audit` |
| Coverage below 70% on `internal/context` | Add tests; check with `make test-coverage`  |

## Output Format

After running checks, report:

1. **Result**: pass or fail
2. **Failures**: what failed and how to fix (if any)
3. **Files touched**: list of files that were auto-formatted

## Quality Checklist

Before reporting QA results, verify:
- [ ] Ran the appropriate checks for what changed (not more,
      not fewer)
- [ ] Fixed auto-fixable issues (gofmt) rather than just
      reporting them
- [ ] All failures have a clear fix action (not just "it
      failed")
- [ ] If smoke tests were needed (CLI behavior changed),
      they were run
