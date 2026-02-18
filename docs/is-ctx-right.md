---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Is ctx Right for Me?
icon: lucide/compass
---

## Good Fit

`ctx` shines when context matters more than code. If any of these sound like
your project, it's worth trying:

- **Multi-session AI work** — you use AI across many sessions on the same
  codebase, and re-explaining is slowing you down
- **Architectural decisions that matter** — your project has non-obvious
  choices (database, auth strategy, API design) that the AI keeps
  second-guessing
- **"Why" matters as much as "what"** — you need the AI to understand
  *rationale*, not just current code
- **Team handoffs** — multiple people (or multiple AI tools) work on the
  same project and need shared context
- **AI-assisted development across tools** — you switch between Claude Code,
  Cursor, Copilot, or other tools and want context to follow the project,
  not the tool
- **Long-lived projects** — anything you'll work on for weeks or months,
  where accumulated knowledge has compounding value

---

## Not the Right Fit

`ctx` adds overhead that isn't worth it for every project. Be honest about
when to skip it:

- **One-off scripts** — if the project is a single file you'll finish today,
  there's nothing to remember
- **Pure RAG workflows** — if your AI setup already retrieves context from a
  knowledge base and that's sufficient, `ctx` solves a different problem
- **No AI involvement** — `ctx` is designed for human–AI workflows; without
  an AI consumer, the files are just documentation
- **Enterprise-managed context platforms** — if your organization provides
  centralized context services, `ctx` may duplicate that layer

For a deeper technical comparison with RAG, prompt management tools, and
agent frameworks, see [ctx and Similar Tools](comparison.md).

---

## Project Size Guide

### Solo developer, single repo

This is `ctx`'s sweet spot. You get the most value here: one person, one
project, decisions and learnings accumulating over time. Setup takes 5
minutes and the `.context/` directory is small enough to never think about.

### Small team, one or two repos

Works well. Context files commit to git, so the whole team shares the same
decisions and conventions. Each person's AI sessions benefit from what
others have recorded. Merge conflicts on `.context/` files are rare and
easy to resolve (they're just markdown).

### Multiple repos or larger teams

`ctx` works per-repository. Each repo gets its own `.context/` directory
with its own decisions, tasks, and learnings. There's no cross-repo context
sharing built in — each project is self-contained. For organizations that
need centralized knowledge, `ctx` can complement (but not replace) a
platform-level solution.

---

## 5-Minute Trial

Zero commitment. Try it, and delete `.context/` if it's not for you.

```bash
# 1. Initialize
cd your-project
ctx init

# 2. Add one real decision from your project
ctx add decision "Your actual architectural choice" \
  --context "What prompted this decision" \
  --rationale "Why you chose this approach" \
  --consequences "What changes as a result"

# 3. Check what the AI will see
ctx status

# 4. Start an AI session and ask: "Do you remember?"
```

If the AI cites your decision back to you, it's working.
If it doesn't add value for your workflow, clean up is one command:

```bash
rm -rf .context/
```

No dependencies to uninstall. No configuration to revert. Just files.

---

**Ready to try it?**

- [Getting Started →](getting-started.md) — full installation and setup
- [ctx and Similar Tools →](comparison.md) — detailed comparison with other approaches
