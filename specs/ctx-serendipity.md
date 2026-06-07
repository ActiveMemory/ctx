# ctx-serendipity (the garden walk)

*Companion to `specs/ctx-dream.md`. The dream **proposes**; serendipity is
the **human gate** that reviews those proposals and bridges accepted ones
into tracked artifacts. Split from ctx-dream per the resolved decision
(session 2263caef, 2026-06-07). v1 reviews discipline-mode proposals over
`ideas/`.*

## Problem

The dream emits atomic, provenance-bearing proposals into the gitignored
`dreams/<ts>/` notebook, but — by construction (Option B) — it never acts
on them. Something must bring a human to those proposals and turn
accept/reject/amend decisions into real outcomes (archive, promote to
spec, mark for blog, merge) without ever letting an ungated machine write
canonical memory.

The failure mode to design against is the one the dream brief named:
**the human doesn't show up, proposals rot, and the backlog is reborn as
a queue.** The author already skips chore-shaped maintenance. So the
review cannot feel like a queue to drain — it has to feel like *walking
the garden*: a small, browsable surface, per-entry attention as pleasure,
no completion pressure.

## Approach

A skill, `/ctx-serendipity` (the "garden walk"), that drives
`ctx dream review` over the current proposal set. It is the **human gate**
in the dream architecture and the only sanctioned path by which a dream
proposal becomes a tracked change.

Two locked principles carried from ctx-dream:

- **Same proposals, two consumers, two interfaces.** The CLI exposes a
  terse, action-coded worklist (the agent's view); the skill renders each
  proposal substance-forward for the human — generated summary +
  provenance + "why now" — so the human never file-hunts.
- **Mechanical vs generative dispositions.** `archive` / `keep` /
  `mark-blog` / `reject` apply **instantly, no LLM cost**. `merge` /
  `promote` are **generative** — they drop to the agent, which reads the
  **full source** (never the lossy summary) to draft the spec or merged
  note. `promote → specs/` is the one deliberate declassification across
  the don't-leak boundary; everything else stays inside `ideas/`/`dreams/`.

Every decision is appended to the dream **ledger** (the shared contract),
so a rejected proposal is not re-surfaced unless its source changes
(dedup-against-*seen*, not against-accepted).

## Behavior

### Happy Path

1. **Nag.** `ctx remind` surfaces "a serendipity round is waiting" at
   session start / every N turns (the proven channel).
2. **Open the walk.** The user triggers `/ctx-serendipity`; it reads the
   committed proposal set from `dreams/<ts>/` via `ctx dream review`.
3. **Browse.** Each proposal is shown substance-forward: generated
   summary, `status` + recommended `action`, `evidence` (commit/spec/
   near-neighbor), `confidence`, one-line `rationale`, and "why now".
4. **Decide per entry:** **accept / reject / amend / skip** — no pressure
   to clear the set.
5. **Apply.** Mechanical reactions run instantly (`ctx dream accept`
   resolves `archive`/`keep`/`mark-blog`, `reject` records). Generative
   ones drop to the agent: `promote` drafts `specs/<name>.md` via
   `/ctx-spec` from the full source; `merge` reads the full sources and
   writes the merged `ideas/` note (backup-before-mutate first).
6. **Record.** Every decision appends to `dreams/ledger.md`; the source's
   `status`/`last_surfaced`/`history` update so nothing rots or re-nags.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| No pending proposals | The walk reports "garden's quiet — nothing waiting" and exits; no empty ritual. |
| User skips an entry | No decision recorded; it remains pending and may re-surface next round (not a rejection). |
| `amend` changes the action (e.g. promote→keep) | The amended action is what's applied + logged; provenance of the original proposal is preserved in the ledger. |
| Accept `promote` but `/ctx-spec` fails | Surface the error; leave the idea in `ideas/` untagged; ledger records the attempt, not a success. |
| Accept `merge`/destructive but backup fails | Abort the mutation (backup-before-mutate precondition); item stays pending, surfaced again. |
| A dream pass runs while a review is open | The review reads a committed proposal set; a new pass is serialized by the dream's lock and doesn't disturb the open walk. |
| Proposal's source idea changed since the proposal was generated | Flag as stale in the walk; prefer re-triage over acting on a stale summary. |
| Same proposal already decided in a prior round | Ledger dedup — not re-surfaced unless the source content changed. |

### Validation Rules

- Accept/reject/amend operate by stable proposal `id`; an unknown id is
  rejected with a clear error (no silent no-op).
- `promote` is the **only** action permitted to write a tracked path
  (`specs/`); all other writes stay within gitignored `ideas/`/`dreams/`
  (the don't-leak guard still applies).
- A generative disposition must read the **full source**, not the cached
  summary, before writing.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| No proposal set found | `no pending dream proposals` | run a dream pass (`ctx dream`) or wait for cron |
| Unknown proposal id | `unknown proposal: <id>` | re-list with `ctx dream review` |
| `/ctx-spec` fails on promote | surface the `/ctx-spec` error verbatim | retry promote, or amend to `keep` |
| backup failed before merge | `[dream] backup failed; skipping merge for <file>` | item left untouched; re-surfaced next round |

## Interface

### CLI

Serendipity drives the primitives defined in `specs/ctx-dream.md`:

```
ctx dream review                       # interactive ~15-min round
ctx dream accept|reject|amend <id>     # primitives the skill (and agent) drive
```

### Skill

```
/ctx-serendipity   (a.k.a. the "garden walk")
```

Trigger phrases: "serendipity round", "review my dreams", "walk the
garden", "what did the dream find?". Sibling to `/ctx-remember`,
`/ctx-wrap-up`.

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `skills/ctx-serendipity/` (tracked) | The garden-walk skill: render proposals substance-forward, drive accept/reject/amend, route generative items |
| `internal/cli/dream/` | `review`/`accept`/`reject`/`amend` subcommands (shared with ctx-dream) |
| `ctx remind` wiring | The "serendipity round waiting" nag cadence |

### Helpers to Reuse

- `ctx dream review|accept|reject|amend` (the dream's CLI primitives).
- `/ctx-spec` (promote → spec), `/ctx-blog` (mark-blog → draft), the
  archive convention (`ideas/done/`).
- `ctx remind` for the nag cadence — reuse, do not reinvent.
- The dream `ledger` + proposal schema (the shared contract).

## Configuration

- Uses the same `dream:` `.ctxrc` section as ctx-dream (notably the
  `ctx remind` cadence/wording). No separate config surface.

## Testing

- **Unit:** accept resolves a mechanical action with no LLM call; reject
  appends to the ledger and the item does not re-surface; amend changes
  the applied action while preserving original provenance; unknown id
  errors.
- **Integration:** a fixture proposal set → `accept promote` lands a spec
  in tracked `specs/` and tags the idea; `reject` → ledger → absent next
  round; `merge` reads full source and backs up before mutating.
- **Edge:** empty proposal set exits cleanly; stale-source proposal is
  flagged not acted on.

## Non-Goals (v1)

- **No autonomous application.** Every disposition into a tracked
  artifact passes through the human in this skill; serendipity is the
  gate, never an auto-approver.
- **No creative/garden-mode resurfacing** — that's the deferred creative
  dream mode; v1 serendipity reviews discipline proposals only.
- **No new proposal generation** — serendipity only reviews/acts on what
  the dream produced; it does not classify or ground (that's the dream's
  job).
- **No web UI** — the CLI + skill are the surface.
