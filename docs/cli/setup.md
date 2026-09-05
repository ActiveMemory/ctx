---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Setup
icon: lucide/toy-brick
---

![ctx](../images/ctx-banner.png)

## `ctx setup`

Generate AI tool integration configuration.

```bash
ctx setup <tool> [flags]
```

**Flags**:

| Flag      | Short | Description                                                                 |
|-----------|-------|-----------------------------------------------------------------------------|
| `--write` | `-w`  | Write the generated config to disk (e.g. `.github/copilot-instructions.md`) |

**Supported tools**:

| Tool          | Description                                                        |
|---------------|--------------------------------------------------------------------|
| `agents`      | Generic `AGENTS.md` (read natively by Codex, OpenCode, and others) |
| `claude-code` | Redirects to plugin install instructions                           |
| `codex`       | OpenAI Codex CLI (hooks, MCP, skills, `AGENTS.md`)                 |
| `cursor`      | Cursor IDE                                                         |
| `kiro`        | Kiro IDE                                                           |
| `cline`       | Cline (VS Code extension)                                          |
| `aider`       | Aider CLI                                                          |
| `copilot`     | GitHub Copilot                                                     |
| `copilot-cli` | GitHub Copilot CLI (instructions, skills, agent, MCP)              |
| `opencode`    | OpenCode (terminal-first AI coding agent)                          |
| `pi`          | Pi coding-agent CLI (pi.dev)                                       |
| `windsurf`    | Windsurf IDE                                                       |

!!! note "Claude Code Uses the Plugin System"
    Claude Code integration is now provided via the `ctx` plugin.
    Running `ctx setup claude-code` prints plugin install instructions.

**`ctx setup codex`** without `--write` prints the integration overview
and what it detected (Codex binary on PATH, plugin installed, plugin
enabled). With `--write` it deploys the project-local integration:

| File | Behavior |
|------|----------|
| `.codex/hooks.json` | Create, or merge: foreign hook groups preserved, `ctx`-managed groups replaced; an unparseable file is left untouched with a warning |
| `.codex/config.toml` | Append a `[mcp_servers.ctx]` table when the header is absent; skip when present |
| `AGENTS.md` | Marker-merged (same deployer as `ctx setup agents --write`) |
| `.agents/skills/ctx-*/SKILL.md` | Create, refresh when stale, skip with a warning when not `ctx`-managed |

When the `ctx` plugin is enabled in `~/.codex/config.toml`, the
deployer skips hooks, MCP, and skills (Codex would run them twice) and
deploys `AGENTS.md` only. Either way, the hooks do not run until you
trust them: start `codex`, run `/hooks`, and trust the `ctx` entries.
See [`ctx` for Codex](../home/codex.md).

**Examples**:

```bash
# Print hook instructions to stdout
ctx setup cursor
ctx setup aider

# Generate and write .github/copilot-instructions.md
ctx setup copilot --write

# Generate MCP config and sync steering files
ctx setup kiro --write
ctx setup cursor --write
ctx setup cline --write

# Generate OpenCode plugin, skills, AGENTS.md, and global MCP config
ctx setup opencode --write

# Show Codex integration state; then deploy .codex/hooks.json,
# .codex/config.toml, AGENTS.md, and .agents/skills/
ctx setup codex
ctx setup codex --write

# Generate Pi extension, skills, and AGENTS.md
ctx setup pi --write
```
