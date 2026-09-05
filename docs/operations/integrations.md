---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: AI Tools
icon: lucide/plug
---

![ctx](../images/ctx-banner.png)

## AI Tools

Context works with any AI tool that can read files. This guide covers setup 
for popular AI coding assistants.

!!! tip "Run From the Project Root"
    `ctx` reads `$PWD/.context/`. Run the commands on this page from
    the project root (the directory that holds `.context/` and
    `.git/`).

## Claude Code (Full Integration)

Claude Code has the deepest integration via the **`ctx` plugin**.

### Setup

First, install `ctx` and initialize your project:

```bash
ctx init
```

Then, install the `ctx` plugin in Claude Code:

```bash
# From the ctx repository
claude /plugin install ./internal/assets/claude

# Or from the marketplace
claude /plugin marketplace add ActiveMemory/ctx
claude /plugin install ctx@activememory-ctx
```

!!! warning "Ensure the Plugin Is Enabled"
    Installing a plugin registers it, but local installs may not
    auto-enable it globally. Verify `~/.claude/settings.json` contains:

    ```json
    { "enabledPlugins": { "ctx@activememory-ctx": true } }
    ```

    Without this, the plugin's hooks and skills won't appear in other
    projects. Running `ctx init` auto-enables the plugin; use
    `--no-plugin-enable` to skip this step.

This gives you:

| Component                     | Purpose                |
|-------------------------------|------------------------|
| `.context/`                   | All context files      |
| `CLAUDE.md`                   | Bootstrap instructions |
| Plugin hooks                  | Lifecycle automation   |
| Plugin skills                 | Agent Skills           |

### How It Works

```mermaid
graph TD
    A[Session Start] --> B[Claude reads CLAUDE.md]
    B --> C[PreToolUse hook runs]
    C --> D[ctx agent loads context]
    D --> E[Work happens]
    E --> F[Session End]
```

1. **Session start**: Claude reads `CLAUDE.md`, which tells it to check `.context/`
2. **First tool use**: `PreToolUse` hook runs `ctx agent` and emits the context
   packet (*subsequent invocations within the cooldown window are silent*)
3. **Next session**: Claude reads context files and continues with context

### Plugin Hooks

<!-- drift-check: cat internal/assets/claude/hooks/hooks.json | grep -c '"command"' -->
<!-- drift-check: diff <(grep -oP '"`ctx` system \K[^"]+' internal/assets/claude/hooks/hooks.json | sort) <(sed -n '/Plugin Hooks/,/Hook Configuration/p' docs/operations/integrations.md | grep -oP '`ctx` system \K[a-z-]+' | sort) -->
The `ctx` plugin provides lifecycle hooks implemented as Go subcommands
(`ctx system *`):

| Hook                             | Event                        | Purpose                                          |
|----------------------------------|------------------------------|--------------------------------------------------|
| `ctx system context-load-gate`   | PreToolUse (`.*`)            | Auto-inject context on first tool use            |
| `ctx system block-non-path-ctx`  | PreToolUse (`Bash`)          | Block `./ctx` or `go run`: force `$PATH` install |
| `ctx system qa-reminder`         | PreToolUse (`Bash`)          | Remind agent to lint/test before committing      |
| `ctx system specs-nudge`         | PreToolUse (`EnterPlanMode`) | Nudge agent to use project specs when planning   |
| `ctx system check-context-size`  | UserPromptSubmit             | Nudge context assessment as sessions grow        |
| `ctx system check-ceremonies`    | UserPromptSubmit             | Nudge /ctx-remember and /ctx-wrap-up adoption    |
| `ctx system check-persistence`   | UserPromptSubmit             | Remind to persist learnings/decisions            |
| `ctx system check-journal`       | UserPromptSubmit             | Remind to export/enrich journal entries          |
| `ctx system check-reminders`     | UserPromptSubmit             | Relay pending reminders at session start         |
| `ctx system check-version`       | UserPromptSubmit             | Warn when binary/plugin versions diverge         |
| `ctx system check-resources`     | UserPromptSubmit             | Warn when memory/swap/disk/load hit DANGER level |
| `ctx system check-knowledge`     | UserPromptSubmit             | Nudge when knowledge files grow large            |
| `ctx system check-map-staleness` | UserPromptSubmit             | Nudge when ARCHITECTURE.md is stale              |
| `ctx system heartbeat`           | UserPromptSubmit             | Session-alive signal with prompt count metadata  |
| `ctx system post-commit`         | PostToolUse (`Bash`)         | Nudge context capture and QA after git commits   |

A catch-all `PreToolUse` hook also runs `ctx agent` on every tool use
(with cooldown) to autoload context.

### Hook Configuration

The plugin's `hooks.json` wires everything automatically: no manual
configuration in `settings.local.json` needed:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": ".*",
        "hooks": [
          { "type": "command", "command": "ctx system context-load-gate" }
        ]
      },
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "ctx system block-non-path-ctx" }
        ]
      },
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "ctx system qa-reminder" }
        ]
      },
      {
        "matcher": "EnterPlanMode",
        "hooks": [
          { "type": "command", "command": "ctx system specs-nudge" }
        ]
      },
      {
        "matcher": ".*",
        "hooks": [
          { "type": "command", "command": "ctx agent --budget 4000 2>/dev/null || true" }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "ctx system post-commit" }
        ]
      }
    ],
    "UserPromptSubmit": [
      {
        "hooks": [
          { "type": "command", "command": "ctx system check-context-size" },
          { "type": "command", "command": "ctx system check-ceremonies" },
          { "type": "command", "command": "ctx system check-persistence" },
          { "type": "command", "command": "ctx system check-journal" },
          { "type": "command", "command": "ctx system check-reminders" },
          { "type": "command", "command": "ctx system check-version" },
          { "type": "command", "command": "ctx system check-resources" },
          { "type": "command", "command": "ctx system check-knowledge" },
          { "type": "command", "command": "ctx system check-map-staleness" },
          { "type": "command", "command": "ctx system heartbeat" }
        ]
      }
    ]
  }
}
```

### Customizing Token Budget and Cooldown

Edit the `PreToolUse` command to change the token budget or cooldown:

```text
"command": "ctx agent --budget 8000 --session $PPID >/dev/null || true"
"command": "ctx agent --budget 4000 --cooldown 5m --session $PPID >/dev/null || true"
```

The `--session $PPID` flag isolates the cooldown per session: `$PPID` resolves
to the Claude Code process PID, so concurrent sessions don't interfere.
The default cooldown is 10 minutes; use `--cooldown 0` to disable it.

### Verifying Setup

1. Start a new Claude Code session;
2. Ask: **"Do you remember?"**
3. Claude should cite specific context:
     * Current tasks from `.context/TASKS.md`;
     * Recent **decisions** or **learnings**;
     * Recent **session history** from `ctx journal`.

### Local Plugin Development

When developing `ctx` locally (*adding skills, hooks, or changing plugin
behavior*), Claude Code caches the plugin by version. You must bump
the version in **both** files and update the marketplace for changes
to take effect:

1. **Bump version** in both:
   * `internal/assets/claude/.claude-plugin/plugin.json` (*plugin manifest*),
   *`.claude-plugin/marketplace.json` (*marketplace listing*);

2. **Update the marketplace** in Claude Code:
   * Open the Plugins UI (*`/plugins` or Esc menu*),
   * Go to **Marketplaces** tab,
   * Select the `activememory-ctx` Marketplace,
   * Choose **Update marketplace**;

3. **Start a new Claude Code session**: skill changes aren't
   reflected in existing sessions.

!!! warning "Both Version Files Must Match"
    If you only bump `plugin.json` but not `marketplace.json` (or
    vice versa), Claude Code may not detect the update. Always bump
    both together.

### Troubleshooting

| Issue                 | Solution                                                     |
|-----------------------|--------------------------------------------------------------|
| Context not loading   | Check `ctx` is in PATH: `which ctx`                          |
| Hook errors           | Verify plugin is installed: `claude /plugin list`            |
| New skill not visible | Bump version in both `plugin.json` files, update marketplace |

### Manual Context Load

If hooks aren't working, manually load context:

```bash
# Get context packet
ctx agent --budget 4000

# Or paste into conversation
cat .context/TASKS.md
```

<!-- drift-check: ls internal/assets/claude/skills/ | wc -l -->
<!-- drift-check: diff <(ls internal/assets/claude/skills/ | sort) <(sed -n '/Agent Skills/,/Usage Examples/p' docs/operations/integrations.md | grep -oP '/\Kctx-[a-z-]+' | sort -u) -->
### Agent Skills

The `ctx` plugin ships Agent Skills following the
[agentskills.io specification](https://agentskills.io).

These are invoked in Claude Code with `/skill-name`.

#### Session Lifecycle Skills

| Skill                  | Description                                          |
|------------------------|------------------------------------------------------|
| `/ctx-remember`        | Recall project context at session start (ceremony)   |
| `/ctx-wrap-up`         | End-of-session context persistence (ceremony)        |
| `/ctx-status`          | Show context summary (tasks, decisions, learnings)   |
| `/ctx-agent`           | Get AI-optimized context packet                      |
| `/ctx-next`            | Suggest 1-3 concrete next actions from context       |
| `/ctx-commit`          | Commit with integrated context capture               |
| `/ctx-reflect`         | Review session and suggest what to persist           |
| `/ctx-remind`          | Manage session-scoped reminders                      |
| `/ctx-pause`           | Pause context hooks for this session                 |
| `/ctx-resume`          | Resume context hooks after a pause                   |

#### Context Persistence Skills

| Skill                 | Description                                        |
|-----------------------|----------------------------------------------------|
| `/ctx-task-add`       | Add a task to TASKS.md                             |
| `/ctx-learning-add`   | Add a learning to LEARNINGS.md                     |
| `/ctx-decision-add`   | Add a decision with context/rationale/consequence  |
| `/ctx-convention-add` | Add a coding convention to CONVENTIONS.md          |
| `/ctx-archive`        | Archive completed tasks                            |

#### Scratchpad Skills

| Skill                 | Description                                        |
|-----------------------|----------------------------------------------------|
| `/ctx-pad`            | Manage encrypted scratchpad entries                |

#### Session History Skills

| Skill                     | Description                                  |
|---------------------------|----------------------------------------------|
| `/ctx-history`             | Browse AI session history                    |
| `/ctx-journal-enrich`     | Enrich a journal entry with frontmatter/tags |
| `/ctx-journal-enrich-all` | Full journal pipeline: export if needed, then batch-enrich |

#### Blogging Skills

!!! tip "Blogging Is a Better Way of Creating Release Notes"
    The blogging workflow can also double as generating release notes:

    AI reads your git commit history and creates a "*narrative*",
    which is essentially what a *release note* is for.

| Skill                 | Description                                     |
|-----------------------|-------------------------------------------------|
| `/ctx-blog`           | Generate blog post from recent activity         |
| `/ctx-blog-changelog` | Generate blog post from commit range with theme |

#### Auditing & Health Skills

| Skill                  | Description                                                |
|------------------------|------------------------------------------------------------|
| `/ctx-doctor`          | Troubleshoot `ctx` behavior with structural health checks    |
| `/ctx-drift`           | Detect and fix context drift (structural + semantic)       |
| `/ctx-consolidate`     | Merge redundant learnings or decisions into denser entries |
| `/ctx-alignment-audit` | Audit doc claims against playbook instructions             |
| `/ctx-prompt-audit`    | Analyze session logs for vague prompts                     |
| `/check-links`         | Audit docs for dead internal and external links            |

#### Planning & Execution Skills

| Skill               | Description                                      |
|---------------------|--------------------------------------------------|
| `/ctx-loop`         | Generate a Ralph Loop iteration script           |
| `/ctx-task-out`     | Decompose a committed spec into a milestone plan |
| `/ctx-implement`    | Execute a plan step-by-step with checks          |
| `/ctx-plan-import` | Import Claude Code plan files into project specs |
| `/ctx-worktree`     | Manage git worktrees for parallel agents         |
| `/ctx-architecture`          | Build and maintain architecture maps             |

#### Usage Examples

```text
/ctx-status
/ctx-learning-add "Token refresh requires explicit cache invalidation"
/ctx-journal-enrich twinkly-stirring-kettle
```

Skills support partial matching where applicable (e.g., session slugs).

---

## OpenAI Codex

Codex CLI ships lifecycle hooks and plugins whose contract is a near
clone of Claude Code's (same event names, same stdin JSON, same
`decision: block` output). `ctx` reuses the same `ctx system` hook
runtime and delivers it to Codex through a plugin, project-local
files, and a rollout-transcript parser.

!!! tip "Full guide: [`ctx` for Codex](../home/codex.md)"
    The home-page guide covers both install routes step by step, the
    trust step, skills, MCP tools, journal import, known limitations,
    and troubleshooting. This section is the reference for the hook
    manifest and the drift checks that keep it honest.

### Setup

Install the plugin once (every project gets hooks, skills, and the
MCP server):

```bash
codex plugin marketplace add ActiveMemory/ctx
codex plugin add ctx@activememory-ctx

# Or, from a local checkout of the ctx repository
make codex-plugin-install
```

Or deploy project-local files instead (teams, CI, `codex exec`):

```bash
ctx setup codex --write
ctx init
```

Then start `codex`, run `/hooks`, and trust the `ctx` entries. Codex
does not run untrusted hooks.

!!! warning "One Route per Project"
    Codex loads every matching hook from every source, so an enabled
    plugin plus a project-local `.codex/hooks.json` runs each hook twice.
    `ctx setup codex --write` skips hooks, MCP, and skills when it finds
    the plugin enabled in `~/.codex/config.toml` and deploys `AGENTS.md`
    only.

### What Gets Created

| Route | File | Purpose |
|-------|------|---------|
| Plugin | `$CODEX_HOME/plugins/cache/activememory-ctx/ctx/<version>/` | Installed copy of `internal/assets/codex/` (`.codex-plugin/plugin.json`, `hooks/hooks.json`, `.mcp.json`, `skills/`) |
| Plugin | `~/.codex/config.toml` | `[plugins."ctx@activememory-ctx"] enabled = true`, written by Codex |
| Project | `.codex/hooks.json` | Lifecycle hooks (same manifest the plugin ships) |
| Project | `.codex/config.toml` | `[mcp_servers.ctx]` registering `ctx mcp serve` |
| Project | `AGENTS.md` | Agent instructions, read natively by Codex; marker-merged |
| Project | `.agents/skills/ctx-*/SKILL.md` | `ctx` skills, invoked as `$ctx-<name>` |

Project-local files load only for **trusted** projects
(`[projects."<abs path>"] trust_level = "trusted"` in
`~/.codex/config.toml`). `CODEX_HOME` relocates `~/.codex` for both
plugin detection and session discovery.

### How It Works

```mermaid
graph TD
    A[Session Start] --> B[SessionStart hook runs ctx agent]
    B --> C[Codex reads AGENTS.md]
    C --> D[Work happens: PreToolUse / PostToolUse / UserPromptSubmit nudges]
    D --> E[Session End]
    E --> F[SessionEnd hook runs ctx journal import]
```

1. **Session start**: the `SessionStart` hook prints `ctx agent --budget 8000`;
   Codex injects plain stdout as developer context. It re-fires after
   `compact`, so the packet survives compaction.
2. **During the session**: the same gates and nudges as Claude Code, on
   Codex's tool names (`update_plan` for planning, `apply_patch` for edits).
3. **Session end**: `ctx journal import --all -y` captures the rollout
   transcript into `.context/journal/` (incremental; Codex caps the hook
   at 3 seconds).

<!-- drift-check: grep -c '"type": "command"' internal/assets/codex/hooks/hooks.json -->
<!-- drift-check: diff <(grep -oP 'ctx system \K[a-z-]+' internal/assets/codex/hooks/hooks.json | sort -u) <(sed -n '/### Codex Hooks/,/### Codex Skills/p' docs/operations/integrations.md | grep -oP '`ctx system \K[a-z-]+' | sort -u) -->
### Codex Hooks

Every command in `internal/assets/codex/hooks/hooks.json` is prefixed
with `cd "$(git rev-parse --show-toplevel)" &&`: Codex runs hooks with
the session cwd, and `ctx` reads `$PWD/.context/`, so the anchor keeps a
subdirectory cwd from pointing `ctx` at the wrong project. A compliance
test checks the anchor, the event names, and that each command resolves
to a registered subcommand.

| Hook                              | Event                                   | Purpose                                                    |
|-----------------------------------|-----------------------------------------|------------------------------------------------------------|
| `ctx agent --budget 8000`         | SessionStart (all sources)              | Inject the context packet as developer context             |
| `ctx system context-load-gate`    | PreToolUse (`.*`)                       | Auto-inject context on first tool use                      |
| `ctx system block-non-path-ctx`   | PreToolUse (`Bash`)                     | Block `./ctx` or `go run`: force `$PATH` install           |
| `ctx system qa-reminder`          | PreToolUse (`Bash`)                     | Remind agent to lint/test before committing                |
| `ctx system specs-nudge`          | PreToolUse (`update_plan`)              | Nudge agent to use project specs when planning             |
| `ctx system post-commit`          | PostToolUse (`Bash`)                    | Nudge context capture and QA after git commits             |
| `ctx system check-task-completion`| PostToolUse (`apply_patch\|Edit\|Write`) | Detect silently completed tasks after a file edit          |
| `ctx system check-context-size`   | UserPromptSubmit                        | Nudge context assessment as sessions grow                  |
| `ctx system check-ceremony`       | UserPromptSubmit                        | Nudge `$ctx-remember` and `$ctx-wrap-up` adoption          |
| `ctx system check-persistence`    | UserPromptSubmit                        | Remind to persist learnings/decisions                      |
| `ctx system check-journal`        | UserPromptSubmit                        | Remind to import/enrich journal entries                    |
| `ctx system check-reminder`       | UserPromptSubmit                        | Relay pending reminders                                    |
| `ctx system check-version`        | UserPromptSubmit                        | Warn when binary/plugin versions diverge                   |
| `ctx system check-resource`       | UserPromptSubmit                        | Warn when memory/swap/disk/load hit DANGER level           |
| `ctx system check-knowledge`      | UserPromptSubmit                        | Nudge when knowledge files grow large                      |
| `ctx system check-map-staleness`  | UserPromptSubmit                        | Nudge when ARCHITECTURE.md is stale                        |
| `ctx system check-memory-drift`   | UserPromptSubmit                        | Nudge when auto-memory drifts from `.context/`             |
| `ctx system check-freshness`      | UserPromptSubmit                        | Warn when technology-dependent constants go unreviewed     |
| `ctx system check-skill-discovery`| UserPromptSubmit                        | One-shot mid-session tip surfacing easy-to-forget skills   |
| `ctx system heartbeat`            | UserPromptSubmit                        | Session-alive signal with prompt count metadata            |
| `ctx journal import --all -y`     | SessionEnd (`timeout: 3`)               | Import the session transcript into `.context/journal/`     |

Differences from the Claude Code manifest, and why:

| Claude Code                          | Codex                               | Reason                                                        |
|--------------------------------------|-------------------------------------|---------------------------------------------------------------|
| `PreToolUse` `.*` runs `ctx agent`   | `SessionStart` runs `ctx agent`     | Codex ignores plain text on `PreToolUse`; `SessionStart` stdout becomes developer context and re-fires on `compact` |
| `EnterPlanMode` matcher              | `update_plan` matcher               | Codex's planning tool is `update_plan`                        |
| `Edit` / `Write` matchers            | `apply_patch\|Edit\|Write`          | Codex file edits are `apply_patch`; the aliases keep parity   |
| `${CLAUDE_PROJECT_DIR}` anchor       | `$(git rev-parse --show-toplevel)`  | Codex exposes no project-dir variable to hooks                |
| `SessionEnd` has no timeout          | `timeout: 3`                        | Codex caps `SessionEnd` hooks at 3 seconds                    |

Not wired: `PermissionRequest`, `PreCompact`, `PostCompact`,
`SubagentStart`, `SubagentStop`, `Stop`. None carries a `ctx` behavior
today (`PreCompact`/`PostCompact` cannot return additional context;
`Stop` requires JSON-only output).

<!-- drift-check: ls internal/assets/codex/skills/ | wc -l -->
<!-- drift-check: diff <(comm -13 <(ls internal/assets/codex/skills/ | sort) <(ls internal/assets/claude/skills/ | sort)) <(sed -n '/### Codex Skills/,/### Codex Journal Import/p' docs/operations/integrations.md | grep -oP '`\Kctx-[a-z-]+(?=`)' | sort -u) -->
### Codex Skills

`internal/assets/codex/skills/` is generated from the Claude Code skills
by `hack/sync-codex-skills.sh`, which strips the Claude-only
`allowed-tools:` frontmatter key and omits the skills whose body only
makes sense inside Claude Code:

| Excluded skill            | Why                                    |
|---------------------------|----------------------------------------|
| `ctx-permission-sanitize` | Audits `.claude/settings.local.json`   |
| `ctx-plan-import`         | Imports `~/.claude/plans/`             |
| `ctx-dream`               | Headless `claude -p` cron + guard script |
| `ctx-skill-create`        | Authors Claude Code skills and plugins |

Everything else in the [Agent Skills](#agent-skills) list above is
available in Codex under the same name with a `$` prefix. `make
check-codex-skills` fails when the generated tree is stale; `make
sync-codex-skills` regenerates it (`make build` runs the sync).

### Codex Journal Import

`ctx journal import` scans `$CODEX_HOME/sessions/YYYY/MM/DD/rollout-*.jsonl`
(default `~/.codex/sessions`), matches rollouts to the current project by
the session's working directory, and imports them with the tool id
`codex`. Codex-injected user items (`<environment_context>`,
`<user_instructions>`, skill and permission preambles) are filtered out;
a rollout with no real user message is skipped. Filter with
`ctx journal source --tool codex`.

### Local Plugin Development

Codex caches plugins by version, like Claude Code. Bump `VERSION` and run
`make sync-version`: it updates
`internal/assets/codex/.codex-plugin/plugin.json` and
`.agents/plugins/marketplace.json` together with the Claude manifests
(`make check-version-sync` fails if any of them disagree). Then re-run
`codex plugin add ctx@activememory-ctx` and start a new session.

### Troubleshooting

| Issue                              | Solution                                                                                                  |
|------------------------------------|-----------------------------------------------------------------------------------------------------------|
| Nothing fires                      | Run `/hooks` in `codex` and trust the `ctx` entries                                                        |
| Project-local hooks ignored        | Trust the project: `[projects."<abs path>"] trust_level = "trusted"` in `~/.codex/config.toml`            |
| Every nudge appears twice          | Plugin and `.codex/hooks.json` both active; remove one (`codex plugin remove ctx@activememory-ctx`)       |
| `SessionEnd` timeout               | First import of a large backlog; run `ctx journal import --all` once by hand, later runs are incremental  |
| No Codex sessions in the journal   | Check `$CODEX_HOME`; rollouts match by working directory                                                  |

---

## Cursor IDE

Cursor can use context files through its system prompt or by reading 
files directly.

### Setup

```bash
# Generate Cursor configuration
ctx setup cursor

# Initialize context
ctx init --minimal
```

### Configuration

Add to Cursor settings (`.cursor/settings.json`):

```text
// split to multiple lines for readability
{
  "ai.systemPrompt": "Read .context/TASKS.md and 
  .context/CONVENTIONS.md before responding. 
  Follow rules in .context/CONSTITUTION.md.",
}
```

### Usage

1. Open your project in Cursor
2. Context files are available in the file tree
3. Reference them in prompts: 
   "Check .context/DECISIONS.md for our approach to..."

### Manual Context Injection

For more control, paste context directly:

```bash
# Get AI-ready packet
ctx agent --budget 4000 | pbcopy  # macOS
ctx agent --budget 4000 | xclip  # Linux
```

Paste into Cursor's chat.

---

## Aider

Aider works well with context files through its `--read` flag.

### Setup

```bash
# Generate Aider configuration
ctx setup aider

# Initialize context
ctx init
```

### Configuration

Create `.aider.conf.yml`:

```yaml
read:
  - .context/CONSTITUTION.md
  - .context/TASKS.md
  - .context/CONVENTIONS.md
  - .context/DECISIONS.md
```

### Usage

```bash
# Start Aider (reads context files automatically)
aider

# Or specify files explicitly
aider --read .context/TASKS.md --read .context/CONVENTIONS.md
```

### With Watch Mode

Run `ctx watch` alongside Aider to capture context updates:

```bash
# Terminal 1: Run Aider
aider 2>&1 | tee /tmp/aider.log

# Terminal 2: Watch for context updates
ctx watch --log /tmp/aider.log
```

---

## GitHub Copilot

GitHub Copilot integrates with `ctx` at three levels: an automated
instructions file, a VS Code Chat extension, and manual patterns.

### Setup

```bash
# Initialize context
ctx init

# Generate .github/copilot-instructions.md
ctx setup copilot --write
```

The `--write` flag creates `.github/copilot-instructions.md`, which
Copilot reads automatically at the start of every session. This file
contains your project's constitution rules, current tasks, conventions,
and architecture: giving Copilot persistent context without manual
copy-paste.

Re-run `ctx setup copilot --write` after updating your `.context/` files
to regenerate the instructions.

### VS Code Chat Extension (`@ctx`)

The **`ctx` VS Code extension** adds a `@ctx` chat participant to
GitHub Copilot Chat, giving you direct access to 45 context commands
from within the editor, plus automatic hooks on file save / git commit /
`.context/` changes / dependency-file edits, and a reminder status-bar
indicator.

!!! tip "Full guide: [`ctx` for VS Code](../home/vscode.md)"
    The home-page guide covers daily workflows, the full command list,
    natural-language routing, auto-bootstrap of the `ctx` CLI, troubleshooting,
    and "Verify It Works." This subsection is the install-and-pointers
    overview; the dedicated page is the authoritative reference.

#### Installation

The extension ships to the [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=activememory.ctx-context)
under publisher `activememory` (display name: *`ctx`: Persistent Context
for AI*). Install via the Extensions view or `code --install-extension`.

To build from source instead (requires Node.js 20+):

```bash
cd editors/vscode
npm ci
npm run build
npx @vscode/vsce package
code --install-extension ctx-context-<version>.vsix
```

Reload VS Code. Type `@ctx` in Copilot Chat to verify.

#### What Gets Created

| File | Purpose |
|------|---------|
| `.context/` | Project-local context directory (created by `ctx init`, not by the extension) |
| `.github/copilot-instructions.md` | Repository instructions Copilot reads natively; regenerated automatically when `.context/` files change |

The extension itself lives in VS Code's extension storage; no project
files beyond `.context/` and the Copilot instructions are added.

#### How It Works

- **Chat participant:** `@ctx` is registered with VS Code's Chat API; 45 slash commands route to dedicated handlers that shell out to the `ctx` CLI.
- **Automatic hooks:** file save → task-completion check; git commit → decision/learning prompt; `.context/` change → regenerate Copilot instructions; dependency-file change → `/map` prompt.
- **Status-bar reminder:** a `$(bell) ctx` indicator surfaces pending session reminders, refreshing every 5 minutes.
- **Natural language:** plain English after `@ctx` is routed to the nearest matching command.
- **Auto-bootstrap:** if the `ctx` CLI isn't on PATH, the extension downloads the correct platform binary from GitHub Releases and caches it.

#### Configuration

| Setting              | Default | Description                                                      |
|----------------------|---------|------------------------------------------------------------------|
| `ctx.executablePath` | `ctx`   | Path to the `ctx` binary. Set this if `ctx` is not in your `PATH`. |

### Session Persistence

`ctx init` creates a `.context/sessions/` directory for storing
session data from non-Claude tools. The Markdown session parser scans
this directory during `ctx journal`, enabling session history for
Copilot and other tools.

### Manual Patterns

These patterns work without the extension, using Copilot's built-in
file awareness:

**Pattern 1: Keep context files open**

Open `.context/CONVENTIONS.md` in a split pane. Copilot will reference it.

**Pattern 2: Reference in comments**

```typescript
// See .context/CONVENTIONS.md for naming patterns
// Following decision in .context/DECISIONS.md: Use PostgreSQL

function getUserById(id: string) {
  // Copilot now has context
}
```

**Pattern 3: Paste context into Copilot Chat**

```bash
ctx agent --budget 2000
```

Paste output into Copilot Chat for context-aware responses.

---

## OpenCode

OpenCode is a terminal-first AI coding agent. `ctx` integrates via
a thin lifecycle plugin, MCP server, and `AGENTS.md` instructions.

### Setup

```bash
# Generate OpenCode plugin, global MCP config, skills, and AGENTS.md
ctx setup opencode --write

# Initialize context
ctx init
```

### What Gets Created

| File | Purpose |
|------|---------|
| `.opencode/plugins/ctx.ts` | Lifecycle plugin (hooks to `ctx system`) |
| `~/.config/opencode/opencode.json` | Global MCP server registration (or `$OPENCODE_HOME/opencode.json`) |
| `AGENTS.md` | Agent instructions (read natively) |
| `.opencode/skills/ctx-*/SKILL.md` | `ctx` skills |

### How It Works

The plugin wires OpenCode lifecycle events to `ctx system`:

- **`session.created`**: warms `ctx` state in the background (bootstrap + agent packet) so MCP queries are fast on first use.
- **`tool.execute.after` (shell, on `git commit`)**: runs `ctx system post-commit`.
- **`tool.execute.after` (edit/write)**: runs `ctx system check-task-completion`.
- **`session.idle`**: runs persistence and task-completion checks (silent: output is buffered, not surfaced to the TUI).
- **`shell.env`**: ensures the agent's shell starts in the project root so `ctx` commands resolve to the right project.
- **`experimental.session.compacting`**: pushes `ctx system bootstrap` output into the compaction context so the agent keeps breadcrumbs back to `.context/`.

The plugin is a single file with no runtime dependencies; no `bun install`
needed. OpenCode loads it automatically on launch.

### Context Updates

```bash
# Get AI-optimized context packet
ctx agent

# Check context health
ctx status
```

---

## Pi

Pi (earendil-works/pi, pi.dev) is a self-extensible coding-agent CLI.
By design it has no built-in MCP, so `ctx` integrates via a thin
TypeScript extension, Agent-Skills-standard skills, and `AGENTS.md`
instructions. All logic stays in the `ctx` binary via `ctx system`
subcommands.

### Setup

```bash
# Generate Pi extension, skills, and AGENTS.md
ctx setup pi --write

# Initialize context
ctx init
```

### What Gets Created

| File | Purpose |
|------|---------|
| `.pi/extensions/ctx.ts` | Lifecycle extension (hooks to `ctx system`) |
| `AGENTS.md` | Agent instructions (read natively) |
| `.pi/skills/ctx-*/SKILL.md` | `ctx` skills (invoke as `/skill:ctx-*`) |

Project-local `.pi/` files load only after the project is trusted;
the first `pi` launch in a new project prompts for trust.

### How It Works

The extension wires Pi lifecycle events to `ctx system`:

- **`session_start`**: warms the ctx agent packet off the prompt path.
- **`before_agent_start`**: injects the packet as a persistent message
  when no ctx message exists in the live context (summary + kept tail
  + post-compaction entries): fresh sessions inject on the first turn,
  and a packet folded into the compaction summary is re-injected as a
  breadcrumb into `.context/`.
- **`tool_result` (bash, on `git commit`, not `isError`)**: runs
  `ctx system post-commit`.
- **`tool_result` (edit/write, not `isError`)**: runs
  `ctx system check-task-completion`.
- **`agent_settled`**: runs `ctx system check-persistence`.

The extension is a single file with no runtime dependencies; no
`npm install` needed. Pi loads it automatically on launch.

---

## Windsurf IDE

Windsurf supports custom instructions and file-based context.

### Setup

```bash
# Generate Windsurf configuration
ctx setup windsurf

# Initialize context
ctx init
```

### Configuration

Add to Windsurf settings:

```text
// Split to multiple lines for readability
{
  "ai.customInstructions": "Always read .context/CONSTITUTION.md first. 
  Check .context/TASKS.md for current work. 
  Follow patterns in .context/CONVENTIONS.md."
}
```

### Usage

Context files appear in the file tree. Reference them when chatting:

* "What's in our task list?" → AI reads `.context/TASKS.md`
* "What convention do we use for naming?" → AI reads `.context/CONVENTIONS.md`

---

## Generic Integration

For any AI tool that can read files, use these patterns:

### Manual Context Loading

```bash
# Get full context
ctx load

# Get AI-optimized packet
ctx agent --budget 8000

# Get specific file
cat .context/TASKS.md
```

### System Prompt Template

```
You are working on a project with persistent context in .context/

Before responding:
1. Read .context/CONSTITUTION.md - NEVER violate these rules
2. Check .context/TASKS.md for current work
3. Follow .context/CONVENTIONS.md patterns
4. Reference .context/DECISIONS.md for architectural choices

When you learn something new, note it for .context/LEARNINGS.md
When you make a decision, document it for .context/DECISIONS.md
```

### Automated Updates

If your AI tool outputs to a log, use `ctx watch`:

```bash
# Watch log file for context-update commands
your-ai-tool 2>&1 | tee /tmp/ai.log &
ctx watch --log /tmp/ai.log
```

The AI can emit updates like:

```xml
<context-update type="complete">implement caching</context-update>
<context-update type="learning"
  context="Implementing caching layer"
  lesson="Important thing learned today"
  application="Apply this insight going forward"
>Caching Insight</context-update>
```

---

## Context Update Commands

The `ctx watch` command parses update commands from AI output. Use this format:

```text
<context-update type="TYPE" [attributes]>Content</context-update>
```

### Supported Types

| Type         | Target File    | Required Attributes                           |
|--------------|----------------|-----------------------------------------------|
| `task`       | TASKS.md       | None                                          |
| `decision`   | DECISIONS.md   | `context`, `rationale`, `consequence`         |
| `learning`   | LEARNINGS.md   | `context`, `lesson`, `application`            |
| `convention` | CONVENTIONS.md | None                                          |
| `complete`   | TASKS.md       | None                                          |

### Simple Format (Tasks, Conventions, Complete)

```xml
<context-update type="task">Implement rate limiting</context-update>
<context-update type="convention">Use kebab-case for files</context-update>
<context-update type="complete">rate limiting</context-update>
```

### Structured Format (*Learnings, Decisions*)

Learnings and decisions support structured attributes for better documentation:

**Learning with full structure**:

```xml
<context-update type="learning"
  context="Debugging Claude Code hooks"
  lesson="Hooks receive JSON via stdin, not environment variables"
  application="Parse JSON stdin with the host language (Go, Python, etc.): no jq needed"
>Hook Input Format</context-update>
```

**Decision with full structure**:

```xml
<context-update type="decision"
  context="Need a caching layer for API responses"
  rationale="Redis is fast, well-supported, and team has experience"
  consequence="Must provision Redis infrastructure; team training on Redis patterns"
>Use Redis for caching</context-update>
```

Learnings require: `context`, `lesson`, `application` attributes.
Decisions require: `context`, `rationale`, `consequence` attributes.
Updates missing required attributes are rejected with an error.

## Further Reading

* [Skills That Fight the Platform](../blog/2026-02-04-skills-that-fight-the-platform.md): 
  Common pitfalls in skill design that work against the host tool
* [The Anatomy of a Skill That Works](../blog/2026-02-07-the-anatomy-of-a-skill-that-works.md):
  What makes a skill reliable: the E/A/R framework and quality gates
