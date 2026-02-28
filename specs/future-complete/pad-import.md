# Pad Import

## Overview

Bulk-import lines from a file into the scratchpad. Each non-empty line becomes
a separate entry. Complements `pad add` (single entry) with a batch path.

## Command

```
ctx pad import FILE
ctx pad import -        # read from stdin
```

## Behavior

1. Open `FILE` (or stdin when `-` is given).
2. Read line-by-line. Skip empty lines (whitespace-only counts as empty).
3. Append each line as a new scratchpad entry.
4. Write all entries in a single `writeEntries` call (one encrypt/write cycle).
5. Print summary: `Imported N entries.`

If the file is empty or contains only blank lines, print
`No entries to import.` and exit 0.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| (none yet) | | | |

Reserved for future consideration:
- `--dry-run`: print entries that would be added without writing
- `--skip-dupes`: skip lines already present in the scratchpad

## Errors

| Condition | Message | Exit |
|-----------|---------|------|
| Missing FILE argument | cobra `ExactArgs(1)` error | 1 |
| FILE not found | `open FILE: no such file or directory` | 1 |
| Scratchpad key missing (encrypted mode, existing data) | existing `errNoKey` | 1 |
| Decryption failure | existing `errDecryptFail` | 1 |

## Implementation

- New file: `internal/cli/pad/import.go`
- Register in `pad.go`: `cmd.AddCommand(importCmd())`
- Reuse `readEntries()` / `writeEntries()` from `store.go`
- Use `bufio.Scanner` for line reading
- Stdin detection: `args[0] == "-"` → `os.Stdin`

### Pseudocode

```go
func runImport(cmd *cobra.Command, file string) error {
    var r io.Reader
    if file == "-" {
        r = os.Stdin
    } else {
        f, err := os.Open(file)
        ...
        defer f.Close()
        r = f
    }

    entries, err := readEntries()
    ...

    var count int
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        entries = append(entries, line)
        count++
    }

    if count == 0 {
        cmd.Println("No entries to import.")
        return nil
    }

    if err := writeEntries(entries); err != nil {
        return err
    }

    cmd.Printf("Imported %d entries.\n", count)
    return nil
}
```

## Tests

| Test | Scenario |
|------|----------|
| `TestImport_FromFile` | File with 3 lines → 3 entries added |
| `TestImport_SkipsEmpty` | File with blank lines → only non-empty imported |
| `TestImport_EmptyFile` | Empty file → "No entries to import." |
| `TestImport_AppendsToExisting` | Scratchpad has 2 entries, import 3 → total 5 |
| `TestImport_Stdin` | Read from stdin via `-` argument |
| `TestImport_FileNotFound` | Missing file → error |
| `TestImport_Plaintext` | Works in plaintext mode |
| `TestImport_WhitespaceOnly` | Lines with only spaces/tabs → skipped |

## Design Decisions

- **Single write**: Read all lines first, then one `writeEntries` call.
  Avoids N encrypt/write cycles for N lines.
- **TrimSpace**: Trim leading/trailing whitespace from each line to avoid
  invisible entries. The scratchpad is for one-liners, not whitespace art.
- **No `--file` flag**: Unlike `pad add --file` (blob ingestion), import
  treats the file as a list of text entries. The positional arg is the
  source file, not a blob. These are different operations.
- **Stdin support**: `-` convention is standard Unix. Enables piping
  (`grep pattern file | ctx pad import -`).
