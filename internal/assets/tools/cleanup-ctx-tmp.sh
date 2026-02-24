#!/bin/bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

#
# Manual cleanup of ctx temp files.
# Usage: .context/tools/cleanup-ctx-tmp.sh [max_age_days]
#   max_age_days: Remove files older than this (default: 15)
#
# Safe to run while sessions are active — only removes stale files.
# Can be scheduled via cron for long-running servers.
#

MAX_AGE="${1:-15}"
CTX_TMPDIR="${XDG_RUNTIME_DIR:-/tmp}/ctx-$(id -u)"

if [ ! -d "$CTX_TMPDIR" ]; then
  echo "No ctx temp directory found at $CTX_TMPDIR"
  exit 0
fi

COUNT=$(find "$CTX_TMPDIR" -type f -mtime +"$MAX_AGE" 2>/dev/null | wc -l)
if [ "$COUNT" -eq 0 ]; then
  echo "No stale files (older than ${MAX_AGE} days) in $CTX_TMPDIR"
  exit 0
fi

find "$CTX_TMPDIR" -type f -mtime +"$MAX_AGE" -delete 2>/dev/null
echo "Removed $COUNT stale file(s) from $CTX_TMPDIR"
