//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package state

import (
	"time"

	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
)

// today returns today's date as YYYY-MM-DD.
//
// Returns:
//   - string: current date formatted per cfgTime.DateFormat
func today() string {
	return time.Now().Format(cfgTime.DateFormat)
}
