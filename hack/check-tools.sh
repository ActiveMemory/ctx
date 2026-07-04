#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# check-tools.sh — verify tooling dependency versions in one pass.
#
# Reads hack/tool-versions.txt and probes every entry: PATH binaries
# (with minimum-version enforcement), npm-installed tools (with
# registry drift detection), and MCP servers registered with the
# claude CLI. Prints one table; the answer to "is this machine ready
# to work on ctx?".
#
# Verdicts: OK, MISSING, OUTDATED, SKIP (probe impossible — no
# claude CLI, no network). Only failures of tools marked `required`
# in the manifest affect the exit code.
#
# Usage:
#   ./hack/check-tools.sh            # exit 1 only on required-tool failures
#   ./hack/check-tools.sh --strict   # optional-tool failures also fatal
#
# Spec: specs/check-tools.md

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MANIFEST="${SCRIPT_DIR}/tool-versions.txt"

STRICT=false
[ "${1:-}" = "--strict" ] && STRICT=true

if [ ! -f "$MANIFEST" ]; then
    echo "manifest not found: $MANIFEST" >&2
    exit 1
fi

# Run a command under a timeout when coreutils timeout exists
# (absent on stock macOS); otherwise run it bare.
with_timeout() {
    local secs="$1"
    shift
    if command -v timeout >/dev/null 2>&1; then
        timeout "$secs" "$@"
    else
        "$@"
    fi
}

# First dotted-number token in a version banner ("go1.26.3",
# "v24.18.0", "jq-1.7.1" all yield the bare number).
extract_version() {
    grep -oE '[0-9]+(\.[0-9]+)+' | head -1
}

# True when $2 >= $1 under version ordering.
version_gte() {
    [ "$(printf '%s\n%s\n' "$1" "$2" | sort -V | head -1)" = "$1" ]
}

# `claude mcp list` is slow (it health-checks every server), so run
# it at most once and only if the manifest has mcp rows.
MCP_LIST=""
MCP_LIST_STATE="unloaded"   # unloaded | ok | no-cli | failed
load_mcp_list() {
    [ "$MCP_LIST_STATE" != "unloaded" ] && return 0
    if ! command -v claude >/dev/null 2>&1; then
        MCP_LIST_STATE="no-cli"
        return 0
    fi
    if MCP_LIST="$(with_timeout 30 claude mcp list 2>/dev/null)"; then
        MCP_LIST_STATE="ok"
    else
        MCP_LIST_STATE="failed"
    fi
}

bin_version() {
    local name="$1" out=""
    case "$name" in
        go) out="$(go version 2>/dev/null || true)" ;;
        *)  out="$("$name" --version 2>/dev/null || true)" ;;
    esac
    printf '%s' "$out" | extract_version
}

ROWS=""
FAILED_REQUIRED=0
FAILED_OPTIONAL=0

add_row() {
    ROWS="${ROWS}$(printf '%s\t%s\t%s\t%s\t%s' "$1" "$2" "$3" "$4" "$5")"$'\n'
}

fail() {
    local requirement="$1"
    if [ "$requirement" = "required" ]; then
        FAILED_REQUIRED=$((FAILED_REQUIRED + 1))
    else
        FAILED_OPTIONAL=$((FAILED_OPTIONAL + 1))
    fi
}

check_bin() {
    local name="$1" requirement="$2" min="$3"
    if ! command -v "$name" >/dev/null 2>&1; then
        add_row "$name" "$requirement" "MISSING" "-" "not on PATH"
        fail "$requirement"
        return 0
    fi
    local ver
    ver="$(bin_version "$name")"
    if [ "$min" != "-" ] && [ -n "$ver" ] && ! version_gte "$min" "$ver"; then
        add_row "$name" "$requirement" "OUTDATED" "$ver" "minimum is $min"
        fail "$requirement"
        return 0
    fi
    local note="-"
    [ "$min" != "-" ] && note="min $min"
    add_row "$name" "$requirement" "OK" "${ver:-present}" "$note"
}

check_npm() {
    local name="$1" requirement="$2"
    if ! command -v "$name" >/dev/null 2>&1; then
        add_row "$name" "$requirement" "MISSING" "-" "npm install -g $name"
        fail "$requirement"
        return 0
    fi
    local ver latest
    ver="$(bin_version "$name")"
    if ! command -v npm >/dev/null 2>&1; then
        add_row "$name" "$requirement" "OK" "${ver:-present}" "no npm; drift check skipped"
        return 0
    fi
    latest="$(with_timeout 10 npm view "$name" version 2>/dev/null || true)"
    if [ -z "$latest" ]; then
        add_row "$name" "$requirement" "OK" "${ver:-present}" "registry unreachable; drift check skipped"
    elif [ "$ver" = "$latest" ]; then
        add_row "$name" "$requirement" "OK" "$ver" "latest"
    else
        add_row "$name" "$requirement" "OUTDATED" "$ver" "latest is $latest"
        fail "$requirement"
    fi
}

check_mcp() {
    local name="$1" requirement="$2"
    load_mcp_list
    case "$MCP_LIST_STATE" in
        no-cli)
            add_row "mcp:$name" "$requirement" "SKIP" "-" "claude CLI not on PATH"
            return 0
            ;;
        failed)
            add_row "mcp:$name" "$requirement" "SKIP" "-" "claude mcp list failed"
            return 0
            ;;
    esac
    if printf '%s\n' "$MCP_LIST" | grep -q "^${name}:"; then
        local note="registered"
        if [ "$name" = "gemini-search" ] && [ -z "${GEMINI_API_KEY:-}" ]; then
            note="registered; GEMINI_API_KEY not set in this shell"
        fi
        add_row "mcp:$name" "$requirement" "OK" "-" "$note"
    else
        add_row "mcp:$name" "$requirement" "MISSING" "-" "make register-mcp"
        fail "$requirement"
    fi
}

while read -r name type requirement min _; do
    case "$name" in ''|\#*) continue ;; esac
    case "$type" in
        bin) check_bin "$name" "$requirement" "$min" ;;
        npm) check_npm "$name" "$requirement" ;;
        mcp) check_mcp "$name" "$requirement" ;;
        *)
            echo "manifest error: unknown type '$type' for '$name'" >&2
            exit 1
            ;;
    esac
done < "$MANIFEST"

printf '%s' "$ROWS" | awk -F'\t' '
BEGIN {
    printf "%-18s %-10s %-9s %-12s %s\n", "TOOL", "REQ", "STATUS", "VERSION", "NOTE"
    printf "%-18s %-10s %-9s %-12s %s\n", "----", "---", "------", "-------", "----"
}
{ printf "%-18s %-10s %-9s %-12s %s\n", $1, $2, $3, $4, $5 }
'

echo ""
if [ "$FAILED_REQUIRED" -gt 0 ]; then
    echo "FAIL: $FAILED_REQUIRED required tool(s) missing or outdated."
    exit 1
fi
if [ "$FAILED_OPTIONAL" -gt 0 ]; then
    echo "WARN: $FAILED_OPTIONAL optional tool(s) missing or outdated."
    if [ "$STRICT" = true ]; then
        echo "FAIL: --strict promotes optional failures to fatal."
        exit 1
    fi
    exit 0
fi
echo "All tools OK."
