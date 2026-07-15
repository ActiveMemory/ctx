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
//
// Fields:
//   - Header: Table header written when the artifact is created
//   - NextID: Allocates the next monotonic ID from existing bytes
//   - Render: Formats one markdown table row (with trailing newline)
//   - ErrMkdir: Wraps a parent-directory mkdir failure
//   - ErrRead: Wraps an artifact read failure
//   - ErrOpen: Wraps an open-for-append failure
//   - ErrWrite: Wraps a row-write failure
type KBRowHooks struct {
	Header   string
	NextID   func(existing []byte) (string, error)
	Render   func(id string) string
	ErrMkdir func(error) error
	ErrRead  func(error) error
	ErrOpen  func(error) error
	ErrWrite func(error) error
}
