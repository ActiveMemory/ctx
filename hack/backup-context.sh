#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Backup .context/, .claude/, and ideas/ to SMB share with timestamped archives.
# Usage: ./hack/backup-context.sh [project_dir]
#
# Creates a timestamped tarball and copies it to the remote share.
# Keeps local copy in /tmp for verification.
#
# Required environment variables (set in ~/.bashrc or similar):
#   CTX_BACKUP_SMB_URL    - SMB share URL       (e.g. smb://myhost/myshare)
#   CTX_BACKUP_SMB_SUBDIR - Subdirectory on share (e.g. ctx-sessions)

set -euo pipefail

PROJECT_DIR="${1:-$(cd "$(dirname "$0")/.." && pwd)}"
TIMESTAMP="$(date +%Y%m%d-%H%M%S)"
ARCHIVE_NAME="ctx-backup-${TIMESTAMP}.tar.gz"

# Validate required env vars
if [ -z "${CTX_BACKUP_SMB_URL:-}" ]; then
  echo "ERROR: CTX_BACKUP_SMB_URL is not set." >&2
  echo "  Set it in ~/.bashrc, e.g.:" >&2
  echo "    export CTX_BACKUP_SMB_URL=\"smb://myhost/myshare\"" >&2
  exit 1
fi

SMB_URL="${CTX_BACKUP_SMB_URL}"
SMB_SUBDIR="${CTX_BACKUP_SMB_SUBDIR:-ctx-sessions}"

# Derive GVFS mount path from SMB URL: smb://host/share -> server=host,share=share
SMB_HOST=$(echo "$SMB_URL" | sed -n 's|smb://\([^/]*\)/.*|\1|p')
SMB_SHARE=$(echo "$SMB_URL" | sed -n 's|smb://[^/]*/\(.*\)|\1|p')
GVFS_MOUNT="/run/user/$(id -u)/gvfs/smb-share:server=${SMB_HOST},share=${SMB_SHARE}"
DEST="${GVFS_MOUNT}/${SMB_SUBDIR}"

echo "==> Creating archive: ${ARCHIVE_NAME}"
tar czf "/tmp/${ARCHIVE_NAME}" \
  --exclude='.context/journal-site' \
  -C "${PROJECT_DIR}" \
  .context/ \
  .claude/ \
  ideas/ \
  -C "$HOME" \
  .bashrc

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

# Mark successful backup (used by check-backup-age.sh hook)
mkdir -p ~/.local/state
touch ~/.local/state/ctx-last-backup

# Show what's on the share
echo ""
echo "Backups on share:"
find "${DEST}" -maxdepth 1 -name 'ctx-backup-*.tar.gz' -printf '  %f %s\n' 2>/dev/null | sort
echo ""
echo "Done. Local copy: /tmp/${ARCHIVE_NAME}"
