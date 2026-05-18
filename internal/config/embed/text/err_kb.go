//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for `ctx kb` CLI error wrappers.
const (
	// DescKeyErrKbGroundingMissing wraps a missing
	// grounding-sources.md error.
	DescKeyErrKbGroundingMissing = "err.kb.grounding-missing"
	// DescKeyErrKbGroundingEmpty wraps an empty
	// grounding-sources.md error.
	DescKeyErrKbGroundingEmpty = "err.kb.grounding-empty"
	// DescKeyErrKbTopicExists wraps a topic-already-exists
	// refusal.
	DescKeyErrKbTopicExists = "err.kb.topic-exists"
	// DescKeyErrKbMkdirIngest wraps `os.MkdirAll` for the
	// ingest directory.
	DescKeyErrKbMkdirIngest = "err.kb.mkdir-ingest"
	// DescKeyErrKbOpenFindings wraps `os.OpenFile` for the
	// findings log.
	DescKeyErrKbOpenFindings = "err.kb.open-findings"
	// DescKeyErrKbWriteFinding wraps a write to the findings
	// log.
	DescKeyErrKbWriteFinding = "err.kb.write-finding"
	// DescKeyErrKbReadKBIndex wraps `os.ReadFile` for the kb
	// landing page during reindex.
	DescKeyErrKbReadKBIndex = "err.kb.read-kb-index"
	// DescKeyErrKbWriteKBIndex wraps `os.WriteFile` for the
	// kb landing page during reindex.
	DescKeyErrKbWriteKBIndex = "err.kb.write-kb-index"
	// DescKeyErrKbReadTopicsDir wraps `os.ReadDir` of the
	// topics directory during reindex.
	DescKeyErrKbReadTopicsDir = "err.kb.read-topics-dir"
	// DescKeyErrKbMkdirTopic wraps `os.MkdirAll` for a new
	// topic directory.
	DescKeyErrKbMkdirTopic = "err.kb.mkdir-topic"
	// DescKeyErrKbReadTopicTemplate wraps `fs.ReadFile` for
	// the embedded topic template.
	DescKeyErrKbReadTopicTemplate = "err.kb.read-topic-template"
	// DescKeyErrKbWriteTopicIndex wraps `os.WriteFile` for
	// the topic index.md.
	DescKeyErrKbWriteTopicIndex = "err.kb.write-topic-index"
	// DescKeyErrKbAskNoQuestion is the text key for the
	// empty-question-arg sentinel.
	DescKeyErrKbAskNoQuestion = "err.kb.ask-no-question"
	// DescKeyErrKbIngestNoSources is the text key for the
	// empty-sources-arg sentinel.
	DescKeyErrKbIngestNoSources = "err.kb.ingest-no-sources"
	// DescKeyErrKbNoteNoText is the text key for the
	// empty-note-arg sentinel.
	DescKeyErrKbNoteNoText = "err.kb.note-no-text"
	// DescKeyErrKbTopicEmptyName is the text key for the
	// empty-slug topic-new sentinel.
	DescKeyErrKbTopicEmptyName = "err.kb.topic-empty-name"
	// DescKeyErrKbReindexMissingBlock is the text key for
	// the missing-CTX:KB:TOPICS-block sentinel.
	DescKeyErrKbReindexMissingBlock = "err.kb.reindex-missing-block"
)
