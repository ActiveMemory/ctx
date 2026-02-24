# ctx pad — Encrypted Scratchpad

## Problem

Developers need a place for short, sensitive one-liners that travel with the
project via git but remain opaque at rest. Existing context files (TASKS.md,
DECISIONS.md) are plaintext by design — the AI reads them directly. There's no
place for notes that should be encrypted in the repo and decrypted only by
those who hold the key.

## Solution

A scratchpad stored in `.context/scratchpad.enc`, encrypted with AES-256-GCM
using a symmetric key at `.context/.context.key`. The key is gitignored;
the encrypted file is committed. Only machines with the key can read or write.

A plaintext fallback (`.context/scratchpad.md`) is available via config for
users who don't need encryption.

## Design Principles

- **Simple**: one-liners, no types, no categories. It's a scratchpad.
- **Encrypted by default**: opt out via `.contextrc`, not opt in.
- **Key never in stdout**: the key file is a file. Copy it. Never print it.
- **Git-native**: encrypted file commits like any other file.
- **Skill-mediated**: Claude interacts via `ctx pad` CLI, never reads the
  file directly.

## Storage

```
.context/scratchpad.enc      # committed to git (encrypted, default)
.context/scratchpad.md       # committed to git (plaintext, config override)
.context/.context.key     # gitignored, 0600, created at ctx init
```

Only one of `.enc` / `.md` exists at a time, depending on config.

## Encryption

- Algorithm: AES-256-GCM (Go stdlib `crypto/aes`, `crypto/cipher`)
- Key: 256-bit random (`crypto/rand`), stored raw in `.context.key`
- Nonce: 12-byte random, prepended to ciphertext
- File format: `[12-byte nonce][ciphertext + 16-byte GCM tag]`
- Plaintext format: newline-delimited UTF-8 lines (entries)
- Each write re-encrypts the entire file (scratchpad is small)

No external dependencies. Go stdlib only.

## Configuration

In `.contextrc` (YAML):

```yaml
scratchpad_encrypt: true   # default: true
```

When `false`:
- No key is generated
- Scratchpad is stored as `.context/scratchpad.md` (plain markdown)
- Same CLI, same skill, same commands — only storage differs

## CLI Commands

```
ctx pad                       # list all entries (numbered, 1-based)
ctx pad show N                # output raw text of entry N (1-based, no numbering prefix)
ctx pad add "..."             # append entry
ctx pad rm N                  # delete entry at position N
ctx pad edit N "..."          # replace entry at position N
ctx pad edit N --append "..."  # append text to entry at position N
ctx pad edit N --prepend "..." # prepend text to entry at position N
ctx pad mv N M                # move entry from position N to position M
```

### Append / Prepend Flags

`ctx pad edit` supports `--append` and `--prepend` flags as an alternative to
full replacement:

- `--append "text"` — concatenates a space and the given text to the **end** of
  the existing entry: `old + " " + text`.
- `--prepend "text"` — concatenates the given text and a space to the
  **beginning** of the existing entry: `text + " " + old`.

The flags are **mutually exclusive** with each other and with the positional
replacement text argument. Providing more than one is an error.

Examples:

```
# Full replacement (existing behavior):
ctx pad edit 2 "new text"

# Append to entry 2:
ctx pad edit 2 --append "updated at 3pm"
# If entry 2 was "check DNS", it becomes "check DNS updated at 3pm"

# Prepend to entry 2:
ctx pad edit 2 --prepend "URGENT:"
# If entry 2 was "check DNS", it becomes "URGENT: check DNS"
```

### Output Format

`ctx pad` (list):
```
  1. remember to check DNS config on staging
  2. API rate limit is 1000/min not 500
  3. deploy window is Tuesday 2-4am UTC
```

`ctx pad show 3`:
```
deploy window is Tuesday 2-4am UTC
```

Raw text only, no numbering prefix. Designed for pipe composability:
```
ctx pad edit 1 --append "$(ctx pad show 3)"
```

`ctx pad add "..."`:
```
Added entry 4.
```

`ctx pad rm 2`:
```
Removed entry 2.
```

### Error Cases

- No key, encrypted file exists: `Encrypted scratchpad found but no key. Place your key at .context/.context.key`
- No scratchpad file exists: `Scratchpad is empty.` (or create on first add)
- Index out of range: `Entry N does not exist. Scratchpad has M entries.`
- Decryption failure (wrong key): `Decryption failed. Wrong key?`

## ctx init Behavior

When `scratchpad_encrypt: true` (default):
1. Generate 256-bit random key → `.context/.context.key` (mode 0600)
2. Add `.context/.context.key` to `.gitignore` if not present
3. Print: `Scratchpad key created at .context/.context.key`
4. Print: `Copy this file to your other machines at the same path.`
5. If key already exists: skip generation (idempotent)
6. If `scratchpad.enc` exists but no key: warn, don't overwrite

When `scratchpad_encrypt: false`:
1. No key generated
2. Create empty `.context/scratchpad.md` if not exists

## Key Distribution

The key is a file. The user copies it. No CLI commands to print, export,
or import the key. The key never appears in stdout, session transcripts,
journal entries, or any AI-visible output.

Suggested methods (documented, not enforced):
- `scp .context/.context.key user@host:project/.context/`
- Copy via password manager
- Any secure file transfer

## Merge Conflicts

The encrypted file is binary to git. Conflicts manifest as "both modified."

Resolution strategy (implemented in the skill, not the CLI):
1. Claude detects conflict state
2. Both versions are decryptable (same key on both machines)
3. `ctx pad resolve` decrypts ours + theirs, presents both lists
4. Claude asks the user: keep yours / take theirs / merge interactively
5. User picks, Claude writes the merged result

The CLI provides `ctx pad resolve` which outputs both sides. The skill
mediates the human decision.

## Skill: `/ctx-pad`

Wraps `ctx pad` commands. Claude maps natural language to CLI:

- "add a note: check DNS" → `ctx pad add "check DNS"`
- "show my scratchpad" → `ctx pad`
- "delete the third one" → `ctx pad rm 3`
- "move the last one to the top" → `ctx pad mv N 1`
- "change entry 2 to ..." → `ctx pad edit 2 "..."`
- "add 'URGENT' to the start of entry 3" → `ctx pad edit 3 --prepend "URGENT"`
- "append 'done' to entry 1" → `ctx pad edit 1 --append "done"`

The skill has `allowed-tools: Bash(ctx:*)`.

## Package Structure

```
internal/crypto/          # encrypt, decrypt, key generation
internal/crypto/crypto_test.go
internal/cli/pad/         # cobra commands: root, add, rm, edit, mv, resolve
internal/cli/pad/doc.go
internal/cli/pad/pad_test.go
internal/tpl/claude/skills/ctx-pad/SKILL.md
```

## Permissions and Constants

New constants in `internal/config/`:
- `FileScratchpadEnc = "scratchpad.enc"`
- `FileScratchpadMd = "scratchpad.md"`
- `FileContextKey = ".context.key"`
- `PermSecret = 0600` (for key file)

New field in `internal/rc/types.go`:
- `ScratchpadEncrypt bool \`yaml:"scratchpad_encrypt"\`` (default true)

New permission in `DefaultClaudePermissions`:
- `"Bash(ctx pad:*)"`

## Non-Goals

- No categories or types (it's a scratchpad)
- No key export/import commands (copy the file)
- No automatic merge (skill-mediated only)
- No encryption of existing context files (separate concern)
- No passphrase-derived keys (file-based only)
