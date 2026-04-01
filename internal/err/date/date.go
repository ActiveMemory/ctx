//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package date

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// InvalidValue returns an error for an invalid date string.
//
// Parameters:
//   - value: the invalid date string.
//
// Returns:
//   - error: "invalid date <value> (expected YYYY-MM-DD)"
func InvalidValue(value string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDateInvalidDateValue), value,
	)
}

// Invalid returns an error for an invalid date flag value.
//
// Parameters:
//   - flag: the flag name (e.g. "--since", "--until").
//   - value: the invalid date string.
//   - cause: the underlying parse error.
//
// Returns:
//   - error: "invalid <flag> date <value> (expected YYYY-MM-DD): <cause>"
func Invalid(flag, value string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrDateInvalidDate), flag, value, cause,
	)
}
