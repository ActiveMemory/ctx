//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	errAudit "github.com/ActiveMemory/ctx/internal/ctxctl/err/audit"
)

// Audit-channel error text. ctxctl owns its user-facing voice as
// plain English Go constants; internal/ctxctl/err/audit returns
// typed, text-free errors that these templates render. The values
// were lifted verbatim from the former err/audit constants so the
// migration preserves output exactly. (%w became %v: these build
// display strings via Sprintf, not wrapped errors.) See
// specs/ctxctl-bootstrap.md and DECISIONS.md (2026-05-27).
const (
	fmtReadReport     = "read audit report %s: %v"
	fmtParseReport    = "parse audit report %s: %v"
	fmtWriteDismissal = "write audit dismissal ledger: %v"
	fmtReadDismissal  = "read audit dismissal ledger: %v"
	fmtUnknownID      = "unknown audit id: %s"
	msgIDRequired     = "audit id required " +
		"(use --all to dismiss every report)"
	msgNoFrontmatter = "audit report missing yaml frontmatter"
	msgUnterminated  = "audit report frontmatter is " +
		"not terminated by a closing delimiter"
)

// errPrefix mirrors cobra's default "Error:" prefix so ctxctl's own
// printer produces output identical to the silenced cobra path.
const errPrefix = "Error:"

// printErr is ctxctl's sole error printer (analogous to ctx's
// internal/write/err.With). The root command silences cobra's own
// error path so this renders audit typed errors in ctxctl's English
// voice before display.
//
// Parameters:
//   - cmd: cobra command whose stderr stream receives the message;
//     nil cmd or nil err is a no-op
//   - err: the error to render and print after the "Error:" prefix
func printErr(cmd *cobra.Command, err error) {
	if cmd == nil || err == nil {
		return
	}
	cmd.PrintErrln(errPrefix, auditErrText(err))
}

// auditErrText renders an audit-channel error in ctxctl's
// user-facing English. Errors ctxctl does not own (e.g. the shared
// context-dir errors from the ctx module) fall through to their own
// message.
//
// Parameters:
//   - err: the error to render
//
// Returns:
//   - string: the user-facing message
func auditErrText(err error) string {
	if e, ok := errors.AsType[*errAudit.ReadReportError](err); ok {
		return fmt.Sprintf(fmtReadReport, e.Name, e.Cause)
	}
	if e, ok := errors.AsType[*errAudit.ParseReportError](err); ok {
		return fmt.Sprintf(fmtParseReport, e.Name, auditCause(e.Cause))
	}
	if e, ok := errors.AsType[*errAudit.WriteDismissalError](err); ok {
		return fmt.Sprintf(fmtWriteDismissal, e.Cause)
	}
	if e, ok := errors.AsType[*errAudit.ReadDismissalError](err); ok {
		return fmt.Sprintf(fmtReadDismissal, e.Cause)
	}
	if e, ok := errors.AsType[*errAudit.UnknownIDError](err); ok {
		return fmt.Sprintf(fmtUnknownID, e.ID)
	}
	if msg, ok := sentinelText(err); ok {
		return msg
	}
	return err.Error()
}

// auditCause renders a wrapped cause: a known audit sentinel in its
// user-facing English, otherwise the cause's own message. Used so a
// parse failure reads "parse audit report X: <sentinel prose>"
// rather than leaking the sentinel's diagnostic code.
//
// Parameters:
//   - cause: the wrapped cause error
//
// Returns:
//   - string: the user-facing cause text
func auditCause(cause error) string {
	if msg, ok := sentinelText(cause); ok {
		return msg
	}
	return cause.Error()
}

// sentinelText maps an audit sentinel error to its user-facing
// English.
//
// Parameters:
//   - err: the error to test
//
// Returns:
//   - string: the user-facing message (empty when not a sentinel)
//   - bool: true when err matches a known audit sentinel
func sentinelText(err error) (string, bool) {
	switch {
	case errors.Is(err, errAudit.ErrIDRequired):
		return msgIDRequired, true
	case errors.Is(err, errAudit.ErrNoFrontmatter):
		return msgNoFrontmatter, true
	case errors.Is(err, errAudit.ErrUnterminatedFrontmatter):
		return msgUnterminated, true
	}
	return "", false
}
