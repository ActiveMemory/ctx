---
title: "Keeping Context in a Separate Repo"
icon: lucide/folder-symlink
---

![ctx](../images/ctx-banner.png)

## The Problem

`ctx` files contain project-specific **decisions**, **learnings**,
**conventions**, and **tasks**. By default, they live in
`.context/` inside the project tree, and that works well when the context
can be public.

But sometimes you need the context *outside* the project:

* **Open-source projects with private context**: Your architectural notes,
  internal task lists, and scratchpad entries shouldn't ship with the public
  repo.
* **Compliance or IP concerns**: Context files reference sensitive design
  rationale that belongs in a separate access-controlled repository.
* **Personal preference**: You want to keep notes separate from code.

`ctx` supports this by letting you point `CTX_DIR` anywhere. This recipe
shows how to set that up and how to tell your AI assistant where to find the
context.

!!! warning "One `.context/` per project"
    The parent of the context directory is the project root by contract.
    `ctx sync`, `ctx drift`, and the memory-drift hook all read the
    codebase at `filepath.Dir(ContextDir())`. Pointing two projects at
    the same directory corrupts their journals, state, and secrets. To
    share knowledge (CONSTITUTION / CONVENTIONS / ARCHITECTURE) across
    projects, use [`ctx hub`](hub-overview.md), not a shared `.context/`.

## TL;DR

Create the external context directory, initialize it, and bind it:

```bash
mkdir -p ~/repos/myproject-context && cd ~/repos/myproject-context && git init
cd ~/repos/myproject

# Bind CTX_DIR to the external location, then init creates files there.
export CTX_DIR=~/repos/myproject-context/.context
ctx init
```

All `ctx` commands now use the external directory. If you share the
setup across shells, add the `export CTX_DIR=...` line to your
shell rc, or source a per-project `.envrc` with direnv.

## What Works, What Quietly Degrades

The single-source-anchor contract states that
`filepath.Dir(CTX_DIR)` is the project root. When the context
lives outside the project tree, ctx still resolves correctly for
every operation that reads or writes inside `.context/`. But any
operation that scans the **codebase** scans the wrong tree, and
does so silently:

| Operation                       | Behavior with external `.context/`                |
|---------------------------------|---------------------------------------------------|
| `ctx status`, `agent`, `add`    | ✅ Works. Operates on files inside `CTX_DIR`.     |
| Journal, scratchpad, hub        | ✅ Works. Same reason.                            |
| `ctx sync`                      | ⚠️ Scans the *context repo*, not the code repo.   |
| `ctx drift`                     | ⚠️ Same. Reports nothing useful.                  |
| Memory-drift hook (`MEMORY.md`) | ⚠️ Looks for `MEMORY.md` next to the external `.context/`, not the code. |

Nothing errors. The code-aware operations just find an empty or
unrelated tree where the project root should be.

### Workaround: symlink the `.context/` into the code tree

If you want both the privacy of an external git repo *and* working
`ctx sync` / `drift` / memory-drift, symlink the external
`.context/` into the code repo and point `CTX_DIR` at the symlink:

```bash
# External repo holds the real files
mkdir -p ~/repos/myproject-context && cd ~/repos/myproject-context && git init

# Symlink it into the code repo
ln -s ~/repos/myproject-context/.context ~/repos/myproject/.context

# Bind CTX_DIR to the symlink path; ctx init will follow it
export CTX_DIR=~/repos/myproject/.context
ctx init
```

Now `filepath.Dir(CTX_DIR)` is the **code repo**, so code-aware
operations scan the right tree. The actual files still live in
the external repo and commit there. Add `.context` to the code
repo's `.gitignore` (or `.git/info/exclude`) so the symlink itself
isn't tracked by the code repo.

The basename guard is permissive about symlinks: it checks the
declared name, not the resolved target, so a `.context` symlink
pointing anywhere is accepted as long as the declared basename is
`.context`.

## Commands and Skills Used

| Tool            | Type         | Purpose                                 |
|-----------------|--------------|-----------------------------------------|
| `ctx init`      | CLI command  | Initialize context directory            |
| `ctx activate`  | CLI command  | Emit `export CTX_DIR=...` for the shell |
| `CTX_DIR`       | Env variable | Declare context directory per-session   |
| `.ctxrc`        | Config file  | Per-project configuration               |
| `/ctx-status`   | Skill        | Verify context is loading correctly     |

## The Workflow

### Step 1: Create the Private Context Repo

Create a separate repository for your context files. This can live anywhere:
a private GitHub repo, a shared drive, a sibling directory:

```bash
# Create the context repo
mkdir -p ~/repos/myproject-context
cd ~/repos/myproject-context
git init
```

### Step 2: Initialize ctx Pointing at It

From your project root, declare `CTX_DIR` pointing to the external
location, then initialize:

```bash
cd ~/repos/myproject
CTX_DIR=~/repos/myproject-context/.context ctx init
```

This creates the canonical `.context/` file set inside
`~/repos/myproject-context/` instead of `~/repos/myproject/.context/`.

### Step 3: Make It Stick

Declaring `CTX_DIR` on every command is tedious. Pick one of these
methods to make the configuration permanent. The context directory
itself must be declared via `CTX_DIR`; `.ctxrc` does not carry the
path.

#### Option A: `CTX_DIR` Environment Variable (*Recommended*)

```bash
# Direct path. Works for ctx status / agent / add but degrades
# code-aware operations. See "What Works, What Quietly Degrades".
export CTX_DIR=~/repos/myproject-context/.context

# Or, with the symlink approach above, point at the symlink path
# inside the code repo so code-aware operations stay healthy.
export CTX_DIR=~/repos/myproject/.context
```

Put either form in your shell profile (`~/.bashrc`, `~/.zshrc`)
or a direnv `.envrc`.

For a single session, run `eval "$(ctx activate)"` from any
directory inside the project where exactly one `.context/`
candidate is visible (the symlink counts). `activate` does not
accept a path argument; bind a specific path by exporting
`CTX_DIR` directly instead.

#### Option B: `.ctxrc` for Other Settings

Put any settings (token budget, priority order, freshness files) in a
`.ctxrc` at the project root (`dirname(CTX_DIR)`), which here is the
parent of the external `.context/`:

```yaml
# ~/repos/myproject-context/.ctxrc
token_budget: 16000
```

`.ctxrc` is always read from the parent of `CTX_DIR`, so this file is
picked up whenever `CTX_DIR` points at
`~/repos/myproject-context/.context`.

#### Resolution

`ctx` reads the context directory from a single channel: the
`CTX_DIR` environment variable. When `CTX_DIR` is unset, `ctx`
errors with a "no context directory specified" hint pointing at
`ctx activate` and this recipe. When set, the value must be an
absolute path with `.context` as its basename; relative paths and
other names are rejected on first use.

See
[Activating a Context Directory](activating-context.md) for the full
recipe.

### Step 4: Agent Auto-Discovery via Bootstrap

When context lives outside the project tree, your AI assistant needs to know
where to find it. The `ctx system bootstrap` command resolves the configured
context directory and communicates it to the agent automatically:

```bash
$ ctx system bootstrap
ctx system bootstrap
====================

context_dir: /home/user/repos/myproject-context/.context

Files:
  CONSTITUTION.md, TASKS.md, DECISIONS.md, ...
```

The `CLAUDE.md` template generated by `ctx init` already instructs the agent to
run `ctx system bootstrap` at session start. Because `CTX_DIR` is inherited
by child processes, your agent picks up the external path automatically.

Here is the relevant section from `CLAUDE.md` for reference:

```markdown
<!-- CLAUDE.md -->
1. **Run `ctx system bootstrap`**: CRITICAL, not optional.
   This tells you where the context directory is. If it returns any
   error, relay the error output to the user verbatim, point them at
   https://ctx.ist/recipes/activating-context/ for setup, and STOP.
   Do not try to recover; the user decides.
```

Moreover, every nudge (*context checkpoint, persistence reminder, etc.*) also
includes a `Context: /home/user/repos/myproject-context/.context` footer, so
the agent remains anchored to the correct directory even in long sessions.

Export `CTX_DIR` in your shell profile so every hook process inherits it:

```bash
export CTX_DIR=~/repos/myproject-context/.context
```

### Step 5: Share with Teammates

Teammates clone both repos and export `CTX_DIR`:

```bash
# Clone the project
git clone git@github.com:org/myproject.git
cd myproject

# Clone the private context repo
git clone git@github.com:org/myproject-context.git ~/repos/myproject-context
export CTX_DIR=~/repos/myproject-context/.context
```

If teammates use different paths, each developer sets their own `CTX_DIR`.

For encryption key distribution across the team, see the
[Syncing Scratchpad Notes](scratchpad-sync.md) recipe.

### Step 6: Day-to-Day Sync

The external context repo has its own git history. Treat it like any other
repo: commit and push after sessions:

```bash
cd ~/repos/myproject-context

# After a session
git add -A
git commit -m "Session: refactored auth module, added rate-limit learning"
git push
```

Your AI assistant can do this too. When ending a session:

```text
You: "Save what we learned and push the context repo."

Agent: [runs ctx add learning, then commits and pushes the context repo]
```

You can also set up a post-session habit: project code gets committed to the
project repo, context gets committed to the context repo.

----

## Conversational Approach

You don't need to remember the flags; simply ask your assistant:

### Set Up Your System Using Natural Language

```text
You: "Set up ctx to use ~/repos/myproject-context as the context directory."

Agent: "I'll set CTX_DIR to that path, run ctx init to materialize
       it, and show you the export line to add to your shell
       profile. Want me to seed the core context files too?"
```

### Configure Separate Repo for `.context` Folder Using Natural Language

```text
You: "My context is in a separate repo. Can you load it?"

Agent: [reads CTX_DIR, loads context from the external dir]
       "Loaded. You have 3 pending tasks, last session was about the auth
       refactor."
```

----

## Tips

* **Start simple**. If you don't need external context yet, don't set it up.
  The default `.context/` in-tree is the easiest path. Move to an external
  repo when you have a concrete reason.
* **One context repo per project**. Sharing a single context directory across
  multiple projects corrupts journals, state, and secrets. Use `ctx hub` for
  cross-project knowledge sharing.
* **Export `CTX_DIR` in your shell profile** so hooks and tools inherit the
  path without per-command flags.
* **Commit both repos at session boundaries**. Context without code history
  (*or code without context history*) loses half the value.

----

## Next Up

**[The Complete Session →](session-lifecycle.md)**: Walk through a
full ctx session from start to finish.

## See Also

* [Setting Up ctx Across AI Tools](multi-tool-setup.md): initial setup recipe
* [Syncing Scratchpad Notes Across Machines](scratchpad-sync.md): distribute
  encryption keys when context is shared
* [CLI Reference](../cli/index.md): full command list and global options
