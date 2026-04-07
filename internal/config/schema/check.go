//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

// JSONL field keys used during validation.
const (
	// FieldJSON is the finding name for malformed JSON.
	FieldJSON = "json"
	// FieldType is the JSONL top-level type key.
	FieldType = "type"
	// FieldMessage is the JSONL message envelope key.
	FieldMessage = "message"
	// FieldContent is the message content key.
	FieldContent = "content"
	// ArrayPrefix is the JSON array opening bracket.
	ArrayPrefix = "["
)

// Finding key prefixes for deduplication.
const (
	// PrefixField is the finding key prefix for unknown fields.
	PrefixField = "field"
	// PrefixMissing is the finding key prefix for missing fields.
	PrefixMissing = "missing"
	// PrefixRecord is the finding key prefix for unknown records.
	PrefixRecord = "record"
	// PrefixBlock is the finding key prefix for unknown blocks.
	PrefixBlock = "block"
	// PrefixMalformed is the finding key prefix for bad JSON.
	PrefixMalformed = "malformed"
	// PrefixUnknown is the fallback finding key prefix.
	PrefixUnknown = "unknown"
)

// Write output strings for the schema commands.
const (
	// MsgNoDirs is printed when no session directories exist.
	MsgNoDirs = "No session directories found."
	// MsgNoFiles is printed when no JSONL files are found.
	MsgNoFiles = "No session files found."
	// FmtClean is the no-drift summary format.
	FmtClean = "No schema drift. %d files, %d lines scanned."
	// FmtDumpVersion is the schema version display format.
	FmtDumpVersion = "Schema version: %s"
	// FmtDumpCCRange is the CC version range display format.
	FmtDumpCCRange = "CC version range: %s"
	// FmtDumpMetadata is the metadata record display format.
	FmtDumpMetadata = "- %s (metadata, no field validation)"
	// FmtDumpRecordType is the record type heading format.
	FmtDumpRecordType = "- %s"
	// FmtDumpRequired is the required fields format.
	FmtDumpRequired = "  Required: %s"
	// FmtDumpOptional is the optional fields format.
	FmtDumpOptional = "  Optional: %s"
	// FmtDumpBlock is the block type display format.
	FmtDumpBlock = "- %s (%s)"
	// LabelParsed is the block kind label for parsed types.
	LabelParsed = "parsed"
	// LabelKnown is the block kind label for unparsed types.
	LabelKnown = "known (not parsed)"
	// HeadingRecordTypes is the record types section heading.
	HeadingRecordTypes = "## Record Types"
	// HeadingBlockTypes is the block types section heading.
	HeadingBlockTypes = "## Content Block Types"
)

// Error message constants.
const (
	// ErrMsgDrift is the error message for schema drift.
	ErrMsgDrift = "schema drift detected"
)
