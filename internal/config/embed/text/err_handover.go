//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for handover error wrappers. The matching YAML
// entries live in commands/text/errors.yaml; constructors in
// internal/err/handover/ resolve them via desc.Text at error
// construction time.
const (
	// DescKeyErrHandoverLatest is the text key for the
	// latest-handover lookup failure wrapper.
	DescKeyErrHandoverLatest = "err.handover.latest"
	// DescKeyErrHandoverListCloseouts is the text key for the
	// fold-time closeout-listing failure wrapper.
	DescKeyErrHandoverListCloseouts = "err.handover.list-closeouts"
	// DescKeyErrHandoverMarshalFrontmatter is the text key for
	// the frontmatter marshal failure wrapper.
	DescKeyErrHandoverMarshalFrontmatter = "err.handover.marshal-frontmatter"
	// DescKeyErrHandoverMkdirHandovers is the text key for the
	// handovers-dir mkdir failure wrapper.
	DescKeyErrHandoverMkdirHandovers = "err.handover.mkdir-handovers"
	// DescKeyErrHandoverWriteHandover is the text key for the
	// new-handover file write failure wrapper.
	DescKeyErrHandoverWriteHandover = "err.handover.write-handover"
	// DescKeyErrHandoverArchiveFoldedCloseouts is the text key
	// for the post-fold archive failure wrapper.
	DescKeyErrHandoverArchiveFoldedCloseouts = "err.handover.archive-folded-closeouts"
	// DescKeyErrHandoverReadHandover is the text key for the
	// handover-from-disk read failure wrapper.
	DescKeyErrHandoverReadHandover = "err.handover.read-handover"
	// DescKeyErrHandoverReadHandoversDir is the text key for
	// the handovers-directory enumeration failure wrapper.
	DescKeyErrHandoverReadHandoversDir = "err.handover.read-handovers-dir"
	// DescKeyErrHandoverParseFrontmatter is the text key for
	// the frontmatter parse failure wrapper.
	DescKeyErrHandoverParseFrontmatter = "err.handover.parse-frontmatter"
	// DescKeyErrHandoverResolveHead is the text key for the
	// git-head resolution failure wrapper.
	DescKeyErrHandoverResolveHead = "err.handover.resolve-head"
)
