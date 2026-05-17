//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handover

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgHandover "github.com/ActiveMemory/ctx/internal/config/handover"
)

// ErrTitleRequired signals an empty Title supplied to
// [github.com/ActiveMemory/ctx/internal/write/handover.Write].
var ErrTitleRequired = errors.New(cfgHandover.ErrMsgTitleRequired)

// ErrSummaryRequired signals an empty Summary supplied to
// [github.com/ActiveMemory/ctx/internal/write/handover.Write].
var ErrSummaryRequired = errors.New(cfgHandover.ErrMsgSummaryRequired)

// ErrNextRequired signals an empty Next supplied to
// [github.com/ActiveMemory/ctx/internal/write/handover.Write].
var ErrNextRequired = errors.New(cfgHandover.ErrMsgNextRequired)

// ErrMissingFrontmatter signals a handover file that does not
// open with `---`.
var ErrMissingFrontmatter = errors.New(cfgHandover.ErrMsgMissingFrontmatter)

// ErrMissingClosingDelim signals a handover whose frontmatter
// is never closed by a second `---`.
var ErrMissingClosingDelim = errors.New(cfgHandover.ErrMsgMissingClosingDelim)

// ErrMissingGeneratedAt signals a handover whose frontmatter
// parsed but has no generated-at value.
var ErrMissingGeneratedAt = errors.New(cfgHandover.ErrMsgMissingGeneratedAt)

// Latest wraps a failure encountered while reading the
// latest handover during fold.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func Latest(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrHandoverLatest), cause)
}

// ListCloseouts wraps a closeout-listing failure encountered
// during fold.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func ListCloseouts(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrHandoverListCloseouts), cause)
}

// MarshalFrontmatter wraps a `yaml.Marshal` failure while
// encoding a new handover's YAML header.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func MarshalFrontmatter(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHandoverMarshalFrontmatter), cause,
	)
}

// MkdirHandovers wraps an `os.MkdirAll` failure for the
// handovers directory.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func MkdirHandovers(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHandoverMkdirHandovers), cause,
	)
}

// WriteFailed wraps an `os.WriteFile` failure for the new
// handover file.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func WriteFailed(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHandoverWriteHandover), cause,
	)
}

// ArchiveFoldedCloseouts wraps the closeout archival pass
// following a fold.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func ArchiveFoldedCloseouts(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHandoverArchiveFoldedCloseouts), cause,
	)
}

// ReadFailed wraps an `os.ReadFile` failure while loading a
// handover from disk.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func ReadFailed(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHandoverReadHandover), cause,
	)
}

// ReadHandoversDir wraps an `os.ReadDir` failure while
// enumerating handovers.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func ReadHandoversDir(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHandoverReadHandoversDir), cause,
	)
}

// ParseFrontmatter wraps a `yaml.Unmarshal` failure while
// parsing the handover YAML header.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func ParseFrontmatter(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHandoverParseFrontmatter), cause,
	)
}

// ResolveHead wraps a
// [github.com/ActiveMemory/ctx/internal/git_meta.ResolveHead]
// failure when stamping sha / branch into new handovers.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func ResolveHead(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrHandoverResolveHead), cause,
	)
}
