//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package drift provides functionality for detecting stale or invalid context.
package drift

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

const staleAgeDays = 30

var staleAgeExclude = []string{config.FileConstitution}

// Status returns the overall status of the report.
//
// Returns:
//   - StatusType: StatusViolation if any violations, StatusWarning if only
//     warnings, StatusOk otherwise
func (r *Report) Status() StatusType {
	if len(r.Violations) > 0 {
		return StatusViolation
	}
	if len(r.Warnings) > 0 {
		return StatusWarning
	}
	return StatusOk
}

// Detect runs all drift detection checks on the given context.
//
// Performs multiple validation checks including path references, staleness
// indicators, constitution compliance, and required file presence.
//
// Parameters:
//   - ctx: Loaded context containing files to check
//
// Returns:
//   - *Report: Drift report with warnings, violations, and passed checks
func Detect(ctx *context.Context) *Report {
	report := &Report{
		Warnings:   []Issue{},
		Violations: []Issue{},
		Passed:     []CheckName{},
	}

	// Check path references in context files
	checkPathReferences(ctx, report)

	// Check for staleness indicators
	checkStaleness(ctx, report)

	// Check constitution rules (basic heuristics)
	checkConstitution(ctx, report)

	// Check for empty required files
	checkRequiredFiles(ctx, report)

	// Check for files not modified recently
	checkFileAge(ctx, report)

	// Check for excessive entry counts in knowledge files
	checkEntryCount(ctx, report)

	// Check for undocumented internal packages
	checkMissingPackages(ctx, report)

	return report
}

// checkPathReferences scans ARCHITECTURE.md and CONVENTIONS.md for dead paths.
//
// Looks for backtick-enclosed file paths and verifies they exist on disk.
// Skips URLs, template patterns, and glob patterns.
//
// Parameters:
//   - ctx: Loaded context containing files to scan
//   - report: Report to append warnings to (modified in place)
func checkPathReferences(ctx *context.Context, report *Report) {
	foundDeadPaths := false

	for _, f := range ctx.Files {
		if f.Name != config.FileArchitecture && f.Name != config.FileConvention {
			continue
		}

		lines := strings.Split(string(f.Content), config.NewlineLF)
		for lineNum, line := range lines {
			matches := config.RegExPath.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				path := m[1]
				// Skip URLs and common non-file patterns
				if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "//") {
					continue
				}
				// Skip template patterns
				if strings.Contains(path, "{") || strings.Contains(path, "*") {
					continue
				}
				// Check if the file exists
				if _, err := os.Stat(path); os.IsNotExist(err) {
					report.Warnings = append(report.Warnings, Issue{
						File:    f.Name,
						Line:    lineNum + 1,
						Type:    IssueDeadPath,
						Message: "references path that does not exist",
						Path:    path,
					})
					foundDeadPaths = true
				}
			}
		}
	}

	if !foundDeadPaths {
		report.Passed = append(report.Passed, CheckPathReferences)
	}
}

// checkStaleness detects signs that context files need maintenance.
//
// Currently checks for excessive completed tasks (>10) in TASKS.md,
// which indicates the file should be compacted.
//
// Parameters:
//   - ctx: Loaded context containing files to scan
//   - report: Report to append warnings to (modified in place)
func checkStaleness(ctx *context.Context, report *Report) {
	staleness := false

	if f := ctx.File(config.FileTask); f != nil {
		// Count completed tasks
		completedCount := strings.Count(string(f.Content), "- [x]")
		if completedCount > 10 {
			report.Warnings = append(report.Warnings, Issue{
				File:    f.Name,
				Type:    IssueStaleness,
				Message: "has many completed items (consider archiving)",
				Path:    "",
			})
			staleness = true
		}
	}

	if !staleness {
		report.Passed = append(report.Passed, CheckStaleness)
	}
}

// checkConstitution performs heuristic checks for constitution violations.
//
// Currently, it scans the working directory for files that may contain secrets
// (e.g., .env, credentials, api_key) and flags them as violations.
//
// Parameters:
//   - ctx: Loaded context (currently unused, reserved for future checks)
//   - report: Report to append violations to (modified in place)
func checkConstitution(_ *context.Context, report *Report) {
	// Basic heuristic checks for constitution violations
	// Check for potential secrets in common config files

	secretPatterns := []string{
		".env",
		"credentials",
		"secret",
		"api_key",
		"apikey",
		"password",
	}

	// Look for common secret file patterns in the working directory
	entries, readErr := os.ReadDir(".")
	if readErr != nil {
		return
	}

	foundViolation := false
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		for _, pattern := range secretPatterns {
			if strings.Contains(name, pattern) &&
				!strings.HasSuffix(name, ".example") &&
				!strings.HasSuffix(name, ".sample") {
				// Check if it contains actual content (not just template)
				content, readFileErr := os.ReadFile(entry.Name())
				if readFileErr != nil {
					continue
				}
				if len(content) > 0 && !isTemplateFile(content) {
					report.Violations = append(report.Violations, Issue{
						File:    entry.Name(),
						Type:    IssueSecret,
						Message: "may contain secrets (constitution violation)",
						Rule:    "no_secrets",
					})
					foundViolation = true
				}
			}
		}
	}

	if !foundViolation {
		report.Passed = append(report.Passed, CheckConstitution)
	}
}

// checkRequiredFiles verifies that all required context files are present.
//
// Checks against config.FilesRequired and adds a warning for each missing file.
//
// Parameters:
//   - ctx: Loaded context containing existing files
//   - report: Report to append warnings to (modified in place)
func checkRequiredFiles(ctx *context.Context, report *Report) {
	allPresent := true

	existingFiles := make(map[string]bool)
	for _, f := range ctx.Files {
		existingFiles[f.Name] = true
	}

	for _, name := range config.FilesRequired {
		if !existingFiles[name] {
			report.Warnings = append(report.Warnings, Issue{
				File:    name,
				Type:    IssueMissing,
				Message: "required context file is missing",
			})
			allPresent = false
		}
	}

	if allPresent {
		report.Passed = append(report.Passed, CheckRequiredFiles)
	}
}

// checkFileAge flags context files whose ModTime is older than staleAgeDays.
//
// Files listed in staleAgeExclude (e.g., CONSTITUTION.md) are skipped because
// they are expected to be static.
//
// Parameters:
//   - ctx: Loaded context containing files to check
//   - report: Report to append warnings to (modified in place)
func checkFileAge(ctx *context.Context, report *Report) {
	foundStale := false
	cutoff := time.Now().AddDate(0, 0, -staleAgeDays)

	for _, f := range ctx.Files {
		excluded := false
		for _, ex := range staleAgeExclude {
			if f.Name == ex {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		if f.ModTime.Before(cutoff) {
			days := int(time.Since(f.ModTime).Hours() / 24)
			report.Warnings = append(report.Warnings, Issue{
				File:    f.Name,
				Type:    IssueStaleAge,
				Message: fmt.Sprintf("last modified %d days ago", days),
			})
			foundStale = true
		}
	}

	if !foundStale {
		report.Passed = append(report.Passed, CheckFileAge)
	}
}

// checkEntryCount warns when LEARNINGS.md or DECISIONS.md have too many entries.
//
// Uses index.ParseEntryBlocks for counting and rc thresholds for limits.
// A threshold of 0 disables the check for that file.
//
// Parameters:
//   - ctx: Loaded context containing files to check
//   - report: Report to append warnings to (modified in place)
func checkEntryCount(ctx *context.Context, report *Report) {
	checks := []struct {
		file      string
		threshold int
	}{
		{config.FileLearning, rc.EntryCountLearnings()},
		{config.FileDecision, rc.EntryCountDecisions()},
	}

	found := false
	for _, c := range checks {
		if c.threshold <= 0 {
			continue // disabled
		}
		f := ctx.File(c.file)
		if f == nil {
			continue
		}
		blocks := index.ParseEntryBlocks(string(f.Content))
		if len(blocks) > c.threshold {
			report.Warnings = append(report.Warnings, Issue{
				File: f.Name,
				Type: IssueEntryCount,
				Message: fmt.Sprintf(
					"has %d entries (recommended: ≤%d)",
					len(blocks), c.threshold,
				),
			})
			found = true
		}
	}

	if !found {
		report.Passed = append(report.Passed, CheckEntryCount)
	}
}

// reInternalPkg matches backtick-quoted paths starting with "internal/".
var reInternalPkg = regexp.MustCompile("`(internal/[^`]+)`")

// checkMissingPackages warns about internal/ directories not referenced in ARCHITECTURE.md.
//
// Extracts backtick-quoted internal/ paths from ARCHITECTURE.md, normalizes
// them to top-level packages (e.g., internal/cli/pad → internal/cli), then
// compares against actual internal/ subdirectories. Missing coverage is
// reported as a warning.
//
// Parameters:
//   - ctx: Loaded context containing files to scan
//   - report: Report to append warnings to (modified in place)
func checkMissingPackages(ctx *context.Context, report *Report) {
	f := ctx.File(config.FileArchitecture)
	if f == nil {
		return
	}

	// Extract referenced internal/ paths and normalize to top-level packages.
	referenced := make(map[string]bool)
	matches := reInternalPkg.FindAllStringSubmatch(string(f.Content), -1)
	for _, m := range matches {
		pkg := normalizeInternalPkg(m[1])
		referenced[pkg] = true
	}

	// Scan actual internal/ subdirectories (one level deep, directories only).
	entries, err := os.ReadDir("internal")
	if err != nil {
		return
	}

	found := false
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pkg := "internal/" + entry.Name()
		if !referenced[pkg] {
			report.Warnings = append(report.Warnings, Issue{
				File:    f.Name,
				Type:    IssueMissingPackage,
				Message: fmt.Sprintf("package %s is not documented", pkg),
				Path:    pkg,
			})
			found = true
		}
	}

	if !found {
		report.Passed = append(report.Passed, CheckMissingPackages)
	}
}

// normalizeInternalPkg normalizes an internal package path to its top-level
// directory. For example, "internal/cli/pad" becomes "internal/cli" and
// "internal/recall/parser" becomes "internal/recall".
func normalizeInternalPkg(path string) string {
	// Split: ["internal", "cli", "pad"] → join first two: "internal/cli"
	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 2 {
		return path
	}
	return parts[0] + "/" + parts[1]
}

// isTemplateFile checks if file content appears to be a template.
//
// Looks for common template markers like YOUR_, {{, REPLACE_, TODO, CHANGEME.
// Used to avoid flagging template files as containing secrets.
//
// Parameters:
//   - content: File content to check
//
// Returns:
//   - bool: True if content contains template markers
func isTemplateFile(content []byte) bool {
	s := string(content)
	// Check for common template markers
	templateMarkers := []string{
		"YOUR_",
		"<your",
		"{{",
		"REPLACE_",
		"TODO",
		"CHANGEME",
		"FIXME",
	}
	for _, marker := range templateMarkers {
		if strings.Contains(strings.ToUpper(s), marker) {
			return true
		}
	}
	return false
}
