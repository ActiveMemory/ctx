//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package kb

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Mkdir wraps `os.MkdirAll` for a scaffolded directory.
//
// Parameters:
//   - path: the directory path that failed.
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func Mkdir(path string, cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrInitKbMkdir), path, cause)
}

// CopyIngestTemplates wraps a failure to copy the embedded
// ingest-side templates.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func CopyIngestTemplates(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrInitKbCopyIngestTemplates), cause,
	)
}

// CopySchemas wraps a failure to copy the embedded schemas.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func CopySchemas(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrInitKbCopySchemas), cause)
}

// CopyLanding wraps a failure to copy the embedded kb
// landing page.
//
// Parameters:
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func CopyLanding(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrInitKbCopyLanding), cause)
}

// ReadEmbed wraps a failure to read an embedded asset path.
//
// Parameters:
//   - path: embedded asset path that failed.
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func ReadEmbed(path string, cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrInitKbReadEmbed), path, cause)
}

// MkdirFor wraps a parent-dir create failure for a specific
// destination file.
//
// Parameters:
//   - dst: destination path whose parent could not be created.
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func MkdirFor(dst string, cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrInitKbMkdirFor), dst, cause)
}

// WriteFile wraps a write failure for a specific destination
// path.
//
// Parameters:
//   - dst: destination path that failed.
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func WriteFile(dst string, cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrInitKbWriteFile), dst, cause)
}

// ReadDir wraps a `os.ReadDir` failure for a scaffolded
// directory.
//
// Parameters:
//   - dir: directory whose listing failed.
//   - cause: underlying error.
//
// Returns:
//   - error: wrapped failure.
func ReadDir(dir string, cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrInitKbReadDir), dir, cause)
}
