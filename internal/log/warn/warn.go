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

// SetSink redirects warning output to w and returns a function that
// restores the sink in effect before the call. It exists so tests can
// capture or discard warnings; production code never calls it. Callers
// must not invoke it from parallel tests — sink is process-global.
//
// Parameters:
//   - w: writer to receive subsequent warnings
//
// Returns:
//   - func(): restores the previous sink
func SetSink(w io.Writer) func() {
	prev := sink
	sink = w
	return func() { sink = prev }
}
