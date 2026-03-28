//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

// ConfigPattern pairs a glob pattern with its localizable topic description.
type ConfigPattern struct {
	Pattern string
	Topic   string
}
