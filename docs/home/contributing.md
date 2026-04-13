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
your working tree: no reinstall needed after edits:

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
    `internal/assets/claude/`, **clear the cache and restart**:

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
│   ├── memory/         # Memory bridge (discover, mirror, import, publish)
│   ├── notify/         # Webhook notifications
│   ├── rc/             # .ctxrc parsing
│   ├── journal/        # Session history, parsers, and state
│   ├── sysinfo/        # System resource monitoring
│   ├── task/           # Task management
│   └── validation/     # Input validation
├── .claude/
│   └── skills/         # Dev-only skills (not distributed)
├── assets/             # Static assets (banners, logos)
├── docs/               # Documentation site source
├── editors/            # Editor extensions (VS Code)
├── examples/           # Example configurations
├── hack/               # Build scripts and [runbooks](../operations/index.md#runbooks)
├── specs/              # Feature specifications
└── .context/           # ctx's own context (dogfooding)
```

### Skills: Two Directories, One Rule

<!-- drift-check: ls internal/assets/claude/skills/ | wc -l -->

| Directory                        | What lives here                                 | Distributed to users? |
|----------------------------------|-------------------------------------------------|-----------------------|
| `internal/assets/claude/skills/` | The 39 `ctx-*` skills that ship with the plugin | Yes                   |
| `.claude/skills/`                | Dev-only skills (release, QA, backup, etc.)     | No                    |

**`internal/assets/claude/skills/`** is the single source of truth for
user-facing skills. If you are adding or modifying a `ctx-*` skill,
edit it there.

**`.claude/skills/`** holds skills that only make sense inside this
repository (*release automation, QA checks, backup scripts*). These are
never distributed to users.

#### Dev-Only Skills Reference

<!-- drift-check: ls .claude/skills/ -->

| Skill                        | When to use                                                   |
|------------------------------|---------------------------------------------------------------|
| `/_ctx-absorb`               | Merge deltas from a parallel worktree or separate checkout    |
| `/_ctx-audit`                | Detect code-level drift after YOLO sprints or before releases |
| `/_ctx-backup`               | Backup context and Claude data to SMB share                   |
| `/_ctx-qa`                   | Run QA checks before committing                               |
| `/_ctx-release`              | Run the full release process                                  |
| `/_ctx-release-notes`        | Generate release notes for `dist/RELEASE_NOTES.md`            |
| `/_ctx-alignment-audit`      | Audit doc claims against agent instructions                   |
| `/_ctx-update-docs`          | Check docs/code consistency after changes                     |
| `/_ctx-command-audit`        | Audit CLI surface after renames, moves, or deletions          |

Six skills previously in this list have been promoted to bundled plugin skills
and are now available to all ctx users: `/ctx-brainstorm`, `/ctx-link-check`,
`/ctx-permission-sanitize`, `/ctx-skill-create`, `/ctx-spec`.

----

## How To Add Things

### Adding a New CLI Command

1. Create a package under `internal/cli/<name>/` with `doc.go`, `cmd.go`,
   and `run.go`;
2. Implement `Cmd() *cobra.Command` as the entry point;
3. Add `Use*` and `DescKey*` constants in `internal/config/embed/cmd/<name>.go`;
4. Add command descriptions in `internal/assets/commands/commands.yaml`;
5. Add examples in `internal/assets/commands/examples.yaml`;
6. Add flag descriptions in `internal/assets/commands/flags.yaml`;
7. Register the command in `internal/bootstrap/group.go` (add import +
   entry in the appropriate group function);
8. Create an output package at `internal/write/<name>/` for all
   user-facing output (see [Package Taxonomy](#package-taxonomy));
9. Create error constructors at `internal/err/<name>/` for
   domain-specific errors;
10. Add tests in the same package (`<name>_test.go`);
11. Add a doc page at `docs/cli/<name>.md` and update
    `docs/cli/index.md`;
12. Add the page to `zensical.toml` nav.

Pattern to follow: `internal/cli/pad/pad.go` (parent with subcommands) or
`internal/cli/drift/` (single command).

### Package Taxonomy

ctx separates concerns into a strict package taxonomy. Knowing where
things go prevents code review friction and keeps the AST lint tests
happy.

#### Output: `internal/write/`

Every CLI command's user-facing output lives in its own sub-package
under `internal/write/<domain>/`. Output functions accept
`*cobra.Command` and call `cmd.Println(...)` — never `fmt.Print*`
directly. All text strings are loaded from YAML via
`desc.Text(text.DescKey*)`, never inline.

```
internal/write/add/add.go       # output for ctx add
internal/write/stat/stat.go     # output for ctx usage
internal/write/resource/        # output for ctx sysinfo
```

Exception: `write/rc/` writes to `os.Stderr` because rc loads before
cobra is initialized.

#### Errors: `internal/err/`

Domain-specific error constructors live under `internal/err/<domain>/`.
Each package mirrors the write structure. Functions return `error`
(never custom error types) and load messages from YAML via
`desc.Text(text.DescKey*)`.

```
internal/err/add/add.go         # errors for ctx add
internal/err/config/config.go   # errors for configuration
internal/err/cli/cli.go         # errors for CLI argument validation
```

#### Config constants: `internal/config/`

Pure-constant leaf packages with zero internal dependencies (stdlib
only). Over 60 sub-packages, organized by domain. See
`internal/config/README.md` for the full decision tree.

| What you're adding              | Where it goes                     |
|---------------------------------|-----------------------------------|
| File names, extensions, paths   | `config/file/`, `config/dir/`     |
| Regex patterns                  | `config/regex/`                   |
| CLI flag names (`--flag-name`)  | `config/flag/flag.go`             |
| Flag description YAML keys      | `config/embed/flag/<cmd>.go`      |
| Command Use/DescKey strings     | `config/embed/cmd/<cmd>.go`       |
| User-facing text YAML keys      | `config/embed/text/<domain>.go`   |
| Time durations, thresholds      | `config/<domain>/`                |

#### The assets pipeline

User-facing text flows through a three-level chain:

1. **Go constant** (`config/embed/text/`) defines a string key:
   `DescKeyWriteAddedTo = "write.added-to"`
2. **Call site** resolves it: `desc.Text(text.DescKeyWriteAddedTo)`
3. **YAML** (`internal/assets/commands/text/write.yaml`) holds the
   actual text: `write.added-to: { short: "Added to %s" }`

The same pattern applies to command descriptions (`commands.yaml`),
flag descriptions (`flags.yaml`), and examples (`examples.yaml`).
The `TestDescKeyYAMLLinkage` test verifies every constant resolves
to a non-empty YAML value.

### Adding a New Session Parser

The journal system uses a `SessionParser` interface. To add support for a
new AI tool (e.g. Aider, Cursor):

1. Create `internal/journal/parser/<tool>.go`;
2. Implement parsing logic that returns `[]*Session`;
3. Register the parser in `FindSessions()` / `FindSessionsForCWD()`;
4. Use `config.Tool*` constants for the tool identifier;
5. Add test fixtures and parser tests.

Pattern to follow: the Claude Code JSONL parser in `internal/journal/parser/`.

!!! note "Multilingual session headers"
    The Markdown parser recognizes session header prefixes configured via
    `session_prefixes` in `.ctxrc` (default: `Session:`). To support a new
    language, users add a prefix to their `.ctxrc` - no code change needed.
    New parser implementations can use `rc.SessionPrefixes()` if they also
    need prefix-based header detection.

### Adding a Bundled Skill

1. Create `internal/assets/claude/skills/<skill-name>/SKILL.md`;
2. Follow the skill format: trigger, negative triggers, steps, quality gate;
3. Run `make plugin-reload` and restart Claude Code to test;
4. Add a `Skill` entry to `.claude-plugin/plugin.json` if user-invocable;
5. Document in `docs/reference/skills.md`.

Pattern to follow: any skill in `internal/assets/claude/skills/ctx-status/`.

### Test Expectations

- **Unit tests**: colocated with source (`foo.go` → `foo_test.go`);
- **Test helpers**: use `t.Helper()` so failures point to callers;
- **HOME isolation**: use `t.TempDir()` + `t.Setenv("HOME", ...)` for
  tests that touch `~/.claude/` or `~/.ctx/`;
- **rc.Reset()**: call after `os.Chdir` in tests that change working
  directory (rc caches on first access);
- **No network**: all tests run offline, use fixtures.

Run `make test` before submitting. Target: no failures, no skips.

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

**Clear the cache and restart**:

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

### Configuration Profiles

The repo ships two `.ctxrc` source profiles. The working copy (`.ctxrc`)
is gitignored and swapped between them:

| File          | Purpose                                                   |
|---------------|-----------------------------------------------------------|
| `.ctxrc.base` | Golden baseline: all defaults, no logging                 |
| `.ctxrc.dev`  | Dev profile: notify events enabled, verbose logging       |
| `.ctxrc`      | Working copy (*gitignored*: copied from one of the above) |

Use ctx commands to switch:

```bash
ctx config switch dev      # switch to dev profile
ctx config switch base     # switch to base profile
ctx config status          # show which profile is active
```

After cloning, run `ctx config switch dev` to get started with full logging.

See [Configuration](configuration.md) for the full `.ctxrc` option reference.

### Backups

Back up project context and global Claude Code data with:

```bash
ctx backup                    # both project + global (default)
ctx backup --scope project    # .context/, .claude/, ideas/ only
ctx backup --scope global     # ~/.claude/ only
```

Archives are saved to `/tmp/`. When `CTX_BACKUP_SMB_URL` is configured,
they are also copied to an SMB share. See
[CLI Reference: backup](../cli/backup.md) for details.

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
* Handle errors explicitly — use descriptive names (`readErr`,
  `writeErr`) not repeated `err`;
* No magic strings — all repeated literals go in `internal/config/`;
* Output goes through `internal/write/` packages, not `fmt.Print*`;
* Errors go through `internal/err/` constructors, not inline
  `fmt.Errorf`;
* See [Package Taxonomy](#package-taxonomy) and
  `.context/CONVENTIONS.md` for the full reference.

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
