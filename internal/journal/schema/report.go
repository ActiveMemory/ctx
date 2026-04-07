//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"fmt"
	"strings"

	cfgSchema "github.com/ActiveMemory/ctx/internal/config/schema"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Report formats accumulated findings as a Markdown drift
// report. The report is deterministic and diffable.
//
// Parameters:
//   - c: Collector with accumulated findings
//
// Returns:
//   - string: Markdown-formatted drift report
func Report(c *Collector) string {
	if !c.Drift() {
		return ""
	}

	var b strings.Builder

	writeHeading(&b, c)
	writeFindings(&b, c)
	suggestAdoption(&b, c)

	return b.String()
}

// Summary returns a short summary suitable for stderr.
//
// Parameters:
//   - c: Collector with accumulated findings
//
// Returns:
//   - string: short summary for terminal display
func Summary(c *Collector) string {
	if !c.Drift() {
		return ""
	}

	findings := c.SortedFindings()

	var parts []string
	appendSummaryLine(
		&parts, findings, UnknownField,
		cfgSchema.FmtUnknownFields,
	)
	appendSummaryLine(
		&parts, findings, MissingField,
		cfgSchema.FmtMissingExpected,
	)
	appendSummaryLine(
		&parts, findings, UnknownRecordType,
		cfgSchema.FmtUnknownRecords,
	)
	appendSummaryLine(
		&parts, findings, UnknownBlockType,
		cfgSchema.FmtUnknownBlocks,
	)

	fileCount := uniqueFileCount(findings)
	indent := token.NewlineLF + token.Space + token.Space
	return fmt.Sprintf(
		cfgSchema.FmtDriftDetected, fileCount) +
		indent +
		strings.Join(parts, indent) +
		token.NewlineLF +
		cfgSchema.FmtCheckHint
}
