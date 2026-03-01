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
* **Personal preference**: You want a single context repo that covers
  multiple projects, or you just prefer keeping notes separate from code.

`ctx` supports this through three configuration methods. This recipe shows how
to set them up and how to tell your AI assistant where to find the context.

## TL;DR

First `--allow-outside-cwd` in your project:

```bash
mkdir ~/repos/myproject-context && cd ~/repos/myproject-context && git init
cd ~/repos/myproject
ctx --context-dir ~/repos/myproject-context --allow-outside-cwd init
```

Then, [create a `.ctxrc`](../home/configuration.md) in your **project root**
to specify the new `.context` folder location:

```yaml
context_dir: ~/repos/myproject-context
allow_outside_cwd: true
```

All `ctx` commands now use the external directory automatically.

## Commands and Skills Used

| Tool                  | Type         | Purpose                                 |
|-----------------------|--------------|-----------------------------------------|
| `ctx init`            | CLI command  | Initialize context directory            |
| `--context-dir`       | Global flag  | Point ctx at a non-default directory    |
| `--allow-outside-cwd` | Global flag  | Permit context outside the project root |
| `.ctxrc`              | Config file  | Persist the context directory setting   |
| `CTX_DIR`             | Env variable | Override context directory per-session  |
| `/ctx-status`         | Skill        | Verify context is loading correctly     |

## The Workflow

### Step 1: Create the Private Context Repo

Create a separate repository for your context files. This can live anywhere:
a private GitHub repo, a shared drive, a sibling directory:

```bash
# Create the context repo
mkdir ~/repos/myproject-context
cd ~/repos/myproject-context
git init
```

### Step 2: Initialize ctx Pointing at It

From your project root, initialize ctx with `--context-dir` pointing to the
external location. Because the directory is outside your project tree, you also
need `--allow-outside-cwd`:

```bash
cd ~/repos/myproject
ctx --context-dir ~/repos/myproject-context \
    --allow-outside-cwd \
    init
```

This creates the full `.context/`-style file set inside
`~/repos/myproject-context/` instead of `~/repos/myproject/.context/`.

!!! warning "Boundary Validation"
    `ctx` validates that the `.context` directory is within the current working
    directory. 

    If your external directory is truly outside the project root:

    * Either every `ctx` command needs `--allow-outside-cwd`, 
    * or you can persist the setting in `.ctxrc` (*next step*).

### Step 3: Make It Stick

Typing `--context-dir` and `--allow-outside-cwd` on every command is tedious.
Pick one of these methods to make the configuration permanent.

#### Option A: `.ctxrc` (*Recommended*)

Create a `.ctxrc` file in your project root:

```yaml
# .ctxrc: committed to the project repo
context_dir: ~/repos/myproject-context
allow_outside_cwd: true
```

ctx reads `.ctxrc` automatically. Every command now uses the external
directory without extra flags:

```bash
ctx status          # reads from ~/repos/myproject-context
ctx add learning "Redis MULTI doesn't roll back on error"
```

!!! tip "Commit `.ctxrc`"
    `.ctxrc` belongs in the project repo. It contains no secrets: It's just a
    path and a boundary override. 

    `.ctxrc` lets teammates share the same configuration.

#### Option B: `CTX_DIR` Environment Variable

Good for CI pipelines, temporary overrides, or when you don't want to commit
a `.ctxrc`:

```bash
# In your shell profile (~/.bashrc, ~/.zshrc)
export CTX_DIR=~/repos/myproject-context
```

Or for a single session:

```bash
CTX_DIR=~/repos/myproject-context ctx status
```

#### Option C: Shell Alias

If you prefer a shell alias over `.ctxrc`:

```bash
# ~/.bashrc or ~/.zshrc
alias ctx='ctx --context-dir ~/repos/myproject-context --allow-outside-cwd'
```

#### Priority Order

When multiple methods are set, `ctx` resolves the context directory in this
order (*highest priority first*):

1. `--context-dir` flag
2. `CTX_DIR` environment variable
3. `context_dir` in `.ctxrc`
4. Default: `.context/`

### Step 4: Agent Auto-Discovery via Bootstrap

When context lives outside the project tree, your AI assistant needs to know
where to find it. The `ctx system bootstrap` command resolves the configured
context directory and communicates it to the agent automatically:

```bash
$ ctx system bootstrap
ctx bootstrap
=============

context_dir: /home/user/repos/myproject-context

Files:
  CONSTITUTION.md, TASKS.md, DECISIONS.md, ...
```

The `CLAUDE.md` template generated by `ctx init` already instructs the agent to
run `ctx system bootstrap` at session start. Because `.ctxrc` is in the
project root, your agent inherits the external path automatically via
the `ctx system boostrap` call instruction.

Here  is the relevant section from `CLAUDE.md` for reference:

```markdown
<!-- CLAUDE.md -->
1. **Run `ctx system bootstrap`**: CRITICAL, not optional.
   This tells you where the context directory is. If it fails or returns
   no context_dir, STOP and warn the user.
```

Moreover, every nudge (*context checkpoint, persistence reminder, etc.*) also 
includes a `Context: /home/user/repos/myproject-context` footer, so the agent
remains anchored to the correct directory even in long sessions.

If you use `CTX_DIR` instead of `.ctxrc`, export it in your shell
profile so the hook process inherits it:

```bash
export CTX_DIR=~/repos/myproject-context
```

### Step 5: Share with Teammates

Teammates clone both repos and set up `.ctxrc`:

```bash
# Clone the project
git clone git@github.com:org/myproject.git
cd myproject

# Clone the private context repo
git clone git@github.com:org/myproject-context.git ~/repos/myproject-context
```

If `.ctxrc` is already committed to the project, they're done: `ctx`
commands will find the external context automatically.

If teammates use different paths, each developer sets their own `CTX_DIR`:

```bash
export CTX_DIR=~/my-own-path/myproject-context
```

For encryption key distribution across the team, see the
[Syncing Scratchpad Notes](scratchpad-sync.md) recipe.

### Step 6: Day-to-Day Sync

The external context repo has its own git history. Treat it like any other
repo: Commit and push after sessions:

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

Agent: "I'll create a .ctxrc in the project root pointing to that path.
       I'll also update CLAUDE.md so future sessions know where to find
       context. Want me to initialize the context files there too?"
```

### Configure Separate Repo for `.context` Folder Using Natural Language

```text
You: "My context is in a separate repo. Can you load it?"

Agent: [reads .ctxrc, finds the path, loads context from the external dir]
       "Loaded. You have 3 pending tasks, last session was about the auth
       refactor."
```

----

## Tips

* **Start simple**. If you don't need external context yet, don't set it up.
  The default `.context/` in-tree is the easiest path. Move to an external
  repo when you have a concrete reason.
* **One context repo per project**. Sharing a single context directory across
  multiple projects creates confusion. Keep the mapping 1:1.
* **Use `.ctxrc` over env vars** when the path is stable. It's committed,
  documented, and works for the whole team without per-developer shell setup.
* **Don't forget the boundary flag**. The most common error is
  `Error: context directory is outside the project root`. Set
  `allow_outside_cwd: true` in `.ctxrc` or pass `--allow-outside-cwd`.
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
* [CLI Reference](../cli/index.md): all global flags including
  `--context-dir` and `--allow-outside-cwd`
