#!/usr/bin/env bash
# block-hack-scripts.sh — PreToolUse hook for Bash tool
# Blocks direct invocation of hack/ scripts; nudges toward make targets.
# Reading hack/ files (cat, head, grep, etc.) is allowed.

set -euo pipefail

input=$(cat)
command=$(printf '%s' "$input" | sed -n 's/.*"command" *: *"\(.*\)".*/\1/p' | head -1)

# Empty command — nothing to check
[ -z "$command" ] && exit 0

# Allow read-only operations on hack/ files
if printf '%s' "$command" | grep -qP '^\s*(cat|head|tail|less|read|ls|grep|rg|diff|wc|file|stat)\b'; then
    exit 0
fi

# Pattern: hack script invocation at start of command or after a separator
# Matches: ./hack/foo.sh, hack/foo.sh, bash ./hack/foo.sh, sh hack/foo.sh
# Also after && ; || |
if ! printf '%s' "$command" | grep -qP '(^|\s|&&|;|\|\||\|)\s*(bash\s+|sh\s+)?(\.\/)?hack\/\S+\.sh'; then
    exit 0
fi

# Extract the script name(s) that matched
script=$(printf '%s' "$command" | grep -oP '(\.\/)?hack\/\S+\.sh' | head -1 | sed 's|^\./||')

# Map scripts to make targets
case "$script" in
    hack/release.sh)        target="make release" ;;
    hack/build-all.sh)      target="make build-all" ;;
    hack/lint-drift.sh)     target="make lint-drift" ;;
    hack/lint-docs.sh)      target="make lint-docs" ;;
    hack/plugin-reload.sh)  target="make plugin-reload" ;;
    hack/reinstall.sh)      target="make install" ;;
    hack/gpg-fix.sh)        target="make gpg-fix / make gpg-test" ;;
    *)                      target="" ;;
esac

if [ -n "$target" ]; then
    reason="Use \`${target}\` instead of invoking \`${script}\` directly."
else
    reason="Direct hack/ script invocation blocked. Ask the user to run it manually, or create a make target first."
fi

printf '{"decision":"block","reason":"%s"}\n' "$reason"
