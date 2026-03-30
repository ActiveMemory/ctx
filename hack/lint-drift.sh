#!/usr/bin/env bash
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0


# lint-drift.sh — catch code-level drift that static analyzers miss.
#
# Checks:
#   1. Literal "\n" in non-test .go files (should use config.NewlineLF)
#   2. Printf/PrintErrf with trailing \n (should use Println)
#   3. Magic directory strings that have config.Dir* constants
#   4. Literal ".md" (should use config.ExtMarkdown)
#   5. Go DescKey ↔ YAML key linkage (embed/{cmd,flag,text} vs commands/)
#      Only DescKey* constants are checked — Use* and ExampleKey* are excluded.
#   6. Inline error strings outside internal/err/ (should use err/ constructors)
#   7. err.go files outside internal/err/ (broken-window magnet)
#   8. strings.Join with inline separator (should use token.* constants)
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
hits=$(drift_grep '"\\n"' 'whitespace.go')
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
hits=$(drift_grep '"\.md"' 'ext.go')
count=$(drift_count "$hits")
if [ "$count" -gt 0 ]; then
  echo "==> Literal \".md\" found ($count occurrences, use config.ExtMarkdown):"
  echo "$hits"
  echo ""
  issues=$((issues + count))
fi

# ── 6. Inline error strings outside internal/err/ ─────────────────────
# fmt.Errorf("literal") and errors.New("literal") should live in
# internal/err/ with YAML-backed text keys, not inline in core/cmd/.
# Exclude internal/err/ (that's where constructors belong) and lines
# that already use desc.Text/lookup.TextDesc (already externalized).
hits=$(grep -rn --include='*.go' --exclude='*_test.go' \
  -E 'fmt\.Errorf\s*\(\s*["`]' internal/ 2>/dev/null \
  | grep -v 'internal/err/' \
  || true)
count=$(drift_count "$hits")
if [ "$count" -gt 0 ]; then
  echo "==> Inline fmt.Errorf strings outside internal/err/ ($count, use err/ constructors):"
  echo "$hits"
  echo ""
  issues=$((issues + count))
fi

hits=$(grep -rn --include='*.go' --exclude='*_test.go' \
  -E 'errors\.New\s*\(\s*"' internal/ 2>/dev/null \
  | grep -v 'internal/err/' \
  || true)
count=$(drift_count "$hits")
if [ "$count" -gt 0 ]; then
  echo "==> Inline errors.New strings outside internal/err/ ($count, use err/ constructors):"
  echo "$hits"
  echo ""
  issues=$((issues + count))
fi

# ── 7. err.go files outside internal/err/ ─────────────────────────────
# Convention: error constructors belong in internal/err/, never in
# per-package err.go files. An err.go outside err/ is a broken-window
# magnet that invites agents to add local error constructors.
hits=$(find internal/ -name 'err.go' -not -path 'internal/err/*' 2>/dev/null || true)
count=$(drift_count "$hits")
if [ "$count" -gt 0 ]; then
  echo "==> err.go files outside internal/err/ ($count, move to internal/err/):"
  echo "$hits" | sed 's/^/    /'
  echo ""
  issues=$((issues + count))
fi

# ── 8. strings.Join with inline separator ─────────────────────────────
# Separators like ", " should use token.CommaSpace and friends.
# Skip token/ where the constants are defined.
hits=$(grep -rn --include='*.go' --exclude='*_test.go' \
  -E 'strings\.Join\([^)]+,\s*"' internal/ 2>/dev/null \
  | grep -v 'internal/config/token/' \
  || true)
count=$(drift_count "$hits")
if [ "$count" -gt 0 ]; then
  echo "==> strings.Join with inline separator ($count, use token.* constants):"
  echo "$hits"
  echo ""
  issues=$((issues + count))
fi

# ── 5. Go DescKey ↔ YAML key linkage ──────────────────────────────────
# Every Go constant in internal/config/embed/{cmd,flag,text}/ must have
# a matching YAML key, and vice versa. Orphans in either direction mean
# broken lookups or dead entries.
check_linkage() {
  local label="$1" go_dir="$2" yaml_file="$3"
  local go_keys yaml_keys missing_yaml missing_go

  go_keys=$(grep -rh '= "' "$go_dir"/*.go 2>/dev/null \
    | grep -v '//' \
    | sed 's/.*= "//; s/"//' \
    | sort -u)

  yaml_keys=$(grep -E '^[a-z]' "$yaml_file" \
    | sed 's/:$//' \
    | sort -u)

  missing_yaml=$(comm -23 <(echo "$go_keys") <(echo "$yaml_keys"))
  missing_go=$(comm -13 <(echo "$go_keys") <(echo "$yaml_keys"))

  local count=0
  if [ -n "$missing_yaml" ]; then
    local n; n=$(echo "$missing_yaml" | wc -l | tr -d ' ')
    echo "==> $label: $n Go constant(s) missing from YAML:"
    echo "$missing_yaml" | sed 's/^/    /'
    echo ""
    count=$((count + n))
  fi
  if [ -n "$missing_go" ]; then
    local n; n=$(echo "$missing_go" | wc -l | tr -d ' ')
    echo "==> $label: $n YAML key(s) missing from Go:"
    echo "$missing_go" | sed 's/^/    /'
    echo ""
    count=$((count + n))
  fi
  issues=$((issues + count))
}

# cmd constants ↔ commands.yaml
# Only DescKey* constants are YAML lookup keys.
# Use* constants are cobra Use fields; ExampleKey* map to examples.yaml.
cmd_go_dir="internal/config/embed/cmd"
cmd_go_keys=$(grep -rh 'DescKey.*= "' "$cmd_go_dir"/*.go 2>/dev/null \
  | grep -v '//' \
  | sed 's/.*= "//; s/"//' \
  | sort -u)
cmd_yaml_keys=$(grep -E '^[a-z]' "internal/assets/commands/commands.yaml" \
  | sed 's/:$//' \
  | sort -u)
cmd_missing_yaml=$(comm -23 <(echo "$cmd_go_keys") <(echo "$cmd_yaml_keys"))
cmd_missing_go=$(comm -13 <(echo "$cmd_go_keys") <(echo "$cmd_yaml_keys"))
if [ -n "$cmd_missing_yaml" ]; then
  n=$(echo "$cmd_missing_yaml" | wc -l | tr -d ' ')
  echo "==> cmd↔commands.yaml: $n Go constant(s) missing from YAML:"
  echo "$cmd_missing_yaml" | sed 's/^/    /'
  echo ""
  issues=$((issues + n))
fi
if [ -n "$cmd_missing_go" ]; then
  n=$(echo "$cmd_missing_go" | wc -l | tr -d ' ')
  echo "==> cmd↔commands.yaml: $n YAML key(s) missing from Go:"
  echo "$cmd_missing_go" | sed 's/^/    /'
  echo ""
  issues=$((issues + n))
fi

# flag constants ↔ flags.yaml
check_linkage "flag↔flags.yaml" \
  "internal/config/embed/flag" \
  "internal/assets/commands/flags.yaml"

# text constants ↔ text/*.yaml (merge all text YAML files)
text_yaml_merged=$(mktemp)
cat internal/assets/commands/text/*.yaml > "$text_yaml_merged"
check_linkage "text↔text/*.yaml" \
  "internal/config/embed/text" \
  "$text_yaml_merged"
rm -f "$text_yaml_merged"

# Cross-namespace duplicates: cmd, flag, and text are separate YAML files
# with separate Go packages — the same key in each is by design (e.g. "agent"
# appears as a command description, a flag description, and a text string).
# No check needed.

# ── Summary ──────────────────────────────────────────────────────────
if [ "$issues" -eq 0 ]; then
  echo "lint-drift: clean"
else
  echo "lint-drift: $issues issues found"
fi

exit "$issues"
