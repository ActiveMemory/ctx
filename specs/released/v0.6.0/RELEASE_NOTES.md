<img src="https://ctx.ist/images/ctx-banner.png" />

# Context CLI v0.6.0

This release replaces shell hook scripts with native Go subcommands and ships hooks and skills as a Claude Code plugin. `ctx init` is now tool-agnostic: it no longer scaffolds `.claude/hooks/` or `.claude/skills/`. The plugin model eliminates the jq dependency, enables `go test` coverage for all hook logic, and makes distribution a single install command.

## Highlights

- **Plugin-based Distribution**: Hooks and skills ship as a Claude Code plugin. Install with `/plugin marketplace add ActiveMemory/ctx` then `/plugin install ctx@activememory-ctx`. No build step required; the plugin is served directly from source.
- **Shell Hooks to Go Subcommands**: All 6 shell scripts replaced by `ctx system *` commands compiled into the binary. Zero external dependencies.
- **Obsidian Vault Export**: `ctx journal obsidian` generates a full Obsidian vault from enriched journal entries with wikilinks, MOC pages, and graph-optimized cross-linking.
- **Encrypted Scratchpad**: `ctx pad` provides a git-friendly encrypted scratchpad (AES-256-GCM) for sensitive one-liners that travel with the project.
- **Security Hardening**: Path boundary validation, symlink detection, and user-specific temp directories close the medium-severity findings from the security audit.

## Breaking Changes

- `ctx init` no longer creates `.claude/hooks/` or `.claude/skills/`. Install the ctx plugin instead.
- `ctx hook claude-code` now prints plugin install instructions instead of generating shell scripts.
- Version jumps from 0.3.0 to 0.6.0 to signal the magnitude of the plugin conversion.

## Features

- Add `ctx system` subcommands: `check-context-size`, `check-persistence`, `check-journal`, `post-commit`, `block-non-path-ctx`, `cleanup-tmp`
- Add Claude Code plugin with marketplace.json, hooks.json, and 25 skills
- Serve plugin directly from `internal/tpl/claude/`; eliminate `make plugin` build step
- Add `ctx journal obsidian` command with wikilink conversion, frontmatter transformation, MOC generation, and related-sessions footer
- Add `ctx pad` command suite: `show`, `edit`, `clear`, `edit --append`, `edit --prepend`
- Add `ctx permissions snapshot` and `ctx permissions restore` for settings.local.json management
- Add `allow_outside_cwd` option to `.contextrc` for path boundary override
- Add `ctx init` auto-append of recommended `.gitignore` entries
- Add `Context.File()` lookup method for programmatic context file access
- Add journal reminder hook detecting unexported sessions and unenriched entries
- Add SessionEnd cleanup hook removing stale temp files
- Add persistence nudge hook with adaptive frequency based on prompt count
- Add `/check-links` skill for dead link auditing
- Add `/ctx-pad` skill for scratchpad interaction
- Add `/ctx-worktree` skill for parallel agent development with git worktrees
- Add `/ctx-borrow` skill for extracting and applying deltas between project copies
- Add `/sanitize-permissions` skill for settings.local.json security auditing

## Bug Fixes

- Hooks no longer create partial `.context/` (logs only) before `ctx init` runs
- `ctx init` treats `.context/` with only logs as uninitialized; skips overwrite prompt
- Fix CodeQL int64-to-int truncation warning in persistence state parser
- Fix UTF-8 safe string truncation preventing mid-rune splits
- Fix 18 golangci-lint issues across pad, compact, crypto, and validation packages
- Fix hook output channels: stderr is invisible for UserPromptSubmit hooks
- Fix outdated context-update XML syntax in docs
- Remove all stale session/save references from docs, skills, and source

## Security

- Add path boundary validation on `--context-dir` / `CTX_DIR` preventing operations outside project root (M-1)
- Add symlink detection with `Lstat()` before file read/write in `.context/` (M-2)
- Use `$XDG_RUNTIME_DIR/ctx` or user-specific temp subdirectory for state files (M-3)
- Add `/sanitize-permissions` skill for auditing dangerous Bash permissions

## Refactoring

- Split `config/tpl.go` (349 lines) into 4 feature-area files
- Convert `formatToolUse` switch to dispatch map
- Move 6 utility functions from `recall/run.go` to `recall/fmt.go`
- Extract `findSessions()` and `deployHookScript()` helpers to reduce duplication
- Unify task archiving into `compact.WriteArchive` helper
- Fix stale godocs and add `doc.go` for 6 packages
- Update `internal/tpl` package doc to reflect dual purpose (templates + plugin assets)

## Documentation

- Add migration guide for upgrading from shell hooks to plugin model
- Add first-session walkthrough with end-to-end interaction examples
- Add agent security docs and defense-in-depth blog post
- Add scratchpad docs, recipes, and autonomous loop hardening guide
- Add agent team decision framework recipe
- Reorder docs nav: promote recipes, demote adopting/upgrading
- Add recipe TL;DR admonitions for long recipes
- Add cross-references between blog posts and documentation pages
- 8 new blog posts covering skills anatomy, IRC bouncers, worktrees, and more

---

Full changelog: https://github.com/ActiveMemory/ctx/compare/v0.3.0...v0.6.0
