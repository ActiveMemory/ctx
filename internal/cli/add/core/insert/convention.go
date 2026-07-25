//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package insert

// Convention inserts a convention bullet at the top of an H2 section.
//
// A CONVENTIONS root groups bullets under "## " sections, and the
// digesting pass addresses those sections — so a bullet is only reachable
// if it sits inside one. Appending at EOF (the pre-M4 behavior) put new
// bullets below "## Themes" once a root was folded, breaking the
// entry-below-themes invariant; inserting above the first section would
// put them in the preamble, where nothing enumerates them.
//
// Target resolution:
//   - an explicit section name is honored, and created when absent
//   - otherwise the root's first section receives the bullet
//   - a root with no sections at all gets [cfgDisc.SectionDefaultConvention]
//
// A created section is placed above "## Themes", keeping it in the
// staging zone.
//
// Parameters:
//   - content: Existing file content
//   - entry: Formatted convention bullet ("- text\n")
//   - section: Target section name, or "" to resolve one
//
// Returns:
//   - []byte: Modified content with the bullet inserted
func Convention(content, entry, section string) []byte {
	header := conventionTarget(content, section)
	if idx, ok := conventionBodyStart(content, header); ok {
		return []byte(content[:idx] + entry + content[idx:])
	}
	return createConventionSection(content, entry, header)
}
