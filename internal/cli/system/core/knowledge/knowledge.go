//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package knowledge

import (
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
	"github.com/ActiveMemory/ctx/internal/disclosure"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Health scans the canonical knowledge roots and returns the two M5
// signals: foldable roots (staging accretion → /ctx-digest) and heavy
// pages (byte weight over the ceiling → split / extract-to-tooling). One
// scan feeds the check-knowledge hook and the skill report path alike, so
// they cannot drift.
//
// Foldability is the staging-zone count via disclosure.StagedEntries, so
// it reads a never-migrated root correctly and quiets once a root is
// folded. Weight is bytes over the root AND every theme file, closing the
// blindness the root-only measure had to the bulk folding relocates.
//
// Findings are ordered foldable-first: when a root trips both, folding is
// the single move that reduces both.
//
// Parameters:
//   - contextDir: absolute path to the context directory
//   - t: the four thresholds (0 disables a check)
//
// Returns:
//   - []finding: signals found, or nil when every root is within limits
func Health(contextDir string, t Thresholds) []finding {
	roots := []struct {
		file      string
		threshold int
	}{
		{ctx.Learning, t.Learnings},
		{ctx.Decision, t.Decisions},
		{ctx.Convention, t.Conventions},
	}

	var foldables, heavies []finding
	for _, r := range roots {
		data, readErr := io.SafeReadFile(contextDir, r.file)
		if readErr != nil {
			continue
		}
		kind, ok := disclosure.KindFor(r.file)
		if !ok {
			continue
		}

		// Signal 1 — foldable root (staging count).
		if r.threshold > 0 {
			n := len(disclosure.StagedEntries(disclosure.Parse(string(data), kind)))
			if n > r.threshold {
				foldables = append(foldables, finding{
					Kind: foldable, File: r.file, Count: n,
					Threshold: r.threshold, Unit: foldUnit(kind),
				})
			}
		}

		// Signal 2 — heavy pages (bytes): the root, then its theme files.
		if t.PageBytes > 0 {
			if len(data) > t.PageBytes {
				heavies = append(heavies, heavyFinding(r.file, len(data), t.PageBytes))
			}
			heavies = append(heavies, heavyThemeFiles(contextDir, kind, t.PageBytes)...)
		}
	}

	if len(foldables) == 0 && len(heavies) == 0 {
		return nil
	}
	return append(foldables, heavies...)
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
	var foldables, heavies []finding
	for _, f := range findings {
		if f.Kind == heavy {
			heavies = append(heavies, f)
			continue
		}
		foldables = append(foldables, f)
	}

	findingFmt := desc.Text(text.DescKeyCheckKnowledgeFindingFormat)
	var b strings.Builder
	writeGroup := func(group []finding, remedyKey string) {
		if len(group) == 0 {
			return
		}
		if b.Len() > 0 {
			b.WriteString(token.NewlineLF)
		}
		for _, f := range group {
			io.SafeFprintf(&b, findingFmt, f.File, f.Count, f.Unit, f.Threshold)
		}
		b.WriteString(desc.Text(remedyKey))
	}
	// Foldable first: when a root trips both, folding is the move that
	// reduces both, so its remedy leads.
	writeGroup(foldables, text.DescKeyCheckKnowledgeRemedyFoldable)
	writeGroup(heavies, text.DescKeyCheckKnowledgeRemedyHeavy)
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
	// fileWarnings already carries its per-kind remedy lines (see
	// FormatWarnings), so it is a complete fallback on template-load
	// failure — no generic remedy to append.
	content := message.Load(hook.CheckKnowledge, hook.VariantWarning,
		map[string]any{knowledge.VarFileWarnings: fileWarnings}, fileWarnings)
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
	t := thresholds()

	// All disabled - nothing to check
	if t.Learnings == 0 && t.Decisions == 0 &&
		t.Conventions == 0 && t.PageBytes == 0 {
		return "", false, nil
	}

	findings := Health(ctxDir, t)
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

// Report returns the formatted knowledge-health findings for on-demand
// display by the /ctx-remember and /ctx-wrap-up skills. Unlike
// [CheckHealth] it neither throttles nor relays: the skills call it
// deliberately, want the current state every time, and must not spam the
// hook audit trail. Empty string when every root is within limits.
//
// Parameters:
//   - ctxDir: absolute path to the context directory
//
// Returns:
//   - string: the formatted findings + per-kind remedies, or "" if clean
func Report(ctxDir string) string {
	findings := Health(ctxDir, thresholds())
	if len(findings) == 0 {
		return ""
	}
	return FormatWarnings(findings)
}
