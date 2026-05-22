# ctx ai — extraction and recall (supplementary, sketched)

> **Status:** supplementary vision spec (blocks B + C).
> **Not a contract.** Sketched intent only; **re-debate after the
> block A spec ships** ([specs/ctx-ai-backend.md](./ctx-ai-backend.md)).
> A community developer reads this for context — to avoid
> foreclosing B/C with A's interface choices — but does **not**
> bind to it.
> **Brief:** [.context/briefs/20260522T001011Z-ctx-ai-backend.md](../.context/briefs/20260522T001011Z-ctx-ai-backend.md)

## Problem

Block A delivers the AI backend abstraction; it does not yet do
anything with it. The two real workloads ctx wants to run *over*
that backend are:

- **B. Structured extraction.** Turn meeting transcripts, design
  docs, code comments, PR discussions, and Slack exports into
  reviewable context artifacts — proposed patches against
  `.context/DECISIONS.md`, `.context/TASKS.md`,
  `.context/LEARNINGS.md`, `.context/GLOSSARY.md`,
  `.context/RISKS.md` (new file), and the kb evidence-index.
- **C. Embedding-backed recall.** Add a semantic-discovery
  dimension over project memory — find similar prior decisions,
  detect duplicate learnings, suggest related architecture docs,
  cluster notes into topics, find stale context that no longer
  matches code.

Both blocks share three failure modes (per the brief): silent
required-dep creep, deterministic/semantic path drift, and
ceremony substitution. Both must remain optional, additive, and
fail-closed. Both must land changes as proposed patches into a
review queue, never directly into authoritative state.

## Approach

### Block B — Structured extraction

vLLM (and most OpenAI-compatible servers) support
schema-constrained outputs — JSON schema, regex, grammar,
structural tags. ctx leverages this to turn unstructured inputs
into **typed proposals** for the canonical files.

The user's framing — lift verbatim from the brief:

> The CLI can still validate, diff and ask for (human or agent)
> acceptance.

Pattern shape (illustrative — final command/flag names TBD):

```text
ctx ingest meeting-notes.md \
  --extract decisions,tasks,risks,terms \
  --model local/qwen

ctx compact ~/.cursor/sessions/latest.jsonl \
  --emit decisions,learnings,tasks,open-questions

ctx evidence extract ./recordings/transcript.md \
  --topic cursor-hooks \
  --write-proposals .context/proposals/evidence.json
```

Schema example (decision extraction):

```json
{
  "type": "object",
  "properties": {
    "decision": { "type": "string" },
    "rationale": { "type": "string" },
    "owner": { "type": "string" },
    "evidence": { "type": "array", "items": { "type": "string" } },
    "confidence": { "type": "number" }
  },
  "required": ["decision", "rationale", "evidence"]
}
```

Hard CLI enforcement (verbatim from the brief):

- no citation → no evidence row
- no source span → reject
- duplicate semantic claim → merge proposal
- unsupported claim → quarantine

### Block C — Embedding-backed recall

vLLM serves embedding models through an OpenAI-compatible
Embeddings API. ctx uses that to build a semantic index over the
project's `.context/`, `docs/`, and (optionally) `src/` trees.

Pattern shape (illustrative — final command/flag names TBD):

```text
ctx index .context ./docs ./src
ctx search "why did we reject strict-CWD?"
ctx related DECISIONS.md:D-017
ctx doctor --semantic
ctx drift --semantic
ctx assemble "review PR 184 for Velero scalability risk" \
  --rank-with vllm \
  --budget 24000
```

Storage: **SQLite by default** ("a filesystem wearing a DB hat"
— user's framing; preserves Invariant 1's spirit because the
index is a derived view, not authoritative state). Other vector
stores (LanceDB / Qdrant / pgvector / FAISS / HNSW) are opt-in
escape hatches the user wires themselves.

Semantic `ctx assemble` classifies each chunk as:

- `must_include`
- `useful_if_budget`
- `background_only`
- `irrelevant`
- `dangerous_stale`

(Verbatim from the brief.) This keeps the system auditable.

### Cross-cutting design discipline

Both B and C must respect the same invariants A enforces:

- **Additive only.** No existing command degrades or changes
  shape when the AI backend is unconfigured.
- **Sibling commands for semantic paths.** `ctx agent`
  (deterministic, Tier 2) stays untouched; `ctx assemble
  --rank-with <backend>` is a *new sibling*, not a flag on
  `ctx agent`.
- **Proposed patches, never auto-merges.** Every AI-produced
  artifact lands in a review queue. The five canonical files
  remain human-ratified (Invariant 4).
- **Prefix-cache-aware prompt structure.** Both blocks should
  structure prompts as `[stable project context prefix][task
  context][user request]` to let vLLM's prefix caching amortise
  the expensive prefix work. This is performance shape, not a
  contract.

## Behavior

### Happy Path (sketch)

**B example — `ctx compact`:**

1. User runs `ctx compact ~/.cursor/sessions/latest.jsonl --emit
   decisions,learnings,tasks,open-questions`.
2. ctx loads the input, structures a prompt with stable project
   prefix + transcript body + JSON-schema-constrained
   instructions, dispatches through the configured backend.
3. Backend returns schema-constrained JSON. ctx validates,
   deduplicates against existing entries, and writes a proposal
   artifact (e.g., `.context/proposals/20260522T-compact.md`).
4. User runs `/ctx-wrap-up` (or a future `/ctx-proposals review`
   skill) which surfaces the proposal, lets the user accept /
   edit / reject each row, and only then writes to
   `.context/DECISIONS.md` etc.

**C example — `ctx search`:**

1. User has previously run `ctx index .context ./docs`.
2. User runs `ctx search "why did we reject strict-CWD?"`.
3. ctx embeds the query via the backend, queries the SQLite
   vector index, returns top-N matching chunks with source
   citations.
4. No state mutation.

### Edge Cases (sketch — to be enumerated in the contract version)

- AI backend unreachable mid-extraction: discard partial output,
  fail closed, no partial proposal artifact.
- Schema-constrained response is well-formed but semantically
  empty (e.g., zero decisions extracted from a transcript that
  clearly contains decisions): surface to the user; do not
  silently write an empty proposal.
- Embedding model dimensions change between `ctx index` runs:
  refuse subsequent index/search ops with a clear "model drift
  detected; reindex required" error.
- Proposal queue conflicts (two extraction runs propose the same
  decision row): merge per the brief's "duplicate semantic claim
  → merge proposal" rule; the contract version must specify the
  merge algorithm.
- SQLite vector index grows large: TBD (vacuum strategy, sharding
  by source-tree subdirectory, opt-in pruning).
- User asks `ctx search` over a path that was indexed yesterday
  but has since changed: surface staleness ("index is N hours
  old; run `ctx index <path>` to refresh") rather than silently
  returning possibly-stale matches.

### Validation Rules (sketch)

- B: every proposed row must carry a source citation (file:line
  or transcript timestamp). No citation → no row.
- B: extracted entries diff-checked against existing canonical
  files before proposal; duplicates surface as merge proposals.
- C: embedding model identity (name + dimensions + tokenizer)
  recorded in the SQLite index; mismatches refuse to query.
- C: indexed paths recorded with last-indexed-at timestamp;
  search results surface staleness.

### Error Handling (sketch)

Inherits A's error sentinel pattern. Specific error conditions
to be enumerated in the contract version. Notable categories:

- Schema-constrained output rejected by upstream (invalid JSON,
  schema violation) → surface verbatim; offer to retry with a
  more capable model.
- Embedding index corruption → refuse to query; suggest
  reindex.
- Proposal queue write failure → fail closed; no partial state.

## Interface (sketch)

CLI shape is illustrative until A's namespace decision lands.

If A picks **Option 1 (`ctx ai <verb>`):** B and C verbs nest
under `ctx ai ingest`, `ctx ai compact`, `ctx ai search`,
`ctx ai index`, etc.

If A picks **Option 2 (flags on existing commands):** B and C
extend existing commands (`ctx compact --emit ...`, `ctx assemble
--rank-with ...`) and add new top-level verbs only where no
existing command fits (`ctx index`, `ctx search`).

The B+C contract spec inherits this choice from A; this
supplementary does **not** pre-commit.

## Implementation (sketch)

Likely package shape (subject to A's interface decision):

- `internal/cli/ingest/` — block B command surface.
- `internal/cli/index/`, `internal/cli/search/` — block C
  command surface.
- `internal/extract/` — schema-driven extraction pipeline;
  dedup, merge, quarantine logic.
- `internal/embedding/` — embedding dispatch + index abstraction.
- `internal/embedding/sqlite/` — SQLite vector store
  implementation (default).
- `internal/embedding/opaque/` — adapter shape for opt-in
  external stores (LanceDB, Qdrant, pgvector, FAISS, HNSW);
  out-of-band wiring, never a hard dep.
- `internal/proposals/` — proposal queue read/write; ratification
  hooks for ceremony skills.
- `internal/assets/claude/skills/ctx-proposals/` — new skill
  family for surfacing and ratifying proposals (analogous to
  `/ctx-kb-*` for kb closeouts).

## Configuration (sketch)

Extends A's `.ctxrc` `[backends]` table with B+C-specific
sections:

```toml
[ai.extract]
default_model = "openai/gpt-oss-120b"
proposal_dir = ".context/proposals"
require_citations = true   # enforces "no citation → no row"

[ai.embeddings]
default_model = "text-embedding-3-small"   # or local equivalent
index_path = ".context/state/embeddings.db"
indexed_paths = [".context", "docs", "src"]
```

## Testing (sketch)

- B: schema-constrained extraction unit tests against a fake
  backend; dedup/merge/quarantine logic; proposal artifact
  shape; ratification round-trip via the proposals skill.
- C: embedding dispatch unit tests; SQLite vector store
  correctness (insert, query, dimension-mismatch refusal);
  staleness detection; `ctx assemble --rank-with` produces
  the documented chunk classes.
- Cross-cutting: confirm `ctx agent`, `ctx status`, and all
  ceremony commands are bit-for-bit identical with and without
  B/C installed; AI commands fail closed when backend
  unreachable.

## Non-Goals

Inherits A's Non-Goals (cost, secrets, observability, proxy,
auto-application, ceremony replacement). Plus B+C-specific:

- **Online learning / fine-tuning.** The system uses
  off-the-shelf models; no in-loop model updates.
- **Vector DB as a service-shape dep.** SQLite is the default;
  any other vector store is opt-in and user-wired.
- **Semantic rewriting of canonical files.** AI may *propose*
  changes; it never *rewrites*.
- **Auto-acceptance of any AI-produced artifact.** Even
  high-confidence proposals require human or agent ratification
  through the existing ceremony surface.
- **Replacing deterministic `ctx agent` assembly.** Semantic
  assembly is a sibling command, never the default.

## Open Questions

(All of these are intentionally open. The point of this
supplementary spec is to make the shape *visible* so A's
implementation doesn't accidentally foreclose any of them.)

1. **Proposal queue location and format.** `.context/proposals/`
   with one file per extraction run? A unified
   `proposals/inbox.md` per the kb closeout pattern? Inherits
   from B's contract version.
2. **Ratification UX.** A new `/ctx-proposals` skill family, or
   absorbed into `/ctx-wrap-up` and `/ctx-kb-ingest`?
3. **Embedding model recommendations.** Recipe vs. baked into
   the spec? Probably recipe — too fast-moving for spec
   commitment.
4. **Reindex triggers.** Manual only? File-watcher-driven?
   Hook-driven on commit?
5. **Search-result citation format.** Always include source line
   ranges? Configurable?
6. **`ctx drift --semantic` vs. existing `ctx drift`.** Sibling
   subcommand, flag on existing, or a separate `ctx doctor
   --semantic`?
7. **Cross-language inputs.** B accepts transcripts in any
   language ctx already supports (per `multilingual-sessions`
   recipe); does C's embedding model need per-language
   selection?
8. **Privacy boundary.** When the configured backend is *not*
   local (e.g., the user pointed at OpenAI), should B/C refuse
   to send certain file paths (e.g., anything matching a
   `.ctxrc` `[ai].redact` glob)? A "do not exfiltrate"
   discipline likely belongs here, not in A.
