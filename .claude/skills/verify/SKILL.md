---
name: verify
description: "Verify before claiming completion. Use before saying work is done, tests pass, or builds succeed."
---

Run the relevant verification command before claiming a result.

## Workflow

1. **Identify** what command proves the claim
2. **Run** the command (fresh, not a previous run)
3. **Read** the full output — check exit code, count failures
4. **Report** actual results with evidence

## Claim-to-Evidence Map

| Claim             | Required Evidence                          |
|-------------------|--------------------------------------------|
| Tests pass        | Test command output showing 0 failures     |
| Linter clean      | `golangci-lint run` output showing 0 errors |
| Build succeeds    | `go build` exit 0 (linter passing is not enough) |
| Bug fixed         | Original symptom no longer reproduces      |
| Regression tested | Red-green cycle: test fails without fix, passes with it |
| All checks pass   | `make audit` output showing all steps pass |

## Examples

```
Good:  ran `make audit` → "All checks pass (format, vet, lint, test)"
Good:  ran `go test ./...` → "34/34 tests pass"
Avoid: "Should pass now" / "Looks correct" (without running anything)
```

## Relationship to /qa

`/qa` tells you *what to run*. `/verify` reminds you to *actually run it*
before claiming the result.
