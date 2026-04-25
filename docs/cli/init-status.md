---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Init and Status
icon: lucide/rocket
---

![ctx](../images/ctx-banner.png)

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
| `--merge`   |       | Auto-merge `ctx` content into existing `CLAUDE.md`                          |

**Creates**:

- `.context/` directory with all template files
- `.claude/settings.local.json` with pre-approved ctx permissions
- `CLAUDE.md` with bootstrap instructions (or merges into existing)

Claude Code hooks and skills are provided by the **`ctx` plugin**
(see [Integrations](../operations/integrations.md#claude-code-full-integration)).

**Example**:

```bash
# Standard init
ctx init

# Minimal setup (just core files)
ctx init --minimal

# Force overwrite existing
ctx init --force

# Merge into existing files
ctx init --merge
```

After `ctx init` succeeds, the final output includes a hint showing
the exact `eval "$(ctx activate)"` line to bind the new directory
for your shell. Every other `ctx` command requires that binding
(or an equivalent direct `CTX_DIR=/abs/path/.context` export) before
it will run.

---

### `ctx activate`

Emit a shell-native `export CTX_DIR=...` line for the target
`.context/` directory. `ctx` does not walk the filesystem during
operating commands; every non-exempt command requires `CTX_DIR`
set before it will run. `activate` is the convenience that figures
out the path and lets you bind it with one line.

```bash
# Walk up from CWD, emit if exactly one candidate visible.
eval "$(ctx activate)"
```

**Flags**:

| Flag      | Description                                                                              |
|-----------|------------------------------------------------------------------------------------------|
| `--shell` | Shell dialect override. POSIX-family (`bash`, `zsh`, `sh`) all share one syntax today; the flag exists for future fish/nushell/powershell support. Auto-detected from `$SHELL`. |

**Resolution**:

| Candidate count from CWD | Behavior                                                                 |
|--------------------------|--------------------------------------------------------------------------|
| Zero                     | Error. Use `ctx init` to create one, or `cd` closer to the project root. |
| One                      | Emit `export CTX_DIR=<path>` for that candidate.                         |
| Two or more              | Refuse. List every candidate. Re-run from a more specific cwd.           |

`activate` is args-free under the single-source-anchor model; the
explicit-path mode was removed because hub-client / hub-server
scenarios store at `~/.ctx/hub-data/` and never read `.context/`,
so they activate from the project root like everyone else. Direct
binding without a project-local scan is still available via
`export CTX_DIR=/abs/path/.context` or the inline form.

If the parent shell already has `CTX_DIR` set to a different value,
the output gains a leading `# ctx: replacing stale CTX_DIR=...`
comment so the user sees the change in `eval` output before the
replacement takes effect.

**See also**: [Activating a Context Directory](../recipes/activating-context.md)
for the full recipe including direnv setup and CI patterns.

---

### `ctx deactivate`

Emit a shell-native `unset CTX_DIR` line. Pairs with `activate`.

```bash
eval "$(ctx deactivate)"
```

**Flags**:

| Flag      | Description                                                                                                                                                                       |
|-----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--shell` | Shell dialect override. POSIX-family (`bash`, `zsh`, `sh`) all share one `unset` syntax today; the flag exists for future fish/nushell/powershell support. Auto-detected from `$SHELL`. |

`deactivate` does not touch the filesystem, doesn't require a
declared context directory, and never fails under normal operation;
unsetting an already-unset variable is a no-op across supported
shells.

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
| `--include-hub` | false | Include hub entries from `.context/hub/`             |

**How budget works**:

The budget controls how much context is included. Entries are selected
in priority tiers:

1. **Constitution**: always included in full (*inviolable rules*)
2. **Tasks**: all active tasks, up to 40% of budget
3. **Conventions**: all conventions, up to 20% of budget
4. **Decisions**: scored by recency and relevance to active tasks
5. **Learnings**: scored by recency and relevance to active tasks
6. **[Steering](steering.md)**: applicable steering file bodies,
   scored by their `inclusion` mode and description match
   against the active prompt
7. **Skill**: named skill content (from `--skill`)
8. **Hub**: entries from `.context/hub/` (with `--include-hub`,
   see [`ctx connect`](connection.md))

Decisions and learnings are ranked by a combined score (how recent + how
relevant to your current tasks). High-scoring entries are included with
their full body. Entries that don't fit get title-only summaries in an
"Also Noted" section. Superseded entries are excluded.

**Output Sections**:

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
