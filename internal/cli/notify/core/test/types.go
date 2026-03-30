//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package test

// Result holds the outcome of a test notification.
//
// Fields:
//   - NoWebhook: No webhook URL is configured
//   - Filtered: Event was excluded by the filter
//   - StatusCode: HTTP response status code
type Result struct {
	NoWebhook  bool
	Filtered   bool
	StatusCode int
}
