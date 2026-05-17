//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package evidence

import (
	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	errKbEvidence "github.com/ActiveMemory/ctx/internal/err/kb/evidence"
)

// validateBand rejects unknown confidence bands.
//
// Parameters:
//   - band: confidence band string.
//
// Returns:
//   - error: wrapped [errKbEvidence.ErrInvalidBand] or nil.
func validateBand(band string) error {
	switch band {
	case cfgKB.ConfidenceHigh,
		cfgKB.ConfidenceMedium,
		cfgKB.ConfidenceLow,
		cfgKB.ConfidenceSpeculative:
		return nil
	}
	return errKbEvidence.InvalidBand(band)
}
