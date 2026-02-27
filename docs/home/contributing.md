---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Contributing
icon: lucide/git-pull-request
---

![ctx](../images/ctx-banner.png)

## Development Setup

### Prerequisites

* [Go](https://go.dev/) (*version defined in [`go.mod`](https://github.com/ActiveMemory/ctx/blob/main/go.mod)*)
* [Claude Code](https://docs.anthropic.com/en/docs/claude-code/overview)
* [Git](https://git-scm.com/)
* [GNU Make](https://www.gnu.org/software/make/)
* [Zensical](https://github.com/zensical/zensical)

### 1. Fork (*or Clone*) the Repository

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

The repository ships a Claude Code plugin under `internal/assets/claude/`.
Point Claude Code at your local copy so that skills and hooks reflect
your working tree — no reinstall needed after edits:

1. Launch `claude`;
2. Type `/plugin` and press Enter;
3. Select **Marketplaces** → **Add Marketplace**
4. Enter the **absolute path** to the root of your clone,
   e.g. `~/WORKSPACE/ctx`
   (*this is where `.claude-plugin/marketplace.json` lives: it points
   Claude Code to the actual plugin in `internal/assets/claude`*);
5. Back in `/plugin`, select **Install** and choose `ctx`.

!!! warning "Claude Code Caches Plugin Files"
    Even though the marketplace points at a directory on disk, Claude Code
    **caches** skills and hooks. After editing files under
    `internal/assets/claude/`, clear the cache and restart:

    ```bash
    make plugin-reload   # then restart Claude Code
    ```

    See [Skill or Hook Changes](#skill-or-hook-changes) for details.

### 4. Verify

```bash
ctx --version       # binary is in PATH
claude /plugin list # plugin is installed
```

You should see the `ctx` plugin listed, sourced from your local path.

----

## Project Layout

<!-- drift-check: ls -d cmd/ internal/*/ .claude/ docs/ editors/ hack/ specs/ assets/ examples/ .context/ -->
```
ctx/
├── cmd/ctx/            # CLI entry point
├── internal/
│   ├── assets/claude/  # ← Claude Code plugin (skills, hooks)
│   ├── bootstrap/      # Project initialization templates
│   ├── claude/         # Claude Code integration helpers
│   ├── cli/            # Command implementations
│   ├── config/         # Configuration loading
│   ├── context/        # Core context logic
│   ├── crypto/         # Scratchpad encryption
│   ├── drift/          # Drift detection
│   ├── index/          # Context file indexing
│   ├── journal/        # Journal site generation
│   ├── notify/         # Webhook notifications
│   ├── rc/             # .ctxrc parsing
│   ├── recall/         # Session history and parsers
│   ├── sysinfo/        # System resource monitoring
│   ├── task/           # Task management
│   └── validation/     # Input validation
├── .claude/
│   └── skills/         # Dev-only skills (not distributed)
├── assets/             # Static assets (banners, logos)
├── docs/               # Documentation site source
├── editors/            # Editor extensions (VS Code)
├── examples/           # Example configurations
├── hack/               # Build scripts and runbooks
├── specs/              # Feature specifications
└── .context/           # ctx's own context (dogfooding)
```

### Skills: Two Directories, One Rule

<!-- drift-check: ls internal/assets/claude/skills/ | wc -l -->

| Directory                        | What lives here                                 | Distributed to users? |
|----------------------------------|-------------------------------------------------|-----------------------|
| `internal/assets/claude/skills/` | The 29 `ctx-*` skills that ship with the plugin | Yes                   |
| `.claude/skills/`                | Dev-only skills (release, QA, backup, etc.)     | No                    |

**`internal/assets/claude/skills/`** is the single source of truth for
user-facing skills. If you are adding or modifying a `ctx-*` skill,
edit it there.

**`.claude/skills/`** holds skills that only make sense inside this
repository (*release automation, QA checks, backup scripts*). These are
never distributed to users.

#### Dev-Only Skills Reference

<!-- drift-check: ls .claude/skills/ -->

| Skill | When to use |
|-------|-------------|
| `/_ctx-absorb` | Merge deltas from a parallel worktree or separate checkout |
| `/_ctx-audit` | Detect code-level drift after YOLO sprints or before releases |
| `/_ctx-backup` | Backup context and Claude data to SMB share |
| `/_ctx-brainstorm` | Structured design dialogue before implementation |
| `/_ctx-check-links` | Audit docs for dead links before releases |
| `/_ctx-qa` | Run QA checks before committing |
| `/_ctx-release` | Run the full release process |
| `/_ctx-release-notes` | Generate release notes for `dist/RELEASE_NOTES.md` |
| `/_ctx-sanitize-permissions` | Audit `settings.local.json` for dangerous permissions |
| `/_ctx-skill-creator` | Create or improve a skill |
| `/_ctx-spec` | Scaffold a feature spec from `specs/spec-template.md` |
| `/_ctx-update-docs` | Check docs/code consistency after changes |
| `/_ctx-verify` | Verify before claiming work is done |

----

## Day-to-Day Workflow

### Go Code Changes

After modifying Go source files, rebuild and reinstall:

```bash
make build && sudo make install
```

The `ctx` binary is statically compiled. There is no hot reload.
You must rebuild for Go changes to take effect.

### Skill or Hook Changes

Edit files under `internal/assets/claude/skills/` or
`internal/assets/claude/hooks/`.

Claude Code caches plugin files, so edits aren't picked up automatically.
Clear the cache and restart:

```bash
make plugin-reload   # nukes ~/.claude/plugins/cache/activememory-ctx/
# then restart Claude Code
```

The plugin will be re-installed from your local marketplace on startup.
No version bump is needed during development.

!!! tip "Version bumps are for releases, not iteration"
    Only bump `VERSION`, `plugin.json`, and `marketplace.json` when
    cutting a release. During development, `make plugin-reload` is
    all you need.

### Running Tests

```bash
make test   # fast: all tests
make audit  # full: fmt + vet + lint + drift + docs + test
make smoke  # build + run basic commands end-to-end
```

### Running the Docs Site Locally

```bash
make site-setup  # one-time: install zensical via pipx
make site-serve  # serve at localhost
```

----

## Submitting Changes

### Before You Start

1. Check existing issues to avoid duplicating effort;
2. For large changes, open an issue first to discuss the approach;
3. Read the specs in `specs/` for design context.

### Pull Request Process

Respect the maintainers' time and energy:
Keep your pull requests **isolated** and strive to minimze code changes.

If you Pull Request solves more than one distinct issues, it's better to create
separate pull requests instead of sending them in one large bundle.

1. Create a feature branch: `git checkout -b feature/my-feature`;
2. Make your changes;
3. Run `make audit` to catch issues early;
4. Commit with a **clear message**;
5. Push and open a pull request.

!!! tip "Audit Your Code Before Submitting"
    Run `make audit` before submitting:

    `make audit` covers formatting, vetting, linting, drift checks, 
    doc consistency, and tests in one pass.

### Commit Messages

Following conventional commits is recommended but not required:

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`

Examples:

* `feat(cli): add ctx export command`
* `fix(drift): handle missing files gracefully`
* `docs: update installation instructions`

### Code Style

* Follow Go conventions (`gofmt`, `go vet`);
* Keep functions **focused** and **small**;
* Add tests for new functionality;
* Handle errors explicitly.

----

## Code of Conduct

A clear context requires **respectful** collaboration.

`ctx` follows the
[Contributor Covenant](https://github.com/ActiveMemory/ctx/blob/main/CODE_OF_CONDUCT.md).

----

## Boring Legal Stuff

### Developer Certificate of Origin (*DCO*)

By contributing, you agree to the
[Developer Certificate of Origin](https://github.com/ActiveMemory/ctx/blob/main/CONTRIBUTING_DCO.md).

All commits must be signed off:

```bash
git commit -s -m "feat: add new feature"
```

### License

Contributions are licensed under the
[Apache 2.0 License](https://github.com/ActiveMemory/ctx/blob/main/LICENSE).
