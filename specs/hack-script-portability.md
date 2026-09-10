# Spec: hack/ Script Portability — macOS Default Toolchain

## Problem

`make audit` fails before running a single real check on a stock
macOS machine (bash 3.2, BSD grep):

`hack/lint-drift.sh` — `"${exclude_args[@]}"` on an empty array
aborts under `set -u` on bash 3.2 ("unbound variable"; bash 4.4+
treats it as empty).

This failure predates any feature work and masks real findings: the
audit gate cannot run at all on contributor machines with the default
macOS toolchain.

The sibling `hack/lint-docstrings.sh` portability bugs (apostrophe in
a `$( … )` comment, `grep -cP` on BSD grep) were fixed upstream; see
`specs/lint-docstrings-macos-portability.md`.

## Fix

Minimal, behavior-preserving substitution that runs identically under
GNU and BSD toolchains:

- `${arr[@]+"${arr[@]}"}` guard for empty-array expansion.

## Non-Goals

- Rewriting the lint scripts in Go (tracked in TASKS.md: "Replace
  hack/lint-drift.sh with AST-based Go tests"; "Rewrite lint-style
  scripts in Go as ctxctl subcommands"). This spec only unblocks the
  gate until that lands.
- A full portability audit of every script under `hack/` (only the
  `make audit` chain is in scope).

## Verification

`./hack/lint-drift.sh` completes on macOS (bash 3.2.57, BSD grep)
with the same findings as on a GNU toolchain: `lint-drift: clean`.
