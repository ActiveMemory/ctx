//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

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
	return errors.New(
		assets.TextDesc(assets.TextDescKeyErrMemoryNotFound),
	)
}

// DiscoverFailed wraps a MEMORY.md discovery failure.
//
// Parameters:
//   - cause: the underlying discovery error.
//
// Returns:
//   - error: "MEMORY.md not found: <cause>"
func DiscoverFailed(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryDiscoverFailed), cause,
	)
}

// DiffFailed wraps a memory diff computation failure.
//
// Parameters:
//   - cause: the underlying diff error.
//
// Returns:
//   - error: "computing diff: <cause>"
func DiffFailed(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryDiffFailed), cause,
	)
}

// SelectContentFailed wraps a content selection failure.
//
// Parameters:
//   - cause: the underlying selection error.
//
// Returns:
//   - error: "selecting content: <cause>"
func SelectContentFailed(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemorySelectContentFailed), cause,
	)
}

// PublishFailed wraps a publish operation failure.
//
// Parameters:
//   - cause: the underlying publish error.
//
// Returns:
//   - error: "publishing: <cause>"
func PublishFailed(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryPublishFailed), cause,
	)
}

// Read wraps a failure to read MEMORY.md.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "reading MEMORY.md: <cause>"
func Read(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryReadMemory), cause,
	)
}

// Write wraps a failure to write MEMORY.md.
//
// Parameters:
//   - cause: the underlying write error.
//
// Returns:
//   - error: "writing MEMORY.md: <cause>"
func Write(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryWriteMemoryTop), cause,
	)
}

// Sync wraps a sync operation failure.
//
// Parameters:
//   - cause: the underlying error from the sync operation.
//
// Returns:
//   - error: "sync failed: <cause>"
func Sync(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemorySyncFailed), cause,
	)
}

// DiscoverResolveRoot wraps a project root resolution failure.
func DiscoverResolveRoot(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryDiscoverResolveRoot), cause,
	)
}

// DiscoverResolveHome wraps a home directory resolution failure.
func DiscoverResolveHome(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryDiscoverResolveHome), cause,
	)
}

// DiscoverNoMemory returns an error when no auto memory file exists.
func DiscoverNoMemory(path string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryDiscoverNoMemory), path,
	)
}

// ReadSource wraps a source file read failure during sync.
func ReadSource(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryReadSource), cause,
	)
}

// MemoryArchivePrevious wraps a failure to archive the previous mirror.
func MemoryArchivePrevious(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryArchivePrevious), cause,
	)
}

// MemoryCreateDir wraps a failure to create the memory directory.
func MemoryCreateDir(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryCreateDir), cause,
	)
}

// MemoryWriteMirror wraps a failure to write the mirror file.
func MemoryWriteMirror(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryWriteMirror), cause,
	)
}

// ReadMirrorArchive wraps a failure to read the mirror for archiving.
func ReadMirrorArchive(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryReadMirrorArchive), cause,
	)
}

// CreateArchiveDir wraps a failure to create the archive directory.
func CreateArchiveDir(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryCreateArchiveDir), cause,
	)
}

// WriteArchive wraps a failure to write an archive file.
func WriteArchive(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryWriteArchive), cause,
	)
}

// ReadMirror wraps a failure to read the mirror file.
func ReadMirror(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryReadMirror), cause,
	)
}

// ReadDiffSource wraps a failure to read the source for diff.
func ReadDiffSource(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryReadDiffSource), cause,
	)
}

// SelectContent wraps a failure to select publish content.
func SelectContent(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemorySelectContent), cause,
	)
}

// WriteMemory wraps a failure to write MEMORY.md.
func WriteMemory(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrMemoryWriteMemory), cause,
	)
}
