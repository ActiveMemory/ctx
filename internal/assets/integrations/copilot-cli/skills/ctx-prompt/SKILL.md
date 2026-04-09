---
name: ctx-prompt
description: "Apply, list, and manage saved prompt templates from .context/prompts/. Use when the user asks to apply, list, or create a reusable template like code-review or refactor."
tools: [bash, read, write]
---

Apply reusable prompt templates from `.context/prompts/`.

## When to Use

- User says "use the code-review prompt" or "apply the refactor template"
- User asks to list, create, or manage prompt templates
- User mentions "prompt template" or "reusable prompt"

## When NOT to Use

- For structured context entries (use `ctx add` instead)
- For full workflow automation (use a dedicated skill instead)
- For scratchpad notes (use `ctx pad` instead)

## Command Mapping

| User intent                      | Command                         |
|----------------------------------|---------------------------------|
| "list my prompts"                | `ctx prompt list`               |
| "show the code-review prompt"    | `ctx prompt show code-review`   |
| "create a new prompt"            | `ctx prompt add <name> --stdin` |
| "delete the debug prompt"        | `ctx prompt rm debug`           |

## Execution

**When no name is given:**
```bash
ctx prompt list
```

**When a name is given:**
```bash
ctx prompt show <name>
```

Read the prompt content, then follow the instructions in the
prompt applied to the user's current context.

## Quality Checklist

- [ ] Used correct subcommand for user intent
- [ ] Prompt content was applied, not just displayed
- [ ] If prompt not found, suggested `ctx prompt list`
