//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for closeout error wrappers. The matching YAML
// entries live in commands/text/errors.yaml; constructors in
// internal/err/closeout/ resolve them via desc.Text at error
// construction time.
const (
	// DescKeyErrCloseoutReadCloseout is the text key for the
	// closeout-file read failure wrapper.
	DescKeyErrCloseoutReadCloseout = "err.closeout.read-closeout"
	// DescKeyErrCloseoutParseFrontmatter is the text key for the
	// frontmatter parse failure wrapper.
	DescKeyErrCloseoutParseFrontmatter = "err.closeout.parse-frontmatter"
	// DescKeyErrCloseoutMarshalFrontmatter is the text key for
	// the frontmatter marshal failure wrapper.
	DescKeyErrCloseoutMarshalFrontmatter = "err.closeout.marshal-frontmatter"
	// DescKeyErrCloseoutReadCloseoutsDir is the text key for the
	// closeouts-directory enumeration failure wrapper.
	DescKeyErrCloseoutReadCloseoutsDir = "err.closeout.read-closeouts-dir"
	// DescKeyErrCloseoutResolveHead is the text key for the
	// git-head resolution failure wrapper.
	DescKeyErrCloseoutResolveHead = "err.closeout.resolve-head"
	// DescKeyErrCloseoutMkdirCloseouts is the text key for the
	// closeouts-dir mkdir failure wrapper.
	DescKeyErrCloseoutMkdirCloseouts = "err.closeout.mkdir-closeouts"
	// DescKeyErrCloseoutWriteCloseout is the text key for the
	// closeout-file write failure wrapper.
	DescKeyErrCloseoutWriteCloseout = "err.closeout.write-closeout"
	// DescKeyErrCloseoutMkdirArchive is the text key for the
	// archive-dir mkdir failure wrapper.
	DescKeyErrCloseoutMkdirArchive = "err.closeout.mkdir-archive"
	// DescKeyErrCloseoutArchiveMove is the text key for the
	// single-closeout archive-move failure wrapper.
	DescKeyErrCloseoutArchiveMove = "err.closeout.archive-move"
	// DescKeyErrCloseoutMissingFields is the text key for the
	// missing-fields sentinel-wrap with the joined field list.
	DescKeyErrCloseoutMissingFields = "err.closeout.missing-fields"
	// DescKeyErrCloseoutMissingFieldsMsg is the text key for the
	// missing-fields sentinel's own `.Error()` string (the prefix
	// that the wrapper format above interpolates via `%w`).
	DescKeyErrCloseoutMissingFieldsMsg = "err.closeout.missing-fields-msg"
	// DescKeyErrCloseoutMissingFrontmatter is the text key for the
	// missing-frontmatter parse sentinel.
	DescKeyErrCloseoutMissingFrontmatter = "err.closeout.missing-frontmatter"
	// DescKeyErrCloseoutModeRequired is the text key for the
	// empty-mode sentinel returned by Write.
	DescKeyErrCloseoutModeRequired = "err.closeout.mode-required"
)
