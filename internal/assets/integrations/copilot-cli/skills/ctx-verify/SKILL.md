---
name: ctx-verify
description: "Verify before claiming completion. Use before saying work is done, tests pass, or builds succeed."
tools: [bash, read, glob, grep]
---

Run the relevant verification command before claiming a result.

## When to Use

- Before saying "tests pass", "build succeeds", or "bug fixed"
- Before reporting completion of any task with a testable outcome
- When the user asks "does it work?" or "is it done?"

## When NOT to Use

- For documentation-only changes with no testable outcome
- When the user explicitly says "skip verification"
- For exploratory work with no pass/fail criterion

## Workflow

1. **Identify** what command proves the claim
2. **Think through** what passing looks like (and false positives)
3. **Run** the command (fresh, not a previous run)
4. **Read** full output; check exit code, count failures
5. **Report** actual results with evidence

## Claim-to-Evidence Map

| Claim             | Required Evidence                          |
|-------------------|--------------------------------------------|
| Tests pass        | Test command output showing 0 failures     |
| Linter clean      | `golangci-lint run` showing 0 errors       |
| Build succeeds    | `go build` exit 0                          |
| Bug fixed         | Original symptom no longer reproduces      |
| All checks pass   | `make audit` showing all steps pass        |

## Self-Audit Questions

Before presenting any artifact as complete:
- What assumptions did I make?
- What did I NOT check?
- Where am I least confident?
- What would a reviewer question first?

## Quality Checklist

- [ ] Verification command was run fresh (not reused)
- [ ] Exit code was checked
- [ ] Claim matches evidence (build ≠ tests)
- [ ] If multiple claims, each has its own evidence
