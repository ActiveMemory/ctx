//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

// SMBConfig holds parsed SMB share connection details.
type SMBConfig struct {
	Host      string
	Share     string
	Subdir    string
	GVFSPath  string
	SourceURL string
}
