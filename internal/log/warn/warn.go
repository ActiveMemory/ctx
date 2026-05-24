//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package warn

import (
	"fmt"
	"io"
	"os"

	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// sink receives warning messages from best-effort operations
// whose errors would otherwise be silently discarded. Production
// code writes to os.Stderr; tests replace it with io.Discard.
var sink io.Writer = os.Stderr

// Warn formats and writes a warning to sink. It is intended
// for errors that are not actionable by the caller but should
// not be silently swallowed (file close, remove, state writes).
//
// The output is prefixed with "ctx: " and terminated with a
// newline. Sink write failures are silently dropped; there is
// nowhere else to report them.
//
// Parameters:
//   - format: Printf-style format string
//   - args: Format arguments
func Warn(format string, args ...any) {
	_, _ = fmt.Fprintf(
		sink, cfgCtx.StderrPrefix+format+token.NewlineLF, args...)
}

// SetSinkForTesting swaps the warn sink for the duration of a
// test and returns a restore function the caller must defer.
// The package-level sink is otherwise unexported because there
// is no production reason to redirect it; tests that assert
// against captured warning output go through this helper.
//
// Parameters:
//   - w: writer to receive warnings during the test
//
// Returns:
//   - func(): restores the previous sink when called
func SetSinkForTesting(w io.Writer) func() {
	prev := sink
	sink = w
	return func() { sink = prev }
}
