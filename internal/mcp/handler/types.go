//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

// EntryOpts holds optional fields for entry creation.
type EntryOpts struct {
	Priority    string
	Context     string
	Rationale   string
	Consequence string
	Lesson      string
	Application string
}
