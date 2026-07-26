//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package knowledge

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/disclosure"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// thresholds builds the four M5 limits from rc, shared by CheckHealth
// (the throttled hook path) and Report (the on-demand skill path).
//
// Returns:
//   - Thresholds: the current rc-configured limits
func thresholds() Thresholds {
	return Thresholds{
		Learnings:   rc.EntryCountLearnings(),
		Decisions:   rc.EntryCountDecisions(),
		Conventions: rc.ConventionSectionCount(),
		PageBytes:   rc.ThemePageByteCeiling(),
	}
}

// foldUnit is the foldability count's unit: sections for a convention
// root (its identity is the section title), entries for the timestamped
// kinds.
//
// Parameters:
//   - kind: the root kind being measured
//
// Returns:
//   - string: the resolved unit noun ("sections" | "entries")
func foldUnit(kind disclosure.Kind) string {
	if kind == disclosure.KindConvention {
		return desc.Text(text.DescKeyWriteKnowledgeUnitSections)
	}
	return desc.Text(text.DescKeyWriteKnowledgeUnitEntries)
}

// heavyFinding builds a heavy-page finding for a page of size bytes.
//
// Parameters:
//   - file: the page's display name (root basename or theme-file path)
//   - size: the page's byte size
//   - ceiling: the configured byte ceiling it exceeded
//
// Returns:
//   - finding: a heavy-kind finding for the page
func heavyFinding(file string, size, ceiling int) finding {
	return finding{
		Kind: heavy, File: file, Count: size,
		Threshold: ceiling,
		Unit:      desc.Text(text.DescKeyWriteKnowledgeUnitBytes),
	}
}

// heavyThemeFiles returns heavy findings for a kind's theme files that
// exceed the byte ceiling, sized by stat rather than a full read. A kind
// with no theme directory (never folded) yields nothing.
//
// Parameters:
//   - contextDir: absolute path to the context directory
//   - kind: the root kind whose theme directory is scanned
//   - ceiling: the byte ceiling a theme file must exceed to be flagged
//
// Returns:
//   - []finding: heavy findings for oversized theme files, nil if none
func heavyThemeFiles(
	contextDir string, kind disclosure.Kind, ceiling int,
) []finding {
	noun, ok := disclosure.ThemeDir(kind)
	if !ok {
		return nil
	}
	entries, readErr := os.ReadDir(filepath.Join(contextDir, noun))
	if readErr != nil {
		return nil
	}
	var out []finding
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), cfgFile.ExtMarkdown) {
			continue
		}
		info, infoErr := e.Info()
		if infoErr != nil {
			continue
		}
		if int(info.Size()) > ceiling {
			out = append(out, heavyFinding(
				filepath.Join(noun, e.Name()), int(info.Size()), ceiling))
		}
	}
	return out
}
