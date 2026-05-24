# Out-of-Band Audit Channel

A discipline-enforcement channel for findings that an
in-session agent reliably misses: an out-of-band auditor
(run by a different Claude Code session, human-triggered)
drops structured reports into `.context/audit/`, and a
verbatim-relay hook surfaces them to the next interactive
session before it picks up momentum.

## Problem

The pad-undo Phase 1 commit (2026-05-24, `6bcaf889`) shipped
a new user-facing CLI subcommand (`ctx pad undo`) with no
matching SKILL.md or recipe updates. The agent labeled the
gap "Phase 2" — which is exactly the deferral the
Constitution forbids ("I can create a follow-up task"). The
gap was caught only because the user noticed.

The failure mode generalizes:

- Mid-task tunnel vision suppresses recall of
  cross-cutting rules. CONVENTIONS.md is loaded at session
  start but doesn't survive contact with feature focus.
- Adding more advisory prose to CONVENTIONS.md doesn't
  help. The CONVENTIONS rule added by the follow-up commit
  (`77fa01f9`) is itself susceptible to the same tunnel
  vision.
- Programmatic test gates (`internal/audit`,
  `internal/compliance`) catch what's mechanical (an
  exported func missing a doc section, a SubcommandDrift
  drift) but cannot judge *which* recipe should mention a
  new flag, or *whether* a captured Learning is materially
  useful. That's a judgment call.
- In-band AI gates (auditor running synchronously inside
  the commit / push flow) work but add latency and API
  cost to every commit. Wrong tradeoff for a discipline
  reminder.

The one discipline pattern in this codebase that
empirically survives agent tunnel vision is **verbatim
relay**: the system reminders that wrap every prompt
("┌─ Journal Reminder ─" boxes etc.) get echoed by the
agent without filtering, every turn, because the relay
bypasses agent judgment entirely. That channel works.

## Approach

Combine three proven patterns:

1. **A different agent does the audit.** An out-of-band
   Claude Code session — triggered manually by the user
   on their plan-billed Claude Code, not in-band on the
   working session's API — runs an audit skill. Fresh
   context means no implementer bias. Manual trigger
   means no per-commit API cost.

2. **The audit drops a report.** Reports land at a known
   path (`.context/audit/<kind>.md`), one file per audit
   kind. Format is machine-parseable header + verbatim-
   relay-shaped body. Multiple audit kinds coexist
   (surface coverage, spec-trailer truthfulness, captured
   Decisions/Learnings coverage, etc.) under the same
   directory and same hook.

3. **A hook relays the report verbatim.** A
   UserPromptSubmit hook (`ctx system checkaudit`) reads
   `.context/audit/*.md`, filters by *not-dismissed AND
   not-stale*, and emits each unprocessed report inside
   the standard verbatim-relay envelope (the
   `┌─ Title ─ ... └─` box). The hook does not interpret
   findings; it relays the body the audit wrote.

Dismissal is explicit (`ctx audit dismiss <id>`) so the
relay stops once the user/agent has acted, mirroring the
proven `ctx remind` pattern.

## Behavior

### Happy Path

1. User finishes a feature on `feat/foo` (commits, lint,
   tests).
2. User opens a *different* Claude Code session in the
   same project worktree and invokes
   `/ctx-surface-audit`.
3. The skill (fresh context, no implementer bias)
   compares the branch against `main`, identifies
   user-facing surfaces added/changed (new subcommands in
   `commands.yaml`, new flags, new behavior), checks
   coverage against SKILL.md command-mapping tables and
   recipe tables.
4. Skill writes `.context/audit/surface.md` with status
   `findings` + a per-surface verbatim body listing the
   specific files to edit.
5. User returns to the original working session (or a
   new one). The UserPromptSubmit hook
   `ctx system checkaudit` reads `.context/audit/`,
   detects the unprocessed `surface.md`, emits a
   verbatim-relay box at the top of the next response.
6. Agent (per the verbatim-relay invariant) echoes the
   box and prompts the user/itself to act.
7. User addresses the findings and runs
   `ctx audit dismiss surface` (or the relay-box footer
   command). The hook goes silent until the next audit
   run.

### Auto-Dismissal on Resolution

When a follow-up commit changes the surfaces named in the
audit body (e.g. SKILL.md edited for the cited
subcommand), `ctx system checkaudit` (or a separate
`checkauditstale` hook) marks the audit dismissed
automatically. Implementation: digest the surface list
into the report header; on hook fire, re-derive the
current state and compare; mark dismissed when the gap is
closed.

Auto-dismissal is a Phase 2 nicety. Phase 1 ships manual
dismissal only.

### Empty Audit Run

When the audit finds no issues, the skill writes
`.context/audit/<kind>.md` with status `clean` and an
empty body. The hook treats `clean` as no-op (no relay).
A `clean` report is still useful as evidence that an
audit ran — `ctx audit list` shows it with a timestamp.

### Stale Report

A report whose `commit-range` no longer matches the
current branch tip is *stale*: it audited code that has
since moved. The hook either:

- Suppresses the relay AND warns once that the report is
  stale (asking the user to re-run the auditor), OR
- Relays anyway with a "STALE — audited at <ref>, branch
  now at <newer-ref>" prefix so the user can decide.

The second option is friendlier; Phase 1 ships that.

### Conflict-of-Interest Guard

The audit skill must refuse to run inside a session that
has uncommitted working-tree changes to the audit target
range, because that session is the implementer and the
point of the channel is fresh-context judgment. Refuse
loudly: "Run this audit from a separate session (commit
or stash here first)." This is the spec's central trust
boundary.

## Interface

### Skill

`/ctx-surface-audit` (the first audit-family skill):

- Inputs: optional git ref range (default: `main..HEAD`).
- Conflict-of-interest guard: refuse if `git status` is
  non-clean.
- Scans diff for: new entries in
  `internal/assets/commands/commands.yaml`, new flag
  registrations, new files under `internal/cli/*/cmd/`,
  new exported behavior visible from CLI help.
- For each surface, checks coverage in:
  `internal/assets/claude/skills/<ctx-area>/SKILL.md`
  (command-mapping table + Execution section), then
  `docs/recipes/*.md` (commands table + walked-through
  steps), then `docs/cli/<command>.md` if a per-command
  page exists.
- Writes `.context/audit/surface.md` with header +
  per-surface findings body.

The skill family scaffolding is generic. Future siblings
under the same scaffolding:
`/ctx-spec-trailer-audit` (does each commit's `Spec:`
trailer point at a spec that genuinely covers that
commit's scope?), `/ctx-capture-audit` (was a Decision
or Learning persisted for non-trivial work that ended
without one?), etc.

### CLI

```
ctx audit list                 # show all reports + status + age + dismissed-state
ctx audit show <id>            # print one report's body
ctx audit dismiss <id>         # mark a report dismissed; hook stops relaying it
ctx audit dismiss --all        # dismiss everything in the audit dir
```

`<id>` is the report basename without extension (e.g.
`surface` for `.context/audit/surface.md`).

The CLI is a sibling of `ctx remind`. Same shape, same
mental model. Wired through `internal/cli/audit/{cmd,
core}/` to match the rest of the CLI taxonomy.

### Hook

`ctx system checkaudit` — UserPromptSubmit hook.

- Reads `.context/audit/*.md`.
- For each report: parse header, check `status` (skip
  `clean`), check `dismissed` flag in
  `.context/state/audit-dismissed.json` (or similar),
  check staleness against `git rev-parse HEAD`.
- Emit the verbatim-relay box with the report's body
  prepended by `┌─ <Kind> Audit ─` and appended with
  `Dismiss: ctx audit dismiss <id>`.
- Same throttling/quietude rules as
  `ctx system checkreminder`.

## Report Format

```yaml
---
kind: surface
status: findings   # findings | clean
commit-range: main..6bcaf889
generated-at: 2026-05-24T14:30:12Z
generator: /ctx-surface-audit
digest: a4c1b7e2  # for auto-dismissal staleness detection
---
Commit 6bcaf889 added user-facing surface without docs:

  • New subcommand `ctx pad undo`
    - SKILL.md: no row in internal/assets/claude/skills/ctx-pad/SKILL.md
      command-mapping table
    - Recipe: docs/recipes/scratchpad-with-claude.md unchanged

Fix:
  - edit internal/assets/claude/skills/ctx-pad/SKILL.md
  - edit docs/recipes/scratchpad-with-claude.md
```

Body is rendered as-is inside the relay box. No
post-processing.

## State

- `.context/audit/<kind>.md` — reports, one per kind.
- `.context/state/audit-dismissed.json` — dismissal
  ledger (id → dismissed-at-commit). Reset when a fresh
  audit overwrites a report (new digest = fresh state =
  not dismissed).

## Files to Create / Modify

Phase 1 scope:

- `internal/assets/claude/skills/ctx-surface-audit/SKILL.md`
  — skill prose (auditor instructions, refuse on dirty
  worktree, write the report).
- `internal/cli/audit/` — new CLI package
  (`audit.go`, `cmd/{list,show,dismiss}/`,
  `core/{store,parse,format}/`).
- `internal/cli/system/cmd/checkaudit/` — UserPromptSubmit
  hook (parallel to `checkreminder`).
- `internal/config/audit/` — path constants
  (`.context/audit/`, `.context/state/audit-dismissed.json`,
  staleness threshold).
- `internal/err/audit/`, `internal/write/audit/` — error
  constructors + writer functions, plus i18n yaml
  additions (errors.yaml + write.yaml) per CONVENTIONS.
- `internal/config/embed/cmd/audit.go` — `UseAudit*` +
  `DescKeyAudit*` constants.
- `internal/assets/commands/{commands,examples}.yaml` —
  `audit.list`, `audit.show`, `audit.dismiss` entries.
- `internal/cli/system/system.go` — register
  `checkaudit` subcommand.
- `.claude/settings.local.json` (project-level) — wire
  the hook as a UserPromptSubmit handler.
- Tests for: CLI commands (list/show/dismiss),
  hook-format output (verbatim box shape matches the
  remind hook), report parser, dismissal-state
  persistence, refuse-on-dirty-worktree skill guard.

Phase 2 (separate spec / commit):

- Auto-dismissal on detected resolution (digest the
  surface list; re-derive on hook fire; suppress when
  closed).
- Additional audit-family skills:
  `/ctx-spec-trailer-audit`, `/ctx-capture-audit`.
- Stale-report graceful escalation.

## Testing

- `TestSurfaceAuditSkill_RefusesOnDirtyTree` — invoke
  the skill in a worktree with uncommitted changes,
  assert refusal message.
- `TestAuditListShowsFindings` — drop a report file
  manually, run `ctx audit list`, assert output includes
  kind, status, age.
- `TestAuditDismiss_StopsRelay` — drop a report,
  dismiss, run the hook, assert no relay emitted.
- `TestAuditDismissAll` — multiple reports, `--all`
  dismisses every one.
- `TestCheckAuditHook_RelaysVerbatimBox` — drop a
  report with a known body, run the hook, assert the
  output is the exact body wrapped in the standard
  relay envelope.
- `TestCheckAuditHook_SilentOnCleanReport` — drop a
  `status: clean` report, assert hook emits nothing.
- `TestCheckAuditHook_StalenessNotice` — drop a report
  whose commit-range no longer matches HEAD, assert
  relay includes the STALE prefix.

## Non-Goals

- **No in-band / pre-commit audit.** The whole point is
  decoupling from the working-session commit cadence.
- **No automated trigger** (Phase 1). Cron / post-commit
  hook / file-watcher are future work. Manual invocation
  via a separate Claude Code session is the entry path.
- **No remediation by the auditor.** The auditor writes a
  report; it does not edit code. The in-session agent
  (or human) does the fixes.
- **No mandatory schedule.** A user who never runs the
  auditor gets no nags — there is no enforcement of
  audit-cadence itself. This is a choice: the channel
  exists for users who want it.
- **No multi-tenant report queue.** One report per kind.
  Re-running an audit overwrites the prior report of the
  same kind. History (if wanted) lives in the dismissal
  ledger or a Phase 2 archive directory.

## Open Questions

1. **Naming collision with `internal/audit/`.** That
   package is an internal-only AST-based-tests package
   (no CLI surface, not imported by anything). The new
   user-facing `ctx audit` CLI and `.context/audit/`
   directory don't collide in any compilable sense, but
   a future reader might confuse them. Mitigation
   options: (a) accept the colloquial reuse, document
   the distinction in `.context/audit/README.md`; (b)
   rename the new channel (`notice`, `inbox`, `relay`,
   `dispatch`); (c) rename the existing internal
   package. Defer to the user.

2. **Where do dismissals live?** Phase 1 design assumes
   `.context/state/audit-dismissed.json`. The
   `.context/state/` directory is per-session tombstone
   space; a cross-session dismissal record might want
   `.context/audit/.dismissed.json` instead (kept with
   the reports). Slight preference for the
   audit-adjacent file so a user nuking `.context/state/`
   does not silently re-surface old findings.

3. **Should the skill family share a common
   scaffolding library?** Phase 2 question. If
   `/ctx-surface-audit`, `/ctx-spec-trailer-audit`,
   `/ctx-capture-audit` all share dirty-tree guard,
   report-write-format, header-shape, etc., a shared
   `internal/audit/skill-helpers/` package becomes
   worthwhile. Phase 1 keeps it per-skill.

## Source

User request, 2026-05-24 session (follow-on to
`feat/pad-undo-snapshot` Phase 1 doc-gap fixup): in-band
discipline rules in CONVENTIONS.md don't survive agent
tunnel vision; the verbatim-relay channel does. Move
discipline enforcement onto the relay channel via
out-of-band-skill → drop-report → hook-relay,
human-triggered for now (no per-commit API cost), with
the same Claude Code plan billing the user already pays
for.
