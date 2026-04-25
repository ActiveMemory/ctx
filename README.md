![ctx](assets/ctx-banner.png)

## `ctx` (*Context*)

> **`ctx` is a system, not a prompt.**

A lightweight, file-based system that enables AI coding assistants to persist,
structure, and rehydrate project context across sessions.

`ctx` works with **any AI tool** that can read files; no model or 
vendor lock-in.

**Full documentation: [ctx.ist](https://ctx.ist)**

## The `ctx` Manifesto

> **Creation, not code. Context, not prompts. Verification, not vibes.**

`ctx` is infrastructure for preserving intent under scale. Without durable
context, intelligence resets. With `ctx`, creation compounds.

Read the full **[Manifesto](MANIFESTO.md)** | **[ctx.ist/manifesto](https://ctx.ist/manifesto/)**

## The Thesis

> **Context as State: A Persistence Layer for Human-AI Cognition**

AI-assisted development systems assemble context at prompt time using heuristic
retrieval from mutable sources. These approaches optimize relevance at the moment
of generation but provide no mechanism for persistence, verification, or
accumulated learning across sessions. `ctx` treats context as deterministic state.

Read the full **[Thesis](https://ctx.ist/thesis/)**

## Core Documents

| Document                                                          | Context                                             |
|-------------------------------------------------------------------|-----------------------------------------------------|
| [Manifesto](https://ctx.ist/manifesto/)                           | Philosophy: creation, context, verification         |
| [The Thesis](https://ctx.ist/thesis/)                             | Whitepaper: context as deterministic state          |
| [Design Invariants](https://ctx.ist/reference/design-invariants/) | System properties that must always hold             |
| [Tool Comparison](https://ctx.ist/reference/comparison/)          | How `ctx` differs from .cursorrules, Aider, Copilot |
| [`ctx` Blog](https://ctx.ist/blog/)                               | Deep dives, architecture notes, learnings           |

## The Problem

Most LLM-driven development fails not because models are weak: They fail because
**context is ephemeral**. Every new session starts near zero:

* You re-explain architecture
* The AI repeats past mistakes
* Decisions get rediscovered instead of remembered

## The Solution

`ctx` treats context as infrastructure:

* **Persist**: Tasks, decisions, learnings survive session boundaries
* **Reuse**: Decisions don't get rediscovered; lessons stay learned
* **Align**: Context structure mirrors how engineers actually think
* **Integrate**: Works with any AI tool that can read files

Here's what that looks like in practice:

```text
❯ "Do you remember?"

● Yes. The PreToolUse hook runs ctx agent, and CLAUDE.md tells me to
  read the context files. I have context.

❯ "What have we been working on recently?"

● Yes. I can run ctx journal source and review recent activity:
    - 2025-01-20: The meta-experiment that started it all
    - 2025-01-21: The ctx rename + Claude hooks session
```

That's the whole point: **Temporal continuity across sessions**.

## Installation

Download pre-built binaries from the
[releases page](https://github.com/ActiveMemory/ctx/releases), or build from
source:

```bash
git clone https://github.com/ActiveMemory/ctx.git
cd ctx
CGO_ENABLED=0 go build -o ctx ./cmd/ctx
sudo mv ctx /usr/local/bin/
```

See [installation docs](https://ctx.ist/#installation) for platform-specific
instructions.

## Quick Start

```bash
# Initialize context directory in your project
ctx init

# Activate it for the current shell (binds CTX_DIR). Required
# before every other command: ctx no longer walks up the
# filesystem looking for .context/.
eval "$(ctx activate)"

# Check context status
ctx status

# Get an AI-ready context packet
ctx agent --budget 4000

# Add tasks, decisions, learnings
ctx add task "Implement user authentication"
ctx add decision "Use PostgreSQL for primary database" \
  --context "Need a reliable database for production workloads" \
  --rationale "PostgreSQL offers ACID compliance, JSON support, and team familiarity" \
  --consequence "Team needs PostgreSQL training; must set up replication"
ctx add learning "Mock functions must be hoisted in Jest"
```

`ctx activate` emits `export CTX_DIR=...` for your shell; one-shot
callers can prefix the binding inline as `CTX_DIR=<abs-path> ctx ...`.
The value must be an absolute path with `.context` as its basename;
relative paths and other names are rejected on first use. A small
allowlist (`init`, `activate`, `deactivate`, `version`, `help`,
`system bootstrap`, `doctor`, `guide`, `why`, `config switch/status`,
`hub *`) runs without CTX_DIR declared; every other command exits
with a next-step hint when it is unset.

## Documentation

This README is a map, not the territory. The full documentation
lives at **[ctx.ist](https://ctx.ist)** and carries the recipes,
runbooks, threat model, and design rationale that this file
intentionally doesn't try to fit. If you're past install and
wondering "*how do I actually use this in a real session,*" the
recipes are the right next stop.

| Guide                                           | Description                            |
|-------------------------------------------------|----------------------------------------|
| [Getting Started](https://ctx.ist)              | Installation, quick start, first steps |
| [Recipes](https://ctx.ist/recipes/)             | Practical workflow guides              |
| [CLI Reference](https://ctx.ist/cli-reference/) | All commands and options               |
| [Context Files](https://ctx.ist/context-files/) | File formats and structure             |
| [Integrations](https://ctx.ist/integrations/)   | Claude Code, Cursor, Aider setup       |
| [Operations](https://ctx.ist/operations/)       | Runbooks, day-to-day, hub deployment   |
| [Security](https://ctx.ist/security/)           | Trust model, audit trail, permissions  |

## Contributing

Contributions welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

All commits must be signed off (`git commit -s`) to certify the
[DCO](CONTRIBUTING_DCO.md).

## Community

**Open source is better together**.

Join the community to ask questions, share feedback, and connect with
other users:

[Join the `ctx` Discord](https://ctx.ist/discord)

## License

[Apache 2.0](LICENSE)
