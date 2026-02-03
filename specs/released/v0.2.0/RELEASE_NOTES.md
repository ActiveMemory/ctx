<img src="https://ctx.ist/images/ctx-banner.png" />

# Context CLI v0.2.0

This release brings significant improvements to session recall, code quality, and documentation. The new journal system enables AI-powered analysis of exported sessions, while extensive refactoring consolidates configuration constants and adds thread safety throughout the codebase.

## Highlights

- **Session Recall & Journal System**: Browse, search, and export AI session history with `ctx recall`, then analyze sessions with `ctx journal`
- **Quick Reference Indexes**: DECISIONS.md and LEARNINGS.md now include auto-generated indexes for faster scanning
- **Improved CLI Flags**: Global `--context-dir` and `--no-color` flags, plus required structured flags for decisions and learnings
- **Code Quality**: Consolidated configuration constants, thread-safe runtime config, and comprehensive test coverage

## Features

- Add `ctx recall` command for browsing AI session history across projects
- Add `ctx journal` command with site generation for session analysis
- Add quick reference index to DECISIONS.md and LEARNINGS.md with `ctx decisions reindex` and `ctx learnings reindex`
- Add global flags: `--context-dir` to override context directory, `--no-color` to disable colored output
- Add `.contextrc` configuration file support for project-level settings
- Add structured attributes to `<context-update>` XML format for richer metadata
- Require `--context`, `--rationale`, `--consequences` flags for `ctx add decision`
- Require `--context`, `--lesson`, `--application` flags for `ctx add learning`
- Add shell completion support via `ctx completion` (bash, zsh, fish, powershell)

## Bug Fixes

- Fix `ctx tasks archive` to handle nested content correctly
- Fix shell script linter warnings in release and tag scripts

## Refactoring

- Consolidate hardcoded strings into config constants (file names, env vars, Claude API types)
- Add thread safety with RWMutex for runtime configuration
- Extract shared helpers to eliminate code duplication (ReindexFile, ScanDirectory)
- Rename internal/templates to internal/tpl
- Use iota for enum-like constants
- Add CRLF-aware newline handling

## Documentation

- Add security page with vulnerability reporting guidelines
- Add version history page with release documentation links
- Update demo project with AGENT_PLAYBOOK.md, LEARNINGS.md, and specs examples
- Standardize Go docstrings with Parameters/Returns/Fields sections
- Add CLI output convention (use cmd.Print* instead of fmt.Print*)
- Document `ctx completion` command and `--all-projects` flag for recall commands

## Dependencies

- Bump golangci/golangci-lint-action from 6 to 9
- Bump actions/setup-go from 5 to 6
- Bump actions/checkout from 4 to 6

---

Full changelog: https://github.com/ActiveMemory/ctx/compare/v0.1.2...v0.2.0
