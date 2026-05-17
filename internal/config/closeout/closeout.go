//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package closeout

// Sentinel error-message constants. These back `errors.New`
// values in `internal/err/closeout/` and are matched via
// `errors.Is` at the call site. They cannot use desc.Text
// because the sentinels are package-level vars evaluated
// before the embedded YAML lookup is populated; wrapping
// format strings live in commands/text/errors.yaml instead.
const (
	// ErrMsgMissingFrontmatter signals a closeout file missing
	// the `---` open delimiter on line 1.
	ErrMsgMissingFrontmatter = "closeout missing frontmatter"
	// ErrMsgMissingFields signals a closeout frontmatter
	// missing one of the required fields (sha, branch, mode,
	// generated-at).
	ErrMsgMissingFields = "closeout frontmatter missing required fields"
	// ErrMsgModeRequired signals an empty Mode supplied to
	// Write. Every closeout must declare a mode.
	ErrMsgModeRequired = "closeout mode is required"
)

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
