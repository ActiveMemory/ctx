//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handover

import (
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	cfgHandover "github.com/ActiveMemory/ctx/internal/config/handover"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	errHandover "github.com/ActiveMemory/ctx/internal/err/handover"
)

// composeMarkdown builds the handover's full file content:
// frontmatter YAML block plus body.
//
// Parameters:
//   - fm: frontmatter.
//   - body: rendered body markdown.
//
// Returns:
//   - string: full file content with frontmatter delimiters.
//   - error: non-nil on YAML marshal failure.
func composeMarkdown(fm Frontmatter, body string) (string, error) {
	var sb strings.Builder
	sb.WriteString(token.Separator)
	sb.WriteString(token.NewlineLF)
	enc, yamlErr := yaml.Marshal(fm)
	if yamlErr != nil {
		return "", errHandover.MarshalFrontmatter(yamlErr)
	}
	sb.Write(enc)
	sb.WriteString(token.Separator)
	sb.WriteString(token.DoubleNewline)
	sb.WriteString(strings.TrimRight(body, token.NewlineLF))
	sb.WriteString(token.NewlineLF)
	return sb.String(), nil
}

// renderBody composes the handover's markdown body sections.
//
// Parameters:
//   - entry: caller-supplied content fields.
//   - folded: closeouts folded into this handover.
//
// Returns:
//   - string: rendered markdown body.
func renderBody(entry Entry, folded []entity.CloseoutFile) string {
	var sb strings.Builder
	sb.WriteString(cfgHandover.SectionSummary)
	sb.WriteString(token.DoubleNewline)
	sb.WriteString(strings.TrimSpace(entry.Summary))
	sb.WriteString(token.DoubleNewline)
	sb.WriteString(cfgHandover.SectionNext)
	sb.WriteString(token.DoubleNewline)
	sb.WriteString(strings.TrimSpace(entry.Next))
	sb.WriteString(token.NewlineLF)

	if h := strings.TrimSpace(entry.Highlights); h != "" {
		sb.WriteString(token.NewlineLF)
		sb.WriteString(cfgHandover.SectionHighlights)
		sb.WriteString(token.DoubleNewline)
		sb.WriteString(h)
		sb.WriteString(token.NewlineLF)
	}
	if q := strings.TrimSpace(entry.OpenQuestions); q != "" {
		sb.WriteString(token.NewlineLF)
		sb.WriteString(cfgHandover.SectionOpenQuestions)
		sb.WriteString(token.DoubleNewline)
		sb.WriteString(q)
		sb.WriteString(token.NewlineLF)
	}

	if len(folded) > 0 {
		sb.WriteString(token.NewlineLF)
		sb.WriteString(cfgHandover.SectionFoldedCloseouts)
		sb.WriteString(token.DoubleNewline)
		for _, f := range folded {
			sb.WriteString(cfgHandover.FoldEntryPrefix)
			sb.WriteString(filepath.Base(f.Path))
			sb.WriteString(cfgHandover.FoldEntryModePrefix)
			sb.WriteString(f.Frontmatter.Mode)
			if f.Frontmatter.PassMode != "" {
				sb.WriteString(cfgHandover.FoldEntryPassModePrefix)
				sb.WriteString(f.Frontmatter.PassMode)
			}
			sb.WriteString(cfgHandover.FoldEntryGeneratedAtPrefix)
			sb.WriteString(f.Frontmatter.GeneratedAt.Format(time.RFC3339))
			sb.WriteString(token.NewlineLF)
		}
	}

	return sb.String()
}
