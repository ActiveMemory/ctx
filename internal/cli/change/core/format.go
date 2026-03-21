//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/format"
)

// RenderChanges renders the full CLI output for `ctx changes`.
//
// Parameters:
//   - refLabel: Human-readable reference time label
//   - ctxChanges: Context file changes since reference time
//   - code: Code change summary since reference time
//
// Returns:
//   - string: Formatted Markdown output
func RenderChanges(
	refLabel string, ctxChanges []ContextChange, code CodeSummary,
) string {
	var b strings.Builder

	b.WriteString(
		desc.TextDesc(text.DescKeyChangesHeading) +
			token.NewlineLF + token.NewlineLF,
	)
	b.WriteString(fmt.Sprintf(
		desc.TextDesc(text.DescKeyChangesRefPoint)+
			token.NewlineLF+token.NewlineLF, refLabel,
	),
	)

	if len(ctxChanges) > 0 {
		b.WriteString(
			desc.TextDesc(text.DescKeyChangesCtxHeading) + token.NewlineLF,
		)
		for _, c := range ctxChanges {
			b.WriteString(fmt.Sprintf(
				desc.TextDesc(text.DescKeyChangesCtxLine)+token.NewlineLF,
				c.Name, c.ModTime.Format(cfgTime.DateTimeFormat)))
		}
		b.WriteString(token.NewlineLF)
	}

	if code.CommitCount > 0 {
		b.WriteString(
			desc.TextDesc(text.DescKeyChangesCodeHeading) + token.NewlineLF,
		)
		b.WriteString(fmt.Sprintf(
			desc.TextDesc(text.DescKeyChangesCodeCommits)+token.NewlineLF,
			format.Pluralize(code.CommitCount, desc.TextDesc(text.DescKeyTimeCommit))))
		if code.LatestMsg != "" {
			b.WriteString(fmt.Sprintf(
				desc.TextDesc(
					text.DescKeyChangesCodeLatest)+token.NewlineLF, code.LatestMsg,
			),
			)
		}
		if len(code.Dirs) > 0 {
			b.WriteString(fmt.Sprintf(
				desc.TextDesc(text.DescKeyChangesCodeDirs)+token.NewlineLF,
				strings.Join(code.Dirs, token.CommaSpace)))
		}
		if len(code.Authors) > 0 {
			b.WriteString(fmt.Sprintf(
				desc.TextDesc(text.DescKeyChangesCodeAuthors)+token.NewlineLF,
				strings.Join(code.Authors, token.CommaSpace)))
		}
		b.WriteString(token.NewlineLF)
	}

	if len(ctxChanges) == 0 && code.CommitCount == 0 {
		b.WriteString(desc.TextDesc(text.DescKeyChangesNone) + token.NewlineLF)
	}

	return b.String()
}

// RenderChangesForHook renders a compact summary for hook injection.
//
// Parameters:
//   - refLabel: Human-readable reference time label
//   - ctxChanges: Context file changes since reference time
//   - code: Code change summary since reference time
//
// Returns:
//   - string: Compact single-line summary, or empty if no changes
func RenderChangesForHook(
	refLabel string, ctxChanges []ContextChange, code CodeSummary,
) string {
	var parts []string

	if len(ctxChanges) > 0 {
		names := make([]string, len(ctxChanges))
		for i, c := range ctxChanges {
			names[i] = c.Name
		}
		parts = append(parts, fmt.Sprintf(
			desc.TextDesc(text.DescKeyChangesHookCtxFiles),
			refLabel, strings.Join(
				names, token.CommaSpace),
		),
		)
	}

	if code.CommitCount > 0 {
		msg := fmt.Sprintf(
			desc.TextDesc(text.DescKeyChangesHookCommits),
			format.Pluralize(code.CommitCount, desc.TextDesc(text.DescKeyTimeCommit)))
		if code.LatestMsg != "" {
			msg += fmt.Sprintf(
				desc.TextDesc(text.DescKeyChangesHookCommitsExtra), code.LatestMsg,
			)
		}
		parts = append(parts, msg)
	}

	if len(parts) == 0 {
		return ""
	}

	return desc.TextDesc(
		text.DescKeyChangesHookPrefix,
	) + strings.Join(
		parts, token.PeriodSpace,
	) + token.NewlineLF
}
