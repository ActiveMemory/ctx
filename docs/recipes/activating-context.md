---
title: "Activating a Context Directory"
icon: lucide/plug-zap
---

![ctx](../images/ctx-banner.png)

## The Problem

You ran a `ctx` command and got:

```
Error: no context directory specified for this project
```

This means ctx doesn't know which `.context/` directory to operate
on. It will not guess, and it will not walk up from your current
working directory looking for one; that behavior was removed
deliberately, because silent inference was the source of several
bugs (stray agent-created directories, cross-project bleed-through,
webhook-route misrouting, sub-agent fragmentation). Every `ctx`
command requires you to declare the target directory explicitly.

This page shows you the three ways to do that and when to use each.

## TL;DR

If the project has already been initialized and you just need to
bind it for your shell:

```bash
eval "$(ctx activate)"
```

That's 95% of the time. Add it to `.zshrc` / `.bashrc` per project
with direnv, or run it once per terminal.

## When You See the Error

The exact error message depends on how many `.context/` directories
are visible from the current directory:

### Zero Candidates

```
Error: no context directory specified for this project
```

Either you haven't initialized this project yet (run `ctx init`)
or you're in a directory that doesn't belong to a ctx-tracked
project. If you know the project lives elsewhere, use one of the
declaration methods below with its absolute path.

### One Candidate

```
Error: no context directory specified; a likely candidate is at
    /Users/you/repos/myproject/.context
```

ctx found a single `.context/` on the way up from here but won't
bind to it automatically. Run `eval "$(ctx activate)"` and ctx
will emit the `export` for the candidate. Or set `CTX_DIR` by hand.

### Multiple Candidates

```
Error: no context directory specified; multiple candidates visible:
  /Users/you/repos/myproject/.context
  /Users/you/repos/myproject/packages/web/.context
```

You're inside nested projects. Pick the one you mean:

```bash
ctx activate /Users/you/repos/myproject/.context
# …copy and paste the `export` line it prints, or wrap in eval:
eval "$(ctx activate /Users/you/repos/myproject/.context)"
```

## Three Ways to Declare

### 1. `ctx activate` (Recommended for Shells)

`ctx activate` emits a shell-native `export CTX_DIR=...` line to
stdout. Wrap it in `eval` and the binding takes effect for the
current shell:

```bash
# Walk up from current dir and bind the single visible candidate:
eval "$(ctx activate)"

# Bind a specific path explicitly:
eval "$(ctx activate /abs/path/to/.context)"

# Clear the binding:
eval "$(ctx deactivate)"
```

`ctx activate` validates paths strictly: the target must exist, be
a directory, and contain at least one canonical context file
(`CONSTITUTION.md` or `TASKS.md`). It refuses to emit for multiple
upward candidates; pick one explicitly in that case.

Under the hood, the emitted line is just:

```bash
export CTX_DIR='/abs/path/to/.context'
```

So you can copy it into your `.zshrc` / `.bashrc` if you want the
binding permanent for a given shell setup. Better: use
[direnv](https://direnv.net/) with a per-project `.envrc`.

### 2. `CTX_DIR` Env Var

If you already know the path, export it directly:

```bash
export CTX_DIR=/abs/path/to/.context
ctx status
```

`CTX_DIR` is the same variable `ctx activate` writes; `activate`
is just a convenience that figures out the path for you.

### 3. Inline One-Shot

For one-shot commands (CI jobs, scripts, debugging a specific
project without changing your shell state), prefix the binding
inline:

```bash
CTX_DIR=/abs/path/to/.context ctx status
```

This binds `CTX_DIR` for that invocation only.

`CTX_DIR` must be an absolute path with `.context` as its basename.
Relative paths and other names are rejected on first use; the
basename guard catches the common footgun
(`export CTX_DIR=$(pwd)`) before stray writes can leak to the
project root.

## For CI and Scripts

Do not rely on shell activation in automated flows. Set `CTX_DIR`
explicitly at the top of the script:

```bash
#!/usr/bin/env bash
set -euo pipefail

export CTX_DIR="$GITHUB_WORKSPACE/.context"
ctx status
ctx drift
```

## For Claude Code Users

The ctx plugin's hooks are generated with
`CTX_DIR="$CLAUDE_PROJECT_DIR/.context"` prefixed to each command,
so hook-driven ctx invocations resolve correctly without any
per-session setup. You only need to activate manually when running
`ctx` yourself in a terminal.

## One Project, One `.context/`

The context directory is not a free-floating bag of files. It is
pinned to a project by contract: **`filepath.Dir(ContextDir())` is
the project root.** That parent directory is what `ctx sync`,
`ctx drift`, and the memory-drift hook scan for code, secret files,
and `MEMORY.md` respectively.

The practical consequences:

- **Don't share one `.context/` across multiple projects.** It holds
  per-project journals, per-session state, and per-project secrets.
  Pointing two codebases at the same directory corrupts all three.
- **If you want to share knowledge** (CONSTITUTION, CONVENTIONS,
  ARCHITECTURE) across projects, use `ctx hub`. It cherry-picks
  entries at the right granularity and keeps the per-project bits
  where they belong.
- **The `CTX_DIR` you activate is implicitly a project-root
  declaration.** Setting `CTX_DIR=/weird/place/.context` means
  you're telling ctx the project root is `/weird/place/`. That's
  your call to make; ctx does not police it.

### Recommended Layout

```
~/WORKSPACE/my-to-do-list
  ├── .git
  ├── .context          ← owned by this project; do not share
  ├── ideas
  │   └── ...
  ├── Makefile
  ├── Makefile.ctx
  └── specs
      └── ...
```

`.context/` sits at the project root, next to `.git`. `ctx activate`
binds to it; every ctx subsystem reads the project from its parent.

## Why Not Walk Up Automatically?

Nested projects, submodules, rogue agent-created `.context/`
directories, and sub-agent sessions all produced silent misrouting
under the old walk-up model. See the
[explicit-context-dir spec](https://github.com/ActiveMemory/ctx/blob/main/specs/explicit-context-dir.md)
and [the analysis doc](https://github.com/ActiveMemory/ctx/blob/main/specs/context-resolution-analysis.md)
for the full reasoning.

The short version: ctx decided to stop guessing and require the
caller to declare. Every other decision flows from there.
