# Dream Executor Contract

The ctx-dream **executor** is the thing that actually runs an out-of-band
dream pass: it reads `ideas/`, classifies and grounds each idea, and
writes proposals into the `dreams/` notebook. ctx ships **cron `claude
-p`** as the reference executor (see [Run the Dream](../recipes/run-the-dream.md)),
but the executor is a **documented contract, not a hardcoded
assumption** — any harness (a different AI CLI, a raw model-API loop, a
CI runner) can implement it.

This page is the contract. If you are wiring the dream into a non-Claude-
Code harness, implement everything below.

## What ctx owns (executor-agnostic)

The Go package `internal/dream` owns the parts that must behave
identically regardless of executor:

- The **data contract** — the proposal schema, the per-source state
  record (`dreams/state.json`), and the append-only ledger
  (`dreams/ledger.md`).
- **Delta selection** — the hash-based "discipline clock" that decides
  which ideas are new or changed since last triage.
- **The two structural guards** as callable logic — `WriteScope` and
  `Leak`.

Your executor must use these, not reimplement them.

## What an executor must do

1. **Run one bounded pass.** Honor the `max` ideas and step/token
   `budget` from the `dream:` `.ctxrc` section. Read only the idea delta.
2. **Propose, never act, never touch canonical.** The pass writes
   provenance-bearing proposals as a single JSON array to
   `dreams/<ts>/proposals.json` (the run directory is handed to the
   executor) and nothing else. It must not archive/merge/promote/tag
   ideas and must never write the five canonical files. (Acting on
   proposals is the human's `/ctx-serendipity` step, out of band from
   the pass.)
3. **Enforce the three guards structurally — not via prompt text.** This
   is the load-bearing portability requirement:
   - **Write-scope** — a write is allowed only under `dreams/` during a
     pass.
   - **Don't-leak** — every write target must be gitignored (`git
     check-ignore`); a write that resolves to a tracked path is refused.
   - **Sources-as-data** — idea text is wrapped as untrusted and is never
     executed as instructions.
   The Claude Code reference enforces write-scope and don't-leak with a
   **PreToolUse hook** (`guard.sh`) and sources-as-data via the skill's
   `<<<UNTRUSTED>>>` wrapping. A harness without hook interception must
   call the same checks in its own **tool executor** before every write
   — that is where `internal/dream.WriteScope` and `internal/dream.Leak`
   move. A prompt instruction is not enforcement.
4. **Fail loud.** On auth failure, a missing executor binary, or a
   PATH/env problem, write a failmark (`dreams/.failed`) and exit
   non-zero. Never silently no-op — a dream that quietly does nothing is
   indistinguishable from a healthy one that found nothing, and that
   ambiguity rots trust.
5. **Serialize passes.** Take the `dreams/.lock` before a pass; if it is
   held, exit cleanly. A review in progress reads a committed proposal
   set and is unaffected.
6. **Defer on a dirty tree.** If the working tree under the dream's paths
   is dirty, defer the pass to avoid torn reads.

## The proposal contract

Proposals are a JSON **array** in `dreams/<ts>/proposals.json`, each
element matching the `internal/dream.Proposal` schema:

```json
{
  "id": "<stable-id>",
  "targets": ["ideas/<file>.md"],
  "status": "implemented|duplicate|meritorious|sidenote|blog-candidate",
  "action": "archive|merge|promote|mark-blog|keep",
  "evidence": "<commit / spec path / near-neighbor + why>",
  "confidence": "high|med|low",
  "rationale": "<one-line why>"
}
```

`id` must be **stable** (so a re-run does not duplicate an already-decided
proposal, and so v2 canonical supersession is not foreclosed). An
executor must not re-emit a proposal whose `id` already appears in
`dreams/ledger.md` unless the source content changed.

## Why the contract, not just cron

The ctx dev team is multi-tool, and ctx's users are more so. Hardcoding
"the dream is cron + Claude Code" would exclude everyone else and couple
a memory feature to one harness. Keeping the cognition in a skill and the
invariants in `internal/dream` means the same dream — same guards, same
ledger, same proposals — runs anywhere the contract is met. See
`specs/ctx-dream.md` and the decision record in `.context/DECISIONS.md`
("ctx-dream executor is a documented contract").
