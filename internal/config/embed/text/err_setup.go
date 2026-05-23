//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for setup operations errors.
const (
	// DescKeyErrSetupCreateDir is the text key for err setup create dir messages.
	DescKeyErrSetupCreateDir = "err.setup.create-dir"
	// DescKeyErrSetupMarshalConfig is the text key for err setup marshal config
	// messages.
	DescKeyErrSetupMarshalConfig = "err.setup.marshal-config"
	// DescKeyErrSetupFileWrite is the text key for err setup file write messages.
	DescKeyErrSetupFileWrite = "err.setup.write-file"
	// DescKeyErrSetupSyncSteering is the text key for err setup sync steering
	// messages.
	DescKeyErrSetupSyncSteering = "err.setup.sync-steering"
	// DescKeyErrSetupMissingEmbeddedAsset is the text key for the
	// "embedded asset missing" setup error.
	DescKeyErrSetupMissingEmbeddedAsset = "err.setup.missing-embedded-asset"
	// DescKeyErrSetupMissingToolOrBackend is the sentinel
	// message key for ErrMissingToolOrBackend.
	DescKeyErrSetupMissingToolOrBackend = "err.setup.missing-tool-or-backend"
	// DescKeyErrSetupBackendAndToolConflict is the sentinel
	// message key for ErrBackendAndToolConflict.
	DescKeyErrSetupBackendAndToolConflict = "err.setup.backend-and-tool-conflict"
	// DescKeyErrSetupBackendNameRequired is the sentinel
	// message key for ErrBackendNameRequired.
	DescKeyErrSetupBackendNameRequired = "err.setup.backend-name-required"
	// DescKeyErrSetupReadCtxrc is the wrapper format key
	// for ReadCtxrc.
	DescKeyErrSetupReadCtxrc = "err.setup.read-ctxrc"
	// DescKeyErrSetupParseCtxrc is the wrapper format key
	// for ParseCtxrc.
	DescKeyErrSetupParseCtxrc = "err.setup.parse-ctxrc"
	// DescKeyErrSetupMarshalCtxrc is the wrapper format
	// key for MarshalCtxrc.
	DescKeyErrSetupMarshalCtxrc = "err.setup.marshal-ctxrc"
)
