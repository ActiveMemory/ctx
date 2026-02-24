![ctx](assets/ctx-banner.png)

## `ctx` (*Context*)

> **`ctx` is a system, not a prompt.**

A lightweight, file-based system that enables AI coding assistants to persist,
structure, and rehydrate project context across sessions.

`ctx` works with **any AI tool** that can read files; no model or 
vendor lock-in.

**Full documentation: [ctx.ist](https://ctx.ist)**

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

● Yes. I can run ctx recall list and review recent activity:
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

# Check context status
ctx status

# Get an AI-ready context packet
ctx agent --budget 4000

# Add tasks, decisions, learnings
ctx add task "Implement user authentication"
ctx add decision "Use PostgreSQL for primary database" \
  --context "Need a reliable database for production workloads" \
  --rationale "PostgreSQL offers ACID compliance, JSON support, and team familiarity" \
  --consequences "Team needs PostgreSQL training; must set up replication"
ctx add learning "Mock functions must be hoisted in Jest"
```

## Documentation

| Guide                                           | Description                            |
|-------------------------------------------------|----------------------------------------|
| [Getting Started](https://ctx.ist)              | Installation, quick start, first steps |
| [CLI Reference](https://ctx.ist/cli-reference/) | All commands and options               |
| [Context Files](https://ctx.ist/context-files/) | File formats and structure             |
| [Integrations](https://ctx.ist/integrations/)   | Claude Code, Cursor, Aider setup       |
| [Ralph Loop](https://ctx.ist/ralph-loop/)       | Autonomous AI development workflows    |

## Manifesto

> **Creation, not code. Context, not prompts. Verification, not vibes.**

`ctx` is infrastructure for preserving intent under scale. Without durable
context, intelligence resets. With `ctx`, creation compounds.

Read the full **[Manifesto](MANIFESTO.md)** | **[ctx.ist/manifesto](https://ctx.ist/manifesto/)**

## Contributing

Contributions welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

All commits must be signed off (`git commit -s`) to certify the
[DCO](CONTRIBUTING_DCO.md).

## Community

**Open source is better together**.

Join the community to ask questions, share feedback, and connect with
other users:

[Join `#ctx` on `irc.libera.chat`](https://web.libera.chat/#ctx)

## License

[Apache 2.0](LICENSE)
