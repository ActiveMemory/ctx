# ctx-dream (v1: disciplined `ideas/` triage)

*Source: debated brief `.context/briefs/20260606T203414Z-ctx-dream-disciplined-consolidator.md` and the decision "ctx-dream: standalone proposing memory consolidator (Option B), human-gated via serendipity" in `DECISIONS.md`. Scoped to **v1 — discipline mode over `ideas/`**. Creative/garden mode and the v2 canonical-consolidation pipeline are **sketched-not-contracted**, re-debatable after v1. A working skill draft + the 28-source research corpus already exist under the gitignored `ideas/ctx-dreams/` (referenced by path only).*

## Problem

`ideas/` is a gitignored personal goldmine — hundreds of markdown notes (plus large binaries) — that has grown **too large to triage by hand**. The author can't tell which ideas are already implemented, which are duplicates, which still have merit, which are throwaway sidenotes, and which are blog material. Good ideas sink and are forgotten; stale ones clutter.

More broadly, ctx memory accumulates over time, and for ctx's eventual users canonical bloat dilutes agent context. But the research is unambiguous that *naive, continuous* LLM consolidation makes memory **worse, not better** (utility can fall below the no-memory baseline). So the need is a **gated, periodic, proposing** process — never an autonomous rewriter — that surfaces a *small, high-signal* set of triage proposals (and, later, serendipitous re-encounters) for a human to act on.

**Why now:** `ideas/` is concretely overwhelming today, and that felt pain is the cheapest place to validate the mechanism before pointing it at higher-stakes canonical memory.

## Approach

A scheduled, standalone **"dream"** — a sleep-time process that runs headless/background and **only ever proposes**. It never autonomously mutates canonical memory (**Option B**). One skill, **two execution modes**: `discipline` (default, the only mode built in v1) and `creative/exploration` (a safe relaxation, deferred). v1 points discipline mode at `ideas/`.

The dream **classifies** each idea (implemented / duplicate / still-meritorious / sidenote / blog-candidate) and proposes a **disposition** (archive / merge / promote-to-spec / mark-blog / keep), grounded against the codebase + specs, semantically deduped, each proposal carrying provenance + evidence + a confidence signal. A `ctx remind` nag pulls the human into a ~15-minute **serendipity** review round to **accept / reject / amend**. Accepted items route to existing skills (`archive`, `/ctx-spec`, `/ctx-blog`). Rejections are recorded in a ledger so future dreams don't re-surface them (dedup-against-*seen*, not against-accepted).

**Where it fits / locked principles:**
- **Decouple the cognition, reuse the plumbing.** The dream owns its triage/synthesis logic; it reuses `ctx remind` and (in v2) import/enrich/kb-ingest via a stable data contract. v1 touches none of that plumbing — it reads raw `ideas/`, classifies, proposes.
- **`dreams/` is the dream's notebook *about* `ideas/`** — a **gitignored, root-level** working area. Ideas never leave `ideas/`. `dreams/` holds derived summaries, per-source state, the ledger, per-run digests, and pre-mutation backups.
- **Three structural safety invariants, enforced by guards (not prompts):**
  1. **Write-scope** — the dream may only write under `dreams/` and `ideas/` (and `specs/` *only* via an accepted `promote`).
  2. **Sources-as-data** — idea text is never executed as instructions; `ideas/` is an indirect-prompt-injection surface even when only read.
  3. **Don't-leak** — privacy class propagates from source to every derived artifact; every dream write target is verified with `git check-ignore` and **refused if it resolves to a tracked path**. The one legitimate boundary crossing is the human's explicit `promote` (deliberate declassification into `specs/`).

## Behavior

### Happy Path

1. **Trigger.** A cron job (or a session-start fallback when cron is stale) runs `ctx dream` in `discipline` mode, bounded to `max N` files per pass (default ~50).
2. **Scan + delta.** The dream lists `ideas/**.md` (excluding its own `dreams/` notebook and large binaries), computes each file's content hash, and selects work by the **discipline clock**: files whose hash is new or changed since last triage (skip unchanged-and-already-dispositioned).
3. **Summarize + ground.** For each selected idea it generates/refreshes a cached summary and grounds a classification against the codebase + specs ("implemented?" → cite the commit/spec; "duplicate?" → cite the near-neighbor).
4. **Propose.** It writes atomic, provenance-bearing proposals into `dreams/<ts>/` (gitignored), practising **ruthless self-rejection** — surfacing a few high-confidence items, not everything.
5. **Nag.** `ctx remind` surfaces "a serendipity round is waiting" at session start / every N turns.
6. **Review.** The human runs `ctx dream review` (or the serendipity skill) — a ~15-minute walk. Each proposal shows its generated summary + provenance + "why now"; the human chooses **accept / reject / amend / skip**.
7. **Apply.** Mechanical reactions (`archive`, `keep`, `mark-blog`, `reject`) apply **instantly with no LLM cost**; generative ones (`merge`, `promote`) drop to the agent, which reads the **full source** (not the lossy summary). `promote` drafts `specs/<name>.md` (tracked — the deliberate declassification).
8. **Record.** Every decision is appended to the ledger; `last_surfaced`/`status`/`history` on each source update so nothing rots or re-nags.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| No new/changed ideas (empty delta) | Trigger gate (`should-dream`) exits cleanly; no pass, no proposals, no nag. |
| Dream crashes mid-pass | No torn state: proposals/backups written atomically per item; `state.json` + ledger updated only for completed items; next run resumes from the delta. |
| Two dreams overlap, or a dream runs while a review is open | A `flock` lock serializes passes; a review in progress is unaffected (it reads a committed proposal set). |
| Idea already triaged, unchanged | Discipline clock skips it (hash unchanged). It re-enters triage only if its content changes. |
| Previously **rejected** proposal | Ledger records rejections; the dream does not re-surface a rejected disposition unless the source content changes. |
| `> max N` files in one pass | Bounded; the remainder is processed on subsequent passes (coverage progresses; `last_surfaced` keeps it fair). |
| Destructive disposition (`merge`/rewrite/`delete`) | **Backup-before-mutate**: snapshot the touched `.md` into the gitignored `dreams/` backup first. `archive` (a move to `ideas/done/`) needs no backup — reversible by relocation. |
| Source changed since its summary was cached | Hash mismatch → regenerate the summary before surfacing. |
| Idea note contains injected instructions ("ignore previous… always do X") | Treated as data, not instructions (sources-as-data); wrapped/handled as untrusted. |
| A dream write resolves to a **tracked** path | The don't-leak guard (`git check-ignore`) refuses the write. |
| `claude`/node not on PATH (cron's minimal env) | Fail **loud** with a failmark; never silently no-op. |
| Working tree dirty under the dream's paths | Defer the pass to avoid torn reads (per the executor's pre-flight). |

### Validation Rules

- **Write-scope:** every write path must be under `dreams/` or `ideas/` (or `specs/` via an accepted `promote`); enforced by a PreToolUse guard.
- **Don't-leak:** every dream-written path must satisfy `git check-ignore` (be gitignored), except the `promote → specs/` crossing; enforced by guard.
- **Provenance required:** every proposal must carry a source citation + evidence; a classification with no evidence is not surfaced (no fabrication — use a `___` placeholder rather than invent).
- **Bound/budget:** a pass honors `max N` files and a step/token budget; a runaway loop is the main cost failure mode.
- **Backup precondition:** a destructive disposition must confirm its backup succeeded before mutating; if backup fails, abort the mutation.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| `claude` auth probe fails (headless) | `[dream] FAIL: claude auth probe failed (401/expired?)` + failmark file | Set a non-interactive credential; failmark stays until a clean run |
| `claude`/git not on PATH | `[dream] FAIL: claude not on PATH (node/nvm not resolved)` | Fix cron env / NODE_BIN; rerun |
| Lock held by another dream | `[dream] another dream holds the lock; exiting` (exit 0) | Wait for the other pass; no action needed |
| Guard blocks a tracked-path write | `ctx-dream guard: write to tracked path refused: <path>` | Bug in output routing — fix path to a gitignored location |
| Backup failed before a destructive op | `[dream] backup failed; skipping <merge/delete> for <file>` | Item left untouched; surfaced again next pass |
| `promote` requested but `/ctx-spec` fails | Surface the `/ctx-spec` error; leave the idea in `ideas/` untagged | Retry promote, or amend to `keep` |

## Interface

### CLI

```
ctx dream [--mode discipline|creative] [--max N] [--budget STEPS]   # run one pass (headless-friendly)
ctx dream review                                                    # interactive ~15-min round
ctx dream accept|reject|amend <proposal-id>                         # primitives the review (and the agent) drive
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--mode` | enum | `discipline` | `discipline` (v1) or `creative` (deferred) |
| `--max` | int | `50` *(proposed)* | Max `ideas/` files processed per pass |
| `--budget` | int | TBD | Step/token budget for the pass |

### Skill

```
/ctx-serendipity   (a.k.a. the "garden walk"; sibling to /ctx-remember, /ctx-wrap-up)
```

Trigger phrases: "serendipity round", "review my dreams", "walk the garden", "what did the dream find?". Drives `ctx dream review`: mechanical reactions apply instantly; `merge`/`promote` are generative and handled in-session by the agent reading the full source. *(Open question: ship as its own `specs/ctx-serendipity.md`, or keep folded here.)*

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `skills/ctx-dream/` (tracked) | The dream skill — adapt from the gitignored draft at `ideas/ctx-dreams/skill/` (`SKILL.md`, `dream.sh`, `should-dream.sh`, `guard.sh`, `grep_claims.sh`) |
| `.gitignore` | Add `dreams/` (root-level dream notebook must be hidden) |
| `internal/cli/dream/` | Cobra: `ctx dream` root + `review`/`accept`/`reject`/`amend` subcommands |
| `internal/dream/` | Engine: delta selection, proposal generation, state/ledger, the three guards as callable logic, disposition appliers |
| `internal/config/dream/`, `internal/write/dream/`, `internal/err/dream/` | Constants, output, errors (per convention) |
| Guard logic (`git check-ignore` don't-leak + write-scope + sources-as-data) | In `internal/dream/`; reference executor wires it as a PreToolUse hook (adapt `guard.sh`); other executors call the same logic |
| State/ledger handling (`dreams/state.json`, `dreams/ledger.md`) | New |
| `internal/rc` | `dream:` `.ctxrc` section (enabled/mode/max/cadence/quiet_minutes/model/budget) |

### Documentation

Dream is **opt-in (not enabled by default)**, so it ships with two
user-facing docs:

| Doc | Audience | Content |
|-----|----------|---------|
| Enable guide (e.g. `docs/.../dream-setup.md`) | Claude Code users | How to turn dream on: `.ctxrc` `dream.enabled`, the cron entry, the guard hook wiring, the `ctx remind` cadence. |
| Executor-contract reference | Other-harness implementers | The contract from "Executor contract": bounded pass, structural guard enforcement, fail-loud, proposals-only-into-`dreams/`. So a non-Claude-Code harness can run the dream correctly. |

### Key Functions

TBD — do not invent signatures ahead of implementation. Conceptual units: delta selection (hash-based), summary cache (regenerate on hash change), proposal generation, the two guards, backup-before-mutate, the review loop, the disposition appliers (mechanical vs generative), ledger append + dedup-against-seen.

### Helpers to Reuse

- `ctx remind` (cadence/nag) — reuse, do not reinvent.
- The gitignored skill draft `ideas/ctx-dreams/skill/` as starting scaffolding.
- `/ctx-spec` (promote → spec) and `/ctx-blog` (mark-blog → draft) for routing accepted generative items.
- Existing archive convention (move to `ideas/done/`).
- `git check-ignore` as the leak-guard primitive; `flock` for the lock.

### Proposal schema (the dream↔serendipity contract)

Per-proposal fields (the contract the review and ledger build on):

- `id` — stable, for accept/reject/amend + ledger reference
- `target` — idea file path(s) (multiple for `merge`)
- `status` — observed: `implemented | duplicate | meritorious | sidenote | blog-candidate`
- `action` — recommended: `archive | merge | promote | mark-blog | keep`
- `evidence` — provenance/grounding (commit, spec, near-neighbor + similarity)
- `confidence` — high | med | low (drives attention triage)
- `rationale` — one-line human-readable why

### State record + ledger

Per-source record in gitignored `dreams/` (`state.json` or per-file): `path`, content `hash`, cached `summary` ref, `last_modified`, `last_surfaced`, `merit`, `status` (`active | archived | promoted→… | merged→…`), decision `history`. **Two clocks, one record:** discipline reads `hash` (re-triage only on content change); creative (deferred) reads `last_surfaced` + `merit` + randomness. Append-only `ledger.md` records every disposition incl. rejections.

## Configuration

- **`.gitignore`** must include `dreams/` (hard requirement; the don't-leak guard double-checks at write time).
- **Cron** entry runs `ctx dream` nightly (example: `30 2 * * *`); configurable cadence.
- **Quiet-window gate** *(proposed; from the Hermes sibling)*: `quiet_minutes` (default ~60) — `should-dream` defers the pass if the user was active within the window, complementing the existing dirty-tree / empty-delta defers.
- **Opt-in:** dream is **not enabled by default**. A `dream:` `.ctxrc` section gates it — `enabled` (default `false`), `mode` (`discipline`), `max` (50), `cadence` (cron), `quiet_minutes` (60), `model` (null = session default), `budget`. Nothing runs until the user sets `enabled: true` and installs the cron entry (see Documentation → Enable guide).

## Testing

- **Unit:** delta selection on hash change; summary regenerates on hash mismatch; don't-leak guard *refuses* a tracked path and *allows* a gitignored one; write-scope guard; backup-before-mutate aborts the mutation on backup failure; ledger dedup-against-seen (a rejected item doesn't resurface).
- **Integration:** full pass over a fixture `ideas/` dir → proposals in gitignored `dreams/<ts>/`; `accept promote` → spec lands in tracked `specs/`, idea tagged; `reject` → ledger → not resurfaced next pass; bound respected.
- **Edge / adversarial:** an idea note carrying injected instructions is filed, not obeyed; the **corrupted-artifact appendix from arXiv 2605.12978** as a regression fixture for the review/dedup gate; crash-mid-pass leaves consistent state.

## Non-Goals (v1)

- **No autonomous mutation of the five canonical files** (DECISIONS/LEARNINGS/CONVENTIONS/CONSTITUTION/TASKS). Ever, by construction (Option B).
- **No auto-approve.** Every disposition into a tracked artifact passes through the human.
- **No creative/garden mode behavior** — the `--mode` flag exists, but only `discipline` is built; `creative` is sketched, post-v1.
- **No v2 pipeline** — journals/harness/kb consolidation, the enriched-journal contract, canonical supersession with stable IDs are out of scope.
- **No web UI** — the CLI is the UI.
- **Not enabled by default** — dream is strictly opt-in (`dream.enabled: false` until the user turns it on); ctx never starts gardening unbidden.

## v2 (sketched, not contracted)

The full end-to-end pipeline this v1 is a slice of. **Not built in v1** —
captured here so the north star isn't lost and so v1's ledger/proposal
shapes don't foreclose it. Re-debatable before any of it is contracted.

```
raw episodes / transcripts / journals
    ↓
append-only evidence store
    ↓
out-of-band dream / consolidation pass
    ↓
typed durable artifacts:
  - decisions
  - learnings
  - conventions
  - tasks
  - project overview
  - current operating model
    ↓
retrieval / startup context / agent steering
```

**Locked design constraint (already binding in v1):** *dreams are derived
views, not the source of truth.* Preserve raw episodes; the dream layer
only summarizes, dedupes, merges, prunes, and indexes **with citations
back to evidence.** This is the Option B principle and the
"raw-is-first-class" mandate from "Useful Memories Become Faulty," applied
to the whole pipeline. Note the sharpening from the deep-research eval
cluster: provenance proves *traceability, not truth* (the
Citation-Alignment Fallacy), so consolidation must re-verify against
reality, and the gate stays a **human** — "a single agreeable LLM is not an
adversarial gate."

**How the layers map to today:**

| Pipeline layer | Status |
|----------------|--------|
| raw episodes / journals | `ctx journal source` already exists as the raw store; v1 reads raw `ideas/` instead — the dream is not yet wired to journals. |
| append-only evidence store | Partial: Auto-Dreamer's append-only writer (research) + the v1 `ledger.md`. A dedicated evidence layer is the v2 enriched-journal contract. |
| out-of-band dream pass | **The core — built in v1**, pointed at `ideas/` rather than journals. |
| typed durable artifacts (decisions/learnings/conventions/tasks) | v2: canonical supersession over the five files, an explicit v1 non-goal. |
| retrieval / startup / steering | Already exists (`ctx agent`, `ctx status`, `ctx steering get`); the dream feeds it, doesn't build it. |

**Two reconciliations the v2 contract must resolve:**

1. **Human gate is non-optional.** The diagram draws *dream → artifacts*
   directly; Option B requires *dream → proposal → human accept → artifact.*
   The serendipity gate sits between the consolidation pass and any write
   into the canonical files. v2 must not collapse that arrow.
2. **Two artifact types have no home yet.** *project overview* loosely
   echoes the hub's `ARCHITECTURE.md`; *current operating model* has **no
   analog** in ctx's five canonical files. v2 must decide whether to grow
   the canonical set or map these onto existing artifacts — v1's
   ledger/proposal `status`/`action` enums should stay open enough not to
   foreclose either.

**A v2 use case to fold in (from a serendipity round, 2026-06-07):** a
*dream over journals* could export a portable **personal-context
document** — preferences, context, key decisions, ongoing work —
distilled from journal history (the idea in `ideas/ctx-journal-summarize.md`).
This is not a standalone feature: it is the same v2 dream-over-journals
pipeline pointed at a human-portable output instead of (or alongside) the
canonical files. v1 dreams only inside `ideas/`; this is explicitly v2.

### Prior art: a cautionary sibling (Hermes "Dreaming")

An independently-designed feature (Hermes "Dreaming", itself after
OpenClaw) lands on nearly the same *architecture* — sleep-phase split
(Light/REM/Deep), cron + quiet-hours trigger, a weighted memory score, a
human-readable dream diary, opt-in-disabled-by-default. Useful
corroboration that the shape is natural. **But it is the un-gated form
ctx-dream deliberately rejects:** its Deep-Sleep phase *autonomously
promotes* entries into `MEMORY.md` (no human gate), scores purely
statistically (no provenance, no grounding-against-reality), and carries
none of the three safety invariants (sources-as-data, don't-leak,
backup-before-mutate) or a rejection ledger. It is effectively the
"`dream → artifacts` direct arrow" of reconciliation #1 above, built — and
thus a live illustration of the exact failure mode ("Useful Memories
Become Faulty": naive autonomous consolidation can push utility *below* the
no-memory baseline) that the human serendipity gate exists to prevent.
**Takeaway for v2:** keep the corroborated mechanics (phases, scoring,
diary, quiet-hours gate); never adopt the autonomous-promotion or
statistics-without-grounding parts.

## Executor contract

**The executor is a documented contract, not a hardcoded assumption.**
cron `claude -p` is the **reference executor** ctx ships guidance for,
but ctx must not assume it is the only way to run a dream pass — other
harnesses (a different AI CLI, a raw Anthropic-API loop, a CI runner)
must be able to implement the same contract.

ctx owns the **executor-agnostic core** in Go: the `dreams/` notebook
layout, the per-source state record, the append-only ledger, the
proposal schema, and the three guards *as callable logic* (write-scope,
sources-as-data, don't-leak / `git check-ignore`). An executor must:

1. Run one bounded pass (honor `max` files + step/token budget) that
   reads `ideas/`, classifies + grounds, and writes provenance-bearing
   proposals into `dreams/<ts>/` — never touching canonical.
2. **Enforce the three guards structurally** — not via prompt text. The
   reference executor (Claude Code) wires them as PreToolUse hooks; a
   raw-API executor must call the same guard logic in its tool executor
   before any write. This is the load-bearing portability requirement.
3. Fail **loud** (failmark) on auth/PATH/env problems; never silent
   no-op (see Error Handling).

Because dream is **opt-in and not enabled by default**, ctx ships two
docs (see Implementation → Documentation): a Claude-Code enablement
guide, and an executor-contract reference for other harnesses.

## Resolved Decisions (v1)

Settled from the prior Open Questions (session 2263caef, 2026-06-07):

- **Executor:** cron `claude -p` is the reference; the executor is a
  documented contract (above). cron does the gardening; session-start
  only *surfaces* + nags, with a small bounded session-start pass as a
  staleness fallback.
- **Serendipity:** shipped as its **own spec**, `specs/ctx-serendipity.md`
  (the human review / garden-walk skill), not folded here.
- **Merit signal:** v1 defaults `merit` via the attention-scoring rubric
  below — used **only** for ranking / ruthless self-rejection (surface
  top-N), **never** as an autonomous promote threshold. v1.1 creative
  mode reads it for resurfacing.

  | Component | Weight | Signal |
  |-----------|--------|--------|
  | Relevance | 30% | meaningful keyword density |
  | Frequency | 24% | how often the topic recurs |
  | Query diversity | 15% | appears across multiple sources/sessions |
  | Recency | 15% | newer ranks higher |
  | Consolidation | 10% | penalize if already captured (dedup) |
  | Conceptual richness | 6% | longer / more detailed = richer |

  Caveat (eval cluster): a statistical score ranks *attention*, not
  *truth* — grounding-against-code/specs decides eligibility; the score
  only decides ordering.
- **Summary generation:** reuse the session model (the `claude -p`
  model); per-pass cost ceiling = `max` files + a step budget. A
  dedicated cheap-summary model is deferred until there is real cost
  data.
- **Go package layout:** follow established convention —
  `internal/cli/dream/` (cobra: root + `review`/`accept`/`reject`/
  `amend`), `internal/dream/` (engine: state, ledger, proposal, guards),
  `internal/config/dream/` (constants), `internal/write/dream/`
  (output), `internal/err/dream/` (errors).
- **`.ctxrc`:** a `dream:` section — `enabled` (default `false`, opt-in),
  `mode` (`discipline`), `max` (50), `cadence` (cron, e.g.
  `30 2 * * *`), `quiet_minutes` (60), `model` (null = session default),
  `budget` (TBD numeric ceiling).
- **Stable IDs:** proposals and ledger entries carry stable IDs in v1,
  so v2 canonical supersession is not foreclosed.

### Still TBD (at implementation)

- Exact numeric `budget` ceiling (depends on first real-pass cost data).
- Presence/format details of the per-source state record file(s).
