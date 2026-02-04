---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Session Journal
icon: lucide/book-open
---

![ctx](images/ctx-banner.png)

!!! danger "Important Security Note"

    Session journals contain **sensitive data** such as
    file contents, commands, API keys, internal discussions, 
    error messages with stack traces, and more. 
    
    The `.context/journal-site/` directory **MUST** be `..gitignore`d.

    * **DO NOT** host your journal publicly.
    * **DO NOT** commit your journal files to version control.

## Browse Your Session History

`ctx`'s **Session Journal** turns your AI coding sessions into a **browsable**, 
**searchable**, and **editable** archive.

## Quick Start

After using `ctx` for a couple of sessions, you can generate a 
journal site with:

```bash
# Export all sessions to markdown
ctx recall export --all

# Generate and serve the journal site
ctx journal site --serve
```

Then open [http://localhost:8000](http://localhost:8000) to browse your sessions.

## What You Get

The Session Journal gives you:

* **Browsable history**: Navigate through all your AI sessions by date
* **Full conversations**: See every message, tool use, and result
* **Token usage**: Track how many tokens each session consumed
* **Search**: Find sessions by content, project, or date
* **Dark mode**: Easy on the eyes for late-night archaeology

Each session page includes the following sections:

| Section      | Content                                          |
|--------------|--------------------------------------------------|
| Metadata     | Date, time, duration, model, project, git branch |
| Summary      | Space for your notes (editable)                  |
| Tool Usage   | Which tools were used and how often              |
| Conversation | Full transcript with timestamps                  |

## The Workflow

### 1. Export Sessions

```bash
# Export all sessions from current project
ctx recall export --all

# Export sessions from all projects
ctx recall export --all --all-projects

# Export a specific session by ID
ctx recall export abc123

# Re-export (NOTE: overwrites existing files; back up first!)
ctx recall export --all --force
```

Exported sessions go to `.context/journal/` as editable Markdown files.

### 2. Generate the Site

```bash
# Generate site structure
ctx journal site

# Generate and build static HTML
ctx journal site --build

# Generate and serve locally
ctx journal site --serve

# Custom output directory
ctx journal site --output ~/my-journal
```

The site is generated in `.context/journal-site/` by default.

### 3. Browse and Search

Open [http://localhost:8000](http://localhost:8000) after running `--serve`.

- Use the sidebar to navigate by date
- Use search (`/` key) to find specific content
- Click any session to see the full conversation

## Editing Sessions

Exported sessions are plain Markdown in `.context/journal/`. You can:

- **Add summaries** - Fill in the `## Summary` section
- **Add notes** - Insert your own commentary anywhere
- **Highlight key moments** - Use Markdown formatting
- **Delete noise** - Remove irrelevant tool outputs

After editing, regenerate the site:

```bash
ctx journal site --serve
```

!!! warning "Edits are preserved unless you use `--force`"

    Running `ctx recall export --all` **skips existing files** by default,
    preserving your edits.

    However, `ctx recall export --all --force`
    **overwrites everything**, and your edits will be **lost**.

    Back up your `.context/journal/` directory before using `--force`.

## Large Sessions

Sessions with many messages (200+) are automatically split into multiple parts for better browser performance. Navigation links connect the parts:

```
session-abc123.md      (Part 1 of 3)
session-abc123-p2.md   (Part 2 of 3)
session-abc123-p3.md   (Part 3 of 3)
```

## Suggestion Sessions

Claude Code generates "suggestion" sessions for auto-complete prompts. These are separated in the index under a "Suggestions" section to keep your main session list focused.

## Tips

**Daily workflow:**
```bash
# At end of day, export and browse
ctx recall export --all && ctx journal site --serve
```

**After a productive session:**
```bash
# Export just that session and add notes
ctx recall export <session-id>
# Edit .context/journal/<session>.md
# Regenerate: ctx journal site
```

**Searching across all sessions:**
```bash
# Use grep on the journal directory
grep -r "authentication" .context/journal/
```

## Requirements

The journal site uses [zensical](https://pypi.org/project/zensical/) for static site generation:

```bash
pip install zensical
```

## See Also

- [ctx recall](../cli-reference.md#ctx-recall) - Session discovery and listing
- [ctx journal](../cli-reference.md#ctx-journal) - Site generation commands
- [Context Files](../context-files.md) - The `.context/` directory structure
