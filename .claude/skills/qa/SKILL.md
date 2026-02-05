---
name: qa
description: "Run QA checks before committing. Use after writing code, before commits, or when CI might fail."
---

Run the project's QA pipeline locally to catch issues before they hit CI.

## What to Run

Run these checks **in order** — each depends on the previous passing:

### 1. Format

```bash
gofmt -l .
```

If files are listed, they need formatting. Fix with `gofmt -w .` and
include the formatted files in the commit.

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

Builds the binary and exercises `ctx init`, `ctx status`, `ctx agent`,
`ctx drift`, `ctx add task`, and `ctx session save` in a temp directory.

## Shortcut

`make audit` runs steps 1-4 in sequence. Use it when you want a single
pass/fail answer.

## When to Run What

| Changed              | Minimum Check          |
|----------------------|------------------------|
| Any `.go` file       | `make audit`           |
| CLI command behavior | `make audit` + `make smoke` |
| Only docs/config     | Nothing                |

## Common Failures

| Failure                          | Fix                                             |
|----------------------------------|--------------------------------------------------|
| `gofmt -l` lists files           | `gofmt -w .`                                     |
| `fmt.Printf` in CLI code         | Use `cmd.Printf` (enforced by AST test)          |
| golangci-lint unused variable     | Remove it — don't rename to `_`                  |
| Test needs `CTX_SKIP_PATH_CHECK`  | Already set in `make test` and `make audit`      |
| Coverage below 70% on `internal/context` | Add tests — check with `make test-coverage` |

## Output Format

After running checks, report:

1. **Result**: Pass or fail
2. **Failures**: What failed and how to fix (if any)
3. **Files touched**: List of files that were auto-formatted (if any)
