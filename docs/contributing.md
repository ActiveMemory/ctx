---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Contributing to ctx
icon: lucide/git-pull-request
---

## Development Setup

### Prerequisites

- [Go 1.25+](https://go.dev/)
- [Claude Code](https://docs.anthropic.com/en/docs/claude-code/overview)
- Git

### 1. Fork (or Clone) the Repository

```bash
# Fork on GitHub, then:
git clone https://github.com/<you>/ctx.git
cd ctx

# Or, if you have push access:
git clone https://github.com/ActiveMemory/ctx.git
cd ctx
```

### 2. Build and Install the Binary

```bash
make build
sudo make install
```

This compiles the `ctx` binary and places it in `/usr/local/bin/`.

### 3. Install the Plugin from Your Local Clone

The repository ships a Claude Code plugin under `internal/tpl/claude/`.
Point Claude Code at your local copy so that skills and hooks reflect
your working tree — no reinstall needed after edits:

1. Launch `claude`
2. Type `/plugin` and press Enter
3. Select **Marketplaces** → **Add Marketplace**
4. Enter the **absolute path** to the plugin directory in your clone,
   e.g. `~/WORKSPACE/ctx/internal/tpl/claude`
5. Back in `/plugin`, select **Install** and choose `ctx`

!!! tip "Local Plugin = Live Edits"
    Because the marketplace points at a directory on disk, any change
    you make to a skill or hook under `internal/tpl/claude/` takes
    effect the next time Claude Code loads it. No rebuild, no
    reinstall.

### 4. Verify

```bash
ctx --version       # binary is in PATH
claude /plugin list # plugin is installed
```

You should see the `ctx` plugin listed, sourced from your local path.

----

## Project Layout

```
ctx/
├── cmd/ctx/                  # CLI entry point
├── internal/
│   ├── cli/                  # Command implementations
│   ├── context/              # Core context logic
│   ├── drift/                # Drift detection
│   ├── claude/               # Claude Code integration helpers
│   └── tpl/                  # Embedded templates and plugin
│       └── claude/           # ← Claude Code plugin (skills, hooks)
│           └── skills/       #   Source of truth for distributed skills
├── .claude/
│   └── skills/               # Dev-only skills (not distributed)
├── specs/                    # Feature specifications
├── docs/                     # Documentation site source
└── .context/                 # ctx's own context (dogfooding)
```

### Skills: Two Directories, One Rule

| Directory | What lives here | Distributed to users? |
|---|---|---|
| `internal/tpl/claude/skills/` | The 25 `ctx-*` skills that ship with the plugin | Yes |
| `.claude/skills/` | Dev-only skills (release, QA, backup, etc.) | No |

**`internal/tpl/claude/skills/`** is the single source of truth for
user-facing skills. If you are adding or modifying a `ctx-*` skill,
edit it there.

**`.claude/skills/`** holds skills that only make sense inside this
repository (release automation, QA checks, backup scripts). These are
never distributed to users.

----

## Day-to-Day Workflow

### Go Code Changes

After modifying Go source files, rebuild and reinstall:

```bash
make build && sudo make install
```

The `ctx` binary is statically compiled — there is no hot reload.
You must rebuild for Go changes to take effect.

### Skill or Hook Changes

Edit files under `internal/tpl/claude/skills/` or
`internal/tpl/claude/hooks/`. Changes are picked up by Claude Code
automatically (no rebuild needed).

### Running Tests

```bash
make test            # fast: all tests
make audit           # full: fmt + vet + lint + drift + docs + test
make smoke           # build + run basic commands end-to-end
```

### Running the Docs Site Locally

```bash
make site-setup      # one-time: create venv, install zensical
make site-serve      # serve at localhost
```

----

## Submitting Changes

### Before You Start

1. Check existing issues to avoid duplicating effort
2. For large changes, open an issue first to discuss the approach
3. Read the specs in `specs/` for design context

### Pull Request Process

1. Create a feature branch: `git checkout -b feature/my-feature`
2. Make your changes
3. Run `make audit` to catch issues early
4. Commit with a clear message (see below)
5. Push and open a pull request

### Commit Messages

Follow conventional commits:

```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`

Examples:

- `feat(cli): add ctx export command`
- `fix(drift): handle missing files gracefully`
- `docs: update installation instructions`

### Code Style

- Follow Go conventions (`gofmt`, `go vet`)
- Keep functions focused and small
- Add tests for new functionality
- Handle errors explicitly

----

## Legal

### Developer Certificate of Origin (DCO)

By contributing, you agree to the
[Developer Certificate of Origin](https://github.com/ActiveMemory/ctx/blob/main/CONTRIBUTING_DCO.md).

All commits must be signed off:

```bash
git commit -s -m "feat: add new feature"
```

### License

Contributions are licensed under the
[Apache 2.0 License](https://github.com/ActiveMemory/ctx/blob/main/LICENSE).

### Code of Conduct

This project follows the
[Contributor Covenant Code of Conduct](https://github.com/ActiveMemory/ctx/blob/main/CODE_OF_CONDUCT.md).
