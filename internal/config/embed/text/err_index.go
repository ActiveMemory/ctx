//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for index-block guard errors.
const (
	// DescKeyErrIndexEntriesInBlock is the text key for the error raised when
	// entry bodies are found between the INDEX:START/INDEX:END markers.
	DescKeyErrIndexEntriesInBlock = "err.index.entries-in-block"
	// DescKeyErrIndexMalformedMarkers is the text key for the error raised when
	// the INDEX:START/INDEX:END markers are missing, duplicated, or out of order.
	DescKeyErrIndexMalformedMarkers = "err.index.malformed-markers"
)
