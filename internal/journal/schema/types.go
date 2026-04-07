//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

// Schema declares the expected shape of Claude Code JSONL records.
//
// Fields:
//   - Version: Schema version identifier
//   - CCVersionRange: Claude Code version range covered
//   - RecordTypes: Expected fields per record type
//   - BlockTypes: Known content block types
type Schema struct {
	Version        string
	CCVersionRange string
	RecordTypes    map[string]RecordSchema
	BlockTypes     map[string]BlockKind
}

// RecordSchema declares expected fields for a record type.
//
// Fields:
//   - Required: Fields that must be present
//   - Optional: Fields that may appear
type RecordSchema struct {
	Required []string
	Optional []string
}

// BlockKind classifies a content block type.
type BlockKind int

const (
	// BlockParsed is a block type the parser extracts.
	BlockParsed BlockKind = iota
	// BlockKnown is a recognized but skipped block type.
	BlockKnown
)

// FindingType classifies a schema drift observation.
type FindingType int

const (
	// UnknownField is an unrecognized top-level field.
	UnknownField FindingType = iota
	// MissingField is a required field absent from a record.
	MissingField
	// UnknownRecordType is an unrecognized record type.
	UnknownRecordType
	// UnknownBlockType is an unrecognized content block type.
	UnknownBlockType
	// MalformedLine is a line that failed JSON parsing.
	MalformedLine
)

// Finding represents a single schema drift observation.
//
// Fields:
//   - Type: Classification of the drift
//   - Name: The field, type, or block name that drifted
//   - Files: Which files exhibited this finding
//   - Count: How many lines exhibited this finding
type Finding struct {
	Type  FindingType
	Name  string
	Files map[string]bool
	Count int
}

// ScanMeta holds metadata about a validation scan run.
//
// Fields:
//   - FilesScanned: Total JSONL files examined
//   - LinesScanned: Total lines examined
//   - MalformedLines: Lines that failed JSON parsing
//   - SchemaVersion: Version of the schema used
type ScanMeta struct {
	FilesScanned   int
	LinesScanned   int
	MalformedLines int
	SchemaVersion  string
}

// Collector accumulates findings across multiple files.
//
// Fields:
//   - Findings: Deduplicated findings keyed by "type:name"
//   - Meta: Scan metadata
type Collector struct {
	Findings map[string]*Finding
	Meta     ScanMeta
}
