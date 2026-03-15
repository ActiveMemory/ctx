//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// InvalidDateValue returns an error for an invalid date string.
//
// Parameters:
//   - value: the invalid date string.
//
// Returns:
//   - error: "invalid date <value> (expected YYYY-MM-DD)"
func InvalidDateValue(value string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrDateInvalidDateValue), value,
	)
}

// InvalidDate returns an error for an invalid date flag value.
//
// Parameters:
//   - flag: the flag name (e.g. "--since", "--until").
//   - value: the invalid date string.
//   - cause: the underlying parse error.
//
// Returns:
//   - error: "invalid <flag> date <value> (expected YYYY-MM-DD): <cause>"
func InvalidDate(flag, value string, cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrDateInvalidDate), flag, value, cause,
	)
}
