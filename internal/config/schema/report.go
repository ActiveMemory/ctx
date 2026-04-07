//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

// Report heading and metadata format strings.
const (
	// FmtReportHeading is the report title.
	FmtReportHeading = "# Schema Drift Report"
	// FmtSchemaVersion is the schema version line.
	FmtSchemaVersion = "Schema version: %s"
	// FmtScanStats is the scan statistics line.
	FmtScanStats = "Files scanned: %d | Lines scanned: %d"
	// FmtMalformed is the malformed line count suffix.
	FmtMalformed = " | Malformed: %d"
)

// Section titles for finding groups.
const (
	// TitleUnknownFields is the heading for unknown fields.
	TitleUnknownFields = "Unknown Fields"
	// TitleMissingFields is the heading for missing fields.
	TitleMissingFields = "Missing Required Fields"
	// TitleUnknownRecordTypes is the heading for unknown record types.
	TitleUnknownRecordTypes = "Unknown Record Types"
	// TitleUnknownBlockTypes is the heading for unknown block types.
	TitleUnknownBlockTypes = "Unknown Block Types"
	// TitleMalformedLines is the heading for malformed lines.
	TitleMalformedLines = "Malformed Lines"
)

// Markdown table format strings.
const (
	// FmtSectionHeading is a Markdown H2 heading.
	FmtSectionHeading = "## %s"
	// TableHeader is the table header row.
	TableHeader = "| Name | Occurrences | Files |"
	// TableSeparator is the table separator row.
	TableSeparator = "|------|-------------|-------|"
	// FmtTableRow is a table data row.
	FmtTableRow = "| `%s` | %d | %d |"
)

// Summary format strings for terminal output.
const (
	// FmtUnknownFields is the summary line for unknown fields.
	FmtUnknownFields = "Unknown fields: %s"
	// FmtMissingExpected is the summary line for missing fields.
	FmtMissingExpected = "Missing expected: %s"
	// FmtUnknownRecords is the summary line for unknown record types.
	FmtUnknownRecords = "Unknown record types: %s"
	// FmtUnknownBlocks is the summary line for unknown block types.
	FmtUnknownBlocks = "Unknown block types: %s"
	// FmtDriftDetected is the summary header.
	FmtDriftDetected = "Schema drift detected in %d file(s):"
	// FmtCheckHint is the trailing hint to run the check command.
	FmtCheckHint = "Run `ctx journal schema check` for full report."
)

// Suggestion format strings.
const (
	// FmtSuggestAdoption is the suggestion for universal fields.
	FmtSuggestAdoption = "**Suggestion:** All files contain field(s) `%s`."
	// SuggestAdd is the suggestion follow-up line.
	SuggestAdd = "Consider adding to the schema."
)

// Schema version and CC version range.
const (
	// Version is the current schema version.
	Version = "1.0.0"
	// CCVersionRange is the CC version range tested.
	CCVersionRange = "2.1.2\u20132.1.92"
)
