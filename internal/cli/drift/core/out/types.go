//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package out

import (
	cfgDrift "github.com/ActiveMemory/ctx/internal/config/drift"
	"github.com/ActiveMemory/ctx/internal/drift"
)

// JSONOutput represents the JSON structure for
// machine-readable drift output.
//
// Fields:
//   - Timestamp: RFC3339-formatted UTC time
//   - Status: Overall drift status
//   - Warnings: Issues that should be addressed
//   - Violations: Constitution violations
//   - Passed: Names of checks that passed
type JSONOutput struct {
	Timestamp  string               `json:"timestamp"`
	Status     cfgDrift.StatusType  `json:"status"`
	Warnings   []drift.Issue        `json:"warnings"`
	Violations []drift.Issue        `json:"violations"`
	Passed     []cfgDrift.CheckName `json:"passed"`
}
