#!/usr/bin/env bash
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0


# check-go-work-sum.sh — go.work.sum must already hold every checksum
# the workspace needs.
#
# dependabot bumps each module's go.mod/go.sum but never go.work.sum,
# so the workspace sum goes stale after every Go dependency PR, and the
# next contributor whose tooling runs `go mod download` (or
# `go list -m all`, `go mod verify`, gopls) finds it dirty. This gate
# runs that download and fails if the file had to change.
#
# Snapshot-then-regenerate, like check-steering: the comparison is
# against the working copy, not HEAD, so an uncommitted refresh passes.
# On failure the refreshed go.work.sum is left in place, ready to commit.
#
# Portable: bash 3.2 + BSD diff/sed (see specs/hack-script-portability.md).
#
# Exit code: 0 = complete, 1 = stale (go.work.sum now refreshed).

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

snapshot="$(mktemp)"
trap 'rm -f "$snapshot"' EXIT

cp go.work.sum "$snapshot"
go mod download

if ! cmp -s "$snapshot" go.work.sum; then
  added="$(diff "$snapshot" go.work.sum | grep -c '^>' || true)"
  echo "FAIL: go.work.sum was missing $added workspace checksum line(s);"
  echo "      it has been refreshed in place. Commit it."
  diff "$snapshot" go.work.sum | sed -n 's/^> /  + /p' | head -n 10 || true
  exit 1
fi
echo "go.work.sum is complete."
