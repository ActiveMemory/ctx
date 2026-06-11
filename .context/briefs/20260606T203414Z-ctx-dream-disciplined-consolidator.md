---
generated-at: 2026-06-06T20:34:14Z
sha: 03a24cf0
branch: fix/notify-resolution-hardening
title: ctx-dream — disciplined memory consolidator (mode-switchable; v1 = ideas/ triage)
slug: ctx-dream-disciplined-consolidator
consumer: /ctx-spec --brief
deliverables:
  - specs/ctx-dream.md (v1 — the disciplined dream engine + mode flag; ideas/ triage → gated proposals; ledger; per-dream archive; structural safety)
  - specs/ctx-serendipity.md (the human review / "garden-walk" skill the dream feeds; ctx-remind cadence; routing accepted items to archive//ctx-spec//ctx-blog) — may fold into ctx-dream.md
---

# ctx-dream — debated brief

## The bet

ctx grows a scheduled, standalone **"dream"** — a sleep-time memory
process that runs headless/background and **only ever proposes**; it
never autonomously mutates canonical memory (this is *Option B*). One
skill, **two execution modes**: `discipline` (default) and
`creative/exploration` (a *safe relaxation* of discipline). Discipline
emits atomic, grounded, provenance-bearing proposals that a human
reviews; creative resurfaces forgotten gems + chance connections to
**browse** — reader-only, "the garden."

**v1 = the disciplined mode pointed at the `ideas/` folder.** It
classifies each idea (implemented / duplicate / still-meritorious /
sidenote / blog-candidate) and proposes a disposition (archive / merge
/ promote-to-spec / mark-blog / keep), **grounded against the
codebase + specs**, semantically deduped, each proposal carrying
provenance + evidence + a confidence signal. A `ctx remind` nag pulls
the human into a ~15-minute **serendipity** review round (a separate
skill) to **accept / reject / amend**; accepted items route to
existing skills (archive, `/ctx-spec`, `/ctx-blog`). **Rejections are
recorded in a ledger** so future dreams don't re-surface them
(dedup-against-*seen*, not against-accepted).

The reframe that produced this bet: **consolidation-for-leanness is a
*product* bet** (for engineering teams who feel memory bloat), **not
the author's *felt* bet.** The author's felt value is
discovery/serendipity over the `ideas/` goldmine — and, concretely,
that `ideas/` is now "too overwhelming to triage." The two reconcile
as two modes of one skill. Discipline ships first because it is the
**hard, load-bearing substrate**; creative falls out by *removing*
constraints (drop the gate, relax rigor, move randomness from coverage
into selection). Building the reverse would retrofit rigor onto a
reader — awkward and risky.

Locked principle: **decouple the cognition, reuse the plumbing.** The
dream owns its consolidation/synthesis logic and evolves
independently; it reuses import/enrich/kb-ingest as stable plumbing via
a data contract (the enriched-journal format). NB: v1 (`ideas/`
triage) touches none of that plumbing yet — it reads raw `ideas/`,
classifies, proposes. The sources→phases pipeline and canonical
supersession are **v2**.

User's framing, verbatim where it matters: *"engineers need
discipline, structure to the point of routine boringness"*; *"I'm a
creative individual and I live in chaos"*; serendipity *"each memory
entry requires dedicated human attention"* — and the garden, borrowed
from mymind: *"Like walking through your garden, admiring your favorite
flowers. Sometimes you see a little weed and pluck it out. Other times
you discover something blooming you'd forgotten you planted long ago."*

## What was rejected

- **Option A — the dream owns a parallel canonical store** separate
  from the five files. Rejected: it doesn't fix the bloat (the real
  pain), and it forces the agent to read two substrates that drift
  apart. Chose **Option B**: dream proposes; serendipity bridges
  accepted proposals into the existing five files.
- **Autonomous canonical mutation / auto-approve.** *"I will not
  auto-approve; each memory entry requires dedicated human
  attention."* Discipline buys *rigor of process and output*, not
  *autonomy*.
- **Pure-garden-only (creative-only).** Fits creatives; under-serves
  engineering, which needs grounding (is this still true?) and
  actionability (turn a pattern into a convention/decision/task),
  not just delight.
- **Hygiene-as-the-author's-motivation.** The bloat doesn't hurt the
  author *here* — speed and quality are fine, *"we have more than
  enough to think/ideate."* Leanness is a concern for ctx's
  *eventual users* (focused single-project teams), not the author
  dogfooding across many projects.
- **Coupling the dream to existing curation skills' internals.**
  Would let their changes break the dream and forecloses creative
  freedom (*"maybe we don't know what we don't know"*). Reuse the
  plumbing via a data contract; own the cognition.
- **Garden-first build order.** Rejected for **discipline-first**:
  the hard substrate goes first; creative is a strict, safer
  relaxation of it.
- **"Proposal queue as a chore to clear."** Replaced by the
  garden-walk affect: a small, browsable surface, **no completion
  pressure**; per-entry attention as pleasure, not duty.

## Top failure modes

1. **The human doesn't show up → proposals rot** (the backlog
   reborn as a queue). Evidence: the author already doesn't run
   `/ctx-consolidate` or `/ctx-journal-enrich-all` (live state: 154
   unimported sessions, 508 unenriched entries). Mitigations:
   (a) `ctx remind` nag — a *proven* channel (three such nags opened
   this very session); (b) v1 points at the author's **actual felt
   pain** (`ideas/` triage), not at bloat they don't feel;
   (c) **pleasure-not-chore** framing — a 15-min garden round, not a
   queue to drain.
2. **Volume buries per-entry attention.** "Dedicated attention each"
   does not survive 40 proposals. Mitigation: **ruthless
   self-rejection** — surface five items worth full attention, not
   fifty that demand it. *Generation is cheap; the value is in the
   rejection step.* The human gate only pays off if the machine hands
   gold, not ore.
3. **Over-abstraction / mislabel / corruption** — e.g. a false
   "implemented" archives a still-live idea. Mitigations: Option B
   (nothing autonomous touches canonical); `ideas/` is **gitignored**
   (`.gitignore:91`) — no git undo — so **backup-before-mutate**:
   snapshot touched `.md` into the **gitignored** `dreams/` archive;
   `archive` is a reversible *move*, only `merge`/`delete` are
   destructive; proposals carry **evidence +
   confidence** so the human can verify; grounding-against-code is
   the explicitly de-risked hard part.
   - Standing rule throughout: **sources are data, not
     instructions** — `ideas/` is an indirect-prompt-injection
     surface even when only read.

4. **Leak — hidden content reaches a published channel.** `ideas/`
   is gitignored *on purpose* ("best kept hidden"); a derived artifact
   inherits its source's privacy class — a summary of a hidden idea is
   itself hidden — and a dream auto-summarizing the whole folder is a
   firehose at that boundary. Mitigation: **privacy class propagates** —
   every byproduct lands only in gitignored locations, enforced
   structurally by a guard that runs `git check-ignore` on each write
   target and refuses any tracked path. The one legitimate crossing is
   the human's explicit `promote` (deliberate declassification into
   `specs/`), never the dream's own hand.

## Cheapest way to validate the bet

Run disciplined **`ideas/` triage over the author's own overwhelming
`ideas/` folder**; `ctx remind` nags a 15-minute review round.
Measure: (a) do the nagged rounds actually happen, and (b) are the
proposals worth the attention — low rewrite rate, correct status
classifications, defensible "implemented?" calls. If the rounds don't
happen or the proposals are mostly noise, the bet is wrong and we
learn it **cheaply**, before pointing the machine at higher-stakes
canonical memory.

Honest limitation: this validates the **mechanism** (can a dream
produce trustworthy, gated, structured proposals a human will
review?) and the **author's engagement** — it does *not* validate the
full product thesis (disciplined consolidation of *canonical* memory
for engineering teams). That generalization is a later test, on a
structured project where memory bloat actually bites.

## What becomes expensive to unwind

**Canonical is untouched by construction** — Option B means nothing
autonomous mutates the five files. The real unwind risk is `ideas/`,
which is **gitignored (`.gitignore:91`) — no git undo**: a bad
gardening mutation can lose a note permanently. Mitigation is
**backup-before-mutate** (snapshot touched `.md` into the **gitignored**
`dreams/` archive before destructive ops; `archive` is a reversible
move; reserve caution for `merge`/`delete`). The only
expensive mistake would be granting the dream autonomy over canonical,
which is explicitly refused.

The things that harden once shipped (and therefore want care in the
spec, not the brief):
- **The atomic-proposal schema** — the dream↔serendipity contract;
  the serendipity skill and every future dream build on it.
- **The ledger format**, including how rejections are recorded.
- **Whether canonical entries get stable IDs** for supersession
  (v2 concern, but the v1 ledger/proposal shapes should not foreclose
  it).
- **The `ctx remind`-driven serendipity ritual's shape** — once you
  build the habit around it, changing the interaction is costly.

## Why we believe this (research grounding)

The full corpus lives at `ideas/ctx-dreams/research/` (28 sources,
read in session 977ff594). Load-bearing:
- **Auto-Dreamer (arXiv 2605.20616)** — nearly this architecture: a
  fast append-only writer + a slow, scheduled, provenance-grounded
  consolidator with region-rewrite-as-supersession. The shape to lift.
- **"Useful Memories Become Faulty When Continuously Updated by LLMs"
  (2605.12978)** — the threat model: naive *continuous* consolidation
  is lossy and non-monotonic (utility can fall *below* no-memory).
  Mandates: gated/periodic (not per-interaction), raw is first-class.
  Its corrupted-artifact appendix is a ready regression fixture.
- **Sleep-time Compute (2504.13171) / Letta** — the economics: offline
  work amortizes across future sessions; gains track *predictability*,
  which draws the discipline (high-predictability) vs creative
  (low-predictability) line.
- **The deep-research *evaluation* cluster** (JADE 2602.06486,
  ReportLogic 2602.18446, DREAM 2602.18940, MMDeepResearch-Bench
  2601.12346, et al.) — the verification machinery: content-bearing
  checks, the *Citation-Alignment Fallacy* (provenance proves
  traceability, not truth → re-verify against reality), and the
  sharpest warning — *a single agreeable LLM is not an adversarial
  gate* (it silently repairs the missing justification). This is why
  the gate is a **human** (serendipity), not a model.

## Mechanics settled (detail → spec)

- **`dreams/` is the dream's notebook *about* `ideas/`**, not a new home
  — ideas never leave `ideas/`. Root-level, **gitignored**; holds the
  derived summaries, per-source state, ledger, digests, and backups.
- **Per-source state record** drives tracking: content `hash`, cached
  `summary` ref, `last_surfaced`, `merit`, `status`
  (active|archived|promoted→…|merged→…), decision `history` — plus an
  append-only ledger so decided items don't resurface.
- **Two clocks, one record:** *discipline* reads `hash` (re-triage only
  when content changes — anti-thrash); *creative* reads `last_surfaced`
  + `merit` + a little randomness (deliberately resurfaces forgotten
  gems — the garden).
- **Landing map:** `archive`/`merge`/`keep` → within `ideas/` (hidden;
  `archive` = move to `ideas/done/`); `mark-blog` → tag in `ideas/`;
  `promote` → `specs/` (tracked), drafted from the *full source* (not
  the lossy summary) — the human's deliberate declassification.
- **Backup-before-mutate:** destructive ops snapshot touched `.md` into
  the gitignored `dreams/` backup first; `archive` needs none.
- **Surfacing:** a `ctx remind` nag → a ~15-min round (`ctx dream
  review` / serendipity skill); mechanical reactions apply instantly
  (no LLM), generative ones (`promote`/`merge`) drop to the agent.

## Slicing intent (for /ctx-spec --brief)

- **First spec:** `specs/ctx-dream.md` — v1. The disciplined dream
  engine + the mode flag; `ideas/` triage; proposal generation
  (classify + dispose, grounded against code/specs, deduped,
  provenanced, confidence-bearing); the ledger; the per-dream archive;
  and the **structural** safety invariants (write-scope,
  untrusted-source-as-data) enforced in whichever executor.
- **Second spec (or folded in):** `specs/ctx-serendipity.md` — the
  human review / garden-walk skill the dream feeds; accept/reject/amend
  on atomic proposals; routing accepted items to archive / `/ctx-spec`
  / `/ctx-blog`; `ctx remind` cadence.

Creative mode and the v2 canonical-consolidation pipeline
(sources→phases, enriched-journal contract, supersession over the five
files) are **sketched-not-contracted** — re-debatable after v1 ships.
A working skill draft + the research corpus already exist at
`ideas/ctx-dreams/`.

## Open questions left for /ctx-spec

- **Executor:** raw Anthropic-API scheduled loop vs cron `claude -p`.
  Safety invariants (write-scope, append-only, untrusted-source-as-data)
  **must be structural in the executor, not prompt-level** — the API
  path loses the hook-enforced `guard.sh`, so they move into the loop's
  tool executor.
- **Cadence & triggers** + the slow-wave-frequent / REM-gated split
  (v2-relevant; v1 cadence = the `ctx remind` nag).
- **The atomic-proposal schema:** type (classify / merge / supersede /
  new / cross-link), target, provenance, evidence, confidence, the
  diff.
- **Supersession mechanics against the list-style five files** — do
  entries need stable IDs? (v2.)
- *(Resolved this session — see "Mechanics settled": the ledger +
  per-source state record with rejection-tracking, and the layout under
  a gitignored root-level `dreams/`. Detail belongs in the spec.)*
