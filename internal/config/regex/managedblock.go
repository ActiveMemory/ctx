//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// ManagedKBTopics matches the CTX:KB:TOPICS managed block inside
// `.context/kb/index.md`. The pattern is greedy + multi-line so
// the reindex command can replace the body wholesale.
var ManagedKBTopics = regexp.MustCompile(
	`(?s)<!-- CTX:KB:TOPICS START -->.*?<!-- CTX:KB:TOPICS END -->`,
)
