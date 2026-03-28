//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

// Item represents a single blob ready for export.
type Item struct {
	Label   string
	Data    []byte
	OutPath string
	AltName string // Non-empty when collision renamed
	Exists  bool   // True when outPath already exists
}
