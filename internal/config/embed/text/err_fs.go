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
	// DescKeyErrContextDirNotFound is the text key for err context dir not found
	// messages.
	DescKeyErrContextDirNotFound = "err.context.dir-not-found"
	// DescKeyErrContextNotDeclaredZero is the text key used when CTX_DIR
	// is not set and no .context/ candidate is visible from CWD.
	DescKeyErrContextNotDeclaredZero = "err.context.not-declared-zero"
	// DescKeyErrContextNotDeclaredOne is the text key used when CTX_DIR
	// is not set and exactly one .context/ candidate is visible from CWD.
	DescKeyErrContextNotDeclaredOne = "err.context.not-declared-one"
	// DescKeyErrContextNotDeclaredMany is the text key used when CTX_DIR
	// is not set and two or more .context/ candidates are visible from CWD.
	DescKeyErrContextNotDeclaredMany = "err.context.not-declared-many"
	// DescKeyErrContextRelativeNotAllowed is the text key for the
	// "CTX_DIR must be absolute" rejection.
	DescKeyErrContextRelativeNotAllowed = "err.context.relative-not-allowed"
	// DescKeyErrContextNonCanonicalBasename is the text key for the
	// "CTX_DIR basename must be .context" rejection.
	DescKeyErrContextNonCanonicalBasename = "err.context.non-canonical-basename"
	// DescKeyErrContextDirNotADirectory is the text key for the
	// "CTX_DIR points at a file, not a directory" rejection.
	DescKeyErrContextDirNotADirectory = "err.context.dir-not-a-directory"
	// DescKeyErrContextDirStat is the text key for stat failures
	// other than not-exist (permission denied, I/O error).
	DescKeyErrContextDirStat = "err.context.dir-stat"
)

// DescKeys for filesystem write output.
const (
	// DescKeyWritePathExists is the text key for write path exists messages.
	DescKeyWritePathExists = "write.path-exists"
)
