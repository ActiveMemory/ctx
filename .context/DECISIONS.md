# Decisions

<!-- INDEX:START -->
| Date | Decision |
|------|--------|
| 2026-04-26 | OpenCode tool.execute.before omission is permanent; block-dangerous-commands will not become a ctx Go subcommand |
| 2026-04-26 | Editor-integration plugins must filter post-commit to actual git commit invocations |
| 2026-04-26 | OpenCode plugin ships without tool.execute.before hook |
| 2026-04-25 | Use t.Setenv for subprocess env in tests, not append(os.Environ(), ...) |
| 2026-04-25 | Tighten state.Dir / rc.ContextDir to (string, error) with sentinel errors |
<!-- INDEX:END -->

<!-- DECISION FORMATS

## Quick Format (Y-Statement)

For lightweight decisions, a single statement suffices:

> "In the context of [situation], facing [constraint], we decided for [choice]
> and against [alternatives], to achieve [benefit], accepting that [trade-off]."

## Full Format

For significant decisions:

## [YYYY-MM-DD] Decision Title

**Status**: Accepted | Superseded | Deprecated

**Context**: What situation prompted this decision? What constraints exist?

**Alternatives Considered**:
- Option A: [Pros] / [Cons]
- Option B: [Pros] / [Cons]

**Decision**: What was decided?

**Rationale**: Why this choice over the alternatives?

**Consequence**: What are the implications? (Include both positive and negative)

**Related**: See also [other decision] | Supersedes [old decision]

## When to Record a Decision

✓ Trade-offs between alternatives
✓ Non-obvious design choices
✓ Choices that affect architecture
✓ "Why" that needs preservation

✗ Minor implementation details
✗ Routine maintenance
✗ Configuration changes
✗ No real alternatives existed

-->
## [2026-04-26-231517] OpenCode tool.execute.before omission is permanent; block-dangerous-commands will not become a ctx Go subcommand

**Status**: Accepted

**Context**: The 2026-04-26-152858 decision shipped the OpenCode plugin without a tool.execute.before hook and noted "Re-add when block-dangerous-commands is promoted to the ctx Go binary." Revisited: that promotion is no longer planned. Keeping the open task on the books makes future sessions believe a re-add is pending.

**Decision**: We will not promote block-dangerous-commands to a ctx system Go subcommand. The OpenCode plugin's missing tool.execute.before hook is permanent, not deferred.

**Rationale**: The Cobra exit-1 / `{ blocked: true }` interaction makes any shim hostile to users without the Claude wrapper, and the safety-hook gap is acceptable given OpenCode's positioning. Recording this avoids the tax of a perpetually-pending follow-up that no one intends to land.

**Consequences**: TASKS.md item "Promote 'block-dangerous-commands' to a real ctx system Go subcommand…" marked `[-]` skipped. The 2026-04-26-152858 rationale's "Re-add when…" clause is void; the underlying ship-without-the-hook decision remains in force. Other (non-OpenCode) editor integrations that want a dangerous-command safety net will need a different mechanism.

**Related**: Amends [2026-04-26-152858] OpenCode plugin ships without tool.execute.before hook (rationale's deferred re-add is now closed).

---

## [2026-04-26-152905] Editor-integration plugins must filter post-commit to actual git commit invocations

**Status**: Accepted

**Context**: Original PR #72 OpenCode plugin ran 'ctx system post-commit' after every shell tool call, not only after real commits

**Decision**: Editor-integration plugins must filter post-commit to actual git commit invocations

**Rationale**: post-commit is meaningful only after a real commit lands; firing on every shell call is noise that trains users to ignore the resulting nudges

**Consequences**: Editor plugins always sniff the actual command string (regex on the extracted command) before triggering capture nudges that target specific commands. Same pattern applies to any future hook that targets a specific porcelain command.

---

## [2026-04-26-152858] OpenCode plugin ships without tool.execute.before hook

**Status**: Accepted

**Context**: The natural fit (block-dangerous-commands) doesn't exist as a ctx system Go subcommand; shimming to it would block every shell call on installs without the Claude wrapper because Cobra's unknown-command exit 1 is read as { blocked: true } by OpenCode

**Decision**: OpenCode plugin ships without tool.execute.before hook

**Rationale**: Better to ship a feature-narrower plugin than one that bricks the editor for users without the wrapper. Re-add when block-dangerous-commands is promoted to the ctx Go binary.

**Consequences**: OpenCode users get bootstrap, persistence, post-commit, and task-completion nudges but no dangerous-command safety net. specs/opencode-integration.md records the deliberate omission.

---

## [2026-04-25-014704] Use t.Setenv for subprocess env in tests, not append(os.Environ(), ...)

**Status**: Accepted

**Context**: TestBinaryIntegration spawns subprocesses; the prior helper did append(os.Environ(), CTX_DIR=...) to override the developer-shell value. Wrong abstraction.

**Decision**: Use t.Setenv for subprocess env in tests, not append(os.Environ(), ...)

**Rationale**: t.Setenv mutates the live process env, exec.Cmd with nil Env inherits it, and cleanup is automatic at test end. One line replaces the helper.

**Consequence**: Helper deleted, six call sites simplified, no env-dedup logic to maintain. Pattern reusable for other subprocess tests.

---

## [2026-04-25-014704] Tighten state.Dir / rc.ContextDir to (string, error) with sentinel errors

**Status**: Accepted

**Context**: Old single-return form returned ('', nil) when CTX_DIR was undeclared. Callers that filtered only on err != nil joined empty stateDir with relative names and wrote state files into CWD instead of .context/state/.

**Decision**: Tighten state.Dir / rc.ContextDir to (string, error) with sentinel errors

**Rationale**: Returning a sentinel ErrDirNotDeclared makes the empty-path case unrepresentable in a 'looks fine' branch. Forces every caller through the same explicit gate.

**Consequence**: All callers needed migration; tests had to declare CTX_DIR explicitly. In return, the filepath.Join('', rel) trap is closed by construction.
