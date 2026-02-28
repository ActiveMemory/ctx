#!/usr/bin/env bash
# ctxrc-swap — toggle .ctxrc between base (defaults) and dev (verbose logging)
#
# Usage:
#   hack/ctxrc-swap.sh          # swap to the other profile
#   hack/ctxrc-swap.sh dev      # switch to dev profile
#   hack/ctxrc-swap.sh base     # switch to base profile (all defaults)
#   hack/ctxrc-swap.sh status   # show which profile is active
#
# Source files (.ctxrc.base, .ctxrc.dev) are committed to git.
# The working copy (.ctxrc) is gitignored and treated as a sink.

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

CTXRC=".ctxrc"
CTXRC_BASE=".ctxrc.base"
CTXRC_DEV=".ctxrc.dev"

is_dev() {
    grep -q '^notify:' "$CTXRC" 2>/dev/null
}

ensure_exists() {
    if [[ ! -f "$CTXRC" ]]; then
        cp "$CTXRC_BASE" "$CTXRC"
        echo "created $CTXRC from base profile"
    fi
}

status() {
    if [[ ! -f "$CTXRC" ]]; then
        echo "active: none (.ctxrc does not exist — run 'make rc-base' or 'make rc-dev')"
    elif is_dev; then
        echo "active: dev (verbose logging enabled)"
    else
        echo "active: base (defaults)"
    fi
}

case "${1:-swap}" in
    status)
        status
        ;;
    dev)
        if [[ -f "$CTXRC" ]] && is_dev; then
            echo "already on dev profile"
        else
            cp "$CTXRC_DEV" "$CTXRC"
            echo "switched to dev profile"
        fi
        ;;
    base|prod)
        if [[ -f "$CTXRC" ]] && ! is_dev; then
            echo "already on base profile"
        else
            cp "$CTXRC_BASE" "$CTXRC"
            echo "switched to base profile"
        fi
        ;;
    swap)
        ensure_exists
        if is_dev; then
            cp "$CTXRC_BASE" "$CTXRC"
            echo "swapped: dev → base"
        else
            cp "$CTXRC_DEV" "$CTXRC"
            echo "swapped: base → dev"
        fi
        ;;
    *)
        echo "usage: hack/ctxrc-swap.sh [dev|base|swap|status]" >&2
        exit 1
        ;;
esac
