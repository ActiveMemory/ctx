//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package contradiction

import (
	"fmt"
	"strconv"

	cfgKbC "github.com/ActiveMemory/ctx/internal/config/kb/contradiction"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	errKbC "github.com/ActiveMemory/ctx/internal/err/kb/contradiction"
)

// nextID scans existing for `C-###` tokens and returns the
// next formatted ID. Empty / no-matches input yields `C-001`.
//
// Parameters:
//   - existing: prior file contents (may be nil).
//
// Returns:
//   - string: the next zero-padded ID.
//   - error: wrapped [errKbC.ParseCNumber] on overflow.
func nextID(existing []byte) (string, error) {
	high := 0
	for _, m := range regex.KBContradictionID.FindAllSubmatch(
		existing, -1,
	) {
		digits := string(m[1])
		n, parseErr := strconv.Atoi(digits)
		if parseErr != nil {
			return "", errKbC.ParseCNumber(digits, parseErr)
		}
		if n > high {
			high = n
		}
	}
	return fmt.Sprintf(cfgKbC.IDFormat, cfgKbC.IDPrefix, high+1), nil
}
