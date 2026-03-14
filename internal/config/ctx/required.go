//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ctx

// FilesRequired lists the essential context files that must be present.
//
// These are the files created with `ctx init --minimal` and checked by
// drift detection for missing files.
var FilesRequired = []string{
	Constitution,
	Task,
	Decision,
}
