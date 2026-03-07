//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import "github.com/ActiveMemory/ctx/internal/config"

// Status constants — aliased from config for local use.
const (
	StatusOK      = config.StatusOK
	StatusWarning = config.StatusWarning
	StatusError   = config.StatusError
	StatusInfo    = config.StatusInfo
)

// Result represents a single check outcome.
//
// Fields:
//   - Name: Machine-readable identifier for the check
//   - Category: Grouping label (Structure, Quality, Plugin, etc.)
//   - Status: One of StatusOK, StatusWarning, StatusError, StatusInfo
//   - Message: Human-readable description of the outcome
type Result struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Status   string `json:"status"` // "ok", "warning", "error", "info"
	Message  string `json:"message"`
}

// Report is the complete doctor output.
//
// Fields:
//   - Results: All individual check results
//   - Warnings: Count of results with StatusWarning
//   - Errors: Count of results with StatusError
type Report struct {
	Results  []Result `json:"results"`
	Warnings int      `json:"warnings"`
	Errors   int      `json:"errors"`
}
