//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package closeout

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

const (
	// ErrMissingFrontmatter signals a closeout file missing the
	// `---` open delimiter on line 1.
	ErrMissingFrontmatter = entity.Sentinel(
		text.DescKeyErrCloseoutMissingFrontmatter,
	)
	// ErrMissingFields signals a closeout frontmatter missing
	// one of the required fields (sha, branch, mode,
	// generated-at). Constructor [MissingFields] wraps it with
	// the actual field names.
	ErrMissingFields = entity.Sentinel(
		text.DescKeyErrCloseoutMissingFieldsMsg,
	)
	// ErrModeRequired signals a
	// [github.com/ActiveMemory/ctx/internal/write/closeout.Write]
	// call with an empty mode string.
	ErrModeRequired = entity.Sentinel(text.DescKeyErrCloseoutModeRequired)
)

// MissingFields wraps the sentinel [ErrMissingFields] with a
// comma-separated list of the missing field names.
//
// Parameters:
//   - missing: ordered slice of field names that were absent
//     or empty in the parsed frontmatter.
//
// Returns:
//   - error: wrapping [ErrMissingFields] for [errors.Is]
//     matches at the call site.
func MissingFields(missing []string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrCloseoutMissingFields),
		ErrMissingFields,
		strings.Join(missing, token.CommaSpace),
	)
}

// ReadFailed wraps an `os.ReadFile` failure encountered while
// reading a closeout from disk.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ReadFailed(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrCloseoutReadCloseout), cause)
}

// ParseFrontmatter wraps a `yaml.Unmarshal` failure while
// decoding the closeout's YAML header.
//
// Parameters:
//   - cause: the underlying YAML parser error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ParseFrontmatter(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrCloseoutParseFrontmatter), cause,
	)
}

// MarshalFrontmatter wraps a `yaml.Marshal` failure while
// encoding a new closeout's YAML header.
//
// Parameters:
//   - cause: the underlying YAML marshal error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func MarshalFrontmatter(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrCloseoutMarshalFrontmatter), cause,
	)
}

// ReadCloseoutsDir wraps an `os.ReadDir` failure while
// enumerating closeouts under `.context/ingest/closeouts/`.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ReadCloseoutsDir(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrCloseoutReadCloseoutsDir), cause,
	)
}

// ResolveHead wraps a
// [github.com/ActiveMemory/ctx/internal/git_meta.ResolveHead]
// failure when stamping sha / branch into new closeouts.
//
// Parameters:
//   - cause: the underlying gitmeta error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ResolveHead(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrCloseoutResolveHead), cause)
}

// MkdirCloseouts wraps `os.MkdirAll` for the closeouts
// directory.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func MkdirCloseouts(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrCloseoutMkdirCloseouts), cause)
}

// WriteFailed wraps `os.WriteFile` for the closeout file
// itself.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func WriteFailed(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrCloseoutWriteCloseout), cause)
}

// MkdirArchive wraps `os.MkdirAll` for the archive destination
// directory at `.context/archive/closeouts/`.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func MkdirArchive(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrCloseoutMkdirArchive), cause)
}

// ArchiveMove wraps `os.Rename` of a single closeout into the
// archive directory.
//
// Parameters:
//   - name: the source file's basename (operator-facing).
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped with operator-friendly prefix.
func ArchiveMove(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrCloseoutArchiveMove), name, cause,
	)
}
