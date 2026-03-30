//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import "github.com/ActiveMemory/ctx/internal/config/file"

// Task archive/snapshot constants.
const (
	// ScopeTasks is the scope identifier for task archives.
	ScopeTasks = "tasks"
	// DefaultSnapshotName is the default name when no snapshot name is provided.
	DefaultSnapshotName = "snapshot"
	// SnapshotFilenameFormat is the filename template for task snapshots.
	// Args: name, formatted timestamp.
	SnapshotFilenameFormat = "tasks-%s-%s" + file.ExtMarkdown
	// SnapshotTimeFormat is the compact timestamp layout for snapshot filenames.
	SnapshotTimeFormat = "2006-01-02-1504"
)

// Backup archive writer identifiers for error reporting.
const (
	// WriterGzip identifies the gzip compression writer.
	WriterGzip = "gzip"
	// WriterTar identifies the tar archive writer.
	WriterTar = "tar"
)
