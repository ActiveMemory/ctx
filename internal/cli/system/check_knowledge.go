//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// checkKnowledgeCmd returns the "ctx system check-knowledge" hook command.
//
// Counts entries in DECISIONS.md and LEARNINGS.md and outputs a VERBATIM
// relay nudge when either exceeds the configured threshold. Daily throttle
// prevents repeated nudges within the same day.
//
// Hook event: UserPromptSubmit
// Output: VERBATIM relay (when thresholds exceeded), silent otherwise
// Silent when: below thresholds, already nudged today, or uninitialized
func checkKnowledgeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-knowledge",
		Short: "Knowledge file growth nudge",
		Long: `Counts entries in DECISIONS.md and LEARNINGS.md and lines in
CONVENTIONS.md, and outputs a VERBATIM relay nudge when any file exceeds
the configured threshold. Throttled to once per day.

  Learnings threshold:   entry_count_learnings   (default 30)
  Decisions threshold:   entry_count_decisions    (default 20)
  Conventions threshold: convention_line_count    (default 200)

Hook event: UserPromptSubmit
Output: VERBATIM relay (when thresholds exceeded), silent otherwise
Silent when: below thresholds, already nudged today, or uninitialized`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckKnowledge(cmd)
		},
	}
}

func runCheckKnowledge(cmd *cobra.Command) error {
	if !isInitialized() {
		return nil
	}

	markerPath := filepath.Join(secureTempDir(), "check-knowledge")
	if isDailyThrottled(markerPath) {
		return nil
	}

	lrnThreshold := rc.EntryCountLearnings()
	decThreshold := rc.EntryCountDecisions()
	convThreshold := rc.ConventionLineCount()

	// All disabled — nothing to check
	if lrnThreshold == 0 && decThreshold == 0 && convThreshold == 0 {
		return nil
	}

	contextDir := rc.ContextDir()

	type finding struct {
		file      string
		count     int
		threshold int
		unit      string
	}
	var findings []finding

	if decThreshold > 0 {
		decPath := filepath.Join(contextDir, config.FileDecision)
		if data, err := os.ReadFile(decPath); err == nil { //nolint:gosec // project-local path
			count := len(index.ParseEntryBlocks(string(data)))
			if count > decThreshold {
				findings = append(findings, finding{
					file: config.FileDecision, count: count, threshold: decThreshold, unit: "entries",
				})
			}
		}
	}

	if lrnThreshold > 0 {
		lrnPath := filepath.Join(contextDir, config.FileLearning)
		if data, err := os.ReadFile(lrnPath); err == nil { //nolint:gosec // project-local path
			count := len(index.ParseEntryBlocks(string(data)))
			if count > lrnThreshold {
				findings = append(findings, finding{
					file: config.FileLearning, count: count, threshold: lrnThreshold, unit: "entries",
				})
			}
		}
	}

	if convThreshold > 0 {
		convPath := filepath.Join(contextDir, config.FileConvention)
		if data, err := os.ReadFile(convPath); err == nil { //nolint:gosec // project-local path
			lineCount := bytes.Count(data, []byte("\n"))
			if lineCount > convThreshold {
				findings = append(findings, finding{
					file: config.FileConvention, count: lineCount, threshold: convThreshold, unit: "lines",
				})
			}
		}
	}

	if len(findings) == 0 {
		return nil
	}

	cmd.Println("IMPORTANT: Relay this knowledge health notice to the user VERBATIM before answering their question.")
	cmd.Println()
	cmd.Println("\u250c\u2500 Knowledge File Growth \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500")
	for _, f := range findings {
		cmd.Println(fmt.Sprintf("\u2502 %s has %d %s (recommended: \u2264%d).", f.file, f.count, f.unit, f.threshold))
	}
	cmd.Println("\u2502")
	cmd.Println("\u2502 Large knowledge files dilute agent context. Consider:")
	cmd.Println("\u2502  \u2022 Review and remove outdated entries")
	cmd.Println("\u2502  \u2022 Use /ctx-consolidate to merge overlapping entries")
	cmd.Println("\u2502  \u2022 Use /ctx-drift for semantic drift (stale patterns)")
	cmd.Println("\u2502  \u2022 Move stale entries to .context/archive/ manually")
	if line := contextDirLine(); line != "" {
		cmd.Println("\u2502 " + line)
	}
	cmd.Println("\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500")

	_ = notify.Send("nudge", "check-knowledge: Knowledge file growth detected", "")
	_ = notify.Send("relay", "check-knowledge: Knowledge file growth detected", "")

	touchFile(markerPath)

	return nil
}
