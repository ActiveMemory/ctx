#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# gitnexus-docker.sh — run GitNexus via its arm64 Docker image.
#
# Vendored from the sibling os repo's scripts/ directory.
#
# Why Docker: GitNexus pins tree-sitter@0.21.1, a native addon that does NOT
# build against this box's Node 24 (V8 API removed in Node 24). The official
# image bakes Node 22 + prebuilt arm64 native deps, so it Just Works — and
# stays off the host (no Node downgrade, no extra host attack surface).
#
# Subcommands:
#   index <repo>   analyze + register a repo. The repo is mounted at its REAL
#                  host path so the registry's absolute paths stay consistent
#                  for the MCP server. Writes <repo>/.gitnexus (gitignored).
#   mcp            stdio MCP server for Claude Code. Auto-mounts every registered
#                  repo at its real path (read from the registry) + the data
#                  volume. This is the command you register with `claude mcp add`.
#   <subcommand>   passthrough to `gitnexus <subcommand>` against the current
#                  git repo + data volume (e.g. list, status, context, impact,
#                  query, remove, doctor).
#
# Network policy (embeddings need a one-time download):
#   index  → network ENABLED so the ONNX embedding model (huggingface.co) and
#            the LadybugDB VECTOR extension (extension.ladybugdb.com) download
#            ONCE, cached in the persistent `gitnexus-cache` volume (mounted at
#            the container HOME) so it is genuinely one-time.
#   mcp    → --network none. The always-on server Claude Code spawns — the one
#            that actually chews on repo content — never gets egress; it reads
#            the cached model. That's the deliberate line: online only for the
#            conscious index step against known hosts; air-gapped for the server.
#   other  → --network none (local reads).
# NOTE: `index` currently gets full egress (Docker bridge) for the run; scope it
# to huggingface.co + extension.ladybugdb.com once the box's egress allowlist
# exists. Runs as the image's non-root `node` user.
set -euo pipefail

IMG="${GITNEXUS_IMAGE:-ghcr.io/abhigyanpatwari/gitnexus:latest}"
DATA=(-v "${GITNEXUS_VOLUME:-gitnexus-data}:/data/gitnexus")
# persists the embedding model + VECTOR extension across runs (container HOME)
CACHE=(-v "${GITNEXUS_CACHE_VOLUME:-gitnexus-cache}:/home/node")

cmd="${1:-}"; [ "$#" -gt 0 ] && shift || true

# Populate $mounts with a -v flag per registered repo path so those paths
# resolve inside the container. gitnexus PRUNES registry entries whose paths
# it cannot stat — observed live with both `index` and `list` — so every
# registry-touching invocation must see every registered repo, not just the
# current one ($1 = path to exclude, mounted separately).
registry_mounts() {
  local exclude="${1:-}" strict="${2:-}"
  mounts=()

  # Read registry.json out of the data volume (empty string when absent).
  local registry
  registry="$(docker run --rm --network none --entrypoint sh "${DATA[@]}" "$IMG" \
    -c 'cat /data/gitnexus/registry.json 2>/dev/null' || true)"

  # A fresh setup has no registry (or an empty list): legitimately no mounts.
  case "${registry//[[:space:]]/}" in
    ''|'[]') return 0 ;;
  esac

  # The registry has entries, so we MUST parse it to mount those paths. If we
  # cannot parse (python3 missing) or the parse yields nothing (corrupt JSON),
  # returning empty mounts makes a registry-touching gitnexus run PRUNE every
  # entry it can't stat in-container — the data loss this helper exists to stop.
  if ! command -v python3 >/dev/null 2>&1; then
    echo "gitnexus-docker.sh: python3 not found but registry.json has entries;" >&2
    echo "  refusing to compute empty mounts (would risk pruning the registry)." >&2
    if [ "$strict" = strict ]; then exit 1; fi
    return 0
  fi

  local paths
  paths="$(printf '%s' "$registry" | python3 -c 'import sys,json
try:
    for r in json.load(sys.stdin):
        p=r.get("path")
        if p: print(p)
except Exception:
    pass')"

  if [ -z "$paths" ]; then
    echo "gitnexus-docker.sh: registry.json has entries but none parsed to a" >&2
    echo "  path (corrupt registry?) — not mounting anything." >&2
    if [ "$strict" = strict ]; then exit 1; fi
    return 0
  fi

  local p
  while IFS= read -r p; do
    # if-form, not a bare && chain: a false filter on the LAST registry line
    # would be the loop body's exit status and set -e would kill the script
    if [ -n "$p" ] && [ "$p" != "$exclude" ] && [ -d "$p" ]; then
      mounts+=(-v "$p:$p")
    fi
  done <<< "$paths"
}

case "$cmd" in
  index)
    [ "${1:-}" ] || { echo "usage: $(basename "$0") index <repo-path>" >&2; exit 2; }
    repo="$(cd "$1" && pwd)"                       # resolve to absolute host path
    echo ">> analyze $repo (embeddings + skills, no CLAUDE.md/AGENTS.md injection; network ON for one-time model fetch)"
    # analyze can exit non-zero on a late non-fatal warning (e.g. VECTOR/embeddings)
    # even after writing a valid index — so don't let `set -e` skip the register.
    docker run --rm -w "$repo" -v "$repo:$repo" "${DATA[@]}" "${CACHE[@]}" \
      -e GITNEXUS_LBUG_EXTENSION_INSTALL=auto \
      "$IMG" gitnexus analyze "$repo" --embeddings --skills --skip-agents-md \
      || echo "   (analyze exited non-zero — will still register if the index was written)"
    if [ ! -f "$repo/.gitnexus/meta.json" ]; then
      echo ">> no .gitnexus/meta.json — analyze truly failed, not registering" >&2; exit 1
    fi
    echo ">> register $repo into the global registry"
    registry_mounts "$repo" strict
    docker run --rm --network none -w "$repo" -v "$repo:$repo" ${mounts[@]+"${mounts[@]}"} \
      "${DATA[@]}" "${CACHE[@]}" "$IMG" gitnexus index "$repo"
    echo ">> done. Remember to add '.gitnexus/' to $repo/.gitignore"
    ;;

  mcp)
    # Mount every registered repo at its real path so the registry's absolute
    # paths resolve inside the MCP container.
    registry_mounts
    exec docker run -i --rm --network none ${mounts[@]+"${mounts[@]}"} "${DATA[@]}" "${CACHE[@]}" \
      -e GITNEXUS_LBUG_EXTENSION_INSTALL=never "$IMG" gitnexus mcp
    ;;

  ""|-h|--help)
    echo "usage: $(basename "$0") {index <repo> | mcp | <gitnexus-subcommand> ...}" >&2
    [ "$cmd" = "" ] && exit 2 || exit 0
    ;;

  *)
    # passthrough (list/status/context/impact/query/remove/doctor/...).
    # strict: `list` (and other registry-touching subcommands) prune
    # entries they cannot stat, so refuse to run with empty mounts when
    # the registry is non-empty but unparseable — same data-loss guard
    # as the index branch.
    repo="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
    registry_mounts "$repo" strict
    exec docker run --rm --network none -w "$repo" -v "$repo:$repo" ${mounts[@]+"${mounts[@]}"} \
      "${DATA[@]}" "${CACHE[@]}" \
      -e GITNEXUS_LBUG_EXTENSION_INSTALL=never "$IMG" gitnexus "$cmd" "$@"
    ;;
esac
