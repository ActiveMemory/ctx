---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Init and Status
icon: lucide/rocket
---

### `ctx init`

Initialize a new `.context/` directory with template files.

```bash
ctx init [flags]
```

**Flags**:

| Flag        | Short | Description                                                                 |
|-------------|-------|-----------------------------------------------------------------------------|
| `--force`   | `-f`  | Overwrite existing context files                                            |
| `--minimal` | `-m`  | Only create essential files (`TASKS.md`, `DECISIONS.md`, `CONSTITUTION.md`) |
| `--merge`   |       | Auto-merge `ctx` content into existing `CLAUDE.md` and `PROMPT.md`          |
| `--ralph`   |       | Agent works autonomously without asking questions                           |

**Creates**:

- `.context/` directory with all template files
- `.claude/settings.local.json` with pre-approved ctx permissions
- `PROMPT.md` with session prompt (autonomous mode with `--ralph`)
- `IMPLEMENTATION_PLAN.md` with high-level project direction
- `CLAUDE.md` with bootstrap instructions (or merges into existing)

Claude Code hooks and skills are provided by the **ctx plugin**
(see [Integrations](../operations/integrations.md#claude-code-full-integration)).

**Example**:

```bash
# Collaborative mode (agent asks questions when unclear)
ctx init

# Autonomous mode (agent works independently)
ctx init --ralph

# Minimal setup (just core files)
ctx init --minimal

# Force overwrite existing
ctx init --force

# Merge into existing files
ctx init --merge
```

---

### `ctx status`

Show the current context summary.

```bash
ctx status [flags]
```

**Flags**:

| Flag        | Short | Description                   |
|-------------|-------|-------------------------------|
| `--json`    |       | Output as JSON                |
| `--verbose` | `-v`  | Include file contents summary |

**Output**:

- Context directory path
- Total files and token estimate
- Status of each file (*loaded, empty, missing*)
- Recent activity (*modification times*)
- Drift warnings if any

**Example**:

```bash
ctx status
ctx status --json
ctx status --verbose
```

---

### `ctx agent`

Print an AI-ready context packet optimized for LLM consumption.

```bash
ctx agent [flags]
```

**Flags**:

| Flag         | Default | Description                                                          |
|--------------|---------|----------------------------------------------------------------------|
| `--budget`   | 8000    | Token budget: controls content selection and prioritization          |
| `--format`   | md      | Output format: `md` or `json`                                        |
| `--cooldown` | 10m     | Suppress repeated output within this duration (requires `--session`) |
| `--session`  | (none)  | Session ID for cooldown isolation (e.g., `$PPID`)                    |

**How budget works**:

The budget controls how much context is included. Entries are selected
in priority tiers:

1. **Constitution**: always included in full (*inviolable rules*)
2. **Tasks**: all active tasks, up to 40% of budget
3. **Conventions**: all conventions, up to 20% of budget
4. **Decisions**: scored by recency and relevance to active tasks
5. **Learnings**: scored by recency and relevance to active tasks

Decisions and learnings are ranked by a combined score (how recent + how
relevant to your current tasks). High-scoring entries are included with
their full body. Entries that don't fit get title-only summaries in an
"Also Noted" section. Superseded entries are excluded.

**Output sections**:

| Section          | Source            | Selection                             |
|------------------|-------------------|---------------------------------------|
| Read These Files | all `.context/`   | Non-empty files in priority order     |
| Constitution     | `CONSTITUTION.md` | All rules (*never truncated*)         |
| Current Tasks    | `TASKS.md`        | All unchecked tasks (*budget-capped*) |
| Key Conventions  | `CONVENTIONS.md`  | All items (*budget-capped*)           |
| Recent Decisions | `DECISIONS.md`    | Full body, scored by relevance        |
| Key Learnings    | `LEARNINGS.md`    | Full body, scored by relevance        |
| Also Noted       | overflow          | Title-only summaries                  |

**Example**:

```bash
# Default (8000 tokens, markdown)
ctx agent

# Smaller packet for tight context windows
ctx agent --budget 4000

# JSON format for programmatic use
ctx agent --format json

# Pipe to file
ctx agent --budget 4000 > context.md

# With cooldown (hooks/automation: requires --session)
ctx agent --session $PPID
```

**Use case**: Copy-paste into AI chat, pipe to system prompt, or use in hooks.

---

### `ctx load`

Load and display assembled context as AI would see it.

```bash
ctx load [flags]
```

**Flags**:

| Flag                | Description                               |
|---------------------|-------------------------------------------|
| `--budget <tokens>` | Token budget for assembly (default: 8000) |
| `--raw`             | Output raw file contents without assembly |

**Example**:

```bash
ctx load
ctx load --budget 16000
ctx load --raw
```
