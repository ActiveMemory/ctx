//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// WriteJSON marshals v as JSON and writes it to w, followed by a
// newline. The mutex serialises concurrent writes.
//
// Parameters:
//   - w: output stream
//   - mu: mutex guarding w
//   - v: value to marshal and write
//
// Returns:
//   - error: non-nil on marshal or write failure
func WriteJSON(w io.Writer, mu *sync.Mutex, v any) error {
	data, marshalErr := json.Marshal(v)
	if marshalErr != nil {
		return marshalErr
	}
	mu.Lock()
	nl := token.NewlineLF[0]
	_, writeErr := w.Write(append(data, nl))
	mu.Unlock()
	return writeErr
}
