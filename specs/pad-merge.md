# Pad Merge

## Overview

Merge entries from one or more scratchpad files into the current pad. Reads
each input file, auto-detects whether it is encrypted or plaintext, extracts
entries, deduplicates across all sources (including the current pad), and writes
the union back.

## Command

```
ctx pad merge FILE...
ctx pad merge scratchpad.enc               # single encrypted file
ctx pad merge scratchpad.md                 # single plaintext file
ctx pad merge pad-a.enc pad-b.md           # mixed inputs
ctx pad merge --key /other/.scratchpad.key foreign.enc
ctx pad merge --dry-run pad-a.enc
```

## Behavior

1. Load the current pad entries via `readEntries()`.
2. For each input FILE:
   a. Read the raw bytes.
   b. Auto-detect format (see Format Detection below).
   c. Parse into entries via `parseEntries()`.
3. Build a union of all entries (current pad + all inputs).
4. Deduplicate: keep only the first occurrence of each unique entry.
   Position does not matter — two identical entries are duplicates regardless
   of which file or line they come from.
5. Write the deduplicated entries via `writeEntries()`.
6. Print summary: `Merged N new entries (M duplicates skipped).`

If all input entries are duplicates of existing pad entries:
`No new entries to merge (M duplicates skipped).`

If no entries exist in any input file:
`No entries to merge.`

## Format Detection

Each input file is auto-detected as encrypted or plaintext. The detection
strategy uses a try-decrypt-first approach:

1. If a key is available (project key or `--key`):
   - Attempt decryption with the key.
   - If decryption succeeds → treat as encrypted, use decrypted entries.
   - If decryption fails → fall back to plaintext parsing.
2. If no key is available:
   - Treat as plaintext.

After plaintext parsing, if the file produced entries and any entry contains
bytes that are not valid UTF-8, print a warning:
`  ! FILE appears to contain binary data; it may be encrypted (use --key)`

This handles the common cases:
- `.enc` file from the same project → decrypts with project key.
- `.md` file → decryption fails, falls back to plaintext. Works.
- `.enc` file from another project → decryption fails with wrong key, falls
  back to plaintext (garbage). The UTF-8 warning catches this.
- Plaintext file, no key configured → plaintext parsing. Works.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--key` | `-k` | (project key) | Path to key file for decrypting input files |
| `--dry-run` | | false | Print what would be merged without writing |

## Arguments

| Arg | Required | Description |
|-----|----------|-------------|
| FILE... | Yes (1+) | Scratchpad files to merge |

## Deduplication

Deduplication is purely content-based using exact string comparison:

- **Text entries**: exact string match after parsing (entries are already
  trimmed by `parseEntries`).
- **Blob entries**: the full `label:::base64data` string is compared. Two
  blobs with the same label but different data are NOT duplicates — both
  are kept. This is correct because the base64 content differs, so the
  full entry string differs.

A seen-set (Go `map[string]bool`) tracks all entries encountered. Order of
first appearance is preserved: current pad entries first, then input files
in argument order.

### Same-label-different-data blobs

When a blob label appears in multiple inputs with different base64 data,
both entries are kept (they are distinct entries). A warning is printed:

```
  ! blob "config.json" has different content across sources; both kept
```

This is informational only — both entries are preserved.

## Output

### Normal mode

```
$ ctx pad merge worktree-a/.context/scratchpad.enc notes.md
  + check DNS config                     (from notes.md)
  + review staging deploy                (from notes.md)
  = remember to update docs              (duplicate, skipped)
  + backup.sh [BLOB]                     (from worktree-a/.context/scratchpad.enc)
Merged 3 new entries (1 duplicate skipped).
```

### Dry-run mode

```
$ ctx pad merge --dry-run notes.md
  + check DNS config                     (from notes.md)
  = remember to update docs              (duplicate, skipped)
Would merge 1 new entry (1 duplicate skipped).
```

## Errors

| Condition | Message | Exit |
|-----------|---------|------|
| No FILE arguments | cobra `MinimumNArgs(1)` error | 1 |
| FILE not found | `open FILE: no such file or directory` | 1 |
| Current pad key missing (encrypted mode, existing data) | existing `errNoKey` | 1 |
| Current pad decryption failure | existing `errDecryptFail` | 1 |
| `--key` file not found | `open KEY: no such file or directory` | 1 |
| `--key` file wrong size | `invalid key size: got N bytes, want 32` | 1 |

Input file decryption failures are NOT fatal — they fall back to plaintext
parsing. Only failures reading the current pad are fatal.

## Implementation

- New file: `internal/cli/pad/merge.go`
- Register in `pad.go`: `cmd.AddCommand(mergeCmd())`
- Reuse `readEntries()` / `writeEntries()` from `store.go`
- Reuse `parseEntries()` from `store.go`
- Reuse `displayEntry()` from `blob.go`
- Reuse `crypto.LoadKey()` / `crypto.Decrypt()` for input file decryption

### Core function: readFileEntries

```go
func readFileEntries(path string, key []byte) ([]string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    // Try decryption if key is available.
    if key != nil {
        plaintext, decErr := crypto.Decrypt(key, data)
        if decErr == nil {
            return parseEntries(plaintext), nil
        }
    }

    // Fall back to plaintext.
    return parseEntries(data), nil
}
```

### Core function: runMerge

```go
func runMerge(cmd *cobra.Command, files []string, keyFile string, dryRun bool) error {
    // 1. Load current pad.
    current, err := readEntries()
    ...

    // 2. Load key for input file decryption.
    key := loadMergeKey(keyFile)  // nil if no key available

    // 3. Build seen set from current entries.
    seen := make(map[string]bool, len(current))
    for _, e := range current {
        seen[e] = true
    }

    // 4. Track blob labels for conflict detection.
    blobLabels := buildBlobLabelMap(current)

    // 5. Process each input file.
    var added, dupes int
    var newEntries []string
    for _, file := range files {
        entries, err := readFileEntries(file, key)
        ...
        warnIfBinary(cmd, file, entries)
        for _, entry := range entries {
            if seen[entry] {
                dupes++
                cmd.Printf("  = %-40s (duplicate, skipped)\n", displayEntry(entry))
                continue
            }
            seen[entry] = true
            checkBlobConflict(cmd, entry, blobLabels)
            newEntries = append(newEntries, entry)
            added++
            cmd.Printf("  + %-40s (from %s)\n", displayEntry(entry), file)
        }
    }

    // 6. Write merged entries.
    if added == 0 { ... print summary, return }
    if dryRun { ... print would-merge summary, return }

    merged := append(current, newEntries...)
    if err := writeEntries(merged); err != nil {
        return err
    }

    cmd.Printf("Merged %d new entries (%d duplicates skipped).\n", added, dupes)
    return nil
}
```

### Key loading strategy

```go
func loadMergeKey(keyFile string) []byte {
    // Explicit --key flag takes priority.
    if keyFile != "" {
        key, err := crypto.LoadKey(keyFile)
        if err != nil {
            return nil  // will be caught later if needed
        }
        return key
    }

    // Try the project's scratchpad key.
    key, err := crypto.LoadKey(keyPath())
    if err != nil {
        return nil  // no key available, will fall back to plaintext
    }
    return key
}
```

## Tests

| Test | Scenario |
|------|----------|
| `TestMerge_Basic` | Merge file with 3 entries, 1 duplicate → 2 new added |
| `TestMerge_AllDuplicates` | All entries already in pad → "No new entries" |
| `TestMerge_EmptyFile` | Empty input file → "No entries to merge." |
| `TestMerge_MultipleFiles` | Two files merged, cross-file dedup works |
| `TestMerge_EncryptedInput` | Encrypted .enc file decrypted and merged |
| `TestMerge_PlaintextFallback` | Decryption fails, falls back to plaintext |
| `TestMerge_MixedEncPlain` | One encrypted, one plaintext in same call |
| `TestMerge_DryRun` | --dry-run prints summary without writing |
| `TestMerge_CustomKey` | --key flag used for foreign encrypted file |
| `TestMerge_BlobEntries` | Blob entries merged and deduplicated |
| `TestMerge_BlobConflict` | Same label, different data → warning, both kept |
| `TestMerge_BinaryWarning` | Non-UTF-8 content triggers warning |
| `TestMerge_FileNotFound` | Missing file → error |
| `TestMerge_EmptyPadMerge` | Empty current pad + file with entries → entries added |
| `TestMerge_PlaintextMode` | Works when scratchpad encryption is disabled |
| `TestMerge_PreservesOrder` | Current entries stay first, new entries appended |

## Design Decisions

- **Try-decrypt-first**: Simpler than requiring the user to specify format.
  Decryption either works or it doesn't — no ambiguity. The UTF-8 warning
  catches the case where an encrypted file is mistakenly parsed as plaintext.
- **Content-based dedup**: Position doesn't matter, only content. Two entries
  "foo" at line 1 in file A and line 5 in file B are the same entry. This
  matches the user's mental model of a scratchpad as an unordered set of notes.
- **Non-fatal input decryption**: If one input file can't be decrypted, fall
  back rather than abort. The user may be merging a mix of formats.
- **Append-only**: New entries are appended after existing entries. This
  preserves the user's current pad ordering.
- **Single write cycle**: All entries (old + new) are written in one call to
  `writeEntries`, avoiding multiple encrypt/write cycles.
- **`--key` deferred to v1**: Including it from the start because it's simple
  to implement and the workaround (manually decrypt + import) is clunky.

## Non-Goals

- **Interactive conflict resolution**: Same-label-different-data blobs are
  both kept with a warning. Interactive selection is a future enhancement.
- **Three-way merge**: This is a union merge, not a diff-based merge. There
  is no concept of a "base" version.
- **Entry editing during merge**: The merge only adds entries. It does not
  modify or remove existing entries.
