# Prompt Template System Removal

## Problem

`ctx` ships two overlapping systems for agent instructions: **skills**
(`.claude/skills/`, with frontmatter, trigger descriptions, and tool
permissions) and **prompt templates** (`.context/prompts/`, plain markdown,
no metadata). Three of the four bundled prompts (code-review, explain,
refactor) are functionally lightweight skills — they give the agent working
instructions to apply. But because they live in the prompt system, the agent
doesn't discover them at session start, the playbook has no trigger mappings,
and even ctx's creator didn't know they existed. Adding metadata to prompts
to fix discoverability would recreate the skill system. The fourth prompt
(loop.md) is not an agent instruction at all — it's a file consumed by
`ctx loop` shell scripts.

## Approach

Eliminate the prompt template system as an agent instruction mechanism.
Skills are the single concept for "instructions the agent follows."

1. **Promote** `code-review.md`, `explain.md`, `refactor.md` to proper
   skills in `.claude/skills/` with SKILL.md frontmatter.
2. **Retain** `loop.md` as a `ctx loop` configuration file at
   `.context/loop.md` — it's not an agent skill, it's input for
   shell-script-driven autonomous loops.
3. **Remove** the `ctx prompt` CLI (list/show/add/rm) and the
   `/ctx-prompt` skill entirely. Not deprecated — deleted.
4. **Remove** `.context/prompts/` creation from `ctx init`.
5. **Update** AGENT_PLAYBOOK.md with a single generic instruction:
   check available skills. No per-skill hardcoding needed since skills
   are self-describing via frontmatter.
6. **Update** all user-facing docs and recipes.

## Behavior

### Happy Path

1. User runs `ctx init` on a new project
2. `.context/prompts/` is no longer created; no prompt templates deployed
3. Three new skills are available: `/ctx-code-review`, `/ctx-explain`,
   `/ctx-refactor`
4. `.context/loop.md` is created for `ctx loop` consumption
5. User says "review this code" → agent recognizes the trigger, invokes
   `/ctx-code-review`
6. User wants a custom agent instruction → uses `/ctx-skill-creator`,
   not `ctx prompt add`
7. User runs `ctx loop` → reads `.context/loop.md` by default

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Existing project has `.context/prompts/` with custom prompts | `ctx init --force` does not delete user's custom prompts; migration note in GETTING_STARTED.md tells users to convert custom prompts to skills via `/ctx-skill-creator` |
| User runs `ctx prompt` after removal | Command not found — removed from binary, not deprecated |
| User references `/ctx-prompt` in their own docs/scripts | Skill no longer exists; clear error. Release notes document the removal |
| `ctx loop` without `.context/loop.md` | `ctx loop --prompt` flag still accepts any file path; error message if default path missing |
| User has overridden loop.md content in `.context/prompts/loop.md` | Migration: `ctx init` now writes to `.context/loop.md`; old path no longer checked |

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| `ctx loop` default prompt missing | `error: prompt file not found: .context/loop.md` | Run `ctx init` or pass `--prompt <file>` |

## Interface

### Skills (new)

| Skill | Trigger phrases |
|-------|----------------|
| `/ctx-code-review` | "review this code", "review this change", "code review" |
| `/ctx-explain` | "explain this code", "explain this", "what does this do" |
| `/ctx-refactor` | "refactor this", "clean this up" |

### CLI (removed)

- `ctx prompt` (all subcommands: list, show, add, rm) — deleted
- `/ctx-prompt` skill — deleted

### CLI (changed)

- `ctx loop --prompt` default: `.context/prompts/loop.md` → `.context/loop.md`

## Implementation

### Files to Create

| File | Change |
|------|--------|
| `internal/assets/claude/skills/ctx-code-review/SKILL.md` | New skill — promotes code-review.md |
| `internal/assets/claude/skills/ctx-explain/SKILL.md` | New skill — promotes explain.md |
| `internal/assets/claude/skills/ctx-refactor/SKILL.md` | New skill — promotes refactor.md |

### Files to Delete

| File | Reason |
|------|--------|
| `internal/assets/claude/skills/ctx-prompt/SKILL.md` | Skill removed |
| `internal/assets/prompt-templates/` | Entire directory — templates promoted to skills or moved |
| `internal/assets/read/prompt/` | Package — no consumers |
| `internal/cli/prompt/` | Entire command tree (list/show/add/rm) |
| `internal/cli/initialize/core/prompt/prompt_tpl.go` | Deploys prompt templates |
| `internal/cli/initialize/core/prompt/doc.go` | Package now empty |
| `internal/write/prompt/` | Output functions for prompt commands |
| `docs/recipes/prompt-templates.md` | Recipe for removed system |

### Files to Modify

| File | Change |
|------|--------|
| `internal/bootstrap/bootstrap.go` | Remove `prompt` command registration |
| `internal/assets/embed.go` | Remove `prompt-templates/*.md` glob |
| `internal/config/dir/dir.go` | Remove `Prompts` constant |
| `internal/config/asset/asset.go` | Remove `DirPromptTemplates` constant |
| `internal/config/loop/prompt.go` | Change default to `.context/loop.md` |
| `internal/config/embed/text/err_prompt.go` | Remove prompt-specific error keys |
| `internal/err/prompt/` | Remove prompt-specific error functions (keep if used elsewhere) |
| `internal/assets/commands/commands.yaml` | Remove `prompt` command entry |
| `internal/assets/commands/flags.yaml` | Remove prompt-related flags |
| `internal/assets/commands/text/write.yaml` | Remove prompt output strings |
| `internal/assets/commands/text/ui.yaml` | Remove prompt UI strings |
| `internal/assets/commands/text/err.yaml` | Remove prompt error strings |
| `internal/assets/embed_test.go` | Remove prompt template tests |
| `.context/AGENT_PLAYBOOK.md` | Add generic "check available skills" instruction |
| `docs/cli/tools.md` | Remove `ctx prompt` section, update `ctx loop` default |
| `docs/cli/index.md` | Remove `ctx prompt` entry |
| `docs/recipes/index.md` | Remove prompt-templates recipe reference |
| Other docs referencing prompt templates | Update per search |

## Testing

- **Unit**: Remove `TestPrompt*` tests in `embed_test.go`; remove
  `internal/cli/prompt/` tests
- **Integration**: Verify `ctx init` no longer creates `.context/prompts/`;
  verify `ctx loop` reads `.context/loop.md`
- **Edge case**: Verify build compiles clean with no dangling imports;
  verify `ctx loop` errors clearly when `.context/loop.md` is missing

## Non-Goals

- **Migration tooling**: no automated `ctx prompt migrate` command. Users
  with custom prompts convert manually via `/ctx-skill-creator`.
- **Prompt-as-concept removal from docs entirely**: the word "prompt" is
  fine; what's removed is the `.context/prompts/` template system as a
  parallel to skills.
- **Changing how `ctx loop` works**: it still reads a file and passes it
  to the AI tool. Only the default path changes.
