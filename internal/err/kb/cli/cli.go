//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
)

const (
	// ErrAskNoQuestion signals an empty `ctx kb ask`
	// invocation.
	ErrAskNoQuestion = entity.Sentinel(
		text.DescKeyErrKbCliAskNoQuestion,
	)
	// ErrIngestNoSources signals an empty `ctx kb ingest`
	// invocation.
	ErrIngestNoSources = entity.Sentinel(
		text.DescKeyErrKbCliIngestNoSources,
	)
	// ErrNoteNoText signals an empty `ctx kb note` invocation.
	ErrNoteNoText = entity.Sentinel(text.DescKeyErrKbCliNoteNoText)
	// ErrTopicEmptyName signals a `ctx kb topic new`
	// invocation whose name reduces to an empty slug.
	ErrTopicEmptyName = entity.Sentinel(
		text.DescKeyErrKbCliTopicEmptyName,
	)
	// ErrReindexMissingBlock signals a kb landing page that is
	// missing the CTX:KB:TOPICS managed block.
	ErrReindexMissingBlock = entity.Sentinel(
		text.DescKeyErrKbCliReindexMissingBlock,
	)
)

// GroundingMissing wraps a missing grounding-sources.md error
// with the resolved path.
//
// Parameters:
//   - path: absolute path to the missing grounding file.
//
// Returns:
//   - error: descriptive refusal.
func GroundingMissing(path string) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbCliGroundingMissing), path)
}

// GroundingEmpty wraps an empty grounding-sources.md error
// with the resolved path.
//
// Parameters:
//   - path: absolute path to the empty grounding file.
//
// Returns:
//   - error: descriptive refusal.
func GroundingEmpty(path string) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbCliGroundingEmpty), path)
}

// TopicExists wraps a topic-already-exists refusal with the
// slug and the indexPath that would have been written.
//
// Parameters:
//   - slug: the topic slug.
//   - indexPath: the path of the existing index.md.
//
// Returns:
//   - error: descriptive refusal.
func TopicExists(slug, indexPath string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbCliTopicExists), slug, indexPath,
	)
}

// MkdirIngest wraps an ingest-dir create failure.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func MkdirIngest(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbCliMkdirIngest), cause)
}

// OpenFindings wraps a findings-file open failure.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func OpenFindings(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbCliOpenFindings), cause)
}

// WriteFinding wraps a findings-file write failure.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func WriteFinding(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbCliWriteFinding), cause)
}

// ReadKBIndex wraps a kb-index read failure during reindex.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func ReadKBIndex(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbCliReadKBIndex), cause)
}

// WriteKBIndex wraps a kb-index write failure during reindex.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func WriteKBIndex(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbCliWriteKBIndex), cause)
}

// ReadTopicsDir wraps a topics-dir read failure during reindex.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func ReadTopicsDir(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbCliReadTopicsDir), cause)
}

// MkdirTopic wraps a topic-dir create failure.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func MkdirTopic(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbCliMkdirTopic), cause)
}

// ReadTopicTemplate wraps an embedded-template read failure.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func ReadTopicTemplate(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbCliReadTopicTemplate), cause,
	)
}

// WriteTopicIndex wraps a topic-index write failure.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func WriteTopicIndex(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbCliWriteTopicIndex), cause,
	)
}
