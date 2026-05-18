//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// Glossary writer DescKeys.
const (
	// DescKeyErrKbGlossaryReadFile wraps a stat/read failure.
	DescKeyErrKbGlossaryReadFile = "err.kb.glossary.read-file"
	// DescKeyErrKbGlossaryMkdirDir wraps a parent-dir mkdir
	// failure.
	DescKeyErrKbGlossaryMkdirDir = "err.kb.glossary.mkdir-dir"
	// DescKeyErrKbGlossaryOpenFile wraps an open-for-append
	// failure.
	DescKeyErrKbGlossaryOpenFile = "err.kb.glossary.open-file"
	// DescKeyErrKbGlossaryWriteRow wraps a row-write failure.
	DescKeyErrKbGlossaryWriteRow = "err.kb.glossary.write-row"
)

// Timeline writer DescKeys.
const (
	// DescKeyErrKbTimelineReadFile wraps a stat/read failure.
	DescKeyErrKbTimelineReadFile = "err.kb.timeline.read-file"
	// DescKeyErrKbTimelineMkdirDir wraps a parent-dir mkdir
	// failure.
	DescKeyErrKbTimelineMkdirDir = "err.kb.timeline.mkdir-dir"
	// DescKeyErrKbTimelineOpenFile wraps an open-for-append
	// failure.
	DescKeyErrKbTimelineOpenFile = "err.kb.timeline.open-file"
	// DescKeyErrKbTimelineWriteRow wraps a row-write failure.
	DescKeyErrKbTimelineWriteRow = "err.kb.timeline.write-row"
)

// Source-map writer DescKeys.
const (
	// DescKeyErrKbSourcemapReadFile wraps a stat/read failure.
	DescKeyErrKbSourcemapReadFile = "err.kb.sourcemap.read-file"
	// DescKeyErrKbSourcemapMkdirDir wraps a parent-dir mkdir
	// failure.
	DescKeyErrKbSourcemapMkdirDir = "err.kb.sourcemap.mkdir-dir"
	// DescKeyErrKbSourcemapOpenFile wraps an open-for-append
	// failure.
	DescKeyErrKbSourcemapOpenFile = "err.kb.sourcemap.open-file"
	// DescKeyErrKbSourcemapWriteRow wraps a row-write failure.
	DescKeyErrKbSourcemapWriteRow = "err.kb.sourcemap.write-row"
)

// Relationship-map writer DescKeys.
const (
	// DescKeyErrKbRelationshipReadFile wraps a stat/read
	// failure.
	DescKeyErrKbRelationshipReadFile = "err.kb.relationship.read-file"
	// DescKeyErrKbRelationshipMkdirDir wraps a parent-dir mkdir
	// failure.
	DescKeyErrKbRelationshipMkdirDir = "err.kb.relationship.mkdir-dir"
	// DescKeyErrKbRelationshipOpenFile wraps an open-for-append
	// failure.
	DescKeyErrKbRelationshipOpenFile = "err.kb.relationship.open-file"
	// DescKeyErrKbRelationshipWriteRow wraps a row-write
	// failure.
	DescKeyErrKbRelationshipWriteRow = "err.kb.relationship.write-row"
)

// Contradiction writer DescKeys.
const (
	// DescKeyErrKbContradictionReadFile wraps a file-read
	// failure.
	DescKeyErrKbContradictionReadFile = "err.kb.contradiction.read-file"
	// DescKeyErrKbContradictionMkdirDir wraps a parent-dir
	// mkdir failure.
	DescKeyErrKbContradictionMkdirDir = "err.kb.contradiction.mkdir-dir"
	// DescKeyErrKbContradictionOpenFile wraps an
	// open-for-append failure.
	DescKeyErrKbContradictionOpenFile = "err.kb.contradiction.open-file"
	// DescKeyErrKbContradictionWriteRow wraps a row-write
	// failure.
	DescKeyErrKbContradictionWriteRow = "err.kb.contradiction.write-row"
	// DescKeyErrKbContradictionParseCNumber wraps a
	// strconv.Atoi failure on a C-### digit string.
	DescKeyErrKbContradictionParseCNumber = "err.kb.contradiction.parse-c-number"
)

// Decision writer DescKeys.
const (
	// DescKeyErrKbDecisionReadFile wraps a file-read failure.
	DescKeyErrKbDecisionReadFile = "err.kb.decision.read-file"
	// DescKeyErrKbDecisionMkdirDir wraps a parent-dir mkdir
	// failure.
	DescKeyErrKbDecisionMkdirDir = "err.kb.decision.mkdir-dir"
	// DescKeyErrKbDecisionOpenFile wraps an open-for-append
	// failure.
	DescKeyErrKbDecisionOpenFile = "err.kb.decision.open-file"
	// DescKeyErrKbDecisionWriteRow wraps a row-write failure.
	DescKeyErrKbDecisionWriteRow = "err.kb.decision.write-row"
	// DescKeyErrKbDecisionParseDDNumber wraps a strconv.Atoi
	// failure on a DD-### digit string.
	DescKeyErrKbDecisionParseDDNumber = "err.kb.decision.parse-dd-number"
)

// Question writer DescKeys.
const (
	// DescKeyErrKbQuestionReadFile wraps a file-read failure.
	DescKeyErrKbQuestionReadFile = "err.kb.question.read-file"
	// DescKeyErrKbQuestionMkdirDir wraps a parent-dir mkdir
	// failure.
	DescKeyErrKbQuestionMkdirDir = "err.kb.question.mkdir-dir"
	// DescKeyErrKbQuestionOpenFile wraps an open-for-append
	// failure.
	DescKeyErrKbQuestionOpenFile = "err.kb.question.open-file"
	// DescKeyErrKbQuestionWriteRow wraps a row-write failure.
	DescKeyErrKbQuestionWriteRow = "err.kb.question.write-row"
	// DescKeyErrKbQuestionParseQNumber wraps a strconv.Atoi
	// failure on a Q-### digit string.
	DescKeyErrKbQuestionParseQNumber = "err.kb.question.parse-q-number"
)
