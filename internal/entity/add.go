//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// EntryParams contains all parameters needed to add an entry to a context file.
type EntryParams struct {
	Type        string
	Content     string
	Section     string
	Priority    string
	Context     string
	Rationale   string
	Consequence string
	Lesson      string
	Application string
	ContextDir  string
}

// AddConfig holds all flags for the add command.
type AddConfig struct {
	Priority    string
	Section     string
	FromFile    string
	Context     string
	Rationale   string
	Consequence string
	Lesson      string
	Application string
}
