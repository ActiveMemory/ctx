//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import "fmt"

// MemoryNotFound returns an error indicating that MEMORY.md was not
// discovered. Used by all memory subcommands (sync, status, diff).
//
// Returns:
//   - error: "MEMORY.md not found"
func MemoryNotFound() error {
	return fmt.Errorf("MEMORY.md not found")
}

// MemoryDiscoverFailed wraps a MEMORY.md discovery failure.
//
// Parameters:
//   - cause: the underlying discovery error.
//
// Returns:
//   - error: "MEMORY.md not found: <cause>"
func MemoryDiscoverFailed(cause error) error {
	return fmt.Errorf("MEMORY.md not found: %w", cause)
}

// MemoryDiffFailed wraps a memory diff computation failure.
//
// Parameters:
//   - cause: the underlying diff error.
//
// Returns:
//   - error: "computing diff: <cause>"
func MemoryDiffFailed(cause error) error {
	return fmt.Errorf("computing diff: %w", cause)
}

// SelectContentFailed wraps a content selection failure.
//
// Parameters:
//   - cause: the underlying selection error.
//
// Returns:
//   - error: "selecting content: <cause>"
func SelectContentFailed(cause error) error {
	return fmt.Errorf("selecting content: %w", cause)
}

// PublishFailed wraps a publish operation failure.
//
// Parameters:
//   - cause: the underlying publish error.
//
// Returns:
//   - error: "publishing: <cause>"
func PublishFailed(cause error) error {
	return fmt.Errorf("publishing: %w", cause)
}

// ReadMemory wraps a failure to read MEMORY.md.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "reading MEMORY.md: <cause>"
func ReadMemory(cause error) error {
	return fmt.Errorf("reading MEMORY.md: %w", cause)
}

// WriteMemory wraps a failure to write MEMORY.md.
//
// Parameters:
//   - cause: the underlying write error.
//
// Returns:
//   - error: "writing MEMORY.md: <cause>"
func WriteMemory(cause error) error {
	return fmt.Errorf("writing MEMORY.md: %w", cause)
}

// SyncFailed wraps a sync operation failure.
//
// Parameters:
//   - cause: the underlying error from the sync operation.
//
// Returns:
//   - error: "sync failed: <cause>"
func SyncFailed(cause error) error {
	return fmt.Errorf("sync failed: %w", cause)
}

// DiscoverResolveRoot wraps a project root resolution failure.
func DiscoverResolveRoot(cause error) error {
	return fmt.Errorf("resolving project root: %w", cause)
}

// DiscoverResolveHome wraps a home directory resolution failure.
func DiscoverResolveHome(cause error) error {
	return fmt.Errorf("resolving home directory: %w", cause)
}

// DiscoverNoMemory returns an error when no auto memory file exists.
func DiscoverNoMemory(path string) error {
	return fmt.Errorf("no auto memory found at %s", path)
}

// MemoryReadSource wraps a source file read failure during sync.
func MemoryReadSource(cause error) error {
	return fmt.Errorf("reading source: %w", cause)
}

// MemoryArchivePrevious wraps a failure to archive the previous mirror.
func MemoryArchivePrevious(cause error) error {
	return fmt.Errorf("archiving previous mirror: %w", cause)
}

// MemoryCreateDir wraps a failure to create the memory directory.
func MemoryCreateDir(cause error) error {
	return fmt.Errorf("creating memory directory: %w", cause)
}

// MemoryWriteMirror wraps a failure to write the mirror file.
func MemoryWriteMirror(cause error) error {
	return fmt.Errorf("writing mirror: %w", cause)
}

// MemoryReadMirrorArchive wraps a failure to read the mirror for archiving.
func MemoryReadMirrorArchive(cause error) error {
	return fmt.Errorf("reading mirror for archive: %w", cause)
}

// MemoryCreateArchiveDir wraps a failure to create the archive directory.
func MemoryCreateArchiveDir(cause error) error {
	return fmt.Errorf("creating archive directory: %w", cause)
}

// MemoryWriteArchive wraps a failure to write an archive file.
func MemoryWriteArchive(cause error) error {
	return fmt.Errorf("writing archive: %w", cause)
}

// MemoryReadMirror wraps a failure to read the mirror file.
func MemoryReadMirror(cause error) error {
	return fmt.Errorf("reading mirror: %w", cause)
}

// MemoryReadDiffSource wraps a failure to read the source for diff.
func MemoryReadDiffSource(cause error) error {
	return fmt.Errorf("reading source: %w", cause)
}

// MemorySelectContent wraps a failure to select publish content.
func MemorySelectContent(cause error) error {
	return fmt.Errorf("selecting content: %w", cause)
}

// MemoryWriteMemory wraps a failure to write MEMORY.md.
func MemoryWriteMemory(cause error) error {
	return fmt.Errorf("writing MEMORY.md: %w", cause)
}
