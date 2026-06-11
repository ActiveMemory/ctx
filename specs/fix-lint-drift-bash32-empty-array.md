# Fix lint-drift Empty-Array Expansion on Bash 3.2

`hack/lint-drift.sh` aborted on stock macOS before running a
single check, which silently broke `make audit` (and therefore
the contributing guide's mandatory pre-PR gate) for every Mac
contributor.

## Problem

The `drift_grep` helper builds its `--exclude` flags in an
array and expands it unconditionally:

```bash
local exclude_args=()
for ex in "$@"; do
  exclude_args+=(--exclude="$ex")
done
grep -rn --include='*.go' --exclude='*_test.go' "${exclude_args[@]}" \
  -E "$pattern" internal/ 2>/dev/null || true
```

The script runs under `set -euo pipefail`. On bash 4.4+ an
empty `"${arr[@]}"` expands to zero words; on bash 3.2 — the
newest bash Apple ships, frozen at the GPLv2 boundary — the
same expansion is an **unbound variable** error under `set -u`:

```
./hack/lint-drift.sh: line 39: exclude_args[@]: unbound variable
```

Two `drift_grep` call sites pass no exclude globs (checks 2 and
3). Check 8 (`strings.Join`) scans with a *direct* `grep`, not
`drift_grep`, so it never touches the array and is unaffected.
The script dies on the first no-exclude `drift_grep` call, so
`make lint-style` → `make audit` fail before any drift check
executes.

## Solution

Guard the expansion with the parameter-expansion alternate
form, the canonical bash-3.2-safe idiom:

```bash
${exclude_args[@]+"${exclude_args[@]}"}
```

When the array is empty the outer expansion produces nothing;
when populated it reproduces the original quoted expansion
verbatim. Behavior on bash 4+ is unchanged. A comment at the
call site documents why the guard exists so a future cleanup
doesn't "simplify" it back.

Verified: `bash hack/lint-drift.sh` reproduced the
`line 39: exclude_args[@]: unbound variable` abort on stock
macOS bash 3.2.57 before the guard, and exits 0
(`lint-drift: clean`) after — so `make lint-style` → `make audit`
again run end-to-end on a stock Mac.

Note on automated coverage: no CI or `make` gate actually *lints*
this file. `hack/lint-shellcheck.sh` explicitly excludes `hack/`
scripts ("they run only on developer machines"), and CI runs
`lint-shellcheck` but not `lint-drift`/`audit`, so a Linux/bash-5
runner could not reproduce a 3.2-only bug anyway. The guarded
`${arr[@]+"${arr[@]}"}` is the canonical shellcheck-clean idiom;
a contributor with `shellcheck` installed can confirm with
`shellcheck hack/lint-drift.sh`. The inline call-site comment
plus the `LEARNINGS.md` entry are the durable guard against a
future "simplification" reintroducing the bug.

## Out of Scope

- Hardening the other `hack/` scripts' array expansions. A
  `grep -rn '\[@\]' hack/` sweep plus empirical bash 3.2
  checks show the remaining sites are all safe today, though
  for different reasons:
  - `lint-shellcheck.sh` and `lint-powershell.sh` both exit on
    a `${#TARGETS[@]} -eq 0` count guard (count expansion does
    not trip `set -u` on 3.2) before reaching their
    `"${TARGETS[@]}"` element expansion.
  - `build-all.sh` expands a `TARGETS` array populated from a
    hardcoded non-empty literal.
  - `detect-ai-typography.sh` expands `"${EXTS[@]}"`, parsed
    from `EXT`. `EXT` is *defaulted* (`EXT="${EXT:-md}"`), not
    a literal, so even `--ext ""` collapses to `md` and the
    array is never empty — safe, but via the default rather
    than a count guard.
  These are latent-only hazards (a future edit that lets an
  array reach the expansion empty), not failures today, and
  belong to a separate sweep if ever.
- Requiring bash 4+ (e.g. a version check or `#!/usr/bin/env
  bash4`). Contributors should not need a homebrew bash to run
  the project's own audit gate.
