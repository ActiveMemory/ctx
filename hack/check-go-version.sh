#!/usr/bin/env bash
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0


# check-go-version.sh — one Go version, many pins, zero drift.
#
# The `go` directive in the root go.mod is the single source of truth.
# Every other place the toolchain version is written down must agree
# with it, or a hurried bump ships a stale pin:
#
#   exact (major.minor.patch must match go.mod):
#     go.work
#     tools/ctxctl/go.mod
#   floor (major.minor must match go.mod):
#     .github/workflows/*.yml     every `go-version:` pin
#     hack/tool-versions.txt      the `go bin required <min>` row
#     .context/steering/tech.md   the "**Go X.Y+**" prose (tool-native
#                                 copies are covered by `make check-steering`)
#
# Portable: bash 3.2 + BSD awk/grep (see specs/hack-script-portability.md).
#
# Exit code: number of mismatches (0 = clean).

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

issues=0

fail() {
  echo "FAIL: $1"
  issues=$((issues + 1))
}

# First `go <version>` directive in a go.mod / go.work file.
go_directive() {
  awk '$1 == "go" { print $2; exit }' "$1"
}

FULL="$(go_directive go.mod)"
if [ -z "$FULL" ]; then
  echo "FAIL: no 'go' directive in go.mod" >&2
  exit 1
fi
MINOR="$(printf '%s\n' "$FULL" | awk -F. '{ print $1 "." $2 }')"

# --- exact pins -----------------------------------------------------------

for f in go.work tools/ctxctl/go.mod; do
  got="$(go_directive "$f")"
  if [ "$got" != "$FULL" ]; then
    fail "$f: go $got != go.mod $FULL"
  fi
done

# --- floor pins -----------------------------------------------------------

# Workflows: every go-version pin, quoted or bare.
found=0
for wf in .github/workflows/*.yml; do
  while IFS= read -r line; do
    found=$((found + 1))
    got="$(printf '%s\n' "$line" | sed -e "s/.*go-version:[[:space:]]*//" -e "s/['\"]//g" -e 's/[[:space:]]*$//')"
    if [ "$got" != "$MINOR" ]; then
      fail "$wf: go-version '$got' != go.mod $MINOR"
    fi
  done < <(grep -E '^[[:space:]]*go-version:' "$wf" || true)
done
if [ "$found" -eq 0 ]; then
  fail ".github/workflows: no go-version pins found (grep pattern drifted?)"
fi

# Tool manifest: the go row's minimum.
got="$(awk '$1 == "go" && $2 == "bin" { print $4; exit }' hack/tool-versions.txt)"
if [ "$got" != "$MINOR" ]; then
  fail "hack/tool-versions.txt: go minimum '$got' != go.mod $MINOR"
fi

# Steering source prose: "**Go X.Y+**".
got="$(grep -o -E 'Go [0-9]+\.[0-9]+\+' .context/steering/tech.md | head -1 | sed -e 's/^Go //' -e 's/+$//')"
if [ "$got" != "$MINOR" ]; then
  fail ".context/steering/tech.md: 'Go ${got:-?}+' != go.mod $MINOR"
fi

# --- verdict --------------------------------------------------------------

if [ "$issues" -gt 0 ]; then
  echo "Go version drift: $issues mismatch(es) — go.mod says $FULL; align the sites above."
  exit "$issues"
fi
echo "Go version sync OK ($FULL / floor $MINOR)."
