//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

// Export configuration.
const (
	// MaxMessagesPerPart is the maximum number of messages per exported
	// journal file. Sessions with more messages are split into multiple
	// parts for browser performance.
	MaxMessagesPerPart = 200
)
