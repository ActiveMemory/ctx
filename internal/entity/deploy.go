//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// DeployParams holds configuration for deploying embedded templates to a
// subdirectory.
//
// Fields:
//   - SubDir: Target subdirectory within the deploy root
//   - ListErrKey: Error key for directory listing failures
//   - ReadErrKey: Error key for file read failures
type DeployParams struct {
	SubDir     string
	ListErrKey string
	ReadErrKey string
}
