---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "ctx for Codex"
icon: lucide/terminal
---

![ctx](../images/ctx-banner.png)

## The Problem

Every Codex session starts from zero. You re-explain your architecture,
the AI repeats mistakes it made yesterday, and decisions get rediscovered
instead of remembered.

**Without `ctx`:**

```
> "Add the validation middleware we discussed"

I don't have context about previous discussions. Could you describe
what validation middleware you're referring to?
```

**With `ctx`:**

```
> "Add the validation middleware we discussed"

Yes. From the Jan 15 session. You decided on Zod schemas at the
route level (DECISIONS.md #12), and the pattern is in
CONVENTIONS.md. I'll follow the existing middleware in
src/middleware/auth.ts as a reference.
```

That's the whole pitch: **your AI remembers**.

## Setup

Install the `ctx` binary first ([installation docs](getting-started.md#installation)),
then pick **one** of the two routes below. Both deliver the same hooks,
the same skills, and the same MCP server; they differ only in where the
files live.

| Route | Files live in | Best for |
|-------|---------------|----------|
| [Plugin](#route-a-the-ctx-plugin) | Codex's plugin cache (`$CODEX_HOME/plugins/`) | One install, every project |
| [Project-local](#route-b-project-local-files) | Your repository (`.codex/`, `.agents/`) | Teams, CI, `codex exec` |

!!! warning "Pick One Route, Not Both"
    Codex loads every matching hook from every source. If the plugin is
    enabled **and** the project has `.codex/hooks.json`, each hook runs
    twice. `ctx setup codex --write` detects an enabled plugin and skips
    hooks, MCP, and skills (it still deploys `AGENTS.md`), but a plugin
    installed *after* a project-local deploy is not detected by anything.

### Route A: The `ctx` Plugin

Register the `ctx` marketplace and install the plugin:

```bash
codex plugin marketplace add ActiveMemory/ctx
codex plugin add ctx@activememory-ctx
```

Working from a local checkout of the `ctx` repository? The Makefile wraps
both commands (it registers the checkout itself as a local marketplace):

```bash
make codex-plugin-install
```

Then initialize your project:

```bash
cd your-project
ctx init
```

The installed copy lives under
`$CODEX_HOME/plugins/cache/activememory-ctx/ctx/<version>/` (the version
segment is `local` for a local marketplace). Codex records the enabled
state in `~/.codex/config.toml`:

```toml
[plugins."ctx@activememory-ctx"]
enabled = true
```

### Route B: Project-Local Files

From your project root:

```bash
ctx setup codex --write && ctx init
```

`ctx setup codex` without `--write` prints the integration overview and
what it detected (Codex binary on PATH, plugin installed, plugin enabled)
without touching anything.

#### What Gets Created

| File | Purpose |
|------|---------|
| `.codex/hooks.json` | Lifecycle hooks (same manifest the plugin ships) |
| `.codex/config.toml` | `[mcp_servers.ctx]` table registering the `ctx` MCP server |
| `AGENTS.md` | Agent instructions (Codex reads this natively); marker-merged into an existing file |
| `.agents/skills/ctx-*/SKILL.md` | 50 `ctx` skills, invoked as `$ctx-<name>` |

Re-running `--write` is safe: existing foreign hook groups in
`.codex/hooks.json` are preserved and only the `ctx`-managed groups are
replaced; an existing `[mcp_servers.ctx]` table is left alone; a
`SKILL.md` that is not `ctx`-managed is skipped with a warning instead of
overwritten. An unparseable `.codex/hooks.json` is left untouched (with a
warning) while the other steps still run.

!!! note "Project Layers Need a Trusted Project"
    Codex only loads `.codex/hooks.json` and `.codex/config.toml` for
    **trusted** projects. If Codex has not asked you to trust the
    directory yet, add it to `~/.codex/config.toml`:

    ```toml
    [projects."/absolute/path/to/your-project"]
    trust_level = "trusted"
    ```

    The entry must live in the real `~/.codex/config.toml`: a
    `-c 'projects."...".trust_level="trusted"'` command-line
    override does **not** unlock project layers (verified against
    Codex 0.148).

### Trust the Hooks (Both Routes)

Codex refuses to run hooks it has not been told to trust. After either
route, start `codex` in the project, run `/hooks`, and trust the `ctx`
entries. Until you do, nothing fires and no context is injected.

For a single non-interactive run (CI, smoke tests) you can skip the
review with `codex exec --dangerously-bypass-hook-trust`; the flag
applies to that invocation only.

## What Happens Automatically

Once the hooks are trusted, `ctx` is wired into Codex's lifecycle. Every
hook command starts with `cd "$(git rev-parse --show-toplevel)" &&`
because Codex runs hooks with the session cwd and `ctx` reads
`$PWD/.context/`; the anchor makes a subdirectory cwd harmless.

| Codex event | Matcher | What runs | What it does |
|-------------|---------|-----------|--------------|
| `SessionStart` | all sources | `ctx agent --budget 8000` | Injects the context packet as developer context. Re-fires on `compact`, so context survives compaction |
| `PreToolUse` | `.*` | `ctx system context-load-gate` | Autoload gate on first tool use |
| `PreToolUse` | `Bash` | `ctx system block-non-path-ctx` | Blocks `./ctx` and `go run` invocations; forces the `$PATH` install |
| `PreToolUse` | `Bash` | `ctx system qa-reminder` | Lint/test reminder before a commit |
| `PreToolUse` | `update_plan` | `ctx system specs-nudge` | Nudges toward project specs when Codex plans (`update_plan` is Codex's planning tool) |
| `PostToolUse` | `Bash` | `ctx system post-commit` | Context-capture and QA nudge after `git commit` |
| `PostToolUse` | `apply_patch\|Edit\|Write` | `ctx system check-task-completion` | Detects silently completed tasks after a file edit |
| `UserPromptSubmit` | | 12 `ctx system check-*` hooks plus `heartbeat` | Context-size, ceremony, persistence, journal, reminder, version, resource, knowledge, map-staleness, memory-drift, freshness, and skill-discovery nudges; the same list as Claude Code |
| `SessionEnd` | | `ctx journal import --all -y` (`timeout: 3`) | Imports the session into `.context/journal/` |

`PermissionRequest`, `PreCompact`, `PostCompact`, `SubagentStart`,
`SubagentStop`, and `Stop` are not wired: no `ctx` behavior maps onto
them today.

### What Is Different from Claude Code

- The context packet arrives at **`SessionStart`** rather than on the
  first tool call. Codex ignores plain text on `PreToolUse`, and
  `SessionStart` output becomes developer context directly.
- The planning matcher is `update_plan` (Claude Code: `EnterPlanMode`),
  and file edits arrive as `apply_patch`.
- Codex caps `SessionEnd` hooks at **3 seconds**. The journal import is
  incremental (one `stat` per already-imported session), so it normally
  finishes well inside that; if it does not, Codex reports a hook failure
  and the next `check-journal` nudge or the next session end picks up
  where it left off.

## Skills

The plugin and the project-local route both ship the `ctx` skills as
Codex skills. Invoke them with a `$` prefix:

| Skill | When to use |
|-------|-------------|
| `$ctx-agent` | Load the full context packet. Use when context feels stale. |
| `$ctx-remember` | "Do you remember?"; reads tasks, decisions, learnings, and recent journal entries. Returns a structured readback. |
| `$ctx-status` | Context summary at a glance: file count, token estimate, recent activity. |
| `$ctx-wrap-up` | End-of-session ceremony. Captures learnings, decisions, conventions, and outstanding tasks to `.context/` files. |
| `$ctx-commit` | Commit with integrated context capture. |

The Codex skill set is generated from the Claude Code skills
(`hack/sync-codex-skills.sh` strips the Claude-only `allowed-tools:`
frontmatter). Four skills are Claude Code-only and are not shipped:
`ctx-permission-sanitize` (audits `.claude/settings.local.json`),
`ctx-plan-import` (reads `~/.claude/plans/`), `ctx-dream` (headless
`claude -p` cron), and `ctx-skill-create` (authors Claude Code skills).
Skill bodies that mention `/ctx-remember` and friends refer to the same
skill under its `$ctx-remember` name.

## MCP Tools

Both routes register the `ctx` MCP server (`ctx mcp serve`): the plugin
through its bundled `.mcp.json`, the project-local route through
`[mcp_servers.ctx]` in `.codex/config.toml`. The server exposes these
tools to the agent:

| Tool | Purpose |
|------|---------|
| `ctx_add` | Add a task, decision, learning, or convention |
| `ctx_complete` | Mark a task done by number or text match |
| `ctx_search` | Full-text search across all `.context/` files |
| `ctx_next` | Suggest the next pending task by priority |
| `ctx_drift` | Detect stale context: dead paths, missing files |
| `ctx_compact` | Archive completed tasks, clean empty sections |
| `ctx_remind` | List pending session-scoped reminders |
| `ctx_status` | Context health: file count, token estimate |
| `ctx_steering_get` | Retrieve steering files applicable to the current prompt |
| `ctx_journal_source` | Query recent AI session history |
| `ctx_sessionevent` | Signal session start/end lifecycle events |
| `ctx_watch_update` | Apply structured updates to `.context/` files |
| `ctx_checktaskcompletion` | After a write, detect silently completed tasks |

You don't invoke these yourself. The agent uses them as needed.

## Session History

Codex writes a rollout transcript per session under
`$CODEX_HOME/sessions/YYYY/MM/DD/rollout-<timestamp>-<uuid>.jsonl`
(default `~/.codex/sessions`). `ctx journal import` discovers them, matches
them to the current project by the session's working directory, and
imports them with the tool id `codex`:

```bash
ctx journal import --all              # Codex sessions land alongside Claude Code ones
ctx journal source --tool codex       # list only Codex sessions
```

Developer-injected items (`<environment_context>`, `<user_instructions>`,
skill and permission preambles) are filtered out; a rollout with no real
user message is skipped. The `SessionEnd` hook runs the same import on
the way out of every session.

## Steering Files

Codex is not a steering sync target. It receives `inclusion: always`
steering files inside the `SessionStart` context packet and can fetch
`auto`/`manual` files on demand through the `ctx_steering_get` MCP tool.
With `tool: codex` in `.ctxrc`, `ctx steering sync` prints an info line
and exits 0 instead of writing anything. See
[How Claude Code and Codex Consume Steering](../cli/steering.md#how-claude-code-and-codex-consume-steering).

## Already on Claude Code?

Codex ships a `/import` command that copies Claude Code hooks, skills,
and MCP servers into Codex. It works for `ctx` too, but prefer
`ctx setup codex --write` or the plugin: the Codex manifest anchors
commands to the git root instead of `${CLAUDE_PROJECT_DIR}`, moves the
context packet to `SessionStart`, uses Codex's `update_plan` and
`apply_patch` matchers, and sets the 3-second `SessionEnd` timeout.
An imported Claude manifest carries none of that.

## Known Limitations

- **Hooks load only in trusted projects.** Project-local files are
  ignored until the directory is trusted (see above). The plugin route
  does not have this constraint.
- **`SessionEnd` is capped at 3 seconds** by Codex. The import is
  incremental, so a timeout only delays the import to the next sweep.
- **No Windows `commandWindows` override** is shipped, matching the
  Claude Code manifest; hooks require a POSIX shell with `git` on PATH.
- **Codex memories are not bridged.** `~/.codex/memories/` is generated,
  opaque state with no documented file contract, so `ctx` leaves it
  alone (unlike the Claude Code `MEMORY.md` bridge).
- **No status line.** Codex has no statusline hook.

## Refreshing the Integration

- **Plugin route:** Codex caches plugins by version. After bumping
  `VERSION` (`make sync-version` updates
  `internal/assets/codex/.codex-plugin/plugin.json` and
  `.agents/plugins/marketplace.json` together), re-run
  `codex plugin add ctx@activememory-ctx` and start a new session.
- **Project-local route:** re-run `ctx setup codex --write`. Stale
  `ctx`-managed hook groups and skills are refreshed in place; your own
  hook groups and `config.toml` content are preserved.
- **Skills drift check:** `make check-codex-skills` fails when
  `internal/assets/codex/skills/` is out of sync with the Claude Code
  skills; `make sync-codex-skills` regenerates it.

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| No context packet at session start | Hooks not trusted yet | Run `/hooks` in `codex` and trust the `ctx` entries |
| Hooks trusted but nothing fires in this project | Project not trusted, so `.codex/` layers are ignored | Add `[projects."<abs path>"] trust_level = "trusted"` to `~/.codex/config.toml` |
| Plugin install delivered Claude hooks (cache has `.claude-plugin/` but no `.codex-plugin/`, hook commands reference `CLAUDE_PROJECT_DIR`) | The marketplace source revision predates the dual-manifest Claude plugin root, and Codex fell back to the legacy `.claude-plugin/marketplace.json` | Reinstall from a current ref: the Claude plugin root now also carries `.codex-plugin/plugin.json` pointing at Codex-format hooks (`hooks/codex.json`), and Codex prefers a `.codex-plugin` manifest when both exist (verified live) — so either marketplace file yields working Codex hooks. On a stale install, `ctx setup codex --write` detects the wrong variant, warns, and deploys the project-local route anyway |
| Hooks run twice after installing the plugin | Project-local `.codex/hooks.json` and the plugin both load | Pick one route: delete the project's `.codex/hooks.json` and `.agents/skills/` (keep `AGENTS.md`) when moving to the plugin |
| Every nudge appears twice | Plugin enabled **and** `.codex/hooks.json` present | Remove one: `codex plugin remove ctx@activememory-ctx` or delete the `ctx` groups from `.codex/hooks.json` |
| `ctx: command not found` inside a hook | `ctx` not on the PATH Codex inherits | `which ctx`; install to a PATH directory (the `block-non-path-ctx` hook exists for exactly this) |
| `SessionEnd` hook reports a timeout | First import of a long backlog exceeded 3 s | Run `ctx journal import --all` once by hand; later runs are incremental |
| Sessions missing from `ctx journal source` | `CODEX_HOME` points elsewhere, or the rollout `cwd` is not this project | Check `$CODEX_HOME`; sessions match by working directory |

## Verify It Works

Start a new Codex session in the project and ask:

```
Do you remember?
```

The AI should cite specific context: current tasks, recent decisions, or
previous session topics. If it says "I don't have memory" or "Let me
check," something went wrong; confirm the hooks are trusted and
`.context/` has files in it.

## What's Next

- [Your First Session](first-session.md): step-by-step walkthrough from
  `ctx init` to verified recall.
- [Common Workflows](common-workflows.md): day-to-day commands for
  tracking context, checking health, and browsing history.
- [Context Files](context-files.md): what lives in `.context/` and how
  each file is used.
- [AI Tools](../operations/integrations.md#openai-codex): the hook
  manifest, the event table, and the drift checks that keep this page
  honest.
