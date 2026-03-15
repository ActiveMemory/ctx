//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	ctxerr "github.com/ActiveMemory/ctx/internal/err"
)

// SafeReadFile resolves filename within baseDir, verifies the result
// stays within the base directory boundary, and reads the file content.
//
// Unlike [SafeReadUserFile], this function enforces containment: the
// resolved path must remain under baseDir. Use it when the path is
// constructed from a trusted base and a filename component.
//
// Parameters:
//   - baseDir: trusted root directory
//   - filename: file name (or relative path) to join and validate
//
// Returns:
//   - []byte: file content
//   - error: non-nil if resolution fails, path escapes baseDir, or read fails
func SafeReadFile(baseDir, filename string) ([]byte, error) {
	absBase, absErr := filepath.Abs(baseDir)
	if absErr != nil {
		return nil, fmt.Errorf("resolve base: %w", absErr)
	}

	safe := filepath.Join(absBase, filepath.Base(filename))

	if !strings.HasPrefix(safe, absBase+string(os.PathSeparator)) {
		return nil, fmt.Errorf("path escapes base directory: %s", filename)
	}

	data, readErr := os.ReadFile(safe) //nolint:gosec // validated by boundary check above
	if readErr != nil {
		return nil, readErr
	}

	return data, nil
}

// SafeOpenUserFile opens a file for reading after cleaning the path
// and rejecting system directory prefixes.
//
// Parameters:
//   - path: file path to open
//
// Returns:
//   - *os.File: open file handle (caller must close)
//   - error: non-nil on validation or open failure
func SafeOpenUserFile(path string) (*os.File, error) {
	clean, validateErr := cleanAndValidate(path)
	if validateErr != nil {
		return nil, validateErr
	}
	return os.Open(clean) //nolint:gosec // validated by cleanAndValidate
}

// SafeReadUserFile reads a file after cleaning the path and rejecting
// system directory prefixes.
//
// Parameters:
//   - path: file path to read
//
// Returns:
//   - []byte: file content
//   - error: non-nil on validation or read failure
func SafeReadUserFile(path string) ([]byte, error) {
	clean, validateErr := cleanAndValidate(path)
	if validateErr != nil {
		return nil, validateErr
	}
	return os.ReadFile(clean) //nolint:gosec // validated by cleanAndValidate
}

// SafeAppendFile opens a file for appending after cleaning the path
// and rejecting system directory prefixes. Creates the file if it does
// not exist.
//
// Parameters:
//   - path: file path to open
//   - perm: file permission bits used when creating the file
//
// Returns:
//   - *os.File: open file handle in append mode (caller must close)
//   - error: non-nil on validation or open failure
func SafeAppendFile(path string, perm os.FileMode) (*os.File, error) {
	clean, validateErr := cleanAndValidate(path)
	if validateErr != nil {
		return nil, validateErr
	}
	//nolint:gosec // validated by cleanAndValidate
	return os.OpenFile(clean, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
}

// SafeWriteFile writes data to a file after cleaning the path and
// rejecting system directory prefixes.
//
// Parameters:
//   - path: file path to write
//   - data: content to write
//   - perm: file permission bits
//
// Returns:
//   - error: non-nil on validation or write failure
func SafeWriteFile(path string, data []byte, perm os.FileMode) error {
	clean, validateErr := cleanAndValidate(path)
	if validateErr != nil {
		return validateErr
	}
	return os.WriteFile(clean, data, perm) //nolint:gosec // validated by cleanAndValidate
}

// maxRedirects caps the number of HTTP redirects the client will follow.
const maxRedirects = 3

// SafePost sends an HTTP POST with the given content type and body.
//
// Designed for static endpoint URLs that originate from trusted,
// user-configured sources (e.g. webhook URLs stored in AES-256-GCM
// encrypted storage). Centralizes gosec suppression so callers don't
// each need their own nolint pragma.
//
// Protections applied:
//   - Scheme validation: rejects everything except http and https,
//     preventing file://, gopher://, and other protocol smuggling.
//   - Redirect cap: follows at most 3 redirects (Go default is 10).
//     Limits open-redirect abuse where a trusted URL bounces to an
//     unintended destination.
//   - Caller-specified timeout: bounds total request duration
//     including redirects.
//
// Threats explicitly not mitigated (and why that is acceptable):
//   - SSRF to private IPs: the URL is a static, user-configured
//     endpoint — not attacker-controlled input. Blocking RFC 1918
//     ranges would break legitimate local webhook receivers.
//   - Response body size: callers are fire-and-forget (close body
//     immediately), so unbounded reads are not a concern.
//   - TLS certificate pinning: the endpoint is user-chosen; standard
//     system CA validation is appropriate.
//
// Parameters:
//   - rawURL: destination endpoint (trusted, user-configured origin)
//   - contentType: MIME type for the Content-Type header
//   - body: request payload
//   - timeout: per-request timeout (includes redirect hops)
//
// Returns:
//   - *http.Response: the HTTP response (caller must close Body)
//   - error: on scheme validation failure, redirect cap, or HTTP error
func SafePost(rawURL, contentType string, body []byte, timeout time.Duration) (*http.Response, error) {
	if schemeErr := validateHTTPScheme(rawURL); schemeErr != nil {
		return nil, schemeErr
	}

	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return ctxerr.TooManyRedirects()
			}
			return nil
		},
	}

	//nolint:gosec // URL originates from trusted, encrypted storage; scheme validated above
	return client.Post(rawURL, contentType, bytes.NewReader(body))
}

// validateHTTPScheme parses the URL and rejects any scheme other than
// http or https.
func validateHTTPScheme(rawURL string) error {
	parsed, parseErr := url.Parse(rawURL)
	if parseErr != nil {
		return fmt.Errorf("parse URL: %w", parseErr)
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return ctxerr.UnsafeURLScheme(parsed.Scheme)
	}
	return nil
}
