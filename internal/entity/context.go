//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// Context represents the loaded context from a .context/ directory.
//
// Fields:
//   - Dir: Path to the context directory
//   - Files: All loaded context files with their metadata
//   - TotalTokens: Sum of estimated tokens across all files
//   - TotalSize: Sum of file sizes in bytes
type Context struct {
	Dir         string
	Files       []FileInfo
	TotalTokens int
	TotalSize   int64
}

// File returns the FileInfo with the given name, or nil if not found.
//
// Fields:
//   - name: Name of the file to search for.
//
// Returns:
//   - *FileInfo: Pointer to the found FileInfo, or nil if not found.
func (c *Context) File(name string) *FileInfo {
	for i := range c.Files {
		if c.Files[i].Name == name {
			return &c.Files[i]
		}
	}
	return nil
}
