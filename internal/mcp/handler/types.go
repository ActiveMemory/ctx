//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

// violation represents a single governance violation recorded by the
// VS Code extension's detection ring.
//
// Fields:
//   - Kind: violation category identifier
//   - Detail: human-readable description of what was violated
//   - Timestamp: ISO-8601 timestamp of when the violation occurred
type violation struct {
	Kind      string `json:"kind"`
	Detail    string `json:"detail"`
	Timestamp string `json:"timestamp"`
}

// violationsData is the JSON structure of the violations file.
//
// Fields:
//   - Entries: list of recorded violations
type violationsData struct {
	Entries []violation `json:"entries"`
}
