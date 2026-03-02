---
name: ctx-verify
description: "Verify before claiming completion. Use before saying work is done, tests pass, or builds succeed."
---

Run the relevant verification command before claiming a result.

## When to Use

- Before saying "tests pass", "build succeeds", or "bug fixed"
- Before reporting completion of any task with a testable
  outcome
- When the user asks "does it work?" or "is it done?"
- After fixing a failing test or build error

## When NOT to Use

- For documentation-only changes with no testable outcome
- When the user explicitly says "trust me, skip verification"
- For exploratory work where there is no pass/fail criterion

## Usage Examples

```text
/ctx-verify
/ctx-verify (before claiming the refactor is done)
```

## Workflow

1. **Identify** what command proves the claim
2. **Think through** what a passing result looks like — and what
   a false positive would look like — before running
3. **Run** the command (fresh, not a previous run)
4. **Read** the full output; check exit code, count failures
5. **Report** actual results with evidence

Never reuse output from a previous run. Always run fresh.

## Claim-to-Evidence Map

| Claim             | Required Evidence                                       |
|-------------------|---------------------------------------------------------|
| Tests pass        | Test command output showing 0 failures                  |
| Linter clean      | `golangci-lint run` output showing 0 errors             |
| Build succeeds    | `go build` exit 0 (linter passing is not enough)        |
| Bug fixed         | Original symptom no longer reproduces                   |
| Regression tested | Red-green cycle: test fails without fix, passes with it |
| All checks pass   | `make audit` output showing all steps pass              |
| Files match       | `diff` showing no differences (e.g., template vs live)  |
| Design is sound   | Assumptions listed, failure modes identified, alternatives considered |
| Doc is accurate   | Claims traced to source code or config; no stale references |
| Skill works       | Trigger conditions tested, output matches spec, edge cases covered |
| Config is correct  | Values validated against schema or runtime; no stale references |

## Self-Audit Questions

Before presenting any artifact (code, design, doc, config) as
complete, run this checklist on your own output:

- What assumptions did I make?
- What did I NOT check?
- Where am I least confident?
- What would a reviewer question first?

If any answer reveals a gap, address it before reporting done.
This applies to all artifact types — not just code.

## Transform Vague Tasks into Verifiable Goals

Before starting, rewrite the task as a testable outcome:

| Task as given         | Verifiable goal                                     |
|-----------------------|-----------------------------------------------------|
| "Add validation"      | Write tests for invalid inputs, then make them pass |
| "Fix the bug"         | Write a test that reproduces it, then make it pass  |
| "Refactor X"          | Ensure tests pass before and after                  |
| "Improve performance" | Measure before, change, measure after, compare      |

For multi-step work, pair each step with its check:

```
1. [Step] -> verify: [check]
2. [Step] -> verify: [check]
```

Strong success criteria let you loop independently.
Weak criteria ("make it work") require constant clarification.

## Examples

### Good

- Ran `make audit`: "All checks pass (format, vet, lint, test)"
- Ran `go test ./...`: "34/34 tests pass"
- Ran `diff live.md template.md`: "no differences"
- Ran `go build -o /dev/null ./cmd/ctx`: "exit 0"

### Bad

- "Should pass now" (without running anything)
- "Looks correct" (visual inspection is not verification)
- "Tests passed earlier" (stale result; code changed since)
- "The build works" (did you actually run it?)

## Relationship to /ctx-qa

`/ctx-qa` tells you *what to run*. `/ctx-verify` reminds you to
*actually run it* before claiming the result.

## Quality Checklist

Before reporting a claim as verified:
- [ ] The verification command was run fresh (not reused)
- [ ] Exit code was checked (not just output scanned)
- [ ] The claim matches the evidence (build exit 0 does not
      prove tests pass)
- [ ] If multiple claims, each has its own evidence
- [ ] Self-audit questions answered (no unaddressed gaps)
- [ ] For non-code artifacts: relevant artifact verification
      criteria met
