//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for validation errors.
const (
	// DescKeyErrValidateContextDirSymlink is the text key for err validate
	// context dir symlink messages.
	DescKeyErrValidateContextDirSymlink = "err.validate.context-dir-symlink"
	// DescKeyErrValidateContextFileSymlink is the text key for err validate
	// context file symlink messages.
	DescKeyErrValidateContextFileSymlink = "err.validate.context-file-symlink"
	// DescKeyErrValidateInvalidSelection is the text key for err validate invalid
	// selection messages.
	DescKeyErrValidateInvalidSelection = "err.validate.invalid-selection"
	// DescKeyErrValidateUnknownDocument is the text key for err validate unknown
	// document messages.
	DescKeyErrValidateUnknownDocument = "err.validate.unknown-document"
	// DescKeyErrValidateArgRequired is the text key for err validate arg required
	// messages.
	DescKeyErrValidateArgRequired = "err.validation.arg-required"
	// DescKeyErrValidateCtxNotInPath is the text key for err validate ctx not in
	// path messages.
	DescKeyErrValidateCtxNotInPath = "err.validation.ctx-not-in-path"
	// DescKeyErrValidateDriftViolations is the text key for err validate drift
	// violations messages.
	DescKeyErrValidateDriftViolations = "err.validation.drift-violations"
	// DescKeyErrValidateFlagRequired is the text key for err validate flag
	// required messages.
	DescKeyErrValidateFlagRequired = "err.validation.flag-required"
	// DescKeyErrValidateParseFile is the text key for err validate parse file
	// messages.
	DescKeyErrValidateParseFile = "err.validation.parse-file"
	// DescKeyErrValidateWorkingDirectory is the text key for err validate working
	// directory messages.
	DescKeyErrValidateWorkingDirectory = "err.validation.working-directory"
	// DescKeyErrValidateFlagEmpty is the text key for an empty or
	// whitespace-only required body flag.
	DescKeyErrValidateFlagEmpty = "err.validation.flag-empty"
	// DescKeyErrValidateFlagPlaceholder is the text key for a
	// placeholder-valued required body flag.
	DescKeyErrValidateFlagPlaceholder = "err.validation.flag-placeholder"
	// DescKeyErrValidateFlagUnregistered is the text key for a
	// programming-error case where RequireBodyFlags is asked to
	// require a flag that does not exist on the command.
	DescKeyErrValidateFlagUnregistered = "err.validation.flag-unregistered"
	// DescKeyErrValidateMarkRequired is the text key for a failure
	// when calling cobra.Command.MarkFlagRequired.
	DescKeyErrValidateMarkRequired = "err.validation.mark-required"
)
