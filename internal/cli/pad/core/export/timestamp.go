//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"strconv"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// tsWithLabel returns a timestamp-prefixed label for collision avoidance.
//
// Parameters:
//   - label: Suffix to append after the Unix timestamp
//
// Returns:
//   - string: Label in the form "<unix_epoch>-<label>"
func tsWithLabel(label string) string {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	return ts + token.Dash + label
}
