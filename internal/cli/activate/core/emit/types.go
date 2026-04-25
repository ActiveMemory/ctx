//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package emit

// emitter produces a shell-specific one-line statement for the given
// key and pre-quoted value, terminated by a newline. Concrete
// emitters live in posix.go; the dispatch table is in emit.go.
type emitter func(key, quotedValue string) string
