//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"net/url"
	"path/filepath"
	"strings"

	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errHTTP "github.com/ActiveMemory/ctx/internal/err/http"
)

// rejectDangerousPath returns an error if the resolved absolute path
// falls under a system directory that ctx should never touch.
//
// Parameters:
//   - absPath: Resolved absolute path to check
//
// Returns:
//   - error: Non-nil if the path is root or under a dangerous prefix
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
//
// Parameters:
//   - path: Raw path to clean and validate
//
// Returns:
//   - string: Cleaned path on success
//   - error: Non-nil if resolution fails or the path is dangerous
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

// validateHTTPScheme parses the URL and rejects any scheme other than
// http or https.
//
// Parameters:
//   - rawURL: URL string to validate
//
// Returns:
//   - error: Non-nil if the URL is unparseable or uses a non-HTTP scheme
func validateHTTPScheme(rawURL string) error {
	parsed, parseErr := url.Parse(rawURL)
	if parseErr != nil {
		return errHTTP.ParseURL(parseErr)
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return errHTTP.UnsafeURLScheme(parsed.Scheme)
	}
	return nil
}
