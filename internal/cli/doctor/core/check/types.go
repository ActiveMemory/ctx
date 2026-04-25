//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check

// Result represents a single check outcome.
//
// Fields:
//   - Name: Machine-readable identifier for the check
//   - Category: Grouping label (Structure, Quality, etc.)
//   - Status: One of stats.StatusOK, stats.StatusWarning,
//     stats.StatusError, stats.StatusInfo
//   - Message: Human-readable description of the outcome
type Result struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Status   string `json:"status"`
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

// Entry pairs a check function with the name/category to attribute
// a failure to. The runner uses an ordered slice of Entry values to
// produce a uniform "did not run" line when a check returns an
// error, instead of every check having to emit its own failure
// Result for the same cause.
//
// Fields:
//   - Name: Machine-readable identifier to attribute failures to
//   - Category: Grouping label (Structure, Quality, etc.)
//   - Fn: The check function itself
type Entry struct {
	Name     string
	Category string
	Fn       func(*Report) error
}
