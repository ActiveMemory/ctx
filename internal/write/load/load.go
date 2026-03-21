//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package load

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxToken "github.com/ActiveMemory/ctx/internal/context/token"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/entity"
)

// Raw outputs context files without assembly or headers.
//
// Files are output in read order, separated by blank lines.
// Content is printed as-is without modification.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - files: Context files sorted by read order
//
// Returns:
//   - error: Always nil (included for interface consistency)
func Raw(cmd *cobra.Command, files []entity.FileInfo) error {
	for i, f := range files {
		if i > 0 {
			cmd.Println()
		}
		cmd.Print(string(f.Content))
	}
	return nil
}

// Assembled outputs context as formatted Markdown with token budgeting.
//
// Assembles context files into a single Markdown document with headers,
// respecting the token budget. Files are included in read order until the
// budget is exhausted. Truncated files are noted in the output.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - files: Context files sorted by read order
//   - budget: Maximum token count for the output
//   - totalTokens: Total available tokens in context
//   - titleFn: Function to convert filename to display title
//
// Returns:
//   - error: Always nil (included for interface consistency)
func Assembled(
	cmd *cobra.Command,
	files []entity.FileInfo,
	budget, totalTokens int,
	titleFn func(string) string,
) error {
	var sb strings.Builder
	nl := token.NewlineLF
	sep := token.Separator

	sb.WriteString(desc.Text(text.DescKeyHeadingContext) + nl + nl)
	_, _ = fmt.Fprintf(&sb, tpl.LoadBudget+nl+nl, budget, totalTokens)
	sb.WriteString(sep + nl + nl)

	tokensUsed := ctxToken.EstimateTokensString(sb.String())

	for _, f := range files {
		if f.IsEmpty {
			continue
		}

		fileTokens := f.Tokens
		if tokensUsed+fileTokens > budget {
			_, _ = fmt.Fprintf(&sb, nl+sep+nl+nl+tpl.LoadTruncated+nl, f.Name)
			break
		}

		_, _ = fmt.Fprintf(&sb, tpl.LoadSectionHeading+nl+nl, titleFn(f.Name))
		sb.Write(f.Content)
		if !strings.HasSuffix(string(f.Content), nl) {
			sb.WriteString(nl)
		}
		sb.WriteString(nl + sep + nl + nl)

		tokensUsed += fileTokens
	}

	cmd.Print(sb.String())
	return nil
}
