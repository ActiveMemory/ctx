#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Check if .context backup is stale (>2 days old) or SMB share is unmounted.
# Outputs non-blocking warnings to stderr, once per day.
#
# Depends on: hack/backup-context.sh touching ~/.local/state/ctx-last-backup
#             on successful backup (set -euo pipefail guarantees marker is only
#             touched after a successful copy to the SMB share).
#
# Required environment variables (same as hack/backup-context.sh):
#   CTX_BACKUP_SMB_URL - SMB share URL (e.g. smb://myhost/myshare)
#
# If CTX_BACKUP_SMB_URL is not set, the SMB mount check is skipped.

MARKER="$HOME/.local/state/ctx-last-backup"

# Use a user-specific temp directory to prevent symlink race attacks (M-3).
CTX_TMPDIR="${XDG_RUNTIME_DIR:-/tmp}/ctx-$(id -u)"
mkdir -p "$CTX_TMPDIR" && chmod 700 "$CTX_TMPDIR"

REMINDED="${CTX_TMPDIR}/backup-reminded"
MAX_AGE_DAYS=2

# Only remind once per day to avoid spam
if [ -f "$REMINDED" ]; then
  REMINDED_DATE=$(date -r "$REMINDED" +%Y%m%d 2>/dev/null || echo "")
  if [ "$REMINDED_DATE" = "$(date +%Y%m%d)" ]; then
    exit 0
  fi
fi

WARNINGS=""

# Check 1: Is the SMB share mounted?
if [ -n "${CTX_BACKUP_SMB_URL:-}" ]; then
  SMB_HOST=$(echo "$CTX_BACKUP_SMB_URL" | sed -n 's|smb://\([^/]*\)/.*|\1|p')
  SMB_SHARE=$(echo "$CTX_BACKUP_SMB_URL" | sed -n 's|smb://[^/]*/\(.*\)|\1|p')
  GVFS_MOUNT="/run/user/$(id -u)/gvfs/smb-share:server=${SMB_HOST},share=${SMB_SHARE}"
  if [ ! -d "$GVFS_MOUNT" ]; then
    WARNINGS="${WARNINGS}│ SMB share (${SMB_HOST}) is not mounted.\n"
    WARNINGS="${WARNINGS}│ Backups cannot run until it's available.\n"
  fi
fi

# Check 2: Is the backup stale?
if [ ! -f "$MARKER" ]; then
  WARNINGS="${WARNINGS}│ No backup marker found — backup may have never run.\n"
  WARNINGS="${WARNINGS}│ Run: hack/backup-context.sh\n"
else
  AGE_SECONDS=$(( $(date +%s) - $(date -r "$MARKER" +%s) ))
  AGE_DAYS=$(( AGE_SECONDS / 86400 ))
  if [ "$AGE_DAYS" -ge "$MAX_AGE_DAYS" ]; then
    WARNINGS="${WARNINGS}│ Last .context backup is ${AGE_DAYS} days old.\n"
    WARNINGS="${WARNINGS}│ Run: hack/backup-context.sh\n"
  fi
fi

# Output warnings if any
if [ -n "$WARNINGS" ]; then
  echo "" >&2
  echo "┌─ Backup Warning ────────────────────────────────" >&2
  echo -e "$WARNINGS""└──────────────────────────────────────────────────" >&2
  echo "" >&2
  touch "$REMINDED"
fi

exit 0
