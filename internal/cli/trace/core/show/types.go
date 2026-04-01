//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

// JSONRef represents a resolved context reference for JSON output.
type JSONRef struct {
	Raw    string `json:"raw"`
	Type   string `json:"type"`
	Number int    `json:"number,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Found  bool   `json:"found"`
}

// JSONCommit represents a commit with its context refs for JSON output.
type JSONCommit struct {
	Commit  string    `json:"commit"`
	Message string    `json:"message"`
	Refs    []JSONRef `json:"refs"`
}
