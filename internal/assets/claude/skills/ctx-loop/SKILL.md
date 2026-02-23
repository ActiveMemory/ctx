---
name: ctx-loop
description: "Generate autonomous loop script. Use when setting up unattended iteration for a project."
allowed-tools: Bash(ctx:*)
---

Generate a ready-to-use autonomous loop shell script.

## Before Generating

1. **Check for existing loop script**: look for `loop.sh` in the
   project root; confirm before overwriting
2. **Verify PROMPT.md exists**: the generated script defaults to
   reading `PROMPT.md`; if missing, ask the user what prompt file
   to use
3. **Verify the context directory exists**: the loop pattern depends
   on persistent context; run `ctx init` first if needed

## When to Use

- When setting up a project for autonomous iteration
- When the user wants to run unattended AI development
- When switching AI tools (e.g., Claude to Aider) and need a
  new loop script
- When customizing loop parameters (max iterations, completion
  signal, prompt file)

## When NOT to Use

- For interactive pair-programming sessions (just use the AI
  tool directly)
- When the user already has a working loop script and has not
  asked for changes
- When the project lacks a context directory and `PROMPT.md` (set
  those up first with `ctx init --ralph`)

## Usage Examples

```text
/ctx-loop
/ctx-loop --tool aider
/ctx-loop --prompt TASKS.md --max-iterations 10
/ctx-loop --completion SYSTEM_BLOCKED --output my-loop.sh
```

## Flags

| Flag               | Short | Default            | Purpose                       |
|--------------------|-------|--------------------|-------------------------------|
| `--prompt`         | `-p`  | `PROMPT.md`        | Prompt file the loop reads    |
| `--tool`           | `-t`  | `claude`           | AI tool: claude, aider, generic |
| `--max-iterations` | `-n`  | `0` (unlimited)    | Stop after N iterations       |
| `--completion`     | `-c`  | `SYSTEM_CONVERGED` | Signal that ends the loop     |
| `--output`         | `-o`  | `loop.sh`          | Output script filename        |

## Supported Tools

| Tool      | Command generated                    |
|-----------|--------------------------------------|
| `claude`  | `claude --print "$(cat <prompt>)"`   |
| `aider`   | `aider --message-file <prompt>`      |
| `generic` | Template stub for custom AI CLI      |

## Completion Signals

The loop watches AI output for these signals:

| Signal               | Meaning                              |
|----------------------|--------------------------------------|
| `SYSTEM_CONVERGED`   | All tasks complete; loop exits       |
| `SYSTEM_BLOCKED`     | Needs human input; loop exits        |
| `BOOTSTRAP_COMPLETE` | Initial scaffolding done; loop exits |

## Execution

```bash
ctx loop $ARGUMENTS
```

The command writes a shell script (default `loop.sh`) and makes
it executable. Report the generated path and how to run it:

```bash
chmod +x loop.sh   # already done by ctx loop
./loop.sh
```

## Safety Notes

- The generated script includes `set -e` and a 1-second sleep
  between iterations to prevent runaway loops
- `--max-iterations` is strongly recommended for first runs;
  suggest a reasonable default (e.g., 10) if the user omits it
- The script captures AI tool errors with `|| true` so one
  failed iteration does not kill the loop
- Autonomous agents benefit from explicit reasoning prompts in
  PROMPT.md — adding "think step-by-step before each change"
  to the iteration prompt significantly improves accuracy and
  reduces cascading mistakes in unattended runs

## Quality Checklist

Before reporting success, verify:
- [ ] Generated script exists at the output path
- [ ] Script is executable
- [ ] Prompt file referenced in the script actually exists
- [ ] If `--max-iterations 0`, user is aware it runs until
      a completion signal (warn them)
