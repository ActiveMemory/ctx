//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package decision

import (
	"fmt"
	"strconv"

	cfgKbDD "github.com/ActiveMemory/ctx/internal/config/kb/decision"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	errKbDD "github.com/ActiveMemory/ctx/internal/err/kb/decision"
)

// nextID scans existing for `DD-###` tokens and returns the
// next formatted ID. Empty / no-matches input yields `DD-001`.
//
// Parameters:
//   - existing: prior file contents (may be nil).
//
// Returns:
//   - string: the next zero-padded ID.
//   - error: wrapped [errKbDD.ParseDDNumber] on overflow.
func nextID(existing []byte) (string, error) {
	high := 0
	for _, m := range regex.KBDecisionID.FindAllSubmatch(
		existing, -1,
	) {
		digits := string(m[1])
		n, parseErr := strconv.Atoi(digits)
		if parseErr != nil {
			return "", errKbDD.ParseDDNumber(digits, parseErr)
		}
		if n > high {
			high = n
		}
	}
	return fmt.Sprintf(
		cfgKbDD.IDFormat, cfgKbDD.IDPrefix, high+1,
	), nil
}
