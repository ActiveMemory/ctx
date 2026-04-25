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
|----|--------|
| 2026-04-25 | Confident code comments can pull an LLM away from first-principles knowledge |
| 2026-04-25 | filepath.Join('', rel) returns rel as CWD-relative, not error |
| 2026-04-25 | Parallel go test ./... packages can race on ~/.claude/settings.json |
<!-- INDEX:END -->

<!-- Add gotchas, tips, and lessons learned here -->
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
