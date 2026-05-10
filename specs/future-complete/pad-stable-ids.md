# Stable IDs and Batch Operations for Pad and Remind

## Problem

`ctx pad rm` shifts entry indices after deletion. Chaining
multiple `rm` commands (`rm 10; rm 11; rm 12`) deletes the
wrong entries because indices shift between operations.

`ctx remind dismiss` only accepts a single ID or `--all`.
No batch dismiss (`dismiss 3 4 5`) or range (`dismiss 3-5`).

## Approach

Add stable auto-incrementing IDs to pad entries. IDs are line
prefixes in the format `[N] content`. IDs never shift or get
reused. Gaps are expected after deletions.

## Behavior

### Entry format

```
[1] zoom setup for kevin
[2] don't forget to back up broadcom vm's gpg keys
[5] check recent 20 journal entries
```

### ID assignment

- New entries get `max(existing IDs) + 1`
- `--prepend` inserts at file position 1 with the next ID
- `--append` inserts at end with the next ID

### Parsing rules

- `^\[(\d+)\] ` matches a valid ID prefix
- No match → assign next available ID (auto-repair)
- Duplicate ID → keep first occurrence, reassign later ones
- IDs are display-stable: `pad show 5` always means entry 5

### Multi-delete

- `ctx pad rm 10 11 12` — resolve all IDs before any deletion
- `ctx pad rm 10-12` — range syntax, inclusive

### Migration

First load of an old pad (no ID prefixes): assign IDs 1..N to
existing entries in file order, write back. One-time, automatic.

### Display

```
  1. zoom setup for kevin
  2. don't forget to back up broadcom vm's gpg keys
  5. check recent 20 journal entries
```

Current display already uses `N.` format. Just use the stable
ID instead of the line number.

### Normalize

`ctx pad normalize` reassigns IDs as 1..N in current file order,
closing all gaps. This is a deliberate user action — it
invalidates any previously-seen IDs, so it should not run
automatically.

```
Before:               After:
[1] first             [1] first
[5] second            [2] second
[12] third            [3] third
```

## Remind Enhancements

### Batch dismiss

`ctx remind dismiss 3 4 5` — resolve all IDs, then dismiss.
`ctx remind dismiss 3-5` — range syntax, inclusive.

Both resolve IDs before any deletion. Already has stable IDs.

### Normalize

`ctx remind normalize` — reassign IDs as 1..N in current
order, closing gaps from dismissed reminders.

## Non-Goals

- Separate index file for pad (one file is simpler)
- Counter header line for pad (max + 1 is sufficient)
- Changing the encryption scheme
