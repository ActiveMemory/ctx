# Tasks

<!--
UPDATE WHEN:
- New work is identified → add task with #added timestamp
- Starting work → add #in-progress or #started timestamp
- Work completes → mark [x]
- Work is blocked → add to Blocked section with reason
- Scope changes → update task description inline

DO NOT UPDATE FOR:
- Reorganizing or moving tasks (violates CONSTITUTION)
- Removing completed tasks (use ctx task archive instead)

STRUCTURE RULES (see CONSTITUTION.md):
- Tasks stay in their Phase section permanently: never move them
- Use inline labels: #in-progress, #blocked, #priority:high
- Mark completed: [x], skipped: [-] (with reason)
- Never delete tasks, never remove Phase headers

TASK STATUS LABELS:
  `[ ]`: pending
  `[x]`: completed
  `[-]`: skipped (with reason)
  `#in-progress`: currently being worked on (add inline, don't move task)
-->

## Phase 0 Grounding

- [x] Add TypeScript type-check step (`tsc --noEmit`) for embedded
  editor-plugin assets to CI; nothing currently
  checks `.opencode/plugins/ctx/index.ts` before embedding
  #priority:low #added:2026-04-26-152912 #completed:2026-05-11
  Implementation: `tools/typecheck/opencode/` (package.json,
  tsconfig.json, README, package-lock.json); CI job
  `typecheck-opencode-plugin` in `.github/workflows/ci.yml`; embed
  contract documented in `internal/assets/README.md` and decision
  recorded 2026-05-11. Investigation also surfaced three
  related-but-distinct gap tasks (shellcheck, PSScriptAnalyzer,
  skill frontmatter validation) listed below.

- [ ] Human: It's about time to go through the entire codebase check for
  inconsistencies, and move useful functions that are utility and/or reusable
  to relevant convenience packages.

- [ ] Create a typography.md somewhere so that we don't have to remind tha
  Agent things like this: "❯"## What the editorial pipeline is NOT" our headings 
  are Title Case, it always has been; it always will be. Do a full sweep. 
  -- in addition (not checked, just to make sure); `ctx` is always in backticks 
  whenever possible; it's part of the branding."


- [ ] Bug: Fresh folder: git init; eval $(ctx activate); ctx init
  will catch parent .context folder and raise a warning
  expectation: ctx activate should bail if there is no .context
  folder in the same level and ask user to run `ctx init` first.
  discuss this with the Agent too.
- [ ] Bug: if context is active (eval ctx activate); `ctx init`
  on a brand new project can (and probably will) fail. 
  Probably need to nudge user to ctx deactivate first.

- [-] Add TypeScript `tsc --noEmit` gate for the embedded OpenCode
  plugin (`internal/assets/integrations/opencode/plugin/index.ts`).
  Place tooling (`package.json`, `tsconfig.json`) in a sibling
  directory outside `internal/assets/` so it does not pollute the
  embed surface; wire a CI step that installs Bun, `bun install`,
  then `bunx tsc --noEmit`. Spec: respawn from
  `specs/internal-assets-readme.md` open work.
  #priority:low #added:2026-05-11 #grounding-gap
  #skipped:2026-05-11 reason: duplicate of original line-30 task
  above, which has now been completed end-to-end. CI uses `npm ci`
  + `npx tsc --noEmit` (matching the existing `editors/vscode/`
  convention) rather than Bun; `tsc` is the same compiler either
  way and `@types/bun` provides Bun globals to the type-checker.

- [ ] Add `shellcheck` gate for embedded shell scripts
  (`internal/assets/integrations/copilot-cli/scripts/*.sh` and
  `internal/assets/hooks/trace/*.sh`). Run in CI; fail on findings
  at severity `warning` and above. #priority:low #added:2026-05-11
  #grounding-gap

- [ ] Add `PSScriptAnalyzer` gate for embedded PowerShell scripts
  (`internal/assets/integrations/copilot-cli/scripts/*.ps1`). Run in
  CI on a Windows or pwsh-enabled runner; fail on findings at
  severity `Warning` and above. #priority:low #added:2026-05-11
  #grounding-gap

- [ ] Add skill frontmatter validity test covering every embedded
  `SKILL.md` (Claude skills, OpenCode skills, Copilot CLI skills):
  assert required keys present and values typed correctly. Extend
  `internal/assets/embed_test.go` or add a dedicated test under
  `internal/assets/read/skill/`. #priority:medium #added:2026-05-11
  #grounding-gap

- [x] Add CI guardrails for the VS Code extension at
  `editors/vscode/` (separately-published deliverable, ships via
  VS Code Marketplace under publisher `activememory`, not embedded
  into the ctx binary). CI job `vscode-extension` runs `npm ci`,
  `npm run build` (esbuild bundle), and `npx tsc --noEmit
  -p tsconfig.ci.json` (production code only; test file excluded
  pending separate fix). README docs added: `internal/assets/README.md`
  gained an "Embedded vs. Separately-Published" comparison; the
  extension's own README gained a "Release" section documenting
  the manual `vsce publish` flow and the CI gates that protect it.
  #priority:medium #added:2026-05-11 #completed:2026-05-11
  #grounding-gap

- [x] Close VS Code documentation parity gap: ctx had a dedicated
  `docs/home/opencode.md` (185 lines) and `site/home/opencode/`
  published page, but no equivalent for VS Code — the extension
  was reduced to a 24-line install snippet inside
  `docs/operations/integrations.md` plus the marketplace README.
  A docs-site reader had no path to day-to-day usage. Created
  `docs/home/vscode.md` mirroring the opencode page shape
  (problem, setup, what gets created, automatic hooks,
  status bar, slash commands by category, natural language,
  auto-bootstrap, prerequisites, configuration, troubleshooting,
  verification, what's next). Registered in `zensical.toml`'s
  "Get Started" nav. Expanded the integrations.md VS Code
  subsection to point at the new home page and added a
  parallel "First-Class Citizen" block in
  `docs/recipes/multi-tool-setup.md`. #priority:medium
  #added:2026-05-11 #completed:2026-05-11 #grounding-gap

- [ ] Fix `editors/vscode/src/extension.test.ts` type errors and
  re-enable test-file type-checking + vitest in CI. Two distinct
  bugs: (1) tests import handlers (`handleComplete`, `handleTasks`,
  `handleRemind`, `handlePad`, `handleNotify`, `handleSystem`,
  `handleSpec`) that are no longer exported from `extension.ts`
  (only `activate` and `deactivate` are exported now) — the test
  suite is rotting against the actual extension surface; (2) the
  `fakeToken` helper's `onCancellationRequested` mock signature
  is `(cb: () => void) => …` but the VS Code API expects
  `(e: any) => any` with at least one argument. Once fixed,
  remove the `tsconfig.ci.json` carve-out and add `npm test` to
  the `vscode-extension` CI job. Also worth adding `npm run lint`
  (eslint) and a `vsce package` dry-run step.
  #priority:medium #added:2026-05-11 #grounding-gap

- [ ] The target project (to be given to the Agent) has a good "phasing"
  mechanism for tasks; implement that; maybe `ctx task add` can have a
  `--phase` flag too, and we can have a auditor/normalizer for the current
  task document; or a skill that does a semantic pass, or both too.

- [ ] Localize the placeholder set used by `RejectPlaceholder`
  (decision add / learning add and any future body-flag validators).
  Move the shipped defaults out of `internal/config/validate/placeholder.go`
  Go constants into an embedded YAML asset, add a `.ctxrc placeholders:`
  override with EXTEND semantics (user list is appended to defaults, not
  replacing them — Tarzan Turkish is the dominant case), and replace the
  current `strings.ToLower` with proper Unicode case folding via
  `golang.org/x/text/cases` so İ/i, ß/SS, and similar fold correctly.
  Ship `en` only in v1; ctx has no locale-specific assets yet, so the
  structure is established but no `tr.yaml` lands in this work.
  Spec: `specs/placeholder-i18n.md` #priority:high #added:2026-05-11
  #prerequisite-for-locale-work

### Misc

### Agents

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

### Runbooks

### Misc

- [ ] Human: Read the entire documentation page-by-page, line-by-line, with a
  critical mind, including blog posts. Take notes for agent to rectify, or
  directly update the docs whenever it makes sense.

- [ ] Human: Do a documentation audit for AI-generated artifacts. #important
  #not-urgent

- [ ] Human: test `ctx init` on a fresh ubuntu install.

- [ ] Improve hub failover client: distinguish auth errors
  (Unauthenticated/PermissionDenied) from connection errors. Fail fast on auth
  failures instead of cycling through all peers with the same invalid token.
  #priority:low #added:2026-04-08-194612

- [ ] Add file locking to ctx connect sync state to prevent concurrent sync
  races. Two sync processes (hook + manual) can both load the same LastSequence,
  process the same entries, and write duplicate content to .context/shared/.
  #priority:medium #added:2026-04-08-194557

- [ ] Fix fanout broadcast entry loss: non-blocking send drops entries to slow
  listeners silently. Log when entries are dropped. Consider per-listener
  backpressure or disconnect-on-lag. Buffer of 64 is too small for busy hubs.
  #priority:medium #added:2026-04-08-194542

- [ ] Prevent duplicate client registration in hub store: RegisterClient should
  reject if ProjectName already exists. Add token revocation support (delete
  client by ID/project). Currently tokens are valid forever with no way to
  disable compromised ones. #priority:medium #added:2026-04-08-194529

- [ ] Fix hub cluster: NewCluster result is discarded (not stored on Server), so
  Raft runs but leadership status is never queryable. Store cluster reference on
  Server, wire IsLeader/LeaderAddr into Status RPC and hub status command.
  #priority:medium #added:2026-04-08-194511

- [ ] Use crypto/subtle.ConstantTimeCompare for hub token validation instead of
  string equality. Current Store.ValidateToken uses == which is vulnerable to
  timing attacks. Also replace O(n) linear scan with a map[string]*ClientInfo
  for O(1) lookup. #priority:high #added:2026-04-08-194458

- [ ] Fix silent error suppression in hub: (1) ctx add --share silently ignores
  publish failures — warn user on failure, (2) hubsync hook swallows all
  errors — log to event system, (3) replication loop drops errors silently —
  add structured logging for debug. #priority:high #added:2026-04-08-194443

- [ ] Add input validation to hub Publish handler: reject empty ID, validate
  Type against allowed set (decision/learning/convention/task), enforce Content
  length limit (1MB), require non-empty Origin. Prevents garbage data and DoS
  via unbounded content. #priority:high #added:2026-04-08-194430

- [ ] Fix ctx connect listen: currently only does initial sync then blocks on
  ctx.Done() without ever calling the Listen RPC. Must stream entries in
  real-time via the server-streaming Listen RPC, writing to .context/shared/ as
  entries arrive. #priority:high #added:2026-04-08-194415


- [ ] Deprecate and remove `ctx backup`: hub handles cross-machine persistence,
  backup is environment-specific (SMB/GVFS/rsync), and it is the wrong layer
  for ctx to own. Replace with a backup-strategy runbook. About 60 files to
  remove across CLI, config, hooks, docs, skills. Implementation order: runbook
  first, then hook removal, then command removal, then docs cleanup.
  Spec: specs/deprecate-ctx-backup.md #priority:medium
  #added:2026-04-04-010000 #updated:2026-04-16

### Architecture Docs

- [ ] Publish architecture docs to docs/: copy ARCHITECTURE.md,
  DETAILED_DESIGN domain files, and CHEAT-SHEETS.md to docs/reference/.
  Sanitize intervention points into docs/contributing/.
  Exclude DANGER-ZONES.md and ARCHITECTURE-PRINCIPAL.md (internal only).
  Spec: specs/publish-architecture-docs.md #priority:medium
  #added:2026-04-03-150000

- [ ] Update ctx-architecture skill to append discovered terms to GLOSSARY.md
  during Phase 3. Additive only, max 10 terms per run, project-specific only,
  alphabetical insertion, skip if GLOSSARY.md empty. Print added terms in
  convergence report. Spec: specs/publish-architecture-docs.md #priority:low
  #added:2026-04-03-153000

### Code Cleanup Findings

- [ ] Implement journal compaction: Elastic-style tiered storage with tar.gz
  backup. Spec: specs/journal-compact.md #added:2026-03-31-110005

**PD.5 — Validate:**

### Phase -3: DevEx

- [-] Create ctx-docstrings skill: audit and fix docstrings
  against CONVENTIONS.md Documentation section. Superseded by
  TestDocCommentStructure compliance test (68 grandfathered).
  #added:2026-03-20-163413
  #added:2026-03-16-114445

### Phase -2: Task completion nudge:

- [ ] Design UserPromptSubmit hook that runs `make audit` at
  session start and surfaces failures as a consolidation-debt
  warning before the agent acts on stale assumptions.
  Project-level hook (not bundled in ctx), configurable via
  .ctxrc or settings.json. Related: consolidation nudge hook
  spec. #added:2026-03-23-223500

- [ ] Design UserPromptSubmit hook that runs go build and
  surfaces compilation errors before the agent acts on stale
  assumptions #added:2026-03-23-120136

- [ ] Architecture Mapping (Enrichment):
  **Context**: Skill that incrementally builds and maintains
  ARCHITECTURE.md and DETAILED_DESIGN.md. Coverage tracked in
  map-tracking.json. Spec: `specs/ctx-architecture.md`
    - [x] Create ctx-architecture-enrich skill: takes existing
      /ctx-architecture principal-mode artifacts as baseline, runs
      comprehensive enrichment pass via GitNexus MCP (blast radius
      verification, registration site discovery, execution flow
      tracing, domain clustering comparison, shallow module
      deep-dive). Spec: `ideas/spec-architecture-enrich.md`.
      Reference implementation: kubernetes-service enrichment pass
      2026-03-25. #added:2026-03-25-120000

- [ ]: ctx-architecture-failure-analysis
**Context**: Adversarial analysis skill that identifies where
a codebase will silently betray you. Requires
`ctx-architecture` artifacts as input (ARCHITECTURE.md,
DETAILED_DESIGN*.md, map-tracking.json). Does its own
targeted deep reads focusing on mutation points, shared
mutable state, error swallowing, concurrency, implicit
ordering, missing enforcement, and scaling cliffs. Uses
available tooling (GitNexus, Gemini Search) to
cross-reference patterns.

Produces `DANGER-ZONES.md` — a ranked inventory of silent
failure points with: location, failure mode, blast radius,
detection gap, and suggested fix. Two tiers: "most likely to
cause production incidents" and "less likely but equally
dangerous."

Distinct from a security threat model (which would be
`ctx-threat-model` — a separate skill for auth bypass,
injection, privilege escalation, supply chain). This skill
focuses on correctness: race conditions, ordering
assumptions, cache staleness, fan-out amplification,
non-atomic ownership, inverted logic, force-delete orphans,
global state mutation.

- [x] Design SKILL.md for ctx-architecture-failure-analysis:
  inputs (architecture artifacts), analysis phases, output
  format (DANGER-ZONES.md), quality checklist
  #added:2026-03-25-060000
- [x] Define the adversarial analysis framework: categories
  of silent failure (concurrency, ordering, cache,
  amplification, ownership, error swallowing, global state)
  with heuristics for each #added:2026-03-25-060000
- [x] Implement skill with GitNexus integration: use impact
  analysis for blast radius estimation, use context for
  shared-state detection #added:2026-03-25-060000
- [x] Add Gemini Search integration: cross-reference
  discovered patterns against known failure modes in similar
  systems. #added:2026-03-25-060000

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

### Phase CT: Companion Tool Integration

Session-start checks, suppressibility, and registry for companion MCP tools.

- [ ] ctx-remember preflight: verify ctx binary in PATH,
  plugin installed and enabled, binary version matches plugin
  version #priority:medium #added:2026-03-25-234514

- [ ] Design suppressible companion check system: .ctxrc
  configures which companion tools to check (one search MCP,
  one graph MCP), smoke tests only run for configured tools,
  not auto-discovered. Keeps bootstrap fast and predictable.
  #priority:medium #added:2026-03-25-234516

- [ ] Add per-tool suppression for ctx-remember checks: allow
  suppressing individual preflight checks (ctx binary, plugin,
  search MCP, graph MCP) via .ctxrc fields, not just
  companion_check: false blanket toggle
  #priority:low #added:2026-03-25-234518

### Phase CLI-FIX: CLI Infrastructure Fixes

### Phase BLOG: Blog Posts

- [ ] Write blog post about architecture analysis + enrichment two-pass design
  after dogfooding run on ctx itself. Cover: the 5.2x depth observation,
  constraint-as-feature principle, watermelon-rind anti-pattern, and results
  from the ctx self-analysis. #priority:medium #added:2026-03-25-233650

- [ ] Blog post: "Writing a CONSTITUTION for your AI agent" — showcase ctx's
  CONSTITUTION.md as a pattern for hard invariants that agents cannot violate.
  Cover: why advisory rules fail (agents game qualifiers), what belongs in a
  constitution vs conventions, the spec-at-commit enforcement story from this
  session, examples of good rules (absolute, binary, no interpretation needed).
  Include a recipe for writing your own.
  #priority:medium #added:2026-03-27-115500

- [ ] Recipe: "How to write a good CONSTITUTION.md" — practical guide with
  categories (security, quality, process, structure), anti-patterns (vague
  qualifiers, unenforced rules), enforcement mechanisms (hooks, commit gates),
  and a starter template. #priority:medium #added:2026-03-27-115500

- [ ] Import grouping compliance test: parse all .go files, verify imports
  follow stdlib — external — ctx three-group ordering. Add to
  internal/compliance/. Catches violations that goimports misses (it merges
  external and ctx into one group). #priority:medium #added:2026-03-27-120000

- [ ] drift check should notify if claude permissions have insecure stuff in it.

- [ ] task: sync workspace to ARI_INBOX

### Phase -1: Hack Script Absorption

Absorb remaining `hack/` scripts into Go subcommands. Eliminates shell
dependencies, improves portability, and makes the skill layer call `ctx`
directly instead of `make` targets.

### Phase 0.9: Suppress Nudges After Wrap-Up

Spec: `specs/suppress-nudges-after-wrap-up.md`. Read the spec before starting
any P0.9 task.

**Phase 3 — Skill integration:**

- [-] P0.9.2: Split cli-reference.md — moved to Future
  #added:2026-02-24-204208

- [-] P0.9.3: Investigate proactive content suggestions — moved to Future
  #added:2026-02-24-185754

### Phase 0.8: RSS/Atom Feed Generation (`ctx site feed`)

Spec: `specs/rss-feed.md`. Read the spec before starting any P0.8 task.

### Phase 0.4: Hook Message Templates

Spec: `specs/future-complete/hook-message-templates.md`. Read the spec before
starting any P0.4 task.

**Phase 2 — Discoverability + documentation:**

Spec: `specs/future-complete/hook-message-customization.md`.

- [ ] Migrate hook message templates from .txt files to YAML
  localization #added:2026-03-20-163801

### Phase 0.4.9: Injection Oversize Nudge

Spec: `specs/injection-oversize-nudge.md`. Read the spec before starting
any P0.4.9 task.

### Phase 0.4.10: Context Window Token Usage

Spec: `specs/context-window-usage.md`. Read the spec before starting any
P0.4.10 task.

### Phase 0.5 Cleanup

* Human: internal/recall/parser requires a serious refactoring; for example
  the parser object and its private and public methods need to go to its own
  package and other helper functions need to go to a different adjacent package.
* Human: internal/notify/notify.go requires refactoring (all functions bagged in
  one file; types need to go to types.go per convention etc etc)
* Human: split err package into sub packages.


- [ ] Refactor site/cmd/feed: extract helpers and types to core/, make Run
  public #added:2026-03-21-074859

- [ ] Add Use* constants for all cobra subcommand Use
  strings #added:2026-03-20-184639

- [ ] Systematic audit: extract all magic flag name strings across CLI commands
  into config/flag constants #added:2026-03-20-175155

- [-] Move generic string helpers from cli/add/core/strings.go to
  internal/format — file no longer exists, helpers already moved or deleted
  #added:2026-03-20-175046

- [ ] Add missing flag name constants (priority, section, file) and priority
  level constants (high, medium, low) to config/flag #added:2026-03-20-170842

### Phase 0: Ideas

**User-Facing Documentation** (from `ideas/done/REPORT-7-documentation.md`):
Docs are feature-organized, not problem-organized. Key structural improvements:

**Agent Team Strategies** (from `ideas/REPORT-8-agent-teams.md`):
8 team compositions proposed. Reference material, not tasks. Key takeaways:

- [ ] Scan all config/**/* constants and catalog which ones should be ctxrc
  entries for user configurability #priority:medium #added:2026-03-22-095552

- [ ] Update user-facing documentation for changed CLI flag
  shorthands #added:2026-03-21-102755

- [ ] Add Unicode-aware slugification for non-ASCII
  content #added:2026-03-21-070953

- [ ] Make TitleSlugMaxLen configurable via .ctxrc #added:2026-03-21-070944

- [ ] Spec and implement CRLF-to-LF newline normalization for journal and
  context files #added:2026-03-20-224845

- [ ] Test ctx on Windows — validate build, init, agent, drift, journal
  pipeline #added:2026-03-20-224835

- [ ] Evaluate Windows support for sysinfo.Collect and path
  handling #added:2026-03-20-194930

- [ ] Make doctor thresholds configurable via .ctxrc #added:2026-03-20-194923

- [ ] Evaluate cross-platform path handling in change/core/scan.go — git
  always
  uses "/" but UniqueTopDirs should consider filepath.ToSlash for Windows
  robustness #added:2026-03-20-182103

- [ ] Replace English-only Pluralize helper in change/core/detect.go with
  i18n-safe approach #added:2026-03-20-180502

- [ ] Replace ASCII-only alnum check in agent/core/score.go with
  unicode.IsLetter/IsDigit #added:2026-03-20-175943

### Phase S-0: Memory Bridge Groundwork

Prerequisites that unblocked the memory bridge phases.

### Phase MB: Memory Bridge Foundation (`ctx memory`)

Spec: `specs/memory-bridge.md`. Read the spec before starting any MB task.

Bridge Claude Code's auto memory (MEMORY.md) into `.context/` with discovery,
mirroring, and drift detection. Foundation for future import/publish phases.

### Phase MI: Memory Import Pipeline (`ctx memory import`)

Spec: `specs/memory-import.md`. Read the spec before starting any MI task.

Import entries from Claude Code's MEMORY.md into structured `.context/` files
using heuristic classification. Builds on Phase MB foundation (discover,
mirror, state).

- [-] MI.future: `--interactive` mode for agent-assisted classification —
  skipped: `--dry-run` covers review; agents can use `ctx add` directly for
  overrides; interactive CLI prompts don't compose with agent workflows

### Phase S-3: Blog Post — "Agent Memory is Infrastructure"

Spec: `specs/blog-agent-memory-infrastructure.md`.

### Phase MP: Memory Publish (`ctx memory publish`)

Spec: `specs/memory-publish.md`. Read the spec before starting any MP task.

Push curated context from `.context/` into Claude Code's MEMORY.md so the agent
sees structured project context on session start without needing hooks.

### Phase 9: Context Consolidation Skill `#priority:medium`

**Context**: `/ctx-consolidate` skill that groups overlapping entries by keyword
similarity and merges them with user approval. Originals archived, not deleted.
Spec: `specs/context-consolidation.md`
Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 3)

- [ ] Implement consolidation nudge hook: count sessions since last
  consolidation, nudge after 6. Spec:
  `specs/consolidation-nudge-hook.md` #added:2026-03-23-223000

- [ ] Auto-record consolidation baseline commit: `/ctx-consolidate` and `ctx
  system mark-consolidation` should stamp HEAD hash + date into
  `.context/state/consolidation.json` only on first invocation (write-once until
  reset). Subsequent consolidation sessions preserve the original baseline. The
  baseline resets only when the consolidation nudge counter resets (i.e., when a
  new feature cycle begins). This way multi-pass consolidation keeps the true
  starting point. Related:
  `specs/consolidation-nudge-hook.md` #added:2026-03-23-224000

### Phase EM: Extension Map Skill (`/ctx-extension-map`)

question: is this done; or needs planning?

### Phase WC: Write Consolidation

Baseline commit: `4ec5999` (Auto-prune state directory on session start).
Goal: consolidate user-facing messages into `internal/write/` as the central
output package. All CLI commands should route printed output through
this package.

- [ ] Migrate moc.go hardcoded strings to YAML or Go
  templates #added:2026-03-20-214922

- [ ] Design terminal-aware truncation for CLI output #added:2026-03-20-184509

### Phase SP: Configurable Session Prefixes

Spec: `specs/session-prefixes.md`. Read the spec before starting any SP task.

Replace hardcoded `session_prefix` / `session_prefix_alt` pair with a
user-extensible `session_prefixes` list in `.ctxrc`. Parser vocabulary
is not i18n text — it belongs in runtime config.

### Phase EH: Error Handling Audit

Systematic audit of silently discarded errors across the codebase.
Many call sites use `_ =` or `_, _ =` to discard errors without
any feedback. Some are legitimate (best-effort cleanup), most are
lazy escapes that hide failures.

- [ ] EH.1: Catalogue all silent error discards — recursive walk of
  `internal/`
  for patterns: `_ = `, `_, _ = `, `//nolint:errcheck`, bare `return` after
  error-producing calls. Group by category:
  (a) file close in defer — often legitimate but should log on failure
  (b) file write/read — data loss risk, must surface
  (c) os.Remove/Rename — state corruption risk
  (d) fmt.Fprint to stderr — truly best-effort, acceptable
  Commands: `grep -rn '_ =' internal/`, `grep -rn
      'nolint:errcheck' internal/`
  Output: spreadsheet in `.context/` with file, line, expression, category,
  and recommended action (log-stderr, return-error, acceptable-as-is).
  DoD: every `_ =` in the codebase is categorised and has a
  recommended action
  #priority:high #added:2026-03-14

- [ ] EH.2: Address category (b) — file write/read discards. These risk silent
  data loss. Fix: return the error, or at minimum emit to stderr with
  `fmt.Fprintf(os.Stderr, "ctx: ...: %v\n", err)` following the pattern
  established in `internal/log/event.go`.
  DoD: no write/read error is silently discarded
  #priority:high #added:2026-03-14

- [ ] EH.3: Address category (a) — file close in defer. Most are `defer func()
      { _ = f.Close() }()`. For read-only files, close errors are rare but
  should still surface. For write/append files, close can fail if the
  final flush fails — these are data loss. Fix: `if err := f.Close();
      err != nil { fmt.Fprintf(os.Stderr, "ctx: close %s: %v\n", path, err) }`.
  DoD: all defer-close sites log failures to stderr
  #priority:medium #added:2026-03-14

- [ ] EH.4: Address category (c) — os.Remove/Rename discards. These are state
  operations (rotation, pruning, temp file cleanup). Silent failure leaves
  stale state. Fix: stderr warning at minimum; for rotation/rename, consider
  returning the error.
  DoD: no Remove/Rename error is silently discarded
  #priority:medium #added:2026-03-14

- [ ] EH.5: Validate — `grep -rn '_ =' internal/` returns only category (d)
  entries (fmt.Fprint to stderr) and entries explicitly annotated as
  acceptable. Run `make lint && make test` to confirm no regressions.
  DoD: grep output is clean or fully annotated; CI green
  #priority:high #added:2026-03-14

- [ ] Add AST-based lint test to detect exported functions with no external
  callers #added:2026-03-21-070357

- [ ] Audit exported functions used only within their own package and make them
  private #added:2026-03-21-070346

- [ ] Audit and remove side-effect output from error-returning
  functions #added:2026-03-20-212212

### Phase ET: Error Package Taxonomy (`internal/err/`)

`errors.go` is 1995 lines with 188 functions in a single file. Split into
domain-grouped files. No API changes — same package, same function signatures,
just file reorganization.

Taxonomy (from prefix analysis):

| File         | Prefixes / Domain                                                                                                                                                      | ~Count |
|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|
| `memory.go`  | Memory*, Discover*                                                                                                                                                     | 17     |
| `parser.go`  | Parser*                                                                                                                                                                | 7      |
| `crypto.go`  | Crypto*, Encrypt*, Decrypt*, GenerateKey, SaveKey, LoadKey, NoKeyAt                                                                                                    | 14     |
| `task.go`    | Task*, NoTaskSpecified, NoTaskMatch, NoCompletedTasks                                                                                                                  | 8      |
| `journal.go` | LoadJournalState*, SaveJournalState*, ReadJournalDir, NoJournalDir, NoJournalEntries, ScanJournal, UnknownStage, StageNotSet                                           | 10     |
| `session.go` | Session*, FindSessions, NoSessionsFound, All*, Ambiguous*                                                                                                              | 8      |
| `pad.go`     | Edit*, Blob*, ReadScratchpad, OutFlagRequiresBlob, NoConflict*, Resolve*                                                                                               | 10     |
| `recall.go`  | Reindex*, Stats*, EventLog*                                                                                                                                            | 6      |
| `fs.go`      | Read*, Write*, Open*, Stat*, File*, Mkdir*, CreateDir, DirNotFound, NotDirectory, Boundary*                                                                            | 30     |
| `backup.go`  | Backup*, CreateBackup*, CreateArchive*                                                                                                                                 | 6      |
| `prompt.go`  | Prompt*, NoPromptTemplate, ListTemplates, ReadTemplate, NoTemplate                                                                                                     | 7      |
| `hook.go`    | Embedded*, Override*, UnknownHook, UnknownVariant, MarkerNotFound                                                                                                      | 6      |
| `skill.go`   | Skill*                                                                                                                                                                 | 2      |
| `config.go`  | UnknownProfile, ReadProfile, UnknownFormat, UnknownProjectType, InvalidTool, UnsupportedTool, NotInitialized, ContextNotInitialized, ContextDirNotFound, FlagRequires* | 12     |
| `errors.go`  | Remaining general-purpose: WorkingDirectory, CtxNotInPath, ReadInput, InvalidDate*, Reminder*, Drift*, Git*, Webhook*, etc.                                            | ~25    |

- [ ] Add freshness_files to .ctxrc defaults seeded by ctx init — currently
  the
  freshness config is only in the gitignored .ctxrc, so new clones don't get it.
  Consider a .ctxrc.defaults pattern or seeding via ctx init template.
  #priority:medium #added:2026-03-14-105143

- [ ] SEC.1: Security-sensitive file change hook — PostToolUse on Edit/Write
  matching security-critical paths (.claude/settings.local.json,
  .claude/settings.json, CLAUDE.md, .claude/CLAUDE.md,
  .context/CONSTITUTION.md). Three actions: (1) nudge user in-session, (2) relay
  to webhook for out-of-band alerting (autonomous loops), (3) append to
  dedicated security log (.context/state/security-events.jsonl) for forensics.
  Separate from general event log. Spec needed. #priority:high #added:2026-03-13

- [ ] O.5: Session timeline view — add --sessions flag to ctx system events.
  Per-session breakdown of eval/fired counts with hook list. See
  ideas/spec-hook-observability.md Phase 5 #added:2026-03-12-145401

- [ ] O.4: Doctor hook health check — surface hook activity in ctx doctor
  output
  (active/evaluated-never-fired/never-evaluated). See
  ideas/spec-hook-observability.md Phase 4 #added:2026-03-12-145401

- [ ] O.3: Skip reason logging — add eventlog.Skip() with standard reason
  constants (paused, throttled, condition-not-met). Instrument 19 hook
  early-exit paths. See ideas/spec-hook-observability.md Phase
  3 #added:2026-03-12-145401

- [ ] O.2: Event summary view — add --summary flag to ctx system events.
  Aggregates eval/fired counts per hook, shows last-eval/last-fired timestamps,
  lists never-evaluated hooks. See ideas/spec-hook-observability.md Phase
  2 #added:2026-03-12-145401

- [ ] O.1: Hook eval logging — wrap hook cobra commands to log 'eval' events
  on
  every invocation. Refactor Run() signatures from os.Stdin to io.Reader
  (peek+replay pattern). Adds eventlog.Eval(), EventTypeEval constant. See
  ideas/spec-hook-observability.md Phase 1 #added:2026-03-12-145401

- [ ] Companion intelligence recommendation: implement spec from
  ideas/spec-companion-intelligence.md — ctx doctor companion detection, ctx
  init recommendation tip, ctx agent awareness in
  packets #added:2026-03-12-133008

- [ ] Add configurable assets layer: allow users to plug their own YAML files
  for localization (language selection, custom text overrides). Currently all
  user-facing text is hardcoded in commands.yaml; need a mechanism to load
  user-provided YAML that overlays or replaces built-in text. This enables i18n
  without forking. #priority:low #added:2026-03-07-233756

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

- [ ] Make AutoPruneStaleDays configurable via ctxrc. Currently hardcoded to 7
  days in config.AutoPruneStaleDays; add a ctxrc key (e.g., auto_prune_days) and
  fallback to the default. #priority:low #added:2026-03-07-220512

- [-] Refactor check_backup_age/run.go: move consts (lines 23-24) to config,
  magic directories (line 59) to config, symbolic constants for strings (line
  72), messages to assets (lines 79, 90-91), extract non-Run functions to
  system/core, fix docstrings #priority:medium #added:2026-03-07-180020
  **Skipped 2026-04-16**: Superseded by specs/deprecate-ctx-backup.md
  (check_backup_age will be removed entirely, not refactored).

- [ ] Add ctxrc support for recall.list.limit to make the default --limit for
  recall list configurable. Currently hardcoded as config.DefaultRecallListLimit
  (20). #priority:low #added:2026-03-07-164342

- [ ] Extract journal/core into a standalone journal parser package —
  functionally isolated enough for its own package rather than remaining as
  core/ #added:2026-03-07-093815

- [ ] Move PluginInstalled/PluginEnabledGlobally/PluginEnabledLocally from
  initialize to internal/claude — these are Claude Code plugin detection
  functions, not init-specific #added:2026-03-07-091656

- [ ] Move guide/cmd/root/run.go text to assets, listCommands to separate file +
  internal/write #added:2026-03-07-090322

- [ ] Move drift/core/sanitize.go strings to assets #added:2026-03-07-090322

- [ ] Move drift/core/out.go output functions to internal/write per
  convention #added:2026-03-07-090322

- [ ] Move drift/core/fix.go fmt.Sprintf strings to assets — user-facing
  output
  text for i18n #added:2026-03-07-090322

- [ ] Move drift/cmd/root/run.go cmd.Print* output strings to internal/write per
  convention #added:2026-03-07-084152

- [ ] Extract doctor/core/checks.go strings — 105 inline Name/Category/Message
  values to assets (i18n) and config (Name/Category
  constants) #added:2026-03-07-083428

- [ ] Split deps/core builders into per-ecosystem packages — go.go, node.go,
  python.go, rust.go are specific enough for their own packages under deps/core/
  or deps/builders/ #added:2026-03-07-082827

- [ ] Audit git graceful degradation — verify all exec.Command(git) call sites
  degrade gracefully when git is absent, per project guide
  recommendation #added:2026-03-07-081625

- [ ] Fix 19 doc.go quality issues: system (13 missing subcmds), agent (phantom
  refs), load/loop (header typo), claude (stale migration note), 13 minimal
  descriptions (pause, resume, task, notify, decision, learnings, remind,
  context, eventlog, index, rc, recall/parser,
  task/core) #added:2026-03-07-075741

- [ ] Move cmd.Print* output strings in compact/cmd/root/run.go to
  internal/write per convention #added:2026-03-07-074737

- [ ] Extract changes format.go rendering templates to assets — headings,
  labels, and format strings are user-facing text for
  i18n #added:2026-03-07-074719

- [ ] Lift HumanAgo and Pluralize to a common package — reusable time
  formatting, used by changes and potentially
  status/recall #added:2026-03-07-074649

- [ ] Extract isAlnum predicate for localization — currently ASCII-only in
  agent
  keyword extraction (score.go:141) #added:2026-03-07-073900

- [ ] Make stopwords configurable via .ctxrc — currently embedded in assets,
  domain users need custom terms #added:2026-03-07-073900

- [ ] Make recency scoring thresholds and relevance match cap configurable via
  .ctxrc — currently hardcoded in config (7/30/90 days, cap
    3) #added:2026-03-07-073900

- [ ] Make DefaultAgentCooldown configurable via .ctxrc — currently hardcoded
  at
  10 minutes in config #added:2026-03-07-073106

- [ ] Make TaskBudgetPct and ConventionBudgetPct configurable via .ctxrc —
  currently hardcoded at 0.40 and 0.20 in config #added:2026-03-07-072714

- [ ] Localization inventory: audit config constants, write package templates,
  and assets YAML for i18n mapping — low priority, most users are
  English-first
  developers #added:2026-03-06-192419

- [ ] Consider indexing tasks and conventions in TASKS.md and CONVENTIONS.md
  (currently only decisions and learnings have index
  tables) #added:2026-03-06-190225

- [ ] Implement journal compaction: Elastic-style tiered storage with tar.gz
  backup. Spec: specs/journal-compact.md #added:2026-03-31-110005

- [ ] Validate .ctxrc against ctxrc.schema.json at load time — schema is
  embedded but never enforced, doctor does field-level checks without using
  it #added:2026-03-06-174851


- [ ] Add PostToolUse session event capture. Append lightweight event records
  (tool name, files touched, timestamp) to .context/state/session-events.jsonl
  on significant PostToolUse events (file edits, git operations, errors). Not
  SQLite — just JSONL append. This feeds the PreCompact snapshot hook with
  richer input so it can report what the agent was actively working on, not just
  static file state. #added:2026-03-06-185126

- [ ] Add next-step hints to ctx agent and ctx status output. Append actionable
  suggestions based on context health (e.g. stale tasks, high completion ratio,
  drift findings). Pattern learned from GitNexus self-guiding agent
  workflows. #added:2026-03-06-184829

- [ ] Implement PreCompact and SessionStart hooks for session continuity across
  compaction. Wire ctx agent --budget 4000 to both events: PreCompact outputs
  context packet before compaction so compactor preserves key info; SessionStart
  re-injects context packet so fresh/post-compact sessions start oriented. Two
  thin ctx system subcommands, two entries in hooks.json. See
  ideas/gitnexus-contextmode-analysis.md for design
  rationale. #added:2026-03-06-184825

- [ ] Audit fatih/color removal across ~35 files — removed from recall/run.go,
  recall/lock.go, write/validate.go; ~30 files remain. Separate consolidation
  pass. #added:2026-03-06-050140

- [ ] Audit remaining 2006-01-02 usages across codebase — 5+ files still use
  the
  literal instead of config.DateFormat. Incremental
  migration. #added:2026-03-06-050140

- [ ] WC.2: Audit CLI packages for direct fmt.Print/Println usage — candidates
  for migration #added:2026-03-06

### Phase WC2: Write Output Block Consolidation

Spec: `specs/write-output-consolidation.md`. Read the spec before starting any
WC2 task.

Consolidate multi-line imperative `cmd.Println` sequences in `internal/write/`
into pre-computed single-print block patterns. Separates conditional logic from
I/O and replaces 4-8 individual Tpl\* constants per function with one
block template.

- [ ] WC2.1: Tier 1 — Consolidate multi-line functions with no conditionals:
  `InfoInitNextSteps`, `InfoObsidianGenerated`, `InfoJournalSiteGenerated`,
  `InfoDepsNoProject`, `ArchiveDryRun`, `ImportScanHeader`. Add `TplXxxBlock`
  YAML entries, wire through embed.go + config.go, remove replaced individual
  constants. #added:2026-03-17
- [ ] WC2.2: Tier 2a — Consolidate conditional functions in info.go:
  `InfoLoopGenerated` (pre-compute iterLine). Prove the pre-computation pattern
  on the function that motivated this spec. #added:2026-03-17
- [ ] WC2.3: Tier 2b — Consolidate conditional functions in
  sync/recall/notify:
  `SyncResult`, `CtxSyncHeader`, `CtxSyncAction`, `SessionMetadata`,
  `TestResult`, `SyncDryRun`, `PruneSummary`. Each needs 1-3 pre-computed
  strings before the single print call. #added:2026-03-17
- [ ] WC2.4: Constant cleanup — verify all replaced individual `TplXxx*`
  config
  vars, `TextDescKey*` constants, and YAML entries are removed. Run `make lint`
  and `go test ./internal/write/...` to confirm no
  regressions. #added:2026-03-17
- [ ] WC2.5: Update CONVENTIONS.md — add a "Write Package Output" subsection
  documenting the pre-compute-then-print pattern for future functions with 4+
  Printlns and conditionals. #added:2026-03-17

## MCP-related

### Phase MCP-V3: MCP v0.3 Expansion

- [ ] Add drift check: MCP prompt coverage vs bundled skills — programmatic
  check comparing config/mcp/prompt constants against assets.ListSkills() to
  detect skills without MCP prompt equivalents. Pair with the tool coverage
  drift check. @CoderMungan #priority:medium #added:2026-03-15-120519

- [ ] MCP v0.3: expand MCP prompts to cover more skills — current 5 prompts
  (session-start, add-decision, add-learning, reflect, checkpoint) are a subset
  of ~30 bundled skills. Evaluate which skills benefit from protocol-native MCP
  prompt equivalents. Decision 2026-03-06 established 'Skills stay CLI-based;
  MCP Prompts are the protocol equivalent.' @CoderMungan
  #priority:medium #added:2026-03-15-120519

- [ ] Add drift check: MCP tool coverage vs CLI commands — programmatic check
  that compares registered MCP tool names (config/mcp/tool) against ctx CLI
  subcommands to detect newly added CLI commands without MCP equivalents. Could
  be a drift detector check or a compliance test. @CoderMungan
  #priority:medium #added:2026-03-15-120116

- [ ] MCP v0.3: expose additional CLI commands as MCP tools — candidates:
  ctx_load (full context packet), ctx_agent (token-budgeted packet), ctx_reindex
  (rebuild indices), ctx_sync (reconcile docs/code), ctx_doctor (health check).
  Evaluate which provide value over the protocol vs requiring terminal
  interaction. @CoderMungan #priority:medium #added:2026-03-15-120025

- [ ] Make MCP defaults configurable via .ctxrc — add mcp_recall_limit,
  mcp_truncate_len, mcp_truncate_content_len, mcp_min_word_len,
  mcp_min_word_overlap fields to .ctxrc schema; expose via rc.MCP*() with
  fallback to config/mcp/cfg defaults; update tools.go to read from rc instead
  of cfg constants. @CoderMungan #priority:medium #added:2026-03-15-114700

- [ ] MCP tools.go cleanup pass: magic strings, duplicated fragments, nested
  templates. Lines: 461:481 + 186:196 duplicated code; 335 magic number; 382:385
  nested TextDescs → single template; 390+851 magic time literal; 443+499+800
  magic words; 557+892+902 magic numbers; 590+638 nested TextDesc templating;
  820 prefixed %s; 854 suffix %s #priority:high #added:2026-03-15-110429

### Phase MCP-SAN: MCP Server Input Sanitization

[ ] Assignee: @CoderMungan -- https://github.com/ActiveMemory/ctx/issues/49

### Phase MCP-COV: MCP Test Coverage

[ ] Assignee: @CoderMungan -- https://github.com/ActiveMemory/ctx/issues/50

## Later

### Phase PR: State Pruning (`ctx system prune`)

Clean stale per-session state files from `.context/state/`. Files with UUID
session ID suffixes accumulate ~6-8 per session with no cleanup. Strategy:
age-based — prune files older than N days (default 7).

- [ ] Regenerate site/ for state-maintenance recipe
  (docs/recipes/state-maintenance.md added but site not
  rebuilt) #added:2026-03-05-205425

- [ ] Audit remaining global tombstones for session-scoping:
  backup-reminded, ceremony-reminded, check-knowledge,
  journal-reminded, version-checked, ctx-wrapped-up all have
  the same cross-session suppression bug as
  memory-drift-nudged #added:2026-03-05-205425

- [ ] F.2: ctx journal import (remote) — import Claude Code
  session JSONLs from local or remote (~/.claude/projects/)
  into local ~/.claude/projects/. Pure Go: local copy with
  os.CopyFS-style walk, remote via os/exec ssh+scp (no rsync
  dependency). --source flag accepts local path or user@host.
  --dry-run shows what would be copied. Skips existing files
  (content-addressed by UUID filenames). Enables journal export
  from sessions that ran on other machines.
  #added:2026-03-05-141912

- [ ] P0.5: Blog: "Building a Claude Code Marketplace Plugin"
  — narrative from session history, journals, and git diff of
  feat/plugin-conversion branch. Covers: motivation (shell
  hooks to Go subcommands), plugin directory layout,
  marketplace.json, eliminating make plugin, bugs found during
  dogfooding (hooks creating partial .context/), and the fix.
  Use /ctx-blog-changelog with branch diff as source material.
  #added:2026-02-16-111948
- [ ] P9.2: Test manually on this project's LEARNINGS.md (20+ entries).
  #priority:medium #added:2026-02-19
- [ ] P0.8.1: Install golangci-lint on the integration server #for-human
  #priority:medium #added:2026-02-23 #added:2026-02-23-170213
- [ ] PM.3: Review hook diagnostic logs after a long session. Check
  `.context/logs/check-persistence.log` and
  `.context/logs/check-context-size.log` to verify hooks fire correctly.
  Tune nudge frequency if needed. #priority:medium #added:2026-02-09
- [ ] PM.4: Run `/consolidate` to address codebase drift. Considerable drift has
  accumulated (predicate naming, magic strings, hardcoded permissions,
  godoc style). #priority:medium #added:2026-02-06
- [ ] Improve test coverage for core packages at 0% #added:2026-03-20-164324

- [ ] PM.7: Aider/Cursor parser implementations: the recall architecture was
  designed for extensibility (tool-agnostic Session type with
  tool-specific parsers). Adding basic Aider and Cursor parsers would
  validate the parser interface, broaden the user base, and fulfill
  the "works with any AI tool" promise. Aider format is simpler than
  Claude Code's. #priority:medium #source:report-6 #added:2026-02-17

## Future

- [ ] P0.8.5: Enable webhook notifications in worktrees. Currently `ctx notify`
  silently fails because `.context.key` is gitignored and absent in
  worktrees. For autonomous runs with opaque worktree agents, notifications
  are the one feature that would genuinely be useful. Possible approaches:
  resolve the key via `git rev-parse --git-common-dir` to find the main
  checkout, or copy the key into worktrees at creation time (ctx-worktree
  skill). #priority:medium #added:2026-02-22
- [ ] P0.9.2: Split cli-reference.md (1633 lines) into command group pages:
  cli-overview, cli-init-status, cli-context, cli-recall, cli-tools,
  cli-system —
  each page covers a natural command group with its subcommands and flags
  #added:2026-02-24-204208
- [ ] P0.9.3: Investigate proactive content suggestions:
  docs/recipes/publishing.md claims
  agents suggest blog posts and journal rebuilds at natural moments, but no hook
  or playbook mechanism exists to trigger this — either wire it up (e.g.
  post-task-completion nudge) or tone down the docs to match reality
  #added:2026-02-24-185754
- [ ] PG.1: Add agent/tool compatibility matrix to prompting guide —
  document which
  patterns degrade gracefully when agents lack file access, CLI tools, or
  ctx integration. Treat as a "works best with / degrades to" table.
  #priority:medium #added:2026-02-25
- [ ] PG.2: Add versioning/stability note to prompting guide — "these
  principles are
  stable; examples evolve" + doc date in frontmatter. Needed once the guide
  becomes canonical and people start quoting it.
  #priority:low #added:2026-02-25
- [ ] P0.1: Brainstorm: Standardize drift-check comment format and
  integrate with
  `/ctx-drift` — formalize ad-hoc `<!-- drift-check: ... -->` markers, teach
  drift skill to parse/execute them, publish pattern in docs/recipes. Benefits
  tooling/CLI but AI handles ad-hoc fine for now.
  #priority:medium #added:2026-02-28
- [ ] F.1: MCP server integration: expose context as tools/resources via Model
  Context Protocol. Would enable deep integration with any
  MCP-compatible client. #priority:low #source:report-6
- [ ] Q.1: Docstring cross-reference audit — compliance test that
  flags docstrings
  mentioning domains that don't match their callers. Start with `write/**`,
  extend to all `internal/`. Spec: `specs/docstring-cross-reference-audit.md`
  #priority:medium #added:2026-03-17

- [ ] Migrate Sprintf-based templates (tpl_*.go) to Go text/template or embedded
  template files — ObsidianReadme, LoopScript, and other multi-line format
  strings that can't move to YAML #added:2026-03-18-163629

- [ ] Split internal/assets/embed_test.go — tests that call read/ packages
  must
  move to their respective read/ package to avoid import
  cycles #added:2026-03-18-192914

- [ ] Improve recall/core format tests — replace hardcoded string assertions
  (e.g. Contains Tokens) with semantic checks that verify structure and values,
  not label text #added:2026-03-19-194645

### Phase BT: Build Tooling — `cmd/ctxctl`

Replace shell-based build scripts (Makefile shell
expansions, `hack/build-all.sh`,
`hack/release.sh`, `hack/tag.sh`, `sync-*`/`check-*` targets) with a first-class
Go binary at `cmd/ctxctl`. Shares internal packages with `ctx` (version, assets,
embed FS). Installable: `go
install github.com/ActiveMemory/ctx/cmd/ctxctl@latest`.
Eliminates `jq` build dependency. Testable, cross-platform.

- [ ] Bug: release script versions.md table insertion fails silently. The sed
  pattern on line 133 uses `$` anchor but the actual Markdown table header has
  column padding spaces before the trailing `|`. The row is never inserted. Fix:
  relax the header match pattern or switch to a simpler approach (e.g., insert
  after the separator line directly). Also verify the "latest stable" sed
  handles trailing `).\n` correctly. #priority:high #added:2026-03-23-221500

- [ ] Replace hack/lint-drift.sh with AST-based Go tests in internal/audit/.
  Spec: `specs/ast-audit-tests.md` #added:2026-03-23-210000

- [ ] Rewrite lint-style scripts in Go as ctxctl subcommands —
  blocked: prerequisite ctxctl does not exist yet. Deferred.
  #added:2026-03-29-082958

Dividing line: `ctx` is the user/agent tool, `ctxctl` is
the maintainer/contributor
tool. If a developer clones the repo and needs to build, test, release,
or validate
— that's `ctxctl`. If a user is working in a project and needs context —
that's `ctx`.

Strong fits beyond build/release:

- `ctxctl plugin package` — package .claude-plugin for marketplace publishing
- `ctxctl plugin validate` — validate plugin.json, hooks.json, skill structure
- `ctxctl doctor` — contributor pre-flight (Go version, tools, GPG, hooks);
  absorbs `hack/gpg-fix.sh` and `hack/gpg-test.sh`
- `ctxctl changelog` — deterministic release notes from git log

Reasonable fits if project grows:

- `ctxctl test smoke` — replaces the shell pipeline in `make smoke`
- `ctxctl site build/serve` — wraps zensical + feed generation
- `ctxctl mcp register` — replaces `hack/gemini-search.sh` and future
  MCP registrations

Not a fit (keep in `ctx`):

- Anything user-facing in a project context (status, agent, drift, recall)
- Anything Claude Code hooks call — hooks must call `ctx`, not `ctxctl`

- [ ] Design `ctxctl` CLI surface: `ctxctl sync`, `ctxctl build`, `ctxctl
  release`, `ctxctl check`, `ctxctl tag` #added:2026-03-25-050000
- [ ] Implement `ctxctl sync` — stamps VERSION into plugin.json + syncs why
  docs; replaces `sync-version`, `sync-why` #added:2026-03-25-050000
- [ ] Implement `ctxctl check` — drift checks: version sync, why docs,
  lint-drift, lint-docs; replaces `check-*` targets #added:2026-03-25-050000
- [ ] Implement `ctxctl build` — cross-platform builds with version stamping;
  replaces `build-all.sh` #added:2026-03-25-050000
- [ ] Implement `ctxctl release` — full release flow (sync, build, tag,
  checksums); replaces `release.sh` + `tag.sh` #added:2026-03-25-050000
- [ ] Simplify Makefile to thin wrappers: `make build` → `go run ./cmd/ctxctl
  build` #added:2026-03-25-050000
- [ ] Remove `jq` build dependency once ctxctl handles JSON
  natively #added:2026-03-25-050000

- [ ] Implement MCP warm-up in /ctx-remember session ceremony — when a
  graph/RAG
  tool is configured in .ctxrc, run one orientation query at session start to
  build procedural familiarity. Spec:
  `ideas/spec-mcp-warm-up-ceremony.md` #added:2026-03-25-120000

- [ ] Update ctx doctor to check for graph tool availability — detect if a
  graph/RAG MCP is configured in .ctxrc, verify connection status, recommend
  installation if missing #added:2026-03-25-120000

- [ ] Explore pluggable graph tool interface — replace hardcoded GitNexus
  references in skill text with configurable .ctxrc graph_tool key. Skills use
  template placeholder instead of literal tool names. Define minimum interface
  contract (query, context, impact). Spec:
  `ideas/spec-mcp-warm-up-ceremony.md` #added:2026-03-25-120000

### Phase: ctx Hub follow-ups (PR #60)

**Context**: PR #60 `feat: ctx Hub for cross-project knowledge
sharing` (parlakisik) merged despite open review feedback from @bilersan and
a pending review request. Author is heads-down on his Ph.D.; these tasks
capture the cleanup and documentation debt we accepted by merging.
PR: https://github.com/ActiveMemory/ctx/pull/60
Review with findings:
https://github.com/ActiveMemory/ctx/pull/60#pullrequestreview-PRR_kwDOQ9VoNc7ze3nA

#### Build / platform

- [ ] Add Windows job to CI so this class of regression is caught at PR time,
  not by reviewers running local builds. #priority:high #added:2026-04-11 #pr:60
- [ ] Triage the 16 package-level test failures @bilersan reported on Windows
  — classify as platform-specific vs genuine bugs. #added:2026-04-11 #pr:60

#### Convention drift

- [ ] Audit `internal/hub`, `internal/cli/connect`, `internal/cli/hub`,
  `internal/cli/serve` against CONVENTIONS.md (godoc format, import aliases,
  error wrapping, package layout). #added:2026-04-11 #pr:60
- [ ] Run `/ctx-code-review` over the hub subsystem for edge cases missed in
  the merge: token rotation, connection-config migration, Raft leader
  handoff failure modes, sync cursor corruption recovery. #added:2026-04-11
  #pr:60

#### User-facing docs (cornerstone — scope first)

- [ ] Document the auto-sync-on-session-start hook: what it does, how to
  opt out, interaction with existing UserPromptSubmit hooks, performance
  impact on large hubs. Partially covered in connect.md (`check-hub-sync`
  mention); a dedicated section is still owed. #added:2026-04-11 #pr:60
- [ ] Add an **architecture** section to `ARCHITECTURE.md` /
  `DETAILED_DESIGN.md` covering: JSONL append-only store, JSON-over-gRPC
  codec (no protoc), fan-out broadcaster, Raft-lite (election only, data
  via gRPC sync), sequence-based replication. #added:2026-04-11 #pr:60
- [ ] Record a DECISION explaining why we merged PR #60 with known Windows
  breakage and convention drift — trade-off, author context, mitigation
  plan (this task group). #added:2026-04-11 #pr:60
- [ ] Update CONVENTIONS.md if any new patterns from the hub are worth
  canonicalizing (gRPC handler layout, JSONL store access, bearer-token
  middleware). #added:2026-04-11 #pr:60

#### Framing and mental model (2026-04-11 follow-up)

#### Design follow-ups surfaced by the brainstorm (2026-04-11)

- [ ] Decide the product story: "personal cross-project brain",
  "small trusted team", or both — then align the overview, recipes,
  and CONTRIBUTING guidance to match. #priority:high #added:2026-04-11
  #pr:60
- [ ] Server-enforce `Origin` on publish: reject entries whose
  `Origin` does not match the authenticated client's `ProjectName`.
  Closes a spoofing vector and eliminates accidental mislabeling.
  Small change in `internal/hub/handler.go publish()`.
  #priority:high #added:2026-04-11 #pr:60
- [ ] Hash `clients.json` tokens or move them behind the local
  keyring (reuse `internal/crypto`). Removes the plaintext-token
  footgun documented in the security page.
  #priority:high #added:2026-04-11 #pr:60
- [ ] Explore journal-entry → `learning` export path: the density
  users expect from "shared context" lives in enriched journal
  entries, not in manually written `ctx add learning`. Would let
  the hub surface the lessons agents already recorded in sessions
  without actually replicating journals. #added:2026-04-11 #pr:60

#### Phase: Hub identity layer for public-internet usage (2026-04-11)

**Context**: The current hub has no concept of user identity.
Tokens identify **projects**, not humans. `Origin` is
self-asserted on publish. `clients.json` stores tokens in
plaintext. For the "personal" and "small trusted team" stories
(overview.md Stories 1 and 2) this is acceptable — the trust
model is "everyone holding a token is friendly."

For public-internet usage (the "Story 3" shape we explicitly
declared out of scope in the overview) these become real gaps:
no per-user attribution, no way to revoke individual humans, no
audit trail that proves who published what, and `clients.json`
compromise equals total hub compromise.

**Near-term MVP**: a pre-seeded identity registry owned by the
sysadmin. Instead of dynamic token issuance via admin token,
the hub reads a `users.json` file the sysadmin hand-edits, and
client registration validates against that pre-seeded list.
This is simpler than OAuth/OIDC, doesn't require a separate
identity service, and matches how internal services at small
orgs usually start before adopting an SSO.

**Eventual design requirements** (decision record TBD):

- Per-human identity, not per-project
- Tokens tied to a user ID, not a project name
- Server-enforced `Origin` matches the authenticated user (or
  a user's declared project list, with server validation)
- Revocation by removing a user row from the registry and
  forcing token rotation
- Hashed token storage at rest
- Optional: attribution-bearing audit log distinct from
  `entries.jsonl`

The following tasks feed into this track (they already exist
in the "Design follow-ups surfaced by the brainstorm" section
above; do not duplicate here):

- Server-enforce `Origin` on publish (blocks spoofing)
- Hash `clients.json` tokens (blocks plaintext compromise)
- Decide the fate of `Entry.Author` (promote, drop, or keep
  unauthenticated)

Tasks unique to this phase:

- [ ] Write a spec for the sysadmin-curated identity registry:
  filename, format, schema, bootstrap flow, revocation
  procedure, migration path from today's `clients.json`.
  `specs/hub-identity-registry.md`. Must resolve:

    - **Token issuance**: out-of-band on the server
      (`ctx hub users add` prints the plaintext token once
      on stdout; only a hash is persisted).
    - **Client pickup**: user receives the token out-of-band
      and runs `ctx connect register <host> --token
    ctx_cli_... --project <name>`; hub validates against
      the registry.
    - **TTL decision** (pick one, document in the spec):
        * **Option A** (recommended): no TTL, manual revocation
          only. `ctx hub users remove <id>` is the only
          expiry path. Matches today's `clients.json`
          semantics, zero surprise breakage on migration.
        * **Option B**: optional `expires_at` per user row.
          Tokens without it are valid forever (Option A
          behavior); tokens with it are rejected after the
          timestamp. Ship as an additive follow-up.
        * **Option C** (explicitly rejected): rolling
          expiry based on `last_used_at`. Garbage-collects
          dormant tokens but breaks users who take long
          vacations. Not worth the support cost.
    - **Revocation procedure**: sysadmin edits `users.json`,
      signals the hub to reload, affected tokens fail
      immediately on next RPC.
    - **Migration from `clients.json`**: one-shot converter
      that reads today's `clients.json`, prompts the
      sysadmin for a `user_id` per row, and writes
      `users.json`. Leave `clients.json` in place as a
      read fallback during migration, delete once
      everyone is on the new path.

  #priority:high #added:2026-04-11 #pr:60
- [ ] Implement `users.json` format: `{user_id: {project_ids:
  [...], token_hash: "...", created_at: "...", notes: "..."}}`.
  Read on hub start and on each Register RPC. Hot-reload via
  SIGHUP or file watcher. #added:2026-04-11 #pr:60
- [ ] Change `Register` RPC semantics: instead of minting a
  new client token from the admin token, look up the
  requested `ProjectName` in `users.json`. Reject if not
  pre-seeded. Return the pre-hashed token only if the caller
  presents an initial-provisioning credential the sysadmin
  seeded alongside the registry row. #added:2026-04-11 #pr:60
- [ ] Add `ctx hub users` subcommand group for sysadmin
  operations: `add`, `remove`, `rotate`, `list`. These edit
  `users.json` directly and signal the running hub to
  reload. #added:2026-04-11 #pr:60
- [ ] Add per-user audit log (`audits.jsonl` beside
  `entries.jsonl`). Each RPC records user_id, method, result
  status, timestamp. Separate from `entries.jsonl` so it can
  be retained on a different schedule. #added:2026-04-11
  #pr:60
- [ ] Write `docs/security/hub-identity.md` explaining the
  registry-based identity model, the threat model it closes,
  the threats it still doesn't close, and the operational
  procedures (seed the registry, rotate a token, revoke a
  user). #added:2026-04-11 #pr:60
- [ ] Decide whether to ship the identity layer as a
  **breaking change** (existing `clients.json` deployments
  must migrate) or as an **opt-in flag** (`ctx hub start
  --identity users.json`). Document in the spec above.
  #added:2026-04-11 #pr:60
- [ ] Update the hub overview and team recipe to name the
  identity registry as the "upgrade path to larger teams"
  story: "once your team grows past ~10 people or you need
  auditable attribution, enable the identity registry." The
  current overview treats Story 3 as unsupported — with the
  registry this becomes Story 2.5: "small trusted team with
  real attribution." #added:2026-04-11 #pr:60
- [ ] Stretch: OIDC/OAuth bridge. Once the registry layer is
  stable, consider adding an optional provider bridge so
  `users.json` can be auto-populated from an external
  identity source (Google Workspace, GitHub orgs, etc.). Not
  a near-term priority — registry-only covers the first
  order of magnitude of users. #added:2026-04-11 #pr:60
- [ ] Stretch: signed-claim / PKI authentication. The
  sysadmin-registry MVP and the OIDC bridge are both
  **bearer token** models — possession of the token bytes
  is identity. This is fine for trusted orgs but has
  well-known replay/rotation/identity limits for true
  public-internet usage.

  The next tier up is **asymmetric / signed-claim** auth:
  sysadmin holds a private signing key, issues short-lived
  claims `{user, project, expiry}` signed with that key,
  clients present the signed claim on each RPC, server
  verifies with the public key. Benefits:

    - Private key never leaves the sysadmin's machine.
    - Claims expire in minutes → revocation is automatic.
    - Each claim carries identity cryptographically.
    - No per-RPC registry lookup — signature verification
      is cheap.

  Reference designs to evaluate: JWT (RS256/ES256/EdDSA),
  mTLS client certificates, SPIFFE/SPIRE workload
  identities. Decision driver: does ctx ever want to run
  as a real public-internet service, or does "trusted
  team" always remain the upper bound?

  This is the Story 3 → true multi-tenant upgrade. Not a
  near-term priority; captured here so the registry-first
  MVP doesn't get confused for a final-state solution.
  #added:2026-04-11 #pr:60

#### Phase: "dependency-free" claim cleanup (2026-04-11)

**Context**: The design-invariant list in marketing and
reference docs historically included "dependency-free"
as one of five properties (alongside local-first,
file-based, CLI-driven, developer-controlled). This was
accurate when ctx was a single Go binary with no
external services. PR #60 (hub), the zensical
integration (`ctx serve`), the Claude Code plugin +
MCP, and future networked features make the blanket
claim false.

**Replacement framing (adopted 2026-04-11)**:
"**single-binary core**". The context persistence path
(`init`, `add`, `agent`, `status`, `drift`, `load`,
`sync`, `compact`, `task`, `decision`, `learning`, and
siblings) remains a single Go binary with no required
runtime dependencies. Optional integrations — `ctx
trace` (needs `git`), `ctx serve` (needs `zensical`),
`ctx` Hub (needs a running hub), Claude Code plugin
(needs `claude`) — are opt-in and each declares its
dependency explicitly.

This framing is load-bearing: it communicates the
design intent (nothing you don't opt into) without
claiming a literal falsehood.

- [-] `docs/thesis/index.md:412` (the primitive
  comparison table saying "Document: Zero-dependency:
  Yes"): left intact. The claim is about the document
  primitive itself (markdown files have no runtime
  deps), not about ctx as an implementation. Accurate.
  #added:2026-04-11 #skipped:primitive-claim-is-correct
- [ ] Add a design-invariants reference note: the
  blanket claim "dependency-free" MUST NOT be
  reintroduced in new docs. Any new framing should use
  "single-binary core" or name the specific path
  (e.g., "persistence path", "agent packet assembly").
  #priority:medium #added:2026-04-11
- [ ] Pre-release re-sweep: before each minor release,
  grep `docs/`, `README.md`, and any blog drafts for
  `dependency-free|dependency free|zero dependencies|
  no dependencies` and verify each occurrence is
  scoped to a path that is still dependency-free. Add
  to the release runbook. #priority:medium
  #added:2026-04-11
- [ ] Update `docs/reference/design-invariants.md` to
  explicitly list "single-binary core" as an invariant
  with the scope definition, so future doc authors
  have a canonical source to reference instead of
  re-deriving the phrase. #priority:medium
  #added:2026-04-11

#### Phase: Hub security audit (2026-04-11)

**Context**: Full security audit of the hub subsystem,
completed during the PR #60 follow-up brainstorm as a
precondition for any public-internet deployment. 30
findings total — 5 Critical, 12 High, 7 Medium, 4 Low, 2
Info — covering transport security, identity,
attribution, DoS surface, Raft cluster integrity, and
storage integrity.

The audit lives at `specs/hub-security-audit.md` and is
the canonical reference for the rest of the hub security
work. Each finding has a concrete remediation,
complexity estimate, and cross-reference to existing
tasks where applicable. The spec also contains
recommendations grouped by timeline (do-now / short /
medium / long).

**Per-story verdicts from the audit**:

- **Story 1** (personal cross-project brain, localhost):
  acceptable as-is. No adversary in scope.
- **Story 2** (small trusted team on LAN): acceptable
  with documented caveats — LAN private, hub host
  hardened, admin token held only by the sysadmin. The
  `hub-team.md` recipe already names these.
- **Story 3** (public-internet / multi-user): **UNSAFE**.
  Do not deploy. Five critical findings apply, several
  high-severity findings compound catastrophically
  without transport security, and the Raft cluster is
  a remote unauthenticated DoS surface.

**This phase tracks the findings as actionable work**.
Individual findings are numbered H-01 through H-30 in
the spec; this task list references them by number and
links back to the spec for detail.

- [ ] Read and internalize
  [`specs/hub-security-audit.md`](../specs/hub-security-audit.md)
  before starting any hub-security implementation.
  The spec is the single source of truth for findings,
  severity, and remediation patterns. #priority:high
  #added:2026-04-11 #pr:60

**Do-now track** (prerequisites for non-localhost deployments):

- [ ] **H-01** Add server-side TLS: `--tls-cert` and
  `--tls-key` flags on `ctx hub start`, wire into
  `grpc.NewServer` via `grpc.Creds`. Keep plaintext
  default for Story 1. #priority:critical
  #added:2026-04-11 #pr:60 #audit:H-01
- [ ] **H-02** Add client-side TLS: accept `grpc://`
  and `grpcs://` schemes in `hub_addr`. Update
  `NewClient`, `replicateOnce`, `NewFailoverClient` to
  switch credentials per scheme. Optional `--ca-cert`
  for self-signed. Update
  `docs/recipes/hub-multi-machine.md` to document both
  forms (the current nginx-reverse-proxy recommendation
  is un-implementable until this ships). #priority:critical
  #added:2026-04-11 #pr:60 #audit:H-02
- [ ] **H-04** Server-enforce `Origin` on publish:
  `validateBearer` attaches `ClientInfo` to context;
  `handler.go publish()` overwrites `pe.Origin` with
  the authenticated `ClientInfo.ProjectName` before
  store. Add a test that a client authenticated as
  `alpha` cannot publish as `beta`. #priority:high
  #added:2026-04-11 #pr:60 #audit:H-04
- [ ] **H-15** Fix `appendFile` in `internal/hub/persist.go`
  to use real `O_APPEND` instead of read-all-rewrite.
  Closes both a performance bug (O(N²) publishes) and
  a data-loss risk (partial write can truncate history).
  #priority:high #added:2026-04-11 #pr:60 #audit:H-15

**Short-term track** (Story 2 hardening):

- [ ] **H-03** Hash `clients.json` tokens with argon2id.
  One-shot migration reads old file, hashes each token,
  rewrites. Plaintext token only passes through memory
  at registration time; disk only stores hashes.
  Already referenced in the design-follow-ups section
  above; this entry ties it to the audit. #priority:high
  #added:2026-04-11 #pr:60 #audit:H-03
- [ ] **H-08** Per-token Publish rate limiting using
  `golang.org/x/time/rate`. Starting target: 10 entries/sec
  per token, 100 burst. Return `ResourceExhausted` with
  Retry-After hint. #priority:high #added:2026-04-11 #pr:60
  #audit:H-08
- [ ] **H-09** Per-token Listen stream cap (suggested
  limit: 4 concurrent streams per token, 256 total).
  Track in the `fanOut` struct; reject further subscribes
  with `ResourceExhausted`. #priority:high
  #added:2026-04-11 #pr:60 #audit:H-09
- [ ] **H-17** Cap `PublishRequest.Entries` at 32 per
  request; reject larger batches with
  `InvalidArgument`. Document the limit. #priority:high
  #added:2026-04-11 #pr:60 #audit:H-17
- [ ] **H-18** Add `audits.jsonl` as a per-RPC audit log
  distinct from `entries.jsonl`. Records
  `{ts, method, user, project, status, entry_count}`
  per call, including authentication failures. Exposed
  via `ctx hub status --audit`. Independent rotation
  cadence. Already referenced in the identity-layer
  phase; this entry ties it to the audit. #priority:high
  #added:2026-04-11 #pr:60 #audit:H-18
- [ ] **H-19** Implement real revocation: `ctx hub users
  remove <id>` edits the registry and signals the hub
  to reload via `fsnotify`. Revoked tokens fail
  immediately on next RPC. Revocation events logged to
  `audits.jsonl`. Merged with the Hub identity layer
  phase implementation. #priority:high #added:2026-04-11
  #pr:60 #audit:H-19
- [ ] **H-22 (implement)** Implement server-authoritative
  `Entry.Author`. Identical mechanism to H-04 (Origin
  enforcement): `validateBearer` attaches `ClientInfo`
  to the gRPC context; `handler.go publish()` reads
  `ClientInfo` and stamps `entries[i].Author` from the
  server-known identity before calling `store.Append`.
  Pre-registry the stamping source is
  `ClientInfo.ProjectName`; after the registry MVP the
  source becomes `users.json` row's `user_id`; after
  the PKI stretch it becomes the signed-claim `sub`.
  Same commit as H-04 is fine — they share the
  `authFromContext` plumbing. Add a test that a client
  authenticated as project `alpha` cannot publish an
  entry whose stored `Author` differs from `alpha`.
  Audit client-side callers in `ctx connect publish`
  and `ctx add --share` for any that populate
  `pe.Author` from local config and remove them (or
  document them as ignored). #priority:high
  #added:2026-04-11 #pr:60 #audit:H-22
- [ ] **H-22a (server-authoritative Origin stamping)**
  Implement H-04-style server-enforcement for
  `Entry.Origin`: `validateBearer` attaches
  `ClientInfo` to the gRPC context;
  `handler.go publish()` reads `ClientInfo` and
  overwrites `entries[i].Origin` with
  `ClientInfo.ProjectName` before `store.Append`.
  Client's `pe.Origin` becomes advisory and is
  ignored. This is the actual security property
  the Author→Meta split was enabling — the
  schema change made room for it but the
  enforcement still needs to land. Add a test:
  client authenticated as `alpha` cannot publish
  an entry whose stored Origin is `beta`.
  #priority:high #added:2026-04-11 #pr:60 #audit:H-22
- [ ] **H-22b (renderer labels Meta as advisory)**
  Update `internal/cli/connect/core/render/` (and any
  other place that writes fanned-out entries to
  `.context/hub/*.md`) so `Meta`-sourced values are
  labeled as "client label" or "client-reported" in
  prose. The word "Origin" is reserved for the
  server-authoritative project name. Example output:

  ```markdown
  ## [2026-04-11] Use UTC timestamps everywhere
  **Origin**: alpha (client label: Alice via ctx@0.8.1)
  ```

  Add a test verifying that a Meta.DisplayName of
  `"bob"` does NOT cause the rendered output to show
  `Origin: bob`. #priority:high #added:2026-04-11
  #pr:60 #audit:H-22
- [ ] **H-22c (client publish path supports Meta)**
  Update `ctx connect publish` (and `ctx add --share`
  if it reaches the hub) to accept `--display-name`,
  `--host`, `--tool`, `--via` flags (or a single
  `--meta key=val` repeatable flag — implementation
  choice). Defaults: `--tool=ctx@<version>`,
  `--host=<hostname>`, `--via=` left empty,
  `--display-name=` left empty. Document in
  `docs/cli/connect.md`. #priority:medium
  #added:2026-04-11 #pr:60 #audit:H-22
- [ ] **H-22d (docs: `Meta` is advisory)** Add a
  prominent note to `docs/cli/connect.md`,
  `docs/security/hub.md`, and
  `docs/recipes/hub-overview.md` explaining that
  `Meta` fields are client-reported hints, not
  attribution. Cross-reference the decision record
  [2026-04-11-180000]. #added:2026-04-11 #pr:60
  #audit:H-22
- [ ] **H-22e (audit spec update)** Update
  `specs/hub-security-audit.md` H-22 finding to
  reflect the landed schema change: the "decide"
  phase is done, the "meta type" phase is done, the
  remaining work is the Origin stamping (a), the
  renderer labels (b), and the client-side plumbing
  (c). Also note the six regression tests as "partial
  coverage" of the finding. #added:2026-04-11 #pr:60
  #audit:H-22
- [ ] **H-30** gRPC server hardening: `KeepaliveEnforcementPolicy`,
  `KeepaliveParams`, `MaxConcurrentStreams`, total
  concurrent connection limit at the listener level.
  #priority:medium #added:2026-04-11 #pr:60 #audit:H-30

**Medium-term track** (correctness + cluster integrity):

- [ ] **H-12** Deterministic Raft bootstrap: single
  `--bootstrap` node calls `BootstrapCluster`, others
  join via `AddVoter`. Persist a `bootstrapped` flag
  in the raft data dir to avoid double-bootstrapping
  on restart. #priority:medium #added:2026-04-11 #pr:60
  #audit:H-12
- [ ] **H-13** Follower-side replication validation:
  call `validateEntry` on every entry received from
  master before appending. Defense-in-depth against a
  compromised master (which becomes possible under any
  Raft transport compromise — see H-10/H-11).
  #priority:medium #added:2026-04-11 #pr:60 #audit:H-13
- [ ] **H-14** Preserve master sequence on replication:
  add `masterSequence` field to Entry, followers
  remember master-assigned sequences alongside local
  ones. Clients cursor by master sequence so failover
  doesn't re-replicate the entire log. #priority:medium
  #added:2026-04-11 #pr:60 #audit:H-14
- [ ] **H-24** `ctx hub redact <seq>` subcommand: mark
  the entry in `entries_redacted.jsonl`, broadcast a
  redaction notice via Listen, filter on queries, log
  to `audits.jsonl`. #priority:medium #added:2026-04-11
  #pr:60 #audit:H-24
- [ ] **H-29** Bounded in-memory entry cache: LRU over
  `entries.jsonl` with a persistent offset index
  (`entries.idx`). O(log N) seeks without full-file
  reads. Secondary: entries.jsonl rotation at threshold.
  #priority:medium #added:2026-04-11 #pr:60 #audit:H-29

**Long-term track** (Story 3 enablement):

- [ ] **H-10 + H-11** Authenticated + encrypted Raft
  transport. Replace `raft.NewTCPTransport` with a
  TLS-wrapped transport using mTLS between cluster
  peers. Peer certs issued from a cluster CA managed
  by the sysadmin. Precondition for any non-localhost
  multi-node deployment. #priority:critical
  #added:2026-04-11 #pr:60 #audit:H-10,H-11
- [ ] **H-28** Decouple Raft bind port from gRPC port.
  Accept a dedicated `--raft-bind` flag; default to a
  random high port or refuse to start. Makes port
  scanning less productive. #priority:low
  #added:2026-04-11 #pr:60 #audit:H-28
- [ ] Signed-entry mode: publishing clients sign their
  entries with a per-client signing key; followers
  verify on replication. Eliminates the "trust the
  master" assumption even if H-10 fails. Merged with
  the PKI stretch task in the Hub identity layer
  phase. #added:2026-04-11 #pr:60 #audit:H-13

**Low-priority polish** (defense-in-depth):

- [ ] **H-16** Escape / fence `Content` when the
  client-side renderer writes to `.context/hub/*.md`.
  Wrap every entry in explicit markers
  (`<!-- BEGIN ENTRY seq=... -->`) so malicious
  triple-dash patterns can't inject fake frontmatter.
  #added:2026-04-11 #pr:60 #audit:H-16
- [ ] **H-20** Strict constant-time token validation:
  iterate all `ClientInfo` entries and OR the results
  of `subtle.ConstantTimeCompare` instead of a map
  lookup followed by a constant-time compare. Rolled
  into the H-03 hashing work. #added:2026-04-11 #pr:60
  #audit:H-20
- [ ] **H-21** Require exact `Bearer ` prefix in the
  `authorization` header; reject otherwise with
  `Unauthenticated`. Trivial one-line tightening.
  #added:2026-04-11 #pr:60 #audit:H-21
- [ ] **H-23** Offer passphrase-derived admin token
  storage (argon2id) instead of plaintext `admin.token`
  on disk. Optional; document in
  `docs/operations/hub.md`. #added:2026-04-11 #pr:60
  #audit:H-23
- [ ] **H-25** Collapse auth error messages to a single
  generic `Unauthenticated` reason ("authentication
  required"). Log the specific reason server-side
  only. #added:2026-04-11 #pr:60 #audit:H-25

**Informational (no action needed)**:

- H-26: daemon re-exec flag — already fixed earlier in
  this session as part of the `ctx serve --hub` → `ctx
  hub start` split. Recorded in the audit for audit-
  trail completeness.
- H-27: mTLS / asymmetric auth discussion — covered by
  the PKI stretch task in the Hub identity layer
  phase. No separate task needed.

**Out of scope for this audit** (tracked elsewhere):

- Supply chain (Go module pinning, CVE monitoring,
  reproducible builds)
- Build integrity (signed binaries, transparency log)
- Third-party library CVEs (`hashicorp/raft`, `grpc`,
  `raft-boltdb`)
- AI-agent misbehavior (accidental secret publishing
  via `--share` — covered by the "secret-leak runbook"
  task in the PR #60 follow-up section above)
- Per-project read ACLs (still out of scope even after
  the identity layer MVP)

#### Rename "Shared Context Hub" → "`ctx` Hub" (2026-04-11)

Brainstorm outcome: "shared" was overloaded (shared memory,
shared journal, shared state) and actively primed the wrong
mental model in docs. `ctx` Hub is the canonical name; `Hub` is
used alone in nav and operator contexts where surrounding text
disambiguates.

### Later

- [ ] Optional follow-up doc.go pass: a handful of tiny per-subcommand wrappers
  under internal/cli/*/cmd/* still have ~5-line bodies. Most are
  accurate-but-brief; expand only if the brief form proves insufficient in
  review. #session:4b37e2f6 #branch:feat/copilot-cli-skill-parity-rebased
  #commit:edaac81786c9379333b352dae0d55df0ae0f72bb #added:2026-04-14-010311

- [ ] Extend internal/audit/stuttery_functions_test.go to cover *ast.GenDecl
  (consts, vars, types). Current implementation walks *ast.FuncDecl only and
  missed tpl.TplEntryMarkdown (since renamed to HubEntryMarkdown).
  #session:4b37e2f6 #branch:feat/copilot-cli-skill-parity-rebased
  #commit:edaac81786c9379333b352dae0d55df0ae0f72bb #added:2026-04-14-010311

- [ ] Decide whether to delete docs/cli/connect.md — verified dead duplicate
  of docs/cli/connection.md (uses old ctx connect command name; zero inbound
  references; not in zensical.toml). Awaiting explicit user OK before git rm.
  #session:4b37e2f6 #branch:feat/copilot-cli-skill-parity-rebased
  #commit:edaac81786c9379333b352dae0d55df0ae0f72bb #added:2026-04-14-010311

- [-] PROMPT.md design — belongs in another project; skipped here.
  #session:4b37e2f6 #added:2026-04-14-010311 #skipped:2026-04-14

### Phase CP: Ceremony Profiles `#priority:medium #added:2026-04-26`

Spec: `specs/ceremony-profiles.md`

- [ ] Add `Ceremony{Remember,WrapUp}` to `internal/rc/types.go`; apply defaults
  in `internal/rc/rc.go` from
  `internal/config/ceremony/ceremony.go` constants
- [ ] Thread resolved ceremony names into `ScanJournalsForCeremonies` and `Emit`
  in
  `internal/cli/system/core/ceremony/ceremony.go` (replace direct constant
  reads)
- [ ] Convert
  `internal/assets/hooks/messages/check-ceremony/{remember,wrapup,both}.txt` to
  `{REMEMBER}` / `{WRAPUP}`
  sentinels; audit `internal/config/embed/text` ceremony desc keys for the same
- [ ] Add a single sentinel-substitution helper (extend
  `internal/cli/system/core/message.Load` or sibling) so
  substitution happens in one place
- [ ] Show active ceremony profile (one line) in `ctx status` output
- [ ] Tests: default profile renders `/ctx-remember` `/ctx-wrap-up`; project
  with `ceremony.remember: dp-remember`
  renders `/dp-remember` and scanner only counts `dp-remember` as fulfilling the
  open-bookend
- [ ] Document in `docs/recipes/` with the editorial-project (`your-domain`
  knowledgebase) consumer as the worked example

### Phase SK: Skill Surface Polish (Phase 0a; prerequisite for Phase KB) `#priority:high #added:2026-05-09`

Spec: `specs/skill-surface-polish.md` (design ref:
`ideas/002-editorial-pipeline-and-skill-rigor.md` §3 "Reframing the
wishy-washy skills")

Tightens existing capture skills to sibling-project rigor before the editorial
pipeline (Phase KB) lifts that pattern
wholesale. Independent of Phase RG; both can ship in parallel.

- [x] Add `MarkFlagRequired` to `ctx decision add` for `--context`,
  `--rationale`, `--consequence`; reject placeholder
  values (`TBD`, `see chat`, whitespace-only) at CLI level
- [x] Add `MarkFlagRequired` to `ctx learning add` for `--context`, `--lesson`,
  `--application`; same placeholder
  rejection
- [x] Add `--brief <path>` flag to `/ctx-spec` skill: when present, read the
  file as authoritative source per the
  sibling's authority order (frozen contracts > recorded decisions > debrief >
  agent inference labeled `TBD`); skip the
  fresh template Q&A
- [x] Update `/ctx-plan` skill to always offer to write the debated brief to
  `.context/briefs/<TS>-<slug>.md` at the end
  of an interview (creating `.context/briefs/` if absent)
- [x] Add an `Authority boundary (vs other skills)` section to
  `/ctx-decision-add`, `/ctx-learning-add`,
  `/ctx-task-add`, `/ctx-convention-add` skill files (prevent silent promotion
  handover→decision, learning→convention,
  etc., without explicit user ask)
- [x] Standardize "light compression for clarity is allowed; new facts are not"
  wording across capture skills (decide /
  learn primarily); same wording lands in `/ctx-handover` once Phase KB ships
- [x] Document the `--brief` contract in `docs/skills.md` (landed in
  `docs/reference/skills.md`; the actual location)

### Phase RG: Require Git as Architectural Precondition (Phase 0b; prerequisite for Phase KB)
`#priority:high #added:2026-05-09`

Spec: `specs/require-git.md`

Enforces what `ctx` already needs: git. `ctx` works properly only with a
repo present, and this phase makes that a runtime precondition rather than
an assumption. Breaking change for any pre-existing git-less ctx project
(N≈0 in practice). Independent of Phase SK; both can ship in parallel.

- [x] Add `internal/gitmeta/require.go` with `RequireGitTree(projectRoot string)
  error` and typed `MissingGitError`
- [x] Wire `RequireGitTree` into root command PersistentPreRunE; opt-out list
  (via the existing
  `AnnotationSkipInit` mechanism that already covers `--help`, `--version`, `ctx
  system bootstrap`, init,
  activate, deactivate, guide, why, doctor, config switch/status, hub *)
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
- [x] Add CONSTITUTION.md amendment ("Git is required") under Process Invariants
- [x] Add DECISIONS.md entry: "Mandate git as architectural precondition"
  (Accepted; context = LLM-safety + provenance
  honesty + dead-code elimination; consequence = breaking change for
  pre-existing git-less projects, N≈0).
  Filed as DECISIONS.md "Phase KB lifts the current upstream editorial-pipeline
  shape, superseding the 4-phase predecessor"
  (2026-05-16) which folds the git-mandate context into Phase KB's parent
  decision.
- [ ] Update `docs/recipes/bootstrap-a-project.md`, `README.md`,
  `docs/cli/init.md` to show `git init` before `ctx init`
- [ ] Tag as breaking change in `dist/RELEASE_NOTES.md` with one-command
  migration ("Run `git init` in any pre-existing
  git-less ctx projects before upgrading")
- [x] Tests: `.git` dir → nil; `.git` file (worktree pointer) → nil; absent
  → typed error
- [-] Tests: root PreRunE refuses without git; opt-out list allowed. TBD:
  deferred to Phase KB Stage 5 when
  the kb command tree is in place (the existing bootstrap_test.go covers PreRunE
  structure; the gitmeta
  injection's behavioral test runs as part of the kb-ingest smoke)
- [-] Compliance test: no remaining `commit:none` literal in `internal/`. N/A:
  literal never existed

### Phase KB: Editorial Pipeline + Handover (depends on Phase SK + Phase RG) `#priority:high #added:2026-05-09 #revised:2026-05-16`

Spec: `specs/kb-editorial-pipeline.md` (revised 2026-05-16 to current
upstream editorial-pipeline shape: pass-mode contract, completion circuit
breaker, source-coverage state-machine ledger, topic-adjacency
pre-flight, cold-reader rubric, folder-shaped topics from day one).

Comparison input: `ideas/upstream-pipeline-comparison.md`.

Decision record: DECISIONS.md "Phase KB lifts the current
upstream editorial-pipeline shape, superseding the 4-phase predecessor in the
brief" (2026-05-16).

Brief: `ideas/003-editorial-pipeline-debated-brief.md`

Background analysis: `ideas/001-sibling-project-undercover-analysis.md`,
`ideas/002-editorial-pipeline-and-skill-rigor.md`

Validation corpus: `your-project` (live regression
suite; hand-rolled the older 4-phase shape for weeks).
`your-project` is the structural reference for the current
upstream shape applied to a different domain.

Note on task lines below: path-constant locations were originally
specified as `internal/path/path.go`. The revised spec places them
under `internal/cli/kb/core/path/path.go` to match existing ctx
convention (per-subcommand path package, see `internal/cli/task/core/path/`).
Similarly the "store layer" tasks below land under `internal/write/`
(handover, closeout, kb), not `internal/store/`. Task wording kept
historical for traceability; implementation follows the revised spec.

Path constants and embedded templates:

- [x] Constants landed at `internal/config/kb/kb.go` (filenames + subdir names +
  state-machine constants + pass-mode + life-stage) and
  `internal/cli/kb/core/path/path.go` (full-path resolvers: KBDir, KBTopicDir,
  IngestDir, CloseoutsDir, SchemasDir, HandoversDir, ArchiveCloseoutsDir,
  SiteDir, SiteKBDir). Per the per-subcommand convention; not `internal/path/`.
- [x] Templates embedded under `internal/assets/kb/templates/ingest/`:
  KB-RULES.md, 00-GROUND.md, 30-INGEST.md, 40-ASK.md, 50-SITE_REVIEW.md,
  OPERATOR.md, PROMPT.md. INBOX.md and SESSION_LOG.md not pre-seeded; they
  materialise on first skill run (per the contract that skills are the sole
  writers of INBOX.md and SESSION_LOG.md only appears mid-flight).
- [x] Schemas embedded under `internal/assets/kb/templates/ingest/schemas/`:
  evidence-index.md, glossary.md, contradictions.md, outstanding-questions.md,
  domain-decisions.md, timeline.md, source-map.md, source-coverage.md,
  relationship-map.md, session-log.md (10 files; fields + one worked example
  each)
- [x] `internal/assets/embed.go` extended with `//go:embed` lines for the kb
  tree

Store layer (landed under `internal/write/` per the revised spec, not
`internal/store/`):

- [x] `internal/write/handover/` (WriteHandover, Latest, fold via
  closeout.PostdatedBy, archive via closeout.Archive)
- [x] `internal/write/closeout/` (Write, Read, List, PostdatedBy, Archive;
  required frontmatter sha/branch/mode/pass-mode/life-stage/generated-at;
  ErrMissingFrontmatter, ErrMissingFields)
- [x] `internal/write/kb/` split across 9 subpackages, not a single kb.go:
  evidence (no-renumber, ID allocation, ErrDuplicateID/ErrInvalidBand),
  sourcecoverage (state-machine ledger with ValidTransition +
  ErrIllegalTransition + ErrUnknownSource), glossary, contradiction, question,
  decision, timeline, sourcemap, relationship

CLI commands:

- [x] `ctx handover write`: `MarkFlagRequired` on `--summary` and `--next`;
  rejects placeholder values (TBD, see chat, n/a, none); calls
  `internal/write/handover.Write` which folds postdated closeouts (`--no-fold`
  skips fold); supports `--commit`, `--highlights`, `--open-questions`;
  smoke-tested end-to-end
- [x] `ctx kb` parent command + `topic new` (real scaffold writer), `note` (real
  append to findings.md), `reindex` (real CTX:KB:TOPICS managed-block refresh in
  kb/index.md), `ingest`/`ask`/`site-review`/`ground` (skill-driven;
  refuse-on-empty for ingest/ask/ground); smoke-tested
- [x] KB rendering: `.context/kb/` is a tree of Markdown that
  `ctx serve` already serves once a `zensical.toml` is dropped in. Recipe
  Step 5 documents the path. A thin `ctx kb site (build|serve|customize)`
  wrapper that mirrors `internal/cli/journal/cmd/site/` and pre-seeds the
  `zensical.toml` is a follow-up convenience, not a blocker.
- [x] `ctx init` scaffolds `.context/kb/`, `.context/kb/topics/`,
  `.context/ingest/` (with embedded templates copied),
  `.context/ingest/closeouts/`, `.context/ingest/schemas/` (10 schemas copied),
  `.context/handovers/`. Implemented in new `internal/cli/initialize/core/kb/`
  package called from init's run.go. `--upgrade` flag: not added; init's
  existing skip-existing-files behavior is idempotent on byte-identical content
  (the divergence-refusal needs a separate `--upgrade` follow-up).

Skills:

- [x] 6 new SKILL.md files: ctx-handover (280L), ctx-kb-ingest (645L),
  ctx-kb-ask (236L), ctx-kb-site-review (259L), ctx-kb-ground (279L),
  ctx-kb-note (164L)
- [x] Modified `ctx-wrap-up/SKILL.md`: branches on `.context/kb/` existence
  (surfaces pending closeouts + outstanding-questions count); mandatorily drives
  `/ctx-handover` as final step
- [x] Modified `ctx-remember/SKILL.md`: reads latest handover + postdated
  unfolded closeouts; folds KB state into readback when `.context/kb/` exists

Doctor / status / .gitignore:

- [-] Doctor advisories: NOT YET IMPLEMENTED. Spec lists duplicate-`EV-###`,
  `dated:`-source-missing-`occurred:`, malformed-closeout-frontmatter,
  source-coverage-ledger-mismatch (row Updated vs. file mtime),
  closeout-missing-pass-mode-body-block, illegal-ledger-state-transition. Phase
  7 follow-up.
- [-] Mode-aware reads in ctx status / ctx agent / session-start hook: skills
  updated (`/ctx-remember` + `/ctx-wrap-up`); CLI-side `ctx status`/`ctx agent`
  mode-awareness deferred (the skill-side fold covers the user-facing recall;
  CLI text surfaces are v1.1).
- [x] `.gitignore` extended: `.context/site/`, `.context/site-config/`

Tests:

- [ ] Unit tests per package (handover, closeout, kb writers, mode CLIs, doctor
  advisories)
- [ ] Integration: `internal/cli/initcmd/init_test.go` covers full new directory
  tree + `--upgrade` idempotency /
  divergence refusal
- [ ] `hack/smoke-kb.sh`: end-to-end shell smoke (init → kb ingest → kb ask
  → kb site-review → kb ground → handover
  write → archive populated → doctor clean)
- [ ] Edge-case fixtures: aborted-session recovery (closeout without handover);
  temporal misordering (
  occurred-vs-extracted ordering enforces precedence rule); concurrent dupe IDs
  (LLM-resolution fixture); render
  filter (speculative excluded; low paired with outstanding-questions)

Phase KB-2 (validation against live corpus):

- [ ] Port `your-project-*` from its hand-rolled shape to the shipped one. Each
  divergence is either a
  Phase KB bug or a `DECISIONS.md` entry explaining why the formal shape differs
  from what worked manually
- [ ] Document divergences (if any) in `docs/recipes/build-a-knowledge-base.md`

Phase KB-3 (documentation):

- [x] Write `docs/recipes/build-a-knowledge-base.md` (mirrors sibling's recipe
  shape)
- [x] Write `docs/recipes/typical-kb-session.md`
- [x] Write `docs/recipes/recover-aborted-session.md`
- [x] Update CLI reference with new `ctx kb` and `ctx handover` commands (landed
  as separate pages: `docs/cli/kb.md` for the editorial pipeline and
  `docs/cli/handover.md` for the session-glue command; the two are unrelated
  surfaces and now have distinct pages)
- [x] Update `docs/reference/skills.md` with the 6 new skills (table row +
  per-skill section + new "Knowledge Base (Phase KB)" section)
- [x] Update root `README.md` with the Phase KB workflow snippet + git-required
  note
- [x] Update root `CLAUDE.md` and `internal/assets/claude/CLAUDE.md` (the
  user-deployed copy) with the KB trigger table
- [x] Update `dist/RELEASE_NOTES.md` with Phase RG + Phase KB sections
- [ ] Document MemPalace-as-ground-source recipe in
  `docs/recipes/build-a-knowledge-base.md`; uses already-specced
  `mcp:<server>:<resource>` syntax in `grounding-sources.md`; zero new ctx code

### Phase KB-followup: Adversarial design review of parallel skill trees `#priority:medium #added:2026-05-17`

`ctx` ships skills to three host trees:
`internal/assets/claude/skills/` (canonical, full Claude tool surface),
`internal/assets/integrations/copilot-cli/skills/` (Copilot CLI; `tools: [bash]`),
and `internal/assets/integrations/opencode/skills/` (OpenCode; minimal
subset, no `tools` block). Phase KB landed parity across all three trees
by writing each new skill body three times (full content for Claude +
Copilot CLI; terser variant for OpenCode), which works today but
guarantees future drift the next time a canonical skill is revised.

Run an **adversarial design review** to pick the right architecture for
preventing this drift permanently. Candidate shapes:

- **Body-extract + per-host frontmatter wrapper at build time.** Single
  source of truth for behavioral prose; a builder package composes
  host-specific SKILL.md files with the right frontmatter and the
  right capture-skill name swaps (`/ctx-task-add` vs `/ctx-add-task`,
  etc.) at `go generate` or `make build` time. Per-host overrides
  for genuinely different host capabilities live in side files.
- **Write canonical, copy at runtime, make integration trees
  read-only.** Simpler builder; risk is that host-specific tool
  surfaces leak (Claude has Edit/Write/Read; Copilot CLI has bash
  only; OpenCode is more constrained).
- **Convention-only with audit gate.** Keep three independent trees
  but add an audit test that fails CI when a canonical-tree skill
  changes without parallel changes in the integration trees. Cheaper
  but pushes the work onto contributors.
- **Drop one or both integration trees.** OpenCode currently ships
  only a 4-skill subset that the user may or may not want at parity.
  Decide explicitly which trees are first-class.

Deliverables:

- [ ] Adversarial review write-up under `ideas/` enumerating each
  shape with pros / cons / migration cost.
- [ ] DECISIONS.md entry picking the shape, with rationale.
- [ ] Implementation tasks for the chosen shape.
- [ ] A compliance test that fails when the canonical Claude tree
  changes a Phase KB or handover skill without the parallel tree
  being updated, until the builder lands.

Context: filed after Phase KB shipped, when porting the 6 new KB
skills + the 2 updated ceremony skills to copilot-cli and opencode
revealed how brittle the three-tree pattern is.

### Phase JR: Cold-Start Memory Recovery (semantic recall over journal history) `#priority:medium #added:2026-05-10`

Idea: `ideas/004-cold-start-memory-recovery.md`

Pain point: today's "can you check recent journal entries?" workaround forces
brute-force parsing of the journal corpus
or precise user pointers to specific files/dates. ctx has journal management but
no semantic recall layer.
MemPalace (https://github.com/MemPalace/mempalace) does this exact use case at
96.6% R@5 raw on LongMemEval. Three
options to evaluate: A) native ctx journal search (vector-store dep, breaks
single-Go-binary identity); B)
defer-to-MemPalace recipe (zero ctx-side work; coupling to young project); C)
pluggable journal-search hook following
the zensical shell-out pattern (recommended).

- [ ] Spec out cold-start memory recovery: pick approach (A vs B vs C);
  ideas/004 leans toward C. Distinct from Phase KB
  ground-mode `mcp:` source kinds (which cover the KB-grounding angle for free);
  this phase is specifically about
  journal-corpus semantic recall (`ctx journal search "<query>"` shape).
