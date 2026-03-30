//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/ActiveMemory/ctx/internal/config/token"
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

// NewWriter creates a Writer wrapping the given output stream.
//
// Parameters:
//   - w: output stream to write to
//
// Returns:
//   - *Writer: thread-safe JSON writer
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// WriteJSON marshals v as JSON and writes it followed by a newline.
// Concurrent calls are serialized by the internal mutex.
//
// Parameters:
//   - v: value to marshal and write
//
// Returns:
//   - error: non-nil on marshal or write failure
func (sw *Writer) WriteJSON(v any) error {
	data, marshalErr := json.Marshal(v)
	if marshalErr != nil {
		return marshalErr
	}
	nl := token.NewlineLF[0]
	sw.mu.Lock()
	_, writeErr := sw.w.Write(append(data, nl))
	sw.mu.Unlock()
	return writeErr
}
