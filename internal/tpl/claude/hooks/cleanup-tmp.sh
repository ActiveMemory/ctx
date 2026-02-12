#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

#
# Cleanup stale ctx temp files on session end.
# Removes files older than 15 days from the user-specific ctx tmp dir.
#

CTX_TMPDIR="${XDG_RUNTIME_DIR:-/tmp}/ctx-$(id -u)"
[ -d "$CTX_TMPDIR" ] && find "$CTX_TMPDIR" -type f -mtime +15 -delete 2>/dev/null
exit 0
