---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Integrating ctx
icon: lucide/package-plus
---

![ctx](images/ctx-banner.png)

## Adopting ctx in Existing Projects

!!! tip "Claude Code User?"
    You probably want the plugin instead of this page.
    Install ctx from the marketplace (`/plugin` → search "ctx" → Install) and
    you're done — hooks, skills, and updates are handled for you.
    See [Getting Started](getting-started.md) for the full walkthrough.

This guide covers adopting `ctx` in existing projects — regardless of
which tools your team uses.

## Quick Paths

| You have...                        | Command                                           | What happens                                         |
|------------------------------------|---------------------------------------------------|------------------------------------------------------|
| Nothing (*greenfield*)             | `ctx init`                                        | Creates `.context/`, `CLAUDE.md`, permissions        |
| Existing `CLAUDE.md`               | `ctx init --merge`                                | Backs up your file, inserts ctx block after the H1   |
| Existing `CLAUDE.md` + ctx markers | `ctx init --force`                                | Replaces the ctx block, leaves your content intact   |
| `.cursorrules` / `.aider.conf.yml` | `ctx init`                                        | ctx ignores those files — they coexist cleanly       |
| Team repo, first adopter           | `ctx init --merge && git add .context/ CLAUDE.md` | Initialize and commit for the team                   |

---

## Existing `CLAUDE.md`

This is the most common scenario. You have a `CLAUDE.md` with project-specific
instructions and don't want to lose them.

### What `ctx init` Does

When `ctx init` detects an existing `CLAUDE.md`, it checks for ctx markers
(`<!-- ctx:context -->` ... `<!-- ctx:end -->`):

| State                    | Default behavior         | With `--merge`            | With `--force`            |
|--------------------------|--------------------------|---------------------------|---------------------------|
| No `CLAUDE.md`           | Creates from template    | Creates from template     | Creates from template     |
| Exists, no ctx markers   | **Prompts** to merge     | Auto-merges (*no prompt*) | Auto-merges (*no prompt*) |
| Exists, has ctx markers  | Skips (*already set up*) | Skips                     | Replaces ctx block only   |

### The `--merge` Flag

`--merge` auto-merges without prompting. The merge process:

1. **Backs up** your existing `CLAUDE.md` to `CLAUDE.md.<timestamp>.bak`
2. **Finds the H1 heading** (e.g., `# My Project`) in your file
3. **Inserts** the ctx block immediately after it
4. **Preserves** everything else untouched

Your content before and after the ctx block remains exactly as it was.

### Before / After Example

**Before**: your existing `CLAUDE.md`:

```markdown
# My Project

## Build Commands

- `npm run build` — production build
- `npm test` — run tests

## Code Style

- Use TypeScript strict mode
- Prefer named exports
```

**After** `ctx init --merge`:

```markdown
# My Project

<!-- ctx:context -->
<!-- DO NOT REMOVE: This marker indicates ctx-managed content -->

## IMPORTANT: You Have Persistent Memory

This project uses Context (`ctx`) for context persistence across sessions.
...

<!-- ctx:end -->

## Build Commands

- `npm run build` — production build
- `npm test` — run tests

## Code Style

- Use TypeScript strict mode
- Prefer named exports
```

Your build commands and code style sections are untouched. The ctx block sits
between markers and can be updated independently.

### The `--force` Flag

If your `CLAUDE.md` already has ctx markers (from a previous `ctx init`), the
default behavior is to skip it. Use `--force` to replace the ctx block with the
latest template: This is useful after **upgrading** `ctx`:

```bash
ctx init --force
```

This only replaces content between `<!-- ctx:context -->` and `<!-- ctx:end -->`.
Your own content outside the markers is preserved. A timestamped backup is
created before any changes.

### Undoing a Merge

Every merge creates a backup:

```bash
$ ls CLAUDE.md*.bak
CLAUDE.md.1738000000.bak
```

To restore:

```bash
cp CLAUDE.md.1738000000.bak CLAUDE.md
```

Or if you are using `git`, simply:

```bash
git checkout CLAUDE.md
```

---

## Existing .cursorrules / Aider / Copilot

`ctx` doesn't touch tool-specific config files. It creates its own files
(`.context/`, `CLAUDE.md`) and coexists with whatever you already have.

### What ctx Creates vs. What It Leaves Alone

| ctx creates                              | ctx does NOT touch                    |
|------------------------------------------|---------------------------------------|
| `.context/` directory                    | `.cursorrules`                        |
| `CLAUDE.md` (or merges into)             | `.aider.conf.yml`                     |
| `.claude/settings.local.json` (seeded by `ctx init`; the plugin manages hooks and skills) | `.github/copilot-instructions.md` |
|                                          | `.windsurfrules`                      |
|                                          | Any other tool-specific config        |

Claude Code hooks and skills are provided by the **ctx plugin**,
installed from the Claude Code marketplace (`/plugin` → search "ctx" → Install).

### Running ctx Alongside Other Tools

The `.context/` directory is the source of truth. Tool-specific configs point
to it:

- **Cursor**: Reference `.context/` files in your system prompt
  (*see [Cursor setup](integrations.md#cursor-ide)*)
- **Aider**: Add `.context/` files to the `read:` list in `.aider.conf.yml`
  (*see [Aider setup](integrations.md#aider)*)
- **Copilot**: Keep `.context/` files open or reference them in comments
  (*see [Copilot setup](integrations.md#github-copilot)*)

You can generate a tool-specific configuration with:

```bash
ctx hook cursor    # Generate Cursor config snippet
ctx hook aider     # Generate .aider.conf.yml
ctx hook copilot   # Generate Copilot tips
ctx hook windsurf  # Generate Windsurf config
```

### Migrating Content Into .context/

If you have project knowledge scattered across `.cursorrules` or custom
prompt files, consider migrating it:

1. **Rules / invariants** → `.context/CONSTITUTION.md`
2. **Code patterns** → `.context/CONVENTIONS.md`
3. **Architecture notes** → `.context/ARCHITECTURE.md`
4. **Known issues / tips** → `.context/LEARNINGS.md`

You don't need to delete the originals — ctx and tool-specific files
can coexist. But centralizing in `.context/` means every tool gets the
same context.

---

## Team Adoption

### .context/ Is Designed to Be Committed

The `.context/` directory is meant to live in version control. It contains
project knowledge — not secrets or personal preferences.

```bash
# One person initializes
ctx init --merge

# Commit everything
git add .context/ CLAUDE.md
git commit -m "Add ctx context management"
git push
```

Teammates pull and immediately have context. No per-developer setup needed.

### What About .claude/?

The `.claude/` directory contains permissions that `ctx init` seeds.
Hooks and skills are provided by the ctx plugin (not per-project files).

| File                          | Commit? | Why                                  |
|-------------------------------|---------|--------------------------------------|
| `.claude/settings.local.json` | No      | Machine-specific, accumulates session permissions |
| `.claude/settings.golden.json`| Yes     | Curated permission snapshot (via `ctx permissions snapshot`) |

### Merge Conflicts in Context Files

Context files are plain Markdown. Resolve conflicts the same way you would
for any other documentation file:

```bash
# After a conflicting pull
git diff .context/TASKS.md    # See both sides
# Edit to keep both sets of tasks, then:
git add .context/TASKS.md
git commit
```

Common conflict scenarios:

- **TASKS.md**: Two people added tasks: Keep both.
- **DECISIONS.md**: Same decision recorded differently: Unify the entry
- **LEARNINGS.md**: Parallel discoveries: Keep both, remove duplicates

### Gradual Adoption

You don't need the whole team to switch at once:

1. One person runs `ctx init --merge` and commits.
2. `CLAUDE.md` instructions work immediately for Claude Code users.
3. Other tool users can adopt at their own pace using `ctx hook <tool>`.
4. Context files benefit everyone who reads them, even without tool integration.

---

## Verifying It Worked

### Check Status

```bash
ctx status
```

You should see your context files listed with token counts and no warnings.

### Test Memory

Start a new AI session and ask: **"Do you remember?"**

The AI should cite specific context:

* Current tasks from `.context/TASKS.md`;
* Recent decisions or learnings;
* Session history (*if you've had prior sessions*);

If it responds with generic "*I don't have memory*", check that `ctx` is in
your PATH (`which ctx`) and that hooks are configured
(see [Troubleshooting](integrations.md#troubleshooting)).

### Verify the Merge

If you used `--merge`, check that your original content is intact:

```bash
# Your original content should still be there
cat CLAUDE.md

# The ctx block should be between markers
grep -c "ctx:context" CLAUDE.md  # Should print 1
grep -c "ctx:end" CLAUDE.md      # Should print 1
```

---

## Further Reading

* [Getting Started](getting-started.md): Full setup walkthrough
* [Context Files](context-files.md): What each `.context/` file does
* [Integrations](integrations.md): Per-tool setup (*Claude Code, Cursor, Aider, Copilot*)
* [CLI Reference](cli-reference.md): All `ctx` commands and flags
