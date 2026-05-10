# KB Editorial Pipeline (paired with handover mechanism)

> Source: `ideas/003-editorial-pipeline-debated-brief.md` (debated
> brief from `/ctx-plan` session 01d0cf92, 2026-05-09).
>
> Authority order when this spec, `DECISIONS.md`, or `docs/` could
> conflict: frozen contracts in `docs/` > recorded decisions in
> `DECISIONS.md` > the brief at `ideas/003-...` > inference labeled
> `TBD` here. Do not invert.

## Problem

`ctx` users doing knowledge-shaped work (research projects,
domain modeling, post-incident reviews, vendor-spec analysis)
currently have no native surface for it. The canonical files
(`TASKS.md`, `DECISIONS.md`, `LEARNINGS.md`, `CONVENTIONS.md`)
are tuned for code-development context, not for evidence-tracked
knowledge with confidence bands, contradictions, and external
grounding.

The cost is concrete and active:
`things-wtf-disaster-recovery` has been hand-rolling an
editorial pipeline at the repo root (`10-CONSTITUTION.md`,
numbered mode prompts, `20-INBOX.md`, hand-typed 8-item
closeouts) and disabling half of ctx's code-dev skill surface
in `CLAUDE.md` (`/ctx-commit`, `/ctx-implement`, `/ctx-spec`,
`/ctx-architecture`, `/ctx-brainstorm`, `/ctx-wrap-up`) to
avoid name collision with its editorial constitution.

Separately, ctx has no per-session **handover artifact**: the
narrative thread between sessions is currently re-derived
probabilistically from canonical files plus journal. An
ad-hoc one-off (`HANDOVER-2026-04-22.md`) exists but no
discipline shapes it.

This spec adds both surfaces and integrates them via a
closeout/fold mechanism so editorial work flows cleanly into
session-to-session continuity.

## Approach

Lift the validated shape from a sibling tool's editorial
pipeline (4 modes, 9 KB artifacts, closeout/fold, browseable
site) into `ctx` with a non-colliding rename
(`KB-RULES.md`, not `CONSTITUTION.md`) and pair it with a
per-session handover mechanism. The closeout/fold is the
integration point: every editorial pass writes a closeout
with a `generated-at` cursor; the next handover folds
postdated closeouts and archives the sources.

**Two organizing principles** (carried forward verbatim from the
brief because they explain why this is rationally bold rather
than recklessly broad):

- **P1 — The LLM is the migration tool.** Wholesale ID
  renumbering, taxonomy reshuffles, confidence-band remapping,
  cross-file reference rewrites are absorbed by an LLM cleanup
  pass. We commit to specific schemas in v1 instead of hedging
  with abstract types. Be wrong cheaply.
- **P2 — A KB is knowledge; a KB of KBs is a KB.** Recursive
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

1. User runs `ctx init` in a fresh project. Init lays down the
   five canonical files **plus** `.context/handovers/`,
   `.context/kb/` (with `.gitkeep`), `.context/ingest/` (with
   `KB-RULES.md`, mode prompts `00-GROUND.md`, `30-INGEST.md`,
   `40-ASK.md`, `50-SITE_REVIEW.md`, `INBOX.md`,
   `SESSION_LOG.md`, `grounding-sources.md`, `OPERATOR.md`,
   `PROMPT.md`, `schemas/`), and `.context/site/` (gitignored).
2. User runs `ctx setup` to deploy skills. New skills:
   `/ctx-handover`, `/ctx-kb-ingest`, `/ctx-kb-ask`,
   `/ctx-kb-site-review`, `/ctx-kb-ground`, `/ctx-kb-note`.
3. User invokes `/ctx-kb-ingest ./inputs/2026-04-12-call.md`.
   The skill runs the four-phase ingest (triage, extract,
   reconcile, surface), writes evidence rows to
   `.context/kb/evidence-index.md`, updates glossary /
   contradictions / outstanding-questions / timeline as
   needed, appends one line per phase boundary to
   `.context/ingest/SESSION_LOG.md`, and writes a closeout
   `.context/ingest/closeouts/<TS>-ingest-closeout.md` with
   `generated-at` frontmatter.
4. User runs `/ctx-wrap-up`. The skill (a) walks the standard
   capture checklist, (b) detects `.context/kb/` exists and
   surfaces editorial state (pending closeouts, unresolved
   outstanding-questions count), (c) **mandatorily** drives
   `/ctx-handover` as the final step.
5. `/ctx-handover` collects `--summary` (past tense, what
   happened) and `--next` (future tense, specific first
   action), runs `ctx handover write`, which folds postdated
   closeouts into `## Folded closeouts`, archives the source
   closeouts to `.context/archive/closeouts/`, and writes
   `.context/handovers/<TS>-<slug>.md`.
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

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Empty input to `ctx kb ingest` | Refuse cleanly: `no sources provided; pass a folder or describe the materials inline.` Non-zero exit. |
| Empty question to `ctx kb ask` | Refuse cleanly: `no question provided; pass a question or describe it inline.` Non-zero exit. |
| Empty `grounding-sources.md` on `ctx kb ground` | Skill prompts once for sources before running; `NONE` on a line is a per-pass skip (re-prompts next invocation). |
| Concurrent ingest writers producing duplicate `EV-###` | Doctor advisory detects duplicates on next `ctx doctor` run; LLM cleanup pass renumbers and rewrites cross-references (P1). Documented as single-writer convention. |
| Temporal misordering: ingest today's transcript before last week's | Pipeline detects date-stamped filename; demotion policy applies temporal-precedence rule (newer-occurred carries forward over older-occurred even if older was extracted later); doctor warns when `dated:` source has rows missing `occurred:`. |
| Mid-session checkpoint via handover | `ctx handover write --no-fold` writes the handover without consuming closeouts (rare; default is fold). |
| Session aborted before wrap-up | Closeouts stay in place; next session's `/ctx-remember` reads handover **plus** unfolded postdated closeouts. Editorial work survives. |
| `zensical` missing on PATH for `ctx kb site` | Single-line install hint, non-zero exit. No interactive install attempt. |
| `.context/kb/` missing when read-side surfaces look | Skip the "kb state" branch silently; behave as if no KB exists. Mode-awareness is `if exists`, not `must exist`. |
| Speculative-confidence claim ships to rendered site | Render filter excludes `confidence: speculative` content; `low`-confidence content ships only when paired with matching `outstanding-questions.md` entry. |
| Hand-edit to `.context/ingest/INBOX.md` | Silently discarded on next mode-skill run (skill rewrites the inbox). To configure, edit `grounding-sources.md`. Documented in `KB-RULES.md`. |
| `ctx kb note` while no ingest pipeline directory exists | Refuse: `kb not initialized; run \`ctx init\` first`. Non-zero exit. |
| Ingest source path doesn't exist | Refuse cleanly with the missing path; do not attempt partial extraction. |

### Validation Rules

- `ctx handover write` enforces `--summary` and `--next` via
  `MarkFlagRequired`. Empty placeholder values (`TBD`,
  `see chat`, whitespace-only) are rejected by the CLI, not
  just by the skill text.
- Mode skills (`/ctx-kb-ingest`, `/ctx-kb-ask`,
  `/ctx-kb-ground`) refuse on empty input rather than
  prompting; refuse-on-empty is the contract.
- `ctx kb` writes only to `.context/kb/` and
  `.context/ingest/`. Path constants in
  `internal/path/path.go` enforce this (rooted writes).
- Closeout files require frontmatter fields `sha`, `branch`,
  `mode`, `generated-at`. Missing any → site-review mode
  flags it; handover fold skips it with a warning.
- Confidence band must be one of `high|medium|low|speculative`.
  Site-review mode coerces malformed capitalization
  (`High` → `high`); other malformations are flagged.
- Closeouts are append-never-rewrite. Archived closeouts are
  immutable.
- Tasks, decisions, learnings, conventions, reminders
  unchanged: they remain authored by their existing canonical
  CLIs. KB cannot write to canonical files; canonical CLIs
  cannot write to `.context/kb/`.

### Error Handling

| Error condition | User-facing message | Recovery |
|-----------------|---------------------|----------|
| `ctx kb` invoked but `.context/` missing | `context directory not found; run \`ctx init\` first.` | Run `ctx init` |
| `ctx kb` invoked but `.context/ingest/` missing (initialized before this spec shipped) | `kb pipeline not initialized; run \`ctx init --upgrade\` to lay down ingest scaffolding.` | Run `ctx init --upgrade` (idempotent on existing files; refuses to overwrite divergent content) |
| `ctx handover write` with empty `--summary` or `--next` | `--summary and --next are required and must be non-trivial; placeholder values like 'TBD' are rejected.` | Re-run with concrete values |
| `ctx kb ingest` with non-existent path | `source not found: <path>` | Fix path or pass `--inline "..."` description |
| `ctx kb site build` when zensical missing | `zensical not on PATH; install per https://zensical.org/ then re-run.` | Install zensical |
| Closeout fold encounters malformed frontmatter | `warning: skipping malformed closeout (no generated-at): <name>` (proceeds with valid ones) | Hand-edit the malformed file or delete it |
| Doctor advisory: duplicate EV-### detected | Advisory line listing the dupe IDs and files. Non-fatal. | Run an LLM cleanup pass (agent renumbers + rewrites references) |
| Doctor advisory: dated source has rows missing `occurred:` | Advisory line listing the source short-name and row IDs. Non-fatal. | Hand-edit or re-run ingest with corrected source-map dating |

## Interface

### CLI

```
ctx handover write <title> --summary "..." --next "..." [--highlights "..."] [--open-questions "..."] [--no-fold] [--commit <sha>]
ctx kb ingest <folder|paths...>            # 4-phase ingest pass
ctx kb ask "<question>"                    # Q&A from existing KB; read-only on prose
ctx kb site-review                         # mechanical structural audit
ctx kb ground                              # external grounding via grounding-sources.md
ctx kb note "..."                          # lightweight capture; appends to ingest/findings.md
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

### Skill

```
/ctx-handover <title>           # deploys via ctx setup; wraps ctx handover write
/ctx-kb-ingest <folder|paths>   # wraps ctx kb ingest
/ctx-kb-ask "<question>"        # wraps ctx kb ask
/ctx-kb-site-review             # wraps ctx kb site-review
/ctx-kb-ground                  # wraps ctx kb ground
/ctx-kb-note "..."              # wraps ctx kb note
```

Trigger phrases (per-skill):

- `/ctx-handover`: "let's wrap up", "save state", "leave a handover", "before I go", "stepping away" (also: mandatory tail of `/ctx-wrap-up`).
- `/ctx-kb-ingest`: "ingest the transcripts", "pull this into the kb", "add evidence from", explicit slash form for paths.
- `/ctx-kb-ask`: "does the kb say", "according to evidence", explicit slash form for the question.
- `/ctx-kb-site-review`: "audit the kb", "check kb for rot".
- `/ctx-kb-ground`: "re-ground the kb", "check upstream", explicit slash form preferred.
- `/ctx-kb-note`: "drop a note", "capture this for the next ingest", "park this finding".

`/ctx-wrap-up` skill is updated to: (a) detect `.context/kb/`
existence, (b) surface editorial state (pending closeouts,
outstanding-questions count) in its summary, (c) mandatorily
drive `/ctx-handover` as the final step regardless of capture
outcomes.

## Implementation

### Files to Create / Modify

| File | Change |
|------|--------|
| `internal/cli/handover/cmd.go` | New: `ctx handover write` Cobra command. `MarkFlagRequired` on `--summary` and `--next`. Calls `internal/store/handover.go` writer + closeout fold helpers. |
| `internal/cli/kb/cmd.go` | New: `ctx kb` Cobra parent command. |
| `internal/cli/kb/ingest/cmd.go` | New: `ctx kb ingest` four-phase pipeline runner. Refuses on empty input. |
| `internal/cli/kb/ask/cmd.go` | New: `ctx kb ask` Q&A driver. Refuses on empty question. |
| `internal/cli/kb/sitereview/cmd.go` | New: `ctx kb site-review` mechanical audit. |
| `internal/cli/kb/ground/cmd.go` | New: `ctx kb ground` external grounding. Reads `grounding-sources.md`; prompts when empty. |
| `internal/cli/kb/note/cmd.go` | New: `ctx kb note` lightweight capture. Appends to `.context/ingest/findings.md`. |
| `internal/cli/kb/site/cmd.go` | New: `ctx kb site build|serve|customize` mirrors existing `ctx journal site`. |
| `internal/store/handover.go` | New: `WriteHandover`, `LatestHandoverCursor`, `UnconsumedCloseouts`, `ArchiveCloseouts`. Mirrors sibling shape. |
| `internal/store/closeout.go` | New: `WriteCloseout` with required frontmatter; reader for `generated-at` cursor extraction. |
| `internal/store/kb.go` | New: writers for `evidence-index.md` (append, never renumber), `glossary.md`, `contradictions.md`, `outstanding-questions.md`, `domain-decisions.md`, `timeline.md`, `source-map.md`, `relationship-map.md`. |
| `internal/path/path.go` | Extend: `HandoversDir`, `KBDir`, `KBEvidenceFile`, `KBGlossaryFile`, ... `IngestDir`, `IngestRulesFile`, `IngestInboxFile`, `IngestSessionLogFile`, `IngestGroundingSources`, `CloseoutsSubdir`, `ArchiveCloseoutsSubdir`, `SiteDir`, `SiteKBDir`, `SiteConfigDir`. |
| `internal/cli/initcmd/init.go` | Extend: lay down `handovers/`, `kb/.gitkeep`, `ingest/` (full template tree), `site/` (gitignored). Add `--upgrade` flag for repos initialized before this spec. |
| `internal/assets/kb/templates/ingest/*.md` | New embedded templates: `KB-RULES.md`, `00-GROUND.md`, `30-INGEST.md`, `40-ASK.md`, `50-SITE_REVIEW.md`, `INBOX.md`, `SESSION_LOG.md`, `grounding-sources.md`, `OPERATOR.md`, `PROMPT.md`. |
| `internal/assets/kb/templates/ingest/schemas/*.md` | New embedded schemas: `evidence-index.md`, `glossary.md`, `contradictions.md`, `outstanding-questions.md`, `domain-decisions.md`, `timeline.md`, `source-map.md`, `relationship-map.md`, `session-log.md`. Each carries fields list + one worked example, no domain content. |
| `internal/assets/claude/skills/ctx-handover/SKILL.md` | New skill file with input contract, authority boundary, edge cases per spec. |
| `internal/assets/claude/skills/ctx-kb-ingest/SKILL.md` | New: 4-phase pipeline driver; refuses on empty input. |
| `internal/assets/claude/skills/ctx-kb-ask/SKILL.md` | New: Q&A from KB; refuses on empty question. |
| `internal/assets/claude/skills/ctx-kb-site-review/SKILL.md` | New: mechanical audit; no arguments. |
| `internal/assets/claude/skills/ctx-kb-ground/SKILL.md` | New: external grounding; reads `grounding-sources.md`. |
| `internal/assets/claude/skills/ctx-kb-note/SKILL.md` | New: lightweight capture. |
| `internal/assets/claude/skills/ctx-wrap-up/SKILL.md` | Modify: branch on `.context/kb/` existence; mandatorily drive `/ctx-handover` as final step. |
| `internal/assets/claude/skills/ctx-remember/SKILL.md` (or wherever it is) | Modify: read latest handover + any postdated unfolded closeouts; fold KB state into readback if `.context/kb/` exists. |
| `internal/cli/doctor/advisory.go` | Extend: duplicate-`EV-###` detection, `dated:` sources missing `occurred:` check, malformed-closeout-frontmatter detection. |
| `internal/cli/setup/cmd.go` | Already walks skills dir; new skill subdirs picked up automatically. |
| `internal/cli/wrapup/run.go` | Modify: `if KBExists() { printPendingCloseouts; printOutstandingQuestionsCount }`. Mandatory handover wording remains in skill, not CLI. |
| Project root `.gitignore` | Append: `.context/site/` (idempotent — match existing pattern for `.context/journal/.imported.json`). |

### Key Functions

```go
// internal/store/handover.go
type HandoverEntry struct {
    Title           string
    Summary         string
    Next            string
    Highlights      string
    OpenQuestions   string
    Commit          string
    Branch          string
    FoldedCloseouts []CloseoutFile
}
func (s *Store) WriteHandover(e HandoverEntry) (HandoverResult, error)
func (s *Store) LatestHandoverCursor() (cursor time.Time, file string, err error)
func (s *Store) UnconsumedCloseouts(cursor time.Time) (consumed []CloseoutFile, malformed []string, err error)
func (s *Store) ArchiveCloseouts(files []CloseoutFile) error

// internal/store/closeout.go
type CloseoutFrontmatter struct {
    SHA, Branch, Mode string
    GeneratedAt       time.Time
}
type CloseoutFile struct {
    Path string
    Frontmatter CloseoutFrontmatter
    Body string
}
func (s *Store) WriteCloseout(mode string, body string) (CloseoutFile, error)

// internal/store/kb.go
func (s *Store) AppendEvidence(row EvidenceRow) error  // never renumber
func (s *Store) DemoteEvidence(id string, newBand string, reason string) error
// ... per-artifact writers
```

### Helpers to Reuse

- `internal/gitmeta/` — already resolves git HEAD into commit + branch with overrides for CI replay (`CTX_TASK_COMMIT`, `GITHUB_SHA`).
- `internal/cli/journal/site/` — existing zensical shell-out pattern; lift the build/serve runtime-config materialization wholesale for `ctx kb site`.
- `internal/cli/initcmd/` — existing template walk for embedded asset deployment.
- `internal/cli/setup/` — existing skill-dir walker; new skill subdirs land for free.
- `internal/cli/doctor/advisory.go` — existing advisory pattern; add new check functions following the same shape.
- `internal/path/path.go` — existing path constant convention; extend.
- `internal/store/` — existing store helpers (path resolution, `O_CREATE|O_EXCL` writes, idempotency); reuse for new writers.

## Configuration

No new `.ctxrc` keys in v1. Site-config (`.context/site-config/kb.toml`)
is lazy-initialized via `ctx kb site customize` for users who want
to override theme / nav / plugins. Infrastructure paths
(`docs_dir`, `site_dir`) are wrapper-owned and overwritten at
build time per the existing `ctx journal site` pattern.

Environment variables honored:

- `CTX_TASK_COMMIT` — override resolved commit for handover Provenance (CI replay).
- `GITHUB_SHA` (when `GITHUB_ACTIONS=true`) — same purpose.

`ZENSICAL_BIN` is **not** introduced; the binary is resolved
from PATH per the existing journal-site convention.

## Testing

### Unit

- `internal/store/handover_test.go` — `WriteHandover` happy
  path; rejects empty summary/next; rejects placeholder bodies;
  closeout fold cursor logic; archive moves files
  atomically.
- `internal/store/closeout_test.go` — frontmatter parse;
  malformed frontmatter handling; `generated-at` ordering.
- `internal/store/kb_test.go` — evidence-index append never
  renumbers; demotion bands valid only within
  `high|medium|low|speculative`; concurrent-write detection
  surfaces in doctor.
- `internal/cli/kb/<mode>/cmd_test.go` — refuse-on-empty for
  ingest / ask / ground.
- `internal/cli/handover/cmd_test.go` — `MarkFlagRequired`
  enforcement; `--no-fold` behavior; provenance from
  `gitmeta.ResolveHead`.
- `internal/cli/doctor/advisory_test.go` — duplicate-`EV-###`
  detection; `dated:`-without-`occurred:` detection;
  malformed-closeout detection.

### Integration

- `internal/cli/initcmd/init_test.go` — full init creates all
  new dirs and templates; `--upgrade` is idempotent on
  byte-identical existing content; `--upgrade` refuses on
  divergent existing content.
- `internal/cli/setup/setup_test.go` — new skill subdirs deploy.
- `hack/smoke-kb.sh` — end-to-end shell smoke: `ctx init`;
  `ctx kb ingest ./testdata/inputs`; `ctx kb ask "..."`;
  `ctx kb site-review`; `ctx kb ground`; `ctx handover write
  --summary X --next Y`; verify files exist, closeouts folded,
  archive populated, no doctor errors.

### Edge cases

- Aborted-session recovery: write a closeout, do NOT write a
  handover, simulate session restart by re-reading the state;
  verify `/ctx-remember`'s read path picks up the unfolded
  closeout (test the store helper directly).
- Temporal misordering: ingest fixture A (occurred 2026-04-12,
  extracted 2026-05-09) then fixture B (occurred 2026-04-05,
  extracted 2026-05-09); verify demotion does NOT fire because
  temporal-precedence rule wins.
- Concurrent dupe IDs: simulate two parallel writers producing
  `EV-020` against different claims; verify `ctx doctor` flags
  the dupe; verify cross-references survive an LLM-style manual
  resolution (provided as a fixture in tests).
- Render filter: speculative content does NOT appear in built
  HTML; low-confidence content appears only when paired with
  outstanding-questions row.

### Validation corpus

`things-wtf-disaster-recovery` is the live regression suite.
Phase 2 (per the brief) is "port things-wtf to the shipped
shape and document any divergence." Each divergence is either a
bug fix on this spec or a `DECISIONS.md` entry explaining why
the formal shape differs from what worked manually.

## Non-Goals

(Explicit deferrals — referenced from the brief's "What we
rejected" table.)

- **No `/ctx-kb-decide` skill.** KB ontology rejects it: in a
  KB you don't decide, you increase confidence. Pipeline is the
  sole writer; ad-hoc capture flows through `ctx kb note` or
  hand-edit.
- **No team write-coordination layer.** Single-writer convention
  with doctor advisory + LLM cleanup is the v1 stance.
  Multi-writer coordination is deferred until a real team-scale
  user hits the wall.
- **No UUIDs for evidence rows.** `EV-###` aesthetic preserved.
- **No KB-scoped IDs (`research-master/EV-019`) in v1.** P2
  federation handles multi-KB without scoping IDs.
- **No domain-split per user (`.context/kb/<domain>/`).** Single
  KB is simpler; multi-domain via P2 federation when needed.
- **No KB-side merge into canonical files.** `domain-decisions.md`
  stays separate from `DECISIONS.md`. Different schema, different
  write authority, different lifecycle.
- **No automatic demotion-cascade in v1.** When EV-031 demotes,
  affected glossary / domain-decision / timeline rows are NOT
  auto-flagged; the human handles the cascade. v2 may add
  site-review automation.
- **No bundled renderer.** `zensical` is shelled out, not
  vendored or wrapped. Same model as `ctx journal site`.
- **No KB linting.** Confidence-band discipline is rule-driven
  (per `KB-RULES.md`), not enforced programmatically.
- **No bulk migration tooling for repos initialized before this
  spec.** `ctx init --upgrade` lays down the new dirs
  idempotently; pre-existing canonical files are untouched;
  hand-rolled editorial files (e.g. things-wtf's
  `10-CONSTITUTION.md` at repo root) are left alone — porting
  is a manual cutover (Phase 2 of validation).
- **No interactive install of zensical.** Missing-binary case
  fails with a one-line install hint and non-zero exit.
- **No replacement of `/ctx-decision-add`, `/ctx-learning-add`,
  `/ctx-task-add`, `/ctx-convention-add`, `/ctx-wrap-up`.** New
  skills are siblings; existing capture skills unchanged in
  authority.

## Open Questions

(Carried forward from the brief; require pinning during
implementation.)

1. **Naming.** `ctx kb ingest|ask|site-review|ground`
   (kb-prefixed) vs `ctx ingest|ask|site-review|ground`
   (top-level). Lean prefixed. Confirm during implementation
   kickoff.
2. **Brief vs spec storage.** Where do briefs live?
   `.context/briefs/` (debate residue, distinct lifecycle) vs
   `.context/specs/briefs/` (subdir of specs). Lean dedicated
   `briefs/`. Tied to the polish-PR work in `ideas/002` §3.
3. **`ctx kb note` destination.** Single
   `.context/ingest/findings.md` or one file per invocation?
   Lean single file (simpler); per-invocation preserves
   provenance per note.
4. **Pure-research project init.** Should `ctx init` learn
   `--research` / `--kb-only` to suppress code-dev skill
   deployment? Light defer; not v1 critical. Things-wtf
   workaround (CLAUDE.md disabling) remains acceptable for v1.
5. **Confidence bands flowing into `LEARNINGS.md`.** Probably
   no (different truth bases — KB has citations, learnings have
   author intent). Confirm rather than assume.
6. **`relationship-map.md` vs GitNexus.** Different graphs;
   v1 keeps independent. Cross-feed is v2.
7. **Demotion-policy automation.** Auto-flag affected pages on
   demotion? v1 defers to human; v2 may automate via
   site-review.
8. **`--no-fold` flag scope.** Handover-only (sibling's choice)
   vs every artifact-writing command. Lean handover-only.
9. **Polish-PR ordering.** The `MarkFlagRequired` /
   `--brief <path>` / authority-boundary rewrites in
   `ideas/002` §3 should ship **before** this spec is
   implemented (so `/ctx-spec --brief` works for the next round
   of feature specs). Recommended as Phase 0a of the
   implementation task breakdown.
10. **Git-mandate dependency (Phase 0b).** `specs/require-git.md`
    is a hard prerequisite. Once it ships, the `commit:none`
    sentinel is unreachable: the duplicate-`EV-###` advisory in
    `internal/cli/doctor/advisory.go` does NOT need to handle
    `commit:none` cases; `gitmeta.ResolveHead` calls in
    `WriteHandover` / `WriteCloseout` similarly drop their
    `none` fallback paths; closeout frontmatter `sha:` /
    `branch:` are guaranteed populated. Confirm during
    implementation that no Phase KB code paths silently retain
    `commit:none` handling.

---

## Source

`ideas/003-editorial-pipeline-debated-brief.md` (debated brief
from `/ctx-plan` session 01d0cf92, 2026-05-09).

Inputs to the brief:

- `ideas/001-sibling-project-undercover-analysis.md` — handover
  mechanism, closeout/fold mechanism, doctor advisory tier.
- `ideas/002-editorial-pipeline-and-skill-rigor.md` — the lifted-
  pipeline plan + skill ceremony comparison.
- `things-wtf-disaster-recovery/` (sibling workspace) — live
  test corpus; hand-rolled version of the shape.
- Sibling tool's `internal/cli/initcmd/templates/ingest/` —
  source templates referenced for shape lift.
- `.context/journal-site/zensical.toml` — proof that ctx already
  shells out to zensical for journal rendering.
