//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package evidence

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
)

const (
	// ErrDuplicateID signals that Append was called with an
	// explicit row.ID already present in the file.
	// Renumbering is forbidden; callers must reuse the existing
	// row verbatim or mint a new ID by leaving row.ID empty.
	ErrDuplicateID = entity.Sentinel(
		text.DescKeyErrKbEvidenceDuplicateIDMsg,
	)
	// ErrInvalidBand signals a row whose Confidence is not one
	// of the four canonical bands defined in
	// [github.com/ActiveMemory/ctx/internal/config/kb].
	ErrInvalidBand = entity.Sentinel(
		text.DescKeyErrKbEvidenceInvalidBandMsg,
	)
)

// DuplicateID wraps ErrDuplicateID with the offending
// identifier so callers see exactly which ID collided.
//
// Parameters:
//   - id: the EV-### identifier that already existed.
//
// Returns:
//   - error: wraps [ErrDuplicateID] for [errors.Is] matches.
func DuplicateID(id string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbEvidenceDuplicateID),
		ErrDuplicateID, id,
	)
}

// InvalidBand wraps ErrInvalidBand with the offending band
// string.
//
// Parameters:
//   - band: the band string supplied by the caller.
//
// Returns:
//   - error: wraps [ErrInvalidBand] for [errors.Is] matches.
func InvalidBand(band string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbEvidenceInvalidBand),
		ErrInvalidBand, band,
	)
}

// ParseEVNumber wraps a strconv.Atoi failure while parsing the
// numeric portion of an EV-### token.
//
// Parameters:
//   - digits: the raw digit string that failed to parse.
//   - cause: the underlying parse error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ParseEVNumber(digits string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrKbEvidenceParseIDNumber),
		digits, cause,
	)
}

// ReadIndex wraps a file-read failure on the evidence-index
// file.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ReadIndex(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbEvidenceReadIndex), cause)
}

// MkdirDir wraps a directory-create failure on the
// evidence-index parent directory.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func MkdirDir(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbEvidenceMkdirDir), cause)
}

// OpenIndex wraps an open-for-append failure on the
// evidence-index file.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func OpenIndex(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbEvidenceOpenIndex), cause)
}

// WriteRow wraps a row-write failure to the evidence-index
// file.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func WriteRow(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrKbEvidenceAppendRow), cause)
}
