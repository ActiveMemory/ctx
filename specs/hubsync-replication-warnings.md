# Hubsync And Replication Warnings

Governing issue: ActiveMemory/ctx#100.

## Problem

The session-start hub sync hook and hub replication loop intentionally avoid
blocking callers when a best-effort sync fails. Several failure paths currently
return silently, making operators see no synced entries without any diagnostic
that distinguishes "nothing new" from "sync failed."

## Requirements

- Preserve non-blocking behavior: failures still return an empty hubsync status
  or let the replication loop retry later.
- Emit stderr warnings via `internal/log/warn.Warn` on hubsync load, dial, sync,
  and write failures.
- Do not warn when hubsync returns zero entries successfully.
- Emit stderr warnings via `internal/log/warn.Warn` on replication-loop failures,
  including the dial failure path.
- Keep warning text local and literal unless an existing convention clearly
  requires shared constants.
- Document hubsync's warn-not-block behavior in the package documentation.
- Add focused tests for warning output where practical, especially hubsync
  load/dial failures and replication dial failure.
