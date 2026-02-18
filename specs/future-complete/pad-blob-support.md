# Pad Blob Support: Encrypted File Entries

## Context

The `ctx pad` currently only stores short one-liner entries. Users need to carry
portable encrypted files (e.g., `settings.local.json`, plans, credentials) that
travel with `.context/` without scp or git-tracking plaintext secrets. The
enhancement encodes files as base64 blobs within the existing single-line entry
format, preserving the encryption model with zero changes to storage or crypto.

## Design

Entry format: `label:::<base64-encoded-file-content>`

- `ctx pad list` shows: `3. my plan [BLOB]`
- `ctx pad show 3` decodes and prints full file content
- `ctx pad show 3 --out ./recovered.md` writes to disk
- `ctx pad add --file ./plan.md "my plan"` ingests a file with a label
- `ctx pad edit 3 --file ./v2.md` replaces file content, keeps label
- `ctx pad edit 3 --label "new label"` changes label, keeps file
- `ctx pad edit 3 --file ./v2.md --label "new"` replaces both
- `--append`/`--prepend` on blob entries → error (would corrupt base64)
- 64KB pre-encoding hard limit per file

## Implementation Steps

### Step 1: Create `blob.go` (new file)
**File**: `internal/cli/pad/blob.go`

Pure helper functions, no external dependencies beyond `encoding/base64`:
- `const BlobSep = ":::"` and `const MaxBlobSize = 64 * 1024`
- `isBlob(entry string) bool` — checks for `:::` separator
- `splitBlob(entry string) (label string, data []byte, ok bool)` — parses and decodes; returns `ok=false` on malformed base64 (graceful fallback)
- `makeBlob(label string, data []byte) string` — encodes and joins
- `displayEntry(entry string) string` — returns `"label [BLOB]"` for blobs, entry as-is otherwise

### Step 2: Modify `add.go`
**File**: `internal/cli/pad/add.go`

- Add `--file` / `-f` flag (`StringVarP`)
- When `--file` set: read file, check size limit, `makeBlob(label, data)`, append
- New function: `runAddBlob(cmd, label, filePath) error`
- Existing `runAdd` unchanged
- New imports: `"os"`, `"fmt"`

### Step 3: Modify `pad.go` (list display)
**File**: `internal/cli/pad/pad.go`

- In `runList`: change `entry` → `displayEntry(entry)` in the Printf (line 68)
- Update `Long` description to mention blob support

### Step 4: Modify `show.go`
**File**: `internal/cli/pad/show.go`

- Add `--out` flag (`StringVar`)
- Change `runShow` signature: `(cmd, n)` → `(cmd, n, outPath)`
- Auto-detect blob via `splitBlob`; if blob: decode and print content (or write to `--out`)
- `--out` on non-blob entry → error
- Use `cmd.Print` (not `Println`) for blob content — file may have its own trailing newline
- New import: `"os"`

### Step 5: Modify `edit.go`
**File**: `internal/cli/pad/edit.go`

- Add `--file` / `-f` and `--label` flags
- `--file`/`--label` can coexist with each other but conflict with positional/`--append`/`--prepend`
- New function: `runEditBlob(cmd, n, filePath, labelText) error` — reads existing blob, replaces file and/or label
- `--file`/`--label` on non-blob entry → error
- Add blob guard at top of `runEditAppend` and `runEditPrepend`: `isBlob(entries[n-1])` → error
- New import: `"os"`

### Step 6: Modify `resolve.go`
**File**: `internal/cli/pad/resolve.go`

- Lines 70, 77: change `entry` → `displayEntry(entry)` in Printf

### Step 7: Update `doc.go`
**File**: `internal/cli/pad/doc.go`

- Update package doc to mention file blobs and the `label:::data` format

### Step 8: Add tests in `pad_test.go`
**File**: `internal/cli/pad/pad_test.go`

~20 new test functions following existing patterns:
- **Blob helpers**: `TestIsBlob`, `TestSplitBlob_Valid`, `TestSplitBlob_NonBlob`, `TestSplitBlob_MalformedBase64`, `TestMakeBlob_Roundtrip`, `TestDisplayEntry_Blob`, `TestDisplayEntry_Plain`
- **Add blob**: `TestAdd_BlobEncrypted`, `TestAdd_BlobTooLarge`, `TestAdd_BlobFileNotFound`
- **List with blobs**: `TestList_BlobDisplay`
- **Show blobs**: `TestShow_BlobAutoDecodes`, `TestShow_BlobOutFlag`, `TestShow_OutFlagOnPlainEntry`
- **Edit blobs**: `TestEdit_BlobReplaceFile`, `TestEdit_BlobReplaceLabel`, `TestEdit_BlobReplaceBoth`, `TestEdit_AppendOnBlobErrors`, `TestEdit_PrependOnBlobErrors`, `TestEdit_LabelOnNonBlobErrors`, `TestEdit_FileAndPositionalMutuallyExclusive`

## Files Summary

| File | Change | Key Functions |
|------|--------|---------------|
| `blob.go` | NEW | `isBlob`, `splitBlob`, `makeBlob`, `displayEntry` |
| `add.go` | Modified | `runAddBlob`, `--file` flag |
| `pad.go` | Modified | `displayEntry` in list, updated description |
| `show.go` | Modified | Blob auto-decode, `--out` flag |
| `edit.go` | Modified | `runEditBlob`, `--file`/`--label` flags, blob guards |
| `resolve.go` | Modified | `displayEntry` in two print lines |
| `doc.go` | Modified | Package doc update |
| `pad_test.go` | Modified | ~20 new test functions |

**Unchanged**: `store.go`, `crypto.go`, `rm.go`, `mv.go` — blobs are just longer strings.

## Verification

```bash
CGO_ENABLED=0 go build -o ctx ./cmd/ctx
CGO_ENABLED=0 go test ./internal/cli/pad/...

# Manual smoke test
echo "secret plan content" > ./ideas/test-blob.md
./ctx pad add --file ./ideas/test-blob.md "test blob"
./ctx pad                                    # should show: 1. test blob [BLOB]
./ctx pad show 1                             # should print: secret plan content
./ctx pad show 1 --out ./ideas/recovered.md  # should write file
diff ./ideas/test-blob.md ./ideas/recovered.md
./ctx pad edit 1 --label "renamed blob"
./ctx pad                                    # should show: 1. renamed blob [BLOB]
```
