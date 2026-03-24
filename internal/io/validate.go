//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"path/filepath"
	"strings"

	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
)

// rejectDangerousPath returns an error if the resolved absolute path
// falls under a system directory that ctx should never touch.
func rejectDangerousPath(absPath string) error {
	if absPath == "/" {
		return errFs.RefuseSystemPathRoot()
	}
	for _, prefix := range dangerousPrefixes {
		if strings.HasPrefix(absPath, prefix) {
			return errFs.RefuseSystemPath(absPath)
		}
	}
	return nil
}

// cleanAndValidate resolves a path and checks it against dangerous
// system prefixes. Returns the cleaned path.
func cleanAndValidate(path string) (string, error) {
	clean := filepath.Clean(path)
	abs, absErr := filepath.Abs(clean)
	if absErr != nil {
		return "", errFs.ResolvePath(absErr)
	}
	if checkErr := rejectDangerousPath(abs); checkErr != nil {
		return "", checkErr
	}
	return clean, nil
}
