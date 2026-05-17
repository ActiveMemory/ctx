//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// KBEvidenceID matches EV-NNN identifiers in evidence-index
// files. The numeric portion is captured for high-water-mark
// allocation in
// [github.com/ActiveMemory/ctx/internal/write/kb/evidence].
// Three or more digits are tolerated for legacy hand-written
// entries; the writer mints exactly three.
var KBEvidenceID = regexp.MustCompile(`\bEV-(\d{3,})\b`)

// KBContradictionID matches C-NNN identifiers in the
// contradictions ledger. The numeric portion is captured
// for high-water-mark allocation in
// [github.com/ActiveMemory/ctx/internal/write/kb/contradiction].
var KBContradictionID = regexp.MustCompile(`\bC-(\d{3,})\b`)

// KBQuestionID matches Q-NNN identifiers in the outstanding
// questions ledger. The numeric portion is captured for
// high-water-mark allocation in
// [github.com/ActiveMemory/ctx/internal/write/kb/question].
var KBQuestionID = regexp.MustCompile(`\bQ-(\d{3,})\b`)

// KBDecisionID matches DD-NNN identifiers in the domain
// decisions ledger. The numeric portion is captured for
// high-water-mark allocation in
// [github.com/ActiveMemory/ctx/internal/write/kb/decision].
var KBDecisionID = regexp.MustCompile(`\bDD-(\d{3,})\b`)
