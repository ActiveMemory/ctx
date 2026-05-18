# KB Editorial Pipeline (Paired with Handover Mechanism)

> Source: `ideas/003-editorial-pipeline-debated-brief.md` (debated
> brief from `/ctx-plan` session 01d0cf92, 2026-05-09).
>
> Revision 2 (2026-05-16): rewritten to absorb the **current**
> upstream editorial-pipeline shape (pass-mode contract, completion circuit breaker,
> source-coverage ledger as state machine, topic-adjacency pre-flight,
> cold-reader orientation rubric, folder-shaped topics from day one,
> failure analysis). The original 4-phase model is superseded; the
> brief's two organizing principles (LLM as migration tool;
> KB-of-KBs is a KB) carry forward. See DECISIONS.md
> "Lift current upstream editorial-pipeline shape, not the 4-phase predecessor"
> for the migration rationale; this comparison note
> `ideas/upstream-pipeline-comparison.md` is the input.
>
> Authority order when this spec, `DECISIONS.md`, or `docs/`
> could conflict: frozen contracts in `docs/` > recorded decisions
> in `DECISIONS.md` > the brief at `ideas/003-...` > inference
> labeled `TBD` here. Do not invert.

## Problem

`ctx` users doing knowledge-shaped work (research projects,
domain modeling, post-incident reviews, vendor-spec analysis)
currently have no native surface for it. The canonical files
(`TASKS.md`, `DECISIONS.md`, `LEARNINGS.md`, `CONVENTIONS.md`)
are tuned for code-development context, not for evidence-tracked
knowledge with confidence bands, contradictions, and external
grounding.

The cost is concrete and active:
`your-project` has been hand-rolling an
editorial pipeline at the repo root (`10-CONSTITUTION.md`,
numbered mode prompts, `20-INBOX.md`, hand-typed 8-item
closeouts) and disabling half of `ctx`'s code-dev skill surface
in `CLAUDE.md` to avoid name collision with its editorial
constitution.

Separately, `ctx` has no per-session **handover artifact**: the
narrative thread between sessions is currently re-derived
probabilistically from canonical files plus journal. An ad-hoc
one-off (`HANDOVER-2026-04-22.md`) exists but no discipline
shapes it.

This spec adds both surfaces and integrates them via a
closeout/fold mechanism so editorial work flows cleanly into
session-to-session continuity.

## Approach

Lift the **current** upstream editorial-pipeline shape into `ctx` with a
non-colliding rename (`KB-RULES.md`, not `CONSTITUTION.md`).
The shape:

1. **Pass-mode contract**: every `ctx kb ingest` invocation
   declares its mode (`topic-page` / `triage` / `evidence-only`)
   before extraction, with a definition-of-done. Mid-pass
   mode-switching is forbidden; abort and re-invoke instead.
2. **Completion circuit breaker**: `topic-page` mode passes
   four invariants before claiming `produced` / `extended`:
   file exists; cites ≥ 1 `EV-###`; site build clean (or named
   failure); cold-reader rubric passes.
3. **Source-coverage ledger**: state machine over every source
   the kb has touched. `discovered → admitted →
   highlights-extracted → partially-ingested → topic-page-drafted
   → comprehensive`, plus terminals `superseded` / `skipped`.
   Every pass updates the ledger honestly; "lying to the ledger"
   is a hard anti-pattern.
4. **Topic-adjacency pre-flight**: before resolving a topic,
   scan the ledger for sibling-topic rows not in
   `{comprehensive, skipped, superseded}` and force the page
   being authored to acknowledge them. Makes discoverability
   failures mechanically hard.
5. **Cold-reader orientation rubric**: four yes/no questions
   the artifact must satisfy: concept clear? why this project
   cares clear? canonical evidence reachable? boundaries clear?
6. **Folder-shaped topics from day one**:
   `.context/kb/topics/<slug>/index.md` for every topic, with
   optional sibling sub-pages (`security.md`, `multi-surface.md`).
   Lazy split: only when `index.md` outgrows the cold-reader
   rubric's "boundaries clear?" check.
7. **CLI-as-scaffold-authority**: topic-page file creation is
   performed only by `ctx kb topic new`; skills invoke the CLI,
   never synthesize the scaffold.

Pair the editorial pipeline with a per-session **handover
artifact**. The closeout/fold mechanism is the integration
point: every editorial pass writes a closeout with
`generated-at` frontmatter; the next `/ctx-handover` folds
postdated closeouts and archives the sources.

**Two organizing principles** (carried forward verbatim from
the brief because they explain why this is rationally bold
rather than recklessly broad):

- **P1: The LLM is the migration tool.** Wholesale ID
  renumbering, taxonomy reshuffles, confidence-band remapping,
  cross-file reference rewrites are absorbed by an LLM cleanup
  pass. We commit to specific schemas in v1 instead of hedging
  with abstract types. Be wrong cheaply.
- **P2: A KB is knowledge; a KB of KBs is a KB.** Recursive
  composability. Federation is `source-map kind: kb` plus
  pointing the standard ingest pipeline at another KB. No
  special feature needed; v1 schemas accommodate v2 federation
  without lock-in.

**KB ontology** (decided in the brief, load-bearing for the
skill surface): in a KB, you do not *decide*; you *learn* and
your *confidence increases*. A claim with confidence >0.9 is a
fact-by-contract. There is no `/ctx-kb-decide` skill; the
pipeline is the sole writer; ad-hoc capture flows through
`ctx kb note` (lightweight) or hand-edit (escape hatch).

**Mixed-mode separation:** `ctx kb` writes only to
`.context/kb/` and `.context/ingest/`. Canonical capture skills
write only to canonical files. Strict command authority
prevents cross-pollution; the read side (`/ctx-remember`,
`ctx status`, `ctx agent`, `/ctx-wrap-up`, session-start
hooks) becomes mode-aware via "if `.context/kb/` exists, also
fold KB state into the readback."

## Behavior

### Happy Path

**End-to-end: research session through cold restart.**

1. User runs `git init && ctx init` in a fresh project. Init
   lays down the five canonical files **plus**
   `.context/handovers/`, `.context/kb/` (with `.gitkeep` and
   `topics/.gitkeep`), `.context/ingest/` (with `KB-RULES.md`,
   mode prompts, `INBOX.md`, `SESSION_LOG.md`,
   `grounding-sources.md`, `OPERATOR.md`, `PROMPT.md`,
   `schemas/`), and `.context/site/` (gitignored).
2. User runs `ctx setup` to deploy skills. New skills:
   `/ctx-handover`, `/ctx-kb-ingest`, `/ctx-kb-ask`,
   `/ctx-kb-site-review`, `/ctx-kb-ground`, `/ctx-kb-note`.
3. User invokes
   `/ctx-kb-ingest ./inputs/2026-04-12-call.md "cursor hooks"`.
   The skill:
   1. Verifies pre-write gates: `.context/`,
      `.context/kb/` exist, kb scope is declared
      (`.context/kb/index.md` has a non-placeholder
      `## Scope`).
   2. **Declares pass-mode**: emits up-front block
      (`Pass-mode: topic-page`, `Reason: ...`,
      `Definition of done: ...`).
   3. Resolves topic, scans
      `.context/kb/source-coverage.md` for adjacency
      pre-flight, lists sibling-topic gaps.
   4. Resolves sources, advances ledger to `admitted`.
   5. Calls `ctx kb topic new "cursor hooks"` if the topic
      folder does not exist (scaffold authority).
   6. Synthesises prose section by section, mints `EV-###`
      rows in `evidence-index.md`, updates `glossary.md`,
      `timeline.md`, `source-map.md` as needed.
   7. Applies life-stage reconciliation discipline (skipped
      if `< 5` topic pages; full discipline at `>= 5`).
   8. Sets Confidence floor per the cited-bands rule.
   9. Updates `source-coverage.md` to `topic-page-drafted`
      or `comprehensive`.
   10. Runs circuit-breaker check: file exists; cites EV;
       `ctx kb site build` clean; cold-reader rubric records
       `Result: pass`.
   11. Writes closeout to
       `.context/ingest/closeouts/<TS>-ingest-closeout.md`
       with required frontmatter
       (`sha`, `branch`, `mode`, `pass-mode`, `life-stage`,
       `generated-at`).
4. User runs `/ctx-wrap-up`. The skill (a) walks the standard
   capture checklist, (b) **always** delegates to
   `/ctx-handover` as its final step regardless of whether
   `.context/kb/` exists. When `.context/kb/` does exist, the
   wrap-up additionally surfaces editorial state (pending
   closeouts, unresolved outstanding-questions count) before
   delegating; KB presence affects what gets folded, not
   whether the handover is written.
5. `/ctx-handover` (wrap-up's handover step) collects
   `--summary` (past tense, what happened) and `--next`
   (future tense, specific first action), runs
   `ctx handover write`, which folds postdated closeouts into
   `## Folded closeouts`, archives the source closeouts to
   `.context/archive/closeouts/`, and writes
   `.context/handovers/<TS>-<slug>.md`. The filename is
   timestamped so concurrent agent runs never overwrite one
   another.
6. User opens a fresh session next morning. The session-start
   hook reads the latest handover. User asks "do you
   remember?". `/ctx-remember` reads canonical files +
   latest handover + any closeouts whose `generated-at`
   postdates the handover (none, in the happy case) and
   presents structured readback citing Summary + Next session.
7. User runs `ctx kb site serve` to browse the rendered KB at
   `localhost:8000`. Render output lands in `.context/site/kb/`
   via a shell-out to `zensical` (already used for
   `ctx journal site`).

### Pass-Mode Contract

Every `ctx kb ingest` invocation classifies itself as exactly
one of three modes **before any source extraction begins**:

| Mode             | Mints prose? | Mints EV-### ? | Touches topic page? | Default? |
|------------------|--------------|----------------|---------------------|----------|
| `topic-page`     | yes          | yes            | yes (create/extend) | yes      |
| `triage`         | no           | **no**         | no                  | no       |
| `evidence-only`  | no           | yes (tagged)   | no                  | no       |

**Mode selection rules.** Default is `topic-page`. `triage`
fires when the user supplied multiple disparate sources with
no clear single topic, OR explicitly invoked triage language.
`evidence-only` fires only on explicit user request matching
the valid-trigger criteria (*"Just mint EV rows"*,
*"backfill evidence"*); the skill MAY NOT infer it from source
size, ambiguity, time pressure, or operator convenience.

**Up-front declaration (mandatory).** Before extraction begins
(after pre-write gates pass and **before** topic resolution),
the skill MUST surface a visible pre-work declaration:

> **Pass-mode:** `<mode>`
> **Reason:** `<one sentence; required when non-default>`
> **Definition of done:** `<mode-specific completion criterion>`

The declaration is a contract, not a label.

**Mid-pass mode-switching is forbidden.** If the work in flight
no longer fits, abort with a partial closeout citing what was
done, and recommend re-invocation under the correct mode.

### Topic-Page Circuit Breaker

A pass in `topic-page` mode MAY NOT report
`topic-page: produced` or `topic-page: extended` unless ALL of
the following are true at completion:

1. `.context/kb/topics/<slug>/index.md` (or a sibling sub-page
   like `.context/kb/topics/<slug>/<sub>.md`) exists and was
   created or extended in this pass.
2. The page cites at least one `EV-###` row that resolves to
   `evidence-index.md`.
3. `ctx kb site build` ran clean (or its failure is named in
   the closeout's `Next pass hint` AND the pass reports
   `topic-page: deferred`).
4. The cold-reader orientation rubric records
   **`Result: pass`** in the closeout's `What changed` section.
   All four rubric items must be `yes`.

Any failure → `topic-page: deferred` and the source-coverage
ledger advances to `topic-page-drafted` (not `comprehensive`).

### Source-Coverage Ledger (State Machine)

`.context/kb/source-coverage.md` is a state machine over every
source the kb has touched. States and allowed transitions:

| state                   | meaning                                                                  | next states |
|-------------------------|--------------------------------------------------------------------------|-------------|
| `discovered`            | surfaced via discovery / candidate harvest; not yet admitted             | → `admitted`, → `skipped` |
| `admitted`              | admitted against scope; not yet extracted                                | → `highlights-extracted`, → `partially-ingested`, → `topic-page-drafted`, → `comprehensive` |
| `highlights-extracted`  | EV rows minted for highlights only; no topic page or stub only           | → `partially-ingested`, → `topic-page-drafted`, → `comprehensive` |
| `partially-ingested`    | topic page exists but is incomplete relative to source                   | → `topic-page-drafted`, → `comprehensive` |
| `topic-page-drafted`    | page exists; cold-reader recorded; Confidence `speculative` or `TBD-cite` remains | → `comprehensive` |
| `comprehensive`         | page complete; cited bands ≥ `medium`; no `TBD-cite`; cold-reader passed | terminal until source updates |
| `superseded`            | source replaced by a newer canonical version; superseder named           | terminal |
| `skipped`               | source admitted-then-rejected as out-of-scope; reason cited              | terminal until scope changes |

**Row schema** (Markdown table):

```
| Source       | Topic        | State                 | EV coverage     | Residue      | Next action                       | Updated    |
| CURSOR-HOOKS | cursor-hooks | highlights-extracted  | EV-018..EV-034  | examples,... | /ctx-kb-ingest cursor-hooks       | 2026-05-13 |
```

Every pass that touches a source updates the ledger before
writing the closeout. **Lying to the ledger is a hard
anti-pattern.**

### Topic-Adjacency Pre-Flight

Before resolving the topic, scan the ledger for rows whose
state is **not** in `{comprehensive, skipped, superseded}` AND
whose `Topic` is plausibly adjacent. Heuristic:

- **Shared first segment of a slash- or hyphen-separated slug**:
  `cursor/skills` is adjacent to `cursor/hooks`;
  `your-domain-foo` is adjacent to
  `your-domain-bar`.
- **Shared product / vendor / surface in source URL or
  description**: sources under `cursor.com/docs/*` are
  adjacent.
- **Explicit cross-references** in the named topic's existing
  sub-pages or this pass's source set.

For each adjacent incomplete topic surfaced, the pass MUST:

1. Acknowledge it in `## Related concepts in this kb` on the
   topic page being authored.
2. Surface it in the closeout's `Adjacency pre-flight` block.
3. Surface it in the response contract's
   `Adjacent topics noted` field.

**Do NOT enumerate `EV-###` IDs by name in the adjacency block.**
Use *count + location* (*"seventeen rows in
`evidence-index.md`"*); naming an EV row from a lower-
confidence sibling demotes the floor of cited bands.

Silence is not the same as a clean pre-flight; explicit
*"no incomplete adjacent topics surfaced"* required when zero
matches.

### Cold-Reader Orientation Rubric

Four yes/no items recorded in the closeout's `What changed`
section, in `topic-page` mode:

```
Cold-reader orientation:
- Concept clear?                yes|no: <short note>
- Why this kb cares clear?      yes|no: <short note>
- Canonical evidence reachable? yes|no: <short note>
- Boundaries clear?             yes|no: <short note>
Result: pass | fail
```

`Result: pass` requires all four `yes`. Any `no` → `Result:
fail` → circuit-breaker fails → `topic-page: deferred`.

### Life-Stage Check

Count `.context/kb/topics/*/index.md` before this pass begins
synthesizing:

- `< 5` topic pages → **bootstrap** mode. Skip reconciliation
  ceremony; synthesize topic pages aggressively. Exception:
  surface a contradiction even in bootstrap if the new
  material plainly contradicts existing kb claims.
- `>= 5` topic pages → **maintenance** mode. Apply full
  reconciliation discipline (laddering, demotion,
  contradictions).

Document the life-stage call in the closeout's frontmatter
(`life-stage:`) and `What changed` section.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Empty input to `ctx kb ingest` | Refuse cleanly: `no sources provided; pass a folder, a URL, an MCP resource, or describe the materials inline.` Non-zero exit. |
| Empty question to `ctx kb ask` | Refuse cleanly: `no question provided; pass a question or describe it inline.` Non-zero exit. |
| Empty `grounding-sources.md` on `ctx kb ground` | Skill prompts once for sources before running; `NONE` on a line is a per-pass skip (re-prompts next invocation). |
| Concurrent ingest writers producing duplicate `EV-###` | Doctor advisory detects duplicates on next `ctx doctor` run; LLM cleanup pass renumbers and rewrites cross-references (P1). Documented as single-writer convention. |
| Temporal misordering: ingest today's transcript before last week's | Pipeline detects date-stamped filename; demotion policy applies temporal-precedence rule (newer-occurred carries forward over older-occurred even if older was extracted later); doctor warns when `dated:` source has rows missing `occurred:`. |
| Mid-session checkpoint via handover | `ctx handover write --no-fold` writes the handover without consuming closeouts (rare; default is fold). |
| Session aborted before wrap-up | Closeouts stay in place; next session's `/ctx-remember` reads handover **plus** unfolded postdated closeouts. Editorial work survives. |
| Mid-pass mode-switching detected | Skill aborts with partial closeout citing the mismatch; recommends re-invocation under the correct mode. **Never** silent-switches. |
| Cold-reader rubric returns `Result: fail` | Pass reports `topic-page: deferred` AND `validation: deferred (cold-reader orientation failed)`; ledger advances to `topic-page-drafted` (not `comprehensive`); closeout names which rubric items returned `no`. |
| Adjacency pre-flight surfaces zero matches | Closeout records explicit *"no incomplete adjacent topics surfaced"* and response-contract field reads `none surfaced`. Silence is not allowed. |
| Inferring `evidence-only` from source size / time pressure | Hard anti-pattern. Skill refuses to set `evidence-only` without explicit user trigger. |
| Lying to the source-coverage ledger | Doctor advisory: cross-check ledger rows against file existence + last-modified time vs. row's `Updated`. Mismatch → advisory. |
| `zensical` missing on PATH for `ctx kb site` | Single-line install hint, non-zero exit. No interactive install attempt. |
| `.context/kb/` missing when read-side surfaces look | Skip the "kb state" branch silently; behave as if no KB exists. Mode-awareness is `if exists`, not `must exist`. |
| Speculative-confidence claim ships to rendered site | Render filter excludes `confidence: speculative` content; `low`-confidence content ships only when paired with matching `outstanding-questions.md` entry. |
| Hand-edit to `.context/ingest/INBOX.md` | Silently discarded on next mode-skill run (skill rewrites the inbox). To configure, edit `grounding-sources.md`. Documented in `KB-RULES.md`. |
| `ctx kb note` while no ingest pipeline directory exists | Refuse: `kb not initialized; run \`ctx init\` first`. Non-zero exit. |
| Ingest source path doesn't exist | Refuse cleanly with the missing path; do not attempt partial extraction. |
| Sub-page split needed (index.md fails cold-reader "boundaries clear?") | Skill proposes the split once; waits for user confirmation; never auto-splits. |
| `ctx kb topic new` for a slug that already exists | Refuse: `topic <slug> already exists at .context/kb/topics/<slug>/index.md`. Non-zero exit. |

### Validation Rules

- `ctx handover write` enforces `--summary` and `--next` via
  `MarkFlagRequired`. Empty placeholder values (`TBD`,
  `see chat`, whitespace-only) are rejected by the CLI, not
  just by the skill text.
- Mode skills (`/ctx-kb-ingest`, `/ctx-kb-ask`,
  `/ctx-kb-ground`) refuse on empty input rather than
  prompting; refuse-on-empty is the contract.
- `ctx kb` writes only to `.context/kb/` and
  `.context/ingest/`. Path constants enforce this (rooted
  writes).
- `/ctx-kb-ingest` MUST emit the up-front pass-mode
  declaration before any source extraction begins. The skill
  body checks the declaration was emitted before proceeding;
  doctor advisory detects closeouts missing the `Pass-mode`
  body block.
- Topic-page mode passes MUST satisfy all four circuit-breaker
  invariants before claiming `produced` / `extended`. Failure
  on any → `topic-page: deferred`.
- Closeout files require frontmatter fields `sha`, `branch`,
  `mode`, `pass-mode`, `life-stage`, `generated-at`. Missing
  any → site-review mode flags it; handover fold skips it with
  a warning.
- Confidence band must be one of `high|medium|low|speculative`.
  Site-review mode coerces malformed capitalization
  (`High` → `high`); other malformations are flagged.
- Closeouts are append-never-rewrite. Archived closeouts are
  immutable.
- Source-coverage ledger row state must match an allowed
  transition from the prior state. Doctor advisory flags
  illegal transitions (e.g. `comprehensive → highlights-
  extracted` without an explicit `superseded` step).
- Topic pages live at `.context/kb/topics/<slug>/index.md`
  from day one. Sub-page split is **lazy** (only when
  `index.md` fails the cold-reader "boundaries clear?" check)
  and **proposed** (one question, wait for confirmation).
- Tasks, decisions, learnings, conventions, reminders
  unchanged: they remain authored by their existing canonical
  CLIs. KB cannot write to canonical files; canonical CLIs
  cannot write to `.context/kb/`.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| `ctx kb` invoked but `.context/` missing | `context directory not found; run \`ctx init\` first.` | Run `ctx init` |
| `ctx kb` invoked but `.context/ingest/` missing (initialized before this spec shipped) | `kb pipeline not initialized; run \`ctx init --upgrade\` to lay down ingest scaffolding.` | Run `ctx init --upgrade` (idempotent on existing files; refuses to overwrite divergent content) |
| `ctx kb` invoked but kb scope undeclared (placeholder in `.context/kb/index.md`) | `kb scope is undeclared. Open .context/kb/index.md and replace the TODO placeholder with a one-paragraph scope statement.` | Hand-edit `.context/kb/index.md` |
| `ctx handover write` with empty `--summary` or `--next` | `--summary and --next are required and must be non-trivial; placeholder values like 'TBD' are rejected.` | Re-run with concrete values |
| `ctx kb ingest` with non-existent path | `source not found: <path>` | Fix path or pass `--inline "..."` description |
| `ctx kb topic new` for existing slug | `topic <slug> already exists at .context/kb/topics/<slug>/index.md` | Use existing folder; skill will append/extend |
| `ctx kb site build` when zensical missing | `zensical not on PATH; install per https://zensical.org/ then re-run.` | Install zensical |
| Closeout fold encounters malformed frontmatter | `warning: skipping malformed closeout (no generated-at): <name>` (proceeds with valid ones) | Hand-edit the malformed file or delete it |
| Doctor advisory: duplicate EV-### detected | Advisory line listing the dupe IDs and files. Non-fatal. | Run an LLM cleanup pass (agent renumbers + rewrites references) |
| Doctor advisory: dated source has rows missing `occurred:` | Advisory line listing the source short-name and row IDs. Non-fatal. | Hand-edit or re-run ingest with corrected source-map dating |
| Doctor advisory: source-coverage ledger mismatch (row's `Updated` predates file mtime) | Advisory line listing the row + file mtime. Non-fatal. | Re-run a pass that touches the source to refresh the row |
| Doctor advisory: closeout missing Pass-mode body block | Advisory line listing the closeout. Non-fatal. | Hand-edit the closeout to add the body block matching the frontmatter `pass-mode:` field |

## Interface

### CLI

```
ctx handover write <title> --summary "..." --next "..." [--highlights "..."] [--open-questions "..."] [--no-fold] [--commit <sha>]
ctx kb ingest <folder|paths...>            # mode-aware editorial pass
ctx kb ask "<question>"                    # Q&A from existing KB; read-only on prose
ctx kb site-review                         # mechanical structural audit
ctx kb ground                              # external grounding via grounding-sources.md
ctx kb note "..."                          # lightweight capture; appends to ingest/findings.md
ctx kb topic new "<name>"                  # scaffold a topic folder (sole writer of topic index.md)
ctx kb reindex                             # refresh kb/index.md's CTX:KB:TOPICS managed block
ctx kb site build                          # render kb to .context/site/kb/
ctx kb site serve [--addr :8000]           # build + serve
ctx kb site customize                      # lazy-init .context/site-config/kb.toml
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--summary` | string | (required) | What happened this session, past tense |
| `--next` | string | (required) | What the next agent should do FIRST, specific |
| `--highlights` | string | "" | Notable artifacts produced this session |
| `--open-questions` | string | "" | Things that remain undecided |
| `--no-fold` | bool | false | Skip closeout consumption (mid-session checkpoint) |
| `--commit` | string | (resolved) | Override resolved git HEAD for Provenance line |
| `--addr` | string | `:8000` | Bind address for `ctx kb site serve` |
| `--inline` | string | "" | Inline natural-language source description for `ctx kb ingest` (for cases with no folder) |
| `--mode` | string | "" | Optional explicit pass-mode override (`topic-page`/`triage`/`evidence-only`). Default is skill-detected. CLI form respects same precedence rules. |

### Skill

```
/ctx-handover <title>           # deploys via ctx setup; wraps ctx handover write
/ctx-kb-ingest <folder|paths>   # wraps ctx kb ingest; declares pass-mode up front
/ctx-kb-ask "<question>"        # wraps ctx kb ask
/ctx-kb-site-review             # wraps ctx kb site-review
/ctx-kb-ground                  # wraps ctx kb ground
/ctx-kb-note "..."              # wraps ctx kb note
```

Trigger phrases (per-skill):

- `/ctx-wrap-up`: "let's wrap up", "save context", "save state", "leave a handover", "before I go", "stepping away", "end of session". Owns the user-facing session-end trigger; always delegates to `/ctx-handover` as its final step. `/ctx-handover` is a sub-mechanism, not a user-facing trigger; legitimate direct-invocation cases are `--no-fold` mid-session checkpoint and recovery after an aborted session.
- `/ctx-kb-ingest`: "ingest the transcripts", "pull this into the kb", "add evidence from", explicit slash form for paths.
- `/ctx-kb-ask`: "does the kb say", "according to evidence", explicit slash form for the question.
- `/ctx-kb-site-review`: "audit the kb", "check kb for rot".
- `/ctx-kb-ground`: "re-ground the kb", "check upstream", explicit slash form preferred.
- `/ctx-kb-note`: "drop a note", "capture this for the next ingest", "park this finding".

`/ctx-wrap-up` skill is updated to: (a) **always** delegate to
`/ctx-handover` as its final step regardless of whether
`.context/kb/` exists or capture outcomes, (b) when
`.context/kb/` exists, additionally surface editorial state
(pending closeouts, outstanding-questions count) in its
summary before delegating. KB presence affects what gets
folded, not whether the handover is written.

`/ctx-remember` skill is updated to: always read the latest
handover (orthogonal to KB; it is the former-agent-to-next-agent
note left by the previous wrap-up); when `.context/kb/`
exists, additionally fold any postdated unfolded closeouts
into the readback.

## Implementation

### Files to Create / Modify

| File | Change |
|------|--------|
| `internal/gitmeta/require.go` | **NEW** (Phase RG): `RequireGitTree(projectRoot)` + typed `MissingGitError`. Phase KB depends on Phase RG; closeout/handover provenance requires honest git state. |
| `internal/gitmeta/resolvehead.go` | **NEW** (Phase RG): `ResolveHead(projectRoot)` returns commit short SHA + branch + override-via-env (`CTX_TASK_COMMIT`, `GITHUB_SHA`). Removes the `commit:none` fallback (state unreachable). |
| `internal/cli/handover/` | **NEW**: parent + `cmd/write/cmd.go`. `MarkFlagRequired` on `--summary`, `--next` with placeholder rejection (Phase SK pattern, already in `internal/validate/`). |
| `internal/cli/kb/` | **NEW**: parent `cmd.go`, `doc.go`, `run.go`. Subcommand subdirs: `cmd/ingest/`, `cmd/ask/`, `cmd/sitereview/`, `cmd/ground/`, `cmd/note/`, `cmd/topicnew/`, `cmd/reindex/`, `cmd/site/build/`, `cmd/site/serve/`, `cmd/site/customize/`. Refuse-on-empty for ingest/ask/ground. |
| `internal/cli/kb/core/` | **NEW**: shared helpers (`path/path.go` for KB path constants, `passmode/` for mode declaration + validation, `circuitbreaker/`, `ledger/` for source-coverage state machine, `adjacency/`, `coldreader/`, `lifestage/`). |
| `internal/write/handover/` | **NEW**: `WriteHandover`, `LatestHandoverCursor`, `UnconsumedCloseouts`, `ArchiveCloseouts`. Mirrors the upstream shape. |
| `internal/write/closeout/` | **NEW**: `WriteCloseout` with required frontmatter (`sha`, `branch`, `mode`, `pass-mode`, `life-stage`, `generated-at`); cursor-extracting reader. |
| `internal/write/kb/` | **NEW**: per-artifact writers (evidence-index append-never-renumber; glossary; contradictions; outstanding-questions; domain-decisions; timeline; source-map; relationship-map; source-coverage state-machine API); demotion API; `EvidenceRow` includes `occurred:` field; `TopicScaffold` writer (called only by `ctx kb topic new`). |
| `internal/cli/initialize/core/kb/` | **NEW**: scaffolding helper for `ctx init` to lay down `.context/kb/topics/.gitkeep`, `.context/ingest/`, `.context/handovers/`, `.context/site/`. |
| `internal/cli/initialize/cmd/root/cmd.go` | Modify: add `--upgrade` flag (idempotent on byte-identical content; refuse on divergent). |
| `internal/assets/embed.go` | Modify: add `//go:embed kb/templates/ingest/*.md kb/templates/ingest/schemas/*.md kb/templates/.gitkeep` lines. |
| `internal/assets/kb/templates/ingest/*.md` | **NEW**: embedded templates: `KB-RULES.md`, `00-GROUND.md`, `30-INGEST.md`, `40-ASK.md`, `50-SITE_REVIEW.md`, `INBOX.md`, `SESSION_LOG.md`, `grounding-sources.md`, `OPERATOR.md`, `PROMPT.md`. |
| `internal/assets/kb/templates/ingest/schemas/*.md` | **NEW**: embedded schemas: `evidence-index.md`, `glossary.md`, `contradictions.md`, `outstanding-questions.md`, `domain-decisions.md`, `timeline.md`, `source-map.md`, `source-coverage.md` (state-machine row format), `relationship-map.md`, `session-log.md`. Each carries fields list + one worked example, no domain content. |
| `internal/assets/kb/templates/kb/index.md` | **NEW**: kb landing page template with `CTX:KB:TOPICS` managed block + `## Scope` placeholder. |
| `internal/assets/kb/templates/kb/topics/_template/index.md` | **NEW**: topic-page template (Status block, lede, "What it is", "Why this kb cares", "Sources and further reading", "Related concepts in this kb"). |
| `internal/assets/claude/skills/ctx-handover/SKILL.md` | **NEW** with input contract, authority boundary, edge cases per spec. |
| `internal/assets/claude/skills/ctx-kb-ingest/SKILL.md` | **NEW**: mode-aware ingest driver; pass-mode declaration; 4 circuit-breaker invariants; adjacency pre-flight; cold-reader rubric; refuses on empty input. |
| `internal/assets/claude/skills/ctx-kb-ask/SKILL.md` | **NEW**: Q&A from KB; refuses on empty question. |
| `internal/assets/claude/skills/ctx-kb-site-review/SKILL.md` | **NEW**: mechanical audit; no arguments. |
| `internal/assets/claude/skills/ctx-kb-ground/SKILL.md` | **NEW**: external grounding; reads `grounding-sources.md`. |
| `internal/assets/claude/skills/ctx-kb-note/SKILL.md` | **NEW**: lightweight capture. |
| `internal/assets/claude/skills/ctx-wrap-up/SKILL.md` | Modify: branch on `.context/kb/` existence; mandatorily drive `/ctx-handover` as final step. |
| `internal/assets/claude/skills/ctx-remember/SKILL.md` | Modify: read latest handover + any postdated unfolded closeouts; fold KB state into readback if `.context/kb/` exists. |
| `internal/cli/doctor/core/advisory.go` (or equivalent) | Extend: duplicate-`EV-###` detection; `dated:` sources missing `occurred:` check; malformed-closeout-frontmatter detection; **source-coverage-ledger-mismatch** detection (row's Updated vs. file mtime); **closeout-missing-pass-mode-body-block** detection; **illegal-ledger-state-transition** detection. |
| `internal/cli/setup/cmd.go` | Already walks skills dir; new skill subdirs picked up automatically. |
| `internal/cli/wrapup/run.go` (or equivalent) | Modify: `if KBExists() { printPendingCloseouts; printOutstandingQuestionsCount }`. Mandatory handover wording remains in skill, not CLI. |
| Project root `.gitignore` | Append: `.context/site/` (idempotent; match existing pattern for `.context/journal/.imported.json`). |
| `hack/smoke-kb.sh` | **NEW**: end-to-end shell smoke (init → kb ingest → kb ask → kb site-review → kb ground → handover write → archive populated → doctor clean). |

### Key Functions

```go
// internal/gitmeta/require.go
type MissingGitError struct{ ProjectRoot string }
func RequireGitTree(projectRoot string) error

// internal/gitmeta/resolvehead.go
type HeadRef struct{ SHA, Branch string }
func ResolveHead(projectRoot string) (HeadRef, error) // honors CTX_TASK_COMMIT, GITHUB_SHA

// internal/write/handover/handover.go
type HandoverEntry struct {
    Title           string
    Summary         string
    Next            string
    Highlights      string
    OpenQuestions   string
    Commit          string
    Branch          string
    FoldedCloseouts []closeout.File
}
func WriteHandover(projectRoot string, e HandoverEntry) (Result, error)
func LatestHandoverCursor(projectRoot string) (cursor time.Time, file string, err error)
func UnconsumedCloseouts(projectRoot string, cursor time.Time) (consumed []closeout.File, malformed []string, err error)
func ArchiveCloseouts(projectRoot string, files []closeout.File) error

// internal/write/closeout/closeout.go
type Frontmatter struct {
    SHA, Branch, Mode, PassMode, LifeStage string
    GeneratedAt                             time.Time
}
type File struct {
    Path        string
    Frontmatter Frontmatter
    Body        string
}
func WriteCloseout(projectRoot, mode, passMode, lifeStage string, body string) (File, error)

// internal/write/kb/evidence.go
type EvidenceRow struct {
    ID         string // EV-###
    Claim      string
    SourceID   string
    Locator    string
    SHA        string // optional; for in-repo citations
    Confidence string // high|medium|low|speculative
    Tags       []string
    Occurred   *time.Time // optional; temporal-precedence rule
    Extracted  time.Time
}
func AppendEvidence(projectRoot string, row EvidenceRow) error  // never renumber
func DemoteEvidence(projectRoot, id, newBand, reason string) error

// internal/write/kb/sourcecoverage.go
type LedgerState string
const (
    StateDiscovered          LedgerState = "discovered"
    StateAdmitted            LedgerState = "admitted"
    StateHighlightsExtracted LedgerState = "highlights-extracted"
    StatePartiallyIngested   LedgerState = "partially-ingested"
    StateTopicPageDrafted    LedgerState = "topic-page-drafted"
    StateComprehensive       LedgerState = "comprehensive"
    StateSuperseded          LedgerState = "superseded"
    StateSkipped             LedgerState = "skipped"
)
type LedgerRow struct {
    Source, Topic                string
    State                        LedgerState
    EVCoverage, Residue          string
    NextAction                   string
    Updated                      time.Time
}
func AdvanceLedger(projectRoot string, row LedgerRow) error // validates allowed transition
func ReadLedger(projectRoot string) ([]LedgerRow, error)
func ValidTransition(from, to LedgerState) bool

// internal/cli/kb/core/passmode/passmode.go
type Mode string
const (
    ModeTopicPage    Mode = "topic-page"
    ModeTriage       Mode = "triage"
    ModeEvidenceOnly Mode = "evidence-only"
)
type Declaration struct {
    Mode            Mode
    Reason          string
    DefinitionOfDone string
}
func RenderDeclaration(d Declaration) string // 3-line markdown block

// internal/cli/kb/core/circuitbreaker/check.go
type Invariants struct {
    PageExists       bool
    CitesEV          bool
    SiteBuildClean   bool
    ColdReaderPass   bool
}
func (i Invariants) AllPassed() bool
func Check(projectRoot, slug string, coldReader ColdReader) (Invariants, error)

// internal/cli/kb/core/coldreader/rubric.go
type Rubric struct {
    ConceptClear            string // yes|no: note
    WhyClear                string
    EvidenceReachable       string
    BoundariesClear         string
}
func (r Rubric) Result() string // "pass" if all yes; "fail" otherwise

// internal/cli/kb/core/lifestage/lifestage.go
type Stage string
const (
    StageBootstrap   Stage = "bootstrap"
    StageMaintenance Stage = "maintenance"
)
func Detect(projectRoot string, threshold int) (Stage, int, error) // threshold default 5
```

### Helpers to Reuse

- `internal/validate/`: already has `RequireBodyFlags` (Phase SK)
  for `MarkFlagRequired` + placeholder rejection. Reuse for
  `ctx handover write`.
- `internal/cli/journal/cmd/site/`: existing zensical shell-out
  pattern; lift the build/serve runtime-config materialization
  wholesale for `ctx kb site`.
- `internal/config/zensical/`: existing mkdocs.go / toml.go;
  reuse for kb-site config rendering.
- `internal/cli/initialize/core/`: existing template walk for
  embedded asset deployment.
- `internal/cli/setup/`: existing skill-dir walker; new skill
  subdirs land for free.
- `internal/cli/doctor/`: existing advisory pattern; add new
  check functions following the same shape.
- `internal/assets/embed.go`: existing `//go:embed` block;
  extend with kb-template lines.

## Configuration

No new `.ctxrc` keys in v1. Site-config (`.context/site-config/kb.toml`)
is lazy-initialized via `ctx kb site customize` for users who
want to override theme / nav / plugins. Infrastructure paths
(`docs_dir`, `site_dir`) are wrapper-owned and overwritten at
build time per the existing `ctx journal site` pattern.

Environment variables honored:

- `CTX_TASK_COMMIT`: override resolved commit for handover
  Provenance (CI replay).
- `GITHUB_SHA` (when `GITHUB_ACTIONS=true`): same purpose.

`ZENSICAL_BIN` is **not** introduced; the binary is resolved
from PATH per the existing journal-site convention.

## Testing

### Unit

- `internal/gitmeta/require_test.go`: `.git` dir → nil;
  `.git` file (worktree pointer) → nil; absent → typed error;
  override allow-list (root help-shaped commands) skips the
  check.
- `internal/gitmeta/resolvehead_test.go`: happy path; env
  override via `CTX_TASK_COMMIT`; `GITHUB_SHA` honored only
  with `GITHUB_ACTIONS=true`.
- `internal/write/handover/handover_test.go`:
  `WriteHandover` happy path; rejects empty summary/next;
  rejects placeholder bodies; closeout fold cursor logic;
  archive moves files atomically.
- `internal/write/closeout/closeout_test.go`: frontmatter
  parse (all required fields including `pass-mode`,
  `life-stage`); malformed frontmatter handling;
  `generated-at` ordering.
- `internal/write/kb/evidence_test.go`: evidence-index
  append never renumbers; demotion bands valid only within
  `high|medium|low|speculative`; concurrent-write detection
  surfaces in doctor.
- `internal/write/kb/sourcecoverage_test.go`: every allowed
  transition; every disallowed transition refused;
  `ValidTransition(from, to)` truth table; row schema parse
  and write round-trip.
- `internal/cli/kb/core/passmode/passmode_test.go`:
  declaration render; mode set validation; rejects invalid
  modes; rejects empty reason on non-default modes.
- `internal/cli/kb/core/circuitbreaker/check_test.go`: all
  four invariants must hold; any one failing produces a
  reportable failure with the specific item named.
- `internal/cli/kb/core/coldreader/rubric_test.go`: any `no`
  → `fail`; all `yes` → `pass`.
- `internal/cli/kb/core/lifestage/lifestage_test.go`: count
  topic folders; threshold boundary (4 → bootstrap, 5 →
  maintenance).
- `internal/cli/kb/<mode>/cmd_test.go`: refuse-on-empty for
  ingest / ask / ground; refuse on existing slug for topic new.
- `internal/cli/handover/cmd/write/cmd_test.go`:
  `MarkFlagRequired` enforcement; `--no-fold` behavior;
  provenance from `gitmeta.ResolveHead`.
- `internal/cli/doctor/core/advisory_test.go`:
  duplicate-`EV-###`; `dated:`-without-`occurred:`;
  malformed-closeout; source-coverage-ledger-mismatch;
  closeout-missing-pass-mode-body-block;
  illegal-ledger-state-transition.

### Integration

- `internal/cli/initialize/init_test.go`: full init creates
  all new dirs and templates; `--upgrade` is idempotent on
  byte-identical existing content; `--upgrade` refuses on
  divergent existing content.
- `internal/cli/setup/setup_test.go`: new skill subdirs deploy.
- `hack/smoke-kb.sh`: end-to-end shell smoke. `git init &&
  ctx init`; `ctx kb ingest ./testdata/inputs`; `ctx kb ask
  "..."`; `ctx kb site-review`; `ctx kb ground`; `ctx
  handover write --summary X --next Y`; verify files exist,
  closeouts folded, archive populated, no doctor errors.

### Edge Cases

- Aborted-session recovery: write a closeout, do NOT write a
  handover, simulate session restart by re-reading the state;
  verify `/ctx-remember`'s read path picks up the unfolded
  closeout (test the writer helper directly).
- Temporal misordering: ingest fixture A (occurred
  2026-04-12, extracted 2026-05-09) then fixture B (occurred
  2026-04-05, extracted 2026-05-09); verify demotion does NOT
  fire because temporal-precedence rule wins.
- Concurrent dupe IDs: simulate two parallel writers
  producing `EV-020` against different claims; verify
  `ctx doctor` flags the dupe.
- Render filter: speculative content does NOT appear in built
  HTML; low-confidence content appears only when paired with
  outstanding-questions row.
- Mid-pass mode-switch attempted: declared `topic-page`,
  agent realizes it spans 5 topics, must abort with partial
  closeout rather than silently widen.
- Adjacency pre-flight zero matches: must surface
  `none surfaced` explicitly in response contract +
  closeout's `Adjacency pre-flight` block.
- Cold-reader fail: integration test ingests a deliberately
  vague source, verifies rubric reports `fail`, ledger
  advances to `topic-page-drafted` not `comprehensive`,
  closeout names which rubric items failed.

### Validation Corpus

`your-project` is the live regression suite
(older shape, hand-rolled) and the structural reference
(current upstream shape applied to a different domain).
Phase KB-2 stands up a new research workspace using the
shipped `ctx kb` tool; each divergence from manual is either
a Phase KB bug or a `DECISIONS.md` entry explaining why the
formal shape differs from what worked manually.

## Non-Goals

(Explicit deferrals, referenced from the brief's "What we
rejected" table.)

- **No `/ctx-kb-decide` skill.** KB ontology rejects it: in a
  KB you don't decide, you increase confidence. Pipeline is
  the sole writer; ad-hoc capture flows through `ctx kb note`
  or hand-edit.
- **No team write-coordination layer.** Single-writer
  convention with doctor advisory + LLM cleanup is the v1
  stance. Multi-writer coordination is deferred until a real
  team-scale user hits the wall.
- **No UUIDs for evidence rows.** `EV-###` aesthetic preserved.
- **No KB-scoped IDs (`research-master/EV-019`) in v1.** P2
  federation handles multi-KB without scoping IDs.
- **No domain-split per user (`.context/kb/<domain>/`).** Single
  KB is simpler; multi-domain via P2 federation when needed.
  Topic folders provide intra-kb structure.
- **No KB-side merge into canonical files.** `domain-decisions.md`
  stays separate from `DECISIONS.md`. Different schema,
  different write authority, different lifecycle.
- **No automatic demotion-cascade in v1.** When EV-031
  demotes, affected glossary / domain-decision / timeline
  rows are NOT auto-flagged; the human handles the cascade.
  v2 may add site-review automation.
- **No bundled renderer.** `zensical` is shelled out, not
  vendored or wrapped. Same model as `ctx journal site`.
- **No KB linting.** Confidence-band discipline is
  rule-driven (per `KB-RULES.md`), not enforced
  programmatically.
- **No bulk migration tooling for repos initialized before
  this spec.** `ctx init --upgrade` lays down the new dirs
  idempotently; pre-existing canonical files are untouched;
  hand-rolled editorial files (e.g. `your-project`'s
  `10-CONSTITUTION.md` at repo root) are left alone;
  porting is a manual cutover (Phase 2 of validation).
- **No interactive install of zensical.** Missing-binary case
  fails with a one-line install hint and non-zero exit.
- **No replacement of `/ctx-decision-add`,
  `/ctx-learning-add`, `/ctx-task-add`, `/ctx-convention-add`,
  `/ctx-wrap-up`.** New skills are siblings; existing
  capture skills unchanged in authority.
- **No mid-pass mode-switching.** Mode is committed at
  declaration time; if work no longer fits, abort and
  re-invoke. Silent switching is a hard anti-pattern.
- **No pre-emptive sub-page split.** Topic pages stay as
  `index.md` until the cold-reader "boundaries clear?" check
  fails; only then is a split proposed.

## Failure Analysis

Three concrete ways the lifted shape fails badly, with
mitigations baked into the spec:

1. **Pass-mode contract gets ignored under operator
   pressure.** Anti-pattern: inferring `evidence-only` to
   dodge topic-page validation. Mitigation: declaration is
   logged to the closeout's frontmatter (`pass-mode:`) AND
   to the closeout body block (redundancy makes
   false-finish drift visible); doctor advisory detects
   closeouts whose body `Pass-mode` block disagrees with the
   frontmatter `pass-mode:` field.
2. **Source-coverage ledger drifts from reality.**
   Anti-pattern: "lying to the ledger." Mitigation: doctor
   advisory cross-checks ledger rows against file existence
   + last-modified time vs. row's `Updated` cell. Mismatch
   → advisory line. Additionally, illegal state transitions
   (e.g. `comprehensive → highlights-extracted` without an
   explicit `superseded` step) are refused at write time by
   `AdvanceLedger`.
3. **Adjacency pre-flight degenerates into trivia.** A
   surfaced sibling pasted as a footnote satisfies the
   letter, not the spirit. Mitigation: the pre-flight
   result is a structured field in the closeout: the
   `Adjacent topics noted` field must be either
   `none surfaced` or a slug-list, never free prose.
   Doctor parses the field; free-prose values fail
   validation.

These mitigations are part of v1, not v2.

## Open Questions

(Carried forward from the brief; require pinning during
implementation.)

1. **Naming.** `ctx kb ingest|ask|site-review|ground`
   (kb-prefixed) vs `ctx ingest|ask|site-review|ground`
   (top-level). Lean prefixed. Confirm during implementation
   kickoff.
2. **Brief vs spec storage.** Where do briefs live?
   `.context/briefs/` (debate residue, distinct lifecycle) vs
   `.context/specs/briefs/` (subdir of specs). Lean
   dedicated `briefs/`. Tied to the polish-PR work in
   `ideas/002` §3.
3. **`ctx kb note` destination.** Single
   `.context/ingest/findings.md` or one file per invocation?
   Lean single file (simpler); per-invocation preserves
   provenance per note.
4. **Pure-research project init.** Should `ctx init` learn
   `--research` / `--kb-only` to suppress code-dev skill
   deployment? Light defer; not v1 critical. Things-wtf
   workaround (CLAUDE.md disabling) remains acceptable for v1.
5. **Confidence bands flowing into `LEARNINGS.md`.** Probably
   no (different truth bases: KB has citations, learnings
   have author intent). Confirm rather than assume.
6. **`relationship-map.md` vs GitNexus.** Different graphs;
   v1 keeps independent. Cross-feed is v2.
7. **Demotion-policy automation.** Auto-flag affected pages
   on demotion? v1 defers to human; v2 may automate via
   site-review.
8. **`--no-fold` flag scope.** Handover-only (the upstream's
   choice) vs every artifact-writing command. Lean
   handover-only.
9. **Life-stage threshold.** 5 topic pages matches the upstream default.
   Confirm or override.
10. **Path constants location.** `internal/cli/kb/core/path/`
    (per-subcommand pattern matching `internal/cli/task/core/path/`)
    vs a new top-level `internal/path/`. Lean
    per-subcommand (matches existing `ctx` convention).

---

## Source

`ideas/003-editorial-pipeline-debated-brief.md` (debated brief
from `/ctx-plan` session 01d0cf92, 2026-05-09).

`ideas/upstream-pipeline-comparison.md` (revision-2
input; 2026-05-16) is the comparison note that surfaced the
deltas absorbed in this revision.

Inputs to the brief:

- `ideas/001-sibling-project-undercover-analysis.md`:
  handover mechanism, closeout/fold mechanism, doctor
  advisory tier.
- `ideas/002-editorial-pipeline-and-skill-rigor.md`:
  the lifted-pipeline plan + skill ceremony comparison.
- `your-project/` (upstream-reference workspace):
  live test corpus; hand-rolled version of the older 4-phase
  shape.
- `your-project` (upstream-reference workspace): structural
  reference for the current upstream shape applied to a
  different domain.
- the upstream reference's editorial-pipeline source tree:
  source templates and SKILL.md referenced for shape lift.
- `.context/journal-site/zensical.toml`: proof that `ctx`
  already shells out to zensical for journal rendering.
