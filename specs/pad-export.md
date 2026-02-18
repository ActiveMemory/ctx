# Pad Export

## Overview

Export blob entries from the scratchpad to a directory as files. Each blob's
label becomes the filename. Complements `pad import` (batch ingest) with a
batch extract path.

## Command

```
ctx pad export [DIR]
ctx pad export              # defaults to current directory
ctx pad export ./ideas      # exports to ./ideas/
```

## Behavior

1. Call `readEntries()` to load all scratchpad entries.
2. Iterate entries. Skip non-blob entries (`!isBlob(entry)`).
3. For each blob entry, call `splitBlob(entry)` to get label and data.
4. Determine output path: `DIR/label`.
5. If file already exists: prepend unix timestamp (`DIR/1739836200-label`).
6. Write decoded data to file.
7. Print summary: `Exported N blobs.`

If no blob entries exist, print `No blob entries to export.` and exit 0.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--force` | `-f` | false | Overwrite existing files instead of timestamping |
| `--dry-run` | | false | Print what would be exported without writing |

## Arguments

| Arg | Required | Default | Description |
|-----|----------|---------|-------------|
| DIR | No | `.` | Target directory for exported files |

## Collision Handling

Default behavior when `DIR/label` already exists:

- Prepend unix timestamp: `DIR/1739836200-label`
- Print warning: `  ! label exists, writing as 1739836200-label`

With `--force`: overwrite without warning.

## Output

```
$ ctx pad export ./ideas
  + settings.local.json
  + ctx-recall-usage.md
  ! marketing.md exists, writing as 1739836200-marketing.md
Exported 3 blobs.
```

With `--dry-run`:

```
$ ctx pad export ./ideas --dry-run
  settings.local.json → ./ideas/settings.local.json
  ctx-recall-usage.md → ./ideas/ctx-recall-usage.md
  marketing.md → ./ideas/1739836200-marketing.md (exists)
Would export 3 blobs.
```

## Errors

| Condition | Message | Exit |
|-----------|---------|------|
| DIR does not exist and cannot be created | `mkdir DIR: permission denied` | 1 |
| Scratchpad key missing (encrypted mode) | existing `errNoKey` | 1 |
| Decryption failure | existing `errDecryptFail` | 1 |
| Write failure on individual file | Warning, continue to next | 0 |

Individual file write failures warn and continue (same as the bash script
pattern). The command succeeds as long as the scratchpad can be read.

## Implementation

- New file: `internal/cli/pad/export.go`
- Register in `pad.go`: `cmd.AddCommand(exportCmd())`
- Reuse `readEntries()` from `store.go`
- Reuse `isBlob()`, `splitBlob()` from `blob.go`
- Use `os.MkdirAll` for target directory creation
- Collision detection: `os.Stat` before write

### Pseudocode

```go
func runExport(cmd *cobra.Command, dir string, force, dryRun bool) error {
    entries, err := readEntries()
    ...

    if err := os.MkdirAll(dir, 0o750); err != nil {
        return err
    }

    var count int
    for _, entry := range entries {
        label, data, ok := splitBlob(entry)
        if !ok {
            continue
        }

        outPath := filepath.Join(dir, label)
        if !force {
            if _, err := os.Stat(outPath); err == nil {
                ts := fmt.Sprintf("%d", time.Now().Unix())
                outPath = filepath.Join(dir, ts+"-"+label)
                cmd.Printf("  ! %s exists, writing as %s\n", label, filepath.Base(outPath))
            }
        }

        if dryRun {
            cmd.Printf("  %s → %s\n", label, outPath)
            count++
            continue
        }

        if err := os.WriteFile(outPath, data, 0o600); err != nil {
            cmd.PrintErrf("  ! failed to write %s: %v\n", label, err)
            continue
        }

        cmd.Printf("  + %s\n", label)
        count++
    }

    if count == 0 {
        cmd.Println("No blob entries to export.")
        return nil
    }

    verb := "Exported"
    if dryRun {
        verb = "Would export"
    }
    cmd.Printf("%s %d blobs.\n", verb, count)
    return nil
}
```

## Tests

| Test | Scenario |
|------|----------|
| `TestExport_Basic` | 2 blobs + 1 text entry → 2 files written |
| `TestExport_EmptyPad` | No entries → "No blob entries to export." |
| `TestExport_NoBlobsOnly` | Text-only entries → "No blob entries to export." |
| `TestExport_CollisionTimestamp` | Existing file → timestamped filename |
| `TestExport_Force` | Existing file + --force → overwritten |
| `TestExport_DryRun` | --dry-run → no files written, summary printed |
| `TestExport_DirCreated` | Non-existent dir → created automatically |
| `TestExport_WriteError` | Unwritable dir → warning, continues |
| `TestExport_Plaintext` | Works in plaintext mode |
| `TestExport_FilePermissions` | Exported files are 0o600 |

## Design Decisions

- **Default to current directory**: Unlike import (which requires a source),
  export has a sensible default — the directory you're standing in.
- **Timestamp collision avoidance**: Overwriting is dangerous for a security
  tool. The default is safe; `--force` is opt-in.
- **Skip non-blobs silently**: Text entries are not files. No warning needed.
- **Continue on write failure**: One bad file should not prevent exporting
  the rest. Matches the import pattern of warn-and-continue.
- **0o600 permissions**: Blobs may contain sensitive data (that's why they're
  in the encrypted pad). Exported files should be owner-only.
