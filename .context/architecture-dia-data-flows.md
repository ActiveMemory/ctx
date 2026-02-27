# Data Flow Diagrams

Parent: [ARCHITECTURE.md](ARCHITECTURE.md)

## 1. `ctx init` — Initialization

```
User                     cli/initialize           assets           Filesystem
 │                            │                     │                  │
 │  ctx init [--minimal]      │                     │                  │
 │ ─────────────────────────► │                     │                  │
 │                            │  Read templates     │                  │
 │                            │ ──────────────────► │                  │
 │                            │  ◄────────────────  │                  │
 │                            │  Template bytes     │                  │
 │                            │                     │                  │
 │                            │  Create .context/                      │
 │                            │ ──────────────────────────────────────►│
 │                            │  Write CONSTITUTION, TASKS, etc.       │
 │                            │ ──────────────────────────────────────►│
 │                            │  Generate AES-256 key                  │
 │                            │ ──────────────────────────────────────►│
 │                            │  Deploy hooks + skills                 │
 │                            │ ──────────────────────────────────────►│
 │                            │  Merge settings.local.json             │
 │                            │ ──────────────────────────────────────►│
 │                            │  Write/merge CLAUDE.md                 │
 │                            │ ──────────────────────────────────────►│
 │  ◄────────────────────────                                          │
 │  "Initialized with N files"│                     │                  │
```

## 2. `ctx agent` — Context Packet Assembly

```
AI Agent              cli/agent              rc              context          FS
  │                      │                   │                  │              │
  │  ctx agent           │                   │                  │              │
  │  --budget 4000       │                   │                  │              │
  │ ───────────────────► │                   │                  │              │
  │                      │  TokenBudget()    │                  │              │
  │                      │ ────────────────► │                  │              │
  │                      │  ◄──────────────  │                  │              │
  │                      │  4000             │                  │              │
  │                      │                   │                  │              │
  │                      │  Load(dir)                           │              │
  │                      │ ──────────────────────────────────►  │              │
  │                      │                                      │  Read .md    │
  │                      │                                      │ ───────────► │
  │                      │                                      │  ◄─────────  │
  │                      │  ◄──────────────────────────────────                │
  │                      │  Context{files, tokens}              │              │
  │                      │                                      │              │
  │                      │  Score entries by                    │              │
  │                      │  recency + relevance                │              │
  │                      │  ┌──────────────┐                   │              │
  │                      │  │ Sort by score│                   │              │
  │                      │  │ Fit to budget│                   │              │
  │                      │  │ Overflow →   │                   │              │
  │                      │  │ "Also Noted" │                   │              │
  │                      │  └──────────────┘                   │              │
  │  ◄──────────────────                                       │              │
  │  Markdown packet     │                   │                  │              │
```

## 3. `ctx drift` — Drift Detection

```
User              cli/drift            context           drift.Detect           FS
  │                  │                    │                    │                  │
  │  ctx drift       │                    │                    │                  │
  │ ───────────────► │                    │                    │                  │
  │                  │  Load(dir)         │                    │                  │
  │                  │ ─────────────────► │                    │                  │
  │                  │  ◄───────────────  │                    │                  │
  │                  │  Context           │                    │                  │
  │                  │                    │                    │                  │
  │                  │  Detect(ctx)                            │                  │
  │                  │ ──────────────────────────────────────► │                  │
  │                  │                                         │ checkPathRefs   │
  │                  │                                         │ ──────────────► │
  │                  │                                         │ checkStaleness  │
  │                  │                                         │ checkConstitution
  │                  │                                         │ checkRequired   │
  │                  │                                         │ checkFileAge    │
  │                  │                                         │ checkEntryCount │
  │                  │                                         │ checkMissingPkgs│
  │                  │  ◄──────────────────────────────────────                  │
  │                  │  Report{warnings, violations}           │                  │
  │  ◄──────────────                                           │                  │
  │  Formatted report│                    │                    │                  │
```

## 4. `ctx recall export` — Session Export Pipeline

```
User           cli/recall        recall/parser        journal/state          FS
  │                │                    │                    │                 │
  │  ctx recall    │                    │                    │                 │
  │  export --all  │                    │                    │                 │
  │ ─────────────► │                    │                    │                 │
  │                │  FindSessionsFor   │                    │                 │
  │                │  CWD(cwd)          │                    │                 │
  │                │ ─────────────────► │                    │                 │
  │                │                    │  Scan              │                 │
  │                │                    │  ~/.claude/        │                 │
  │                │                    │  projects/         │                 │
  │                │                    │ ──────────────────────────────────►  │
  │                │                    │  ◄────────────────────────────────   │
  │                │                    │  Parse JSONL       │                 │
  │                │  ◄─────────────────                     │                 │
  │                │  []Session          │                    │                 │
  │                │                    │                    │                 │
  │                │  Load(journalDir)                       │                 │
  │                │ ──────────────────────────────────────► │                 │
  │                │  ◄──────────────────────────────────────                  │
  │                │  JournalState       │                    │                 │
  │                │                    │                    │                 │
  │                │  Plan: new/regen/skip/locked             │                 │
  │                │  Format as Markdown                      │                 │
  │                │  Write to .context/journal/              │                 │
  │                │ ──────────────────────────────────────────────────────────►│
  │                │  MarkExported()                          │                 │
  │                │ ──────────────────────────────────────► │                 │
  │  ◄────────────                                           │                 │
  │  "Exported N"  │                    │                    │                 │
```

## 5. Hook Lifecycle (Claude Code Plugin)

```
Claude Code                  ctx system                  .context/
     │                            │                          │
     │  ─── UserPromptSubmit ───  │                          │
     │  check-context-size        │                          │
     │ ─────────────────────────► │  Read/increment counter  │
     │                            │ ────────────────────────►│
     │  ◄─────────────────────────                           │
     │  (checkpoint msg or silent)│                          │
     │                            │                          │
     │  check-ceremonies          │                          │
     │ ─────────────────────────► │                          │
     │  ◄─────────────────────────                           │
     │  (nudge or silent)         │                          │
     │                            │                          │
     │  check-persistence         │                          │
     │ ─────────────────────────► │  Verify .context/ exists │
     │  ◄─────────────────────────                           │
     │                            │                          │
     │  ─── PreToolUse(Bash) ───  │                          │
     │  block-non-path-ctx        │                          │
     │ ─────────────────────────► │  Check tool invocation   │
     │  ◄─────────────────────────                           │
     │  BLOCK/ALLOW JSON          │                          │
     │                            │                          │
     │  ─── PreToolUse(Edit) ───  │                          │
     │  qa-reminder               │                          │
     │ ─────────────────────────► │                          │
     │  ◄─────────────────────────                           │
     │  (reminder or silent)      │                          │
     │                            │                          │
     │  ─── PostToolUse(Bash) ──  │                          │
     │  post-commit               │                          │
     │ ─────────────────────────► │  Detect git commit       │
     │  ◄─────────────────────────                           │
     │  (nudge or silent)         │                          │
     │                            │                          │
     │  ─── SessionEnd ─────────  │                          │
     │  cleanup-tmp               │                          │
     │ ─────────────────────────► │  Remove stale tmp files  │
     │  ◄─────────────────────────                           │
```
