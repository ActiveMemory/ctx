//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// NotFound returns an error indicating that MEMORY.md was not
// discovered. Used by all memory subcommands (sync, status, diff).
//
// Returns:
//   - error: "MEMORY.md not found"
func NotFound() error {
	return errors.New(
		desc.Text(text.DescKeyErrMemoryNotFound),
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
		desc.Text(text.DescKeyErrMemoryDiscoverFailed), cause,
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
		desc.Text(text.DescKeyErrMemoryDiffFailed), cause,
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
		desc.Text(text.DescKeyErrMemorySelectContentFailed), cause,
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
		desc.Text(text.DescKeyErrMemoryPublishFailed), cause,
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
		desc.Text(text.DescKeyErrMemoryReadMemory), cause,
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
		desc.Text(text.DescKeyErrMemoryWriteMemoryTop), cause,
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
		desc.Text(text.DescKeyErrMemorySyncFailed), cause,
	)
}

// DiscoverResolveRoot wraps a project root resolution failure.
//
// Parameters:
//   - cause: the underlying resolution error
//
// Returns:
//   - error: "resolving project root: <cause>"
func DiscoverResolveRoot(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryDiscoverResolveRoot), cause,
	)
}

// DiscoverResolveHome wraps a home directory resolution failure.
//
// Parameters:
//   - cause: the underlying resolution error
//
// Returns:
//   - error: "resolving home directory: <cause>"
func DiscoverResolveHome(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryDiscoverResolveHome), cause,
	)
}

// NoDiscovery returns an error when no auto memory file exists.
//
// Parameters:
//   - path: the path that was checked
//
// Returns:
//   - error: "no auto memory at <path>"
func NoDiscovery(path string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryDiscoverNoMemory), path,
	)
}

// ReadSource wraps a source file read failure during sync.
//
// Parameters:
//   - cause: the underlying read error
//
// Returns:
//   - error: "reading source: <cause>"
func ReadSource(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryReadSource), cause,
	)
}

// ArchivePrevious wraps a failure to archive the previous mirror.
//
// Parameters:
//   - cause: the underlying archive error
//
// Returns:
//   - error: "archiving previous mirror: <cause>"
func ArchivePrevious(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryArchivePrevious), cause,
	)
}

// CreateDir wraps a failure to create the memory directory.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "creating memory directory: <cause>"
func CreateDir(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryCreateDir), cause,
	)
}

// WriteMirror wraps a failure to write the mirror file.
//
// Parameters:
//   - cause: the underlying write error
//
// Returns:
//   - error: "writing mirror: <cause>"
func WriteMirror(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryWriteMirror), cause,
	)
}

// ReadMirrorArchive wraps a failure to read the mirror for archiving.
//
// Parameters:
//   - cause: the underlying read error
//
// Returns:
//   - error: "reading mirror for archive: <cause>"
func ReadMirrorArchive(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryReadMirrorArchive), cause,
	)
}

// CreateArchiveDir wraps a failure to create the archive directory.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "creating archive directory: <cause>"
func CreateArchiveDir(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryCreateArchiveDir), cause,
	)
}

// WriteArchive wraps a failure to write an archive file.
//
// Parameters:
//   - cause: the underlying write error
//
// Returns:
//   - error: "writing archive: <cause>"
func WriteArchive(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryWriteArchive), cause,
	)
}

// ReadMirror wraps a failure to read the mirror file.
//
// Parameters:
//   - cause: the underlying read error
//
// Returns:
//   - error: "reading mirror: <cause>"
func ReadMirror(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryReadMirror), cause,
	)
}

// ReadDiffSource wraps a failure to read the source for diff.
//
// Parameters:
//   - cause: the underlying read error
//
// Returns:
//   - error: "reading diff source: <cause>"
func ReadDiffSource(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryReadDiffSource), cause,
	)
}

// SelectContent wraps a failure to select publish content.
//
// Parameters:
//   - cause: the underlying selection error
//
// Returns:
//   - error: "selecting content: <cause>"
func SelectContent(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemorySelectContent), cause,
	)
}

// Write wraps a failure to write MEMORY.md.
//
// Parameters:
//   - cause: the underlying write error
//
// Returns:
//   - error: "writing MEMORY.md: <cause>"
func WriteMemory(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrMemoryWriteMemory), cause,
	)
}
