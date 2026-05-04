//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// SanEntryHeader matches entry headers like "## [2026-" in
// content sanitization (MCP-SAN.3).
var SanEntryHeader = regexp.MustCompile(
	`(?m)^##\s+\[\d{4}-`,
)

// SanTaskCheckbox matches task checkboxes "- [ ]" and
// "- [x]" in content sanitization.
var SanTaskCheckbox = regexp.MustCompile(
	`(?m)^-\s+\[[x ]\]`,
)

// SanConstitutionRule matches constitution rule format
// "- [ ] **Never" in content sanitization.
var SanConstitutionRule = regexp.MustCompile(
	`(?m)^-\s+\[[x ]\]\s+\*\*[A-Z]`,
)

// SanSessionIDUnsafe matches characters not safe for session
// IDs in file paths: anything outside [a-zA-Z0-9._-].
var SanSessionIDUnsafe = regexp.MustCompile(
	`[^a-zA-Z0-9._-]`,
)
