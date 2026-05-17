//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package evidence

import (
	"fmt"
	"strconv"
	"strings"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	cfgKbEvidence "github.com/ActiveMemory/ctx/internal/config/kb/evidence"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	errKbEvidence "github.com/ActiveMemory/ctx/internal/err/kb/evidence"
)

// maxIDFrom scans raw for `EV-NNN` tokens and returns the
// highest numeric value found. Empty / no-match input returns
// 0.
//
// Parameters:
//   - raw: file contents.
//
// Returns:
//   - int: highest EV number, or 0.
//   - error: wrapped [errKbEvidence.ParseEVNumber] on a
//     malformed digit string.
func maxIDFrom(raw string) (int, error) {
	matches := regex.KBEvidenceID.FindAllStringSubmatch(raw, -1)
	high := 0
	for _, m := range matches {
		n, parseErr := strconv.Atoi(m[1])
		if parseErr != nil {
			return 0, errKbEvidence.ParseEVNumber(m[1], parseErr)
		}
		if n > high {
			high = n
		}
	}
	return high, nil
}

// alreadyExists reports whether the given EV-### identifier is
// already present anywhere in the raw file content.
//
// Parameters:
//   - raw: file contents.
//   - id: EV-### identifier to check.
//
// Returns:
//   - bool: true when id is already minted.
func alreadyExists(raw, id string) bool {
	if id == "" {
		return false
	}
	return strings.Contains(raw, id)
}

// formatID renders an integer as a zero-padded EV-NNN string.
//
// Parameters:
//   - n: positive integer.
//
// Returns:
//   - string: `EV-NNN` (3-digit zero-padded for n ≤ 999).
func formatID(n int) string {
	return fmt.Sprintf(cfgKbEvidence.IDFormat,
		cfgKB.EVIDPrefix, cfgKB.EVIDDigits, n)
}
