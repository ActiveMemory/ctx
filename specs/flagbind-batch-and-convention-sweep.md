# Flagbind Batch Helpers and Convention Sweep

## Problem

Repetitive `flagbind.StringFlagP`, `flagbind.BoolFlag`, and similar
calls across CLI command files trigger linter duplicate-code warnings
and increase maintenance surface. Separately, several hub and
initialize packages accumulated const aliases, magic numbers,
misnamed files, and predicate functions that violated project
conventions.

## Approach

Two parallel tracks in one commit:

1. **Flagbind batch helpers** — parallel-slice functions that
   register multiple flags of the same kind in a single call,
   replacing repetitive one-at-a-time registrations.
2. **Convention sweep** — rename files, remove const aliases, move
   magic numbers to config, fix predicate naming, align docstrings.

## Implementation

### Flagbind batch helpers

New file `internal/flagbind/batch.go` with six functions:

- `BindStringFlagsP` — batches `StringFlagP` calls
- `BindStringFlags` — batches `StringFlag` calls
- `BindBoolFlags` — batches `BoolFlag` calls
- `BindBoolFlagsP` — batches `BoolFlagP` calls
- `BindStringFlagShorts` — batches `StringFlagShort` calls
- `BindStringFlagsPDefault` — batches `StringFlagPDefault` calls

Each takes parallel slices (ptrs, names, shorts, descKeys) and
loops over them calling the corresponding single-flag function.

Applied to 8 CLI call sites: add, journal/source, journal/importer,
initialize, loop, pad/edit, event, notify.

### Convention sweep

| File | Change |
|------|--------|
| `hub/entry_validate.go` | Renamed to `validate_entry.go`; magic numbers moved to `config/hub` |
| `hub/errcheck.go` | Renamed to `err_check.go`; `isAuthErr` → `authErr` |
| `hub/eof.go` | `isEOF` → `eof` (no verb prefix) |
| `hub/grpc.go` | Removed aliased const; magic `clientIDBytes` → config |
| `hub/persist.go` | Removed 4 aliased consts → direct config refs |
| `hub/replicate.go` | Magic interval → `cfgHub.ReplicateInterval` |
| `hub/token.go` | Removed 3 aliased consts → direct config refs |
| `hub/store.go` | `dirPerm` → `fs.PermKeyDir` |
| `server/daemon.go` | Removed aliased `pidFile` const |
| `server/setup.go` | Removed aliased `dataDirPerm` const |
| `claudecheck/` | Renamed to `claude_check/` |
| `details.go` | Renamed to `detail.go` (singular) |
| `steering/types.go` | Docstrings aligned with conventions |
| `config/entry/entry.go` | Added `AllowedTypes` set |
| `config/hub/hub.go` | Added validation limit constants |
| `.golangci.yml` | Extended G101 exclusion to all `embed/text/` |
| `compliance_test.go` | Fixed `TestNoSecretsInTemplates` false positive |
| `sysinfo_darwin.go` | Added missing `//nolint:gosec` for G204 |
| `hub/replicate.go` | Added `var _ = startReplication` for unused scaffold |
| `hub/store_sequence.go` | Added `var _ = (*Store).lastSequence` for unused scaffold |

## Non-Goals

- Refactoring flag registration patterns beyond parallel slices
- Removing scaffolded replication code (reserved for cluster mode)
- Renaming YAML DescKey values containing "token" (user-facing text)
