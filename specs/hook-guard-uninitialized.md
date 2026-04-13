# Hook Guard: Silent Exit in Uninitialized Projects

## Problem

When the ctx plugin is installed globally and Claude runs in a
non-ctx project, unsolicited relay alerts fire from `check_resource`
and `check_backup_age`. Users see "Load Xx CPU count" and backup-age
warnings in projects that don't use ctx at all.

## Approach

Add the `state.Initialized()` guard to the two user-visible relay
hooks that were missing it. Each hook is responsible for its own
no-op behavior when ctx is not initialized — this matches the
existing pattern in 18 other hooks.

Scope limited to hooks that emit user-visible relay alerts. Safety
hooks (`block_dangerous_command`, `block_non_path_ctx`) intentionally
run regardless of ctx state.

## Non-Goals

- Full audit of all 28 hooks (follow-up work)
- Centralized guard in bootstrap (blast radius too wide)
- Changing hooks that already have the guard
