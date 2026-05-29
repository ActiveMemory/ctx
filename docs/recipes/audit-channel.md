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

!!! note "Maintainer tooling: lives in `ctxctl`, not the shipped `ctx` binary"
    `ctxctl audit` and the `ctxctl audit-relay` hook are the
    **generic relay** half: a place for *any* out-of-band tool
    to drop a report and have it relayed. They live in
    `ctxctl` — ctx's separate maintainer/contributor binary —
    **not** in the user-facing `ctx` binary, so end users
    never carry an audit hook they have no producer for. The
    **auditor** that produces the report is project-specific —
    it must know *your* conventions and directory layout. ctx
    dogfoods its own internal auditor (`_ctx-surface-audit`,
    a repo-only skill that scans ctx's `internal/` tree); the
    examples below use it as a concrete reference. To adopt
    the pattern in your own project, build the same
    out-of-band relay plus your own audit skill. ctx
    maintainers build and install `ctxctl` once with `make
    reinstall-ctxctl` (→ `/usr/local/bin/ctxctl`); every
    worktree then shares the one binary.

## TL;DR

```bash
# 1. From a separate Claude Code session, run your project's
#    audit skill (ctx's own internal example shown here):
/_ctx-surface-audit               # default: main..HEAD

# 2. It writes a structured report:
.context/audit/surface.md

# 3. Back in the working session, the next prompt fires the
#    repo-local UserPromptSubmit hook (wired in
#    .claude/settings.local.json, not by `ctx setup`):
ctxctl audit-relay

# 4. The agent / human sees a verbatim-relay box on the next
#    response, listing the specific findings.

# 5. After addressing the findings:
ctxctl audit dismiss surface      # stops the relay
```

## Commands and Skills Used

| Tool                        | Type        | Purpose                                                |
|-----------------------------|-------------|--------------------------------------------------------|
| `ctxctl audit list`         | CLI command | Show all reports with status and age                   |
| `ctxctl audit show ID`      | CLI command | Print one report's body, pipe-friendly                 |
| `ctxctl audit dismiss ID`   | CLI command | Mark a report dismissed against its current digest     |
| `ctxctl audit dismiss --all`| CLI command | Bulk dismissal                                         |
| `ctxctl audit-relay`        | CLI command | UserPromptSubmit hook; verbatim-relays reports         |
| `_ctx-surface-audit`        | Skill       | ctx's **own internal** auditor — reference example, not bundled |

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

The `/_ctx-surface-audit` skill enforces this with a hard
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
You (in session 2): "/_ctx-surface-audit"

Skill: "Scanned 4 commits, 3 surfaces detected.
        Wrote .context/audit/surface.md (status: findings).
        Open a working session — the audit-relay hook will
        relay the findings on the next prompt."
```

The auditor compares the branch against `main`, finds new
subcommands / flags / behavior changes, checks each one
against SKILL.md / recipe / `docs/cli` coverage, and writes
a structured report.

### Step 3: Return to the Working Session

The next time you submit a prompt in your working session,
the `ctxctl audit-relay` hook (a UserPromptSubmit hook wired
in the repo-local `.claude/settings.local.json` — maintainer-
only; `ctx setup` does not install it) reads
`.context/audit/` and emits a verbatim-relay box at the top
of the agent's response:

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
│ Dismiss: ctxctl audit dismiss <id>
│ Dismiss all: ctxctl audit dismiss --all
└──────────────────────────────────────────────────
```

The agent echoes this verbatim — that is the discipline
mechanism. You (or the agent) then address each cited file.

### Step 4: Dismiss

Once you've addressed the findings (or accepted them as
out-of-scope), dismiss the report:

```bash
ctxctl audit dismiss surface
```

Dismissal is bound to the **report digest** at dismiss
time. A subsequent audit that produces the same findings
stays dismissed. A subsequent audit that finds *new* surface
drift produces a fresh digest and re-surfaces the report at
the next prompt.

## Retention

The audit channel keeps **one report per kind**. Re-running
`/_ctx-surface-audit` overwrites the prior `surface.md`.
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

`_ctx-surface-audit` is the first of a family of ctx's own
internal auditors (all `_`-prefixed, repo-only). The
scaffolding they share — channel, ledger, hook, CLI — lives
in the maintainer-only `ctxctl` binary; the auditors
themselves are repo-only skills. Planned
siblings under the same shape:

- `_ctx-spec-trailer-audit` — does each commit's `Spec:`
  trailer point at a spec that genuinely covers that
  commit's scope?
- `_ctx-capture-audit` — was a Decision or Learning
  persisted for non-trivial work that ended without one?

Each lives in its own SKILL.md and writes its own report
file (e.g. `.context/audit/spec-trailer.md`). The hook
relays whatever it finds, with no per-kind plumbing — which
is exactly what lets *your* project's auditors plug in
without touching ctx.

## See Also

- [Spec: out-of-band audit channel](https://github.com/ActiveMemory/ctx/blob/main/specs/audit-channel.md):
  full design rationale + Open Questions
- [CONVENTIONS → User-Facing Surface Completeness](https://github.com/ActiveMemory/ctx/blob/main/.context/CONVENTIONS.md):
  the canonical rule the surface audit enforces
- [Detecting and Fixing Drift](context-health.md):
  programmatic drift detection that complements
  judgment-based audits
