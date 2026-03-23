//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package knowledge

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/knowledge"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"

	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Finding describes a single knowledge file that exceeds its
// configured threshold.
type Finding struct {
	// File is the context filename (e.g., DECISIONS.md).
	File string
	// Count is the actual entry or line count.
	Count int
	// Threshold is the configured maximum.
	Threshold int
	// Unit is the measurement unit ("entries" or "lines").
	Unit string
}

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
) []Finding {
	var findings []Finding

	if decThreshold > 0 {
		if data, readErr := io.SafeReadFile(contextDir, ctx.Decision); readErr == nil {
			count := len(index.ParseEntryBlocks(string(data)))
			if count > decThreshold {
				findings = append(findings, Finding{
					File: ctx.Decision, Count: count, Threshold: decThreshold, Unit: "entries",
				})
			}
		}
	}

	if lrnThreshold > 0 {
		if data, readErr := io.SafeReadFile(contextDir, ctx.Learning); readErr == nil {
			count := len(index.ParseEntryBlocks(string(data)))
			if count > lrnThreshold {
				findings = append(findings, Finding{
					File: ctx.Learning, Count: count, Threshold: lrnThreshold, Unit: "entries",
				})
			}
		}
	}

	if convThreshold > 0 {
		if data, readErr := io.SafeReadFile(contextDir, ctx.Convention); readErr == nil {
			lineCount := bytes.Count(data, []byte(token.NewlineLF))
			if lineCount > convThreshold {
				findings = append(findings, Finding{
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
func FormatKnowledgeWarnings(findings []Finding) string {
	var b strings.Builder
	findingFmt := desc.Text(text.DescKeyCheckKnowledgeFindingFormat)
	for _, f := range findings {
		b.WriteString(fmt.Sprintf(findingFmt, f.File, f.Count, f.Unit, f.Threshold))
	}
	return b.String()
}

// EmitKnowledgeWarning builds the knowledge file growth warning box.
//
// Parameters:
//   - sessionID: session identifier for notifications
//   - fileWarnings: pre-formatted findings text
//
// Returns:
//   - string: formatted nudge box, or empty string if silenced
func EmitKnowledgeWarning(sessionID, fileWarnings string) string {
	fallback := fileWarnings + token.NewlineLF + desc.Text(text.DescKeyCheckKnowledgeFallback)
	content := core.LoadMessage(hook.CheckKnowledge, hook.VariantWarning,
		map[string]any{knowledge.VarFileWarnings: fileWarnings}, fallback)
	if content == "" {
		return ""
	}

	box := core.NudgeBox(
		desc.Text(text.DescKeyCheckKnowledgeRelayPrefix),
		desc.Text(text.DescKeyCheckKnowledgeBoxTitle),
		content)

	ref := notify.NewTemplateRef(hook.CheckKnowledge, hook.VariantWarning,
		map[string]any{knowledge.VarFileWarnings: fileWarnings})
	notifyMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckKnowledge, desc.Text(text.DescKeyCheckKnowledgeRelayMessage))
	core.NudgeAndRelay(notifyMsg, sessionID, ref)
	return box
}

// CheckKnowledgeHealth runs the full knowledge health check: scans files,
// formats warnings, and builds output if any thresholds are exceeded.
//
// Parameters:
//   - sessionID: session identifier for notifications
//
// Returns:
//   - string: formatted nudge box, or empty string if no warnings
//   - bool: true if warnings were found
func CheckKnowledgeHealth(sessionID string) (string, bool) {
	lrnThreshold := rc.EntryCountLearnings()
	decThreshold := rc.EntryCountDecisions()
	convThreshold := rc.ConventionLineCount()

	// All disabled — nothing to check
	if lrnThreshold == 0 && decThreshold == 0 && convThreshold == 0 {
		return "", false
	}

	findings := ScanKnowledgeFiles(rc.ContextDir(), decThreshold, lrnThreshold, convThreshold)
	if len(findings) == 0 {
		return "", false
	}

	fileWarnings := FormatKnowledgeWarnings(findings)
	box := EmitKnowledgeWarning(sessionID, fileWarnings)
	return box, true
}
