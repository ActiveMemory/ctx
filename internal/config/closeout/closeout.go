//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package closeout

// YAML field-name constants used by the closeout frontmatter
// parser to report missing fields. These mirror the YAML tags
// on [entity.CloseoutFrontmatter]; they are structural
// identifiers, not localizable prose.
const (
	// FieldSHA is the YAML key for the short git SHA.
	FieldSHA = "sha"
	// FieldBranch is the YAML key for the git branch name.
	FieldBranch = "branch"
	// FieldMode is the YAML key for the closeout mode.
	FieldMode = "mode"
	// FieldGeneratedAt is the YAML key for the write timestamp.
	FieldGeneratedAt = "generated-at"
)
