//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"io"
	"sync"
)

// Writer serializes concurrent JSON writes to an underlying io.Writer.
//
// Fields:
//   - w: output stream
//   - mu: mutex guarding writes
type Writer struct {
	w  io.Writer
	mu sync.Mutex
}
