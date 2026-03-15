//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// MemoryNotFound returns an error indicating that MEMORY.md was not
// discovered. Used by all memory subcommands (sync, status, diff).
//
// Returns:
//   - error: "MEMORY.md not found"
func MemoryNotFound() error {
	return errors.New(assets.TextDesc(assets.TextDescKeyErrMemoryNotFound))
}

// MemoryDiscoverFailed wraps a MEMORY.md discovery failure.
//
// Parameters:
//   - cause: the underlying discovery error.
//
// Returns:
//   - error: "MEMORY.md not found: <cause>"
func MemoryDiscoverFailed(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryDiscoverFailed), cause)
}

// MemoryDiffFailed wraps a memory diff computation failure.
//
// Parameters:
//   - cause: the underlying diff error.
//
// Returns:
//   - error: "computing diff: <cause>"
func MemoryDiffFailed(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryDiffFailed), cause)
}

// SelectContentFailed wraps a content selection failure.
//
// Parameters:
//   - cause: the underlying selection error.
//
// Returns:
//   - error: "selecting content: <cause>"
func SelectContentFailed(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemorySelectContentFailed), cause)
}

// PublishFailed wraps a publish operation failure.
//
// Parameters:
//   - cause: the underlying publish error.
//
// Returns:
//   - error: "publishing: <cause>"
func PublishFailed(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryPublishFailed), cause)
}

// ReadMemory wraps a failure to read MEMORY.md.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "reading MEMORY.md: <cause>"
func ReadMemory(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryReadMemory), cause)
}

// WriteMemory wraps a failure to write MEMORY.md.
//
// Parameters:
//   - cause: the underlying write error.
//
// Returns:
//   - error: "writing MEMORY.md: <cause>"
func WriteMemory(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryWriteMemoryTop), cause)
}

// SyncFailed wraps a sync operation failure.
//
// Parameters:
//   - cause: the underlying error from the sync operation.
//
// Returns:
//   - error: "sync failed: <cause>"
func SyncFailed(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemorySyncFailed), cause)
}

// DiscoverResolveRoot wraps a project root resolution failure.
func DiscoverResolveRoot(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryDiscoverResolveRoot), cause)
}

// DiscoverResolveHome wraps a home directory resolution failure.
func DiscoverResolveHome(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryDiscoverResolveHome), cause)
}

// DiscoverNoMemory returns an error when no auto memory file exists.
func DiscoverNoMemory(path string) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryDiscoverNoMemory), path)
}

// MemoryReadSource wraps a source file read failure during sync.
func MemoryReadSource(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryReadSource), cause)
}

// MemoryArchivePrevious wraps a failure to archive the previous mirror.
func MemoryArchivePrevious(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryArchivePrevious), cause)
}

// MemoryCreateDir wraps a failure to create the memory directory.
func MemoryCreateDir(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryCreateDir), cause)
}

// MemoryWriteMirror wraps a failure to write the mirror file.
func MemoryWriteMirror(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryWriteMirror), cause)
}

// MemoryReadMirrorArchive wraps a failure to read the mirror for archiving.
func MemoryReadMirrorArchive(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryReadMirrorArchive), cause)
}

// MemoryCreateArchiveDir wraps a failure to create the archive directory.
func MemoryCreateArchiveDir(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryCreateArchiveDir), cause)
}

// MemoryWriteArchive wraps a failure to write an archive file.
func MemoryWriteArchive(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryWriteArchive), cause)
}

// MemoryReadMirror wraps a failure to read the mirror file.
func MemoryReadMirror(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryReadMirror), cause)
}

// MemoryReadDiffSource wraps a failure to read the source for diff.
func MemoryReadDiffSource(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryReadDiffSource), cause)
}

// MemorySelectContent wraps a failure to select publish content.
func MemorySelectContent(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemorySelectContent), cause)
}

// MemoryWriteMemory wraps a failure to write MEMORY.md.
func MemoryWriteMemory(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrMemoryWriteMemory), cause)
}
