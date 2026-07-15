//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

// JSONRef represents a resolved context reference for JSON output.
//
// Fields:
//   - Raw: The reference as written (e.g. "T-3")
//   - Type: Reference kind (task, decision, learning, ...)
//   - Number: Numeric ID parsed from the reference (0 if none)
//   - Title: Resolved entry title (empty if not found)
//   - Detail: Additional resolved detail (empty if not found)
//   - Found: True when the reference resolved to a known entry
type JSONRef struct {
	Raw    string `json:"raw"`
	Type   string `json:"type"`
	Number int    `json:"number,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Found  bool   `json:"found"`
}

// JSONCommit represents a commit with its context refs for JSON output.
//
// Fields:
//   - Commit: The commit hash
//   - Message: The commit subject line
//   - Refs: Context references found in the commit
type JSONCommit struct {
	Commit  string    `json:"commit"`
	Message string    `json:"message"`
	Refs    []JSONRef `json:"refs"`
}
