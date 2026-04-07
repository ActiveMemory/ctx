//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"bufio"

	"github.com/ActiveMemory/ctx/internal/config/parser"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// NewCollector creates a Collector for the given schema version.
//
// Parameters:
//   - schemaVersion: version string for scan metadata
//
// Returns:
//   - *Collector: initialized collector ready for findings
func NewCollector(schemaVersion string) *Collector {
	return &Collector{
		Findings: make(map[string]*Finding),
		Meta:     ScanMeta{SchemaVersion: schemaVersion},
	}
}

// Drift reports whether any drift findings were collected.
//
// Returns:
//   - bool: true if any findings exist
func (c *Collector) Drift() bool {
	return len(c.Findings) > 0
}

// SortedFindings returns findings sorted by type then name.
//
// Returns:
//   - []Finding: sorted copy of all findings
func (c *Collector) SortedFindings() []Finding {
	result := make([]Finding, 0, len(c.Findings))
	for _, f := range c.Findings {
		result = append(result, *f)
	}

	for i := range result {
		for j := i + 1; j < len(result); j++ {
			swap := result[i].Type > result[j].Type ||
				(result[i].Type == result[j].Type &&
					result[i].Name > result[j].Name)
			if swap {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// ValidateFile reads a JSONL file and validates each line
// against the schema.
//
// Parameters:
//   - s: Schema to validate against
//   - path: Path to the JSONL file
//   - c: Collector to accumulate findings
//
// Returns:
//   - error: Non-nil if the file cannot be opened or read
func ValidateFile(
	s *Schema, path string, c *Collector,
) error {
	f, openErr := ctxIo.SafeOpenUserFile(path)
	if openErr != nil {
		return openErr
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			ctxLog.Warn(warn.Close, path, closeErr)
		}
	}()

	c.Meta.FilesScanned++

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, parser.BufInitSize)
	scanner.Buffer(buf, parser.BufMaxSizeSchema)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		c.Meta.LinesScanned++
		validateLine(s, line, path, c)
	}

	return scanner.Err()
}
