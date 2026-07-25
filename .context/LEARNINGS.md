# Learnings

<!--
UPDATE WHEN:
- Discover a gotcha, bug, or unexpected behavior
- Debugging reveals non-obvious root cause
- External dependency has quirks worth documenting
- "I wish I knew this earlier" moments
- Production incidents reveal gaps

DO NOT UPDATE FOR:
- Well-documented behavior (link to docs instead)
- Temporary workarounds (use TASKS.md for follow-up)
- Opinions without evidence
-->


## [2026-07-25-124457] Using the proprietary sibling repo as design evidence leaks its internals into tracked files

**Context**: While deciding the pd-m4 add-path shape, I read the sibling repo's convention file to settle the question, then quoted its guide text and attributed the decision to it in a tracked plan file. An unrelated build warning prompted the sweep that caught it.

**Lesson**: The leak vector is not NAMING the sibling (allowed) — it is citing its designs or guide text as evidence in tracked specs/plans. Both the quote and the attribution are leaks.

**Application**: Cite ctx's OWN files as design evidence and describe designs on their merits; never quote or attribute to the sibling in tracked files; git grep the staged tree before committing any spec authored while consulting it. Per specs/public-repo-hygiene.md.

---

## [2026-07-25-124457] Skill and doc examples of a serialized structure must round-trip through the real parser

**Context**: The ctx-digest skill documented digest-plan entries as JSON strings, but Assignment.Entries is []StagedEntry{timestamp,title}; the documented form failed json.Unmarshal. Latent since M3 — any agent following the skill verbatim failed at the apply step.

**Lesson**: A hand-authored example of a serialized structure drifts from the code silently; nothing tests the prose.

**Application**: Verify skill/doc JSON/YAML examples against the actual unmarshal target — build the binary and run the exact documented command — before shipping the skill.

---

## [2026-07-25-124457] A guard derived from a capability accessor silently lifts when the accessor is extended

**Context**: pd-m4 T04 made ThemeDir(convention) return true; Apply's refusal of the convention kind was DERIVED from ThemeDir returning false. That one-line 'vocabulary' change silently opened the destructive mover to CONVENTIONS.md before the convention parse/validate/dup-title paths existed. Caught by TestApply_Refusals.

**Lesson**: When a refusal is derived from a capability check (if !capable { refuse }) instead of an explicit guard, extending that capability lifts the refusal as an invisible side effect — and the two tasks look independent in the plan.

**Application**: Before extending a capability accessor, grep for guards that read it. When a multi-task arc will open a destructive path, add an EXPLICIT gate that the opening task removes deliberately, so it cannot open early.

---

## [2026-07-19-210439] The disclosure parser is a deliberately dumb line-scanner (skips <!-- --> comments, not code fences)

**Context**: M4 unifies conventions into the ## -section entry model, where ## collides far more than ## [. Fence detection was considered and rejected.

**Lesson**: Fence detection (nested/tilde/indented fences) is a rabbit hole in a parser that does destructive moves. Safety comes from byte-conservation (a mis-cut is a reversible mis-grouping, not loss) plus the human-gated plan review (a false entry shows up at inspect).

**Application**: Don't add fence-awareness to the disclosure scanner; a stray ## inside a fenced code block is fixed at authoring time (rewrite that one line), never by the tool.

---

## [2026-07-19-210439] Measurement gates surface a real bug in every disclosure milestone

**Context**: pd-m1's add-path measurement gate flushed the insert.AfterHeader tail-truncation bug; pd-m3's T17 real rollout surfaced the DECISIONS-template comment-parse bug (a ## [ inside <!-- --> read as staging, so Validate returned ErrStagingUnparsable once the fold emptied real staging).

**Lesson**: For the disclosure work, driving the REAL fold — not unit tests — is where bugs hide. The measurement gate is not ceremony.

**Application**: Treat pd-m4 T20 (drive the digest on the real CONVENTIONS.md) as the authoritative gate; expect it to find something; verify conservation + all invariants on the real file before T21's human-gated rollout.

---

---

## Themes

- go-idioms-and-structure — Go language & package layout: import cycles, constant/helper placement, sync.Once smells, toolchain/build-tag pitfalls, error sentinels, test isolation → [go-idioms-and-structure](learnings/go-idioms-and-structure.md)
- audit-lint-compliance — Mechanical enforcement: the internal/audit AST gauntlet, compliance tests as command style guide, cmd/ purity, docstring floors, gosec triggers, linter quirks → [audit-lint-compliance](learnings/audit-lint-compliance.md)
- cli-command-design — cobra/CLI surface conventions: bare-invocation existence probing, legacyArgs silent-success on groups, commands named after input not output → [cli-command-design](learnings/cli-command-design.md)
- hooks-and-integration — Hook & editor-integration mechanics: output channels, key names, compliance wiring, Cursor/OpenCode plugins, project-local hooks, notify webhooks → [hooks-and-integration](learnings/hooks-and-integration.md)
- context-model-and-state — Context-dir resolution & on-disk state: single-source ContextDir, tombstones/logs/fs hygiene, managed-block guards, handover metadata, hub raft-lite → [context-model-and-state](learnings/context-model-and-state.md)
- text-markdown-serialization — Text/markdown/JSON hazards: line-sweep corruption, insert-helper tail drops, diacritic stripping, title-case brands, RFC3339 sort, key-dropping round-trips → [text-markdown-serialization](learnings/text-markdown-serialization.md)
- model-context-window — LLM model->context-window mapping: silent 200k fallback for new families, ordered prefix matching, doctor token_budget vs context_window → [model-context-window](learnings/model-context-window.md)
- git-and-signing — git CLI wrapping quirks: filter-branch leftover refs, hunk-level feature carving, GPG signing from non-TTY (pinentry) → [git-and-signing](learnings/git-and-signing.md)
- environment-and-platform — Cross-platform gotchas: macOS minimal PATH & /var symlink, Tauri rustc floor, GitNexus in-container pruning, host-pressure alerting → [environment-and-platform](learnings/environment-and-platform.md)
- skills-agents-and-tasks — Skill/agent workflow: skill lifecycle & shipping, agent context-loading/routing/behavior, task exit-criteria, ctx-dream design, refactor mechanics → [skills-agents-and-tasks](learnings/skills-agents-and-tasks.md)
- docs-and-templates — Docs/template/asset drift: tracked build output, magic-string discipline, contributor-PR reintroduction, docs/code divergence, redundant scripts → [docs-and-templates](learnings/docs-and-templates.md)
- editorial-and-product-signals — KB editorial & product signals: single topic-enumeration site, decision-recording, creator-confusion as doc-quality signal, editorial leverage, IDE-is-the-UI → [editorial-and-product-signals](learnings/editorial-and-product-signals.md)
