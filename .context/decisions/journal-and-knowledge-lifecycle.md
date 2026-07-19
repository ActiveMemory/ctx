# journal-and-knowledge-lifecycle

## [2026-05-31-094649] Journal screen renders ctx journal source verbatim instead of parsing it

**Status**: Accepted

**Context**: The Journal timeline needs session data, but ctx journal source emits a whitespace-aligned table with no --json mode, and columns (usage/turns) drop out on short sessions.

**Decision**: Journal screen renders ctx journal source verbatim instead of parsing it

**Rationale**: A column parser misaligns when fields are absent; rendering the text verbatim in a monospace panel is honest and robust for P0.

**Consequence**: Journal is a styled text view, not structured cards; a proper timeline awaits a journal source --json mode upstream (same pattern as the artifact list --json commands).

---

## [2026-04-11-200000] Journal stays local; LEARNINGS.md is the shareable layer

**Status**: Accepted

**Context**: With the hub now carrying shared project context between machines
and eventually between teammates, the question came up whether enriched
journal entries should ride along — either the raw `.context/journal/` files
or an "export enriched entries as shareable learning items" pipeline layered
on top of `/ctx-journal-enrich`. The journal is already gitignored per the
2026-03-05 `.context/memory/` decision and for the same reason: it's a
first-person log of raw prompts, half-formed thoughts, dead ends, personal
names, and things the user talks through with themselves. It sits in the
same trust tier as shell history or a private notebook.

The trade-off is real: shared journals would make it trivial for teammates
(or future-me on another machine) to see the full reasoning trail behind a
decision. But "full reasoning trail" is precisely the thing that makes a
journal journal and not a changelog — it includes the parts the author
hasn't decided to stand behind yet, plus incidental private content.

**Decision**: The journal is **Tier-0 personal** and never leaves the
originating machine. No hub sync, no export-by-default, no
enriched-entries-as-shareable-items pipeline. The enrichment pipeline
(`/ctx-journal-enrich`) stays as-is: journal → human-in-the-loop review →
explicit promotion to LEARNINGS.md / DECISIONS.md / CONVENTIONS.md via the
existing `/ctx-learning-add`, `/ctx-decision-add`, `/ctx-convention-add`
commands. Those distilled artifacts are **Tier-1 shareable** and are what
the hub syncs when a team opts into shared context.

The promotion boundary is therefore the enrichment step, not a new export
pipeline. The user is the gate.

**Rationale**: Any "shareable enriched journal entry" pipeline would have to
re-implement the trust boundary that `/ctx-learning-add` already enforces:
the human decides what's worth sharing, strips incidental private content,
and rewrites it as a standalone artifact. A second pipeline that tries to
do this automatically would either (a) leak private content by accident, or
(b) require the same human review and thus collapse back into
`/ctx-learning-add`. The principled answer is that there is no second
pipeline — LEARNINGS.md *is* the shareable form of the journal.

This also preserves the psychological safety of the journal: the author
can write freely because they know nothing they write is one sync away
from a teammate's screen. Lose that property and the journal stops being a
journal and starts being a changelog draft.

**Consequence**:

- Journal files stay gitignored and stay out of `ctx hub` sync paths. Any
  future code that walks context files for replication must exclude
  `.context/journal/` explicitly and be covered by a test.
- `/ctx-journal-enrich` remains the promotion boundary. Its output targets
  are LEARNINGS.md / DECISIONS.md / CONVENTIONS.md, never a separate
  "shareable journal" bucket.
- Hub docs (`docs/home/hub.md`, `docs/recipes/hub-personal.md`,
  `docs/recipes/hub-team.md`, `docs/security/hub.md`) should state the
  Tier-0 / Tier-1 split explicitly so users building team workflows don't
  assume "shared context" means "shared everything."
- The sync code path in `internal/hub/sync_helper.go` and any future
  replication of context files must enforce this exclusion at the
  code level — a gitignore entry is a user-convenience signal, not a
  hub-trust boundary.
- A potential future "personal multi-machine journal sync" (same human,
  different laptops) is explicitly **out of scope** of this decision. If
  it ever ships, it rides a different transport (encrypted-at-rest,
  single-user, not the team hub) and needs its own decision record.

**Alternatives considered**:

- **Sync raw journal files via hub**: rejected. Inverts the gitignore
  decision, leaks private content by construction, destroys the
  journal's "safe to write freely" property.
- **Auto-export enriched entries as a new shareable artifact type**:
  rejected. Duplicates `/ctx-learning-add` without the human gate, or
  collapses back into it. No real difference from the status quo except
  the opportunity for accidental leakage.
- **Opt-in per-entry "publish to hub" flag in the journal**: rejected as
  premature. If the user wants an entry on the hub, the existing flow is
  one command away — write it as a learning or decision. A second path
  adds surface area without adding capability.

**Related**: Reinforces the 2026-03-05 `.context/memory/` gitignore
decision (same trust-tier reasoning for a different private artifact).

## [2026-04-08-013731] Remove #done tag convention, simplify task archival

**Status**: Accepted

**Context**: Tasks had #done:YYYY-MM-DD timestamps that agents added
inconsistently and nobody read. compact --archive filtered by age using these
timestamps.

**Decision**: Remove #done tag convention, simplify task archival

**Rationale**: [x] checkbox is semantically sufficient. git blame provides the
completion timestamp. Removing #done eliminates redundant ceremony and
simplifies compact --archive to archive all completed tasks regardless of age.

**Consequence**: compact --archive no longer filters by archive_after_days for
tasks. The .ctxrc field is inert but retained for backwards compatibility.
Historical #done tags in archives are preserved.

---

## [2026-03-24-001001] Write-once baseline with explicit end-consolidation for consolidation lifecycle

**Status**: Accepted

**Context**: Designing the consolidation nudge hook; multi-pass consolidation
spans dozens of sessions and you cannot programmatically distinguish feature
from consolidation sessions

**Decision**: Write-once baseline with explicit end-consolidation for
consolidation lifecycle

**Rationale**: First ctx-consolidate stamps baseline (write-once), user runs
end-consolidation when done. Failure mode is silence (no stale nudges), not
wrong behavior

**Consequence**: Requires mark-consolidation, end-consolidation, and
snooze-consolidation plumbing commands. Spec: specs/consolidation-nudge-hook.md

---


## [2026-02-26-100004] Task and knowledge management (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-01-27 to 2026-02-18)

- Tasks must include explicit deliverables, not just implementation steps.
  Parent tasks define WHAT the user gets; subtasks decompose HOW to build it.
  Without explicit deliverables, AI optimizes for checking boxes.
- Use reverse-chronological order (newest first) for DECISIONS.md and
  LEARNINGS.md. Ensures most recent items are read first regardless of token
  budget.
- Add quick reference index to DECISIONS.md: compact table at top allows
  scanning; agents can grep for full timestamp to jump to entry. Auto-updated on
  `ctx add decision`.
- Knowledge scaling via archive path for decisions and learnings: follow the
  task archive pattern, move old entries to `.context/archive/`, extend `ctx
  compact --archive` to cover all three file types.

---

