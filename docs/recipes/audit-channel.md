---
title: Out-of-Band Audit Channel
icon: lucide/scan-eye
---

![ctx](../images/ctx-banner.png)

## The Problem

The agent that just shipped a feature is the worst possible
reviewer of its own discipline. It will mark its own work
complete, label deferred docs as "Phase 2," and skip past
its own CONVENTIONS.md rule with conviction. Mid-task tunnel
vision suppresses the rules it read at session start.

You cannot fix that with more advisory prose: the same
convention that didn't stop the agent the first time won't
stop the next agent either. What works is **mechanical
verbatim relay** — the same channel ctx already uses for
`ctx remind`, journal-import nudges, and knowledge-growth
warnings. Agents echo those without filtering, every turn,
because the relay bypasses judgment.

This recipe shows how to run discipline audits **out of
band** (from a separate Claude Code session, on your
plan-billed subscription, not the working session's API) and
drop their findings onto the verbatim-relay channel so the
next interactive session sees them at the top of its next
turn.

## TL;DR

```bash
# 1. From a separate Claude Code session in the same project:
/ctx-surface-audit                # default: main..HEAD

# 2. The skill writes:
.context/audit/surface.md         # structured report

# 3. Back in the working session, the next prompt fires the
#    UserPromptSubmit hook:
ctx system check-audit            # invoked by the hook config

# 4. The agent / human sees a verbatim-relay box on the next
#    response, listing the specific files to fix.

# 5. After addressing the findings:
ctx audit dismiss surface         # stops the relay
```

## Commands and Skills Used

| Tool                        | Type        | Purpose                                                |
|-----------------------------|-------------|--------------------------------------------------------|
| `/ctx-surface-audit`        | Skill       | Out-of-band auditor; scans diff for surface drift      |
| `ctx audit list`            | CLI command | Show all reports with status and age                   |
| `ctx audit show ID`         | CLI command | Print one report's body (pipe-friendly)                |
| `ctx audit dismiss ID`      | CLI command | Mark a report dismissed against its current digest     |
| `ctx audit dismiss --all`   | CLI command | Bulk dismissal                                         |
| `ctx system check-audit`    | CLI command | UserPromptSubmit hook; verbatim-relays unread reports  |

## Why a Separate Session

Two reasons, both load-bearing:

1. **Fresh-context judgment.** The auditor must not inherit
   the implementer's working memory of "what we tried, what
   we decided to defer, why this is fine." The audit only
   works if the reviewer reads the diff cold.
2. **Cost shape.** A per-commit AI gate burns API tokens on
   every commit, regardless of branch maturity. Running the
   auditor manually from a separate Claude Code session
   bills against your interactive plan, not the API, and
   lets you decide *when* to spend the cycles (typically
   right before a PR, not on every micro-commit).

The `/ctx-surface-audit` skill enforces this with a hard
dirty-tree refusal: invoking it in a working session with
uncommitted changes returns

> Run this audit from a separate Claude Code session.

There is no override flag, by design.

## The Workflow

### Step 1: Land Your Work, Then Open a Second Session

Finish the feature on your working branch (commit, lint,
test). Open a second Claude Code window in the same project
worktree. The audit runs against `main..HEAD` by default.

### Step 2: Invoke the Auditor

```text
You (in session 2): "/ctx-surface-audit"

Skill: "Scanned 4 commits, 3 surfaces detected.
        Wrote .context/audit/surface.md (status: findings).
        Open a working session — the check-audit hook will
        relay the findings on the next prompt."
```

The auditor compares the branch against `main`, finds new
subcommands / flags / behavior changes, checks each one
against SKILL.md / recipe / `docs/cli` coverage, and writes
a structured report.

### Step 3: Return to the Working Session

The next time you submit a prompt in your working session,
the `ctx system check-audit` hook (a UserPromptSubmit hook
in `.claude/settings.local.json`) reads `.context/audit/`
and emits a verbatim-relay box at the top of the agent's
response:

```text
┌─ Audit Reports ──────────────────────────────────────
│ [surface] main..HEAD
│ Commit 6bcaf889 added user-facing surface without docs:
│
│   • New subcommand `ctx pad undo`
│     - SKILL.md: internal/assets/claude/skills/ctx-pad/SKILL.md
│       command-mapping table is missing the row
│     - Recipe: docs/recipes/scratchpad-with-claude.md unchanged
│
│ Fix:
│   - edit internal/assets/claude/skills/ctx-pad/SKILL.md
│   - edit docs/recipes/scratchpad-with-claude.md
│
│ Dismiss: ctx audit dismiss <id>
│ Dismiss all: ctx audit dismiss --all
└──────────────────────────────────────────────────
```

The agent echoes this verbatim — that is the discipline
mechanism. You (or the agent) then address each cited file.

### Step 4: Dismiss

Once you've addressed the findings (or accepted them as
out-of-scope), dismiss the report:

```bash
ctx audit dismiss surface
```

Dismissal is bound to the **report digest** at dismiss
time. A subsequent audit that produces the same findings
stays dismissed. A subsequent audit that finds *new* surface
drift produces a fresh digest and re-surfaces the report at
the next prompt.

## Retention

The audit channel keeps **one report per kind**. Re-running
`/ctx-surface-audit` overwrites the prior `surface.md`.
Reports older than 30 days are still relayed but prefixed
with a `STALE — main..HEAD (audited 32d ago)` marker so the
recipient knows the assessment may not match current code.

History (which audits ran when) is preserved by the
dismissal ledger at `.context/audit/.dismissed.json`. The
ledger lives next to the reports — not under
`.context/state/` — so nuking session state never silently
re-surfaces a dismissed audit.

## When to Run the Auditor

- **Before opening a PR.** The natural cadence. The audit
  exists to catch the gaps you can't see in your own
  branch.
- **After landing a multi-commit feature.** Especially
  when the feature added new subcommands or flags.
- **Periodically on `main`**, with a longer range like
  `HEAD~50..HEAD`, to catch surface drift that crept in
  before this channel existed.

There is **no automated trigger** in Phase 1. The cost shape
is intentional: cron and post-commit-hook drivers stay on
the deferred list until the user-driven workflow proves out.

## Other Audit Skills

`/ctx-surface-audit` is the first member of a family. The
scaffolding (channel, ledger, hook, CLI) is generic. Future
siblings under the same shape:

- `/ctx-spec-trailer-audit` — does each commit's `Spec:`
  trailer point at a spec that genuinely covers that
  commit's scope?
- `/ctx-capture-audit` — was a Decision or Learning
  persisted for non-trivial work that ended without one?

Each lives in its own SKILL.md and writes its own report
file (e.g. `.context/audit/spec-trailer.md`). The hook
relays whatever it finds, with no per-kind plumbing.

## See Also

- [Spec: out-of-band audit channel](https://github.com/ActiveMemory/ctx/blob/main/specs/audit-channel.md):
  full design rationale + Open Questions
- [CONVENTIONS → User-Facing Surface Completeness](https://github.com/ActiveMemory/ctx/blob/main/.context/CONVENTIONS.md):
  the canonical rule the surface audit enforces
- [Detecting and Fixing Drift](context-health.md):
  programmatic drift detection that complements
  judgment-based audits
