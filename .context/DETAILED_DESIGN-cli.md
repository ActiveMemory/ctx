# Detailed Design: CLI Layer

Modules: bootstrap, cli/parent, 34 command packages

## internal/bootstrap

**Purpose**: Create root Cobra command, register 34 subcommands
in 8 visible groups + 1 hidden group.

**Key types**: `registration` struct (cmd function + group ID)

**Exported API**:
- `RootCmd()`: creates root command with global flags
  (--context-dir, --allow-outside-cwd)
- `Initialize()`: builds group metadata, registers all commands

**PersistentPreRunE guards**:
1. Apply context-dir override from flag
2. Validate boundary (unless --allow-outside-cwd)
3. Require .context/ initialization (unless annotated SkipInit)
4. Skip for cobra builtins (completion commands)

**Command groups** (defined in group.go):
1. Getting Started: initialize, status, guide
2. Context: add, load, agent, sync, drift, compact
3. Artifacts: decision, learning, task
4. Sessions: journal, memory, remind, pad
5. Runtime: config, permission, pause, resume
6. Integration: setup, mcp, watch, notify, loop
7. Diagnostics: doctor, change, dep, why, trace
8. Utilities: reindex
9. Hidden: serve, site, system

**Danger zones**:
1. PersistentPreRunE runs for ALL subcommands — adding a command
   that should work without init requires AnnotationSkipInit.
2. Group ordering in help output is hardcoded in group.go.
3. Global flags are parsed before subcommand flags — conflicts
   between global and local flag names cause silent shadowing.

**Extension points**:
- Add new command: create package, add registration in group.go
- Add new group: define in group.go with display order

**Dependencies**: all cli/* packages, rc, validate, context/validate

---

## internal/cli/parent

**Purpose**: Generic factory for creating pure grouping commands
(commands with subcommands but no own Run handler).

**Exported API**: `Cmd(descKey, use string, subs ...*cobra.Command)`

Used by 10 commands: config, decision, journal, learning, memory,
mcp, permission, site, system, task.

**Dependencies**: assets/read/desc, spf13/cobra

---

## Command Taxonomy

All commands follow one of two patterns:

**Pattern A: cmd/root with Run()** (24 commands)
```
command/
  command.go         -- package Cmd() returns root.Cmd()
  cmd/root/
    cmd.go           -- creates cobra.Command with flags
    run.go           -- implements Run() logic
  core/              -- (optional) reusable logic
```

**Pattern B: parent.Cmd() with subcommands** (10 commands)
```
command/
  command.go         -- uses parent.Cmd(descKey, use, ...subs)
  cmd/
    sub1/cmd.go      -- subcommand 1
    sub2/cmd.go      -- subcommand 2
  core/              -- (optional) shared logic
```

---

## New Commands (since 2026-02-24)

### cli/change
**Purpose**: Detect code and context changes since a point in time.
**Flags**: --since (duration or date)
**Core**: detect/, scan/, render/ — git-based change detection

### cli/config
**Purpose**: Runtime config management.
**Subcommands**: switch (change profile), status (show active),
schema (dump .ctxrc JSON Schema)
**Core**: profile/ — profile loading and switching

### cli/dep
**Purpose**: Dependency analysis and graph rendering.
**Core**: builder/, golang/, node/, python/, rust/, render/ —
multi-ecosystem dependency tree with interface-based GraphBuilder

### cli/doctor
**Purpose**: Health checks for ctx installation and environment.
**Flags**: --json
**Annotation**: SkipInit (works without .context/)
**Core**: check/, output/ — pluggable health checks

### cli/guide
**Purpose**: Display tutorials and command help.
**Flags**: --skills, --commands
**Annotation**: SkipInit

### cli/mcp
**Purpose**: MCP server entry point.
**Subcommands**: serve

### cli/memory
**Purpose**: Bridge Claude Code auto memory into .context/.
**Subcommands**: sync, status, diff, importer, publish, unpublish
**Core**: discovery, mirroring, drift detection, bidirectional sync

### cli/parent
**Purpose**: Generic grouping command factory (see above).

### cli/pause / cli/resume
**Purpose**: Pause/resume context hooks for current session.
**Flags**: --session-id

### cli/reindex
**Purpose**: Regenerate indices for DECISIONS.md and LEARNINGS.md
in one invocation (convenience wrapper).

### cli/setup
**Purpose**: Generate AI tool integration configs.
**Flags**: --write
**Annotation**: SkipInit
Replaces old `ctx hook` command (Decision 2026-04-01).
Supports: claude, cursor, aider, copilot, windsurf.

### cli/site
**Purpose**: Site generation subcommands. Hidden.
**Subcommands**: feed (RSS generation)

### cli/trace
**Purpose**: Commit context tracing and event inspection.
**Subcommands**: show, collect, file, hook, tag
**Flags**: --last, --json

### cli/why
**Purpose**: Explain dependencies and relationships between
context files and codebase.

---

## Existing Commands (restructured)

All previously flat CLI packages were restructured into cmd/root +
core/ taxonomy per Decision 2026-03-06. Key structural changes:

- **add**: split into cmd/root/, core/entry/, core/example/,
  core/extract/, core/format/, core/insert/, core/normalize/
- **agent**: split into cmd/root/, core/budget/, core/cooldown/,
  core/extract/, core/score/, core/sort/
- **journal**: consumed recall (Decision 2026-03-30); subcommands:
  source, importer, lock, unlock, sync, site, obsidian
- **system**: expanded to 34 hook subcommands including various
  check_*, block_*, post_commit, mark_*, cleanup_*
- **recall**: deleted — merged into journal

**Danger zones**:
1. system/ has 34 subcommands, mostly hidden hook plumbing —
   modifying hook behavior affects all agent integrations.
2. journal/ is the largest CLI package — 24+ files in core/
   with two separate pipelines (site + obsidian).
3. add/ core packages are imported by entry/ and mcp/handler —
   changes to format or insert logic affect 3 callers.

**Extension points**:
- New CLI command: create package, register in bootstrap/group.go
- New journal format: implement SessionParser interface
- New AI tool integration: add template in setup command

**Dependencies per command**: each imports its own write/* and err/*
packages, plus domain packages as needed.
