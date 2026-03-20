//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// regInternalPkg matches backtick-quoted paths starting with "internal/".
var regInternalPkg = regexp.MustCompile("`(internal/[^`]+)`")

// staleAgeExclude lists context files that are expected to be static
// and should not trigger file-age warnings.
var staleAgeExclude = []string{ctxCfg.Constitution}

// checkPathReferences scans ARCHITECTURE.md and CONVENTIONS.md for dead paths.
//
// Looks for backtick-enclosed file paths and verifies they exist on disk.
// Skips URLs, template patterns, and glob patterns.
//
// Parameters:
//   - ctx: Loaded context containing files to scan
//   - report: Report to append warnings to (modified in place)
func checkPathReferences(ctx *entity.Context, report *Report) {
	foundDeadPaths := false

	for _, f := range ctx.Files {
		if f.Name != ctxCfg.Architecture && f.Name != ctxCfg.Convention {
			continue
		}

		lines := strings.Split(string(f.Content), token.NewlineLF)
		for lineNum, line := range lines {
			matches := regex.CodeFencePath.FindAllStringSubmatch(line, -1)
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
				// Skip illustrative examples: bare filenames (no /)
				// and shallow paths whose top-level directory doesn't
				// exist in the project tree. Real references point
				// into actual directories (internal/, cmd/, docs/).
				topDir := strings.SplitN(path, "/", 2)[0]
				if _, dirErr := os.Stat(topDir); os.IsNotExist(dirErr) {
					continue
				}
				// Check if the file exists
				if _, err := os.Stat(path); os.IsNotExist(err) {
					report.Warnings = append(report.Warnings, Issue{
						File:    f.Name,
						Line:    lineNum + 1,
						Type:    IssueDeadPath,
						Message: desc.TextDesc(text.DescKeyDriftDeadPath),
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
func checkStaleness(ctx *entity.Context, report *Report) {
	staleness := false

	if f := ctx.File(ctxCfg.Task); f != nil {
		// Count completed tasks
		completedCount := strings.Count(string(f.Content), marker.PrefixTaskDone)
		if completedCount > 10 {
			report.Warnings = append(report.Warnings, Issue{
				File:    f.Name,
				Type:    IssueStaleness,
				Message: desc.TextDesc(text.DescKeyDriftStaleness),
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
func checkConstitution(_ *entity.Context, report *Report) {
	// Basic heuristic checks for constitution violations
	// Check for potential secrets in common config files

	secretPatterns := token.SecretPatterns

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
				if len(content) > 0 && !templateFile(content) {
					report.Violations = append(report.Violations, Issue{
						File:    entry.Name(),
						Type:    IssueSecret,
						Message: desc.TextDesc(text.DescKeyDriftSecret),
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
func checkRequiredFiles(ctx *entity.Context, report *Report) {
	allPresent := true

	existingFiles := make(map[string]bool)
	for _, f := range ctx.Files {
		existingFiles[f.Name] = true
	}

	for _, name := range ctxCfg.FilesRequired {
		if !existingFiles[name] {
			report.Warnings = append(report.Warnings, Issue{
				File:    name,
				Type:    IssueMissing,
				Message: desc.TextDesc(text.DescKeyDriftMissingFile),
			})
			allPresent = false
		}
	}

	if allPresent {
		report.Passed = append(report.Passed, CheckRequiredFiles)
	}
}

// checkFileAge flags context files whose ModTime is older than
// rc.StaleAgeDays.
//
// Files listed in staleAgeExclude (e.g., CONSTITUTION.md) are skipped
// because they are expected to be static. The check is skipped entirely
// when stale_age_days is 0 in .ctxrc.
//
// Parameters:
//   - ctx: Loaded context containing files to check
//   - report: Report to append warnings to (modified in place)
func checkFileAge(ctx *entity.Context, report *Report) {
	days := rc.StaleAgeDays()
	if days == 0 {
		return
	}
	foundStale := false
	cutoff := time.Now().AddDate(0, 0, -days)

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
				Message: fmt.Sprintf(desc.TextDesc(text.DescKeyDriftStaleAge), days),
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
func checkEntryCount(ctx *entity.Context, report *Report) {
	checks := []struct {
		file      string
		threshold int
	}{
		{ctxCfg.Learning, rc.EntryCountLearnings()},
		{ctxCfg.Decision, rc.EntryCountDecisions()},
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
					desc.TextDesc(text.DescKeyDriftEntryCount),
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

// checkMissingPackages warns about internal/ directories not referenced
// in ARCHITECTURE.md.
//
// Extracts backtick-quoted internal/ paths from ARCHITECTURE.md, normalizes
// them to top-level packages (e.g., internal/cli/pad → internal/cli), then
// compares against actual internal/ subdirectories. Missing coverage is
// reported as a warning.
//
// Parameters:
//   - ctx: Loaded context containing files to scan
//   - report: Report to append warnings to (modified in place)
func checkMissingPackages(ctx *entity.Context, report *Report) {
	f := ctx.File(ctxCfg.Architecture)
	if f == nil {
		return
	}

	// Extract referenced internal/ paths and normalize to top-level packages.
	referenced := make(map[string]bool)
	matches := regInternalPkg.FindAllStringSubmatch(string(f.Content), -1)
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
				File: f.Name,
				Type: IssueMissingPackage,
				Message: fmt.Sprintf(
					desc.TextDesc(text.DescKeyDriftMissingPackage), pkg,
				),
				Path: pkg,
			})
			found = true
		}
	}

	if !found {
		report.Passed = append(report.Passed, CheckMissingPackages)
	}
}
