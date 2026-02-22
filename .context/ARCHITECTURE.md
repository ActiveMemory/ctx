# Architecture

## Overview

ctx is a CLI tool that creates and manages a `.context/` directory
containing structured markdown files. These files provide persistent,
token-budgeted, priority-ordered context for AI coding assistants
across sessions.

Design philosophy:

- **Markdown-centric**: all context is plain markdown; no databases,
  no proprietary formats. Files are human-readable and version-
  controlled alongside the code they describe.
- **Token-budgeted**: context assembly respects configurable token
  limits so AI agents receive the most important information first
  without exceeding their context window.
- **Priority-ordered**: files are loaded in a deliberate sequence
  (rules before tasks, conventions before architecture) so agents
  internalize constraints before acting.
- **Convention over configuration**: sensible defaults with optional
  `.ctxrc` overrides. No config file required to get started.

## Package Dependency Graph

```mermaid
graph TD
    config["config<br/>(constants, regex, file names)"]
    tpl["tpl<br/>(embedded templates)"]

    rc["rc<br/>(runtime config)"] --> config
    context["context<br/>(loader)"] --> rc
    context --> config
    drift["drift<br/>(detector)"] --> config
    drift --> context
    index["index<br/>(reindexing)"] --> config
    task["task<br/>(parsing)"] --> config
    validation["validation<br/>(sanitize)"] --> config
    recall_parser["recall/parser<br/>(session parsing)"] --> config
    claude["claude<br/>(hooks, skills)"] --> config
    claude --> tpl

    bootstrap["bootstrap<br/>(CLI entry)"] --> rc
    bootstrap --> cli_all["cli/* (19 commands)"]
    cli_all --> config
    cli_all --> rc
    cli_all --> context
    cli_all --> drift
    cli_all --> index
    cli_all --> task
    cli_all --> validation
    cli_all --> recall_parser
    cli_all --> claude
    cli_all --> tpl
```

`config` and `tpl` are the two foundation packages with zero internal
dependencies. Everything else builds upward from them.

## Component Map

### Foundation Packages

| Package           | Purpose                                                                                                                                                                          | Depends On |
|-------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------|
| `internal/config` | Constants, regex patterns, file names, read order, permissions (9 files covering config, file names, directories, entries, headings, markers, patterns, regex, token formatting) | (none)     |
| `internal/assets`    | Embedded templates via go:embed; provides access to all `.context/` scaffolds, skill definitions, hook scripts, and tools                                                        | (none)     |

### Core Packages

| Package                  | Purpose                                                                                                                                       | Depends On                            |
|--------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------|
| `internal/rc`            | Runtime configuration from `.ctxrc`, env vars, and CLI flags; singleton with sync.Once caching                                            | `internal/config`                     |
| `internal/context`       | Load `.context/` directory: read .md files, estimate tokens, generate summaries, detect empty files                                           | `internal/rc`, `internal/config`      |
| `internal/drift`         | Detect stale or invalid context: dead path references, completed-task buildup, potential secrets, missing required files                      | `internal/config`, `internal/context` |
| `internal/index`         | Generate and update markdown index tables in DECISIONS.md and LEARNINGS.md                                                                    | `internal/config`                     |
| `internal/task`          | Parse task checkboxes in TASKS.md; match state, indentation, content                                                                          | `internal/config`                     |
| `internal/validation`    | Input sanitization (filenames)                                                                                                                | `internal/config`                     |
| `internal/recall/parser` | Parse AI session transcripts (JSONL) into structured data; extensible parser registry supporting Claude Code (and designed for Aider, Cursor) | `internal/config`                     |
| `internal/claude`        | Generate Claude Code integration: hooks, skills, settings, permissions                                                                        | `internal/config`, `internal/assets`     |

### Entry Point

| Package              | Purpose                                                                     | Depends On                                   |
|----------------------|-----------------------------------------------------------------------------|----------------------------------------------|
| `internal/bootstrap` | Create root Cobra command, register global flags, attach all 19 subcommands | `internal/rc`, all `internal/cli/*` packages |

### CLI Commands (`internal/cli/*`)

| Command      | Purpose                                                                         | Key Dependencies                                                     |
|--------------|---------------------------------------------------------------------------------|----------------------------------------------------------------------|
| `add`        | Append entries to context files (decisions, tasks, learnings, conventions)      | `config`, `rc`, `index`, `initialize`                                |
| `agent`      | Generate AI-ready context packets with token budgeting                          | `config`, `rc`, `context`, `task`, `initialize`                      |
| `compact`    | Archive completed tasks, clean up context files                                 | `config`, `rc`, `context`, `task`, `add`, `complete`, `initialize`   |
| `complete`   | Mark tasks as done in TASKS.md                                                  | `config`, `rc`, `task`, `add`, `initialize`                          |
| `decision`   | Manage DECISIONS.md (add, list, reindex)                                        | `config`, `rc`, `index`                                              |
| `drift`      | Detect stale/invalid context and report issues                                  | `config`, `rc`, `context`, `drift`, `task`, `tpl`, `initialize`      |
| `hook`       | Generate AI tool integration configs (Claude, Cursor, Aider, Copilot, Windsurf) | (Cobra only)                                                         |
| `initialize` | Create `.context/` directory, deploy templates, generate hooks, merge settings  | `config`, `rc`, `claude`, `tpl`                                      |
| `journal`    | Export and synthesize sessions; generate static site                            | `config`, `rc`                                                       |
| `learnings`  | Manage LEARNINGS.md (add, list, reindex)                                        | `config`, `rc`, `index`                                              |
| `load`       | Output assembled context in read order with token budgeting                     | `config`, `rc`, `context`, `initialize`                              |
| `loop`       | Generate Ralph loop scripts for iterative AI workflows                          | `config`                                                             |
| `recall`     | Browse AI session history (list, show, export)                                  | `config`, `rc`, `recall/parser`                                      |
| `serve`      | Serve static journal site locally                                               | `rc`                                                                 |
| `session`    | Manage session snapshots (save, list, load, parse)                              | `config`, `rc`, `validation`, `initialize`                           |
| `status`     | Show context summary: file list, tokens, summaries                              | `config`, `rc`, `context`, `initialize`                              |
| `sync`       | Sync project files with context; detect deps, suggest updates                   | `config`, `context`, `initialize`                                    |
| `task`       | Task management utilities; archive completed tasks                              | `config`, `rc`, `task`, `validation`, `add`, `compact`, `initialize` |
| `watch`      | Monitor stdin for context updates; auto-save sessions                           | `config`, `rc`, `context`, `task`, `add`, `initialize`               |

## Data Flow

### 1. `ctx init` -- Initialization

```mermaid
sequenceDiagram
    participant User
    participant CLI as cli/initialize
    participant TPL as tpl
    participant Claude as claude
    participant FS as Filesystem

    User->>CLI: ctx init [--minimal]
    CLI->>FS: Create .context/ directory
    CLI->>TPL: Read embedded templates
    TPL-->>CLI: Template content
    CLI->>FS: Write context files (CONSTITUTION, TASKS, etc.)
    CLI->>Claude: DefaultHooks(projectDir)
    Claude->>TPL: Read hook scripts
    Claude-->>CLI: Settings (hooks + permissions)
    CLI->>FS: Write .claude/settings.local.json
    CLI->>TPL: Read skill templates
    CLI->>FS: Deploy .claude/skills/*/SKILL.md
    CLI->>TPL: Read CLAUDE.md template
    CLI->>FS: Write or merge CLAUDE.md
    CLI->>FS: Deploy Makefile.ctx, context-watch.sh
    CLI-->>User: Initialized .context/ with N files
```

### 2. `ctx agent` -- Context Packet Assembly

```mermaid
sequenceDiagram
    participant Agent as AI Agent
    participant CLI as cli/agent
    participant RC as rc
    participant Ctx as context
    participant FS as Filesystem

    Agent->>CLI: ctx agent --budget 4000
    CLI->>RC: TokenBudget(), PriorityOrder()
    RC-->>CLI: Budget limit, file order
    CLI->>Ctx: Load(contextDir)
    Ctx->>FS: Read all .md files
    FS-->>Ctx: File contents
    Ctx-->>CLI: Context (files, tokens, summaries)
    CLI->>CLI: Sort by priority (FileReadOrder)
    CLI->>CLI: Truncate to budget
    CLI-->>Agent: Markdown context packet
```

### 3. `ctx drift` -- Drift Detection

```mermaid
sequenceDiagram
    participant User
    participant CLI as cli/drift
    participant Ctx as context
    participant Det as drift.Detect
    participant FS as Filesystem

    User->>CLI: ctx drift [--json]
    CLI->>Ctx: Load(contextDir)
    Ctx-->>CLI: Context with all files
    CLI->>Det: Detect(ctx)
    Det->>Det: checkPathReferences (ARCHITECTURE, CONVENTIONS)
    Det->>FS: os.Stat() each backtick path
    Det->>Det: checkStaleness (completed task count)
    Det->>Det: checkConstitution (secret file scan)
    Det->>Det: checkRequiredFiles (CONSTITUTION, TASKS, DECISIONS)
    Det-->>CLI: Report (warnings, violations, passed)
    CLI-->>User: Formatted report or JSON
```

### 4. `ctx recall export` -- Session Export

```mermaid
sequenceDiagram
    participant User
    participant CLI as cli/recall
    participant Parser as recall/parser
    participant FS as Filesystem

    User->>CLI: ctx recall export [session-id]
    CLI->>Parser: FindSessionsForCWD(cwd)
    Parser->>FS: Scan ~/.claude/projects/
    Parser->>Parser: Match by git remote or path
    FS-->>Parser: JSONL session files
    Parser->>Parser: ParseFile (Claude JSONL format)
    Parser-->>CLI: []Session (messages, tools, tokens)
    CLI->>CLI: Format as markdown
    CLI-->>User: Session transcript
```

### 5. Hook Lifecycle (via ctx plugin)

```mermaid
sequenceDiagram
    participant CC as Claude Code
    participant CTX as ctx system
    participant FS as Filesystem

    Note over CC: Event: PreToolUse
    CC->>CTX: ctx system block-non-path-ctx
    CTX->>CTX: Check if tool invokes ctx via absolute path
    CTX-->>CC: BLOCK or ALLOW

    Note over CC: Event: UserPromptSubmit
    CC->>CTX: ctx system check-context-size
    CTX->>CTX: Increment counter, check thresholds
    CTX-->>CC: Checkpoint message (or silent)

    Note over CC: Event: SessionEnd
    CC->>CTX: ctx system cleanup-tmp
    CTX->>FS: Remove stale temp files
    CTX-->>CC: (silent)
```

## Context File Lifecycle

```mermaid
stateDiagram-v2
    [*] --> Empty: Project created
    Empty --> Populated: ctx init
    Populated --> Active: ctx add / manual edits
    Active --> Active: ctx add, ctx complete
    Active --> Stale: No updates, drift detected
    Stale --> Active: ctx compact / consolidation
    Active --> Archived: ctx task archive
    Archived --> [*]
```

## Key Architectural Patterns

### Priority-Based File Ordering

Files load in a deliberate sequence defined by `config.FileReadOrder`:

1. CONSTITUTION (rules the agent must not violate)
2. TASKS (what to work on now)
3. CONVENTIONS (how to write code)
4. ARCHITECTURE (system structure)
5. DECISIONS (why things are this way)
6. LEARNINGS (gotchas and tips)
7. GLOSSARY (domain terms)
8. AGENT_PLAYBOOK (how to use this system)

Overridable via `priority_order` in `.ctxrc`.

### Token Budgeting

Token estimation uses a 4-characters-per-token heuristic
(see the context package). When the total context exceeds the
budget (default 8000, configurable via `CTX_TOKEN_BUDGET` or
`.ctxrc`), lower-priority files are truncated or omitted.
Higher-priority files always get included first.

### Structured Entry Headers

Decisions and learnings use timestamped headers for chronological
ordering and index generation:

```
## [2026-01-28-143022] Use PostgreSQL for primary database
```

The regex `config.RegExEntryHeader` parses these across the codebase.

### Runtime Config Hierarchy

Configuration resolution (highest priority wins):

1. CLI flags (`--context-dir`)
2. Environment variables (`CTX_DIR`, `CTX_TOKEN_BUDGET`)
3. `.ctxrc` file (YAML)
4. Hardcoded defaults in `internal/rc`

Managed by `internal/rc` with sync.Once singleton caching.

### Extensible Session Parsing

`internal/recall/parser` defines a `SessionParser` interface. Each
AI tool (Claude Code, potentially Aider, Cursor) registers its own
parser. Currently only Claude Code JSONL is implemented
(see `internal/recall/parser`). Session matching uses git
remote URLs, relative paths, and exact CWD matching.

### Template and Live Skill Dual-Deployment

Skills exist in two locations:

- **Templates** (`internal/assets/claude/skills/`): embedded in the
  binary, deployed on `ctx init`
- **Live** (`.claude/skills/`): project-local copies that the user
  and agent can edit

`ctx init` deploys templates to live. The `/update-docs` skill
checks for drift between them.

## File Layout

```
ctx/
├── cmd/ctx/                     # Main entry point (main.go)
├── internal/
│   ├── bootstrap/               # CLI initialization, command registration
│   ├── claude/                  # Claude Code hooks, skills, settings
│   ├── cli/                     # 19 command packages
│   │   ├── add/                 #   ctx add
│   │   ├── agent/               #   ctx agent
│   │   ├── compact/             #   ctx compact
│   │   ├── complete/            #   ctx complete
│   │   ├── decision/            #   ctx decision
│   │   ├── drift/               #   ctx drift
│   │   ├── hook/                #   ctx hook
│   │   ├── initialize/          #   ctx init
│   │   ├── journal/             #   ctx journal
│   │   ├── learnings/           #   ctx learnings
│   │   ├── load/                #   ctx load
│   │   ├── loop/                #   ctx loop
│   │   ├── recall/              #   ctx recall
│   │   ├── serve/               #   ctx serve
│   │   ├── session/             #   ctx session
│   │   ├── status/              #   ctx status
│   │   ├── sync/                #   ctx sync
│   │   ├── task/                #   ctx task
│   │   └── watch/               #   ctx watch
│   ├── config/                  # Constants, regex, file names, read order
│   ├── context/                 # Context loading, token estimation
│   ├── drift/                   # Drift detection engine
│   ├── index/                   # Index table generation for DECISIONS/LEARNINGS
│   ├── rc/                      # Runtime config (.ctxrc, env, CLI flags)
│   ├── recall/
│   │   └── parser/              # Session transcript parsing
│   ├── task/                    # Task checkbox parsing
│   ├── tpl/                     # Embedded templates (go:embed)
│   │   ├── claude/
│   │   │   ├── hooks/           #   Hook scripts (4 .sh files)
│   │   │   └── skills/          #   Skill templates (16 directories)
│   │   ├── entry-templates/     #   Decision/learning entry templates
│   │   ├── ralph/               #   Ralph loop PROMPT.md
│   │   └── tools/               #   context-watch.sh, cleanup-ctx-tmp.sh
│   └── validation/              # Input sanitization
├── docs/                        # Documentation site source (mkdocs)
│   ├── cli-reference.md
│   ├── context-files.md
│   ├── integrations.md
│   ├── prompting-guide.md
│   ├── session-journal.md
│   └── blog/
├── site/                        # Generated static site
├── hack/                        # Build and release scripts
├── .context/                    # This project's own context files
│   ├── CONSTITUTION.md
│   ├── TASKS.md
│   ├── CONVENTIONS.md
│   ├── ARCHITECTURE.md          # (this file)
│   ├── DECISIONS.md
│   ├── LEARNINGS.md
│   ├── GLOSSARY.md
│   ├── AGENT_PLAYBOOK.md
│   ├── sessions/                # Session snapshots
│   └── archive/                 # Archived tasks
└── .claude/                     # Claude Code integration
    ├── settings.local.json      # Hooks and permissions
    └── skills/                  # Live skill definitions (22 skills)
```
