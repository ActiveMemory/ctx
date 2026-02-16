---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Common Workflows
icon: lucide/repeat
---

## Track Context

```bash
# Add a task
ctx add task "Implement user authentication"

# Record a decision (full ADR fields required)
ctx add decision "Use PostgreSQL for primary database" \
  --context "Need a reliable database for production" \
  --rationale "PostgreSQL offers ACID compliance and JSON support" \
  --consequences "Team needs PostgreSQL training"

# Note a learning
ctx add learning "Mock functions must be hoisted in Jest" \
  --context "Tests failed with undefined mock errors" \
  --lesson "Jest hoists mock calls to top of file" \
  --application "Place jest.mock() before imports"

# Mark task complete
ctx complete "user auth"
```

## Check Context Health

```bash
# Detect stale paths, missing files, potential secrets
ctx drift

# See full context summary
ctx status
```

## Browse Session History

List and search past AI sessions from the terminal:

```bash
ctx recall list --limit 5
```

### Journal Site

Export session transcripts to a browsable static site with search,
navigation, and topic indices.

!!! info ""
    The `ctx journal` command requires
    [zensical](https://pypi.org/project/zensical/) (**Python >= 3.10**).

    `zensical` is a Python-based static site generator from the
    *Material* for *MkDocs* team.

    (*[why zensical?](blog/2026-02-15-why-zensical.md)*).

If you don't have it on your system,
install `zensical` once with [pipx](https://pipx.pypa.io/):

```bash
# One-time setup
pipx install zensical
```

!!! warning "Avoid `pip install zensical`"
    `pip install` often fails: For example, on macOS, system Python installs a
    non-functional stub (*`zensical` requires `Python >= 3.10`*), and
    Homebrew Python blocks system-wide installs (`PEP 668`).

    `pipx` creates an **isolated environment** with the
    **correct Python version** automatically.

Then, **export and serve**:

```bash
# Export all sessions to .context/journal/
ctx recall export --all

# Generate and serve the journal site
ctx journal site --serve
```

Open [http://localhost:8000](http://localhost:8000) to browse.

To update after new sessions, run the same two commands again;
`recall export` preserves existing YAML frontmatter and only
updates conversation content.

See [Session Journal](session-journal.md) for the full pipeline
including **normalization** and **enrichment**.

## Run an Autonomous Loop

Generate a script that iterates an AI agent until a completion
signal is detected:

```bash
ctx loop
chmod +x loop.sh
./loop.sh
```

See [Autonomous Loops](autonomous-loop.md) for configuration
and advanced usage.

----

**Next Up**: [Context Files →](context-files.md) — what each `.context/` file does and how to use it.

**See Also**:

* [Recipes](recipes/index.md) — targeted how-to guides for specific tasks
* [Knowledge Capture](recipes/knowledge-capture.md) — patterns for recording decisions, learnings, and conventions
* [Context Health](recipes/context-health.md) — keeping your `.context/` accurate and drift-free
* [Session Archaeology](recipes/session-archaeology.md) — digging into past sessions
* [Task Management](recipes/task-management.md) — tracking and completing work items
