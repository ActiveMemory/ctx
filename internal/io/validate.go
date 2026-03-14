//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"fmt"
	"path/filepath"
	"strings"
)

// rejectDangerousPath returns an error if the resolved absolute path
// falls under a system directory that ctx should never touch.
func rejectDangerousPath(absPath string) error {
	if absPath == "/" {
		return fmt.Errorf("refusing to access system path: /")
	}
	for _, prefix := range dangerousPrefixes {
		if strings.HasPrefix(absPath, prefix) {
			return fmt.Errorf("refusing to access system path: %s", absPath)
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
		return "", fmt.Errorf("resolve path: %w", absErr)
	}
	if checkErr := rejectDangerousPath(abs); checkErr != nil {
		return "", checkErr
	}
	return clean, nil
}
