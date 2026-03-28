//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package test

// Result holds the outcome of a test notification.
type Result struct {
	NoWebhook  bool
	Filtered   bool
	StatusCode int
}
