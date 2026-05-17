//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package question

import (
	"fmt"
	"strconv"

	cfgKbQ "github.com/ActiveMemory/ctx/internal/config/kb/question"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	errKbQ "github.com/ActiveMemory/ctx/internal/err/kb/question"
)

// nextID scans existing for `Q-###` tokens and returns the
// next formatted ID. Empty / no-matches input yields `Q-001`.
//
// Parameters:
//   - existing: prior file contents (may be nil).
//
// Returns:
//   - string: the next zero-padded ID.
//   - error: wrapped [errKbQ.ParseQNumber] on overflow.
func nextID(existing []byte) (string, error) {
	high := 0
	for _, m := range regex.KBQuestionID.FindAllSubmatch(
		existing, -1,
	) {
		digits := string(m[1])
		n, parseErr := strconv.Atoi(digits)
		if parseErr != nil {
			return "", errKbQ.ParseQNumber(digits, parseErr)
		}
		if n > high {
			high = n
		}
	}
	return fmt.Sprintf(cfgKbQ.IDFormat, cfgKbQ.IDPrefix, high+1), nil
}
