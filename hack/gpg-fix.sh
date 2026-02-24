#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

set -euo pipefail

# Fix common headless-GPG issues:
# - stale gpg-agent after sleep/rename
# - bad pinentry-program path in gpg-agent.conf (e.g., old username)
# - missing GPG_TTY in current shell
#
# Usage:
#   ./hack/gpg-fix.sh
#   ./hack/gpg-fix.sh --test

TEST=0
if [[ "${1:-}" == "--test" ]]; then TEST=1; fi

GNUPG_HOME="${GNUPGHOME:-$HOME/.gnupg}"
AGENT_CONF="$GNUPG_HOME/gpg-agent.conf"

mkdir -p "$GNUPG_HOME"
chmod 700 "$GNUPG_HOME"

# Prefer a terminal pinentry on servers.
PINENTRY=""
if command -v pinentry-curses >/dev/null 2>&1; then
  PINENTRY="$(command -v pinentry-curses)"
elif command -v pinentry >/dev/null 2>&1; then
  PINENTRY="$(command -v pinentry)"
else
  echo "ERROR: No pinentry found. Install one:"
  echo "  sudo apt-get update && sudo apt-get install -y pinentry-curses"
  exit 1
fi

# If agent conf exists and points to a non-existent pinentry program, rewrite it.
if [[ -f "$AGENT_CONF" ]]; then
  # Extract current pinentry-program if present
  CURRENT="$(awk '/^[[:space:]]*pinentry-program[[:space:]]+/{print $2; exit}' "$AGENT_CONF" || true)"
  if [[ -n "${CURRENT:-}" && ! -x "$CURRENT" ]]; then
    echo "Fixing pinentry-program: '$CURRENT' -> '$PINENTRY'"
    # Remove any existing pinentry-program lines
    grep -vE '^[[:space:]]*pinentry-program[[:space:]]+' "$AGENT_CONF" > "$AGENT_CONF.tmp" || true
    {
      echo "pinentry-program $PINENTRY"
      cat "$AGENT_CONF.tmp" 2>/dev/null || true
    } > "$AGENT_CONF"
    rm -f "$AGENT_CONF.tmp"
  fi
else
  # Create a sane default
  cat > "$AGENT_CONF" <<EOF
pinentry-program $PINENTRY
default-cache-ttl 1800
max-cache-ttl 7200
EOF
fi

chmod 600 "$AGENT_CONF"

# Refresh TTY for this shell (helps pinentry-curses).
if [[ -t 0 ]]; then
  export GPG_TTY
  GPG_TTY="$(tty)"
  export GPG_TTY
  echo "GPG_TTY=$GPG_TTY"
else
  echo "WARNING: stdin is not a TTY; pinentry-curses may not be able to prompt."
fi

# Restart agent
echo "Restarting gpg-agent..."
gpgconf --homedir "$GNUPG_HOME" --kill gpg-agent >/dev/null 2>&1 || true
gpgconf --homedir "$GNUPG_HOME" --launch gpg-agent >/dev/null 2>&1 || true

# Show what pinentry the agent will use (best-effort)
if command -v gpg-connect-agent >/dev/null 2>&1; then
  gpg-connect-agent --homedir "$GNUPG_HOME" 'GETINFO pinentry_program' /bye 2>/dev/null || true
fi

if [[ "$TEST" -eq 1 ]]; then
  echo "Testing signing..."
  echo test | gpg --homedir "$GNUPG_HOME" --clearsign >/dev/null
  echo "OK: gpg signing works."
fi

echo "Done."