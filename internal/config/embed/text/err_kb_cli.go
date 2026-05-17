//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for `ctx kb` CLI error wrappers.
const (
	// DescKeyErrKbCliGroundingMissing wraps a missing
	// grounding-sources.md error.
	DescKeyErrKbCliGroundingMissing = "err.kb.cli.grounding-missing"
	// DescKeyErrKbCliGroundingEmpty wraps an empty
	// grounding-sources.md error.
	DescKeyErrKbCliGroundingEmpty = "err.kb.cli.grounding-empty"
	// DescKeyErrKbCliTopicExists wraps a topic-already-exists
	// refusal.
	DescKeyErrKbCliTopicExists = "err.kb.cli.topic-exists"
	// DescKeyErrKbCliMkdirIngest wraps `os.MkdirAll` for the
	// ingest directory.
	DescKeyErrKbCliMkdirIngest = "err.kb.cli.mkdir-ingest"
	// DescKeyErrKbCliOpenFindings wraps `os.OpenFile` for the
	// findings log.
	DescKeyErrKbCliOpenFindings = "err.kb.cli.open-findings"
	// DescKeyErrKbCliWriteFinding wraps a write to the findings
	// log.
	DescKeyErrKbCliWriteFinding = "err.kb.cli.write-finding"
	// DescKeyErrKbCliReadKBIndex wraps `os.ReadFile` for the kb
	// landing page during reindex.
	DescKeyErrKbCliReadKBIndex = "err.kb.cli.read-kb-index"
	// DescKeyErrKbCliWriteKBIndex wraps `os.WriteFile` for the
	// kb landing page during reindex.
	DescKeyErrKbCliWriteKBIndex = "err.kb.cli.write-kb-index"
	// DescKeyErrKbCliReadTopicsDir wraps `os.ReadDir` of the
	// topics directory during reindex.
	DescKeyErrKbCliReadTopicsDir = "err.kb.cli.read-topics-dir"
	// DescKeyErrKbCliMkdirTopic wraps `os.MkdirAll` for a new
	// topic directory.
	DescKeyErrKbCliMkdirTopic = "err.kb.cli.mkdir-topic"
	// DescKeyErrKbCliReadTopicTemplate wraps `fs.ReadFile` for
	// the embedded topic template.
	DescKeyErrKbCliReadTopicTemplate = "err.kb.cli.read-topic-template"
	// DescKeyErrKbCliWriteTopicIndex wraps `os.WriteFile` for
	// the topic index.md.
	DescKeyErrKbCliWriteTopicIndex = "err.kb.cli.write-topic-index"
)
