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
