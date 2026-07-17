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

- [ ] The target project (to be given to the Agent) has a good "phasing"
  mechanism for tasks; implement that; maybe `ctx task add` can have a
  `--phase` flag too, and we can have a auditor/normalizer for the current
  task document; or a skill that does a semantic pass, or both too.

## Phase CLI-FIX: CLI Infrastructure Fixes

These have priority because other knowledge ingestion projects depend on them.

- [x] BUG (data loss): `ctx learning add` clobbers a dash-bullet-format
  `LEARNINGS.md`. When the `INDEX:START/END` block uses the dash-bullet format
  (what `ctx init` produces) rather than the pipe-table format, `ctx learning
  add` (1) rewrites the index as a table, (2) **duplicates** the
  `<!-- INDEX:START -->` marker, and (3) **drops every existing learning body**
  —
  keeping only the newly added one. Observed live in `things-wtf-hub` (session
  aa32f065): a `LEARNINGS.md` with 4 bodies collapsed to 1 (a -44-line commit,
  2dc4d1a); recovered via `git show <good-sha>:.context/LEARNINGS.md`.
  `ctx decision add` is UNAFFECTED because that repo's `DECISIONS.md` was
  already
  table-format — so the bug is specifically the learning-add path's handling
  of
  the dash-bullet index variant (likely it can't parse dash-bullet entries, so
  it treats the file as empty and regenerates from only the new entry).
    - Repro: `ctx init` a repo, confirm `LEARNINGS.md` index is dash-bullet
      (`- entry`), add 2+ learnings by hand in that format, then run
      `ctx learning add "x" --context … --lesson … --application …`.
      Expect: the
      hand-authored bodies vanish + a duplicated INDEX:START marker.
    - Likely fix: detect the existing index format (dash-bullet vs table) and
      preserve it round-trip, OR parse dash-bullet entries before regenerating;
      never emit a second INDEX:START; never drop bodies the parser didn't
      recognize (fail loud instead of silently regenerating).
    - Guard: a round-trip test for BOTH index formats (dash-bullet + table) that
      asserts existing bodies survive an add and exactly one marker pair
      remains.
- Severity: HIGH — silent destruction of persisted memory, the one thing ctx
      promises to protect; only git made it recoverable.
    - Provenance: things-wtf-hub session aa32f065 wrap-up; full write-up in that
      repo's LEARNINGS.md ("`ctx learning add` clobbers a dash-bullet-format
      LEARNINGS.md"). #priority:high #added:2026-05-30
      #completed:2026-06-01 #branch:fix/learning-add-index-data-loss
      Shipped: new `index.Validate` precondition guard refuses to regenerate the
      index when it would lose data — entry bodies (`## [ts]` headers) trapped
      between the markers, or markers that are missing/duplicated/out-of-order.
      Wired into the two read-before-mutate choke points (`entry.Write` and
      `index.Reindex`), so add and all reindex commands fail loud and leave the
      file byte-identical instead of clobbering it. `index.Update`'s signature
      is
      untouched (kept the CRITICAL blast radius stable). New
      `internal/err/index`
      package + i18n messages. Verified: the real repro now errors with the file
      unchanged; well-formed adds still preserve all bodies and one marker pair.
      Tests: `index.TestValidate` (7 shapes) + `entry` round-trip
      (refused-untouched + well-formed-preserved). Chosen behavior is fail-loud
      +
      manual fix; auto-repair (`reindex --repair`) considered and declined.
      Spec: specs/fix-learning-add-index-data-loss.md.

- [x] Make 'ctx kb reindex' nesting-aware: scan topics/** not topics/* (grouped
  topic folders currently blank the CTX:
  KB:TOPICS block) #priority:medium #session:c3d2dcb1 #branch:
  feat/pad-undo-snapshot #commit:b9ce72e8 #added:
  2026-05-27-182640 #completed:2026-05-28
  Shipped: `ListTopics` now recurses (topic.go + scan.go),
  enumerating `topics/<group>/<slug>` as slashed slugs and excluding
  group-landing pages (a dir whose index.md sits above nested
  topics). Flat / grouped / mixed / arbitrary-depth layouts all
  enumerate; a non-existent dir still yields the placeholder, never
  an error-blank. `RenderBlock` unchanged (the `topics/<slug>/`
  template already renders nested links; sorted slashed slugs cluster
  by group prefix). Tests: topic_test.go (7 cases + nonexistent),
  block_test.go (nested-slug + empty). Per-group headings deferred
  (managed-block format change). Spec: specs/kb-reindex-nesting.md.
    - Problem: `ctx kb reindex` scans `topics/*/index.md` (one level). A
      consumer kb (the DR project, things-wtf-dr)
      reorganized 49 topics into grouped folders
      `topics/<group>/<slug>/index.md`; reindex then finds 0 topics and
      BLANKS the `CTX:KB:TOPICS` managed block in `index.md` (observed live: "
      reindexed 0 topic(s)"). The same one-level
      assumption likely affects the life-stage topic-count glob (
      `topics/*/index.md`) and any other `topics/*/`
      enumeration.
    - Fix: scan `topics/**/index.md` recursively; exclude group-landing pages
      `topics/<group>/index.md` from topic
      enumeration (orientation, not topic pages); ideally emit the managed block
      grouped by parent folder.
      `ctx kb topic new "<group>/<slug>"` already preserves nested slugs, so
      creation is unaffected — only
      reindex/enumeration lags.

- [x] Add `--json <file>` to `ctx decision/learning/task add` (and `convention`
  if it gains structured fields) for
  ingesting a JSON payload that populates the typed fields directly.
    - Driver: this session hit a class of denial we worked around but should fix
      at the root. The project's canonical
      `permissions.deny` set (`.claude/settings.local.json` lines 119-121)
      matches on the literal Bash command string —
      including the *content* of `--rationale`/`--context`/`--consequence` flag
      values. A decision whose rationale
      legitimately describes installing a binary into the system PATH (literal
      substring " /usr/local/bin") gets caught
      by `Bash(* /usr/local/bin*)` and denied, even though the command's intent
      has nothing to do with that path. The
      workaround was Edit-direct into DECISIONS.md/LEARNINGS.md, which bypasses
      the ctx command's schema gates and
      INDEX:START/END maintenance.
    - Shape: `ctx decision add --json /path/to/payload.json` where the JSON is
      `{"title":"…","context":"…","rationale":"…","consequence":"…"}`.
      The flag
      supersedes individual content flags.
      Provenance (--session-id/--branch/--commit) can stay on the command line
      OR be folded into the JSON envelope ({"
      provenance":{"session_id":"…","branch":"…","commit":"…"}}).
      Complements
      the existing `--file` (which only replaces
      the title/body positional).
- Phase 2 (optional): array form `[{...},{...}]` for batch persists — useful
      for `/ctx-wrap-up` writing N
      decisions+learnings in one call instead of N separate invocations.
    - Mirror per command: same shape applies to `ctx learning add --json …` (
      {title,context,lesson,application}) and
      `ctx task add --json …` ({title,body,priority,section}).
    - Surfaced by: this session's persist denials and post-mortem; reference
      handover
      `20260528T201500Z-ctxctl-and-native-pressure-shipped.md`. #priority:medium
      #session:96765858 #branch:
      feat/pad-undo-snapshot #commit:b9ce72e8 #added:2026-05-28-154725

- [ ] Exploratory: Windows-native memory-pressure detection for the
  `check-resource` hook. macOS (
  `kern.memorystatus_vm_pressure_level`) + Linux (PSI `/proc/pressure/memory`)
  native pressure detection landed on
  feat/pad-undo-snapshot, replacing the broken occupancy-% triggers. Windows ("
  other" platform) currently reports
  `PressureSupported=false` → no memory alert.
    - Explore the Windows-native signal: Memory Resource Notifications API (
      `CreateMemoryResourceNotification`/
      `QueryMemoryResourceNotification` → `LowMemoryResourceNotification`),
      perf
      counters (`Memory\Available MBytes`,
      `Committed Bytes`/`Commit Limit`), or `GlobalMemoryStatusEx.dwMemoryLoad`.
    - Open question: Windows aggressively manages working-set/commit and
      surfaces its own low-memory UI, so it likely
      warns the user before ctx can — assess whether a ctx-side signal adds
      value at all before building it.
    - Wire into a build-tagged `internal/sysinfo/memory_windows.go` (currently
      falls through to memory_other.go).
      Provenance: session 96765858; design context in this session's
      swap-detection thread. #priority:medium #session:
      96765858 #branch:feat/pad-undo-snapshot #commit:b9ce72e8 #added:
      2026-05-27-183909

- [x] Realign the installed plugin's hooks.json with the cwd-anchored binary —
  the LIVE fix for the every-prompt
  help-dump pollution.
    - Problem: the cwd-anchored migration (commit fc7db228, spec
      specs/cwd-anchored-context.md) is UNRELEASED — not in
      any 0.8.x tag (only v0.8.0 exists). The installed plugin (~
      /.claude/plugins/cache/activememory-ctx/ctx/0.8.1/hooks/hooks.json) is
      PRE-migration: it injects `CTX_DIR=` and
      wires `ctx system check-anchor-drift` first under UserPromptSubmit. The
      on-PATH binary (0.8.1) is POST-migration:
      check-anchor-drift deleted, cwd-anchored. So the shipped hooks.json calls
      a command the binary no longer has →
      cobra prints the full `system` help and exits 0 → ~52 lines injected on
      EVERY prompt, labelled "hook success".
    - Fix: cut/republish the plugin so its bundled hooks.json comes from the
      same post-fc7db228 commit as the binary (
      cd-based invocation, no check-anchor-drift, includes check-audit).
      Reinstall/update locally and for any users on
      the skewed 0.8.1 plugin.
    - Recurrence guard (acceptance): add a release-time check that every
      `ctx system <verb>` wired in the shipped
      hooks.json resolves to a registered subcommand on the shipped binary (test
      or hack/release.sh step). A
      half-migrated package must not ship again. Pairs with the verbatim-relay
      guard task above — that one makes a
      future skew fail LOUD; this one closes the current gap.
- #in-progress 2026-05-28 (branch feat/hooks-wiring-guard, session 0066d49b):
      Recurrence guard SHIPPED — `TestShippedHooksResolveToRegisteredCommands`
      in internal/compliance walks every `ctx <…>` invocation in the shipped
      hooks.json against the assembled cobra tree; a wired-but-unregistered verb
      fails `go test`. Proven both ways (passes clean, fails on a reintroduced
      `check-anchor-drift`). Spec: specs/hooks-wiring-guard.md. Implemented as a
      Go test, not a hack/release.sh step (cross-platform, no bash, runs in CI).
      STILL OPEN: the live fix — cut/republish a release where plugin
      hooks.json
      and binary share a post-fc7db228 commit, then reinstall for skewed users
      (a tag+publish action, maintainer-owned).
      CORRECTION to the Fix bullet above: shipped hooks must NOT "include
      check-audit" — the existing `TestShippedHooksExcludeCheckAudit` guard
      forbids it (audit channel is maintainer-only, per the ctxctl migration).
      The current asset correctly omits it; the republished package must too.
    - Provenance: check-anchor-drift version-skew investigation. Design notes:
      specs/experiments/acdl-session-start.md (
      §Root Cause, follow-up #1). #priority:high #session:96765858 #branch:
      feat/pad-undo-snapshot #commit:b9ce72e8
      #added:2026-05-27-145715

- [x] `ctx system`: emit a VERBATIM RELAY on unknown subcommand (replace today's
  silent help-dump + exit 0). Scope:
  `ctx system` ONLY. #completed:2026-05-28 #branch:feat/system-unknown-relay
  Shipped: `ctx system <unknown>` now emits a verbatim NudgeBox (via the write
  layer) naming the verb + version-skew hint, best-effort fires the event-log +
  webhook relay (gated on a session ID read TTY-safely from stdin), suppresses
  cobra's help dump, and exits non-zero. Bare `ctx system` and valid subcommands
  unchanged. Handler in internal/cli/system/core/unknown (RunE on system.Cmd()
  only; parent.Cmd untouched). Verified end-to-end against a real build (box +
  EXIT=1). Spec: specs/system-unknown-subcommand-relay.md.
    - Problem: `ctx system <unknown>` prints the full Long help and exits 0 (
      cobra `legacyArgs` only raises "unknown
      command" for the ROOT command, never a non-root group). In a
      UserPromptSubmit hook a non-zero exit alone is
      swallowed by the harness — "loud via exit code" is dead in the water;
      the
      user never sees it.
    - Fix: route unknown `ctx system` subcommands through the existing
      nudge/verbatim-relay path (same mechanism the
      check-* hooks use) so the message actually reaches the user/agent. Name
      the unknown subcommand and hint at the
      likely cause: a hook referencing a command this binary no longer ships (
      version skew between installed plugin
      hooks.json and the on-PATH binary). Then exit non-zero.
    - Scope guard: `ctx system` only. Do NOT change the generic `parent.Cmd` (
      internal/cli/parent/parent.go); other
      groups (`ctx hub`, etc.) keep cobra's default behavior.
    - Tests: `ctx system <bogus>` emits the verbatim relay (assert body content)
      AND exits non-zero; valid subcommands
      unaffected; bare `ctx system` still prints help.
    - Provenance: surfaced by the check-anchor-drift version-skew investigation.
      Design notes:
      specs/experiments/acdl-session-start.md (Root Cause + follow-up #2). Needs
      its own spec before implementation.
      #priority:medium #session:96765858 #branch:feat/pad-undo-snapshot #commit:
      b9ce72e8 #added:2026-05-27-130130
    - DONE 2026-05-28 (branch feat/system-unknown-relay, session 0066d49b).
      Spec: specs/system-unknown-subcommand-relay.md.
      Approach used: add a RunE on system.Cmd() only (legacyArgs lets the
      leftover args reach the group's RunE for non-root); on unknown verb emit a
      message.NudgeBox to stdout, set SilenceUsage (else cobra re-dumps the help
      we're killing), exit non-zero. system is Hidden so RootCmd
      PersistentPreRunE
      early-returns — no context/git preconditions.
      Decisions settled with user: (1) DO fire the event-log + webhook relay leg
      (nudge.Relay), gated on a real session ID read best-effort from stdin via
      session.ReadID (TTY-safe, timeout-guarded → IDUnknown means skip the
      leg);
      (2) scoped to ctx system only, parent.Cmd untouched.
      Follow-up surfaced: ctx hook (and any parent.Cmd group) has the same
      latent
      exit-0-on-unknown behavior — not wired into hooks.json so out of scope
      here;
      capture as its own task if it ever gets hook-wired.

- [x] Generalize the unknown-subcommand guard beyond `ctx system` (deferred from
  the #5 work above). `ctx hook` and any future `parent.Cmd` group still print
  help + exit 0 on an unknown subcommand — the same latent pollution #5 fixed
  for
  `ctx system`. Low priority while no other group is wired into hooks.json; the
  build-time wiring guard (specs/hooks-wiring-guard.md) only checks `ctx system`
  + `ctx agent` today. If a `ctx hook <verb>` ever gets hook-wired, either
  extend
  the guard's coverage or fold a reusable opt-in into `parent.Cmd` (an optional
  unknown-subcommand handler groups opt into). #priority:low #added:2026-05-28
  DONE 2026-05-30 (branch feat/add-json-file-ingest, session 53db2521).
  Rationale refined: the real justification is not the every-prompt
  amplification
  (unique to hooks.json-wired groups) but making CLI drift LOUD — `ctx hook`
  is
  consumed by name from skills/loops (`ctx hook notify|event|pause|...`), and a
  drifted verb silently returns help+exit-0 (agent misreads; for `notify` the
  human is never told). Lifted the handler from `system/core/unknown` into a
  neutral, parameterized `internal/cli/unknown` (Config + HandlerFor); `system`
  and `hook` both opt in via `c.RunE = unknown.HandlerFor(...)`. `ctx hook` is
  user-facing (not Hidden) and previously rode the no-RunE PreRunE exemption, so
  it needed AnnotationSkipInit to stay reachable without an initialized
  context/git (bootstrap regression test added). Did NOT fold into `parent.Cmd`
  (would widen every group's deps). Skill/loop `ctx hook <verb>` build-time
  guard
  left out of scope. Spec: specs/unknown-subcommand-relay-generalization.md.

## Important

Important things that agent (or human) yeeted to the future.

- [x] Nav gap: doc pages absent from zensical.toml's nav are silently
  unreachable from the site sidebar. Discovered + fully fixed 2026-07-06
  (session 7f6de29d, UNCOMMITTED). Swept EVERY docs/ tree, not just
  recipes:
  - Recipes (10): scrutinizing-a-plan + spec-driven-development →
    Knowledge and Tasks; session-changes → Sessions; hook-sequence-
    diagrams → Hooks; state-maintenance → Maintenance; run-the-dream +
    architecture-deep-dive → Agents and Automation; new "Knowledge Base"
    group for build-a-knowledge-base + typical-kb-session + recover-
    aborted-session.
  - CLI: dream → Sessions; handover → Sessions; kb → Context; the real
    subcommand docs event (`ctx hook event`) + message (`ctx hook
    message`) → Runtime by hook; bootstrap (`ctx system bootstrap`) →
    Runtime by system. Verified via bootstrap/group.go that each is a
    live command, not a stale doc.
  - Reference: dream-executor-contract. Operations: backup-strategy.
  - DELETED docs/cli/connect.md — a verified pure-stale duplicate: the
    command's Use string is now "connection" (connect renamed away),
    connection.md (155 lines) fully supersedes it (same 6 subcommands,
    no unique sections, no inbound links). Also fixed connection.md's
    stale `## ctx connect` heading → `ctx connection`.
  - RECURRENCE GUARD: internal/compliance/docs_nav_test.go
    (TestEveryDocPageIsReachableInNav) fails when any non-blog docs/
    page is missing from the nav. Proven both ways (fails on a planted
    orphan naming it; passes clean). blog/ excluded by design (plugin-
    rendered), includes/ excluded (partials).
  Verified: TOML parses; full non-blog orphan sweep returns zero; guard
  + compliance lint green. #priority:low #session:7f6de29d #branch:main
  #commit:a0e5cbf9 #added:2026-07-06

- [x] Adversarial code-review pass over the 2026-07-06 jumbo working
  diff (Phase JI self-healing import + bundles ②③④: strip-gitnexus/
  de-npx, spec-driven-development recipe, .ctxrc config-ification)
  BEFORE committing. Do it at home with a real code-viewing setup —
  remote triple-tmux is too cramped to review a ~48-file diff well,
  and this wants proper human/agent back-and-forth. Run
  `/ctx-code-review` on the working tree. Focus areas: the growth/
  adopt/foreign-edit decision matrix in plan.Import (body-hash-over-
  frontmatter design); the RegenCount-counts-Grown interaction (an
  interactive `ctx journal import --all` without -y now prompts on
  growth — intended? safe but maybe surprising); the new slug.FromTitle
  → rc singleton coupling the ④ agent introduced (FromTitle is no
  longer pure); state v2 migration/adoption on the real corpus; and the
  SessionEnd-hook direct-wiring vs. the spec's original "thin system
  verb". Everything built/tested/lint-clean this session, but nothing
  is committed and nothing has had a second pair of eyes.
  #priority:high #session:7f6de29d #branch:main #commit:a0e5cbf9
  #added:2026-07-06
  DONE 2026-07-06 (session 2cff382a, branch fix/jumbo-diff-review-fixes):
  4-pass adversarial review found 20 findings (3 critical, 4 medium,
  13 low); ALL fixed with regression tests. Highlights: C1 — `ctx
  journal site` rewrote entry bodies without refreshing render_hash
  (self-heal false-flagged every site-normalized entry as foreign-edit);
  C2/C3 — negative `auto_prune_days` deleted all session state and
  negative `title_slug_max_len` panicked (getters guarded ==0 not <=0);
  M1 — MarkSource stamped at plan time stranded failed grown writes
  (moved to execute, post-write). RegenCount concern confirmed → split
  to GrownCount so routine growth no longer prompts. slug.FromTitle
  purity: cooldown/budget-pct promoted to pointer types (0 = explicit).
  Full build/lint/test green; docs updated; site rebuilt. See LEARNINGS
  2026-07-06 entries and specs/gitnexus-docker-fold.md.

- [x] Migrate Sprintf-based templates (tpl_*.go) to Go text/template or embedded
  template files — ObsidianReadme, LoopScript, and other multi-line format
  strings that can't move to YAML #added:2026-03-18-163629
  Spec: specs/tpl-text-template-migration.md
  DONE 2026-05-30 (branch refactor/tpl-text-template-migration). Tier-1 blocks
  + static Zensical + LoopScript + Tier-2 recall HTML (metaTable/details)
  migrated to embedded templates behind handles; Tier-3 single-line format
  strings, pure joins, and the RecallListRow meta-format kept as fmt.Sprintf.
- [x] P0.8.5: Harden notify resolution (reframed 2026-06-02). The original
  premise ("`ctx notify` silently fails in worktrees because the key is
  gitignored and absent") was investigated and largely disproven: with the
  default global key, notify works in worktrees (verified against a built
  binary + isolated repo + fake webhook sink). The failure only reproduces
  with a deprecated project-local key. Real defects to fix: (1) remove the
  implicit `.context/.ctx.key` resolution tier — the sole worktree-divergence
  and a documented security antipattern; (2) surface the silent fire-path
  failure when a CONFIGURED webhook can't be delivered (decrypt/read/POST),
  while keeping legitimate silences (not-configured, event-not-subscribed).
  Whether config reaches a worktree is the user's call via `.ctxrc`
  git-tracking — ctx does not special-case worktrees (it cannot distinguish a
  worktree from N side-by-side terminals). Approaches A (--git-common-dir key
  fallback) and B (copy key at worktree creation) rejected; see DECISIONS.
  Spec: specs/notify-resolution-hardening.md
  #priority:medium #added:2026-02-22 #reframed:2026-06-02
- [ ] P0.9.2: Split cli-reference.md (1633 lines) into command group pages:
  cli-overview, cli-init-status, cli-context, cli-recall, cli-tools,
  cli-system —
  each page covers a natural command group with its subcommands and flags
  #added:2026-02-24-204208
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
- [ ] Q.1: Docstring cross-reference audit — compliance test that
  flags docstrings
  mentioning domains that don't match their callers. Start with `write/**`,
  extend to all `internal/`. Spec: `specs/docstring-cross-reference-audit.md`
  #priority:medium #added:2026-03-17
- [x] Split internal/assets/embed_test.go — tests that call read/ packages
  must
  move to their respective read/ package to avoid import
  cycles #added:2026-03-18-192914
- [ ] Improve recall/core format tests — replace hardcoded string assertions
  (e.g. Contains Tokens) with semantic checks that verify structure and values,
  not label text #added:2026-03-19-194645

## Agents


## Misc






### Architecture Docs



### Code Cleanup Findings

**PD.5 — Validate:**

### Phase -3: DevEx


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


- [ ] ctx-architecture-next — fourth step in the architecture
  pipeline (map → enrich → hunt → **prescribe**).
  **Context**: The three existing skills produce inputs
  (`ARCHITECTURE.md`, `DETAILED_DESIGN*.md` from
  `/ctx-architecture`; enriched verifications from
  `/ctx-architecture-enrich`; ranked failure inventory from
  `/ctx-architecture-failure-analysis`'s `DANGER-ZONES.md`).
  But the agent then has to synthesize "so what do I DO?" on
  its own, every time, from those raw artifacts. The fourth
  step closes the pipeline by producing `NEXT-ACTIONS.md` —
  a sequenced, prioritized fix plan that maps each danger
  zone to a concrete next move (refactor, test, doc,
  escalate, accept) with effort estimates and a suggested
  order.
  **Distinct from ctx-architecture-extend (skipped)**: that
  was about *where features grow*; this is about *what to
  fix first*. Extend overlapped with DETAILED_DESIGN and
  enrich's registration sites. Next has no overlap — pure
  synthesis layer over the prior three artifacts. The
  pipeline is now 4 because each step has a distinct output
  document: map(ARCHITECTURE) → enrich(verified ARCHITECTURE)
  → hunt(DANGER-ZONES) → prescribe(NEXT-ACTIONS).
  **No MCP gateway required**: this skill consumes only the
  three Markdown artifacts produced by the prior skills,
  which already absorbed any GitNexus-derived signal during
  the enrich step. The synthesis is a pure-reasoning pass on
  the agent side. Aligns with the decision that ctx does not
  proxy / gateway companion MCPs; see DECISIONS.md
  "MCP gateway not worth the coupling cost".
  Scope sketch (refine when implementing):
    - [ ] Design SKILL.md: inputs (three artifacts), output
      shape (`NEXT-ACTIONS.md` with ranked sections), quality
      checklist (every action cites a danger zone; every
      danger zone has an action OR an explicit "accepted"
      rationale).
    - [ ] Define the action taxonomy: refactor, test, doc,
      escalate, accept. Each carries effort estimate (S/M/L)
      and a suggested sequence position.
    - [ ] Reference run against ctx itself: produce
      `NEXT-ACTIONS.md` from the existing DANGER-ZONES.md if
      one has been generated; otherwise generate the whole
      4-step pipeline against ctx as the worked example.
    - [ ] Document the pipeline order in
      `docs/recipes/architecture-mapping.md` (or wherever the
      existing 3-step recipe lives): "run all four in
      sequence; each step's output feeds the next".
      #priority:medium #added:2026-05-23


### Phase CT: Companion Tool Integration

- [ ] Backport the registry_mounts prune-fix to the os and orchestrator copies of gitnexus-docker.sh. The ctx copy (hack/gitnexus-docker.sh) now mounts every registered repo at its real path on all registry-touching branches (index-register + passthrough, matching what mcp already did), so a registry-rewriting gitnexus invocation no longer prunes repos it cannot stat in-container. The sibling copies still carry the bug — running their 'index'/'list' targets would silently deregister ctx from the global registry. Backport before using them. See LEARNINGS 'GitNexus prunes registry entries whose repo paths don't resolve in-container'. #priority:medium #session:334b20d1 #branch:main #commit:a0e5cbf9 #added:2026-07-04-212438

- [x] Add a 'make strip-gitnexus' target (backed by a hack/ script) that
  mechanically removes the GitNexus auto-injected block — delimited by <!--
  gitnexus:start --> / <!-- gitnexus:end --> markers — from AGENTS.md and
  CLAUDE.md. Marker-bounded delete (sed range or awk between markers). Must: (1)
  leave AGENTS.md as the redirect stub and CLAUDE.md ending at its Companion
  Tools / GITNEXUS.md pointer; (2) NOT touch GITNEXUS.md (the intended managed
  home for that content); (3) be idempotent (no-op when markers absent). Run it
  after 'npx gitnexus analyze'. Upstream-preferred guard is 'analyze
  --skip-agents-md'; this script is the belt-and-suspenders cleanup when analyze
  runs without that flag. Manual removal was done in 8da165a3; this automates
  it. #priority:medium #session:74c94e3a #branch:fix/notify-resolution-hardening
  #commit:8da165a3 #added:2026-06-02-085625
  DONE 2026-07-06 (session 7f6de29d, UNCOMMITTED). hack/strip-gitnexus.sh
  (awk marker-bounded delete + trailing-blank trim) + `make strip-gitnexus`
  target (added to .PHONY). Ran it: stripped both blocks — CLAUDE.md now
  ends at its Gemini/GITNEXUS.md pointer, AGENTS.md is the 3-line redirect
  stub. Verified idempotent (2nd run no-op) and zero markers/npx remain.

- [x] De-npx the gitnexus instructions across live surfaces: `npx
  gitnexus analyze` fails on machines whose Node cannot build
  gitnexus's pinned tree-sitter native addon (observed on Node 24;
  downgrading the host Node is not an option), yet the wording sits in
  CLAUDE.md's re-injected gitnexus block (note: dcbc8241 re-committed
  the block 8da165a3 deliberately stripped — the strip-gitnexus target
  above is the recurrence guard), AGENTS.md, GITNEXUS.md, the
  ctx-remember SKILL.md (claude + copilot variants — the recall
  readback literally suggests the npx command to users),
  ctx-architecture-enrich SKILL.md, docs/recipes/multi-tool-setup.md,
  and docs/home/getting-started.md. Fix: bare `gitnexus analyze` as
  canonical wording everywhere; repo-local surfaces (GITNEXUS.md,
  CLAUDE.md) additionally recommend the Docker path as the reliable
  runner (make gitnexus-index / hack/gitnexus-docker.sh — image bakes
  a compatible Node, host stays clean); deployed skill assets stay
  generic ("gitnexus analyze — or your Docker wrapper if the npm
  binary isn't viable"). Historical specs stay untouched.
  #priority:medium #session:334b20d1 #branch:main #commit:a0e5cbf9
  #added:2026-07-04
  DONE 2026-07-06 (session 7f6de29d, UNCOMMITTED). Six safe surfaces
  de-npx'd to bare `gitnexus analyze` (GITNEXUS.md + Docker path;
  getting-started.md; multi-tool-setup.md; ctx-remember claude+copilot
  generic wording; ctx-architecture-enrich generic). CLAUDE.md/AGENTS.md
  handled by removal via strip-gitnexus (their npx sat in the managed
  block). Copilot skills back in sync (make check-copilot-skills).
  site/search.json (generated) and historical specs left untouched.

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

- [ ] Refactor site/cmd/feed: extract helpers and types to core/, make Run
  public #added:2026-03-21-074859

- [ ] Add Use* constants for all cobra subcommand Use
  strings #added:2026-03-20-184639

- [ ] Systematic audit: extract all magic flag name strings across CLI commands
  into config/flag constants #added:2026-03-20-175155


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

- [x] Make TitleSlugMaxLen configurable via .ctxrc #added:2026-03-21-070944
  DONE 2026-07-06 (session 7f6de29d, UNCOMMITTED). Part of the config
  batch (with auto_prune_days, agent_cooldown_minutes, task_budget_pct,
  convention_budget_pct). Each: CtxRC field + yaml key + rc accessor
  with zero-guard fallback to the existing config default (no magic
  numbers) + ctxrc.schema.json entry + schema_test mirror + rc_test
  default/override tests; consumers rewired to the accessors. golangci
  0 issues; scoped tests green.

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

- [x] EH.1: Catalogue all silent error discards — recursive walk of
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
  #completed:2026-06-01 #branch:fix/learning-add-index-data-loss
  Done: `.context/audit/eh-silent-errors.md` catalogues all 184 non-test
  discard sites with category + recommended action. Surfaced 4 high-priority
  data-loss/crash findings (memory/publish MergePublished, pad/store
  ReadEntriesWithIDs, hub/replicate Append, memory status nil-deref) plus 11
  write-handle defer-closes. Real fix workload ≈52 sites (B/A/C/SURFACE/
  NIL-DEREF); category (d) fmt.Fprint output is an accepted end-state per EH.5
  DoD. Tests excluded from this pass.

- [x] EH.2: Address category (b) — file write/read discards. These risk silent
  data loss. Fix: return the error, or at minimum emit to stderr with
  `fmt.Fprintf(os.Stderr, "ctx: ...: %v\n", err)` following the pattern
  established in `internal/log/event.go`.
  DoD: no write/read error is silently discarded
  #priority:high #added:2026-03-14
  #completed:2026-06-01 #branch:fix/learning-add-index-data-loss
  Done: pad ReadEntriesWithIDs + hub replicate Append (commit b66816cd);
  marshal-error returns for the vscode/copilot/blocknonpathctx writers,
  journal ScanDirectory + drift reload via logWarn (41e223f5). The
  established sink turned out to be `internal/log/warn` (logWarn.Warn),
  used throughout. Verified each site by reading before editing — two of
  the catalogue's name-inferred B findings were false positives
  (MergePublished bool, LoadState value type).

- [x] EH.3: Address category (a) — file close in defer. Most are `defer func()
      { _ = f.Close() }()`. For read-only files, close errors are rare but
  should still surface. For write/append files, close can fail if the
  final flush fails — these are data loss. Fix: `if err := f.Close();
      err != nil { fmt.Fprintf(os.Stderr, "ctx: close %s: %v\n", path, err) }`.
  DoD: all defer-close sites log failures to stderr
  #priority:medium #added:2026-03-14
  #completed:2026-06-01 #branch:fix/learning-add-index-data-loss
  Done: write/append-handle closes converted to a named-return merge so
  a failed flush fails the op (6 kb appends, trace appendJSONL, skill
  copy out, kb note Run — commit 9d07da57); read-handle and gRPC-client
  closes surfaced via logWarn (commit 06109734). io/security
  SafeWriteFileAtomic was already correct (meaningful close checked;
  error-path closes annotated).

- [x] EH.4: Address category (c) — os.Remove/Rename discards. These are state
  operations (rotation, pruning, temp file cleanup). Silent failure leaves
  stale state. Fix: stderr warning at minimum; for rotation/rename, consider
  returning the error.
  DoD: no Remove/Rename error is silently discarded
  #priority:medium #added:2026-03-14
  #completed:2026-06-01 #branch:fix/learning-add-index-data-loss
  Done (commit 06a88416): os.Remove/RemoveAll surfaced via logWarn where
  failure leaves stale state with real consequences — partial skill
  install, violations file (duplicate alerts), sync lock (blocks sync),
  hub pid file, trace marker. io/security temp-file cleanup discards are
  on already-failed paths and annotated as acceptable.

- [x] EH.5: Validate — `grep -rn '_ =' internal/` returns only category (d)
  entries (fmt.Fprint to stderr) and entries explicitly annotated as
  acceptable. Run `make lint && make test` to confirm no regressions.
  DoD: grep output is clean or fully annotated; CI green
  #priority:high #added:2026-03-14
  #completed:2026-06-01 #branch:fix/learning-add-index-data-loss
  Done (commit f7bf7d8f): `grep -rn '_ = ' internal/` (non-test) = 68
  sites — 47 category-(d) fmt.Fprint (accepted end-state) + 21 explicitly
  annotated/handled. `:=`-form discards (x, _ := …) are outside this
  grep's scope. make lint = 0 issues, make test = 0 failures.
  Whole EH sweep: 7 commits (6ca1198a catalogue → f7bf7d8f), spec
  specs/error-handling-audit.md.

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





- [x] Make AutoPruneStaleDays configurable via ctxrc. Currently hardcoded to 7
  days in config.AutoPruneStaleDays; add a ctxrc key (e.g., auto_prune_days) and
  fallback to the default. #priority:low #added:2026-03-07-220512


- [x] Add ctxrc support for recall.list.limit to make the default --limit for
  recall list configurable. Currently hardcoded as config.DefaultRecallListLimit
  (20). #priority:low #added:2026-03-07-164342
  DONE 2026-07-06 (session 7f6de29d, UNCOMMITTED). Deferred from the
  config batch (its consumer sat in the reserved journal tree), then
  completed once that tree settled. yaml `recall_list_limit` +
  RecallListLimit() accessor (zero-guard → journal.DefaultRecallListLimit)
  + schema entry + schema_test mirror + rc_test default/custom;
  consumer internal/cli/journal/cmd/source/cmd.go now uses the accessor
  as the --limit default (dropped its config/journal import). Build +
  rc/assets/source tests + lint all green. Completes the ④ batch (5/5).

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

- [x] Make DefaultAgentCooldown configurable via .ctxrc — currently hardcoded
  at
  10 minutes in config #added:2026-03-07-073106

- [x] Make TaskBudgetPct and ConventionBudgetPct configurable via .ctxrc —
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
- [ ] F.1: MCP server integration: expose context as tools/resources via Model
  Context Protocol. Would enable deep integration with any
  MCP-compatible client. #priority:low #source:report-6

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

- [ ] Rewrite lint-style scripts in Go as ctxctl subcommands.
  Unblocked 2026-05-28: ctxctl now exists (`tools/ctxctl`, PR #104),
  so the prerequisite is met; would land as `ctxctl check` / lint
  subcommands per the CLI-surface tasks below. #added:2026-03-29-082958


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

### Phase SK: Skill Surface Polish (Phase 0a; prerequisite for Phase KB)
`#priority:high #added:2026-05-09`

Spec: `specs/skill-surface-polish.md` (design ref:
`ideas/002-editorial-pipeline-and-skill-rigor.md` §3 "Reframing the
wishy-washy skills")

Tightens existing capture skills to sibling-project rigor before the editorial
pipeline (Phase KB) lifts that pattern
wholesale. Independent of Phase RG; both can ship in parallel.


### Phase RG: Require Git as Architectural Precondition (Phase 0b; prerequisite for Phase KB)

`#priority:high #added:2026-05-09`

Spec: `specs/require-git.md`

Enforces what `ctx` already needs: git. `ctx` works properly only with a
repo present, and this phase makes that a runtime precondition rather than
an assumption. Breaking change for any pre-existing git-less ctx project
(N≈0 in practice). Independent of Phase SK; both can ship in parallel.

- [ ] Update `docs/recipes/bootstrap-a-project.md`, `README.md`,
  `docs/cli/init.md` to show `git init` before `ctx init`
- [ ] Tag as breaking change in `dist/RELEASE_NOTES.md` with one-command
  migration ("Run `git init` in any pre-existing
  git-less ctx projects before upgrading")

### Phase KB: Editorial Pipeline + Handover (depends on Phase SK + Phase RG)

`#priority:high #added:2026-05-09 #revised:2026-05-16`

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


Store layer (landed under `internal/write/` per the revised spec, not
`internal/store/`):


CLI commands:


Skills:


Doctor / status / .gitignore:


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

- [ ] Document MemPalace-as-ground-source recipe in
  `docs/recipes/build-a-knowledge-base.md`; uses already-specced
  `mcp:<server>:<resource>` syntax in `grounding-sources.md`; zero new ctx code

- [ ] Bug / gap: Phase KB scaffold has no retrofit path for projects
  that pre-date the kb subsystem. `coreKB.Scaffold(contextDir)` is
  only called from `internal/cli/initialize/cmd/root/run.go`'s init
  flow; `ctx init` itself refuses on populated projects without
  `--reset` (destructive). On the ctx project this branch is in,
  `.context/kb/` and `.context/ingest/` were missing entirely until
  hand-rolled on 2026-05-21 by copying
  `internal/assets/kb/templates/{ingest,kb}/*` into place. Add a
  dedicated `ctx kb init` subcommand (or `ctx init --kb-only`) that
  calls `coreKB.Scaffold` and nothing else; existing per-file
  preservation in `Scaffold` already makes it idempotent. Wire the
  command annotation so it bypasses the require-context-dir
  PreRunE gate (the gate already passes when `.context/` exists,
  but a freshly-init'd project in the same shell session must work
  too). Update `docs/recipes/build-a-knowledge-base.md` to point at
  the new subcommand for retrofit. #priority:medium #added:2026-05-21

- [ ] Bug / gap: `ctx init` refuses on a populated project without
  `--reset`, but `--reset` is destructive (it backs up populated
  files then overwrites them). There is no path between "project
  already exists, do nothing" and "blow it all away." Add an
  `--upgrade` mode that runs the scaffolding stages that are
  per-file-existence-preserving (kb, steering foundation, entry
  templates, scratchpad bootstrap if absent, gitignore amends,
  Makefile.ctx, settings.local.json permission merge) but skips
  reset-required stages (CLAUDE.md merge, populated-file refuse).
  Pairs with `ctx kb init` above; same shape, broader surface.
  #priority:medium #added:2026-05-21

#### Adjacent-tool kb ingests `#added:2026-05-21`

The kb's declared scope (`.context/kb/index.md`) covers design
lessons and operational patterns from adjacent / inspirational
AI infrastructure projects. Each entry below is a separate
`/ctx-kb-ingest` pass. Topic slugs follow the
lowercase-kebab-case convention used by `ctx kb topic new`.
Suggested invocation per row is a starting point; the operator
can refine the source URL during the pass. Mark
`[x]` only when the topic page clears the cold-reader rubric and
the source-coverage ledger has the row at `comprehensive` (or
honestly at `topic-page-drafted` if the page is good but the
ledger admits residue).


- [ ] `claude-code` — Anthropic's official CLI for Claude. Surface
  to study: hooks, slash commands, skills (`~/.claude/skills/`),
  settings.json, plugin system. Suggested seed:
  `/ctx-kb-ingest https://docs.claude.com/en/docs/claude-code claude-code`.
  Question: how does Claude Code's hook + skill surface compare
  to ctx's, and what entry points (e.g. settings.json structure)
  is ctx echoing vs diverging from? #priority:medium

- [ ] `opencode` — sst/opencode terminal AI agent. Surface to
  study: plugin model (TypeScript `index.ts`), MCP-server
  registration, skill discovery, command palette shape.
  Suggested seed: `/ctx-kb-ingest https://opencode.ai/docs opencode`.
  Question: ctx already integrates with OpenCode via
  `internal/assets/integrations/opencode/` — what's the kb
  reading of that integration as a pattern, and where could it
  generalise to other host CLIs? #priority:medium

- [ ] `cursor` — Cursor editor. Surface to study: workspace
  hooks (`.cursorrules`, `.cursor/`), MCP integration, the
  cross-IDE settings-leak that motivated ctx's state.Initialized
  gate (spec: `specs/state-dir-no-mkdir-when-uninitialized.md`).
  Suggested seed:
  `/ctx-kb-ingest https://cursor.com/docs cursor`. Question: how
  does Cursor's workspace-level hook discipline shape what ctx
  has to defend against (cross-workspace state leaks), and what
  would ctx-on-Cursor parity look like beyond the current
  defensive gate? #priority:medium

- [ ] `gitnexus` — code-intelligence MCP toolchain that ships as
  a companion to ctx (see `.claude/skills/gitnexus/`,
  `GITNEXUS.md`). Surface to study: MCP tool catalogue (cypher,
  impact, route_map, tool_map, group_*), graph-backed code
  navigation, the impedance match with Go projects of ctx's
  size. Suggested seed:
  `/ctx-kb-ingest GITNEXUS.md gitnexus` plus discovery enabled
  to pull official docs. Question: which GitNexus capabilities
  is ctx *not* using that would meaningfully change how ctx
  develops itself (e.g. blast-radius checks pre-refactor)?
  #priority:medium

- [ ] `mempalace` — memory-palace / spatial-recall AI project
  (operator: confirm the canonical URL; tentative
  `https://github.com/mempalace` or similar). Surface to study:
  whatever the project's memory-persistence model is, and how it
  differs from ctx's file-anchored memory model. Question: is
  there a spatial / graph / vector substrate worth lifting into
  ctx's memory layer, or is the contrast purely contrastive
  (ctx commits to file-anchored; mempalace commits to
  something else)? #priority:low

- [ ] `deepwiki` — Devin's auto-generated wiki for any GitHub
  repo. Surface to study: how it derives a wiki structure from
  code+commits, and whether that output is usable as a ctx
  substrate (per reminder [6]: *"use deepwiki to enhance docs
  of ctx and use it as a substrate for further analysis of
  other stuff"*). Suggested seed:
  `/ctx-kb-ingest https://deepwiki.com/ActiveMemory/ctx deepwiki`.
  Question: is deepwiki's auto-derived structure complementary
  to ctx's hand-authored docs (use both, treat them as
  different views) or competitive (one supersedes the other)?
  Connects to reminders [6, 7]. #priority:medium

- [ ] `zensical` — static-site generator that anchors to
  `zensical.toml` (referenced as the canonical
  config-file-anchored precedent in
  `specs/cwd-anchored-context.md`). Surface to study: the
  anchor-to-config-file pattern, recipe library shape, how
  zensical handles cwd vs config-dir resolution. Question: ctx
  cited zensical as precedent for the cwd-anchored decision;
  what other zensical patterns are worth borrowing or
  rejecting? #priority:low

- [ ] Discuss: rename `ctx kb site build` (referenced in the
  `/ctx-kb-ingest` skill's circuit-breaker item #3 but absent
  from the installed binary) into a top-level family —
  `ctx site kb build`, `ctx site journal build`, etc. The
  motivation: ctx now ships multiple site-shaped surfaces
  (kb topic pages, journal entries, possibly more); the
  current `kb site-review` placement under `kb/` no longer
  generalises. A top-level `site` subcommand would let each
  domain register its own `build` and `review` verbs without
  cross-domain namespace bleed. Open questions: where do the
  per-domain build implementations live (`internal/cli/site/cmd/kb/build/`?
  `internal/cli/kb/cmd/site/build/`?), how does this interact
  with the existing `kb site-review` / `ctx kb reindex`, and
  what becomes of the `/ctx-kb-ingest` skill's circuit-breaker
  reference? Treat this as a naming + topology discussion before
  any code lands; the vllm topic page is `topic-page: deferred`
  partly because of the missing build subcommand, so resolving
  this unblocks the circuit breaker too. #priority:medium #added:2026-05-21

Each row is a single `/ctx-kb-ingest` pass when started; further
follow-ups for that tool (per-category deep dives, sub-page
splits) get tracked on the source-coverage ledger, not as
TASKS.md children. Open a new TASKS row only when a *different*
adjacent tool joins the list.

- [ ] Feature: skill usage tally + ceremony-time nudge.
  Motivation: ctx ships 60+ skills; discoverability is a real
  problem. A usage tally would (a) surface usage patterns, and
  (b) let ceremonies remind the operator about under-used skills
  that might help current work. Two phases:

  **Phase 1 — instrument.** Extend the journal-enrich pipeline
  (`/ctx-journal-enrich-all` or sibling) to scan
  `~/.claude/projects/*/*.jsonl` for `Skill` tool uses and write
  two artifacts:
    - **Time-series** at `~/.ctx/state/skill-usage.jsonl`:
      append-only, one row per invocation, fields
      `{ts, project, session_id, skill_name, source:
      "claude-code"|"opencode"|...}`.
    - **Aggregate** at `~/.ctx/state/skill-usage.json`: derived
      rollup, `{skill_name → {count, first_used, last_used, projects[]}}`.
      Stays in `~/.ctx/state/` (user-global), not per-project, so
      patterns survive across projects.

  **Phase 2 — wire ceremony nudges (NOT auto-prompts).** Surface
  the tally inside two existing ceremonies, never as session-start
  noise:
    - `/ctx-remember`: at the end of the recall readback, add a
      *"unused-but-might-help"* line that names 1-3 skills with
      `last_used > 30d ago` (or `never`) whose descriptions match
      keywords from current TASKS.md focus / branch name / recent
      commits.
    - `/ctx-wrap-up`: in the candidate-proposal phase, include a
      *"this session's skill mix"* line summarising which skills
      fired this session, and surface 1-2 skills that would have
      fit the work but weren't invoked.
    - Explicitly NOT in `/ctx-handover` — that ceremony is for the
      next agent, not introspection.

  Hard anti-patterns: stale-skill-name pollution (when skills
  rename, the tally must reconcile by reading the current skill
  catalogue and dropping unknowns to a `*.deprecated.jsonl`
  archive); skill-nudge inside a tool-use loop (only at ceremony
  invocation, never via PreToolUse hook); LLM-judged matching at
  Phase 1 (start with naive string-match of skill descriptions
  against TASKS.md / branch / recent commits; revisit if the
  signal is too weak).

  Open questions: where exactly the journal-enrich pipeline
  writes the artifacts (does it touch `~/.ctx/` or keep
  per-project state and aggregate at read time?); whether the
  nudge text is rendered by ctx or by the skill itself reading
  the JSON; whether the "match current work" heuristic lives in
  Go or in the skill prompt. Tackle these at spec time, not
  implementation time. #priority:medium #added:2026-05-21

### Phase KB-followup: Adversarial design review of parallel skill trees
`#priority:medium #added:2026-05-17`

`ctx` ships skills to three host trees:
`internal/assets/claude/skills/` (canonical, full Claude tool surface),
`internal/assets/integrations/copilot-cli/skills/` (Copilot CLI;
`tools: [bash]`),
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

### Phase JR: Cold-Start Memory Recovery (semantic recall over journal history)
`#priority:medium #added:2026-05-10`

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

### Phase EVA:
`ctx kb ev append` helper — eliminate Edit-anchor brittleness for append-only
structured rows

`#priority:medium #added:2026-05-23`

**Hub relevance** (flagged 2026-05-28, hub workstream; not a re-file):
cross-tenant ingestion (the `consumed` relation, D3) appends `EV-###`
rows in the *consuming* tenant via the exact
anchor-Edit-on-prior-tail-row + verify-with-awk dance this phase
codifies away. A typed `ctx kb ev append` de-risks the hub's
cross-tenant ingest path directly, not just single-tenant
`/ctx-kb-ingest`.

**Pain point**: agents performing `/ctx-kb-ingest` passes append `EV-###`
rows to `.context/kb/evidence-index.md` via the Edit tool. The append-only
invariant means new rows go at the bottom; never reordered, never
renumbered. Today's append pattern picks an `old_string` anchor on the
NEW row's start (`| EV-NNN | ...`) and prepends the new EV before it — but
when the anchor is the prior tail row, the natural-reading intent
"insert after EV-NNN" gets accidentally implemented as "insert
before EV-NNN" — silently swapping order. Observed 3+ times in a single
DR-kb session (2026-05-23), each requiring a delete + re-insert correction
that burns context and risks deeper mistakes during fixup.

**Why this matters**: the `evidence-index.md` schema is a pipe-delimited
table with append-only ordering as a structural invariant
(`KB-RULES.md` §Source-coverage ledger + glossary expectations).
Mis-ordered rows aren't caught by `ctx kb site build` because the
build only validates references and Markdown syntax, not row ordering.
A future row that cites `EV-948` before `EV-949` appears in the file
would still resolve and build clean — but the resulting reading order
is hostile to humans + future agents diffing the table.

**Why a sort script is the WRONG fix**: sorting after the fact would
normalize my own mistakes silently. If an agent accidentally minted
the right ID with the wrong claim content (e.g., EV-950 carrying what
should have been EV-949's claim), sorting would happily preserve the
broken claim under the wrong ID. The actual issue is using a free-form
text-editing tool (Edit) for what should be a typed append operation.

**Proposed shape**: add `ctx kb ev append` CLI subcommand that takes
structured input (claim summary + source short-name + locator + sha +
confidence band + tags + extracted-date) and appends a correctly-formatted
row to `evidence-index.md` after the highest existing `EV-###` row.
Behaviors:

- Read `evidence-index.md`; find the highest `EV-NNN`; assign new row
  `EV-(NNN+1)`; refuse if `--ev` is supplied and doesn't match (catch
  agents that try to mint a specific ID).
- Validate confidence band is one of `{speculative, low, medium, high}`.
- Validate tags are comma-separated kebab-case slugs.
- Append the row + a newline; preserve all other content byte-for-byte.
- Print the assigned `EV-NNN` to stdout for downstream citation.
- Exit non-zero on any validation failure with a precise error.

**Companion: `ctx kb ev next`** — returns the next available EV number
without appending. Lets skills cite the EV ID inline in topic-page prose
BEFORE the row exists, then mint it via `ctx kb ev append --ev EV-NNN`.

**Skill changes**:

- `/ctx-kb-ingest` skill prose updated to invoke `ctx kb ev append` /
  `ctx kb ev next` instead of Edit-based row insertion.
- Same pattern applies to `Q-###` rows in `outstanding-questions.md` —
  consider `ctx kb question open` as a parallel helper if `outstanding-
  questions.md` exhibits similar issues (no direct evidence yet but
  same shape).

Deliverables:

- [ ] Spec the structured append surface (CLI flags + stdin shape +
  validation rules + output contract). Tradeoff to decide: full
  positional flags vs YAML/JSON stdin payload.
- [ ] Implement `cmd/kb/ev/append.go` + `cmd/kb/ev/next.go`.
- [ ] Add table-driven tests covering: highest-ID detection edge cases
  (empty file, only header rows, malformed prior row); validation
  failures; concurrent-write protection (file lock or
  read-modify-write check); ID-skip detection (refuse if `--ev` would
  create a gap).
- [ ] Update `/ctx-kb-ingest` skill prose to use the new CLI instead
  of Edit anchors.
- [ ] Update `KB-RULES.md` if necessary to make the helper the
  blessed path for `EV-###` appends.
- [ ] Optionally extend to `ctx kb question open` for `Q-###` rows
  with the same anti-pattern protection.

Context: filed 2026-05-23 after observing 3+ Edit-anchor swap mistakes
during a single DR-kb session that drained ~15 EVs across 9 ingest passes.
Each mistake was self-corrected but required a delete + re-insert that
burned context and added rework. The pattern would compound across the
kb's lifetime as more agents append rows; treating it as a tooling gap
rather than a discipline problem is the long-term fix. Source pointer:
DR-kb session a5736210 closeouts under
`~/Desktop/WORKSPACE/<domain>/.context/ingest/closeouts/`
20260523T044000Z + 20260523T060000Z + 20260523T080000Z reference the issue.

- [ ] Pad undo & snapshot safety net: every destructive `ctx pad`
  operation writes an encrypted snapshot to
  `.context/scratchpad.history/` before overwriting the pad, and a
  new `ctx pad undo` subcommand restores the most recent snapshot.
  Snapshot is the existing pad blob byte-for-byte (no re-encryption);
  bounded ring buffer caps storage. Driver: user accidentally `rm`'d
  a blob entry without reading it; recovery via off-host backup is a
  6-step ritual disproportionate to a single fat-finger.
  Spec: specs/pad-undo-snapshot.md #priority:high
  #added:2026-05-24
    - [x] **Phase 1**: snapshot-on-mutate + `ctx pad undo` (no flags) +
      bounded retention (count cap + age cap, defaults hard-coded) +
      unit tests covering snapshot-before-write, first-write-no-snapshot,
      undo-restores-pre-mutation, undo-is-itself-snapshotted (redo),
      empty-history-exits-zero, prune-evicts-oldest. Plaintext and
      encrypted pad modes both covered. Shipped 2026-05-24 in commit
      6bcaf889 (`feat/pad-undo-snapshot`).
    - [ ] **Phase 2**: `ctx pad undo --list` (with sidecar
      `<slot>.meta.json` for entry counts), `--to <slot>`, `--prune`,
      `--clear` (with confirmation prompt). `.ctxrc` `[pad.history]`
      block for retention tuning. Skill `ctx-pad/SKILL.md` and recipe
      `scratchpad-with-claude.md` updates.

- [ ] Out-of-band audit channel: discipline enforcement via verbatim
  relay (the one channel that survives agent tunnel vision). An
  out-of-band auditor (separate Claude Code session) drops structured
  reports into `.context/audit/<kind>.md`; the `ctx system check-audit`
  UserPromptSubmit hook relays unread reports; `ctx audit list/show/
  dismiss` manage the lifecycle. Driver: pad-undo Phase 1 shipped a
  user-facing command without docs and the in-band CONVENTIONS rule
  did not prevent it (agent that read the rule still skipped it).
  Spec: specs/audit-channel.md #priority:high #added:2026-05-24
    - [x] **Phase 1a**: `ctx audit` CLI (list/show/dismiss + --all),
      `ctx system check-audit` hook, report format + parser,
      digest-bound dismissal ledger at `.context/audit/.dismissed.json`,
      full i18n plumbing, 17 tests. Shipped 2026-05-24 in commit
      aefce517 (`feat/pad-undo-snapshot`).
    - [x] **Phase 1b**: `/ctx-surface-audit` skill (refuse-on-dirty-tree
      guard) + `docs/recipes/audit-channel.md` + index registration.
      Shipped 2026-05-24 in commit 71c3dfa4.
    - [ ] **Phase 2**: auto-dismissal on detected resolution (re-derive
      surface state on hook fire, suppress when the gap is closed);
      sibling audit skills `/ctx-spec-trailer-audit` and
      `/ctx-capture-audit`; stale-report graceful escalation; wire the
      hook into `.claude/settings.local.json` as a real UserPromptSubmit
      handler. Open questions in spec: naming collision with
      `internal/audit/` AST-tests package; shared skill-helpers library.
      Partially shipped via the ctxctl migration (PR #104): the
      repo-local `.claude/settings.local.json` hook is wired as a real
      UserPromptSubmit handler (`ctxctl audit-relay`); the
      `internal/audit/` naming collision is resolved (audit logic moved
      under `internal/ctxctl/`, AST checks made parallel-taxonomy-aware);
      stale-report escalation shipped. Remaining: auto-dismissal on
      detected resolution; sibling skills `/ctx-spec-trailer-audit` and
      `/ctx-capture-audit`; shared skill-helpers.

## Future

- [ ] Implement journal compaction: Elastic-style tiered storage with tar.gz
  backup. Spec: specs/journal-compact.md #added:2026-03-31-110005

## Human Review and Consolidation

* [ ] Human: internal/recall/parser requires a serious refactoring; for example
  the parser object and its private and public methods need to go to its own
  package and other helper functions need to go to a different adjacent package.
* [ ] Human: internal/notify/notify.go requires refactoring (all functions
  bagged in
  one file; types need to go to types.go per convention etc etc)
* [ ] Human: split err package into sub packages.

- [ ] Human: It's about time to go through the entire codebase check for
  inconsistencies, and move useful functions that are utility and/or reusable
  to relevant convenience packages.
- [ ] Human: Read the entire documentation page-by-page, line-by-line, with a
  critical mind, including blog posts. Take notes for agent to rectify, or
  directly update the docs whenever it makes sense.
- [ ] Human: Do a documentation audit for AI-generated artifacts. #important
  #not-urgent
- [x] Human: test `ctx init` on a fresh ubuntu install.
  DONE 2026-07-15 (session 87e465a0). This machine is bare-metal fresh
  Ubuntu; ctx 0.8.1 installed at /usr/local/bin/ctx. Smoke-tested in a
  throwaway temp git repo: `ctx init` created all 9 canonical files +
  steering + kb scaffold + templates, wired .claude/settings.local.json
  (plugin enabled, statusline), CLAUDE.md, Makefile, and 9 .gitignore
  entries; detected the Claude plugin (0.8.1, hot-reload). Follow-on
  `ctx status` (9 files, 22 invariants), `ctx agent` (packet rendered,
  unfilled steering tombstones correctly skipped), and `ctx drift`
  (11 checks PASSED, no drift) all clean. Temp repo removed.
- [ ] Human: These shall be done before a release cut. Especially when the
  amount of code generated is around hundreds of thousands of lines of code,
  we need to sit down and spend as much time as needed. For two reasons:
  If we (humans) don't understand the codebase fully, how can we guide AI?
  And secondly, a human scan can detect things that AI cannot find by itself.

### Phase CLI-FIX: CLI Infrastructure Fixes

- [ ] Reindex grouped-emit (ctx-side): RenderBlock should emit the CTX:KB:TOPICS
  managed block grouped by parent folder (### <group> headings) instead of one
  flat sorted list, for grouped kbs like things-wtf-dr (49 topics). ListTopics
  already returns slashed group/slug slugs (PR #106, spec
  specs/kb-reindex-nesting.md) so only RenderBlock + the consumer-facing
  block-format contract change; must still handle ungrouped/flat top-level
  topics. Deferred from the kb-reindex fix (managed-block format change).
  #priority:high active dependent work in the hub/other workstream; natural
  owner is ctx-side (ListTopics already recursive). #session:cf14dd25
  #branch:main #commit:aae42fe8 #added:2026-05-28-215308

### ctx-dream v1

> **STATUS (2026-06-20): NOT DONE.** The dream→serendipity *engine* has landed
> (triage, review CLI, executor contract, enablement docs — all `[x]` below),
but
> the workstream is incomplete. **Next up: the end-user dream UX loop** (the
first
> `[ ]` task) — making `ctx dream` a clean run→digest experience and wiring
the
> serendipity nag. The `dream-guard` consolidation is the second open item.
> Baseline is intentionally being frozen here for a controlled snapshot-VM
> experiment (SDD / spec-kit alignment); do not assume this phase is
shipped.

- [ ] Close the end-user dream UX loop: invoking /ctx-dream interactively is a
  debug-crawl, not a dream. (1) Reconsider de-listing /ctx-dream as a
  user-invocable slash-command — it's the headless executor's instruction set,
  not a user command; the user surfaces are 'ctx dream' (on-demand) + cron +
  /ctx-serendipity. (2) Make 'ctx dream' a clean run->digest experience for the
  terminal user. (3) Wire the 'ctx remind' nag ('a serendipity round is
  waiting') so dream->nag->review closes without the user watching a pass.
  #priority:high #session:2263caef #branch:main #commit:a1624af5
  #added:2026-06-07-142015

- [ ] Replace skills/ctx-dream/guard.sh with a 'ctx system dream-guard'
  subcommand (convention: hook scripts are ctx system subcommands — cf.
  block-non-path-ctx). It reads the PreToolUse tool-call JSON on stdin and
  applies internal/dream.WriteScope + Leak as the SINGLE source of truth
  (eliminating the current Go-vs-shell guard duplication/drift), emitting the
  hook block decision. Then rewire the PreToolUse settings in
  docs/recipes/run-the-dream.md and docs/reference/dream-executor-contract.md to
  call 'ctx system dream-guard', and delete guard.sh. #priority:medium
  #session:2263caef #branch:main #commit:a1624af5 #added:2026-06-07-140456

- [x] Docs: executor-contract reference for non-Claude-Code harnesses —
  bounded pass, structural guard enforcement, fail-loud,
  proposals-only-into-dreams/ #priority:medium #session:2263caef
  #branch:fix/notify-resolution-hardening #commit:ef59aeea
  #added:2026-06-07-112233

- [x] Docs: Claude Code dream enablement guide — opt-in (.ctxrc
  dream.enabled), cron entry, guard hook wiring, ctx remind cadence
  #priority:medium #session:2263caef #branch:fix/notify-resolution-hardening
  #commit:ef59aeea #added:2026-06-07-112233

- [x] Tests: git check-ignore guard refuses tracked path; ledger
  dedup-against-seen; crash-resume; 2605.12978 corrupted-artifact regression
  fixture #priority:medium #session:977ff594
  #branch:fix/notify-resolution-hardening #commit:03a24cf0
  #added:2026-06-06-162238

- [x] Build ctx dream review CLI (accept/reject/amend) plus serendipity skill;
  mechanical applies instantly, generative drops to agent; backup-before-mutate;
  ctx remind cadence #priority:medium #session:977ff594
  #branch:fix/notify-resolution-hardening #commit:03a24cf0
  #added:2026-06-06-162238

- [x] Implement disciplined ideas triage: classify, ground against code and
  specs, semantic dedup; emit atomic provenanced proposals #priority:medium
  #session:977ff594 #branch:fix/notify-resolution-hardening #commit:03a24cf0
  #added:2026-06-06-162238

- [x] Build proposal/ledger/state machinery: per-source state record,
  append-only ledger recording rejections, two-clocks read model #priority:high
  #session:977ff594 #branch:fix/notify-resolution-hardening #commit:03a24cf0
  #added:2026-06-06-162238

- [x] Build the three structural guards: write-scope, sources-as-data, dont-leak
  (git check-ignore refuses tracked paths) #priority:high #session:977ff594
  #branch:fix/notify-resolution-hardening #commit:03a24cf0
  #added:2026-06-06-162238

- [x] Settle executor: cron claude -p bounded scheduled pass; safety invariants
  structural not prompt-level #priority:high #session:977ff594
  #branch:fix/notify-resolution-hardening #commit:03a24cf0
  #added:2026-06-06-162238

### Misc

- [x] [Epic F] ctx index: docs (remove reindex, add ctx index) + final build/lint/test gate (T23-T24). Plan: specs/plans/computed-index-projection.md #session:75be038e #branch:main #commit:f382bee7 #added:2026-07-14-054851

- [x] [Epic E] ctx index: strip INDEX blocks from .context files, remove marker constants, add guards (T19-T22). Plan: specs/plans/computed-index-projection.md #session:75be038e #branch:main #commit:f382bee7 #added:2026-07-14-054851

- [x] [Epic D] ctx index: new 'ctx index <file>' command with --depth/--json + error handling (T13-T18). Plan: specs/plans/computed-index-projection.md #session:75be038e #branch:main #commit:f382bee7 #added:2026-07-14-054851

- [x] [Epic C] ctx index: detach entry-write path from index blocks + delete block-maintenance API (T09-T12). Plan: specs/plans/computed-index-projection.md #session:75be038e #branch:main #commit:f382bee7 #added:2026-07-14-054851

- [x] [Epic B] ctx index: remove reindex command surfaces — ctx reindex / decision reindex / learning reindex (T05-T08). Plan: specs/plans/computed-index-projection.md #session:75be038e #branch:main #commit:f382bee7 #added:2026-07-14-054851

- [x] [Epic A] ctx index: rename internal/index→internal/heading + generic ATX heading matcher (T01-T04). Plan: specs/plans/computed-index-projection.md #session:75be038e #branch:main #commit:f382bee7 #added:2026-07-14-054851

- [x] ctx-remember nudge live-credit: the check-ceremony nudge is journal-driven (ScanJournalsForCeremonies over recent IMPORTED journals), so it can't credit the current live session's /ctx-remember and misfires until the session is imported. Have the /ctx-remember skill (or ctx) touch the ceremony throttle/marker when it runs live, so the signal reflects the current session instead of waiting for journal import. See internal/cli/system/cmd/checkceremony/run.go (remindedFile/ThrottleID) #session:75be038e #branch:main #commit:f382bee7 #added:2026-07-13-220031
  DONE 2026-07-15 (session 87e465a0). Fixed in ctx itself (not the skill):
  checkceremony.Run now parses the live prompt and, when it IS a ceremony
  command, touches the daily marker (credits the live session) — no more
  journal-import lag. Unified with the self-suppress fix below via
  ceremony.InvokedByPrompt. Verified end-to-end against the built binary.
  Spec: specs/ceremony-nudge-live-session.md.

- [x] ctx-remember nudge self-suppress: the check-ceremony hook fires the 'try starting with /ctx-remember' relay even on the prompt that IS /ctx-remember, because entity.HookInput doesn't parse the UserPromptSubmit 'prompt' field (only session_id + tool_input.command). Add Prompt to HookInput and skip the ceremony nudge when the prompt starts with /ctx-remember or /ctx:ctx-remember. See internal/cli/system/cmd/checkceremony/run.go + internal/entity/hook.go #session:75be038e #branch:main #commit:f382bee7 #added:2026-07-13-220031
  DONE 2026-07-15 (session 87e465a0). Added Prompt (json:"prompt") to
  entity.HookInput; checkceremony.Run now returns without nudging when the
  live prompt is a ceremony command (ceremony.InvokedByPrompt matches the
  first token against the bare /ctx-remember|/ctx-wrap-up and plugin
  /ctx:ctx-remember|/ctx:ctx-wrap-up forms — first-token equality, so
  /ctx-remembering and prose mentions don't match). Verified end-to-end:
  bare + plugin forms suppress, normal prompt still nudges. Tests:
  ceremony_test.go (12 cases) + checkceremony/run_test.go. Spec:
  specs/ceremony-nudge-live-session.md.

- [ ] ctx list/search: richer query surface over knowledge files (filtering, full-text) layered on top of the thin `ctx index` heading-projector — successor to the queued 'CLI-projected list/search' idea; index ships first as the projection primitive #session:75be038e #branch:main #commit:f382bee7 #added:2026-07-13-215523

- [x] Re-sign the release tags (v0.1.0 through v0.8.0 and latest): the 2026-07-06 DCO history rewrite stripped their GPG signatures when the tagged commits changed SHA. #session:2cff382a #branch:fix/jumbo-diff-review-fixes #commit:945850af #added:2026-07-06-214523
  DONE 2026-07-15 (session 87e465a0). All 8 annotated tags (latest, v0.1.0,
  v0.1.1, v0.1.2, v0.2.0, v0.3.0, v0.6.0, v0.8.0) were annotated-but-unsigned
  after the rewrite. Re-signed each FAITHFULLY: preserved the exact target
  commit, tagger identity (older tags stay alekhinejose@gmail.com, newer
  jose@ctx.ist — the signing key carries both UIDs), tagger date+tz, and
  message (extracted from the object, rebuilt with --cleanup=verbatim + per-tag
  GIT_COMMITTER_NAME/EMAIL/DATE). Verified all 8: `git tag -v` = Good signature
  AND object-minus-signature byte-identical to the pre-resign backup (only a
  signature was added). LOCAL only — pushing the re-signed tags to origin
  (git push origin --tags --force) is the maintainer's step; force needed
  because origin holds the unsigned versions.

- [x] Create a /ctx-pr skill: scaffold a PR body from the branch's commits, Spec: trailers, and closed TASKS, written to inbox/ (gitignored) for the user to paste. MUST enforce the no-agent-signoff convention: no 'Co-Authored-By' and no 'Generated with ...' footer, per CONSTITUTION Process Invariants. #session:2cff382a #branch:fix/jumbo-diff-review-fixes #commit:945850af #added:2026-07-06-213149
  DONE 2026-07-15 (session 87e465a0). Shipped as REPO-INTERNAL _ctx-pr
  (.claude/skills/_ctx-pr/SKILL.md, `_` prefix = not bundled in the plugin;
  chosen over a shipped ctx-pr because it hard-enforces ctx's own CONSTITUTION
  conventions + writes to the ctx-repo inbox/). Derives the body from
  git log <base>..HEAD (subjects/bodies/Spec: trailers), the deduped specs,
  and the [ ]→[x] TASKS diff; writes inbox/pr-<branch>-<UTCstamp>.md. Hard
  constraints in a self-check block: no Co-Authored-By / agent sign-off /
  "Generated with…" footer, no git push / gh pr create, empty base..HEAD
  refuses to fabricate. Spec: specs/ctx-pr-skill.md.

- [ ] New orchestrator skill /ctx-architecture-deep-dive: wrap the three-pass architecture arc (/ctx-architecture principal → /ctx-architecture-enrich → /ctx-architecture-failure-analysis) plus the synthesis step (milestone-readiness note → /ctx-task-out --milestone <next>) into one parameterized skill with machine-checkable preconditions (code-intel MCP actually serving the repo, index fresh vs HEAD, synced tree, fresh session). Prior art: zhc/os docs/runbooks/architecture-deep-dive.md — a runbook whose pasted prompt rotted within ONE milestone (needed a 'Historical' banner because it hard-codes milestone facts like 'M0b is untasked'); a skill that derives milestone state from specs/plans/ at runtime doesn't rot. The site recipe architecture-deep-dive documents the arc as prose — this skill would be its ceremony (see the 'unceremonied pipeline step' learning). os prototypes a project-local version first and folds lessons back here. #priority:medium #session:6c276362 #branch:main #commit:a0e5cbf9 #added:2026-07-04-210547

- [x] Skill assets hard-code `npx gitnexus analyze` as the stale-index remedy, but on hosts where GitNexus runs via Docker (tree-sitter@0.21.1 native addon vs Node 24 ABI — no arm64 prebuilt) that command is a silent no-op; os/GITNEXUS.md documents this and both os and ctx expose `make gitnexus-index` → hack|scripts/gitnexus-index.sh instead. Fix the suggestion to be project-aware: point at the repo-local indexing entry point (a `gitnexus-index` make target, an indexing script, or the repo's GITNEXUS.md instructions) and fall back to `npx gitnexus analyze` only when nothing repo-local exists. Sites: internal/assets/claude/skills/ctx-remember/SKILL.md:156, internal/assets/integrations/copilot-cli/skills/ctx-remember/SKILL.md:155, internal/assets/claude/skills/ctx-architecture-enrich/SKILL.md:34+112+130. Found live: /ctx-remember in the os project emitted the broken npx suggestion while the working `make gitnexus-index` existed one directory over. #priority:medium #session:6c276362 #branch:main #commit:a0e5cbf9 #added:2026-07-04-205437
  DONE 2026-07-15 (session 87e465a0). Reworded all four sites (ctx-remember
  companion-check suggestion; enrich precondition, no-MCP block, and >5-commit
  hard-stop) to prefer the repo's own indexing entry point (make gitnexus-index
  target / script / GITNEXUS.md) and fall back to bare `gitnexus analyze` only
  when none exists — kept generic since these ship to arbitrary projects. (npx
  itself was already gone via the 2026-07-06 de-npx pass; this is the
  project-aware follow-up.) Copilot ctx-remember copy regenerated via
  sync-copilot-skills; check-copilot-skills green. Spec:
  specs/gitnexus-project-aware-reindex.md.

- [x] Drop the persisted INDEX blocks from DECISIONS.md/LEARNINGS.md;
  project the index on demand via new CLI verbs instead. The headings
  ARE the index (`## [TS] Title` is one grep away); the stored table
  is 6-7% dead weight on every read (measured 2026-07-04: 5,021B of
  76,866B / 7,774B of 107,350B), doubles the merge-conflict surface
  of every contributor PR that adds an entry (PR #117 conflicts in
  two hunks — index row + body), and is the root cause of the
  learning-add clobber bug class that index.Validate exists to guard.
  An in-file index also only helps if the agent stops reading at the
  marker — which nothing enforces; a command's output is bounded by
  construction and can never go stale. Build order: (HL.1)
  `ctx decision list` / `ctx learning list` (+ optional `show <ts>`)
  projecting the headings; (HL.2) `ctx search <needle>` across the
  five canonical files + .context/archive/** with labeled sources;
  (HL.3) entry-add stops touching the index, one-shot tolerant strip
  of INDEX:START/END (present → remove, absent → no-op; downstream
  repos migrate on next write), reindex subcommands retired; (HL.4)
  reword reader surfaces (agent packet, ctx-remember, CLAUDE.md) to
  "run list/search", coordinating with the read-discipline task
  above; (HL.5) audit index-table consumers (packet builder, drift,
  recall format tests, hub rendering) before the strip. Obviates the
  2026-03-06 "Consider indexing tasks and conventions" task (opposite
  direction — do NOT add more indexes). Full analysis + plan:
  inbox/2026-07-04-memory-file-index-drop.md (local only).
  Spec: specs/computed-index-projection.md (debated 2026-07-14 via
  /ctx-plan; brief at .context/briefs/20260714T045513Z-computed-index-projection-over-persisted-blocks.md).
  Design evolved from this task's build order: HL.1's per-noun
  `ctx decision/learning list` is replaced by a single GENERIC
  `ctx index <file>` heading-projector (works for TASKS phases too),
  and HL.2's `ctx search` is deferred to a separate follow-up task
  (richer list/search layers on top of the index primitive). The
  cold-bucket/time-shard idea was killed by a validation pass
  (~1.5% of entries are genuinely superseded — corpus already GC'd
  by consolidation). `internal/index` is RENAMED to `internal/heading`
  (parser retained; block-maintenance half deleted).
  #priority:medium #session:334b20d1 #branch:main #commit:a0e5cbf9
  #added:2026-07-04

- [ ] `ctx hub --mode git`: a git-backed hub mode beside the existing
  server mode — the hub becomes a checkout on disk (same-repo dir or
  dedicated knowledge repo), no daemon, no ports, no tokens. Core
  shape: open-intake immutable inbox/ (mechanical validation floor is
  the only write gate: schema, required evidence, size caps, secret
  scan; entries inert until promoted) → human curation verbs
  (promote/reject/defer) recording every disposition in a per-scope
  ledger → living accepted pages → gitignored, bounded,
  byte-reproducible .context/HUB.md overlay surfaced as a budgeted
  section of the `ctx agent` packet. Precedence: local .context/
  always wins; no conflict resolution by design (human-paced git sync
  replaces consensus — the trust model in internal/hub/doc.go already
  says every token holder is trusted, so election/failover is
  infrastructure the deployment shape doesn't need). Lifecycle without
  wiki rot: supersedes/reverifies/tombstones relations; delete is a
  deliberately expensive secrets runbook. Trust = prevent (pathguard
  on accepted/**) + detect (validate flags gate bypass) + attribute
  (git commits + per-entry actor) + repair (re-curate). Absorbs the
  Future task "Hub curation: immutable promotion ledger + mechanical
  validate floor" (#added:2026-07-04-153004) — that idea, fully
  formed. Full analysis + phased implementation plan (HG.1–HG.6,
  data shapes, testing spine, open questions):
  inbox/2026-07-04-hub-git-mode-analysis-and-plan.md (local only).
  Needs /ctx-spec before implementation; --share semantics change
  needs a decision record. #priority:medium
  #session:334b20d1 #branch:main #commit:a0e5cbf9 #added:2026-07-04

- [ ] Context-file read discipline for agent instruction surfaces:
  CLAUDE.md's "Do You Remember" section, the ctx-remember skills
  (claude + copilot variants), and the `ctx agent` packet's "Read These
  Files" section all instruct FULL reads of TASKS.md / DECISIONS.md /
  LEARNINGS.md. At current sizes (TASKS.md 133KB ≈ 52k tokens,
  LEARNINGS.md 107KB, DECISIONS.md 77KB) a literal-minded agent pages
  through 100k+ tokens at session start; the harness read-cap only
  slows this, it doesn't stop it. Apply the GITNEXUS.md yield pattern
  inward: every instruction surface directs the agent to (1) the
  budgeted `ctx agent` packet, (2) a projected one-line-per-entry
  view — `ctx decision list` / `ctx learning list` once the
  index-drop task below lands (command output is bounded by
  construction; until then, the file's INDEX block), (3) targeted
  entry-body reads by timestamp — never a full-file page-through. Consider `ctx agent` emitting per-file
  size/token hints so agents can self-limit. Name the existing relief
  valves in the instructions: decisions-reference.md /
  learnings-reference.md offload, /ctx-consolidate, ctx task archive.
  #priority:high #session:334b20d1 #branch:main #commit:a0e5cbf9
  #added:2026-07-04

- [-] Fold journal import into /ctx-wrap-up so users never have to remember
  to import journals (original sketch: an --exclude/--skip-active flag so
  the ceremony could skip the live session whose import would otherwise
  freeze a truncated transcript). SKIPPED same-day: the flag treats the
  symptom. Making import growth-aware removes the hazard entirely and
  needs no new flags — superseded by Phase JI below, spec
  specs/journal-import-self-heal.md. #priority:medium
  #session:334b20d1 #branch:main #commit:a0e5cbf9 #added:2026-07-04-164037
  #superseded:2026-07-04

- [x] Implement ctx system statusline per specs/statusline.md: stdin-JSON render (model, ctx%, cost), .ctxrc statusline block, setup merge with backup/restore. Informational only; no gating (spec records why) #priority:medium #session:a31b3e67 #branch:main #commit:687bbd59 #added:2026-07-04-140249

- [x] Companion-tool feature-delta analysis; borrow/skip verdicts recorded in gitignored inbox/ notes (local only) #session:a31b3e67 #branch:main #commit:687bbd59 #added:2026-07-04-135227

- [x] Add ctx-humanize plugin skill (SKILL.md + references/pattern-catalog.md, docs entry) with progressive disclosure and voice guardrails. Spec: specs/ctx-humanize.md #session:a31b3e67 #branch:main #commit:687bbd59 #added:2026-07-04-134612

- [x] Add hack/check-tools.sh tooling dependency checker (manifest: hack/tool-versions.txt, make check-tools) — spec: specs/check-tools.md #session:a31b3e67 #branch:main #commit:687bbd59 #added:2026-07-04-132852

- [x] Recipe: the full design-to-implementation pipeline from the operator's seat (new 'spec-driven-development.md' or a major fold into design-before-coding.md). Must cover, for a newcomer with zero tribal knowledge: (1) the 5-step chain INCLUDING /ctx-plan — design-before-coding.md's TL;DR currently omits the debated-brief step entirely; (2) altitude: the bet is debated once and the spec covers ALL milestones — briefs are per-bet, never per-milestone; (3) plans are just-in-time per milestone behind the rolling-wave gate (tasking distant milestones produces fiction); (4) blocking-TBD gates as the replacement for per-milestone debates — each task-out run forces exactly the decisions that milestone embeds, into DECISIONS.md; (5) two surfaces, one truth: plan = execution ledger (st column), TASKS.md epics = one-way projections over disjoint id ranges; DoD is measurement/Board-confirmed, never derived; (6) when a NEW brief happens: new bet (e.g. deferred machinery returning) or evidence falsifying the committed bet — never relitigating from below; (7) a worked multi-milestone example, not a one-session feature. Origin: zhc/os session a63353a3 (2026-07-03) — the operator had to reverse-engineer all seven from skill texts and agent explanations #priority:high #session:a63353a3 #branch:main #commit:511a609a #added:2026-07-03-232251
  DONE 2026-07-06 (session 7f6de29d, UNCOMMITTED). New file
  docs/recipes/spec-driven-development.md (all 7 points; five-step chain
  with /ctx-plan first-class; worked four-milestone example — `ctx digest`)
  + an entry in docs/recipes/index.md. Kept as a NEW file, not folded
  into design-before-coding.md (different altitude: on-ramp vs operator
  manual). Two flags for later: design-before-coding.md's TL;DR still
  omits /ctx-plan (the new recipe compensates); and point 5's "Board" is
  sibling-project vocab — rendered in ctx's own terms ("measurement or by
  you"). Verified every claim against the skill texts.

- [ ] Tighten /ctx-spec (and /ctx-implement) skills to bake in spec-kit's one
  good discipline without adopting spec-kit: every /ctx-spec must end with (a)
  explicit, testable acceptance criteria and (b) a derived, dependency-ordered
  task breakdown into TASKS.md; and the spec is FROZEN once /ctx-implement
  begins — changes go through a new decision record, not a silent edit.
  Rationale + full ctx-vs-spec-kit analysis in
  inbox/2026-06-27-ctx-vs-spec-kit-sdd-analysis.md (keep ctx's edge: persistent
  memory + adversarial /ctx-plan; avoid two-constitutions double-bookkeeping).
  #priority:medium #session:210b77dd #branch:main #commit:6b0d0107
  #added:2026-06-27-222130

### Phase JI: Self-Healing Journal Import

Spec: `specs/journal-import-self-heal.md`. Read the spec before starting
any JI task. Goal: journal import becomes growth-aware and safe to run
at any moment — from /ctx-wrap-up, a SessionEnd hook, or by hand — so
importing is never something the user has to remember or time. Claude
Code transcripts are append-only, so "source grew" is detectable from
mtime+size and a partial import is just an intermediate state the next
sweep completes. No new flags.

- [x] JI.1: Journal state schema v2 — add per-session source tracking
  to internal/journal/state: `sessions` map keyed by session id
  recording source_file, source_mtime, source_size; add `render_hash`
  to per-file entries. Version bump 1→2; v1 loads tolerantly (missing
  maps initialise empty); first v2 run adopts already-imported sessions
  by recording current source stats WITHOUT re-rendering. Tests: v1
  round-trips to v2 losslessly, adoption is idempotent, missing-file
  sources adopt as zero-stat (next stat marks Grown → captured, per
  the statSource-fails-open principle). #priority:high
  #session:334b20d1 #branch:main #commit:a0e5cbf9 #added:2026-07-04-173100

- [x] JI.2: Growth-aware planning — replace plan.Import's
  exists-check with a New/Grown/Unchanged decision from state v2.
  Grown = recorded mtime OR size differs, including the chosen
  (richest) transcript switching to a larger resume copy. Part-aware
  re-render: growth only appends messages, so re-render the last part
  plus newly created parts; earlier parts untouched by construction.
  ActionRegenerate demoted to explicit-flag edge case (mass re-render
  after format changes; healing pre-v2 truncated entries). Locked
  entries never rewritten (existing invariant). Tests: import → grow
  source → import yields complete entry flag-free; double sweep is
  byte-identical and reports all Unchanged. #priority:high
  #session:334b20d1 #branch:main #commit:a0e5cbf9 #added:2026-07-04-173100

- [x] JI.3: Foreign-edit safety via render hash — every ctx-authored
  write of a journal entry (import, enrich, normalize, fence-verify)
  refreshes render_hash in state; the hash always reflects the last
  ctx-authored write. On Grown: hash matches → splice fresh transcript
  preserving enriched frontmatter (existing keep-frontmatter
  machinery); hash differs or absent (pre-v2) → leave file untouched,
  warn naming the file, suggest `ctx journal lock` or explicit
  --regenerate. Never clobber. Tests: hand-edited entry survives a
  Grown sweep byte-identical + warned; enriched frontmatter survives a
  Grown re-render. #priority:high
  #session:334b20d1 #branch:main #commit:a0e5cbf9 #added:2026-07-04-173100

- [x] JI.4: Live-transcript parse tolerance — the JSONL session
  parser treats a trailing partial line as end-of-input (truncate to
  last complete line), no error, no warning; the missing record is
  captured by the next sweep. Test fixture: transcript cut mid-record.
  #priority:medium
  #session:334b20d1 #branch:main #commit:a0e5cbf9 #added:2026-07-04-173100

- [x] JI.5: Wrap-up integration — /ctx-wrap-up runs
  `ctx journal import --all -y` as a best-effort, NON-BLOCKING step
  before delegating to /ctx-handover; an import failure never blocks
  the handover. Update the wrap-up skill + its docs; enrichment stays
  out of the ceremony (LLM pass, belongs to /ctx-journal-enrich-all).
  Depends on JI.1–JI.3 (safe to run mid-session). #priority:medium
  #session:334b20d1 #branch:main #commit:a0e5cbf9 #added:2026-07-04-173100

- [x] JI.6: SessionEnd hook — hooks.json entry firing the import
  automatically on session end (thin `ctx system` verb per hook
  conventions), making the ceremony step belt-and-suspenders
  (import is idempotent; redundancy costs one stat per session).
  Must be covered by TestShippedHooksResolveToRegisteredCommands.
  Depends on JI.1–JI.3. #priority:medium
  #session:334b20d1 #branch:main #commit:a0e5cbf9 #added:2026-07-04-173100

  DONE 2026-07-06 (session 7f6de29d, branch main; UNCOMMITTED —
  awaiting a GPG-capable terminal). Whole phase shipped:
  - JI.1: state schema v2 (state/types.go: Sessions map + Source,
    File.RenderHash; state.go: CurrentVersion 1→2, tolerant v1 load
    normalising to v2, SessionSource/MarkSource; hash.go: HashRender,
    RenderHash/SetRenderHash). Tests: v1 tolerance, source round-trip,
    accessors, nil-map.
  - JI.2+JI.3 landed together (Grown re-render is unsafe without the
    hash guard): plan.Import decides New/Grown/Unchanged/Adopt/
    ForeignEdit from source mtime+size vs state, hash-guarding every
    Grown re-render. RENDER HASH IS OVER THE BODY (frontmatter
    stripped): enrich/normalize/fence are agent-side (no ctx write to
    refresh a whole-file hash; MarkEnriched/Normalized/FencesVerified
    are dead code), so agent-side enrichment (frontmatter-only) still
    verifies ctx-owned while a hand-edited body reads as foreign. New
    entity.ActionForeignEdit + label.reason-edited i18n; execute.Import
    stamps the hash on write and skips+warns on ForeignEdit.
    plan_test.go covers every decision.
  - JI.4: parser already skipped malformed lines; pinned with
    TestClaudeCodeParser_PartialTrailingLine.
  - JI.5: /ctx-wrap-up Phase 4.5 (best-effort, non-blocking import)
    across all three skill variants (claude/copilot/opencode).
  - JI.6: SessionEnd hook wires `ctx journal import --all -y` DIRECTLY
    (mirrors the `ctx agent` hook), not a thin system verb — the
    hooks-wiring guard resolves any `ctx <path>`, so no new command is
    needed and ceremony+hook share one code path. Compliance guard
    passes. Spec §5 updated to match.
  DEFERRED refinement: Grown re-renders count toward RegenCount, so an
  interactive `ctx journal import --all` without -y prompts on growth
  (safe, hash-guarded; hook/ceremony use -y). A non-prompting Grown
  path is a possible follow-up.
  Spec: specs/journal-import-self-heal.md.

- [x] JI.7: Docs + flag story — cli-reference and journal recipe:
  document self-healing semantics (import any time, sweeps complete
  what they started), reposition --regenerate/--keep-frontmatter as
  edge-case tools, document the one-time --regenerate heal for pre-v2
  truncated entries. No new flags is a feature; say so. #priority:low
  #session:334b20d1 #branch:main #commit:a0e5cbf9 #added:2026-07-04-173100
  DONE 2026-07-06 (session 7f6de29d, UNCOMMITTED). docs/cli/journal.md
  import section rewritten: self-healing model (new + grown + skip-
  unchanged), "edits never clobbered" (foreign-edit warning →
  lock/--regenerate), --regenerate repositioned as the edge-case tool
  (format re-render + one-time pre-self-heal truncation heal), the
  SessionEnd-hook/wrap-up auto-import, and "no new flags is the
  feature" stated. Also de-staled the skip-existing framing in
  docs/recipes/publishing.md and docs/home/common-workflows.md.
  (session-archaeology.md's --regenerate note left as-is — still
  accurate.) Whole Phase JI (JI.1–JI.7) now complete.

### Future

- [ ] Hub curation: immutable promotion ledger (who accepted what, when, why) + mechanical validate floor for shared knowledge; revisit when ctx hub grows team-curation workflows #priority:low #session:a31b3e67 #branch:main #commit:d800734c #added:2026-07-04-153004

- [ ] ctx-spec-views skill: manager-facing read-models (execution plan, spec briefs, task breakdowns) generated FROM specs/, never source of truth; spec it when someone actually needs the leadership view #priority:low #session:a31b3e67 #branch:main #commit:d800734c #added:2026-07-04-153004

- [ ] KB convention: pinned upstream corpus for grounding — document a snapshot mode (dated local copy of high-churn upstream docs as the citable byte-stream) in KB rules; no code needed #priority:low #session:a31b3e67 #branch:main #commit:d800734c #added:2026-07-04-153004

### Phase PD-M1: Progressive Disclosure — Guards, Invariants, Vocabulary

Plan: `specs/plans/pd-m1.md` · Spec: `specs/progressive-disclosure.md`
Milestone 1 builds the refusal machinery and proves the layout premise.
**Nothing moves** — no entry body is relocated and no gist is authored;
the pass (M2+) that moves bodies is the clobber risk class, so guards
land first.

**Completion rule**: an epic below is checked `[x]` only when every task
in its range is `[x]` or `[o]` in `specs/plans/pd-m1.md`. The plan — not
this list — is the single source of truth for milestone progress.

- [ ] [E1] Structural vocabulary, types, error sentinels (T01–T03). Plan: specs/plans/pd-m1.md #priority:medium #session:87e465a0 #branch:design/progressive-disclosure #added:2026-07-16

- [ ] [E2] Root parser + Validate precondition (T04–T06). Plan: specs/plans/pd-m1.md #priority:medium #session:87e465a0 #branch:design/progressive-disclosure #added:2026-07-16

- [ ] [E3] Cross-file invariants: pairing, uniqueness, links (T07–T09). Plan: specs/plans/pd-m1.md #priority:medium #session:87e465a0 #branch:design/progressive-disclosure #added:2026-07-16

- [ ] [E4] Layout proofs — the add-path de-risking; MEASUREMENT GATE: if these fail, the "zero change to add" premise is wrong and the spec's Layout section must be revisited via /ctx-plan (T10–T12). Plan: specs/plans/pd-m1.md #priority:high #session:87e465a0 #branch:design/progressive-disclosure #added:2026-07-16

- [ ] [E5] doc.go, compliance wiring, milestone gate (T13–T15). Plan: specs/plans/pd-m1.md #priority:medium #session:87e465a0 #branch:design/progressive-disclosure #added:2026-07-16

- [ ] Progressive disclosure for canonical context files: the growth warnings (LEARNINGS/DECISIONS/CONVENTIONS over threshold) are NOT redundancy — consolidation only got LEARNINGS 98→88 because the entries are distinct, dense signal. The real lever is a structural pass: canonical files carry a tight summary/index and detail loads on demand (via `ctx index`/`ctx search` projection + an archive/detail tier). Manual design exercise first (/ctx-brainstorm → spec), then codify the repeatable procedure as a new skill (e.g. /ctx-progressive-disclosure). This exercise IS the baseline for the skill. #priority:medium #session:87e465a0 #branch:main #added:2026-07-16
  DESIGN DONE 2026-07-16 (session 87e465a0): /ctx-brainstorm run to
  completion (Understanding Lock → approaches → stress-test → design).
  Spec: specs/progressive-disclosure.md. Decision: DECISIONS.md
  [2026-07-16-215955]. Shape: each canonical root becomes BOUNDED —
  staging zone (where `add` already writes, zero code change) + `## Themes`
  gists+links; bodies roll out to .context/<noun>/<theme>.md; staging IS
  the watermark (no state file); structure self-similar so nesting is
  available but deferred. Scope LEARNINGS/DECISIONS/CONVENTIONS; excludes
  CONSTITUTION + TASKS. IMPLEMENTATION STILL OPEN — phasing sketched in
  the spec (guards+invariants first, then the pass-as-skill dry-run, then
  rollout). The skill itself is the final deliverable.

- [ ] Journal parser schema drift: `ctx journal import` reports "Schema drift detected" — Claude Code transcripts now emit fields/records the parser doesn't recognize: unknown fields classifierMetaLines, session_id, toolDenialKind; unknown record type file-history-delta; unknown block type fallback. Import still works but silently ignores them. Update internal/journal parser to recognize (or explicitly skip) the new CC schema; run `ctx journal schema check` for the full report. #priority:medium #session:87e465a0 #branch:main #added:2026-07-16

- [ ] Drift path-checker false-positives on angle-bracket placeholders: `ctx drift` flags CONVENTIONS.md template placeholders (`internal/assets/claude/skills/ctx-<area>/SKILL.md`, `docs/recipes/<related-recipe>.md`, `docs/cli/<command>.md`) as missing path references. These are illustrative `<…>` placeholders, not real paths — the path checker should skip any path segment containing `<`/`>`. (Same class as the just-consolidated "detection scripts flag illustrative examples" learning.) #priority:low #session:87e465a0 #branch:main #added:2026-07-16

- [ ] ARCHITECTURE.md package-doc drift: `ctx drift` reports ~44 internal/ packages "not documented" in ARCHITECTURE.md (incl. the renamed internal/heading). Either backfill the package coverage via /ctx-architecture, or reframe ARCHITECTURE.md's scope so the drift check doesn't expect exhaustive per-package coverage. #priority:low #session:87e465a0 #branch:main #added:2026-07-16
