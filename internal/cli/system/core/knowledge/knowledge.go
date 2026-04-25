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
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/knowledge"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// ScanFiles checks knowledge files against their configured
// thresholds and returns any that exceed the limits.
//
// Parameters:
//   - contextDir: absolute path to the context directory
//   - decThreshold: max decision entries (0 = disabled)
//   - lrnThreshold: max learning entries (0 = disabled)
//   - convThreshold: max convention lines (0 = disabled)
//
// Returns:
//   - []KnowledgeFinding: files exceeding thresholds,
//     or nil if all within limits
func ScanFiles(
	contextDir string, decThreshold, lrnThreshold, convThreshold int,
) []finding {
	var findings []finding

	if decThreshold > 0 {
		data, readErr := io.SafeReadFile(contextDir, ctx.Decision)
		if readErr == nil {
			count := len(index.ParseEntryBlocks(string(data)))
			if count > decThreshold {
				findings = append(findings, finding{
					File:      ctx.Decision,
					Count:     count,
					Threshold: decThreshold,
					Unit:      desc.Text(text.DescKeyWriteKnowledgeUnitEntries),
				})
			}
		}
	}

	if lrnThreshold > 0 {
		data, readErr := io.SafeReadFile(contextDir, ctx.Learning)
		if readErr == nil {
			count := len(index.ParseEntryBlocks(string(data)))
			if count > lrnThreshold {
				findings = append(findings, finding{
					File:      ctx.Learning,
					Count:     count,
					Threshold: lrnThreshold,
					Unit:      desc.Text(text.DescKeyWriteKnowledgeUnitEntries),
				})
			}
		}
	}

	if convThreshold > 0 {
		data, readErr := io.SafeReadFile(contextDir, ctx.Convention)
		if readErr == nil {
			lineCount := bytes.Count(data, []byte(token.NewlineLF))
			if lineCount > convThreshold {
				findings = append(findings, finding{
					File:      ctx.Convention,
					Count:     lineCount,
					Threshold: convThreshold,
					Unit:      desc.Text(text.DescKeyWriteKnowledgeUnitLines),
				})
			}
		}
	}

	return findings
}

// FormatWarnings builds a pre-formatted findings list string
// from the given findings.
//
// Parameters:
//   - findings: knowledge file threshold violations
//
// Returns:
//   - string: formatted warning lines for template injection
func FormatWarnings(findings []finding) string {
	var b strings.Builder
	findingFmt := desc.Text(text.DescKeyCheckKnowledgeFindingFormat)
	for _, f := range findings {
		io.SafeFprintf(&b, findingFmt, f.File, f.Count, f.Unit, f.Threshold)
	}
	return b.String()
}

// EmitWarning builds the knowledge file growth warning box.
//
// Parameters:
//   - sessionID: session identifier for notifications
//   - fileWarnings: pre-formatted findings text
//
// Returns:
//   - string: formatted nudge box, or empty string if silenced
//   - error: propagated from [nudge.EmitAndRelay] so callers can
//     honor the log-first principle: if the relay audit entry or
//     webhook fails, the nudge box should not be printed.
func EmitWarning(sessionID, fileWarnings string) (string, error) {
	fallback := fileWarnings + token.NewlineLF + desc.Text(
		text.DescKeyCheckKnowledgeFallback,
	)
	content := message.Load(hook.CheckKnowledge, hook.VariantWarning,
		map[string]any{knowledge.VarFileWarnings: fileWarnings}, fallback)
	if content == "" {
		return "", nil
	}

	box := message.NudgeBox(
		desc.Text(text.DescKeyCheckKnowledgeRelayPrefix),
		desc.Text(text.DescKeyCheckKnowledgeBoxTitle),
		content)

	ref := notify.NewTemplateRef(hook.CheckKnowledge, hook.VariantWarning,
		map[string]any{knowledge.VarFileWarnings: fileWarnings})
	notifyMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckKnowledge, desc.Text(text.DescKeyCheckKnowledgeRelayMessage))
	if err := nudge.EmitAndRelay(notifyMsg, sessionID, ref); err != nil {
		return "", err
	}
	return box, nil
}

// CheckHealth runs the full knowledge health check: scans files,
// formats warnings, and builds output if any thresholds are exceeded.
//
// ctxDir is supplied by the caller (typically a FullPreamble-gated
// hook) so this function does not re-resolve it; a second resolution
// would be dead code today and would ambiguously pair (false, err)
// with the genuine "no warnings found" return value.
//
// Parameters:
//   - sessionID: session identifier for notifications
//   - ctxDir: absolute path to the context directory
//
// Returns:
//   - string: formatted nudge box, or empty string if no warnings
//   - bool: true if warnings were found
//   - error: propagated from [EmitWarning] so callers can honour the
//     log-first principle and skip printing the box when the relay
//     audit entry could not be written.
func CheckHealth(sessionID, ctxDir string) (string, bool, error) {
	lrnThreshold := rc.EntryCountLearnings()
	decThreshold := rc.EntryCountDecisions()
	convThreshold := rc.ConventionLineCount()

	// All disabled - nothing to check
	if lrnThreshold == 0 && decThreshold == 0 && convThreshold == 0 {
		return "", false, nil
	}

	findings := ScanFiles(
		ctxDir, decThreshold, lrnThreshold, convThreshold,
	)
	if len(findings) == 0 {
		return "", false, nil
	}

	fileWarnings := FormatWarnings(findings)
	box, emitErr := EmitWarning(sessionID, fileWarnings)
	if emitErr != nil {
		return "", false, emitErr
	}
	return box, true, nil
}
