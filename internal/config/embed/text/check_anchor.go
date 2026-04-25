//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for the check-anchor-drift hook (single-source-anchor
// model: specs/single-source-context-anchor.md).
const (
	// DescKeyCheckAnchorDriftBoxTitle is the text key for the
	// anchor-drift nudge box title.
	DescKeyCheckAnchorDriftBoxTitle = "check-anchor-drift.box-title"
	// DescKeyCheckAnchorDriftContent is the text key for the
	// anchor-drift nudge body. Two %s placeholders: the inherited
	// CTX_DIR and the Claude-injected CTX_DIR.
	DescKeyCheckAnchorDriftContent = "check-anchor-drift.content"
	// DescKeyCheckAnchorDriftRelayMessage is the text key for the
	// short relay-channel message.
	DescKeyCheckAnchorDriftRelayMessage = "check-anchor-drift.relay-message"
	// DescKeyCheckAnchorDriftRelayPrefix is the text key for the
	// VERBATIM-relay prefix line.
	DescKeyCheckAnchorDriftRelayPrefix = "check-anchor-drift.relay-prefix"
)
