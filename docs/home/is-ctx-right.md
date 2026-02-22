---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Is It Right for Me?
icon: lucide/compass
---

![ctx](images/ctx-banner.png)

## Good Fit

`ctx` shines when context matters more than code.

If any of these sound like your project, it's worth trying:

- **Multi-session AI work**: You use AI across many sessions on the same
  codebase, and re-explaining is slowing you down.
- **Architectural decisions that matter**: Your project has non-obvious
  choices (*database, auth strategy, API design*) that the AI keeps
  second-guessing.
- **"Why" matters as much as "what"**: you need the AI to understand
  *rationale*, not just current code
- **Team handoffs**: Multiple people (*or multiple AI tools*) work on the
  same project and need shared context.
- **AI-assisted development across tools**: Uou switch between Claude Code,
  Cursor, Copilot, or other tools and want context to follow the project,
  not the tool.
- **Long-lived projects**: Anything you'll work on for weeks or months,
  where accumulated knowledge has compounding value.

---

## May Not Be the Right Fit

`ctx` adds overhead that isn't worth it for every project. Be honest about
when to skip it:

- **One-off scripts**: If the project is a single file you'll finish today,
  there's nothing to remember.
- **RAG-only workflows**: If retrieval from an external knowledge base already
  gives the agent everything it needs for each session, adding `ctx` may be
  unnecessary. RAG retrieves information; `ctx` defines the project's
  working memory.
- **No AI involvement**: `ctx` is designed for human–AI workflows; without
  an AI consumer, the files are just documentation.
- **Enterprise-managed context platforms**: If your organization provides
  centralized context services, `ctx` may duplicate that layer.

For a deeper technical comparison with RAG, prompt management tools, and
agent frameworks, see [ctx and Similar Tools](comparison.md).

---

## Project Size Guide

### Solo Developer, Single Repo

This is `ctx`'s sweet spot. 

You get the most value here: one person, one project, decisions, and learnings 
accumulating over time. Setup takes 5 minutes and the `.context/` directory
directory stays small, and every session gets faster.

### Small Team, One or Two Repos

Works well. 

Context files commit to git, so the whole team shares the same
decisions and conventions.
Each person's AI starts with the team’s decisions already loaded.
Merge conflicts on `.context/` files are rare and
easy to resolve (*they are just Markdown*).

### Multiple Repos or Larger Teams

`ctx` operates per repository.

Each repo has its own `.context/` directory with its own decisions,
tasks, and learnings. This matches the way code, ownership, and history
already work in `git`.

There is no built-in cross-repo context layer.

For organizations that need centralized, organization-wide knowledge,
`ctx` complements a platform solution by providing durable,
project-local working memory for AI sessions.

---

## 5-Minute Trial

Zero commitment. Try it, and delete `.context/` if it's not for you.

!!! tip "Using Claude Code?"
    Install the `ctx` plugin from the Marketplace for Claude-native hooks, 
    skills, and automatic context loading:

    1. Type `/plugin` and press Enter
    2. Select **Marketplaces** → **Add Marketplace**
    3. Enter `ActiveMemory/ctx`
    4. Back in `/plugin`, select **Install** and choose `ctx`

    You'll still need the `ctx` binary for the CLI: See
    [Getting Started](getting-started.md#installation) for install options.

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

Want to remove it later? One command:

```bash
rm -rf .context/
```

No dependencies to uninstall. No configuration to revert. **Just files**.

---

**Ready to try it out?**

* [Join the Community→](community.md): Open Source is better together.
* [Getting Started →](getting-started.md): Full installation and setup.
* [`ctx` and Similar Tools →](comparison.md): Detailed comparison with other approaches.
