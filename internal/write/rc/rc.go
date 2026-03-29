//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"fmt"
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// ParseWarning prints a YAML parse warning to stderr.
//
// This runs during config loading before any cobra command exists,
// so it writes to os.Stderr directly.
//
// Parameters:
//   - filename: the config file that failed to parse
//   - cause: the parse error
func ParseWarning(filename string, cause error) {
	_, _ = fmt.Fprintf(os.Stderr,
		desc.Text(text.DescKeyRcParseWarning)+token.NewlineLF,
		filename, cause)
}
