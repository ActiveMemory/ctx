//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for steering operations errors.
const (
	DescKeyErrSteeringComputeRelPath      = "err.steering.compute-rel-path"
	DescKeyErrSteeringContextDirMissing   = "err.steering.context-dir-missing"
	DescKeyErrSteeringCreateDir           = "err.steering.create-dir"
	DescKeyErrSteeringFileExists          = "err.steering.file-exists"
	DescKeyErrSteeringInvalidYAML         = "err.steering.invalid-yaml"
	DescKeyErrSteeringMissingClosingDelim = "err.steering.missing-closing-delimiter"
	DescKeyErrSteeringMissingOpeningDelim = "err.steering.missing-opening-delimiter"
	DescKeyErrSteeringNoTool              = "err.steering.no-tool"
	DescKeyErrSteeringOutputEscapesRoot   = "err.steering.output-escapes-root"
	DescKeyErrSteeringParse               = "err.steering.parse"
	DescKeyErrSteeringReadDir             = "err.steering.read-dir"
	DescKeyErrSteeringReadFile            = "err.steering.read-file"
	DescKeyErrSteeringResolveOutput       = "err.steering.resolve-output"
	DescKeyErrSteeringResolveRoot         = "err.steering.resolve-root"
	DescKeyErrSteeringSyncAll             = "err.steering.sync-all"
	DescKeyErrSteeringSyncName            = "err.steering.sync-name"
	DescKeyErrSteeringUnsupportedTool     = "err.steering.unsupported-tool"
	DescKeyErrSteeringWriteFile           = "err.steering.write-file"
	DescKeyErrSteeringWriteSteeringFile   = "err.steering.write-steering-file"
	DescKeyErrSteeringWriteInitFile       = "err.steering.write-init-file"
)
