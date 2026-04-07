//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"fmt"
	"sort"
	"strings"

	cfgSchema "github.com/ActiveMemory/ctx/internal/config/schema"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// writeHeading writes the report title and scan metadata.
//
// Parameters:
//   - b: string builder to append to
//   - c: collector with scan metadata
func writeHeading(b *strings.Builder, c *Collector) {
	b.WriteString(cfgSchema.FmtReportHeading)
	b.WriteString(token.DoubleNewline)
	if _, fmtErr := fmt.Fprintf(b,
		cfgSchema.FmtSchemaVersion,
		c.Meta.SchemaVersion); fmtErr != nil {
		ctxLog.Warn(warn.Write, cfgSchema.FmtSchemaVersion, fmtErr)
	}
	b.WriteString(token.NewlineLF)
	if _, fmtErr := fmt.Fprintf(b,
		cfgSchema.FmtScanStats,
		c.Meta.FilesScanned,
		c.Meta.LinesScanned); fmtErr != nil {
		ctxLog.Warn(warn.Write, cfgSchema.FmtScanStats, fmtErr)
	}
	if c.Meta.MalformedLines > 0 {
		if _, fmtErr := fmt.Fprintf(b,
			cfgSchema.FmtMalformed,
			c.Meta.MalformedLines); fmtErr != nil {
			ctxLog.Warn(warn.Write, cfgSchema.FmtMalformed, fmtErr)
		}
	}
	b.WriteString(token.DoubleNewline)
}

// writeFindings writes grouped finding sections as Markdown
// tables, one section per finding type that has results.
//
// Parameters:
//   - b: string builder to append to
//   - c: collector with accumulated findings
func writeFindings(b *strings.Builder, c *Collector) {
	findings := c.SortedFindings()

	sections := []struct {
		ft    FindingType
		title string
	}{
		{UnknownField, cfgSchema.TitleUnknownFields},
		{MissingField, cfgSchema.TitleMissingFields},
		{UnknownRecordType, cfgSchema.TitleUnknownRecordTypes},
		{UnknownBlockType, cfgSchema.TitleUnknownBlockTypes},
		{MalformedLine, cfgSchema.TitleMalformedLines},
	}

	for _, sec := range sections {
		group := filterByType(findings, sec.ft)
		if len(group) == 0 {
			continue
		}
		if _, fmtErr := fmt.Fprintf(b,
			cfgSchema.FmtSectionHeading,
			sec.title); fmtErr != nil {
			ctxLog.Warn(warn.Write, cfgSchema.FmtSectionHeading, fmtErr)
		}
		b.WriteString(token.DoubleNewline)
		b.WriteString(cfgSchema.TableHeader)
		b.WriteString(token.NewlineLF)
		b.WriteString(cfgSchema.TableSeparator)
		b.WriteString(token.NewlineLF)
		for _, f := range group {
			files := sortedFiles(f.Files)
			if _, fmtErr := fmt.Fprintf(b,
				cfgSchema.FmtTableRow,
				f.Name, f.Count,
				len(files)); fmtErr != nil {
				ctxLog.Warn(warn.Write, cfgSchema.FmtTableRow, fmtErr)
			}
			b.WriteString(token.NewlineLF)
		}
		b.WriteString(token.NewlineLF)
	}
}

// suggestAdoption appends a suggestion when an unknown field
// appears in every scanned file, indicating the schema should
// be updated to include it.
//
// Parameters:
//   - b: string builder to append to
//   - c: collector with findings and file count
func suggestAdoption(b *strings.Builder, c *Collector) {
	if c.Meta.FilesScanned == 0 {
		return
	}
	findings := c.SortedFindings()
	var candidates []string
	for _, f := range findings {
		if f.Type == UnknownField &&
			len(f.Files) == c.Meta.FilesScanned {
			candidates = append(candidates, f.Name)
		}
	}
	if len(candidates) == 0 {
		return
	}
	sort.Strings(candidates)
	b.WriteString(token.Separator)
	b.WriteString(token.DoubleNewline)
	backtickJoin := token.Backtick +
		token.CommaSpace + token.Backtick
	if _, fmtErr := fmt.Fprintf(b,
		cfgSchema.FmtSuggestAdoption,
		strings.Join(candidates, backtickJoin)); fmtErr != nil {
		ctxLog.Warn(warn.Write, cfgSchema.FmtSuggestAdoption, fmtErr)
	}
	b.WriteString(token.NewlineLF)
	b.WriteString(cfgSchema.SuggestAdd)
	b.WriteString(token.NewlineLF)
}

// appendSummaryLine adds a formatted summary line for the
// given finding type if any findings of that type exist.
//
// Parameters:
//   - parts: accumulator slice to append to
//   - findings: all findings to filter
//   - ft: finding type to select
//   - format: printf format string for the line
func appendSummaryLine(
	parts *[]string, findings []Finding,
	ft FindingType, format string,
) {
	if countByType(findings, ft) > 0 {
		*parts = append(*parts,
			fmt.Sprintf(format, nameList(findings, ft)))
	}
}

// filterByType returns only findings matching the given type.
//
// Parameters:
//   - findings: all findings to filter
//   - ft: finding type to select
//
// Returns:
//   - []Finding: subset matching the type
func filterByType(
	findings []Finding, ft FindingType,
) []Finding {
	var result []Finding
	for _, f := range findings {
		if f.Type == ft {
			result = append(result, f)
		}
	}
	return result
}

// countByType returns the number of distinct finding names
// for the given type.
//
// Parameters:
//   - findings: all findings to count
//   - ft: finding type to match
//
// Returns:
//   - int: number of distinct names
func countByType(findings []Finding, ft FindingType) int {
	n := 0
	for _, f := range findings {
		if f.Type == ft {
			n++
		}
	}
	return n
}

// nameList returns a comma-separated string of finding names
// for the given type.
//
// Parameters:
//   - findings: all findings to extract names from
//   - ft: finding type to match
//
// Returns:
//   - string: comma-separated names
func nameList(findings []Finding, ft FindingType) string {
	var names []string
	for _, f := range findings {
		if f.Type == ft {
			names = append(names, f.Name)
		}
	}
	return strings.Join(names, token.CommaSpace)
}

// uniqueFileCount returns the total number of distinct file
// paths across all findings.
//
// Parameters:
//   - findings: all findings to count files from
//
// Returns:
//   - int: number of unique file paths
func uniqueFileCount(findings []Finding) int {
	files := make(map[string]bool)
	for _, f := range findings {
		for path := range f.Files {
			files[path] = true
		}
	}
	return len(files)
}

// sortedFiles returns file paths from a set, sorted
// alphabetically for deterministic output.
//
// Parameters:
//   - files: set of file paths to sort
//
// Returns:
//   - []string: sorted file paths
func sortedFiles(files map[string]bool) []string {
	result := make([]string, 0, len(files))
	for f := range files {
		result = append(result, f)
	}
	sort.Strings(result)
	return result
}
