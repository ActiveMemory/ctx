# Learnings

<!-- INDEX:START -->
| Date | Learning |
|------|--------|
| 2026-03-14 | Stderr error messages are user-facing text that belongs in assets |
| 2026-03-14 | Subagents rename packages and modify unrelated files without being asked |
| 2026-03-14 | Hardcoded _alt suffixes create implicit language favoritism |
| 2026-03-14 | Subagents reorganize file structure without being asked |
| 2026-03-14 | Internal skill rename requires updates across 6+ layers |
| 2026-03-13 | Skills without a trigger mechanism are dead code |
| 2026-03-13 | Variable shadowing causes cascading test failures after package splits |
| 2026-03-13 | sync-why mechanism existed but was not wired to build |
| 2026-03-13 | Linter reverts import-only edits when references still use old package |
| 2026-03-12 | Project-root files vs context files are distinct categories |
| 2026-03-12 | Constants belong in their domain package not in god objects |
| 2026-03-07 | Always search for existing constants before adding new ones |
| 2026-03-07 | SafeReadFile requires split base+filename paths |
| 2026-03-06 | Spawned agents reliably create new files but consistently fail to delete old ones — always audit for stale files, duplicate function definitions, and orphaned imports after agent-driven refactoring |
| 2026-03-06 | Import cycle avoidance: when package A imports package B for logic, B must own shared types — A aliases them. entry imports add/core for insert logic, so add/core owns EntryParams and entry aliases it as entry.Params |
| 2026-03-06 | Stale directory inodes cause invisible files over SSH |
| 2026-03-06 | Stats sort uses string comparison on RFC3339 timestamps with mixed timezones |
| 2026-03-06 | Claude Code supports PreCompact and SessionStart hooks that ctx does not use |
| 2026-03-06 | nolint:goconst for trivial values normalizes magic strings |
| 2026-03-06 | Package-local err.go files invite broken windows from future agents |
| 2026-03-05 | State directory accumulates silently without auto-prune |
| 2026-03-05 | Global tombstones suppress hooks across all sessions |
| 2026-03-05 | Claude Code has two separate memory systems behind feature flags |
| 2026-03-05 | Blog post editorial feedback is higher-leverage than drafting |
| 2026-03-04 | CONSTITUTION hook compliance is non-negotiable — don't work around it |
| 2026-03-04 | nolint:errcheck in tests normalizes unchecked errors for agents |
| 2026-03-04 | golangci-lint v2 ignores inline nolint directives for some linters |
| 2026-03-02 | Hook message registry test enforces exhaustive coverage of embedded templates |
| 2026-03-02 | Existing Projects is ambiguous framing for migration notes |
| 2026-03-02 | Claude Code JSONL model ID does not distinguish 200k from 1M context |
| 2026-03-01 | Gosec G306 flags test file WriteFile with 0644 permissions |
| 2026-03-01 | Converting PersistentPreRun to PersistentPreRunE changes exit behavior |
| 2026-03-01 | Key path changes ripple across 15+ doc files and 2 skills |
| 2026-03-01 | Test HOME isolation is required for user-level path functions |
| 2026-03-01 | Skill enhancement is a documentation-heavy operation across 10+ files |
| 2026-03-01 | Task descriptions can be stale in reverse — implementation done but task not marked complete |
| 2026-03-01 | Elevating private skills requires synchronized updates across 6 layers |
| 2026-03-01 | Model-to-window mapping requires ordered prefix matching |
| 2026-03-01 | Removing embedded asset directories requires synchronized cleanup across 5+ layers |
| 2026-03-01 | Absorbing shell scripts into Go commands creates a discoverability gap |
| 2026-03-01 | TASKS.md template checkbox syntax inside HTML comments is parsed by RegExTaskMultiline |
| 2026-03-01 | Hook logs had no rotation; event log already did |
| 2026-02-28 | ctx pad import, ctx pad export, and ctx system resources make three hack scripts redundant |
| 2026-02-28 | Getting-started docs assumed Claude Code as the only agent |
| 2026-02-28 | Plugin reload script must rebuild cache, not just delete it |
| 2026-02-27 | site/ directory must be committed with docs changes |
| 2026-02-27 | Doctor token_budget vs context_window confusion |
| 2026-02-27 | Drift detector false positives on illustrative code examples |
| 2026-02-27 | Context injection and compliance strategy (consolidated) |
| 2026-02-26 | Webhook silence after ctxrc profile swap is the most common notify debugging red herring |
| 2026-02-26 | Documentation drift and auditing (consolidated) |
| 2026-02-26 | Agent context loading and task routing (consolidated) |
| 2026-02-26 | Go testing patterns (consolidated) |
| 2026-02-26 | PATH and binary handling (consolidated) |
| 2026-02-26 | Task management and exit criteria (consolidated) |
| 2026-02-26 | Agent behavioral patterns (consolidated) |
| 2026-02-26 | Hook compliance and output routing (consolidated) |
| 2026-02-26 | ctx add and decision recording (consolidated) |
| 2026-02-24 | CLI tools don't benefit from in-memory caching of context files |
| 2026-02-22 | Hook behavior and patterns (consolidated) |
| 2026-02-22 | UserPromptSubmit hook output channels (consolidated) |
| 2026-02-22 | Linting and static analysis (consolidated) |
| 2026-02-22 | Permission and settings drift (consolidated) |
| 2026-02-22 | Gitignore and filesystem hygiene (consolidated) |
| 2026-02-19 | Feature can be code-complete but invisible to users |
| 2026-01-28 | IDE is already the UI |
<!-- INDEX:END -->

---

## [2026-03-14-180903] Stderr error messages are user-facing text that belongs in assets

**Context**: Added fmt.Fprintf(os.Stderr) error reporting to event log, initially with inline strings

**Lesson**: Any string that reaches the user, including stderr warnings, routes through assets.TextDesc() for i18n readiness

**Application**: When adding stderr output, create text.yaml entries and asset keys first

---

## [2026-03-14-180902] Subagents rename packages and modify unrelated files without being asked

**Context**: ET agent renamed internal/eventlog/ to internal/log/ and modified 5+ caller files outside the internal/err/ scope

**Lesson**: Always diff agent output against HEAD to catch scope creep before building; revert unrelated changes immediately

**Application**: After any agent-driven refactor, run git diff --name-only HEAD and revert anything outside the intended scope before testing

---

## [2026-03-14-131202] Hardcoded _alt suffixes create implicit language favoritism

**Context**: Session parser had session_prefix_alt hardcoding Turkish as a special case alongside English default

**Lesson**: Naming a constant _alt and hardcoding one non-English language as a built-in default discriminates by giving that language special status. The pattern doesn't scale (alt_2? alt_3?) and signals that adding languages requires code changes.

**Application**: When a feature needs multi-value support, use configurable lists from the start — not hardcoded pairs with _alt suffixes. Default to a single canonical value; all extensions are user-configured equally.

---

## [2026-03-14-110750] Subagents reorganize file structure without being asked

**Context**: Asked subagent to replace os.ReadFile callsites with validation wrappers. It also moved functions to new files, renamed them (ReadUserFile to SafeReadUserFile), and created a new internal/io package.

**Lesson**: Subagents optimize for clean results and will restructure beyond stated scope — moved, renamed, and split files without being asked.

**Application**: After subagent-driven refactors, always verify file organization matches intent. Audit for moved, renamed, and split files, not just the requested callsite changes.

---

## [2026-03-14-093757] Internal skill rename requires updates across 6+ layers

**Context**: Renamed ctx-alignment-audit to _ctx-alignment-audit. The allow list test in embed_test.go failed because it iterates all bundled skills and expects each in the allow list.

**Lesson**: The allow list test needed a strings.HasPrefix(_) skip for internal skills. This was not obvious until tests ran.

**Application**: When converting public to internal skills, audit: allow.txt, embed_test.go allow list test, reference/skills.md, all recipe docs referencing the skill, contributing.md dev-only skills table, and permissions docs.

---

## [2026-03-13-223110] Skills without a trigger mechanism are dead code

**Context**: ctx-context-monitor was a skill documenting how to respond to hook output, but no hook or agent ever loaded it. The hook output already contained sufficient instructions.

**Lesson**: A skill only enters the agent context when explicitly invoked via /skill-name. If the description says not user-invocable and no mechanism loads it automatically, it is unreachable.

**Application**: Audit skills for reachability. If nothing triggers the skill, either add a trigger or delete it.

---

## [2026-03-13-223108] Variable shadowing causes cascading test failures after package splits

**Context**: Large refactoring moved constants from monolithic config package to sub-packages (dir, entry, file). Test files had local variables named dir, entry, file that shadowed the new package imports.

**Lesson**: When splitting a package, audit test files for local variable names that collide with new package names. dir, file, entry are especially common Go variable names.

**Application**: Before committing a package split, run go test ./... and check for undefined errors caused by variable shadowing.

---

## [2026-03-13-151952] sync-why mechanism existed but was not wired to build

**Context**: assets/why/ had drifted from docs/ — the sync targets existed in the Makefile but build did not depend on sync-why

**Lesson**: Freshness checks that are not in the critical path will be forgotten. Wire them as build prerequisites, not optional audit steps

**Application**: Any derived or copied asset should be a prerequisite of build, not just audit

---

## [2026-03-13-151951] Linter reverts import-only edits when references still use old package

**Context**: Moving tpl_entry.go from config/entry to assets — linter kept reverting the import change

**Lesson**: When moving constants between packages, change imports and all references in a single atomic write (use Write not incremental Edit), so the linter never sees an inconsistent state

**Application**: For future package migrations, use full file rewrites when a linter is active

---

## [2026-03-12-133008] Project-root files vs context files are distinct categories

**Context**: Tried moving ImplementationPlan constant to config/ctx assuming it was a context file

**Lesson**: Files created by ctx init in the project root (Makefile, IMPLEMENTATION_PLAN.md) are scaffolding, not context files loaded via ReadOrder. They belong in config/file, not config/ctx

**Application**: Before moving a file constant, check whether it is in ReadOrder (context) or created by init (project-root)

---

## [2026-03-12-133007] Constants belong in their domain package not in god objects

**Context**: file.go held agent scoring constants, budget percentages, cooldown durations — none related to file config

**Lesson**: When a constant is only used by one domain (e.g. agent scoring), it should live in that domain's config package

**Application**: Check callers before placing constants; if all callers are in one domain, the constant belongs there

---

## [2026-03-07-221151] Always search for existing constants before adding new ones

**Context**: Added ExtJsonl constant to config/file.go but ExtJSONL already existed with the same value, causing a duplicate

**Lesson**: Grep for the value (e.g. '.jsonl') across config/ before creating a new constant — naming variations (camelCase vs ALLCAPS) make duplicates easy to miss

**Application**: Before adding any new constant to internal/config, search by value not just by name

---

## [2026-03-07-221148] SafeReadFile requires split base+filename paths

**Context**: During system/core cleanup, persistence.go passed a full path to validation.SafeReadFile which expects (baseDir, filename) separately

**Lesson**: Use filepath.Dir(path) and filepath.Base(path) to split full paths when adapting os.ReadFile calls to SafeReadFile

**Application**: When converting os.ReadFile to SafeReadFile, always check whether the existing code has a full path or separate components

---

## [2026-03-06-200319] Spawned agents reliably create new files but consistently fail to delete old ones — always audit for stale files, duplicate function definitions, and orphaned imports after agent-driven refactoring

**Context**: Multiple agent batches across cmd/ restructuring, color removal, and flag externalization left stale files, duplicate run.go, and unupdated parent imports

**Lesson**: Agent cleanup is a known gap — budget 5-10 minutes for post-agent audit per batch

**Application**: After every agent batch: grep for stale package declarations, check parent imports point to cmd/root not cmd/, verify old files are deleted

---

## [2026-03-06-200237] Import cycle avoidance: when package A imports package B for logic, B must own shared types — A aliases them. entry imports add/core for insert logic, so add/core owns EntryParams and entry aliases it as entry.Params

**Context**: Extracting entry.Params as a standalone struct in internal/entry created a cycle because entry/write.go imports add/core for AppendEntry

**Lesson**: The package that provides implementation logic must own the types; the facade package aliases them

**Application**: When extracting shared types from implementation packages, check the import direction first — the type lives where the logic lives

---

## [2026-03-06-141506] Stale directory inodes cause invisible files over SSH

**Context**: Files created by Claude Code hooks were visible inside the VM but not from the SSH terminal

**Lesson**: If a directory is recreated (e.g. by auto-prune), an SSH shell holding the old directory inode will not see new files — ls returns no such file even though cat with the full path works from other shells

**Application**: After ctx system prune or any state directory recreation, SSH sessions need cd-dot or re-login to pick up the new inode

---

## [2026-03-06-141504] Stats sort uses string comparison on RFC3339 timestamps with mixed timezones

**Context**: ctx system stats showed only old sessions, hiding the current one

**Lesson**: RFC3339 string comparison breaks when entries mix UTC (Z) and offset (-08:00) formats — 13:00-08:00 sorts before 18:00Z lexicographically despite being later in absolute time

**Application**: Always parse to time.Time before comparing RFC3339 timestamps; never rely on lexicographic sort

---

## [2026-03-06-184820] Claude Code supports PreCompact and SessionStart hooks that ctx does not use

**Context**: context-mode proves both hooks work in production across 5 platforms

**Lesson**: ctx's hook architecture only uses UserPromptSubmit, PreToolUse, and PostToolUse — two lifecycle events are untapped

**Application**: PreCompact snapshot plus SessionStart re-injection would eliminate post-compaction disorientation without any new persistence layer since ctx agent already generates the content

---

## [2026-03-06-050126] nolint:goconst for trivial values normalizes magic strings

**Context**: Found 5 callsites suppressing goconst with nolint:goconst trivial user input check for y/yes comparisons

**Lesson**: Suppressing a linter with a trivial justification sets precedent for other agents to do the same. The fix (two constants) costs less than the accumulated tech debt.

**Application**: Use config.ConfirmShort/config.ConfirmLong instead of suppressing goconst. Prefer constants over nolint directives.

---

## [2026-03-06-050125] Package-local err.go files invite broken windows from future agents

**Context**: Found err.go files in 5 CLI packages with heavily duplicated error constructors (errFileWrite, errMkdir, errZensicalNotFound repeated across packages)

**Lesson**: Centralizing errors in internal/err eliminates duplication and prevents agents from continuing the pattern of adding local err.go files when they see one exists

**Application**: New error constructors go to internal/err/errors.go. No err.go files in CLI packages.

---

## [2026-03-05-205422] State directory accumulates silently without auto-prune

**Context**: Found 234 files in .context/state/ from weeks of sessions with no cleanup mechanism

**Lesson**: Session tombstones are write-only. Without auto-prune, the state directory grows unbounded. Added autoPrune(7) to context-load-gate so cleanup happens once per session at startup.

**Application**: Auto-prune is now wired into session start via context-load-gate. Manual prune still available via ctx system prune for aggressive cleanup.

---

## [2026-03-05-205419] Global tombstones suppress hooks across all sessions

**Context**: Memory drift nudge used memory-drift-nudged with no session ID in filename

**Lesson**: Any tombstone file intended to be session-scoped must include the session ID in its filename, otherwise it suppresses across all concurrent and future sessions. Use the UUID pattern so prune can clean them up.

**Application**: Audit all tombstone files for session-scoping; fixed memory-drift, but backup-reminded, ceremony-reminded, check-knowledge, journal-reminded, version-checked, ctx-wrapped-up still have this bug

---

## [2026-03-05-042157] Claude Code has two separate memory systems behind feature flags

**Context**: Filesystem and behavioral analysis of Claude Code v2.1.69

**Lesson**: Claude Code has two separate memory systems behind feature flags. Auto memory writes MEMORY.md to disk (user-visible, toggleable via settings). Session memory is a separate background extraction pipeline with compaction and team sync (push/pull model). The two systems serve different purposes and are independently feature-flagged.

**Application**: ctx memory bridge targets auto memory (MEMORY.md on disk). Session memory is API-side and not directly accessible. Full findings in ideas/claude-code-project-directory-structure.md.

---

## [2026-03-05-023941] Blog post editorial feedback is higher-leverage than drafting

**Context**: Draft of Agent Memory Is Infrastructure was publication-quality on first pass; user editorial feedback (structural emphasis, rhetorical sharpening, amnesia/archaeology bridge) elevated it significantly more than initial generation

**Lesson**: For narrative content, the first draft captures the argument; the editorial pass captures the voice. Both are necessary but the editorial pass has disproportionate impact on quality.

**Application**: For future blog posts, invest more in the editorial cycle (structural feedback then targeted refinements) rather than trying to nail voice on first generation.

---

## [2026-03-04-105239] CONSTITUTION hook compliance is non-negotiable — don't work around it

**Context**: After make build, ran ./ctx deps --help which was blocked by block-non-path-ctx. Instead of asking user to install, tried cp ctx ~/bin/ — escalating workarounds.

**Lesson**: When a hook blocks an action, the correct response is to follow the hook's instruction (ask the user to sudo make install), not to find creative bypasses.

**Application**: Always ask the user to install when testing a freshly built binary. Never attempt alternative install paths to circumvent a hook.

---

## [2026-03-04-040211] nolint:errcheck in tests normalizes unchecked errors for agents

**Context**: User flagged that suppressing errcheck in tests teaches the agent to spread the pattern to production code

**Lesson**: Broken-window theory applies to lint suppressions. Agents learn from test code patterns. Use _ = f.Close() in a closure or check errors with t.Fatal — never suppress with nolint.

**Application**: Handle all errors in test code the same as production: t.Fatal(err) for setup, defer func() { _ = f.Close() }() for best-effort cleanup.

---

## [2026-03-04-040209] golangci-lint v2 ignores inline nolint directives for some linters

**Context**: nolint:errcheck and nolint:gosec comments were present but golangci-lint v2 still reported violations

**Lesson**: In golangci-lint v2, use config-level exclusions.rules for gosec patterns (G204, G301, G306) rather than relying on inline nolint directives. For errcheck, fix the code instead of suppressing.

**Application**: When adding new lint suppressions, prefer config-level rules for gosec false positives on safe paths/args; never suppress errcheck — handle the error.

---

## [2026-03-02-165039] Hook message registry test enforces exhaustive coverage of embedded templates

**Context**: Adding billing.txt to embedded assets without a registry entry caused TestRegistryCoversAllEmbeddedFiles to fail immediately

**Lesson**: Every new .txt file under internal/assets/hooks/messages/ must have a corresponding entry in registry.go — the test acts as an exhaustive bidirectional check

**Application**: When adding new hook message variants, update the registry entry before running tests

---

## [2026-03-02-123613] Existing Projects is ambiguous framing for migration notes

**Context**: A doc admonition said Existing Projects: if you have an older key at X, it auto-migrates. Every project is existing once installed — the framing does not tell you how far behind you need to be.

**Lesson**: Version-anchored framing (Key Folder Change v0.7.0+) is clearer than relative framing (Existing Projects, Legacy). State the version boundary and the concrete action.

**Application**: When writing migration notes, anchor to a version number and give copy-pasteable commands, not vague auto-handled assurances.

---

## [2026-03-02-005217] Claude Code JSONL model ID does not distinguish 200k from 1M context

**Context**: Heartbeat hook was reporting 16% usage at 162k tokens because it assumed claude-opus-4-6 always has 1M context window

**Lesson**: The JSONL model field is identical for both variants (both report claude-opus-4-6). The 1M context requires a beta header, not a different model ID. The user's model selection is stored in ~/.claude/settings.json with a [1m] suffix when 1M is active.

**Application**: Auto-detect context window from ~/.claude/settings.json model field containing [1m]. Default to 200k for all Claude models. The .ctxrc context_window setting is a no-op for Claude Code users.

---

## [2026-03-01-222739] Gosec G306 flags test file WriteFile with 0644 permissions

**Context**: New tests used os.WriteFile(..., 0o644) for temp context files; lint flagged all three occurrences

**Lesson**: Gosec enforces 0600 max on WriteFile even in test code. Use 0o600 for test temp files

**Application**: Default to 0o600 for os.WriteFile in tests; only use wider permissions when testing permission behavior specifically

---

## [2026-03-01-222738] Converting PersistentPreRun to PersistentPreRunE changes exit behavior

**Context**: Boundary violation test used subprocess pattern because original code called os.Exit(1)

**Lesson**: With PersistentPreRunE, errors propagate through Cobra Execute() return — no os.Exit call. Subprocess-based tests that expected exit codes need converting to direct error assertions

**Application**: When converting PreRun to PreRunE in Cobra commands, audit all tests that relied on os.Exit behavior

---

## [2026-03-01-194147] Key path changes ripple across 15+ doc files and 2 skills

**Context**: Updating docs for the .context/.ctx.key → ~/.local/ctx/keys/ → ~/.ctx/.ctx.key migrations

**Lesson**: Key path changes have a long documentation tail — recipes, references, getting-started, operations, CLI docs, and skills all carry path references. The worktree behavior flip (limitation to works automatically) was the highest-value change per line edited. Simplifying from per-project slugs to a single global key eliminated more code and docs than the original migration added.

**Application**: When moving a file path that appears in user-facing docs, grep broadly (not just code) and budget for 15+ file touches

---

## [2026-03-01-161459] Test HOME isolation is required for user-level path functions

**Context**: After adding ~/.ctx/.ctx.key as global key location, test suites wrote real files to the developer home directory

**Lesson**: Any code that uses os.UserHomeDir() needs t.Setenv(HOME, tmpDir) in tests — especially test helpers called by many tests (like setupEncrypted and helper)

**Application**: When adding features that write to user-level paths (~/.ctx/, ~/.config/), always add HOME isolation to test setup functions first

---

## [2026-03-01-144544] Skill enhancement is a documentation-heavy operation across 10+ files

**Context**: Enhancing /ctx-journal-enrich-all to handle export-if-needed touched the skill, hook messages, fallback strings, 5 doc files, 2 Makefiles, and TASKS.md

**Lesson**: Skill behavior changes ripple through hook messages, fallback strings in Go code, doc descriptions, and Makefile hints — all must stay synchronized

**Application**: When modifying a skill's scope, grep for its name across the entire repo and update every description, not just the skill file itself

---

## [2026-03-01-133014] Task descriptions can be stale in reverse — implementation done but task not marked complete

**Context**: ctx recall sync task said 'command is not registered in Cobra' but the code was fully wired and all tests passed. The task description was stale.

**Lesson**: Tasks can become stale in the opposite direction from docs: implementation gets completed but the task is not updated. Always verify with ctx <cmd> --help before assuming work remains.

**Application**: Before starting implementation on a 'code exists but not wired' task, run the command first to check if it already works.

---

## [2026-03-01-125807] Elevating private skills requires synchronized updates across 6 layers

**Context**: Promoted 6 _ctx-* skills to bundled ctx-* plugin skills

**Lesson**: Moving a skill from .claude/skills/ to internal/assets/claude/skills/ touches: (1) SKILL.md frontmatter name field, (2) internal cross-references between skills (slash command paths), (3) external cross-references in other skills and docs, (4) embed_test.go expected skill list, (5) recipe and reference docs that mention the old name, (6) plugin cache rebuild (`hack/plugin-reload.sh`) + session restart — Claude Code snapshots skills from `~/.claude/plugins/cache/` at startup, so new skills are invisible until the cache is refreshed. Also clean stale underscore-prefixed `Skill(_ctx-*)` entries from `.claude/settings.local.json`.

**Application**: When promoting future skills, use grep -r /_ctx-{name} across the whole tree before declaring done. After code changes, run plugin-reload.sh and restart the session to verify the skill appears in autocomplete.

---

## [2026-03-01-124921] Model-to-window mapping requires ordered prefix matching

**Context**: Implementing modelContextWindow() for the three-tier context window fallback. Claude model IDs use nested prefixes (claude-sonnet-4-5 vs claude-sonnet-4-20250514).

**Lesson**: A switch with ordered HasPrefix cases (most specific first) is cleaner and safer than iterating separate prefix lists. The catch-all 'claude-*' returns 200k for unrecognized Claude models.

**Application**: When adding new model families to modelContextWindow() in session_tokens.go, add the most specific prefix first to avoid shadowing shorter prefixes.

---

## [2026-03-01-112538] Removing embedded asset directories requires synchronized cleanup across 5+ layers

**Context**: Deleting .context/tools/ deployment touched embed directive, asset functions, init logic, tests, config constants, Makefile targets, and docs — missing any one layer leaves dead code or build failures.

**Lesson**: Embedded asset removal is a cross-cutting concern: embed directive → accessor functions → callers → tests → config constants → build targets → documentation. Work outward from the embed.

**Application**: When removing an embedded asset category, use the grep-first approach (search for all references to the accessor functions and constants) before deleting anything.

---

## [2026-03-01-102232] Absorbing shell scripts into Go commands creates a discoverability gap

**Context**: Deleted make backup/backup-global/backup-all and make rc-dev/rc-base/rc-status targets when absorbing into ctx system backup and ctx config switch. The Makefile served as self-documenting discovery (make help).

**Lesson**: When eliminating Makefile targets, the CLI reference page alone is not sufficient — contributor-facing docs (contributing.md) and command catalogs (common-workflows.md) must gain explicit entries to compensate.

**Application**: For future hack/ absorptions (e.g. pad-import-ideas.sh, context-watch.sh), audit contributing.md, common-workflows.md CLI-Only table, and the CLI index page as part of the absorption checklist.

---

## [2026-03-01-095709] TASKS.md template checkbox syntax inside HTML comments is parsed by RegExTaskMultiline

**Context**: Template had example checkboxes (- [x], - [ ]) in HTML comments that the line-based regex matched as real tasks, causing TestArchiveCommand_NoCompletedTasks to fail

**Lesson**: RegExTaskMultiline is line-based and has no awareness of HTML comment blocks — checkbox-like patterns inside comments get counted as real tasks

**Application**: Use backtick-quoted or indented references instead of actual checkbox syntax in template comments. When adding examples to TASKS.md templates, avoid patterns that match regExTaskPattern

---

## [2026-03-01-092611] Hook logs had no rotation; event log already did

**Context**: Investigated .context/logs/ and .context/state/ file management

**Lesson**: eventlog already rotates at 1MB with one previous generation. logMessage() in state.go was pure append-only with no size check.

**Application**: When adding new log sinks, follow the established rotation pattern (size-based, single previous generation)

---

## [2026-02-28-184758] ctx pad import, ctx pad export, and ctx system resources make three hack scripts redundant

**Context**: Audited hack/ scripts against ctx CLI surface

**Lesson**: As ctx CLI grew, several hack scripts became wrappers around built-in commands (pad-import.sh -> ctx pad import, pad-export-blobs.sh -> ctx pad export, resource-watch.sh -> watch -n5 ctx system resources)

**Application**: Periodically audit hack/ for scripts that ctx has absorbed

---

## [2026-02-28-184647] Getting-started docs assumed Claude Code as the only agent

**Context**: The installation section opened with 'A full ctx installation has two parts' — binary + Claude Code plugin — leaving non-Claude-Code users without a clear path

**Lesson**: Installation docs should lead with the universal requirement (the binary) and present agent-specific integration as conditional

**Application**: When writing docs for multi-tool projects, frame the common denominator first, then branch by tool

---

## [2026-02-28-150701] Plugin reload script must rebuild cache, not just delete it

**Context**: hack/plugin-reload.sh was deleting ~/.claude/plugins/cache/activememory-ctx/ without repopulating it. Claude Code's installed_plugins.json still referenced the cache path, so the plugin appeared enabled but hooks.json was missing — all plugin hooks silently stopped firing.

**Lesson**: Claude Code snapshots plugin hooks from the cache directory at session startup. If the cache is deleted, plugin hooks vanish silently with no error. The reload script must rebuild the cache from source assets (internal/assets/claude/) after clearing it, and warn that a session restart is required.

**Application**: Always rebuild the plugin cache in hack/plugin-reload.sh. When debugging hooks that don't fire, check ~/.claude/plugins/cache/ first — a missing hooks.json is the most likely cause.

---

## [2026-02-27-231228] site/ directory must be committed with docs changes

**Context**: The site/ directory contains generated HTML served directly from the repo (no CI build step). Multiple sessions have committed docs/ changes without the corresponding site/ output, or ignored site/ as 'generated noise'.

**Lesson**: site/ is intentionally tracked in git — there is no GitHub Pages workflow or CI step to build it. When docs change, the regenerated site/ HTML must be staged and committed alongside the source.

**Application**: Always git add site/ when committing changes under docs/. Never gitignore site/.

---

## [2026-02-27-230741] Doctor token_budget vs context_window confusion

**Context**: ctx doctor reported context size against token_budget (8k) instead of context_window (200k), making 22k tokens look alarming.

**Lesson**: token_budget (ctx agent output trim target) and context_window (model capacity) serve different purposes. Health checks about context fitting should use context_window, with warning threshold proportional (e.g., 20% of window).

**Application**: Doctor now uses rc.ContextWindow() with 20% threshold and shows per-file token breakdown for actionable insight into which files are heavy.

---

## [2026-02-27-230738] Drift detector false positives on illustrative code examples

**Context**: ctx drift flagged 23 warnings for backtick-quoted paths in CONVENTIONS.md and ARCHITECTURE.md that were prose examples (loader.go, session/run.go, sync.Once), not real file references.

**Lesson**: Path reference detection should verify the top-level directory exists on disk before flagging. Bare filenames and paths under non-existent directories are almost always examples in documentation.

**Application**: The fix checks os.Stat(topDir) on the first path component. Future drift checks on documentation-heavy files should use the same heuristic.

---

## [2026-02-27-002830] Context injection and compliance strategy (consolidated)

**Consolidated from**: 3 entries (2026-02-26)

- Verbal summaries with linked diagram files cut ARCHITECTURE.md from ~12K to ~3.8K tokens. Extract diagrams to linked files outside FileReadOrder; keep prose summaries inline. The 4-chars-per-token estimator is accurate — optimize content, not the estimator.
- Soft instructions have a ~75-85% compliance ceiling because "don't apply judgment" is itself evaluated by judgment. When 100% compliance is required, don't instruct — inject via `additionalContext`. Reserve soft instructions for ~80% acceptable compliance.
- Once ~7K tokens are auto-injected (fait accompli), the agent's rationalization inverts from "skip to save effort" to "marginal cost is trivial." Front-load highest-value content as injection, then use sunk cost to motivate on-demand reads for the remainder.

---

## [2026-02-26-003854] Webhook silence after ctxrc profile swap is the most common notify debugging red herring

**Context**: Spent time investigating why webhooks weren't firing — checked binary version, hook configs, notify.Send internals. Actual cause was .ctxrc swapped to prod profile (notify commented out) earlier in session.

**Lesson**: When webhooks stop, check .ctxrc profile first (`ctx config status`). Also: not all tool uses trigger webhook-sending hooks — Read only triggers context-load-gate (one-shot) and ctx agent (no webhook). qa-reminder requires Edit matcher.

**Application**: Before debugging notify internals, run `ctx config status` and verify the event would actually match a hook with notify.Send.

---

## [2026-02-26-100000] Documentation drift and auditing (consolidated)

**Consolidated from**: 6 entries (2026-01-29 to 2026-02-24)

- CLI reference docs can outpace implementation: ctx remind had no CLI, ctx recall sync had no Cobra wiring, key file naming diverged between docs and code. Always verify with `ctx <cmd> --help` before releasing docs.
- Structural doc sections (project layouts, command tables, skill counts) drift silently. Add `<!-- drift-check: <shell command> -->` markers above any section that mirrors codebase structure.
- Agent sweeps for style violations are unreliable (8 found vs 48+ actual). Always follow agent results with targeted grep and manual classification.
- ARCHITECTURE.md missed 4 core packages and 4 CLI commands. The /ctx-drift skill catches stale paths but not missing entries — run /ctx-architecture after adding new packages or commands.
- Documentation audits must compare against known-good examples and pattern-match for the COMPLETE standard, not just presence of any comment.
- Dead link checking belongs in /consolidate's check list (check 12), not as a standalone concern. When a new audit concern emerges, check if it fits an existing audit skill first.

---

## [2026-02-26-100002] Agent context loading and task routing (consolidated)

**Consolidated from**: 5 entries (2026-01-20 to 2026-01-25)

- `ctx agent` is optimized for task execution (filters pending tasks, surfaces constitution, token-budget aware). Manual file reading is better for exploratory/memory questions (session history, timestamps, completed tasks).
- On "Do you remember?" questions, immediately read .context/ files and run `ctx recall list --limit 5`. Never ask "would you like me to check?" — that is the obvious intent.
- .context/ is NOT a Claude Code primitive. Only CLAUDE.md and .claude/settings.json are auto-loaded. The .context/ directory requires a hook or explicit CLAUDE.md instruction to be discovered.
- Orchestrator (IMPLEMENTATION_PLAN.md) and agent (.context/TASKS.md) task lists must be separate. The orchestrator says "check your mind" — it doesn't maintain a parallel ledger.
- Only CLAUDE.md is auto-loaded by Claude Code. Projects using ctx should rely on the CLAUDE.md -> AGENT_PLAYBOOK.md chain, not AGENTS.md.

---

## [2026-02-26-100005] Go testing patterns (consolidated)

**Consolidated from**: 7 entries (2026-01-19 to 2026-02-26)

- Compiler-driven refactoring misses test files: `go build ./...` catches production callsite breaks but not test files. Always run `go test ./...` after signature changes.
- All runCmd() returns must be consumed in tests: even setup calls need `_, _ = runCmd(...)` to satisfy errcheck.
- Set `color.NoColor = true` in a package-level init function to disable ANSI codes for CLI test string assertions.
- Recall CLI tests isolate via HOME env var: `t.Setenv("HOME", tmpDir)` with `.claude/projects/` structure gives full isolation from real session data.
- `formatDuration` accepts an interface with a Minutes method, not time.Duration directly. Use a stubDuration struct for testing.
- CI tests need `CTX_SKIP_PATH_CHECK=1` env var because init checks if ctx is in PATH.
- CGO must be disabled for ARM64 Linux (`CGO_ENABLED=0`) — CGO causes cross-compilation issues with `-m64` flag.

---

## [2026-02-26-100006] PATH and binary handling (consolidated)

**Consolidated from**: 3 entries (2026-01-21 to 2026-02-17)

- Always use `ctx` from PATH, never `./dist/ctx-linux-arm64` or `go run ./cmd/ctx`. Check `which ctx` if unsure.
- Hooks must use PATH, not hardcoded paths. `ctx init` checks if ctx is in PATH before proceeding. Tests can skip with `CTX_SKIP_PATH_CHECK=1`.
- Agent must never place binaries in any bin directory (not via cp, mv, or go install). Build with `make build`, then ask the user to run the privileged install step. Hooks in block-dangerous-commands.sh enforce this.

---

## [2026-02-26-100007] Task management and exit criteria (consolidated)

**Consolidated from**: 4 entries (2026-01-21 to 2026-02-17)

- Specs get lost without cross-references from TASKS.md. Three-layer defense: (1) playbook instruction, (2) spec reference in Phase header, (3) bold breadcrumb in first task.
- Subtask completion is implementation progress, not delivery. Parent tasks should have explicit deliverables; don't close until deliverable is verified.
- Exit criteria must include verification: integration tests (binary executes correctly), coverage targets, and smoke tests. "All tasks checked off" does not equal "implementation works."
- Reports graduate to ideas/done/ only after all items are tracked or resolved. Cross-reference every item against TASKS.md and the codebase before moving.

---

## [2026-02-26-100008] Agent behavioral patterns (consolidated)

**Consolidated from**: 5 entries (2026-01-25 to 2026-02-22)

- Interaction pattern capture risks softening agent rigor. Do not build implicit user-modeling from session history. Rely on explicit, human-reviewed context (learnings, conventions, hooks) for behavioral shaping.
- Chain-of-thought prompting improves agent reasoning accuracy (17.7% to 78.7%). Added "Reason Before Acting" to AGENT_PLAYBOOK.md and reasoning nudges to 7 skills.
- Say "project conventions" not "idiomatic X" to ensure Claude looks at project files first rather than triggering training priors (stdlib conventions).
- Autonomous "YOLO mode" is effective for feature velocity but accumulates technical debt (magic strings, monolithic tests, hardcoded paths). Schedule periodic consolidation sessions.
- Trust the binary output over source code analysis. A single ambiguous CLI output is not proof of absence — re-run the exact command before claiming something is missing.

---

## [2026-02-26-100009] Hook compliance and output routing (consolidated)

**Consolidated from**: 3 entries (2026-02-22 to 2026-02-25)

- Plain-text hook output is silently ignored by the agent. Claude Code parses hook stdout starting with `{` as JSON directives; plain text is disposable. All hooks should return JSON via `printHookContext()`.
- Hook compliance degrades on narrow mid-session tasks (~15-25% partial skip rate). Root cause: CLAUDE.md's "may or may not be relevant" system reminder competes with hook authority. Fix: CLAUDE.md explicitly elevates hook authority. The mandatory checkpoint relay block is the compliance canary.
- No reliable agent-side before-session-end event exists. SessionEnd fires after the agent is gone. Mid-session nudges and explicit /ctx-wrap-up are the only reliable persistence mechanisms.

---

## [2026-02-26-100010] ctx add and decision recording (consolidated)

**Consolidated from**: 4 entries (2026-01-27 to 2026-02-14)

- `ctx add learning` requires `--context`, `--lesson`, `--application` flags. `ctx add decision` requires `--context`, `--rationale`, `--consequences`. A bare string only sets the title and the command will fail without required flags.
- Structured entries with Context/Lesson/Application are more useful than one-liners. Agents are guided via AGENT_PLAYBOOK.md.
- Always complete decision record sections — placeholder text like "[Add context here]" is a code smell. Decisions without rationale lose their value over time.
- Slash commands using `!` bash syntax require matching permissions in settings.local.json. When adding new /ctx-* commands, ensure ctx init pre-seeds the required `Bash(ctx <subcommand>:*)` permissions.

---

## [2026-02-24-032945] CLI tools don't benefit from in-memory caching of context files

**Context**: Discussed whether ctx should read and cache LEARNINGS.md, DECISIONS.md etc. in memory

**Lesson**: ctx is a short-lived CLI process, not a daemon. Context files are tiny (few KB), sub-millisecond to read. Cache invalidation complexity exceeds the read cost. Caching only makes sense if ctx becomes a long-lived process (MCP server, watch daemon).

**Application**: Don't add caching layers to ctx's file reads. If an MCP server mode is ever added, revisit then.

---

## [2026-02-22-120000] Hook behavior and patterns (consolidated)

**Consolidated from**: 8 entries (2026-01-25 to 2026-02-17)

- Hook scripts receive JSON via stdin (not env vars); parse with `HOOK_INPUT=$(cat)` then jq
- Hook key names are case-sensitive: `PreToolUse` and `SessionEnd` (not `PreToolUseHooks`)
- Use `$CLAUDE_PROJECT_DIR` in hook paths, never hardcode absolute paths
- Hook regex can overfit: `ctx` as binary vs directory name differ; anchor patterns to command-start positions with `(^|;|&&|\|\|)\s*`
- grep patterns match inside quoted arguments — test with `ctx add learning "...blocked words..."` to verify no false positives
- Hook scripts can silently lose execute permission; verify with `ls -la .claude/hooks/*.sh` after edits
- Two-tier output is sufficient: unprefixed (agent context, may or may not relay) and `IMPORTANT: Relay VERBATIM` (guaranteed relay); don't add new severity prefixes
- Repeated injection causes agent repetition fatigue; use `--session $PPID --cooldown 10m` and pair with a readback instruction

---

## [2026-02-22-120001] UserPromptSubmit hook output channels (consolidated)

**Consolidated from**: 2 entries (2026-02-12)

- UserPromptSubmit hook stdout is prepended as AI context (not shown to user); stderr with exit 0 is swallowed entirely
- User-visible output requires `{"systemMessage": "..."}` JSON on stdout (warning banner) or exit 2 (blocks prompt)
- There is no non-blocking user-visible output channel for this hook type
- Design hooks for their actual audience: AI-facing = plain stdout, user-facing = systemMessage JSON

---

## [2026-02-22-120002] Linting and static analysis (consolidated)

**Consolidated from**: 7 entries (2026-01-25 to 2026-02-20)

- Full pre-commit gate: (1) `CGO_ENABLED=0 go build ./cmd/ctx`, (2) `golangci-lint run`, (3) `CGO_ENABLED=0 go test` — all three, every time
- Own the codebase: fix pre-existing lint issues even if you didn't introduce them
- gosec G301/G306: use 0o750 for dirs, 0o600 for files everywhere including tests
- gosec G304 (file inclusion): safe to suppress with `//nolint:gosec` in test files using `t.TempDir()` paths
- golangci-lint errcheck: use `cmd.Printf`/`cmd.Println` in Cobra commands instead of `fmt.Fprintf`
- `defer os.Chdir(x)` fails errcheck; use `defer func() { _ = os.Chdir(x) }()`
- golangci-lint Go version mismatch in CI: use `install-mode: goinstall` to build linter from source

---

## [2026-02-22-120006] Permission and settings drift (consolidated)

**Consolidated from**: 4 entries (2026-02-15)

- Permission drift is distinct from code drift — settings.local.json is gitignored, no review catches stale entries
- `Skill()` permissions don't support name prefix globs — list each skill individually
- Wildcard trusted binaries (`Bash(ctx:*)`, `Bash(make:*)`), but keep git commands granular (never `Bash(git:*)`)
- settings.local.json accumulates session debris; run periodic hygiene via `/sanitize-permissions` and `/ctx-drift`

---

## [2026-02-22-120008] Gitignore and filesystem hygiene (consolidated)

**Consolidated from**: 3 entries (2026-02-11 to 2026-02-15)

- Gitignored directories are invisible to `git status`; stale artifacts persist indefinitely — periodically `ls` gitignored working directories
- Add editor artifacts (*.swp, *.swo, *~) to .gitignore alongside IDE directories from day one
- Gitignore entries for sensitive paths are security controls, not documentation — never remove during cleanup sweeps

---

## [2026-02-19-215200] Feature can be code-complete but invisible to users

**Context**: ctx pad merge was fully implemented with 19 passing tests and binary support, but had zero coverage in user-facing docs (scratchpad.md, cli-reference.md, scratchpad-sync recipe). Only discoverable via --help.

**Lesson**: Implementation completeness \!= user-facing completeness. A feature without docs is invisible to users who don't explore CLI help.

**Application**: After implementing a new CLI subcommand, always check: feature page, cli-reference.md, relevant recipes, and zensical.toml nav (if new page).

---

## [2026-01-28-051426] IDE is already the UI

**Context**: Considering whether to build custom UI for .context/ files

**Lesson**: Discovery, search, and editing of .context/ markdown files works
better in VS Code/IDE than any custom UI we'd build. Full-text search,
git integration, extensions - all free.

**Application**: Don't reinvent the editor. Let users use their preferred IDE.

---

*Module-specific, niche, and historical learnings:
[learnings-reference.md](learnings-reference.md)*
