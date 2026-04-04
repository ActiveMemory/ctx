# Archived Learnings (consolidated 2026-04-03)

Originals replaced by consolidated entries in LEARNINGS.md.

## Group: Subagent scope creep and cleanup

## [2026-03-23-165610] Subagents rename functions and restructure code beyond their scope

**Context**: Agents tasked with fixing em-dashes in comments also renamed exported functions, changed import aliases, and modified function signatures

**Lesson**: Always diff-audit agent output for structural changes before accepting edits, even when the task is narrowly scoped

**Application**: After any agent batch edit, run git diff --stat and scan for non-comment changes before staging

---

## [2026-03-14-180902] Subagents rename packages and modify unrelated files without being asked

**Context**: ET agent renamed internal/eventlog/ to internal/log/ and modified 5+ caller files outside the internal/err/ scope

**Lesson**: Always diff agent output against HEAD to catch scope creep before building; revert unrelated changes immediately

**Application**: After any agent-driven refactor, run git diff --name-only HEAD and revert anything outside the intended scope before testing

---

## [2026-03-14-110750] Subagents reorganize file structure without being asked

**Context**: Asked subagent to replace os.ReadFile callsites with validation wrappers. It also moved functions to new files, renamed them (ReadUserFile to SafeReadUserFile), and created a new internal/io package.

**Lesson**: Subagents optimize for clean results and will restructure beyond stated scope — moved, renamed, and split files without being asked.

**Application**: After subagent-driven refactors, always verify file organization matches intent. Audit for moved, renamed, and split files, not just the requested callsite changes.

---

## [2026-03-06-200319] Spawned agents reliably create new files but consistently fail to delete old ones — always audit for stale files, duplicate function definitions, and orphaned imports after agent-driven refactoring

**Context**: Multiple agent batches across cmd/ restructuring, color removal, and flag externalization left stale files, duplicate run.go, and unupdated parent imports

**Lesson**: Agent cleanup is a known gap — budget 5-10 minutes for post-agent audit per batch

**Application**: After every agent batch: grep for stale package declarations, check parent imports point to cmd/root not cmd/, verify old files are deleted

---


## Group: Bulk rename and replace_all hazards

## [2026-03-20-232518] replace_all on short tokens like function names will also replace their definitions

**Context**: Using replace_all to rename HumanAgo to format.DurationAgo also changed the func declaration line to func format.DurationAgo which is invalid Go

**Lesson**: replace_all matches everywhere in the file including function definitions, not just call sites

**Application**: When renaming functions, delete the old definition separately rather than using replace_all. Or use a more specific match pattern that excludes func declarations.

---

## [2026-03-15-230643] replace_all on short tokens like core. mangles aliased imports

**Context**: Used Edit tool replace_all to change core. to tidy. across handle_tool.go

**Lesson**: Short replacement tokens match inside longer identifiers — remindcore. became remindtidy. silently

**Application**: When doing bulk renames, prefer targeted replacements or grep for collateral damage immediately after. Avoid replace_all on tokens shorter than the full identifier

---

## [2026-03-18-193602] Bulk sed on imports displaces aliased imports

**Context**: Used sed to add desc import to 278 files. Files with aliased imports like ctxCfg config/ctx got the alias stolen by the new import line inserted before it.

**Lesson**: sed insert-before-first-match does not understand Go import aliases. The alias attaches to whatever import line sed inserts, not the original target.

**Application**: When bulk-adding imports, check for aliased imports first and handle them separately. Or use goimports if available.

---


## Group: Import cycles and package splits

## [2026-03-22-220846] Types in god-object files create circular dependencies

**Context**: hook/types.go had 15+ types from 8 domains; session importing hook for SessionTokenInfo created a cycle

**Lesson**: Moving types to their owning domain package breaks import cycles

**Application**: When a type is only used by one domain package, move it there. Check with grep before assuming a type is cross-cutting.

---

## [2026-03-18-193616] Tests in package X cannot import X/sub packages that import X back

**Context**: embed_test.go in package assets kept importing read/ sub-packages that import assets.FS, creating cycles. Recurred 4 times this session.

**Lesson**: Tests that exercise domain accessor packages must live in those packages, not in the parent. The parent test file can only use the shared resource (FS) directly.

**Application**: When splitting functions from a parent package into sub-packages, move the corresponding tests too. Do not leave them in the parent.

---

## [2026-03-13-223108] Variable shadowing causes cascading test failures after package splits

**Context**: Large refactoring moved constants from monolithic config package to sub-packages (dir, entry, file). Test files had local variables named dir, entry, file that shadowed the new package imports.

**Lesson**: When splitting a package, audit test files for local variable names that collide with new package names. dir, file, entry are especially common Go variable names.

**Application**: Before committing a package split, run go test ./... and check for undefined errors caused by variable shadowing.

---

## [2026-03-13-151951] Linter reverts import-only edits when references still use old package

**Context**: Moving tpl_entry.go from config/entry to assets — linter kept reverting the import change

**Lesson**: When moving constants between packages, change imports and all references in a single atomic write (use Write not incremental Edit), so the linter never sees an inconsistent state

**Application**: For future package migrations, use full file rewrites when a linter is active

---

## [2026-03-06-200237] Import cycle avoidance: when package A imports package B for logic, B must own shared types — A aliases them. entry imports add/core for insert logic, so add/core owns EntryParams and entry aliases it as entry.Params

**Context**: Extracting entry.Params as a standalone struct in internal/entry created a cycle because entry/write.go imports add/core for AppendEntry

**Lesson**: The package that provides implementation logic must own the types; the facade package aliases them

**Application**: When extracting shared types from implementation packages, check the import direction first — the type lives where the logic lives

---


## Group: Lint suppression and gosec patterns

## [2026-03-19-194942] Rename constants to avoid gosec G101 false positives

**Context**: Constants named ColTokens, DriftPassed, StatusTokensFormat triggered gosec G101 credential detection. Suppressing with nolint or broadening .golangci.yml exclusion paths is fragile — paths change when files split.

**Lesson**: Rename the constant to avoid the trigger word instead of suppressing the lint. Tokens→Usage, Passed→OK convey the same semantics without false positives. This is cleaner than nolint annotations or path-based exclusions that break on file reorganization.

**Application**: When gosec flags a constant name, ask what the value semantically represents and rename to that. Do not add nolint, nosec, or path exclusions.

---

## [2026-03-06-050126] nolint:goconst for trivial values normalizes magic strings

**Context**: Found 5 callsites suppressing goconst with nolint:goconst trivial user input check for y/yes comparisons

**Lesson**: Suppressing a linter with a trivial justification sets precedent for other agents to do the same. The fix (two constants) costs less than the accumulated tech debt.

**Application**: Use config.ConfirmShort/config.ConfirmLong instead of suppressing goconst. Prefer constants over nolint directives.

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


## Group: Skill lifecycle and promotion

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

## [2026-03-01-125807] Elevating private skills requires synchronized updates across 6 layers

**Context**: Promoted 6 _ctx-* skills to bundled ctx-* plugin skills

**Lesson**: Moving a skill from .claude/skills/ to internal/assets/claude/skills/ touches: (1) SKILL.md frontmatter name field, (2) internal cross-references between skills (slash command paths), (3) external cross-references in other skills and docs, (4) embed_test.go expected skill list, (5) recipe and reference docs that mention the old name, (6) plugin cache rebuild (`hack/plugin-reload.sh`) + session restart — Claude Code snapshots skills from `~/.claude/plugins/cache/` at startup, so new skills are invisible until the cache is refreshed. Also clean stale underscore-prefixed `Skill(_ctx-*)` entries from `.claude/settings.local.json`.

**Application**: When promoting future skills, use grep -r /_ctx-{name} across the whole tree before declaring done. After code changes, run plugin-reload.sh and restart the session to verify the skill appears in autocomplete.

---

## [2026-03-01-144544] Skill enhancement is a documentation-heavy operation across 10+ files

**Context**: Enhancing /ctx-journal-enrich-all to handle export-if-needed touched the skill, hook messages, fallback strings, 5 doc files, 2 Makefiles, and TASKS.md

**Lesson**: Skill behavior changes ripple through hook messages, fallback strings in Go code, doc descriptions, and Makefile hints — all must stay synchronized

**Application**: When modifying a skill's scope, grep for its name across the entire repo and update every description, not just the skill file itself

---


## Group: Cross-cutting change ripple

## [2026-03-01-194147] Key path changes ripple across 15+ doc files and 2 skills

**Context**: Updating docs for the .context/.ctx.key → ~/.local/ctx/keys/ → ~/.ctx/.ctx.key migrations

**Lesson**: Key path changes have a long documentation tail — recipes, references, getting-started, operations, CLI docs, and skills all carry path references. The worktree behavior flip (limitation to works automatically) was the highest-value change per line edited. Simplifying from per-project slugs to a single global key eliminated more code and docs than the original migration added.

**Application**: When moving a file path that appears in user-facing docs, grep broadly (not just code) and budget for 15+ file touches

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

## [2026-02-19-215200] Feature can be code-complete but invisible to users

**Context**: ctx pad merge was fully implemented with 19 passing tests and binary support, but had zero coverage in user-facing docs (scratchpad.md, cli-reference.md, scratchpad-sync recipe). Only discoverable via --help.

**Lesson**: Implementation completeness \!= user-facing completeness. A feature without docs is invisible to users who don't explore CLI help.

**Application**: After implementing a new CLI subcommand, always check: feature page, cli-reference.md, relevant recipes, and zensical.toml nav (if new page).

---


## Group: Dead code detection

## [2026-03-30-003720] internal/cli/recall/ was dead code — never registered in bootstrap

**Context**: The entire recall CLI package existed with tests but was never wired into the command tree. Journal consumed it but nobody deleted the ghost

**Lesson**: Dead package detection requires checking bootstrap registration, not just build success. A package can build and test green while being completely unreachable

**Application**: Add a compliance test that verifies all cli/ packages are registered in bootstrap

---

## [2026-03-25-173339] Dead files accumulate when nothing consumes them

**Context**: IMPLEMENTATION_PLAN.md and PROMPT.md were created by ctx init but no agent, hook, or skill ever read them

**Lesson**: Before adding a file to init scaffolding, verify there is at least one consumer. Periodically audit what init creates vs what the system reads.

**Application**: The prompt deprecation spec documents the reasoning as a papertrail for future removals.

---

## [2026-03-15-101346] Delete legacy code instead of maintaining it — MigrateKeyFile had 5 callers and test coverage but zero users

**Context**: Started by adding constants for legacy key names, then realized nobody uses legacy keys

**Lesson**: When touching legacy compat code, first ask whether the legacy path has real users. If not, delete it entirely rather than improving it

**Application**: Apply to any backward-compat shim: check actual usage before investing in maintenance

---


