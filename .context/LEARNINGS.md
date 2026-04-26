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

<!-- INDEX:START -->
| Date | Learning |
|------|--------|
| 2026-04-26 | make test exit code unreliable due to -cover covdata tooling issue |
| 2026-04-26 | Trailing word boundary in regex matches commit-tree as git commit |
| 2026-04-26 | ctx system help can list project-local hooks not in the Go binary |
| 2026-04-25 | Confident code comments can pull an LLM away from first-principles knowledge |
| 2026-04-25 | filepath.Join('', rel) returns rel as CWD-relative, not error |
| 2026-04-25 | Parallel go test ./... packages can race on ~/.claude/settings.json |
<!-- INDEX:END -->

<!-- Add gotchas, tips, and lessons learned here -->
## [2026-04-26-152850] make test exit code unreliable due to -cover covdata tooling issue

**Context**: make test exited 1 even with all 123 packages passing on this Go install; root cause is missing covdata tool when -cover is enabled

**Lesson**: Don't trust make test exit code alone when verifying changes. The -cover flag in the test target can fail with 'no such tool covdata' even when every package passes.

**Application**: When make test fails, fall back to 'go test ./...' (no -cover) and tally ^ok / ^FAIL counts to distinguish real failures from tooling issues.

---

## [2026-04-26-152842] Trailing word boundary in regex matches commit-tree as git commit

**Context**: First post-commit filter regex \bgit\s+commit\b in the OpenCode plugin would have triggered on git commit-tree because \b matches between t and -

**Lesson**: A trailing word boundary doesn't exclude hyphenated continuations — \b matches every word/non-word transition. Use (?!-) negative lookahead to specifically reject hyphen-suffixed siblings.

**Application**: For any porcelain with hyphenated cousins (commit-tree, commit-graph, for-each-ref), append (?!-) to the boundary.

---

## [2026-04-26-152836] ctx system help can list project-local hooks not in the Go binary

**Context**: PR #72 plugin called 'ctx system block-dangerous-commands'; user's installed ctx 0.7.2 listed it in help, but no directory exists under internal/cli/system/cmd/ — it's a Claude Code plugin-local hook surfaced via wrapper

**Lesson**: ctx system help output is a union of compiled Go subcommands and project-local Claude wrappers; non-Claude integrations only see the Go subset

**Application**: When porting plugin behavior to a new editor, only call subcommands that have a directory under internal/cli/system/cmd/. Don't trust ctx system help output as the canonical surface.

---

## [2026-04-25-014704] Confident code comments can pull an LLM away from first-principles knowledge

**Context**: cli_test.go had a comment claiming 'parent's t.Setenv doesn't propagate to exec'd children unless we build it into cmd.Env' which is wrong. I patched the helper's CTX_DIR dedup instead of questioning the helper itself, despite knowing t.Setenv semantics.

**Lesson**: A comment that explains why a stdlib mechanism 'doesn't work' is doing extra rhetorical work to talk a reader out of the obvious approach. That's exactly when to verify from first principles instead of trusting the surrounding-code frame.

**Application**: When an existing comment justifies a non-canonical approach contradicting stdlib knowledge: pause, verify against memory of the actual API before patching within the existing frame.

---

## [2026-04-25-014704] filepath.Join('', rel) returns rel as CWD-relative, not error

**Context**: Recurring orphan jsonl-path-<sessionID> appeared at project root. Older state.Dir() returned ('', nil) when CTX_DIR was undeclared, so filepath.Join('', 'jsonl-path-XXX') = 'jsonl-path-XXX', writing relative to CWD.

**Lesson**: Functions returning a path-string must never return ('', nil). Sentinel errors force callers to gate, closing the silent CWD-relative write.

**Application**: Audit any (string, error) path-returner that historically had a ('', nil) shortcut. Closed for state.Dir and rc.ContextDir; check remaining resolvers.

---

## [2026-04-25-014704] Parallel go test ./... packages can race on ~/.claude/settings.json

**Context**: make test runs packages in parallel processes. Fourteen test files invoked initialize.Cmd().Execute(), which read-modify-writes ~/.claude/settings.json without HOME isolation.

**Lesson**: Under load the races materialized as flaky 'FAIL coverage: [no statements]' in cli/watch/core. Run alone the package passed; under parallel make test it failed intermittently.

**Application**: testctx.Declare now sets HOME alongside CTX_DIR. Centralized fix; future tests automatically isolate user-home writes.
