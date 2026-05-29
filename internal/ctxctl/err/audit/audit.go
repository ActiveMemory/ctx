//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"errors"
	"fmt"
)

// Audit-channel diagnostic codes. These are stable, namespaced
// identifiers, NOT ctxctl's user-facing voice: ctxctl
// (tools/ctxctl) owns the English the user sees and renders it
// from the typed errors below. A code only surfaces if an audit
// error is printed without going through that formatter (e.g. a
// stray %v in a log); it exists because a Go error must carry a
// string. See specs/ctxctl-bootstrap.md and DECISIONS.md
// (2026-05-27).
const (
	// codeReadReport tags a report-file or directory read failure.
	codeReadReport = "ctxctl/audit: read report"
	// codeParseReport tags a frontmatter / body parse failure.
	codeParseReport = "ctxctl/audit: parse report"
	// codeWriteDismissal tags a dismissal-ledger write failure.
	codeWriteDismissal = "ctxctl/audit: write dismissal ledger"
	// codeReadDismissal tags a dismissal-ledger read failure.
	codeReadDismissal = "ctxctl/audit: read dismissal ledger"
	// codeUnknownID tags an unresolvable audit id.
	codeUnknownID = "ctxctl/audit: unknown id"
	// codeIDRequired tags a dismiss invoked with no id and no --all.
	codeIDRequired = "ctxctl/audit: id required"
	// codeNoFrontmatter tags a report opening without frontmatter.
	codeNoFrontmatter = "ctxctl/audit: report missing frontmatter"
	// codeUnterminated tags a report whose frontmatter is unclosed.
	codeUnterminated = "ctxctl/audit: report frontmatter unterminated"
)

// fmtCodeCause joins a diagnostic code with its underlying cause
// or subject for the Error() fallback string.
const fmtCodeCause = "%s: %v"

// Sentinel errors comparable via errors.Is. ctxctl maps each to a
// user-facing message; the text here is a diagnostic code only.
var (
	// ErrIDRequired is returned when `ctxctl audit dismiss` is
	// invoked with neither an id nor --all.
	ErrIDRequired = errors.New(codeIDRequired)

	// ErrNoFrontmatter is returned when a candidate audit report
	// does not open with the YAML frontmatter delimiter.
	ErrNoFrontmatter = errors.New(codeNoFrontmatter)

	// ErrUnterminatedFrontmatter is returned when the closing
	// frontmatter delimiter is missing.
	ErrUnterminatedFrontmatter = errors.New(codeUnterminated)
)

// ReadReportError reports a failure to read an audit report file
// (or the audit directory itself). ctxctl renders the user-facing
// message from Name and Cause; callers match it with
// errors.AsType[*ReadReportError].
type ReadReportError struct {
	// Name is the report basename (e.g. "surface"), or the audit
	// directory name when the directory read itself failed.
	Name string
	// Cause is the underlying read error.
	Cause error
}

// Error implements the error interface for ReadReportError.
//
// Returns:
//   - string: diagnostic code with the underlying cause
func (e *ReadReportError) Error() string {
	return fmt.Sprintf(fmtCodeCause, codeReadReport, e.Cause)
}

// Unwrap exposes the underlying cause for errors.Is / errors.As.
//
// Returns:
//   - error: the wrapped read error
func (e *ReadReportError) Unwrap() error {
	return e.Cause
}

// ReadReport returns a ReadReportError wrapping a report-file or
// audit-directory read failure.
//
// Parameters:
//   - name: report basename, or the audit directory name
//   - cause: underlying read error
//
// Returns:
//   - *ReadReportError: typed error for errors.AsType matching
func ReadReport(name string, cause error) *ReadReportError {
	return &ReadReportError{Name: name, Cause: cause}
}

// ParseReportError reports a frontmatter / body parse failure for
// a named audit report. ctxctl renders the user-facing message
// from Name and Cause.
type ParseReportError struct {
	// Name is the report basename whose parse failed.
	Name string
	// Cause is the underlying parse error: one of
	// [ErrNoFrontmatter], [ErrUnterminatedFrontmatter], or a YAML
	// unmarshal error.
	Cause error
}

// Error implements the error interface for ParseReportError.
//
// Returns:
//   - string: diagnostic code with the underlying cause
func (e *ParseReportError) Error() string {
	return fmt.Sprintf(fmtCodeCause, codeParseReport, e.Cause)
}

// Unwrap exposes the underlying cause for errors.Is / errors.As.
//
// Returns:
//   - error: the wrapped parse error
func (e *ParseReportError) Unwrap() error {
	return e.Cause
}

// ParseReport returns a ParseReportError wrapping a frontmatter or
// body parse failure.
//
// Parameters:
//   - name: report basename
//   - cause: underlying parse error
//
// Returns:
//   - *ParseReportError: typed error for errors.AsType matching
func ParseReport(name string, cause error) *ParseReportError {
	return &ParseReportError{Name: name, Cause: cause}
}

// WriteDismissalError reports a failure to persist the dismissal
// ledger. ctxctl renders the user-facing message from Cause.
type WriteDismissalError struct {
	// Cause is the underlying directory-create, marshal, or write
	// error.
	Cause error
}

// Error implements the error interface for WriteDismissalError.
//
// Returns:
//   - string: diagnostic code with the underlying cause
func (e *WriteDismissalError) Error() string {
	return fmt.Sprintf(fmtCodeCause, codeWriteDismissal, e.Cause)
}

// Unwrap exposes the underlying cause for errors.Is / errors.As.
//
// Returns:
//   - error: the wrapped write error
func (e *WriteDismissalError) Unwrap() error {
	return e.Cause
}

// WriteDismissal returns a WriteDismissalError wrapping a
// dismissal-ledger write failure.
//
// Parameters:
//   - cause: underlying write error
//
// Returns:
//   - *WriteDismissalError: typed error for errors.AsType matching
func WriteDismissal(cause error) *WriteDismissalError {
	return &WriteDismissalError{Cause: cause}
}

// ReadDismissalError reports a failure to read the dismissal
// ledger. ctxctl renders the user-facing message from Cause.
type ReadDismissalError struct {
	// Cause is the underlying read or JSON-unmarshal error.
	Cause error
}

// Error implements the error interface for ReadDismissalError.
//
// Returns:
//   - string: diagnostic code with the underlying cause
func (e *ReadDismissalError) Error() string {
	return fmt.Sprintf(fmtCodeCause, codeReadDismissal, e.Cause)
}

// Unwrap exposes the underlying cause for errors.Is / errors.As.
//
// Returns:
//   - error: the wrapped read error
func (e *ReadDismissalError) Unwrap() error {
	return e.Cause
}

// ReadDismissal returns a ReadDismissalError wrapping a
// dismissal-ledger read failure.
//
// Parameters:
//   - cause: underlying read or unmarshal error
//
// Returns:
//   - *ReadDismissalError: typed error for errors.AsType matching
func ReadDismissal(cause error) *ReadDismissalError {
	return &ReadDismissalError{Cause: cause}
}

// UnknownIDError reports an audit id that resolves to no report
// file. ctxctl renders the user-facing message from ID.
type UnknownIDError struct {
	// ID is the audit id that could not be resolved.
	ID string
}

// Error implements the error interface for UnknownIDError.
//
// Returns:
//   - string: diagnostic code with the unresolved id
func (e *UnknownIDError) Error() string {
	return fmt.Sprintf(fmtCodeCause, codeUnknownID, e.ID)
}

// UnknownID returns an UnknownIDError for an audit id the CLI
// cannot resolve to a report file.
//
// Parameters:
//   - id: the unknown audit id
//
// Returns:
//   - *UnknownIDError: typed error for errors.AsType matching
func UnknownID(id string) *UnknownIDError {
	return &UnknownIDError{ID: id}
}
