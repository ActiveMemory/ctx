//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for kb evidence-index error wrappers.
const (
	// DescKeyErrKbEvidenceDuplicateID wraps the duplicate-EV-###
	// sentinel with the offending ID.
	DescKeyErrKbEvidenceDuplicateID = "err.kb.evidence.duplicate-id"
	// DescKeyErrKbEvidenceInvalidBand wraps the invalid-confidence
	// sentinel with the offending band string.
	DescKeyErrKbEvidenceInvalidBand = "err.kb.evidence.invalid-band"
	// DescKeyErrKbEvidenceReadIndex wraps an evidence-index read
	// failure.
	DescKeyErrKbEvidenceReadIndex = "err.kb.evidence.read-index"
	// DescKeyErrKbEvidenceMkdirDir wraps a parent-dir mkdir
	// failure for the evidence-index.
	DescKeyErrKbEvidenceMkdirDir = "err.kb.evidence.mkdir-dir"
	// DescKeyErrKbEvidenceOpenIndex wraps an open-for-append
	// failure on the evidence-index.
	DescKeyErrKbEvidenceOpenIndex = "err.kb.evidence.open-index"
	// DescKeyErrKbEvidenceWriteRow wraps a row-write failure on
	// the evidence-index.
	DescKeyErrKbEvidenceWriteRow = "err.kb.evidence.write-row"
	// DescKeyErrKbEvidenceParseIDNumber wraps a strconv.Atoi
	// failure parsing an EV-### number.
	DescKeyErrKbEvidenceParseIDNumber = "err.kb.evidence.parse-id-number"
)
