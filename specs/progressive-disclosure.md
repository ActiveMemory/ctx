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

### Layout (forced by the existing write path)

`ctx X add` must not change. Verified anchors:

| File | Insert | Staging must sit |
|---|---|---|
| DECISIONS, LEARNINGS | `beforeFirstEntry`: before the first line-start `## [`; falls back to `AfterHeader` when none | **above** `## Themes` |
| CONVENTIONS | `AppendAtEnd` (EOF) | **below**, in a trailing `## Recent` |

```
LEARNINGS.md / DECISIONS.md          CONVENTIONS.md
# Learnings                          # Conventions
<!-- UPDATE WHEN … -->               <!-- … -->
                                     ## Themes
## [ts] newest      ← add (STAGING)     - naming — gist… → link
## [ts] entry                           - output — gist… → link
## Themes                            ## Recent
  - hooks — gist… → link                ### new convention  ← add (STAGING)
  - output — gist… → link
```

Because the fallback is `AfterHeader`, an entry lands above `## Themes`
**even when staging is empty**. CONVENTIONS needs the explicit
`## Recent` heading because `###` prose sections would otherwise nest
ambiguously *inside* `## Themes`.

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
5. Create `## Themes` (and `## Recent` for conventions) on first run.

### Triggers — suggestion only

The growth/threshold nudge, `/ctx-remember`, and `/ctx-wrap-up` may
**suggest** the pass. None of them perform it. Wrap-up especially must
stay light: the human is closing the laptop to go live their life, and
semantic work there is against their interest.

## Guards

1. **Append → verify → remove.** Never remove-then-append. Any verify
   failure aborts the whole pass with the root untouched.
2. **Precondition validate** (`index.Validate`-style): exactly one
   `## Themes`; no `## [` below it; staging parses into discrete
   entries. Refuse and fail loud otherwise. **Never regenerate from
   "what I recognized"** — that was the exact root cause of the original
   clobber bug (unparsed content treated as empty).
3. **Crash ordering**: theme-file appends (additive) first, then one
   root rewrite. Worst case = duplication (detectable, recoverable),
   never loss.
4. **Fail loud, no auto-repair** — matching the learning-add clobber-fix
   precedent.

## Invariants (mechanically checkable)

- No line-start `## [` below `## Themes` in a root.
- Root gists ↔ theme files are 1:1 (no orphan file; no gist without a
  file).
- Every theme link resolves (existing `ctx drift` path check).
- Every entry lives in **exactly one** place: staging XOR one theme file.

## Tests

- **Invariant compliance tests** for each rule above.
- **Conservation**: `staging_before == moved + staging_after`; every
  moved body byte-present in exactly one theme file; zero loss, zero
  dups.
- **`add` still works**: `ctx learning add` with populated *and* empty
  staging both land above `## Themes`; `ctx convention add` lands inside
  `## Recent`.
- **Abort**: corrupt the root → pass refuses, file byte-identical.
- **Idempotency**: pass with empty staging = no-op.

## Acceptance

- Each in-scope root stays bounded: gists + links + staging only.
- An agent reading only a root can describe what the corpus says and
  knows where to drill.
- No entry is lost or duplicated across a pass; guards refuse on
  malformed input.
- `ctx decision/learning/convention add` work unchanged, with zero code
  edits to the add path.
- The pass is codified as a reusable skill.

## Non-Goals

- **No time-sharding, no recency-gating, no load-excluded cold bucket** —
  settled in `specs/computed-index-projection.md`.
- **No change to the `add` write path.**
- **No `ctx agent` packet rewire** — boundedness makes it unnecessary.
- **No taxonomy/nesting machinery now** — the structure is self-similar,
  so nesting is available when themes outgrow their file.
- **CONSTITUTION.md and TASKS.md are out of scope** — the former is small
  by design, the latter is auto-archived.
- **KB pipeline untouched.**

## Phasing (sketch — refine via /ctx-task-out)

1. Guards + invariants + the structural vocabulary (`## Themes`,
   `## Recent`), with tests, before any content moves.
2. The pass as a skill, dry-run first (propose themes, move nothing).
3. First real rollout on LEARNINGS (largest corpus), then DECISIONS.
4. CONVENTIONS (prose model, `## Recent`, edits-behind-a-link UX).
5. Wire the suggest-only triggers.
