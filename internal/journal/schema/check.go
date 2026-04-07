//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"encoding/json"
	"strings"

	cfgSchema "github.com/ActiveMemory/ctx/internal/config/schema"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// add records a finding, merging with an existing one if the
// same type+name combination was already seen.
//
// Parameters:
//   - ft: finding classification
//   - name: field, type, or block name that drifted
//   - filePath: source file where the finding occurred
func (c *Collector) add(
	ft FindingType, name, filePath string,
) {
	key := findingKey(ft, name)
	f, ok := c.Findings[key]
	if !ok {
		f = &Finding{
			Type:  ft,
			Name:  name,
			Files: make(map[string]bool),
		}
		c.Findings[key] = f
	}
	f.Files[filePath] = true
	f.Count++
}

// validateLine checks a single JSONL line against the schema,
// recording any drift into the collector.
//
// Parameters:
//   - s: schema to validate against
//   - line: raw JSONL bytes for one line
//   - filePath: source file path for finding attribution
//   - c: collector to accumulate findings into
func validateLine(
	s *Schema, line []byte, filePath string, c *Collector,
) {
	var raw map[string]json.RawMessage
	unmarshalErr := json.Unmarshal(line, &raw)
	if unmarshalErr != nil {
		c.Meta.MalformedLines++
		c.add(MalformedLine, cfgSchema.FieldJSON, filePath)
		return
	}

	recordType := extractString(raw[cfgSchema.FieldType])
	if recordType == "" {
		return
	}

	if !s.KnownRecordType(recordType) {
		c.add(UnknownRecordType, recordType, filePath)
		return
	}

	rs := s.RecordTypes[recordType]

	if len(rs.Required) == 0 && len(rs.Optional) == 0 {
		return
	}

	for field := range raw {
		if !s.KnownField(recordType, field) {
			c.add(UnknownField, field, filePath)
		}
	}

	for _, field := range rs.Required {
		if _, ok := raw[field]; !ok {
			c.add(MissingField, field, filePath)
		}
	}

	validateBlocks(s, raw, filePath, c)
}

// validateBlocks inspects the message.content array for
// unknown content block types.
//
// Parameters:
//   - s: schema with known block types
//   - raw: top-level JSONL fields as raw JSON
//   - filePath: source file for finding attribution
//   - c: collector to accumulate findings into
func validateBlocks(
	s *Schema, raw map[string]json.RawMessage,
	filePath string, c *Collector,
) {
	msgRaw, ok := raw[cfgSchema.FieldMessage]
	if !ok {
		return
	}

	var msg map[string]json.RawMessage
	unmarshalErr := json.Unmarshal(msgRaw, &msg)
	if unmarshalErr != nil {
		return
	}

	contentRaw, ok := msg[cfgSchema.FieldContent]
	if !ok {
		return
	}

	trimmed := strings.TrimSpace(string(contentRaw))
	if !strings.HasPrefix(trimmed, cfgSchema.ArrayPrefix) {
		return
	}

	var blocks []map[string]json.RawMessage
	blockErr := json.Unmarshal(contentRaw, &blocks)
	if blockErr != nil {
		return
	}

	for _, block := range blocks {
		bt := extractString(block[cfgSchema.FieldType])
		if bt != "" && !s.KnownBlockType(bt) {
			c.add(UnknownBlockType, bt, filePath)
		}
	}
}

// extractString extracts a Go string from a raw JSON value.
// Returns empty string if raw is nil or not a JSON string.
//
// Parameters:
//   - raw: JSON-encoded value to decode
//
// Returns:
//   - string: decoded string, or empty on failure
func extractString(raw json.RawMessage) string {
	if raw == nil {
		return ""
	}
	var s string
	if unmarshalErr := json.Unmarshal(raw, &s); unmarshalErr != nil {
		return ""
	}
	return s
}

// findingKey produces a deduplication key for a finding by
// combining a type-based prefix with the finding name.
//
// Parameters:
//   - ft: finding type for prefix selection
//   - name: field, type, or block name
//
// Returns:
//   - string: key in the form "prefix:name"
func findingKey(ft FindingType, name string) string {
	prefixes := [...]string{
		UnknownField:      cfgSchema.PrefixField,
		MissingField:      cfgSchema.PrefixMissing,
		UnknownRecordType: cfgSchema.PrefixRecord,
		UnknownBlockType:  cfgSchema.PrefixBlock,
		MalformedLine:     cfgSchema.PrefixMalformed,
	}
	prefix := cfgSchema.PrefixUnknown
	if int(ft) < len(prefixes) {
		prefix = prefixes[ft]
	}
	return prefix + token.Colon + name
}
