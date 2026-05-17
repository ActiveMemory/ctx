//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcemap

// Markdown rendering constants for the source-map artifact.
// Structural literals stay as Go consts; error wrapping
// format strings have moved to commands/text/errors.yaml.
const (
	// TableHeader is the markdown table header row plus its
	// delimiter row.
	TableHeader = "| Short name | Kind | Locator |" +
		" Admission status | Admission rationale | Dated |\n" +
		"|------------|------|---------|------------------|" +
		"---------------------|-------|\n"
)
