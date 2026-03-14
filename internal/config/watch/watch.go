//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

// Stream scanner buffer sizes.
const (
	// StreamScannerInitCap is the initial capacity for the scanner buffer.
	StreamScannerInitCap = 64 * 1024
	// StreamScannerMaxSize is the maximum size for the scanner buffer.
	StreamScannerMaxSize = 1024 * 1024
)
