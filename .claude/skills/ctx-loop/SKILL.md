---
name: ctx-loop
description: "Generate autonomous loop script. Use when setting up unattended iteration for a project."
allowed-tools: Bash(ctx:*)
---

Generate a ready-to-use autonomous loop shell script.

## When to Use

- When setting up a project for autonomous iteration
- When the user wants to run unattended AI development

## Usage

```
/ctx-loop
/ctx-loop --tool aider
/ctx-loop --prompt PROMPT.md --max-iterations 10
```

Generates a shell script for iterative AI development. Defaults to Claude Code.

## Execution

```bash
ctx loop $ARGUMENTS
```

Report the generated script path and how to run it.
