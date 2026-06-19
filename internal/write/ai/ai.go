//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ai

import (
	"fmt"
	"io"

	cfgAI "github.com/ActiveMemory/ctx/internal/config/ai"
)

// Ping writes backend ping information.
//
// Parameters:
//   - out: destination writer
//   - backend: backend name
//   - endpoint: configured endpoint
//   - firstModel: first model from model listing
//
// Returns:
//   - error: write failure
func Ping(
	out io.Writer,
	backend string,
	endpoint string,
	firstModel string,
) error {
	_, writeErr := fmt.Fprintf(
		out,
		cfgAI.WritePingFormat,
		backend,
		endpoint,
		firstModel,
	)
	return writeErr
}
