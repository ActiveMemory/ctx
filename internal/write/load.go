//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
)

// LoadRaw outputs context files without assembly or headers.
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
func LoadRaw(cmd *cobra.Command, files []context.FileInfo) error {
	for i, f := range files {
		if i > 0 {
			cmd.Println()
		}
		cmd.Print(string(f.Content))
	}
	return nil
}

// LoadAssembled outputs context as formatted Markdown with token budgeting.
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
func LoadAssembled(
	cmd *cobra.Command,
	files []context.FileInfo,
	budget, totalTokens int,
	titleFn func(string) string,
) error {
	var sb strings.Builder
	nl := config.NewlineLF
	sep := config.Separator

	sb.WriteString(config.LoadHeadingContext + nl + nl)
	_, _ = fmt.Fprintf(&sb, config.TplLoadBudget+nl+nl, budget, totalTokens)
	sb.WriteString(sep + nl + nl)

	tokensUsed := context.EstimateTokensString(sb.String())

	for _, f := range files {
		if f.IsEmpty {
			continue
		}

		fileTokens := f.Tokens
		if tokensUsed+fileTokens > budget {
			_, _ = fmt.Fprintf(&sb, nl+sep+nl+nl+config.TplLoadTruncated+nl, f.Name)
			break
		}

		_, _ = fmt.Fprintf(&sb, config.TplLoadSectionHeading+nl+nl, titleFn(f.Name))
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
