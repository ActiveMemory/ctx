#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# lint-shellcheck.sh — run shellcheck against embedded shell scripts.
#
# Scope: scripts that ship inside the binary as embedded assets and
# run on user machines (`internal/assets/hooks/trace/*.sh` and
# `internal/assets/integrations/copilot-cli/scripts/*.sh`). These
# are the highest-stakes scripts — a bug there hits every user, not
# just contributors. `hack/` scripts are out of scope for now (they
# run only on developer machines).
#
# Severity: fails on `warning` and above (matches the project's
# zero-issues posture for other linters). Use --info to see notes
# too.
#
# Exit code:
#   0 — no findings
#   1 — findings or shellcheck not installed

set -euo pipefail

if ! command -v shellcheck >/dev/null 2>&1; then
  echo "shellcheck not installed. Install via:" >&2
  echo "  macOS:        brew install shellcheck" >&2
  echo "  Debian/Ubuntu: apt-get install shellcheck" >&2
  exit 1
fi

SEVERITY="${SEVERITY:-warning}"

# Targets: embedded scripts only. Sorted so the output is stable
# across local runs and CI.
TARGETS=()
while IFS= read -r -d '' f; do
  TARGETS+=("$f")
done < <(
  find \
    internal/assets/hooks/trace \
    internal/assets/integrations/copilot-cli/scripts \
    -type f -name "*.sh" -print0 | sort -z
)

if [[ ${#TARGETS[@]} -eq 0 ]]; then
  echo "No embedded shell scripts found." >&2
  exit 0
fi

echo "Running shellcheck (severity=$SEVERITY) on ${#TARGETS[@]} script(s)..."
shellcheck --severity="$SEVERITY" "${TARGETS[@]}"
echo "shellcheck: no findings at severity >= $SEVERITY."
