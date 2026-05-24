//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// ReadReport wraps a report-file read failure.
//
// Parameters:
//   - name: report basename (e.g. "surface")
//   - cause: underlying read error
//
// Returns:
//   - error: "read audit report <name>: <cause>"
func ReadReport(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrAuditReadReport), name, cause,
	)
}

// ParseReport wraps a frontmatter / body parse failure.
//
// Parameters:
//   - name: report basename
//   - cause: underlying parse error
//
// Returns:
//   - error: "parse audit report <name>: <cause>"
func ParseReport(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrAuditParseReport), name, cause,
	)
}

// WriteDismissal wraps a dismissal-ledger write failure.
//
// Parameters:
//   - cause: underlying write error
//
// Returns:
//   - error: "write audit dismissal ledger: <cause>"
func WriteDismissal(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrAuditWriteDismissal), cause,
	)
}

// ReadDismissal wraps a dismissal-ledger read failure.
//
// Parameters:
//   - cause: underlying read error
//
// Returns:
//   - error: "read audit dismissal ledger: <cause>"
func ReadDismissal(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrAuditReadDismissal), cause,
	)
}

// UnknownID returns an error for an audit id the CLI cannot
// resolve to a report file.
//
// Parameters:
//   - id: the unknown audit id
//
// Returns:
//   - error: "unknown audit id: <id>"
func UnknownID(id string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrAuditUnknownID), id,
	)
}

// IDRequired returns the error surfaced when `ctx audit
// dismiss` is invoked with no ids and no --all.
//
// Returns:
//   - error: "audit id required ..."
func IDRequired() error {
	return errors.New(
		desc.Text(text.DescKeyErrAuditIDRequired),
	)
}

// ErrNoFrontmatter is the sentinel error returned when a
// candidate audit report does not open with the YAML
// frontmatter delimiter. Comparable via errors.Is.
var ErrNoFrontmatter = errors.New(
	desc.Text(text.DescKeyErrAuditNoFrontmatter),
)

// ErrUnterminatedFrontmatter is the sentinel error
// returned when the closing frontmatter delimiter is
// missing. Comparable via errors.Is.
var ErrUnterminatedFrontmatter = errors.New(
	desc.Text(text.DescKeyErrAuditUnterminatedFrontmatter),
)
