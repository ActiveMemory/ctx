//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/tpl"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// ScanKnowledgeFiles checks knowledge files against their configured
// thresholds and returns any that exceed the limits.
//
// Parameters:
//   - contextDir: absolute path to the context directory
//   - decThreshold: max decision entries (0 = disabled)
//   - lrnThreshold: max learning entries (0 = disabled)
//   - convThreshold: max convention lines (0 = disabled)
//
// Returns:
//   - []KnowledgeFinding: files exceeding thresholds, or nil if all within limits
func ScanKnowledgeFiles(
	contextDir string, decThreshold, lrnThreshold, convThreshold int,
) []KnowledgeFinding {
	var findings []KnowledgeFinding

	if decThreshold > 0 {
		if data, readErr := validation.SafeReadFile(contextDir, ctx.Decision); readErr == nil {
			count := len(index.ParseEntryBlocks(string(data)))
			if count > decThreshold {
				findings = append(findings, KnowledgeFinding{
					File: ctx.Decision, Count: count, Threshold: decThreshold, Unit: "entries",
				})
			}
		}
	}

	if lrnThreshold > 0 {
		if data, readErr := validation.SafeReadFile(contextDir, ctx.Learning); readErr == nil {
			count := len(index.ParseEntryBlocks(string(data)))
			if count > lrnThreshold {
				findings = append(findings, KnowledgeFinding{
					File: ctx.Learning, Count: count, Threshold: lrnThreshold, Unit: "entries",
				})
			}
		}
	}

	if convThreshold > 0 {
		if data, readErr := validation.SafeReadFile(contextDir, ctx.Convention); readErr == nil {
			lineCount := bytes.Count(data, []byte(token.NewlineLF))
			if lineCount > convThreshold {
				findings = append(findings, KnowledgeFinding{
					File: ctx.Convention, Count: lineCount, Threshold: convThreshold, Unit: "lines",
				})
			}
		}
	}

	return findings
}

// FormatKnowledgeWarnings builds a pre-formatted findings list string
// from the given findings.
//
// Parameters:
//   - findings: knowledge file threshold violations
//
// Returns:
//   - string: formatted warning lines for template injection
func FormatKnowledgeWarnings(findings []KnowledgeFinding) string {
	var b strings.Builder
	findingFmt := assets.TextDesc(assets.TextDescKeyCheckKnowledgeFindingFormat)
	for _, f := range findings {
		b.WriteString(fmt.Sprintf(findingFmt, f.File, f.Count, f.Unit, f.Threshold))
	}
	return b.String()
}

// EmitKnowledgeWarning builds and prints the knowledge file growth warning box.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionID: session identifier for notifications
//   - fileWarnings: pre-formatted findings text
func EmitKnowledgeWarning(cmd *cobra.Command, sessionID, fileWarnings string) {
	fallback := fileWarnings + token.NewlineLF + assets.TextDesc(assets.TextDescKeyCheckKnowledgeFallback)
	content := LoadMessage(hook.CheckKnowledge, hook.VariantWarning,
		map[string]any{tpl.VarFileWarnings: fileWarnings}, fallback)
	if content == "" {
		return
	}

	cmd.Println(NudgeBox(
		assets.TextDesc(assets.TextDescKeyCheckKnowledgeRelayPrefix),
		assets.TextDesc(assets.TextDescKeyCheckKnowledgeBoxTitle),
		content))

	ref := notify.NewTemplateRef(hook.CheckKnowledge, hook.VariantWarning,
		map[string]any{tpl.VarFileWarnings: fileWarnings})
	notifyMsg := hook.CheckKnowledge + ": " + assets.TextDesc(assets.TextDescKeyCheckKnowledgeRelayMessage)
	NudgeAndRelay(notifyMsg, sessionID, ref)
}

// CheckKnowledgeHealth runs the full knowledge health check: scans files,
// formats warnings, and emits output if any thresholds are exceeded.
// Returns true if warnings were emitted.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionID: session identifier for notifications
//
// Returns:
//   - bool: true if warnings were emitted
func CheckKnowledgeHealth(cmd *cobra.Command, sessionID string) bool {
	lrnThreshold := rc.EntryCountLearnings()
	decThreshold := rc.EntryCountDecisions()
	convThreshold := rc.ConventionLineCount()

	// All disabled — nothing to check
	if lrnThreshold == 0 && decThreshold == 0 && convThreshold == 0 {
		return false
	}

	findings := ScanKnowledgeFiles(rc.ContextDir(), decThreshold, lrnThreshold, convThreshold)
	if len(findings) == 0 {
		return false
	}

	fileWarnings := FormatKnowledgeWarnings(findings)
	EmitKnowledgeWarning(cmd, sessionID, fileWarnings)
	return true
}
