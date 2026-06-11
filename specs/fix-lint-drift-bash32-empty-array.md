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

Several `drift_grep` call sites pass no exclude globs (checks
2, 3, and 8), so the script dies on its first such call and
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

Verified: `make audit` passes end-to-end on macOS
bash 3.2.57 with the guard in place.

## Out of Scope

- Hardening the other `hack/` scripts' array expansions. A
  `grep -rn '\[@\]' hack/` sweep plus empirical bash 3.2
  checks show the remaining sites are all safe today:
  `lint-shellcheck.sh` exits on a `${#TARGETS[@]}` count
  guard (count expansion does not trip `set -u` on 3.2)
  before its element expansion; `build-all.sh` and
  `detect-ai-typography.sh` expand arrays populated from
  hardcoded non-empty literals. Those are latent-only
  hazards (someone emptying a config array), not failures,
  and belong to a separate sweep if ever.
- Requiring bash 4+ (e.g. a version check or `#!/usr/bin/env
  bash4`). Contributors should not need a homebrew bash to run
  the project's own audit gate.
