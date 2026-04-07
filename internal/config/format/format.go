//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

// SI display threshold constants.
const (
	// SIThreshold is the boundary between raw and abbreviated SI display (1000).
	SIThreshold = 1000

	// SIThresholdM is the boundary between K and M display (1,000,000).
	SIThresholdM = 1_000_000

	// IECUnit is the binary unit base for byte formatting (1024).
	IECUnit = 1024

	// HashPrefixLen is the number of bytes used for truncated hex hashes.
	HashPrefixLen = 8

	// TruncateDetail is the max character width for detail strings in
	// governance violation reports and similar summaries.
	TruncateDetail = 120

	// TruncateTitle is the max character width for title/summary lines
	// in import previews and list views.
	TruncateTitle = 60

	// TruncateDescription is the max character width for description
	// text in skill listings and similar compact displays.
	TruncateDescription = 70

	// PreviewLines is the number of content lines shown in status
	// previews and similar compact displays.
	PreviewLines = 5

	// StatusPreviewLines is the number of content lines shown
	// in verbose status file previews.
	StatusPreviewLines = 3

	// StatusRecentFiles is the number of recently modified
	// files shown in the status activity section.
	StatusRecentFiles = 3
)

// IEC binary unit prefix string for byte formatting.
const IECPrefixes = "KMGTPE"
