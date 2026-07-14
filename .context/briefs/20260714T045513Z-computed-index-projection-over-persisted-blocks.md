# Brief: Computed `ctx index` projection replaces persisted INDEX blocks

- **Date**: 2026-07-14
- **Session**: 75be038e | Branch: main @ f382bee7
- **Status**: Debated (via /ctx-plan). Ready for `/ctx-spec --brief`.
- **Supersedes the opening bet**: the original pitch (time-sharded knowledge
  files + a load-excluded "cold" bucket) was **killed by validation** — see
  Rejections.

## The bet

Knowledge files (DECISIONS/LEARNINGS/CONVENTIONS/TASKS) carry hand-maintained
`<!-- INDEX:START -->…<!-- INDEX:END -->` blocks — a stored table of contents.
That stored index is dead weight and pure liability: it drifts, it has caused a
clobber bug and an "entries trapped in the index block" bug, and it is **not even
injected at session start** anymore (the load-gate was narrowed to
CONSTITUTION + AGENT_PLAYBOOK_GATE on 2026-06-07). So the agent only ever sees it
by reading the file top-to-bottom, where it is redundant weight.

An index is just a **projection over the L2 (optionally L3) headings** — a grep,
computable on demand. Maintaining a separate index comment is dead weight.

**Bet:** delete the persisted-index machinery entirely and replace it with a
generic, computed heading projector.

Concretely:
1. **Remove** the stored-index machinery: `internal/index/` (block maintenance),
   the `ctx reindex` command (`internal/cli/reindex/`), `config/marker/index_fmt.go`,
   the `add` insert-path write into the index block, and the now-moot
   `TestWrite_RefusesEntriesTrappedInIndexBlock` guard. Strip the INDEX blocks
   from the files themselves (DECISIONS ×1, LEARNINGS ×2, TASKS ×6). This is
   **net deletion** — on-brand for ctx (cf. the cwd-anchored decision that
   celebrated deleting ~600–1000 LOC).
2. **Add `ctx index <file>`** — a thin, generic projector that greps `##`
   (and optionally `###`) headings and emits them. One mechanism across all
   knowledge files. For TASKS.md it yields the `## Phase` sections — the same
   command doubles as "get all phases/sections" ("multiple birds, one scone").
   The entry headings already are `## [timestamp] Title`, so the grep *is* the
   old `| Date | Decision |` row, for free.
3. **Rewire `ctx agent`** to call the projector, so the agent still gets a TOC —
   now **computed, never stored**, so it can't drift, clobber, or trap entries.

This is Reframe C from the debate: the real hot/cold split is **hot = the fact
that an entry exists (its heading); cold = its body** — orthogonal to whether the
entry is live, so it never risks hiding a live constraint (the agent always sees
every title and pulls bodies on demand).

## Decisions locked in the debate

- **Namespace: `ctx index`, top-level, not `ctx system index`.** A heading TOC is
  genuinely dual-use (human + agent). The "`ctx system` is for plumbing users
  don't type by hand" precedent (2026-04-14 bootstrap decision) is about
  *agent-only* commands; a readable TOC isn't that.
- **Generic projector, not per-noun.** One `ctx index <file>` over any knowledge
  file, confirmed non-colliding: there is **no** existing `ctx decision list` /
  `search` surface today (grep verified), so nothing to reconcile.
- **Ship order: `index` first.** A richer `list/search` (filtering, full-text)
  is a *later* layer on top of the projection primitive. Filed as a follow-up
  task in TASKS.md → Misc (successor to the queued "CLI-projected list/search").

## What we rejected (and what would change our mind)

- **Time-sharded files + a load-excluded "cold" (superseded/deprecated) bucket.**
  Rejected on data. A one-shot supersession-detection pass over the full corpus
  found **~1 genuinely cold entry / 66 decisions and ~1 / 96 learnings** (≤5 even
  in the worst case if code-checkable suspects verify dead). The cold bucket has
  no fuel. *Why the corpus is so live:* it was **already garbage-collected by
  consolidation** — every `consolidated from N entries` block is a tombstone for
  merged-out originals (the June-7 pass took decisions ~109→66, learnings
  ~151→96). The working GC lever in this repo is `/ctx-consolidate`, not
  supersession-exclusion. *Would change our mind:* if a future corpus grew a
  large, genuinely-dead tail that consolidation wasn't run against.
- **Recency-gated selective loading (don't load old entries).** Rejected as
  *more* dangerous after the data, not less: the pass proved **old ≈ live** here
  (e.g. the 4-month-old "path deny-list is a safety net, not a security boundary"
  invariant). Recency-gating would routinely withhold *live constraints* and
  cause the agent to re-litigate settled decisions. *Would change our mind:* a
  reliable liveness signal that doesn't depend on age.
- **A dedicated `/ctx-shard` skill + ceremony wiring** (hook a cleanup pass into
  `/ctx-remember`/`/ctx-wrap-up`, human-acknowledged one-at-a-time with an
  "enough for this session" escape hatch). The mechanism design was *sound* and
  worth remembering — but it exists to feed the cold bucket, which the data
  emptied. Shelved with the cold-bucket bet, not on its own merits.

## Top failure modes

1. **Heading-format hygiene.** The projector greps `##` headings, so a malformed
   heading (missing timestamp, wrong level) silently drops from the index.
   Mitigation: the existing drift-check / kb-site-review already coerce heading
   format; `ctx index` can warn on malformed headings. How we'd notice: an entry
   present in the file but absent from `ctx index` output.
2. **Silent consumer breakage on removal.** Something we didn't find may read the
   INDEX block. Mitigation: grep confirmed the injection gate no longer feeds it;
   `ctx agent` is the one consumer to rewire. Verify no other reader before
   deleting `internal/index/`.
3. **Scope creep into list/search.** The temptation to make `index` filter/query.
   Mitigation: keep `index` a *thin heading grep*; the richer surface is the
   separately-tasked `list/search`.

## Cheapest validation (already done)

- **One agent supersession pass** over both files (zero code) → returned the
  1.5% cold number that killed the shard bet.
- **One grep** → confirmed the `internal/index/` + `ctx reindex` machinery
  exists, is net-deletable, and is no longer session-injected.

Both are in the transcript; no further validation needed before spec.

## Unwind cost

Low. The removal is a net-deletion behind git; the projector is additive. If
`ctx index` proves insufficient, re-add is cheap. No data migration — the
headings that *are* the index already live in the files.

## Next step

`/ctx-spec --brief .context/briefs/20260714T045513Z-computed-index-projection-over-persisted-blocks.md`
