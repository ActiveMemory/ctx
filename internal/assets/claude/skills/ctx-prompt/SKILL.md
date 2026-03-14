---
name: ctx-prompt
description: "Apply, list, and manage saved prompt templates from .context/prompts/. Use when the user asks to apply, list, or create a reusable template like code-review or refactor."
allowed-tools: Bash(ctx:*)
---

Apply reusable prompt templates from `.context/prompts/` to the
current task. Prompt templates are plain Markdown files: no
frontmatter, no trigger rules: for common patterns like code
review, refactoring, or explaining code.

## When to Use

- User says "use the code-review prompt" or "apply the refactor template"
- User invokes `/ctx-prompt` with or without a name argument
- User asks to list, create, or manage prompt templates
- User mentions "prompt template" or "reusable prompt"

## When NOT to Use

- For structured context entries (use `ctx add` instead)
- For full workflow automation (use a dedicated skill instead)
- For scratchpad notes (use `ctx pad` instead)

## Command Mapping

| User intent                                  | Command                                               |
|----------------------------------------------|-------------------------------------------------------|
| "list my prompts" / "what prompts do I have" | `ctx prompt list`                                     |
| "show the code-review prompt"                | `ctx prompt show code-review`                         |
| "create a new prompt called debug"           | `ctx prompt add debug --stdin` (then ask for content) |
| "add the refactor template"                  | `ctx prompt add refactor`                             |
| "delete the debug prompt"                    | `ctx prompt rm debug`                                 |

## Execution

**When no name is given** (or user asks to list):

```bash
ctx prompt list
```

Present the results and ask which prompt to use.

**When a name is given:**

```bash
ctx prompt show <name>
```

Read the prompt content, then **follow the instructions in the prompt**
applied to the user's current context. The prompt template tells you
what to do: treat it as your working instructions.

## Interpreting User Intent

When the user provides a prompt name:

1. Run `ctx prompt show <name>` to retrieve the template
2. Apply the template's instructions to whatever the user is working on
3. If the user also provides a code snippet or file reference, apply
   the prompt to that specific target

When the user wants to create a prompt:

1. Ask what the prompt should contain if not provided
2. Use `echo "content" | ctx prompt add <name> --stdin` or guide them

## Important Notes

- Prompt templates are plain markdown: no frontmatter parsing needed
- Templates live in `.context/prompts/` and are committed to git by default
- `ctx init` stamps starter templates (code-review, refactor, explain)
- If a prompt is not found, suggest running `ctx prompt list` to see
  available prompts
