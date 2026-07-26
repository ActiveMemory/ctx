# Spec: progressive disclosure for canonical knowledge files

> Design brief: `/ctx-brainstorm`, session 87e465a0, 2026-07-16.
> Decision: DECISIONS.md `[2026-07-16-215955]`.
> Builds on — does not re-litigate — `specs/computed-index-projection.md`
> (`ctx index`, the cheap heading projection, and its Non-Goals).

## Problem

Canonical knowledge files grow without bound, and the entries stay
**valid** — so nothing can be dropped. Time-sharding plus a
load-excluded "cold" bucket is already rejected (a supersession pass
found ~1.5% cold across 162 entries; recency-gating is dangerous because
old ≈ live).

At scale this breaks a real workflow: an agent that legitimately wants
system understanding reads every decision, then every learning, and
exhausts its context window. Two existing levers are insufficient:

- **Consumption discipline** (headings-first via `ctx index`) is
  *necessary but not sufficient*: an agent can always choose to read the
  whole file, and will when it wants completeness.
- **Consolidation** does not help: the 2026-07-16 pass moved LEARNINGS
  only 98 → 88, because the corpus is dense with *distinct signal*, not
  redundancy.

The missing piece is **lossless tiering**: compress history into a
compact top layer, keep every body reachable, and descend only as the
task demands.

## Design

### Three self-similar tiers

- **Tier 0 — the root** (`LEARNINGS.md`, `DECISIONS.md`,
  `CONVENTIONS.md`): *bounded*. Preamble + a **staging zone** + a
  `## Themes` section carrying, per theme, a "just enough" gist and a
  markdown link to its theme file.
- **Tier 1 — theme files** (`.context/learnings/<theme>.md`,
  `.context/decisions/<theme>.md`, `.context/conventions/<theme>.md`):
  the entry bodies for that theme. Reachable **only** via the root's
  links — every artifact is reachable from the canonical file by
  following links, however many hops.
- **Tier 2+ — recursion (deferred)**: an overgrown theme file becomes a
  root in its own right (sub-theme gists + its own staging), handled by
  the same pass. Taxonomy emerges only when the corpus demands it;
  nesting is **not precluded**, just not built.

Reading the root **alone** yields compressed history **+** verbatim
recent delta = a complete current picture, with **no staleness gap**,
because staging *is* the un-digested remainder by construction.

### Gist format (resolves the M2-blocking "just enough" TBD)

Each theme is exactly one bullet line under `## Themes`:

```
- <theme-name> — <one-line gist> → [<theme-name>](<noun>/<slug>.md)
```

- `<theme-name>`: a short kebab-or-words label (e.g. `hooks`, `error
  handling`).
- `<one-line gist>`: **one line, a soft ceiling of ~140 chars**, saying
  what the theme *covers* — the shape of its knowledge, not a list of
  its entries ("hook mechanics: output channels, key names, compliance
  wiring", not "entry A; entry B; entry C"). It conveys *whether to
  drill*, nothing more.
- The separator is the em-dash metadata separator (`token.MetaSeparator`)
  before the gist and ` → ` before the link; the link target is
  `<noun>/<slug>.md` relative to the context dir. This is exactly the
  shape `disclosure.parseThemeBullet` already parses.

The gist is **authored by the pass** (an LLM summarizing the theme's
entries), regenerated whenever the theme gains entries. It is stored
(not recomputed on read) precisely because it is expensive to produce —
the reconciling rationale in DECISIONS `[2026-07-16-215955]`.

### Layout (one shared shape across kinds)

`ctx decision/learning add` are unchanged; `ctx convention add` gains the shared prepend anchor. Verified anchors:

| File | Insert | Staging must sit |
|---|---|---|
| DECISIONS, LEARNINGS | `beforeFirstEntry`: before the first line-start `## [`; falls back to `AfterHeader` when none | **above** `## Themes` |
| CONVENTIONS | `beforeFirstEntry`: before the first line-start `## ` (excluding the structural `## Themes`); falls back to `AfterHeader` when none | **above** `## Themes` |

```
LEARNINGS.md / DECISIONS.md          CONVENTIONS.md
# Learnings                          # Conventions
<!-- UPDATE WHEN … -->               <!-- … -->
## [ts] newest      ← add (STAGING)  ## Newest convention  ← add (STAGING)
## [ts] entry                        ## Older convention
## Themes                            ## Themes
  - hooks — gist… → link               - naming — gist… → link
  - output — gist… → link              - output — gist… → link
```

Because the fallback is `AfterHeader`, an entry lands above `## Themes`
**even when staging is empty**. All three kinds share one layout —
preamble | staging | `## Themes` — differing only in the entry prefix
(`## [` for decisions/learnings, `## ` for conventions) and identity
(timestamp vs section title). Conventions **prepend** newest-first like
the other two (correcting the original `AppendAtEnd`); the anchor skips
the exact `## Themes` string, the only structural `## `. No `## Recent`,
no `###` convention model.

### Consequence: no consumption rewire

Because the root itself is bounded, the existing `ctx agent` packet
("Read These Files: …") becomes safe automatically. No packet change.
The only doc change: the playbook notes "drill into theme files as the
task demands."

### The pass (a new skill — the deliverable)

Agent-driven, human-gated, never inline in another ceremony:

1. Read the staging zone.
2. Propose a theme per staged entry — the agent suggests semantically;
   the human may override, rename, or supply themes.
3. Per target theme: **append body to the theme file → verify
   byte-presence → only then remove from staging**.
4. Regenerate the gist of **every theme it touched**; leave untouched
   themes alone.
5. Create `## Themes` on first run.

### Triggers — suggestion only (milestone 5)

The growth/threshold nudge, `/ctx-remember`, and `/ctx-wrap-up` may
**suggest** the pass. None of them perform it. Wrap-up especially must
stay light: the human is closing the laptop to go live their life, and
semantic work there is against their interest.

All three surfaces read one function — `knowledge.Health(ctxDir,
thresholds) → []Finding` — so the hook (Go) and the skills (prose calling
a `ctx` report path) can never drift. It emits **two signal kinds**,
because "should I fold this?" and "is this too heavy to be useful
context?" are different questions with different units. No state file:
the staging zone and on-disk bytes are self-describing (the watermark
principle from the M1 decision).

**Signal 1 — foldable root (a *count*).** The number of staged entries,
via `disclosure.StagedEntries` (`## [` blocks for decisions/learnings,
`## ` sections for conventions, both excluding the structural
`## Themes`). This reads a never-migrated root correctly — with no
`## Themes`, the whole entry region is staging. Over its per-kind
threshold, the finding suggests **`/ctx-digest`** (fold staging into
themes); `/ctx-consolidate` and `/ctx-drift` are demoted to secondary
("if entries overlap rather than just accrete").

**Signal 2 — heavy page (a *byte count*).** Bytes, not lines: a line
hides 10 or 200 characters, but bytes are the true context cost. Scanned
over the root **and every theme file** under `.context/<noun>/*.md` —
the root-only measure was blind to the bloat folding itself relocates
into theme files. Over the byte ceiling, the finding is advisory: *split
the theme, or extract it to actual tooling.* Past a ceiling, more
Markdown is the wrong fix — an LLM is a poor linter and cannot reliably
apply a large ruleset written as prose, so the remedy prompts the human
to decide whether the content still belongs in an LLM-read file at all.
Tier-2 auto-fold is **not** the heavy-page remedy (it stays deferred; see
Non-Goals).

When both signals fire on one root (a large un-migrated file is both
foldable and heavy), the formatter leads with `/ctx-digest` — folding
reduces both.

**Config** (rc-tunable; `0` disables a check, preserving the existing
convention): per-kind staging counts (`entry_count_learnings` 30,
`entry_count_decisions` 20 carry over unchanged; `convention_line_count`
is **replaced** by `convention_section_count`, since a section count is
the watermark unit — a deliberate behavior change framed as user
education), and a shared `theme_page_byte_ceiling` for signal 2.

The existing `check-knowledge` daily-throttle + snooze plumbing and the
`NudgeBox`/`EmitAndRelay` log-first path are reused, so a user who
declines to fold is not nagged every session.

## Guards

1. **Append → verify → remove.** Never remove-then-append. Any verify
   failure aborts the whole pass with the root untouched.
2. **Precondition validate** (`index.Validate`-style): **zero or one**
   `## Themes`; no entry below it (`## [` for decisions/learnings, `## ` for conventions); staging parses into discrete
   entries. Refuse and fail loud otherwise. **Never regenerate from
   "what I recognized"** — that was the exact root cause of the original
   clobber bug (unparsed content treated as empty).

   Zero `## Themes` means the root is **not yet migrated**: this is the
   first run, and the pass creates the section (see step 5 of the pass).
   Two or more is malformed → refuse. Accepting zero is what keeps
   un-migrated roots passing from day one, so the gate is a signal rather
   than noise that trains people to ignore it. The *invariants* need no
   such carve-out: "no `## [` below `## Themes`" and "gists ↔ theme files
   1:1" are vacuously true on an un-migrated root.
3. **Crash ordering**: theme-file appends (additive) first, then one
   root rewrite. Worst case = duplication (detectable, recoverable),
   never loss.
4. **Fail loud, no auto-repair** — matching the learning-add clobber-fix
   precedent.

## Invariants (mechanically checkable)

- No line-start entry below `## Themes` in a root (`## [` for decisions/learnings, `## ` for conventions; only the exact `## Themes` is structural).
- Root gists ↔ theme files are 1:1 (no orphan file; no gist without a
  file).
- Every theme link resolves (existing `ctx drift` path check).
- Every entry lives in **exactly one** place: staging XOR one theme file.
- Convention staging titles are **unique** — identity is the section title, so duplicates fail loud (`ErrDuplicateStagedTitle`).
- **(M5)** The foldability signal counts the *staging zone*, not the
  whole file: a folded root (entries moved to theme files) reports a
  count below threshold and stops nudging; a never-migrated root reports
  its full entry count.
- **(M5)** The weight scan covers the root **and** every theme file; a
  heavy theme file is flagged even when the root is lean.

## Tests

- **Invariant compliance tests** for each rule above.
- **Conservation**: `staging_before == moved + staging_after`; every
  moved body byte-present in exactly one theme file; zero loss, zero
  dups.
- **`add` still works**: `ctx learning add` with populated *and* empty
  staging both land above `## Themes`; `ctx convention add` **prepends**
  above `## Themes` too (correcting its original `AppendAtEnd`).
- **Abort**: corrupt the root → pass refuses, file byte-identical.
- **Idempotency**: pass with empty staging = no-op.
- **(M5) `knowledge.Health` fixtures**: un-migrated large root →
  foldable finding; migrated root + oversized theme file → heavy finding;
  both at once → foldable leads; all under threshold → no findings.
- **(M5) Boundary**: each threshold fires at `> N`, not `>= N`; `0`
  disables the check.
- **(M5) Convention measure**: a folded 5-section / 250-line convention
  root does **not** foldable-nudge (section count, not lines).
- **(M5) Surface parity**: the `check-knowledge` hook and the skill
  report path produce the same findings from the same root (shared
  `Health`).

## Acceptance

- Each in-scope root stays bounded: gists + links + staging only.
- An agent reading only a root can describe what the corpus says and
  knows where to drill.
- No entry is lost or duplicated across a pass; guards refuse on
  malformed input.
- `ctx decision/learning add` are unchanged. `ctx convention add` moves
  from `AppendAtEnd` to the shared prepend anchor (above `## Themes`),
  correcting the original inconsistency — the one deliberate add-path
  change.
- The pass is codified as a reusable skill.
- **(M5)** The growth nudge, `/ctx-remember`, and `/ctx-wrap-up` each
  surface foldable roots (→ `/ctx-digest`) and heavy pages (→ split /
  extract), from one shared `knowledge.Health`, suggestion-only.
- **(M5)** The foldability signal is staging-based across all three
  kinds; `convention_line_count` is retired for `convention_section_count`.
- **(M5)** Theme files are scanned for weight; nothing auto-folds and no
  state file is introduced.

## Non-Goals

- **No time-sharding, no recency-gating, no load-excluded cold bucket** —
  settled in `specs/computed-index-projection.md`.
- **No change to the decision/learning `add` write path** (the convention path changes: `AppendAtEnd` → prepend above `## Themes`, correcting an original inconsistency).
- **No `ctx agent` packet rewire** — boundedness makes it unnecessary.
- **No taxonomy/nesting machinery now** — the structure is self-similar,
  so nesting is available when themes outgrow their file. **(M5)** The
  heavy-page signal deliberately does **not** auto-recurse (tier-2 fold);
  it advises the human to split or extract-to-tooling, because past a
  byte ceiling the right answer may be "this belongs in a linter, not a
  context file," which no automatic subdivision can decide.
- **CONSTITUTION.md and TASKS.md are out of scope** — the former is small
  by design, the latter is auto-archived.
- **KB pipeline untouched.**

## Phasing (sketch — refine via /ctx-task-out)

1. Guards + invariants + the structural vocabulary (`## Themes`), with
   tests, before any content moves.
2. The pass as a skill, dry-run first (propose themes, move nothing).
3. First real rollout on LEARNINGS (largest corpus), then DECISIONS.
4. CONVENTIONS (curated `## `-section model, unified into the entry-kind
   mover + prepend add-path, edits-behind-a-link UX).
5. **Wire the suggest-only triggers (milestone 5).** One
   `knowledge.Health` emitting foldable (staging count → `/ctx-digest`)
   and heavy (bytes over root + theme files → split/extract) findings;
   consumed by the `check-knowledge` hook and a report path the
   `/ctx-remember` and `/ctx-wrap-up` skills call. Retire
   `convention_line_count` for `convention_section_count`; add
   `theme_page_byte_ceiling`. Suggestion-only, no state file, throttle
   reused. Decompose via `/ctx-task-out --milestone pd-m5`.
