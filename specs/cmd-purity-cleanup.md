# cmd/ purity cleanup

Finish clearing the grandfathered map in TestCmdDirPurity so that
no cmd/ directory contains exported non-Cmd/Run functions, unexported
helpers, or type declarations.

## Scope

Move remaining grandfathered symbols from cmd/ packages to their
corresponding core/ packages:

- `BuildVault` in journal/cmd/obsidian → journal/core/obsidian
- Stale entries (`one`, `every` in remind/cmd/dismiss) already migrated;
  remove from grandfathered map

## Outcome

- Grandfathered map deleted entirely from compliance_test.go
- TestCmdDirPurity passes with zero exceptions
- Existing tests continue to pass from their original locations
