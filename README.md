# Active Memory

> **Context is a system, not a prompt.**

A lightweight, file-based approach that lets AI coding assistants persist project knowledge across sessions.

## The Problem

Most AI coding assistants fail not because models are weak—they fail because context is ephemeral. Every new session starts near zero. Architectural decisions, conventions, and lessons learned decay. The user re-explains. The AI repeats mistakes. Progress feels far less cumulative than it should.

## The Solution

Active Memory treats context as infrastructure:

- **Persist** — Tasks, decisions, learnings survive session boundaries
- **Reuse** — Decisions don't get rediscovered; lessons stay learned
- **Align** — Context structure mirrors how engineers actually think
- **Integrate** — Works with any AI tool that can read files

## Installation

### Binary Downloads (Recommended)

Download pre-built binaries from the [releases page](https://github.com/josealekhine/ActiveMemory/releases).

**Linux (x86_64):**
```bash
curl -LO https://github.com/josealekhine/ActiveMemory/releases/latest/download/amem-linux-amd64
chmod +x amem-linux-amd64
sudo mv amem-linux-amd64 /usr/local/bin/amem
```

**Linux (ARM64):**
```bash
curl -LO https://github.com/josealekhine/ActiveMemory/releases/latest/download/amem-linux-arm64
chmod +x amem-linux-arm64
sudo mv amem-linux-arm64 /usr/local/bin/amem
```

**macOS (Intel):**
```bash
curl -LO https://github.com/josealekhine/ActiveMemory/releases/latest/download/amem-darwin-amd64
chmod +x amem-darwin-amd64
sudo mv amem-darwin-amd64 /usr/local/bin/amem
```

**macOS (Apple Silicon):**
```bash
curl -LO https://github.com/josealekhine/ActiveMemory/releases/latest/download/amem-darwin-arm64
chmod +x amem-darwin-arm64
sudo mv amem-darwin-arm64 /usr/local/bin/amem
```

**Windows:**

Download `amem-windows-amd64.exe` from the releases page and add it to your PATH.

### Build from Source

Requires Go 1.22+:

```bash
git clone https://github.com/josealekhine/ActiveMemory.git
cd ActiveMemory
CGO_ENABLED=0 go build -o amem ./cmd/amem
sudo mv amem /usr/local/bin/
```

## Quick Start

```bash
# Initialize context directory in your project
amem init

# Check context status
amem status

# Load full context (what AI sees)
amem load

# Get AI-ready context packet (optimized for LLMs)
amem agent

# Detect stale context
amem drift
```

## Command Reference

| Command | Description |
|---------|-------------|
| `amem init` | Create `.context/` directory with template files |
| `amem status` | Show context summary with token estimate |
| `amem load` | Output assembled context markdown |
| `amem agent [--budget N]` | Print AI-ready context packet (default 4000 tokens) |
| `amem add <type> <content>` | Add decision/task/learning/convention |
| `amem complete <query>` | Mark matching task as done |
| `amem drift [--json]` | Detect stale paths, broken refs |
| `amem sync [--auto]` | Reconcile context with codebase |
| `amem compact` | Archive completed tasks |
| `amem watch [--log FILE]` | Watch for context-update commands |
| `amem hook <tool>` | Generate AI tool integration config |

### Examples

```bash
# Add a new task
amem add task "Implement user authentication"

# Record a decision
amem add decision "Use PostgreSQL for primary database"

# Note a learning
amem add learning "Mock functions must be hoisted in Jest"

# Mark a task complete
amem complete "user auth"

# Get context with custom token budget
amem agent --budget 8000

# Check for stale references (JSON output for automation)
amem drift --json
```

## Context Files

```
.context/
├── CONSTITUTION.md     # Hard invariants — NEVER violate these
├── TASKS.md            # Current and planned work
├── DECISIONS.md        # Architectural decisions with rationale
├── LEARNINGS.md        # Lessons learned, gotchas, tips
├── CONVENTIONS.md      # Project patterns and standards
├── ARCHITECTURE.md     # System overview
├── DEPENDENCIES.md     # Key dependencies and why chosen
├── GLOSSARY.md         # Domain terms and abbreviations
├── DRIFT.md            # Staleness signals and update triggers
└── AGENT_PLAYBOOK.md   # How AI agents should use this system
```

## AI Tool Integration

Active Memory works with any AI tool that can read files. Generate tool-specific configs:

```bash
amem hook claude-code  # Claude Code CLI
amem hook cursor       # Cursor IDE
amem hook aider        # Aider
amem hook copilot      # GitHub Copilot
amem hook windsurf     # Windsurf IDE
```

### Claude Code

Add to your project's `CLAUDE.md`:

```markdown
## Active Memory Context

Before starting any task, load the project context:

1. Read .context/CONSTITUTION.md — These rules are INVIOLABLE
2. Read .context/TASKS.md — Current work items
3. Read .context/CONVENTIONS.md — Project patterns
4. Read .context/ARCHITECTURE.md — System overview
5. Read .context/DECISIONS.md — Why things are the way they are

When you make changes:
- Add decisions: <context-update type="decision">Your decision</context-update>
- Add tasks: <context-update type="task">New task</context-update>
- Add learnings: <context-update type="learning">What you learned</context-update>
- Complete tasks: <context-update type="complete">task description</context-update>

Run 'amem agent' for a quick context summary.
```

### Automated Context Updates

Use `amem watch` to automatically process context-update commands from AI output:

```bash
# Watch stdin (pipe AI output through this)
ai-tool | amem watch

# Watch a log file
amem watch --log /path/to/ai-output.log

# Dry run (preview without making changes)
amem watch --dry-run
```

## Design Philosophy

1. **File-based** — No database, no daemon. Just markdown and convention.
2. **Git-native** — Context versions with code, branches with code, merges with code.
3. **Human-readable** — Engineers can read, edit, and understand context directly.
4. **Token-efficient** — Markdown is cheaper than JSON/XML.
5. **Tool-agnostic** — Works with Claude Code, Cursor, Aider, Copilot, or raw CLI.

## Building with Ralph Wiggum

This project is designed to be built using the [Ralph Wiggum](https://ghuntley.com/ralph/) technique—an iterative AI development loop.

```bash
# Make the loop executable
chmod +x loop.sh

# Planning mode: Generate implementation plan
./loop.sh plan

# Building mode: Implement from plan
./loop.sh 20  # Max 20 iterations

# Unlimited building (Ctrl+C to stop)
./loop.sh
```

### Completion Promises

The loop automatically detects and stops on these signals:

| Promise | Meaning |
|---------|---------|
| `SYSTEM_CONVERGED` | All tasks complete — project is done |
| `BOOTSTRAP_COMPLETE` | Initial context created — ready to build |
| `PLANNING_CONVERGED` | Plan is complete — ready to build |
| `SYSTEM_BLOCKED` | All remaining tasks blocked — needs human input |

### Ralph Files

| File | Purpose |
|------|---------|
| `PROMPT_build.md` | Building mode instructions |
| `PROMPT_plan.md` | Planning mode instructions |
| `AGENTS.md` | Operational guide (how to build/run) |
| `IMPLEMENTATION_PLAN.md` | Generated task list |
| `specs/*.md` | Feature specifications |
| `loop.sh` | The Ralph loop script |

## Specifications

See `specs/` for detailed specifications:

- [Core Architecture](specs/core-architecture.md)
- [Context File Formats](specs/context-file-formats.md)
- [Context Loader](specs/context-loader.md)
- [Context Updater](specs/context-updater.md)
- [CLI](specs/cli.md)
- [AI Tool Integration](specs/ai-tool-integration.md)

## Contributing

1. Read the specs
2. Run `./loop.sh plan` to see current state
3. Run `./loop.sh` to let Ralph build
4. Review and commit changes

## License

MIT

---

*"Ralph is a Bash loop. The technique is deterministically bad in an undeterministic world."*
— Geoffrey Huntley
