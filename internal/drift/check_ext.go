//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/steering"
	"github.com/ActiveMemory/ctx/internal/trigger"
)

// supportedTools lists the valid tool identifiers for ctx.
var supportedTools = []string{"claude", "cursor", "cline", "kiro", "codex"}

// checkSteeringTools validates that all steering files reference only
// supported tool identifiers in their tools list.
//
// Parameters:
//   - report: Report to append warnings to (modified in place)
func checkSteeringTools(report *Report) {
	steeringDir := rc.SteeringDir()

	files, err := steering.LoadAll(steeringDir)
	if err != nil {
		// Directory doesn't exist or can't be read — skip silently.
		report.Passed = append(report.Passed, CheckSteeringTools)
		return
	}

	found := false
	for _, sf := range files {
		for _, tool := range sf.Tools {
			if !slices.Contains(supportedTools, tool) {
				report.Warnings = append(report.Warnings, Issue{
					File: filepath.Base(sf.Path),
					Type: IssueInvalidTool,
					Message: fmt.Sprintf(
						desc.Text(text.DescKeyDriftInvalidTool), tool,
					),
				})
				found = true
			}
		}
	}

	if !found {
		report.Passed = append(report.Passed, CheckSteeringTools)
	}
}

// checkHookPerms scans hook directories for scripts that lack the
// executable permission bit.
//
// Parameters:
//   - report: Report to append warnings to (modified in place)
func checkHookPerms(report *Report) {
	hooksDir := rc.HooksDir()

	// Scan the raw directories to find scripts without the executable bit.
	// We don't use trigger.Discover here because it skips non-executable scripts.
	found := false
	for _, ht := range trigger.ValidTypes() {
		typeDir := filepath.Join(hooksDir, string(ht))
		entries, readErr := os.ReadDir(typeDir)
		if readErr != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			info, infoErr := e.Info()
			if infoErr != nil {
				continue
			}
			if info.Mode().Perm()&fs.ExecBitMask == 0 {
				report.Warnings = append(report.Warnings, Issue{
					File:    filepath.Join(string(ht), e.Name()),
					Type:    IssueHookNoExec,
					Message: desc.Text(text.DescKeyDriftHookNoExec),
					Path:    filepath.Join(typeDir, e.Name()),
				})
				found = true
			}
		}
	}

	if !found {
		report.Passed = append(report.Passed, CheckHookPerms)
	}
}

// checkSyncStaleness compares synced tool-native files against what
// steering.SyncTool would produce. If they differ, the synced file
// is stale.
//
// Parameters:
//   - report: Report to append warnings to (modified in place)
func checkSyncStaleness(report *Report) {
	steeringDir := rc.SteeringDir()

	files, err := steering.LoadAll(steeringDir)
	if err != nil {
		// No steering files — nothing to check.
		report.Passed = append(report.Passed, CheckSyncStaleness)
		return
	}

	if len(files) == 0 {
		report.Passed = append(report.Passed, CheckSyncStaleness)
		return
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		report.Passed = append(report.Passed, CheckSyncStaleness)
		return
	}

	found := false
	// Check each syncable tool.
	syncTools := []string{
		cfgHook.ToolCursor, cfgHook.ToolCline,
		cfgHook.ToolKiro,
	}
	for _, tool := range syncTools {
		stale := steering.StaleFiles(steeringDir, cwd, tool)
		for _, name := range stale {
			report.Warnings = append(report.Warnings, Issue{
				File:    name,
				Type:    IssueStaleSyncFile,
				Message: desc.Text(text.DescKeyDriftStaleSyncFile),
				Path:    fmt.Sprintf("%s (tool: %s)", name, tool),
			})
			found = true
		}
	}

	if !found {
		report.Passed = append(report.Passed, CheckSyncStaleness)
	}
}

// checkRCTool validates that the .ctxrc tool field contains a supported
// tool identifier.
//
// Parameters:
//   - report: Report to append warnings to (modified in place)
func checkRCTool(report *Report) {
	tool := rc.Tool()

	// Empty tool field is valid — it means no tool is configured.
	if tool == "" {
		report.Passed = append(report.Passed, CheckRCTool)
		return
	}

	if !slices.Contains(supportedTools, tool) {
		report.Warnings = append(report.Warnings, Issue{
			File: file.CtxRC,
			Type: IssueInvalidTool,
			Message: fmt.Sprintf(
				desc.Text(text.DescKeyDriftInvalidTool), tool,
			),
		})
		return
	}

	report.Passed = append(report.Passed, CheckRCTool)
}
