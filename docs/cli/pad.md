---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Scratchpad
icon: lucide/notebook-pen
---

![ctx](../images/ctx-banner.png)

## `ctx pad`

Encrypted scratchpad for sensitive one-liners that travel with
the project.

When invoked without a subcommand, lists all entries.

```bash
ctx pad
ctx pad <subcommand>
```

### `ctx pad add`

Append a new entry to the scratchpad.

```bash
ctx pad add <text>
ctx pad add <label> --file <path>
```

**Flags**:

| Flag     | Short | Description                                |
|----------|-------|--------------------------------------------|
| `--file` | `-f`  | Ingest a file as a blob entry (max 64 KB)  |

**Examples**:

```bash
ctx pad add "DATABASE_URL=postgres://user:pass@host/db"
ctx pad add "deploy config" --file ./deploy.yaml
```

### `ctx pad show`

Output the raw text of an entry by number. For blob entries,
prints decoded file content (or writes to disk with `--out`).

```bash
ctx pad show <n>
ctx pad show <n> --out <path>
```

**Arguments**:

- `n`: 1-based entry number

**Flags**:

| Flag    | Description                                       |
|---------|---------------------------------------------------|
| `--out` | Write decoded blob content to a file (blobs only) |

**Examples**:

```bash
ctx pad show 3
ctx pad show 2 --out ./recovered.yaml
```

### `ctx pad rm`

Remove one or more entries by stable ID. Supports individual
IDs and ranges.

```bash
ctx pad rm <id> [id...]
```

**Arguments**:

- `id`: One or more entry IDs (e.g., `3`, `1 4`, `3-5`)

**Examples**:

```bash
ctx pad rm 2
ctx pad rm 1 4
ctx pad rm 3-5
```

### `ctx pad normalize`

Reassign entry IDs as a contiguous sequence 1..N, closing any
gaps left by deletions.

**Examples**:

```bash
ctx pad normalize
```

### `ctx pad edit`

Replace, append to, or prepend to an entry.

```bash
ctx pad edit <n> [text]
```

**Arguments**:

- `n`: 1-based entry number
- `text`: Replacement text (mutually exclusive with
  `--append`/`--prepend`)

**Flags**:

| Flag        | Description                                      |
|-------------|--------------------------------------------------|
| `--append`  | Append text to the end of the entry              |
| `--prepend` | Prepend text to the beginning of entry           |
| `--file`    | Replace blob file content (preserves label)      |
| `--label`   | Replace blob label (preserves content)           |

**Examples**:

```bash
ctx pad edit 2 "new text"
ctx pad edit 2 --append " suffix"
ctx pad edit 2 --prepend "prefix "
ctx pad edit 1 --file ./v2.yaml
ctx pad edit 1 --label "new name"
```

### `ctx pad mv`

Move an entry from one position to another.

```bash
ctx pad mv <from> <to>
```

**Arguments**:

- `from`: Source position (1-based)
- `to`: Destination position (1-based)

**Examples**:

```bash
ctx pad mv 3 1      # promote entry 3 to the top
ctx pad mv 1 5      # bury entry 1 to position 5
```

### `ctx pad resolve`

Show both sides of a merge conflict in the encrypted scratchpad.

**Examples**:

```bash
ctx pad resolve
```

### `ctx pad import`

Bulk-import lines from a file into the scratchpad. Each
non-empty line becomes a separate entry. All entries are
written in a single encrypt/write cycle.

With `--blob`, import all first-level files from a directory
as blob entries. Each file becomes a blob with the filename
as its label. Subdirectories and non-regular files are skipped.

```bash
ctx pad import <file>
ctx pad import -              # read from stdin
ctx pad import --blob <dir>   # import directory files as blobs
```

**Arguments**:

- `file`: Path to a text file, `-` for stdin, or a directory
  (with `--blob`)

**Flags**:

| Flag     | Description                                        |
|----------|----------------------------------------------------|
| `--blob` | Import first-level files from a directory as blobs |

**Examples**:

```bash
ctx pad import notes.txt
grep TODO *.go | ctx pad import -
ctx pad import --blob ./ideas/
```

### `ctx pad export`

Export all blob entries from the scratchpad to a directory as
files. Each blob's label becomes the filename. Non-blob
entries are skipped.

```bash
ctx pad export [dir]
```

**Arguments**:

- `dir`: Target directory (default: current directory)

**Flags**:

| Flag        | Short | Description                                      |
|-------------|-------|--------------------------------------------------|
| `--force`   | `-f`  | Overwrite existing files instead of timestamping |
| `--dry-run` |       | Print what would be exported without writing     |

When a file already exists, a unix timestamp is prepended to
avoid collisions (e.g., `1739836200-label`). Use `--force` to
overwrite instead.

**Examples**:

```bash
ctx pad export ./ideas
ctx pad export --dry-run
ctx pad export --force ./backup
```

### `ctx pad merge`

Merge entries from one or more scratchpad files into the
current pad. Each input file is auto-detected as encrypted or
plaintext. Entries are deduplicated by exact content.

```bash
ctx pad merge FILE...
```

**Arguments**:

- `FILE...`: One or more scratchpad files to merge (encrypted
  or plaintext)

**Flags**:

| Flag        | Short | Description                                 |
|-------------|-------|---------------------------------------------|
| `--key`     | `-k`  | Path to key file for decrypting input files |
| `--dry-run` |       | Print what would be merged without writing  |

**Examples**:

```bash
ctx pad merge worktree/.context/scratchpad.enc
ctx pad merge notes.md backup.enc
ctx pad merge --key /path/to/other.key foreign.enc
ctx pad merge --dry-run pad-a.enc pad-b.md
```
