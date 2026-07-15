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
	// DescKeyErrKbGlossaryAppendRow wraps a row-write failure.
	DescKeyErrKbGlossaryAppendRow = "err.kb.glossary.append-row"
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
	// DescKeyErrKbTimelineAppendRow wraps a row-write failure.
	DescKeyErrKbTimelineAppendRow = "err.kb.timeline.append-row"
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
	// DescKeyErrKbSourcemapAppendRow wraps a row-write failure.
	DescKeyErrKbSourcemapAppendRow = "err.kb.sourcemap.append-row"
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
	// DescKeyErrKbRelationshipAppendRow wraps a row-write
	// failure.
	DescKeyErrKbRelationshipAppendRow = "err.kb.relationship.append-row"
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
	// DescKeyErrKbContradictionAppendRow wraps a row-write
	// failure.
	DescKeyErrKbContradictionAppendRow = "err.kb.contradiction.append-row"
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
	// DescKeyErrKbDecisionAppendRow wraps a row-write failure.
	DescKeyErrKbDecisionAppendRow = "err.kb.decision.append-row"
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
	// DescKeyErrKbQuestionAppendRow wraps a row-write failure.
	DescKeyErrKbQuestionAppendRow = "err.kb.question.append-row"
	// DescKeyErrKbQuestionParseQNumber wraps a strconv.Atoi
	// failure on a Q-### digit string.
	DescKeyErrKbQuestionParseQNumber = "err.kb.question.parse-q-number"
)
