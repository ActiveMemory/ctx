//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package log

import (
	"fmt"
	"io"
	"os"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Sink receives warning messages from best-effort operations
// whose errors would otherwise be silently discarded. Production
// code writes to os.Stderr; tests replace it with io.Discard.
var Sink io.Writer = os.Stderr

// Warn formats and writes a warning to Sink. It is intended
// for errors that are not actionable by the caller but should
// not be silently swallowed (file close, remove, state writes).
//
// The output is prefixed with "ctx: " and terminated with a
// newline. Sink write failures are silently dropped — there is
// nowhere else to report them.
func Warn(format string, args ...any) {
	_, _ = fmt.Fprintf(
		Sink, "ctx: "+format+token.NewlineLF, args...)
}
