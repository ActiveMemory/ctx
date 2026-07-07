# Spec: self-healing journal import

## Problem

`ctx journal import --all` skips any session whose journal file
already exists. That default is safe for static archives but wrong
for the one thing transcripts actually do: grow.

Two concrete failures fall out of it:

1. **Permanent truncation.** Import a session while it is still
   running (or before its final messages flush) and the journal
   entry freezes at whatever prefix existed at import time. The
   skip-existing default guarantees it never completes; the only
   cure is remembering `--regenerate`, which reintroduces exactly
   the "remember to run this" burden the journal pipeline exists
   to remove.
2. **No safe ceremony hook.** Because a mid-session import
   poisons the entry forever, journal import cannot be wired into
   `/ctx-wrap-up` or a session-end hook without special-casing
   the active session — an exclusion flag the user must know
   about, one more thing to think about.

The root cause is that import's unit of memory is the *output
file* ("does it exist?") when it should be the *source
transcript* ("has it changed since I last rendered it?").

## Key observation

Claude Code JSONL transcripts are **append-only**. A transcript
never rewrites history; it only gains lines. Therefore:

- "source grew" is detectable from mtime + size alone — no
  hashing, no diffing;
- a partial import is not a mistake to avoid but an intermediate
  state that a later import completes;
- re-rendering a grown session reproduces everything the earlier
  render produced, plus the new tail.

Truncation stops being a hazard the moment import becomes
growth-aware. No active-session exclusion is needed anywhere —
the live session imports as far as it goes today and self-heals
on the next sweep.

## Design

### 1. Source tracking in journal state (schema v2)

`internal/journal/state` gains a per-session map alongside the
existing per-file stage entries:

```json
{
  "version": 2,
  "entries": { "<filename>.md": { "exported": "…", "render_hash": "…" } },
  "sessions": {
    "<session-id>": {
      "source_file": "<abs path to chosen transcript>",
      "source_mtime": 1234567890,
      "source_size": 123456
    }
  }
}
```

- Keyed by **session id**, not source path: one session can span
  multiple JSONL files (resume copies prior history into a new
  file). The import unit is the session, rendered from the
  richest contributing transcript; the chosen file switching to
  a larger resume copy also counts as growth.
- v1 files load tolerantly (missing maps initialise empty). The
  first v2 run performs a one-time adoption pass: record current
  source stats for already-imported sessions *without*
  re-rendering, so growth detection starts from the next real
  change.

### 2. Growth-aware planning: New / Grown / Unchanged

`plan.Import` replaces its exists-check with a three-way decision
from state:

| Decision  | Condition                                | Action              |
|-----------|------------------------------------------|---------------------|
| New       | session id absent from state.sessions    | render all parts    |
| Grown     | recorded mtime or size differs           | re-render affected  |
| Unchanged | recorded mtime and size match            | skip                |

Part-awareness makes Grown cheap: growth only *appends* messages,
so earlier parts' message ranges are unchanged by construction —
re-render only the last part plus any newly created parts.

`--regenerate` survives as an explicit edge-case tool (mass
re-render after a render-format change, healing pre-v2 truncated
entries whose sources will never grow again). It stops being the
routine healing path. **No new flags.**

### 3. Foreign-edit safety: the render hash

Journal entries are deliberately editable ("add notes, highlight
key moments, clean up the transcript"), and enrichment /
normalization legitimately rewrite them after import. Auto-updates
must never clobber a hand-edited file.

Invariant: **`render_hash` always reflects the last ctx-authored
write of the file.** Every pipeline stage that writes the entry
(import, enrich, normalize, fence-verify) refreshes the hash in
state. On Grown:

- hash matches file → the file is provably ctx-owned → splice the
  fresh transcript (preserving enriched frontmatter via the
  existing keep-frontmatter machinery);
- hash differs → someone edited it → leave the file untouched,
  warn, and suggest `ctx journal lock` (permanent protection) or
  an explicit `--regenerate` (deliberate discard).

Locked entries are never rewritten (existing invariant, unchanged).
Files with no recorded hash (pre-v2) are treated as edited —
warn, never clobber.

### 4. Live-transcript parse tolerance

A transcript read mid-write may end in a partial JSONL line. The
parser treats a trailing incomplete line as end-of-input (truncate
to the last complete line) rather than an error or warning. The
missing record is captured on the next sweep — growth-awareness
makes the cut lossless.

### 5. Delivery: hook first, ceremony as backup

- `/ctx-wrap-up` runs `ctx journal import --all -y` as a
  best-effort, non-blocking step before delegating to
  `/ctx-handover`. A failed import never blocks the handover.
- A **SessionEnd hook** in hooks.json fires the same import
  automatically on the way out of every session. It wires
  `ctx journal import --all -y` directly (output discarded,
  `|| true` so a failure never blocks session teardown),
  mirroring the existing `ctx agent` hook rather than adding a
  redundant `ctx system` passthrough verb: the hooks-wiring
  compliance guard resolves any `ctx <path>`, so
  `ctx journal import` is covered without a new command, and
  the ceremony and hook then share one code path. The ceremony
  step stays as belt-and-suspenders: the import is idempotent,
  so redundancy costs one stat per session.
- Enrichment stays out of both paths: it is an LLM pass and
  belongs to `/ctx-journal-enrich-all`.

## Decisions

- **No active-session exclusion.** Growth-awareness makes the
  live session safe to import at any moment; an exclusion flag
  would be a workaround for a problem this design removes. One
  less flag to think about.
- **Render hash over owned-region markers.** Marker-delimited
  machine-owned regions were considered and rejected: ctx's
  journal contract invites free-form edits anywhere in the file,
  and a marker scheme would require migrating the existing corpus
  and narrowing that contract. The hash gives the same guarantee
  (never clobber human work) with zero file-format change; `lock`
  remains the explicit protection for curated entries.
- **mtime+size, no content hashing of sources.** Append-only
  sources make byte-level comparison redundant; stat is free and
  runs on every sweep.
- **Pre-v2 truncated entries do not auto-heal.** Their sources
  will never grow again, so adoption records current stats and
  moves on; a documented one-time `--regenerate` heals them.
  Silent mass-regeneration on upgrade was rejected as a clobber
  risk (no render hash exists yet to prove the files unedited).

## Acceptance criteria

- Importing a live session, then importing again after it grows,
  yields a complete entry with no flags involved.
- `ctx journal import --all` twice in a row: second run reports
  all Unchanged, writes nothing (byte-identical journal dir).
- A hand-edited entry whose session grows is left byte-identical,
  with a warning naming the file and both escape hatches.
- Locked entries are never rewritten under any decision.
- Enriched frontmatter survives a Grown re-render.
- A transcript ending in a half-written JSONL line imports its
  complete lines without error or warning.
- v1 state files load, adopt, and round-trip to v2 losslessly.
- `/ctx-wrap-up` and the SessionEnd hook both trigger the sweep;
  hook wiring passes the shipped-hooks compliance guard.

## Tasks

Tracked in TASKS.md under **Phase JI: Self-Healing Journal
Import**.
