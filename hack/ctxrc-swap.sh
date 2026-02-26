#!/usr/bin/env bash
# ctxrc-swap — toggle .ctxrc between committed (production) and dev (verbose logging)
#
# Usage:
#   hack/ctxrc-swap.sh        # swap to the other profile
#   hack/ctxrc-swap.sh dev    # switch to dev profile
#   hack/ctxrc-swap.sh prod   # switch to prod profile
#   hack/ctxrc-swap.sh status # show which profile is active
#
# The dev profile lives at .ctxrc.dev. The prod profile is whatever
# git has committed. Swapping copies the current .ctxrc to the inactive
# slot and restores the other.

set -euo pipefail

CTXRC=".ctxrc"
CTXRC_DEV=".ctxrc.dev"

cd "$(git rev-parse --show-toplevel)"

is_dev() {
    # Dev profile has uncommented notify: section
    grep -q '^notify:' "$CTXRC" 2>/dev/null
}

save_dev() {
    cp "$CTXRC" "$CTXRC_DEV"
}

restore_dev() {
    if [[ ! -f "$CTXRC_DEV" ]]; then
        echo "error: no dev profile saved at $CTXRC_DEV" >&2
        exit 1
    fi
    cp "$CTXRC_DEV" "$CTXRC"
}

restore_prod() {
    git checkout -- "$CTXRC"
}

status() {
    if is_dev; then
        echo "active: dev (verbose logging enabled)"
    else
        echo "active: prod (committed defaults)"
    fi
}

case "${1:-swap}" in
    status)
        status
        ;;
    dev)
        if is_dev; then
            echo "already on dev profile"
        else
            restore_dev
            echo "switched to dev profile"
        fi
        ;;
    prod)
        if is_dev; then
            save_dev
            restore_prod
            echo "switched to prod profile (dev saved to $CTXRC_DEV)"
        else
            echo "already on prod profile"
        fi
        ;;
    swap)
        if is_dev; then
            save_dev
            restore_prod
            echo "swapped: dev → prod (dev saved to $CTXRC_DEV)"
        else
            restore_dev
            echo "swapped: prod → dev"
        fi
        ;;
    *)
        echo "usage: hack/ctxrc-swap.sh [dev|prod|swap|status]" >&2
        exit 1
        ;;
esac
