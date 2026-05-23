//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for AI proposal queue errors. The matching
// YAML entries live in commands/text/errors.yaml under
// the `err.proposal.*` namespace; constructors in
// `internal/err/proposal/` resolve them via desc.Text at
// error construction time.
const (
	// DescKeyErrProposalMkdir is the wrapper format key
	// for the proposals-dir mkdir failure constructor.
	DescKeyErrProposalMkdir = "err.proposal.mkdir"
	// DescKeyErrProposalWrite is the wrapper format key
	// for the proposal-file write failure constructor.
	DescKeyErrProposalWrite = "err.proposal.write"
	// DescKeyErrProposalReadInput is the wrapper format
	// key for the extract-input read failure constructor.
	DescKeyErrProposalReadInput = "err.proposal.read-input"
	// DescKeyErrProposalEmptyInput is the message for
	// extract called with empty input after trimming.
	DescKeyErrProposalEmptyInput = "err.proposal.empty-input"
)
