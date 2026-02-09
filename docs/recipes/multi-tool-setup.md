---
title: "Setting Up ctx Across AI Tools"
icon: lucide/wrench
---

![ctx](../images/ctx-banner.png)

## Problem

You have installed ctx and want to set it up with your AI coding assistant so
that context persists across sessions. Different tools have different
integration depths. For example: 

* Claude Code supports native hooks that load and save context automatically
* Cursor injects context via its system prompt
* Aider reads context files through its `--read` flag

This recipe walks through the complete setup for each tool, from initialization
through verification, so you end up with a working memory layer regardless of
which AI tool you use.

## Commands and Skills Used

| Command/Skill       | Role in this workflow                                        |
|---------------------|--------------------------------------------------------------|
| `ctx init`          | Create `.context/` directory, templates, and tool hooks      |
| `ctx hook`          | Generate integration configuration for a specific AI tool    |
| `ctx agent`         | Print a token-budgeted context packet for AI consumption     |
| `ctx load`          | Output assembled context in read order (for manual pasting)  |
| `ctx watch`         | Auto-apply context updates from AI output (non-native tools) |
| `ctx completion`    | Generate shell autocompletion for bash, zsh, or fish         |
| `ctx session parse` | Convert JSONL transcripts to readable markdown               |

## The Workflow

### Step 1: Initialize ctx

Run `ctx init` in your project root. This creates the `.context/` directory
with all template files and, if Claude Code is detected, generates hooks and
Agent Skills automatically.

```bash
cd your-project
ctx init
````

This produces the following structure:

```
.context/
  CONSTITUTION.md     # Hard rules the AI must never violate
  TASKS.md            # Current and planned work
  CONVENTIONS.md      # Code patterns and standards
  ARCHITECTURE.md     # System overview
  DECISIONS.md        # Architectural decisions with rationale
  LEARNINGS.md        # Lessons learned, gotchas, tips
  GLOSSARY.md         # Domain terms and abbreviations
  AGENT_PLAYBOOK.md   # How AI tools should use this system
  sessions/           # Session snapshots

.claude/              # Claude Code integration (auto-generated)
  hooks/              # Auto-save and enforcement scripts
  skills/             # ctx Agent Skills (agentskills.io spec)
  settings.local.json # Hook configuration
```

If you only need the core files (useful for lightweight setups with Cursor or
Copilot), use the `--minimal` flag:

```bash
ctx init --minimal
```

This creates only `TASKS.md`, `DECISIONS.md`, and `CONSTITUTION.md`.

### Step 2: Generate Tool-Specific Hooks

If you are using a tool other than Claude Code (which is configured
automatically by `ctx init`), generate its integration configuration:

```bash
# For Cursor
ctx hook cursor

# For Aider
ctx hook aider

# For GitHub Copilot
ctx hook copilot

# For Windsurf
ctx hook windsurf
```

Each command prints the configuration you need. What you do with it depends on
the tool.

!!! tip "Claude is a First-Class Citizen"
    You don't need any extra steps to integrate with Claude Code.

    `ctx init` already wrote `.claude/settings.local.json` with
    `PreToolUse` and `SessionEnd` hooks.

    The `PreToolUse` hook runs  
    `ctx agent --budget 4000 --session $PPID` on every tool call  
    (*with a 10-minute cooldown so it only fires once per window*).

    The `SessionEnd` hook saves a snapshot to `.context/sessions/`.

**Cursor**: Add the system prompt snippet to `.cursor/settings.json`:

```json
{
  "ai.systemPrompt": "Read .context/TASKS.md and .context/CONVENTIONS.md before responding. Follow rules in .context/CONSTITUTION.md."
}
```

Context files appear in Cursor's file tree. You can also paste a context packet
directly into chat:

```bash
ctx agent --budget 4000 | xclip    # Linux
ctx agent --budget 4000 | pbcopy   # macOS
```

**Aider**: Create `.aider.conf.yml` so context files are loaded on every
session:

```yaml
read:
  - .context/CONSTITUTION.md
  - .context/TASKS.md
  - .context/CONVENTIONS.md
  - .context/DECISIONS.md
```

Then start Aider normally:

```bash
aider
```

Or specify files on the command line:

```bash
aider --read .context/TASKS.md --read .context/CONVENTIONS.md
```

### Step 3: Set Up Shell Completion

Shell completion lets you tab-complete ctx subcommands and flags, which is
especially useful while learning the CLI.

```bash
# Bash (add to ~/.bashrc)
source <(ctx completion bash)

# Zsh (add to ~/.zshrc)
source <(ctx completion zsh)

# Fish
ctx completion fish > ~/.config/fish/completions/ctx.fish
```

After sourcing, typing `ctx a<TAB>` completes to `ctx agent`, and
`ctx session <TAB>` shows `save`, `list`, `load`, and `parse`.

### Step 4: Verify the Setup Works

Start a fresh session in your AI tool and ask:

> **"Do you remember?"**

A correctly configured tool responds with specific context: current tasks from
`TASKS.md`, recent decisions, and previous session topics. It should not say
"I don't have memory" or "Let me search for files."

This question checks the *passive* side of memory. A properly set-up agent is
also **proactive**: it treats context maintenance as part of its job.

* After a debugging session, it offers to save a **learning**
* After a trade-off discussion, it asks whether to record the decision
* After completing a task, it suggests follow-up items

The "*do you remember?*" check verifies both halves: recall **and**
responsibility.

For example, after resolving a tricky bug, a proactive agent might say:

> That Redis timeout issue was subtle. Want me to save this as a learning so
> we don't hit it again?

If you see behavior like this, the setup is working end to end.

In Claude Code, you can also invoke the `/ctx-status` skill:

```text
/ctx-status
```

This prints a summary of all context files, token counts, and recent activity,
confirming that hooks are loading context.

If context is not loading, check the basics:

| Symptom                         | Fix                                                           |
|---------------------------------|---------------------------------------------------------------|
| `ctx: command not found`        | Ensure ctx is in your PATH: `which ctx`                       |
| No sessions saved (Claude Code) | Verify `.claude/settings.local.json` has `SessionEnd` hook    |
| Hook permission errors          | Run `chmod +x .claude/hooks/*.sh`                             |
| Missing sessions directory      | Run `mkdir -p .context/sessions`                              |
| Context not refreshing          | Cooldown may be active; wait 10 minutes or set `--cooldown 0` |

### Step 5: Enable Watch Mode for Non-Native Tools

Tools like Aider, Copilot, and Windsurf do not support native hooks for saving
context automatically. For these, run `ctx watch` alongside your AI tool.

Pipe the AI tool's output through `ctx watch`:

```bash
# Terminal 1: Run Aider with output logged
aider 2>&1 | tee /tmp/aider.log

# Terminal 2: Watch the log for context updates
ctx watch --log /tmp/aider.log
```

Or for any generic tool:

```bash
your-ai-tool 2>&1 | tee /tmp/ai.log &
ctx watch --log /tmp/ai.log
```

When the AI emits structured update commands, `ctx watch` parses and applies
them automatically:

```xml
<context-update type="learning"
  context="Debugging rate limiter"
  lesson="Redis MULTI/EXEC does not roll back on error"
  application="Wrap rate-limit checks in Lua scripts instead"
>Redis Transaction Behavior</context-update>
```

To preview changes without modifying files:

```bash
ctx watch --dry-run --log /tmp/ai.log
```

### Step 6: Parse Session Transcripts (Optional)

If you have JSONL transcripts from Claude Code sessions, convert them to
readable Markdown:

```bash
ctx session parse ~/.claude/projects/.../session.jsonl -o conversation.md
```

To also extract **decisions** and **learnings**:

```bash
ctx session parse transcript.jsonl --extract
```

This scans the conversation and appends relevant entries to `DECISIONS.md` and
`LEARNINGS.md`.

## Putting It Together

Here is the condensed setup for all three tools:

```bash
# -- Common (run once per project) --
cd your-project
ctx init
source <(ctx completion zsh)       # or bash/fish

# -- Claude Code (automatic, just verify) --
# Start Claude Code, then ask: "Do you remember?"

# -- Cursor --
ctx hook cursor
# Add the system prompt to .cursor/settings.json
# Paste context: ctx agent --budget 4000 | pbcopy

# -- Aider --
ctx hook aider
# Create .aider.conf.yml with read: paths
# Run watch mode alongside: ctx watch --log /tmp/aider.log

# -- Verify any tool --
# Ask your AI: "Do you remember?"
# Expect: specific tasks, decisions, session history
```

## Tips

* Start with `ctx init` (not `--minimal`) for your first project. The full
  template set gives the agent more to work with, and you can always delete
  files later.
* For Claude Code, adjust the token budget in `.claude/settings.local.json`
  as your project grows.
* The `--session $PPID` flag isolates cooldowns per Claude Code process, so
  parallel sessions do not suppress each other.
* Commit your `.context/` directory to version control. Several ctx features
  (journals, changelogs, blog generation) rely on git history.
* For Cursor and Copilot, keep `CONVENTIONS.md` visible. These tools treat
  open files as higher-priority context.
* Run `ctx drift` periodically to catch stale references before they confuse
  the agent.
* Once set up, the agent becomes proactive. You should not need to tell it to
  persist context: **it will ask you**.

## See Also

* [The Complete Session](session-lifecycle.md): full session lifecycle recipe
* [CLI Reference](../cli-reference.md): all commands and flags
* [Integrations](../integrations.md): detailed per-tool integration docs
