//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"time"

	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// parseTimeout converts a duration string to
// [time.Duration]. Empty input yields zero. Unparseable
// input emits a warning naming the backend and returns
// zero so the backend's own default applies.
//
// Parameters:
//   - name: backend name, used to make the warning
//     identifiable when multiple backends are configured.
//   - raw: duration string from `.ctxrc`, e.g., "30s".
//
// Returns:
//   - time.Duration: parsed value, or zero on empty or
//     unparseable input.
func parseTimeout(name, raw string) time.Duration {
	if raw == "" {
		return 0
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		logWarn.Warn(cfgWarn.BackendInvalidTimeout, name, raw, err)
		return 0
	}
	return d
}
