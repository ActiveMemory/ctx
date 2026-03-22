//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

const (
	// SIThreshold is the boundary between raw and abbreviated SI display (1000).
	SIThreshold = 1000

	// SIThresholdM is the boundary between K and M display (1,000,000).
	SIThresholdM = 1_000_000

	// IECUnit is the binary unit base for byte formatting (1024).
	IECUnit = 1024

	// HashPrefixLen is the number of bytes used for truncated hex hashes.
	HashPrefixLen = 8
)
