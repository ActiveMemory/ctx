---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Scratchpad
icon: lucide/sticky-note
---

![ctx](../images/ctx-banner.png)

## What Is `ctx` Scratchpad?

A one-liner scratchpad, encrypted at rest, synced via `git`.

Quick notes that don't fit decisions, learnings, or tasks: reminders,
intermediate values, sensitive tokens, working memory during debugging.
Entries are numbered, reorderable, and persist across sessions.

## Encrypted by Default

Scratchpad entries are encrypted with `AES-256-GCM` before touching the disk.

| Component      | Path                      | Git status                     |
|----------------|---------------------------|--------------------------------|
| Encryption key | `.context/.ctx.key`       | Gitignored, `0600` permissions |
| Encrypted data | `.context/scratchpad.enc` | Committed                      |

The key is generated automatically during `ctx init` (256-bit via
`crypto/rand`). The ciphertext format is `[12-byte nonce][ciphertext+tag]`.
No external dependencies: Go stdlib only.

Because the key is `.gitignore`d and the data is committed, you get:

* **At-rest encryption**: the `.enc` file is opaque without the key
* **Git sync**: push/pull the encrypted file like any other tracked file
* **Key separation**: the key never leaves the machine unless you copy it

## Commands

| Command                           | Purpose                                                |
|-----------------------------------|--------------------------------------------------------|
| `ctx pad`                         | List all entries (numbered 1-based)                    |
| `ctx pad show N`                  | Output raw text of entry N (no prefix, pipe-friendly)  |
| `ctx pad add "text"`              | Append a new entry                                     |
| `ctx pad rm N`                    | Remove entry at position N                             |
| `ctx pad edit N "text"`           | Replace entry N with new text                          |
| `ctx pad edit N --append "text"`  | Append text to the end of entry N                      |
| `ctx pad edit N --prepend "text"` | Prepend text to the beginning of entry N               |
| `ctx pad add TEXT --file PATH`    | Ingest a file as a blob entry (TEXT is the label)      |
| `ctx pad show N --out PATH`       | Write decoded blob content to a file                   |
| `ctx pad mv N M`                  | Move entry from position N to position M               |
| `ctx pad resolve`                 | Show both sides of a merge conflict for resolution     |
| `ctx pad import FILE`             | Bulk-import lines from a file (or stdin with `-`)      |
| `ctx pad import --blobs DIR`      | Import directory files as blob entries                  |
| `ctx pad export [DIR]`            | Export all blob entries to a directory as files        |
| `ctx pad merge FILE...`           | Merge entries from other scratchpad files into current |

All commands decrypt on read, operate on plaintext in memory, and
re-encrypt on write. The key file is never printed to stdout.

### Examples

```bash
# Add a note
ctx pad add "check DNS propagation after deploy"

# List everything
ctx pad
#   1. check DNS propagation after deploy
#   2. staging API key: sk-test-abc123

# Show raw text (for piping)
ctx pad show 2
# sk-test-abc123

# Compose entries
ctx pad edit 1 --append "$(ctx pad show 2)"

# Reorder
ctx pad mv 2 1

# Clean up
ctx pad rm 2
```

## Bulk Import and Export

Import lines from a file in bulk (*each non-empty line becomes an entry*):

```bash
# Import from a file
ctx pad import notes.txt

# Import from stdin
grep TODO *.go | ctx pad import -
```

Export all blob entries to a directory as files:

```bash
# Export to a directory
ctx pad export ./ideas

# Preview without writing
ctx pad export --dry-run

# Overwrite existing files
ctx pad export --force ./backup
```

## Merging Scratchpads

Combine entries from other scratchpad files into your current pad. Useful
when merging work from parallel worktrees, other machines, or teammates:

```bash
# Merge from a worktree's encrypted scratchpad
ctx pad merge worktree/.context/scratchpad.enc

# Merge from multiple sources (encrypted and plaintext)
ctx pad merge pad-a.enc notes.md

# Merge a foreign encrypted pad using its key
ctx pad merge --key /other/.ctx.key foreign.enc

# Preview without writing
ctx pad merge --dry-run pad-a.enc pad-b.md
```

Each input file is auto-detected as encrypted or plaintext: decryption is
attempted first, and on failure the file is parsed as plain text. Entries are
deduplicated by exact content, so running merge twice with the same file is
safe.

## File Blobs

The scratchpad can store small files (*up to 64 KB*) as blob entries. Files
are base64-encoded and stored with a human-readable label.

```bash
# Ingest a file: first argument is the label
ctx pad add "deploy config" --file ./deploy.yaml

# Listing shows label with a [BLOB] marker
ctx pad
#   1. check DNS propagation after deploy
#   2. deploy config [BLOB]

# Extract to a file
ctx pad show 2 --out ./recovered.yaml

# Or print decoded content to stdout
ctx pad show 2
```

Blob entries are encrypted identically to text entries. The internal format
is `label:::base64data`: You **never** need to construct this manually.

| Constraint                   | Value                      |
|------------------------------|----------------------------|
| Max file size (pre-encoding) | 64 KB                      |
| Storage format               | `label:::base64(content)`  |
| Display                      | `label [BLOB]` in listings |

!!! tip "When Should You Use Blobs"
    Blobs are for small files you want encrypted and portable: config
    snippets, key fragments, deployment manifests, test fixtures. For
    anything larger than 64 KB, use the filesystem directly.

## Using with AI

!!! tip "Use Natural Language"
    As in many `ctx` features, the `ctx` scratchpad can also be used with
    natural langauge. You don't have to memorize the CLI commands.

    CLI gives you "*precision*", whereas natural language gives you **flow**.

The `/ctx-pad` skill maps natural language to `ctx pad` commands. You
don't need to remember the syntax:

| You say                                        | What happens                                     |
|------------------------------------------------|--------------------------------------------------|
| "*jot down: check DNS after deploy*"           | `ctx pad add "check DNS after deploy"`           |
| "*show my scratchpad*"                         | `ctx pad`                                        |
| "*delete the third entry*"                     | `ctx pad rm 3`                                   |
| "*update entry 2 to include the new endpoint*" | `ctx pad edit 2 "..."`                           |
| "*move entry 4 to the top*"                    | `ctx pad mv 4 1`                                 |
| "*import my notes from notes.txt*"             | `ctx pad import notes.txt`                       |
| "*export all blobs to ./backup*"               | `ctx pad export ./backup`                        |
| "*merge the scratchpad from the worktree*"     | `ctx pad merge worktree/.context/scratchpad.enc` |

The skill handles the translation. You describe what you want in plain
English; the agent picks the right command.

## Worktrees

The encryption key (`.context/.ctx.key`) is gitignored. It only
exists in the main checkout. In a git worktree, `ctx pad` commands fail
gracefully: no key is found, no crash, but the pad is inaccessible.

Use the scratchpad from the main checkout only. This is by design: the
key never leaves the machine where it was generated, and worktrees are
ephemeral working copies.

## Key Distribution

The encryption key (`.context/.ctx.key`) stays on the machine
where it was generated. `ctx` **never** transmits it.

To share the scratchpad across machines:

1. Copy the key manually: `scp`, USB drive, password manager.
2. Push/pull the `.enc` file via git as usual.
3. Both machines can now read and write the same scratchpad.

!!! warning "Never Commit the Key"
    The key is `.gitignore`d by default. If you override this, anyone with
    repo access can decrypt your scratchpad. 

    Treat the key like an SSH private key.

See the [Syncing Scratchpad Notes Across Machines](../recipes/scratchpad-sync.md)
recipe for a step-by-step walkthrough.

## Plaintext Override

For projects where encryption is unnecessary, disable it in `.ctxrc`:

```yaml
scratchpad_encrypt: false
```

In plaintext mode:

* Entries are stored in `.context/scratchpad.md` instead of `.enc`.
* No key is generated or required.
* All `ctx pad` commands work identically.
* The file is human-readable and diffable.

!!! tip "When Should You Use Plaintext"
    Plaintext mode is useful for non-sensitive projects, solo work where
    encryption adds friction, or when you want scratchpad entries visible
    in `git diff`.

## When Should You Use Scratchpad versus Context Files

| Use case                                            | Where it goes    |
|-----------------------------------------------------|------------------|
| Temporary reminders ("*check X after deploy*")      | Scratchpad       |
| Working values during debugging                     | Scratchpad       |
| Sensitive tokens or API keys (short-term)           | Scratchpad       |
| Quick notes that don't fit anywhere else            | Scratchpad       |
| Items that are not directly relevant to the project | Scratchpad       |
| Things that you want to keep near, but also hidden  | Scratchpad       |
| Work items with completion tracking                 | `TASKS.md`       |
| Trade-offs with rationale                           | `DECISIONS.md`   |
| Reusable lessons with context/lesson/application    | `LEARNINGS.md`   |
| Codified patterns and standards                     | `CONVENTIONS.md` |

**Rule of thumb**: 

* If it needs structure or will be referenced months later, use
  a **context file** (*i.e. `DECISIONS.md`, `LEARNINGS.md`, `TASKS.md`*). 
* If it is **working memory** for the **current** session or week, use 
  **the scratchpad**.

## See Also

* [Syncing Scratchpad Notes Across Machines](../recipes/scratchpad-sync.md):
  Key distribution, push/pull workflow, merge conflict resolution
* [Using the Scratchpad](../recipes/scratchpad-with-claude.md):
  Natural language examples, blob workflow, when to use scratchpad vs context files
* [Context Files](../home/context-files.md): Format and conventions for all
  `.context/` files
* [Security](../security/reporting.md): Trust model and permission hygiene
