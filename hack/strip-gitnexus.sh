#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# strip-gitnexus.sh — remove the GitNexus auto-injected block from the
# tracked agent-instruction files (AGENTS.md, CLAUDE.md).
#
# `gitnexus analyze` injects a marker-bounded block
# (<!-- gitnexus:start --> … <!-- gitnexus:end -->) into these files.
# That block is generated output, not source: its canonical home is
# GITNEXUS.md. This script deletes the block (markers included) so the
# tracked files stay clean — AGENTS.md as the CLAUDE.md redirect stub,
# CLAUDE.md ending at its Companion Tools / GITNEXUS.md pointer.
#
# Belt-and-suspenders to `gitnexus analyze --skip-agents-md`; run it
# after any `gitnexus analyze` that ran without that flag.
#
#   ./hack/strip-gitnexus.sh
#
# Idempotent: a no-op when the markers are absent. Never touches
# GITNEXUS.md (the intended managed home for that content).
set -euo pipefail

START='<!-- gitnexus:start -->'
END='<!-- gitnexus:end -->'

# strip_one removes the marker block from a single file, then drops any
# blank lines left dangling at end-of-file so the file ends cleanly.
strip_one() {
  local file="$1"
  [ -f "$file" ] || return 0
  if ! grep -qF "$START" "$file"; then
    return 0 # nothing to strip
  fi

  local tmp
  tmp="$(mktemp)"
  awk -v s="$START" -v e="$END" '
    # Buffer the marked block instead of suppressing lines as we go, so an
    # unbalanced start marker (no matching end) does not silently swallow the
    # rest of the file: on a clean close we drop the buffer, at EOF we flush it.
    index($0, s) {
      skip = 1; buf[++n] = $0
      if (index($0, e)) { skip = 0; n = 0 }   # same-line start+end: drop it
      next
    }
    skip {
      buf[++n] = $0
      if (index($0, e)) { skip = 0; n = 0 }   # clean close: drop the block
      next
    }
    { print }
    END {
      if (skip) {                             # start marker with no end marker
        for (i = 1; i <= n; i++) print buf[i] # keep the block, do not drop tail
        print "strip-gitnexus.sh: unbalanced gitnexus markers; kept block intact" \
          > "/dev/stderr"
      }
    }
  ' "$file" |
    awk '
      /^[[:space:]]*$/ { blanks++; next }
      { while (blanks-- > 0) print ""; blanks = 0; print }
    ' >"$tmp"

  mv "$tmp" "$file"
  echo "stripped GitNexus block: $file"
}

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
strip_one "$repo_root/CLAUDE.md"
strip_one "$repo_root/AGENTS.md"
