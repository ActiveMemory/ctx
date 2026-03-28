//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package wrap

// Wrap-up marker configuration.
const (
	// Marker is the state file name for the wrap-up suppression marker.
	Marker = "ctx-wrapped-up"
	// Content is the content written to the wrap-up marker file.
	Content = "wrapped-up"
	// ExpiryHours is how many hours the marker suppresses nudges.
	ExpiryHours = 2
)
