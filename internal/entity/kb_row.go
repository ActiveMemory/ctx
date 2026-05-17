//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// KBRowHooks supplies the per-artifact pieces of an append
// against a monotonic-ID kb tabular artifact (contradictions,
// domain-decisions, outstanding-questions): the table header,
// the next-ID allocator (scans existing bytes), the row
// renderer (formats one markdown table row including trailing
// newline), and the four error constructors that wrap I/O
// failures in the caller's domain-typed errors.
//
// Consumed by [github.com/ActiveMemory/ctx/internal/write/kb/row.Append].
type KBRowHooks struct {
	Header   string
	NextID   func(existing []byte) (string, error)
	Render   func(id string) string
	ErrMkdir func(error) error
	ErrRead  func(error) error
	ErrOpen  func(error) error
	ErrWrite func(error) error
}
