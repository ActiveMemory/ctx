#!/bin/bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Backup ~/.claude/ (global Claude Code data) to SMB share.
# Usage: ./hack/backup-global.sh
#
# Creates a timestamped tarball of ~/.claude/ (transcripts, settings, etc.)
# and copies it to the remote share. Excludes ephemeral data (todos).
#
# Required environment variables (set in ~/.bashrc or similar):
#   CTX_BACKUP_SMB_URL    - SMB share URL       (e.g. smb://myhost/myshare)
#   CTX_BACKUP_SMB_SUBDIR - Subdirectory on share (e.g. ctx-sessions)

set -euo pipefail

TIMESTAMP="$(date +%Y%m%d-%H%M%S)"
ARCHIVE_NAME="claude-global-backup-${TIMESTAMP}.tar.gz"

# Validate required env vars
if [ -z "${CTX_BACKUP_SMB_URL:-}" ]; then
  echo "ERROR: CTX_BACKUP_SMB_URL is not set." >&2
  echo "  Set it in ~/.bashrc, e.g.:" >&2
  echo "    export CTX_BACKUP_SMB_URL=\"smb://myhost/myshare\"" >&2
  exit 1
fi

SMB_URL="${CTX_BACKUP_SMB_URL}"
SMB_SUBDIR="${CTX_BACKUP_SMB_SUBDIR:-ctx-sessions}"

# Derive GVFS mount path from SMB URL
SMB_HOST=$(echo "$SMB_URL" | sed -n 's|smb://\([^/]*\)/.*|\1|p')
SMB_SHARE=$(echo "$SMB_URL" | sed -n 's|smb://[^/]*/\(.*\)|\1|p')
GVFS_MOUNT="/run/user/$(id -u)/gvfs/smb-share:server=${SMB_HOST},share=${SMB_SHARE}"
DEST="${GVFS_MOUNT}/${SMB_SUBDIR}"

echo "==> Creating archive: ${ARCHIVE_NAME}"
tar czf "/tmp/${ARCHIVE_NAME}" \
  --exclude='.claude/todos' \
  -C "$HOME" \
  .claude/

echo "    $(du -h "/tmp/${ARCHIVE_NAME}" | cut -f1) compressed"

# Mount SMB share if not already mounted
if [ ! -d "${GVFS_MOUNT}" ]; then
  echo "==> Mounting ${SMB_URL} ..."
  gio mount "${SMB_URL}"
  sleep 1
fi

if [ ! -d "${DEST}" ]; then
  echo "==> Creating ${SMB_SUBDIR}/ on share..."
  mkdir -p "${DEST}"
fi

echo "==> Copying to ${DEST}/${ARCHIVE_NAME}"
cp "/tmp/${ARCHIVE_NAME}" "${DEST}/${ARCHIVE_NAME}"

# Show what's on the share
echo ""
echo "Global backups on share:"
find "${DEST}" -maxdepth 1 -name 'claude-global-backup-*.tar.gz' -printf '  %f %s\n' 2>/dev/null | sort
echo ""
echo "Done. Local copy: /tmp/${ARCHIVE_NAME}"
