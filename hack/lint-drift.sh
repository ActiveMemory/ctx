#!/usr/bin/env bash
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# lint-drift.sh — catch code-level drift that static analyzers miss.
#
# Checks:
#   1. Literal "\n" in non-test .go files (should use config.NewlineLF)
#   2. Printf/PrintErrf with trailing \n (should use Println)
#   3. Magic directory strings that have config.Dir* constants
#   4. Literal ".md" (should use config.ExtMarkdown)
#
# Exit code: number of issues found (0 = clean).

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

issues=0

# Helper: grep non-test .go files, excluding specific paths.
# Args: pattern [exclude_glob...]
drift_grep() {
  local pattern="$1"; shift
  local exclude_args=()
  for ex in "$@"; do
    exclude_args+=(--exclude="$ex")
  done
  grep -rn --include='*.go' --exclude='*_test.go' "${exclude_args[@]}" \
    -E "$pattern" internal/ 2>/dev/null || true
}

# Count lines from drift_grep output
drift_count() {
  if [ -z "$1" ]; then
    echo 0
  else
    echo "$1" | wc -l | tr -d ' '
  fi
}

# ── 1. Literal "\n" ─────────────────────────────────────────────────
# Match "\n" as a Go string (not inside comments or imports).
# Skip config/token.go where the constant is defined.
hits=$(drift_grep '"\\n"' 'token.go')
count=$(drift_count "$hits")
if [ "$count" -gt 0 ]; then
  echo "==> Literal \"\\n\" found ($count occurrences, use config.NewlineLF):"
  echo "$hits"
  echo ""
  issues=$((issues + count))
fi

# ── 2. cmd.Printf / cmd.PrintErrf ───────────────────────────────────
# These almost always end with \n; prefer Println(fmt.Sprintf(...)).
hits=$(drift_grep 'cmd\.(Printf|PrintErrf)\(')
count=$(drift_count "$hits")
if [ "$count" -gt 0 ]; then
  echo "==> cmd.Printf/PrintErrf calls ($count occurrences, prefer Println):"
  echo "$hits"
  echo ""
  issues=$((issues + count))
fi

# ── 3. Magic directory strings in filepath.Join ─────────────────────
# These directories have constants in config/dir.go.
for dir in '"sessions"' '"archive"' '"tools"'; do
  hits=$(drift_grep "filepath\.Join\(.*${dir}")
  count=$(drift_count "$hits")
  if [ "$count" -gt 0 ]; then
    echo "==> Magic directory ${dir} in filepath.Join ($count, use config.Dir*):"
    echo "$hits"
    echo ""
    issues=$((issues + count))
  fi
done

# ── 4. Literal ".md" ────────────────────────────────────────────────
# Skip config/file.go where ExtMarkdown is defined.
hits=$(drift_grep '"\.md"' 'file.go')
count=$(drift_count "$hits")
if [ "$count" -gt 0 ]; then
  echo "==> Literal \".md\" found ($count occurrences, use config.ExtMarkdown):"
  echo "$hits"
  echo ""
  issues=$((issues + count))
fi

# ── Summary ──────────────────────────────────────────────────────────
if [ "$issues" -eq 0 ]; then
  echo "lint-drift: clean"
else
  echo "lint-drift: $issues issues found"
fi

exit "$issues"
