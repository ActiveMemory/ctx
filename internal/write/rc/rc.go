//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
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
	logWarn.Warn(cfgWarn.ParseConfig, filename, cause)
}
