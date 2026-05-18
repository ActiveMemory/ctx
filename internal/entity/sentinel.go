//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

import "github.com/ActiveMemory/ctx/internal/assets/read/desc"

// Sentinel is a typed string carrying a `desc.Text` lookup
// key. Every `internal/err/<area>/` package declares its
// identity sentinels as untyped consts of this type, one per
// logical error class.
//
// Two Sentinel values are equal under `==` when they hold
// the same key, which is exactly what `errors.Is` needs:
// as long as each sentinel uses a distinct key, callers can
// match the wrapped error with `errors.Is(err, ErrX)`.
//
// User-facing text never leaks into Go source code:
// `Error()` resolves the key against the embedded YAML
// lookup at call time, so the message is localizable and
// the sentinel value itself remains pure identity. This
// also sidesteps the package-init ordering problem of
// `var ErrX = errors.New(desc.Text(key))`, where the lookup
// table is not yet populated when the var initializer runs.
type Sentinel string

// Error implements the error interface.
//
// Returns:
//   - string: localized sentinel text resolved via
//     `desc.Text` at call time.
func (s Sentinel) Error() string {
	return desc.Text(string(s))
}
