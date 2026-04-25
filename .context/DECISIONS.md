# Decisions

<!-- INDEX:START -->
| Date | Decision |
|----|--------|
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
