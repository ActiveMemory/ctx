---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Getting Started
icon: lucide/rocket
---

![ctx](images/ctx-banner.png)

## `ctx`

`ctx` (*Context*) is a file-based system that enables AI coding assistants to
persist project knowledge across sessions. Instead of re-explaining your
codebase every time, context files let AI tools remember decisions,
conventions, and learnings:

* A session is interactive.
* `ctx` enables **cognitive continuity**.
* **Cognitive continuity** enables durable, *symbiotic-like* human–AI workflows.

!!! quote "The `ctx` Manifesto"
    **Creation, not code. Context, not prompts. Verification, not vibes.**

    Without durable context, intelligence resets.
    With `ctx`, creation compounds.

    **[Read the Manifesto →](https://ctx.ist/)**

## Community

**Open source is better together**.

<!-- the long line is required for zensical to render block quote -->

!!! tip "Help `ctx` Change How AI Remembers"
    **If the idea behind `ctx` resonates, a star helps it reach engineers 
    who run into context drift every day.**

    → https://github.com/ActiveMemory/ctx

    `ctx` is free and open source software, and **contributions are always
    welcome** and appreciated.

Join the community to ask questions, share feedback, and connect with
other users:

- [:fontawesome-brands-stack-exchange: **IRC**](https://web.libera.chat/#ctx): 
   join `#ctx` on `irc.libera.chat`
- [:fontawesome-brands-github: **GitHub**](https://github.com/ActiveMemory/ctx):
  Star the repo, report issues, contribute

## Why?

Most AI-driven development fails not because models are weak—they fail because 
**context is ephemeral**. Every new session starts near zero:

* You re-explain architecture
* The AI repeats past mistakes
* Decisions get rediscovered instead of remembered

`ctx` solves this by treating context as **infrastructure**: 
files that version with your code and persist across sessions.

## Installation

### Build from Source (*Recommended*)

Requires [Go 1.25+](https://go.dev/):

```bash
git clone https://github.com/ActiveMemory/ctx.git
cd ctx
make build
sudo make install
# or:
# sudo mv ctx /usr/local/bin/
```

Building from source gives you the latest features and bug fixes. 

Since `ctx` is predominantly a developer tool, this is the 
**recommended approach**: 

You get the freshest code and can inspect what you are installing.

### Binary Downloads

Pre-built binaries are available from the
[releases page](https://github.com/ActiveMemory/ctx/releases) if you prefer
not to build from source.

=== "Linux (x86_64)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.3.0/ctx-0.3.0-linux-amd64
    chmod +x ctx-0.3.0-linux-amd64
    sudo mv ctx-0.3.0-linux-amd64 /usr/local/bin/ctx
    ```

=== "Linux (ARM64)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.3.0/ctx-0.3.0-linux-arm64
    chmod +x ctx-0.3.0-linux-arm64
    sudo mv ctx-0.3.0-linux-arm64 /usr/local/bin/ctx
    ```

=== "macOS (Apple Silicon)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.3.0/ctx-0.3.0-darwin-arm64
    chmod +x ctx-0.3.0-darwin-arm64
    sudo mv ctx-0.3.0-darwin-arm64 /usr/local/bin/ctx
    ```

=== "macOS (Intel)"

    ```bash
    curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.3.0/ctx-0.3.0-darwin-amd64
    chmod +x ctx-0.3.0-darwin-amd64
    sudo mv ctx-0.3.0-darwin-amd64 /usr/local/bin/ctx
    ```

=== "Windows"

    Download `ctx-0.3.0-windows-amd64.exe` from the releases page and add it to your `PATH`.

#### Verifying Checksums

Each binary has a corresponding `.sha256` checksum file. To verify your download:

```bash
# Download the checksum file
curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.3.0/ctx-0.3.0-linux-amd64.sha256

# Verify the binary
sha256sum -c ctx-0.3.0-linux-amd64.sha256
```

On macOS, use `shasum -a 256 -c` instead of `sha256sum -c`.

Verify installation:

```bash
ctx --version
```

### Version Control (Strongly Recommended)

`ctx` does not require git, but using version control with your `.context/`
directory is strongly recommended. AI sessions occasionally modify or
overwrite context files inadvertently. With git, the AI can check history
and restore lost content — without it, the data is gone. Several `ctx`
features (journal changelog, blog generation) also use git history directly.

## Quick Start

### 1. Initialize Context

```bash
cd your-project
ctx init
```

This creates a `.context/` directory with template files and configures 
AI tool hooks (*for Claude Code*).

### 2. Check Status

```bash
ctx status
```

Shows context summary: files present, token estimate, and recent activity.

### 3. Start Using with AI

With Claude Code, context loads automatically via hooks. For other tools,
paste the output of:

```bash
ctx agent --budget 8000
```

### 4. Verify It Works

Ask your AI: **"Do you remember?"**

It should cite specific context: current tasks, recent decisions,
or previous session topics.

## Your First Session

Here's what a complete first session looks like, from initialization to
the moment your AI cites your project context back to you.

### Step 1: Initialize Your Project

```bash
cd your-project
ctx init
```

```
Context initialized in .context/

  ✓ CONSTITUTION.md
  ✓ TASKS.md
  ✓ DECISIONS.md
  ✓ LEARNINGS.md
  ✓ CONVENTIONS.md
  ✓ ARCHITECTURE.md
  ✓ GLOSSARY.md
  ✓ AGENT_PLAYBOOK.md

Creating project root files...
  ✓ PROMPT.md
  ✓ IMPLEMENTATION_PLAN.md

Setting up Claude Code integration...
  ✓ Claude hooks
  ✓ .scratchpad.key

Next steps:
  1. Edit .context/TASKS.md to add your current tasks
  2. Run 'ctx status' to see context summary
  3. Run 'ctx agent' to get AI-ready context packet
```

This created your `.context/` directory with template files. If you're using
Claude Code, it also configured hooks so context loads automatically.

### Step 2: Populate Your Context

Add a task and a decision — these are the entries your AI will remember:

```bash
ctx add task "Implement user authentication"
```

```
✓ Added to TASKS.md
```

```bash
ctx add decision "Use PostgreSQL for primary database" \
  --context "Need a reliable database for production" \
  --rationale "PostgreSQL offers ACID compliance and JSON support" \
  --consequences "Team needs PostgreSQL training"
```

```
✓ Added to DECISIONS.md
```

These entries are what the AI will recall in future sessions. You don't need
to populate everything now — context grows naturally as you work.

### Step 3: Check Your Context

```bash
ctx status
```

```
Context Status
====================

Context Directory: .context/
Total Files: 8
Token Estimate: 1,247 tokens

Files:
  ✓ CONSTITUTION.md (loaded)
  ✓ TASKS.md (1 items)
  ✓ DECISIONS.md (1 items)
  ○ LEARNINGS.md (empty)
  ✓ CONVENTIONS.md (loaded)
  ✓ ARCHITECTURE.md (loaded)
  ✓ GLOSSARY.md (loaded)
  ✓ AGENT_PLAYBOOK.md (loaded)

Recent Activity:
  - TASKS.md modified 2 minutes ago
  - DECISIONS.md modified 1 minute ago
```

Notice the **token estimate** — this is how much context your AI will load.
The `○` next to LEARNINGS.md means it's still empty; it will fill in as you
capture lessons during development.

### Step 4: Start an AI Session

With **Claude Code**, context loads automatically via hooks — just start a
session and it's there. For other tools, generate a context packet:

```bash
ctx agent --budget 8000
```

```markdown
# Context Packet
Generated: 2026-02-14T15:30:45Z | Budget: 8000 tokens | Used: 1247

## Read These Files (in order)
1. .context/CONSTITUTION.md
2. .context/TASKS.md
3. .context/DECISIONS.md
...

## Current Tasks
- [ ] Implement user authentication

## Recent Decisions
- Use PostgreSQL for primary database
...
```

Paste this output into your AI tool's system prompt or conversation start.

### Step 5: Verify It Works

Ask your AI: **"What are our current tasks?"**

A working setup produces a response like:

> Based on the project context, you have one active task:
>
> - **Implement user authentication** (pending)
>
> There's also a recent architectural decision to **use PostgreSQL for
> the primary database**, chosen for its ACID compliance and JSON support.
> Want me to start on the authentication task?

That's the success moment — the AI is citing your exact context entries
from Step 2, not hallucinating or asking you to re-explain.

## What Gets Created

```
.context/
├── CONSTITUTION.md     # Hard rules — NEVER violate these
├── TASKS.md            # Current and planned work
├── CONVENTIONS.md      # Project patterns and standards
├── ARCHITECTURE.md     # System overview
├── DECISIONS.md        # Architectural decisions with rationale
├── LEARNINGS.md        # Lessons learned, gotchas, tips
├── GLOSSARY.md         # Domain terms and abbreviations
├── AGENT_PLAYBOOK.md   # How AI tools should use this
└── sessions/           # Session snapshots

.claude/                # Claude Code integration (if detected)
├── hooks/              # Lifecycle hooks (enforcement, coaching, cleanup)
├── skills/             # ctx Agent Skills (agentskills.io spec)
└── settings.local.json # Hook configuration
```

See [Context Files](context-files.md) for detailed documentation of each file.

## Common Workflows

### Track Context

```bash
# Add a task
ctx add task "Implement user authentication"

# Record a decision (full ADR fields required)
ctx add decision "Use PostgreSQL for primary database" \
  --context "Need a reliable database for production" \
  --rationale "PostgreSQL offers ACID compliance and JSON support" \
  --consequences "Team needs PostgreSQL training"

# Note a learning
ctx add learning "Mock functions must be hoisted in Jest" \
  --context "Tests failed with undefined mock errors" \
  --lesson "Jest hoists mock calls to top of file" \
  --application "Place jest.mock() before imports"

# Mark task complete
ctx complete "user auth"
```

### Check Context Health

```bash
# Detect stale paths, missing files, potential secrets
ctx drift

# See full context summary
ctx status
```

### Browse Session History

Export AI session transcripts to a browsable journal site:

```bash
# Export all sessions to .context/journal/
ctx recall export --all

# Generate and serve the journal site
ctx journal site --serve
```

Then open [http://localhost:8000](http://localhost:8000).

To update the journal after new sessions, run the same two commands
again; `recall export` preserves existing YAML frontmatter and only
updates conversation content.

See [Session Journal](session-journal.md) for the full pipeline
including enrichment and normalization.

### Browse Session History

```bash
# List recent sessions
ctx recall list --limit 5

# Export sessions to browsable journal
ctx recall export --all
```

### Run an Autonomous Loop

Generate a script that iterates an AI agent until a completion
signal is detected:

```bash
ctx loop
chmod +x loop.sh
./loop.sh
```

See [Autonomous Loops](autonomous-loop.md) for configuration
and advanced usage.

## Next Steps

* [Prompting Guide](prompting-guide.md) — Effective prompts for AI sessions
* [CLI Reference](cli-reference.md) — All commands and options
* [Context Files](context-files.md) — File formats and structure
* [Session Journal](session-journal.md) — Browse and search session history
* [Autonomous Loops](autonomous-loop.md) — Iterative AI development workflows
* [Integrations](integrations.md) — Setup for Claude Code, Cursor, Aider
* [Blog](blog/index.md) — Stories and lessons from building ctx
