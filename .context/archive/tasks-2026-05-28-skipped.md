# Skipped Tasks - Archived 2026-05-28

Tasks deliberately skipped with documented reasons. Preserved for traceability.

## Agents

- [-] Add `ctx explore` command — scaffolds `.arch-explorer/` in a workspace
  directory with manifest.json, PROMPT.md (from
  `hack/agents/architecture-explorer.md`), run-log.md, and a README. Similar to
  `ctx init` but for multi-repo architecture exploration. The prompt template
  lives in `hack/agents/architecture-explorer.md` and ships embedded.
  #priority:low #added:2026-04-13
  **Skipped 2026-04-16**: Superseded by
  `docs/operations/runbooks/architecture-exploration.md`. A runbook is the right
  weight — a CLI scaffolding command was speculative abstraction for a
  workflow
  that's better served by a discoverable doc with an embedded prompt.

## Misc

### Architecture Docs

- [-] Publish architecture docs to docs/: copy ARCHITECTURE.md,
  DETAILED_DESIGN domain files, and CHEAT-SHEETS.md to docs/reference/.
  Sanitize intervention points into docs/contributing/.
  Exclude DANGER-ZONES.md and ARCHITECTURE-PRINCIPAL.md (internal only).
  Spec: specs/publish-architecture-docs.md #priority:medium
  #added:2026-04-03-150000 #skipped:2026-05-23
  Decided not to ship. Reasons: (1) audience is maintainer-focused —
  anyone wanting DETAILED_DESIGN depth can read it on GitHub from the
  canonical `.context/` source; (2) AI-generated content would require
  a permanent human editorial pass before each publish; (3) every
  architecture change forces a re-run + re-publish loop or accepts
  known staleness in the public docs site; (4) the marginal
  discoverability gain doesn't justify importing that maintenance
  burden into the docs pipeline. If discoverability ever becomes a
  real ask, cheap fallback is a one-page
  `docs/contributing/architecture.md` that links to the GitHub-hosted
  `.context/ARCHITECTURE.md` — pointer, not a copy. Replay note: do
  not re-open without revisiting these four reasons.

### Phase -3: DevEx

- [-] Create ctx-docstrings skill: audit and fix docstrings
  against CONVENTIONS.md Documentation section. Superseded by
  TestDocCommentStructure compliance test (68 grandfathered).
  #added:2026-03-20-163413
  #added:2026-03-16-114445

### Phase -2: Task completion nudge:

- [-] ctx-architecture-extend
  Skipped: extension point analysis is covered by /ctx-architecture
  DETAILED_DESIGN (per-module) and /ctx-architecture-enrich
  (registration sites). A fourth skill fragments the pipeline
  without enough distinct value. Three is the right number:
  map, enrich, hunt.
  **Context**: Companion to `ctx-architecture` and
  `ctx-failure-analysis`, completing a trilogy: how does it
  work → where will it break → where does it grow. Reads
  architecture artifacts → identifies registration patterns
  (interfaces, factory functions, plugin systems, ordered
  slices, scheme registrations) → traces recent additions via
  git log to confirm which extension points are actually used
  → produces `EXTENSION-POINTS.md` ranked by frequency, with
  exact file locations, function signatures, and the typical
  feature pattern (e.g., "most features require a variable +
  a mutator + a machine-agent task").

  Valuable for onboarding ("I need to add feature X, where do
  I start?") and architecture review ("are we adding features
  in the right places?").

    - [-] Design SKILL.md for ctx-extension-map
      Skipped: parent task skipped.
      #added:2026-03-25-062000
    - [-] Define extension point detection heuristics
      Skipped: parent task skipped.
      #added:2026-03-25-062000
    - [-] Add git log frequency analysis
      Skipped: parent task skipped.
      #added:2026-03-25-062000
    - [-] Integrate with GitNexus for registration sites
      Skipped: parent task skipped.
      #added:2026-03-25-062000

### Phase 0.9: Suppress Nudges After Wrap-Up

- [-] P0.9.2: Split cli-reference.md — moved to Future
  #added:2026-02-24-204208

- [-] P0.9.3: Investigate proactive content suggestions — moved to Future
  #added:2026-02-24-185754

### Phase 0.5 Cleanup

- [-] Move generic string helpers from cli/add/core/strings.go to
  internal/format — file no longer exists, helpers already moved or deleted
  #added:2026-03-20-175046

### Phase MI: Memory Import Pipeline (`ctx memory import`)

- [-] MI.future: `--interactive` mode for agent-assisted classification —
  skipped: `--dry-run` covers review; agents can use `ctx add` directly for
  overrides; interactive CLI prompts don't compose with agent workflows

### Phase ET: Error Package Taxonomy (`internal/err/`)

- [-] Cleanup internal/cli/system/core/persistence.go: move 10 (base for
  ParseInt) to config constant — not actionable, 10 is stdlib decimal base
  convention, not a magic number #priority:low #added:2026-03-07-220825

- [-] Cleanup internal/cli/system/core/session_tokens.go: move SessionStats from
  state.go to types.go — file and type no longer exist, refactored away
  #priority:low #added:2026-03-07-220825

- [-] SMB mount path support: add `CTX_BACKUP_MOUNT_PATH` env var so
  `ctx backup` can use fstab/systemd automounts instead of requiring GVFS.
  Spec: specs/smb-mount-path-support.md #priority:medium
  #added:2026-04-04-010000
  **Skipped 2026-04-16**: Duplicate of line 214. Superseded by
  specs/deprecate-ctx-backup.md (full removal, not mount path fix).

- [-] Refactor check_backup_age/run.go: move consts (lines 23-24) to config,
  magic directories (line 59) to config, symbolic constants for strings (line
  72), messages to assets (lines 79, 90-91), extract non-Run functions to
  system/core, fix docstrings #priority:medium #added:2026-03-07-180020
  **Skipped 2026-04-16**: Superseded by specs/deprecate-ctx-backup.md
  (check_backup_age will be removed entirely, not refactored).

## Future

### Phase BT: Build Tooling — `cmd/ctxctl`

- [-] Explore pluggable graph tool interface — replace hardcoded GitNexus
  references in skill text with configurable .ctxrc graph_tool key. Skills use
  template placeholder instead of literal tool names. Define minimum interface
  contract (query, context, impact). Spec:
  `ideas/spec-mcp-warm-up-ceremony.md` #added:2026-03-25-120000
  **Skipped 2026-05-23**: contradicts the committed-to-GitNexus
  direction recorded in DECISIONS.md "MCP gateway not worth the
  coupling cost". Pluggable abstraction implies multiple
  candidate graph tools, which in turn implies ctx vouches for
  the interface contract across implementations — exactly the
  ownership coupling we're avoiding. If a second viable graph
  tool emerges that's worth the cost of pluggability, revisit
  by un-skipping; the design sketch in
  `ideas/spec-mcp-warm-up-ceremony.md` stays available as the
  starting point.

### Phase: ctx Hub follow-ups (PR #60)

- [-] `docs/thesis/index.md:412` (the primitive
  comparison table saying "Document: Zero-dependency:
  Yes"): left intact. The claim is about the document
  primitive itself (markdown files have no runtime
  deps), not about ctx as an implementation. Accurate.
  #added:2026-04-11 #skipped:primitive-claim-is-correct

### Later

- [-] PROMPT.md design — belongs in another project; skipped here.
  #session:4b37e2f6 #added:2026-04-14-010311 #skipped:2026-04-14

### Phase RG: Require Git as Architectural Precondition (Phase 0b; prerequisite for Phase KB)

- [-] Update `ctx init` to call `RequireGitTree` first. N/A: `ctx init` is
  `AnnotationSkipInit`; the precondition
  check moved to the init command body in Phase KB Stage 5 alongside `--upgrade`

- [-] Remove `commit:none` fallback from `internal/gitmeta/resolvehead.go`. N/A:
  `resolvehead.go` is
  net-new in this phase; no fallback to remove

- [-] Remove `commit:none` advisory + counts from
  `internal/cli/doctor/advisory.go`. Verified by grep:
  no `commit:none` / `commit=\"none\"` literal exists in `internal/` already

- [-] Audit `internal/cli/<various>/cmd.go` for any other `commit:none` fallback
  handling; remove. Same:
  audit by grep returned zero matches

- [-] Tests: root PreRunE refuses without git; opt-out list allowed. TBD:
  deferred to Phase KB Stage 5 when
  the kb command tree is in place (the existing bootstrap_test.go covers PreRunE
  structure; the gitmeta
  injection's behavioral test runs as part of the kb-ingest smoke)

- [-] Compliance test: no remaining `commit:none` literal in `internal/`. N/A:
  literal never existed

### Phase KB: Editorial Pipeline + Handover (depends on Phase SK + Phase RG)

- [-] Doctor advisories: NOT YET IMPLEMENTED. Spec lists duplicate-`EV-###`,
  `dated:`-source-missing-`occurred:`, malformed-closeout-frontmatter,
  source-coverage-ledger-mismatch (row Updated vs. file mtime),
  closeout-missing-pass-mode-body-block, illegal-ledger-state-transition. Phase
  7 follow-up.

- [-] Mode-aware reads in ctx status / ctx agent / session-start hook: skills
  updated (`/ctx-remember` + `/ctx-wrap-up`); CLI-side `ctx status`/`ctx agent`
  mode-awareness deferred (the skill-side fold covers the user-facing recall;
  CLI text surfaces are v1.1).
