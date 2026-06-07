//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for the dream command and its subcommands.
const (
	// UseDream is the cobra Use string for the dream command.
	UseDream = "dream"
	// UseDreamReview is the cobra Use string for the dream review
	// subcommand.
	UseDreamReview = "review"
	// UseDreamAccept is the cobra Use string for the dream accept
	// subcommand (takes a proposal id argument).
	UseDreamAccept = "accept <id>"
	// UseDreamReject is the cobra Use string for the dream reject
	// subcommand (takes a proposal id argument).
	UseDreamReject = "reject <id>"
	// UseDreamAmend is the cobra Use string for the dream amend
	// subcommand (takes a proposal id argument).
	UseDreamAmend = "amend <id>"
)

// DescKeys for the dream command and its subcommands.
const (
	// DescKeyDream is the description key for the dream command.
	DescKeyDream = "dream"
	// DescKeyDreamReview is the description key for the dream review
	// subcommand.
	DescKeyDreamReview = "dream.review"
	// DescKeyDreamAccept is the description key for the dream accept
	// subcommand.
	DescKeyDreamAccept = "dream.accept"
	// DescKeyDreamReject is the description key for the dream reject
	// subcommand.
	DescKeyDreamReject = "dream.reject"
	// DescKeyDreamAmend is the description key for the dream amend
	// subcommand.
	DescKeyDreamAmend = "dream.amend"
)
