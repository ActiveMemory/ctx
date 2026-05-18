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
		text.DescKeyErrKbAskNoQuestion,
	)
	// ErrIngestNoSources signals an empty `ctx kb ingest`
	// invocation.
	ErrIngestNoSources = entity.Sentinel(
		text.DescKeyErrKbIngestNoSources,
	)
	// ErrNoteNoText signals an empty `ctx kb note` invocation.
	ErrNoteNoText = entity.Sentinel(text.DescKeyErrKbNoteNoText)
	// ErrTopicEmptyName signals a `ctx kb topic new`
	// invocation whose name reduces to an empty slug.
	ErrTopicEmptyName = entity.Sentinel(
		text.DescKeyErrKbTopicEmptyName,
	)
	// ErrReindexMissingBlock signals a kb landing page that is
	// missing the CTX:KB:TOPICS managed block.
	ErrReindexMissingBlock = entity.Sentinel(
		text.DescKeyErrKbReindexMissingBlock,
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
	return fmt.Errorf(desc.Text(text.DescKeyErrKbGroundingMissing), path)
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
	return fmt.Errorf(desc.Text(text.DescKeyErrKbGroundingEmpty), path)
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
		desc.Text(text.DescKeyErrKbTopicExists), slug, indexPath,
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
	return fmt.Errorf(desc.Text(text.DescKeyErrKbMkdirIngest), cause)
}

// OpenFindings wraps a findings-file open failure.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func OpenFindings(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbOpenFindings), cause)
}

// WriteFinding wraps a findings-file write failure.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func WriteFinding(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbWriteFinding), cause)
}

// ReadKBIndex wraps a kb-index read failure during reindex.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func ReadKBIndex(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbReadKBIndex), cause)
}

// WriteKBIndex wraps a kb-index write failure during reindex.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func WriteKBIndex(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbWriteKBIndex), cause)
}

// ReadTopicsDir wraps a topics-dir read failure during reindex.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func ReadTopicsDir(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbReadTopicsDir), cause)
}

// MkdirTopic wraps a topic-dir create failure.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func MkdirTopic(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbMkdirTopic), cause)
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
		desc.Text(text.DescKeyErrKbReadTopicTemplate), cause,
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
		desc.Text(text.DescKeyErrKbWriteTopicIndex), cause,
	)
}
