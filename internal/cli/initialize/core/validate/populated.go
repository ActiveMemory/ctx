//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
)

// PopulatedFiles returns the basenames of the essential context
// files that already exist in contextDir.
//
// Used by ctx init to enumerate the destructive blast radius
// before refusing or, with --reset, before backing up and
// overwriting. Returns nil when none of the essential files
// exist (the directory is effectively empty / never
// initialized).
//
// Parameters:
//   - contextDir: Absolute path to the context directory
//
// Returns:
//   - []string: Basenames of populated essential files
//     (subset of ctx.FilesRequired), in the canonical order
//     declared by FilesRequired.
func PopulatedFiles(contextDir string) []string {
	var present []string
	for _, f := range ctx.FilesRequired {
		if _, statErr := os.Stat(filepath.Join(contextDir, f)); statErr == nil {
			present = append(present, f)
		}
	}
	return present
}
