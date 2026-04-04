//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"strings"

	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errSteering "github.com/ActiveMemory/ctx/internal/err/steering"
)

// syncableTools lists the tool identifiers that support
// native-format sync. Claude and Codex use ctx agent
// directly and do not need synced files.
var syncableTools = []string{
	cfgHook.ToolCursor,
	cfgHook.ToolCline,
	cfgHook.ToolKiro,
}

// SyncTool writes steering files to the tool-native format directory.
// It loads all steering files from steeringDir, filters out files whose
// tools list excludes the target tool, formats each file in the tool's
// native format, and writes it to the appropriate output directory under
// projectRoot.
//
// Files whose content hasn't changed are skipped (idempotent).
// Output paths are validated to resolve within the project root boundary.
//
// Supported tools: cursor, cline, kiro.
func SyncTool(
	steeringDir, projectRoot, tool string,
) (SyncReport, error) {
	if !syncableTool(tool) {
		supported := strings.Join(
			syncableTools, token.CommaSpace,
		)
		return SyncReport{}, errSteering.UnsupportedTool(
			tool, supported,
		)
	}

	files, loadErr := LoadAll(steeringDir)
	if loadErr != nil {
		return SyncReport{}, loadErr
	}

	var report SyncReport
	for _, sf := range files {
		if !matchTool(sf, tool) {
			report.Skipped = append(report.Skipped, sf.Name)
			continue
		}

		outPath := nativePath(projectRoot, tool, sf.Name)

		if validateErr := validateOutputPath(
			outPath, projectRoot,
		); validateErr != nil {
			report.Errors = append(
				report.Errors,
				errSteering.SyncName(sf.Name, validateErr),
			)
			continue
		}

		content := formatNative(tool, sf)

		if unchanged(outPath, content) {
			report.Skipped = append(report.Skipped, sf.Name)
			continue
		}

		if writeErr := writeFile(outPath, content); writeErr != nil {
			report.Errors = append(
				report.Errors,
				errSteering.WriteFile(outPath, writeErr),
			)
			continue
		}

		report.Written = append(report.Written, sf.Name)
	}

	return report, nil
}

// SyncAll syncs steering files to all supported
// tool-native formats. It calls SyncTool for each
// syncable tool and merges the reports.
func SyncAll(
	steeringDir, projectRoot string,
) (SyncReport, error) {
	var merged SyncReport
	for _, tool := range syncableTools {
		r, err := SyncTool(steeringDir, projectRoot, tool)
		if err != nil {
			return merged, errSteering.SyncAll(tool, err)
		}
		merged.Written = append(merged.Written, r.Written...)
		merged.Skipped = append(merged.Skipped, r.Skipped...)
		merged.Errors = append(merged.Errors, r.Errors...)
	}
	return merged, nil
}

// StaleFiles returns the names of steering files whose synced
// tool-native output differs from what SyncTool would produce.
// This is a read-only check — no files are written.
//
// Returns nil if no stale files are found or if the steering
// directory cannot be read.
func StaleFiles(steeringDir, projectRoot, tool string) []string {
	if !syncableTool(tool) {
		return nil
	}

	files, err := LoadAll(steeringDir)
	if err != nil {
		return nil
	}

	var stale []string
	for _, sf := range files {
		if !matchTool(sf, tool) {
			continue
		}
		outPath := nativePath(projectRoot, tool, sf.Name)
		content := formatNative(tool, sf)
		if !unchanged(outPath, content) {
			stale = append(stale, sf.Name)
		}
	}
	return stale
}
