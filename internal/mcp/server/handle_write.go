//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// writeJSON marshals v as JSON and writes it to the output stream,
// followed by a newline.
//
// Safe to call from any goroutine; outMu serialises access.
//
// Parameters:
//   - v: value to marshal and write
//
// Returns:
//   - error: non-nil on marshal or write failure
func (s *Server) writeJSON(v any) error {
	out, marshalErr := json.Marshal(v)
	if marshalErr != nil {
		return marshalErr
	}
	s.outMu.Lock()
	_, writeErr := s.out.Write(append(out, token.NewlineLF[0]))
	s.outMu.Unlock()
	return writeErr
}
