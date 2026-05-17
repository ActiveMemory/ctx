//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for handover CLI output strings. The matching
// YAML entries live in commands/text/write.yaml; the
// `ctx handover write` command resolves them via desc.Text
// at output time.
const (
	// DescKeyWriteHandoverWrote is the text key for the
	// "wrote <path>" success line.
	DescKeyWriteHandoverWrote = "write.handover.wrote"
	// DescKeyWriteHandoverFolded is the text key for the
	// "folded N closeout(s); archived to <path>" line.
	DescKeyWriteHandoverFolded = "write.handover.folded"
	// DescKeyWriteHandoverMalformedWarning is the text key for
	// the stderr warning block opener when malformed closeouts
	// were skipped during fold.
	DescKeyWriteHandoverMalformedWarning = "write.handover.malformed-warning"
	// DescKeyWriteHandoverMalformedLine is the text key for one
	// malformed-closeout filename line inside the warning block.
	DescKeyWriteHandoverMalformedLine = "write.handover.malformed-line"
)
