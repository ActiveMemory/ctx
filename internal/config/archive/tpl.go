//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import "github.com/ActiveMemory/ctx/internal/config/file"

const (
	// TplArchiveFilename is the format for dated archive filenames.
	// Args: prefix, date.
	TplArchiveFilename = "%s-%s" + file.ExtMarkdown
	// ArchiveDateSep is the separator between heading and date in archive headers.
	ArchiveDateSep = " - "
)
