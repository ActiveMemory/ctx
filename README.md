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

Read the full **[Manifesto](MANIFESTO.md)** | **[ctx.ist](https://ctx.ist/)**

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
| [Manifesto](https://ctx.ist/)                                     | Philosophy: creation, context, verification         |
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
# (git is required: ctx refuses to operate without .git/.
# Run `git init` first if the project does not have a repo yet.)
ctx init

# Run subsequent commands from the project root. ctx always
# reads $PWD/.context/; there is no env-var or walk-up.
ctx status

# Get an AI-ready context packet
ctx agent --budget 4000

# Add tasks, decisions, learnings
ctx task add "Implement user authentication"
ctx decision add "Use PostgreSQL for primary database" \
  --context "Need a reliable database for production workloads" \
  --rationale "PostgreSQL offers ACID compliance, JSON support, and team familiarity" \
  --consequence "Team needs PostgreSQL training; must set up replication"
ctx learning add "Mock functions must be hoisted in Jest"
```

### Knowledge-base workflow (Phase KB)

For knowledge-shaped work (research projects, vendor-spec analysis,
post-incident reviews), `ctx init` also lays down an editorial
pipeline distinct from the code-development surface above:

```bash
# Scaffold a folder-shaped topic page
ctx kb topic new "Cursor Hooks"

# Run an editorial ingest pass (in your AI assistant)
/ctx-kb-ingest ./inputs/2026-05-15-call.md "cursor hooks"

# Ask grounded questions, refuses to web-jump
/ctx-kb-ask "does the kb say hooks fire async?"

# Per-session handover; folds postdated closeouts automatically
ctx handover write "First session" \
  --summary "Drafted topic-page; minted EV-018..EV-024" \
  --next "Re-ingest the v1.1 release notes URL when available"
```

See the
[Build a Knowledge Base recipe](https://ctx.ist/recipes/build-a-knowledge-base/)
for the full workflow, including the pass-mode contract,
source-coverage state-machine ledger, and the closeout/fold
mechanism.

`ctx` reads `$PWD/.context/` — run commands from the project root
(the directory that holds both `.git/` and `.context/`). A small
allowlist (`init`, `version`, `help`, `system bootstrap`, `doctor`,
`guide`, `why`, `config switch/status`, `hub *`) runs without
`.context/` present; every other command exits with a next-step
hint when it is missing.

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
| [CLI Reference](https://ctx.ist/cli/)           | All commands and options               |
| [Context Files](https://ctx.ist/home/context-files/) | File formats and structure       |
| [Integrations](https://ctx.ist/operations/integrations/) | Claude Code, Cursor, Aider setup |
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
