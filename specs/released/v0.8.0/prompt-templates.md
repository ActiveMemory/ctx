# Spec: Prompt Templates (`ctx prompt`)

## Problem

ctx skills (SKILL.md with frontmatter, trigger rules, allowed-tools) are the
right weight for complex workflows, but overkill for reusable prompt patterns.
A user who wants a "code review checklist" or "refactoring guard rails" prompt
has two bad options:

- Author a full skill (high friction, needs frontmatter, rebuild to embed)
- Keep prompts in their head or a scratch file (no discoverability, no sharing)

The gap between "raw instruction typed into chat" and "full SKILL.md skill"
has no lightweight option.

## Decision

Add **prompt templates** — plain markdown files in `.context/prompts/` that
users can invoke via `/ctx-prompt <name>` or manage via `ctx prompt` CLI.
No frontmatter, no build step, no trigger rules. Just named markdown.

## Design

### Storage

```
.context/prompts/
  code-review.md      # plain markdown, no frontmatter
  refactor.md
  explain.md
```

**Committed to git** by default — team-shared prompts are more valuable than
personal ones. Users who want private prompts add `.context/prompts/` to
`.gitignore` (one line).

### Starter Templates

`ctx init` stamps a few genuinely useful starters from embedded templates,
same pattern as TASKS.md and CONVENTIONS.md:

| Template | Purpose |
|----------|---------|
| `code-review.md` | Review checklist anchored to project conventions |
| `refactor.md` | Refactoring with guard rails (tests first, preserve behavior) |
| `explain.md` | Explain code for onboarding / knowledge transfer |

Kept minimal. Users add their own.

### CLI

```
ctx prompt list             # list available prompts (names from .context/prompts/)
ctx prompt show <name>      # print prompt content to stdout
ctx prompt add <name>       # create from starter template or stdin
ctx prompt rm <name>        # delete a prompt file
```

All commands read from `.context/prompts/` on disk (not embedded assets).
The embedded templates are only used by `ctx init` and `ctx prompt add` as
scaffolding.

### Skill: `/ctx-prompt`

A single bundled skill that handles all interaction:

- **`/ctx-prompt code-review`** — runs `ctx prompt show code-review`, injects
  the content as working instructions for the current task.
- **`/ctx-prompt` (no argument)** — runs `ctx prompt list`, presents available
  prompts so the user can pick one.

The skill SKILL.md instructs the agent:
1. If no name given, run `ctx prompt list` and present the results.
2. If a name is given, run `ctx prompt show <name>` and follow the prompt
   instructions in the current context.

### Autocomplete Trade-off

The skill name `/ctx-prompt` autocompletes via Claude Code's tab completion.
The slug argument (e.g., `code-review`) does not autocomplete. This is
acceptable because:

- Users create and name their own prompts — muscle memory develops naturally
- For forgotten names, `/ctx-prompt` with no arg lists everything
- The slug namespace is small (typically 5-15 prompts per project)

## Implementation

### Embedded Templates

```
internal/assets/prompt-templates/
  code-review.md
  refactor.md
  explain.md
```

Added to the `//go:embed` directive in `embed.go`. Accessor functions:

```go
func PromptTemplate(name string) ([]byte, error)
func ListPromptTemplates() ([]string, error)
```

### CLI Commands

New package: `internal/cli/prompt/`

| File | Command | Notes |
|------|---------|-------|
| `prompt.go` | `ctx prompt` | Parent command |
| `list.go` | `ctx prompt list` | Reads `.context/prompts/` directory |
| `show.go` | `ctx prompt show <name>` | Prints file content to stdout |
| `add.go` | `ctx prompt add <name>` | Creates from embedded template or stdin (`--stdin`) |
| `rm.go` | `ctx prompt rm <name>` | Deletes with confirmation |

### Bundled Skill

```
internal/assets/claude/skills/ctx-prompt/SKILL.md
```

Frontmatter: `name: ctx-prompt`, `allowed-tools: Bash(ctx:*)`.

### `ctx init` Integration

`ctx init` creates `.context/prompts/` and stamps starter templates if the
directory doesn't exist. Same idempotency pattern as other context files —
skip if already present.

### Config Constants

In `internal/config/`:

```go
DirPrompts   = "prompts"       // subdirectory under context_dir
```

### Files Changed

| File | Change |
|------|--------|
| `internal/assets/embed.go` | Add `prompt-templates/*.md` to embed, add accessors |
| `internal/assets/embed_test.go` | Test `PromptTemplate()` and `ListPromptTemplates()` |
| `internal/cli/prompt/prompt.go` | New — parent command |
| `internal/cli/prompt/list.go` | New — list prompts |
| `internal/cli/prompt/show.go` | New — show prompt |
| `internal/cli/prompt/add.go` | New — add prompt |
| `internal/cli/prompt/rm.go` | New — remove prompt |
| `internal/cli/prompt/*_test.go` | New — tests for each command |
| `internal/cli/root.go` | Register `prompt.Cmd()` |
| `internal/config/dir.go` | Add `DirPrompts` constant |
| `internal/cli/initialize/init.go` | Create prompts dir and stamp templates |
| `internal/assets/claude/skills/ctx-prompt/SKILL.md` | New — bundled skill |
| `internal/assets/prompt-templates/*.md` | New — starter templates |
| `permissions/allow.txt` | Add `Skill(ctx-prompt)` |

## Non-Goals

- Template variable substitution (`{{language}}`, `{{project}}`) — the AI
  adapts generic prompts to context naturally; a parser adds complexity
  without proportional value
- User-level prompts at `~/.ctx/prompts/` — can be added later if demand
  exists; project-level is the right default
- Prompt chaining or composition — each prompt is standalone; composition
  is the agent's job
- Trigger rules or auto-invocation — prompts are explicit, user-initiated

## Dependencies

- None. This is additive — no existing code changes behavior.

## Verification

```bash
make build && make test && make lint

# Init stamps templates
rm -rf /tmp/test-prompt && mkdir /tmp/test-prompt && cd /tmp/test-prompt
git init && ctx init
ls .context/prompts/

# CLI works
ctx prompt list
ctx prompt show code-review
echo "# My Custom Prompt" | ctx prompt add my-prompt --stdin
ctx prompt list
ctx prompt rm my-prompt

# Skill works (in Claude Code session)
# /ctx-prompt              → lists prompts
# /ctx-prompt code-review  → injects code review prompt
```
