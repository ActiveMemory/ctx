//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package path resolves file system paths for the task
// subsystem. It centralizes path construction so that
// sibling packages (archive, complete, count) do not
// need to import config constants directly.
//
// # Path Resolution
//
// [File] returns the absolute path to TASKS.md by
// joining rc.ContextDir with the task filename constant.
// Every task operation that reads or writes TASKS.md
// calls this function for the canonical path.
//
// [ArchiveDir] returns the absolute path to the archive
// directory (.context/archive/) where completed task
// blocks are moved during archival. The archive package
// uses tidy.WriteArchive which handles directory creation,
// so callers do not need to ensure the directory exists.
//
// # Design
//
// Both functions are thin wrappers around filepath.Join
// with config constants. They exist to avoid scattering
// path construction across multiple packages and to
// provide a single place to change if the directory
// layout evolves.
package path
