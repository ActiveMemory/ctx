# Error-Handling Audit (Phase EH)

Systematically surface and resolve every silently discarded error
under `internal/`. "Fix only the risky ones" is not the bar — every
discard is audited and given a verdict; the risky ones are fixed and
the rest are justified in place.

## Problem

The codebase discards errors at ~184 non-test sites via `_ =`,
`_, _ =`, and `x, _ :=`. Some are legitimate (a discarded `ok` bool,
best-effort CLI output), but many hide real failures: a dropped
`Marshal`/`Parse` error writes bad data, a dropped write-handle
`Close` loses the final flush, a dropped `os.Remove` leaves stale
state. The same class of silent loss motivated
`specs/fix-learning-add-index-data-loss.md`.

## Approach

1. **Catalogue (EH.1).** Recursive walk of `internal/`, every discard
   site classified with a recommended action, written to
   `.context/audit/eh-silent-errors.md`. The catalogue is an
   inventory, not a verdict: categories assigned by pattern/name are
   provisional.

2. **Verify before fixing.** Every site is read in its enclosing
   context before any edit. Name-inference is untrustworthy — the
   first cut of the catalogue mislabelled two sites
   (`MergePublished` returns `(string, bool)`, not an error;
   `LoadState` returns a value, not a pointer, so no nil-deref). Per
   the Constitution's Context Integrity Invariants, the discarded
   value's actual type and the call's role decide the category.

3. **Resolve by category.**
   - **Data path** (a dropped error lets bad/empty/partial data get
     written, or an unreadable source gets silently overwritten):
     return the error. Fail loud.
   - **Best-effort** (telemetry, display hints, background loops with
     no return path, file close on read, `os.Remove`/`Rename`
     cleanup): `logWarn.Warn(cfgWarn.<Key>, err)` — the project sink
     for "not actionable by the caller but must not be swallowed"
     (`internal/log/warn`, prefixes `ctx: `, stderr in prod /
     `io.Discard` in tests).
   - **Write-handle close**: surface the close error (a failed final
     flush is data loss), via a named-return merge or `logWarn`.
   - **Category (d)** `fmt.Fprint(cmd.Out/Err …)`: accepted end-state
     per EH.5; CLI output is best-effort by construction.
   - **`ok`-discards / nil-safe / init-time programmer-error /
     false-positives**: annotate intent so a reviewer (and the EH.5
     grep) reads them as deliberate.

4. **Validate (EH.5).** `grep -rn '_ =' internal/` resolves to only
   accepted-and-annotated sites; `make lint && make test` green.

## Settled Decisions

1. `logWarn.Warn` is the canonical sink for best-effort discards;
   new per-call messages live as `cfgWarn` keys, not inline strings.
2. Data-path errors are returned (fail loud), never logged-and-
   continued.
3. The catalogue's pattern-assigned categories are provisional;
   the verdict is set per-site at fix time after reading the code.
4. Tests are out of scope for this pass.

## Out of Scope

- Phase ET (Error Package Taxonomy, `internal/err/`) — separate phase.
- Test-file discards.
