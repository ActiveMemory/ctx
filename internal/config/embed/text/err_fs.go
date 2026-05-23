//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for filesystem operations errors.
const (
	// DescKeyErrFsCreateDir is the text key for err fs create dir messages.
	DescKeyErrFsCreateDir = "err.fs.create-dir"
	// DescKeyErrFsDirNotFound is the text key for err fs dir not found messages.
	DescKeyErrFsDirNotFound = "err.fs.dir-not-found"
	// DescKeyErrFsFileAmend is the text key for err fs file amend messages.
	DescKeyErrFsFileAmend = "err.fs.file-amend"
	// DescKeyErrFsFileRead is the text key for err fs file read messages.
	DescKeyErrFsFileRead = "err.fs.file-read"
	// DescKeyErrFsFileUpdate is the text key for err fs file update messages.
	DescKeyErrFsFileUpdate = "err.fs.file-update"
	// DescKeyErrFsFileWrite is the text key for err fs file write messages.
	DescKeyErrFsFileWrite = "err.fs.file-write"
	// DescKeyErrFsMkdir is the text key for err fs mkdir messages.
	DescKeyErrFsMkdir = "err.fs.mkdir"
	// DescKeyErrFsNoInput is the text key for err fs no input messages.
	DescKeyErrFsNoInput = "err.fs.no-input"
	// DescKeyErrFsNotDirectory is the text key for err fs not directory messages.
	DescKeyErrFsNotDirectory = "err.fs.not-directory"
	// DescKeyErrFsOpenFile is the text key for err fs open file messages.
	DescKeyErrFsOpenFile = "err.fs.open-file"
	// DescKeyErrFsPathEscapesBase is the text key for err fs path escapes base
	// messages.
	DescKeyErrFsPathEscapesBase = "err.fs.path-escapes-base"
	// DescKeyErrFsReadDir is the text key for err fs read dir messages.
	DescKeyErrFsReadDir = "err.fs.read-dir"
	// DescKeyErrFsReadDirectory is the text key for err fs read directory
	// messages.
	DescKeyErrFsReadDirectory = "err.fs.read-directory"
	// DescKeyErrFsReadFile is the text key for err fs read file messages.
	DescKeyErrFsReadFile = "err.fs.read-file"
	// DescKeyErrFsReadInput is the text key for err fs read input messages.
	DescKeyErrFsReadInput = "err.fs.read-input"
	// DescKeyErrFsReadInputStream is the text key for err fs read input stream
	// messages.
	DescKeyErrFsReadInputStream = "err.fs.read-input-stream"
	// DescKeyErrFsRefuseSystemPath is the text key for err fs refuse system path
	// messages.
	DescKeyErrFsRefuseSystemPath = "err.fs.refuse-system-path"
	// DescKeyErrFsRefuseSystemPathRoot is the text key for err fs refuse system
	// path root messages.
	DescKeyErrFsRefuseSystemPathRoot = "err.fs.refuse-system-path-root"
	// DescKeyErrFsResolveBase is the text key for err fs resolve base messages.
	DescKeyErrFsResolveBase = "err.fs.resolve-base"
	// DescKeyErrFsResolvePath is the text key for err fs resolve path messages.
	DescKeyErrFsResolvePath = "err.fs.resolve-path"
	// DescKeyErrFsStatPath is the text key for err fs stat path messages.
	DescKeyErrFsStatPath = "err.fs.stat-path"
	// DescKeyErrFsStdinRead is the text key for err fs stdin read messages.
	DescKeyErrFsStdinRead = "err.fs.stdin-read"
	// DescKeyErrFsWriteBuffer is the text key for err fs write buffer messages.
	DescKeyErrFsWriteBuffer = "err.fs.write-buffer"
	// DescKeyErrFsWriteFileFailed is the text key for err fs write file failed
	// messages.
	DescKeyErrFsWriteFileFailed = "err.fs.write-file-failed"
	// DescKeyErrFsWriteMerged is the text key for err fs write merged messages.
	DescKeyErrFsWriteMerged = "err.fs.write-merged"
)

// DescKeys for context directory errors.
const (
	// DescKeyErrContextDirNotFound is the text key for the legacy
	// NotFoundError type's stringified message. Retained for the
	// context/load.NotFoundError shape; under the cwd-anchored
	// model the canonical sentinel is ErrNoCtxHere.
	DescKeyErrContextDirNotFound = "err.context.dir-not-found"
	// DescKeyErrContextNoContextHere is the text key for the
	// user-facing message wrapped by NoCtxHere(cwd) when
	// $PWD/.context does not exist.
	DescKeyErrContextNoContextHere = "err.context.no-context-here"
	// DescKeyErrContextNoContextHereMsg is the text key for the
	// ErrNoCtxHere sentinel's own `.Error()` string (the
	// prefix interpolated via `%w` by the NoCtxHere wrapper).
	DescKeyErrContextNoContextHereMsg = "err.context.no-context-here-msg"
	// DescKeyErrContextDirNotADirectory is the text key for the
	// "$PWD/.context is a regular file, not a directory" rejection.
	DescKeyErrContextDirNotADirectory = "err.context.dir-not-a-directory"
	// DescKeyErrContextDirStat is the text key for stat failures
	// other than not-exist (permission denied, I/O error).
	DescKeyErrContextDirStat = "err.context.dir-stat"
	// DescKeyErrContextNotInitialized is the text key for the
	// "context directory exists but ctx init has not run" rejection.
	// Used when state.Dir() is invoked in a project that has
	// $PWD/.context present but the required context files absent.
	DescKeyErrContextNotInitialized = "err.context.not-initialized"
	// DescKeyErrContextDirNotADirectoryMsg is the text key for
	// the ErrContextDirNotADirectory sentinel's own `.Error()`
	// string.
	DescKeyErrContextDirNotADirectoryMsg = "err.context.dir-not-a-directory-msg"
	// DescKeyErrContextDirStatMsg is the text key for the
	// ErrContextDirStat sentinel's own `.Error()` string.
	DescKeyErrContextDirStatMsg = "err.context.dir-stat-msg"
	// DescKeyErrContextNotInitializedMsg is the text key for the
	// ErrNotInitialized sentinel's own `.Error()` string.
	DescKeyErrContextNotInitializedMsg = "err.context.not-initialized-msg"
)

// DescKeys for filesystem write output.
const (
	// DescKeyWritePathExists is the text key for write path exists messages.
	DescKeyWritePathExists = "write.path-exists"
)
