//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package glossary

// Markdown rendering constants for the glossary artifact.
// These are structural literals (table header bytes) and stay
// here as Go consts. Error wrapping format strings have moved
// to commands/text/errors.yaml and are resolved via desc.Text
// inside internal/err/kb/glossary/.
const (
	// TableHeader is the markdown table header row plus its
	// delimiter row, written when the file is first created.
	TableHeader = "| Term | Definition | Confidence |" +
		" Evidence | Related terms |\n" +
		"|------|------------|------------|----------|" +
		"---------------|\n"
)
