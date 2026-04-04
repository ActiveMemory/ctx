//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// ComputeRelPath wraps a failure to compute a relative path.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "compute relative path: <cause>"
func ComputeRelPath(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringComputeRelPath), cause,
	)
}

// ContextDirMissing returns an error when the .context/ directory
// does not exist.
//
// Returns:
//   - error: ".context/ directory does not exist; run ctx init first"
func ContextDirMissing() error {
	return errors.New(
		desc.Text(text.DescKeyErrSteeringContextDirMissing),
	)
}

// CreateDir wraps a steering directory creation failure.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "create steering directory: <cause>"
func CreateDir(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringCreateDir), cause,
	)
}

// FileExists returns an error when a steering file already exists.
//
// Parameters:
//   - path: the existing file path
//
// Returns:
//   - error: "steering file already exists: <path>"
func FileExists(path string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringFileExists), path,
	)
}

// InvalidYAML wraps an invalid YAML frontmatter parse failure.
//
// Parameters:
//   - filePath: path to the steering file
//   - cause: the underlying YAML error
//
// Returns:
//   - error: "steering: <filePath>: invalid YAML frontmatter: <cause>"
func InvalidYAML(filePath string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringInvalidYAML), filePath, cause,
	)
}

// MissingClosingDelimiter returns an error for missing closing
// frontmatter delimiter.
//
// Returns:
//   - error: "missing closing frontmatter delimiter (---)"
func MissingClosingDelimiter() error {
	return errors.New(
		desc.Text(text.DescKeyErrSteeringMissingClosingDelim),
	)
}

// MissingOpeningDelimiter returns an error for missing opening
// frontmatter delimiter.
//
// Returns:
//   - error: "missing opening frontmatter delimiter (---)"
func MissingOpeningDelimiter() error {
	return errors.New(
		desc.Text(text.DescKeyErrSteeringMissingOpeningDelim),
	)
}

// NoTool returns an error when no tool is specified for sync.
//
// Returns:
//   - error: "no tool specified: use --tool <tool>, --all, or set
//     the tool field in .ctxrc"
func NoTool() error {
	return errors.New(
		desc.Text(text.DescKeyErrSteeringNoTool),
	)
}

// OutputEscapesRoot returns an error when an output path escapes
// the project root.
//
// Parameters:
//   - outPath: the output path
//   - projectRoot: the project root path
//
// Returns:
//   - error: "output path <outPath> escapes project root <projectRoot>"
func OutputEscapesRoot(outPath, projectRoot string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringOutputEscapesRoot),
		outPath, projectRoot,
	)
}

// Parse wraps a steering file parse failure.
//
// Parameters:
//   - filePath: path to the steering file
//   - cause: the underlying error
//
// Returns:
//   - error: "steering: <filePath>: <cause>"
func Parse(filePath string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringParse), filePath, cause,
	)
}

// ReadDir wraps a steering directory read failure.
//
// Parameters:
//   - dir: the directory path
//   - cause: the underlying error
//
// Returns:
//   - error: "steering: read directory <dir>: <cause>"
func ReadDir(dir string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringReadDir), dir, cause,
	)
}

// ReadFile wraps a steering file read failure.
//
// Parameters:
//   - path: the file path
//   - cause: the underlying error
//
// Returns:
//   - error: "steering: read file <path>: <cause>"
func ReadFile(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringReadFile), path, cause,
	)
}

// ResolveOutput wraps a failure to resolve an output path.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "resolve output path: <cause>"
func ResolveOutput(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringResolveOutput), cause,
	)
}

// ResolveRoot wraps a failure to resolve the project root path.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "resolve project root: <cause>"
func ResolveRoot(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringResolveRoot), cause,
	)
}

// SyncAll wraps a failure during sync-all for a specific tool.
//
// Parameters:
//   - tool: the tool that failed
//   - cause: the underlying error
//
// Returns:
//   - error: "steering: sync <tool>: <cause>"
func SyncAll(tool string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringSyncAll), tool, cause,
	)
}

// SyncName wraps a steering sync error for a named file.
//
// Parameters:
//   - name: steering file name
//   - cause: the underlying error
//
// Returns:
//   - error: "steering: <name>: <cause>"
func SyncName(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringSyncName), name, cause,
	)
}

// UnsupportedTool returns an error for an unsupported sync tool.
//
// Parameters:
//   - tool: the unsupported tool name
//   - supported: comma-separated list of supported tools
//
// Returns:
//   - error: "steering: unsupported sync tool <tool>; supported: <supported>"
func UnsupportedTool(tool, supported string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringUnsupportedTool),
		tool, supported,
	)
}

// WriteFile wraps a steering file write failure during sync.
//
// Parameters:
//   - path: the output path
//   - cause: the underlying error
//
// Returns:
//   - error: "steering: write <path>: <cause>"
func WriteFile(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringWriteFile), path, cause,
	)
}

// Write wraps a steering file write failure.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "write steering file: <cause>"
func Write(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringWriteSteeringFile), cause,
	)
}

// WriteInitFile wraps a steering init file write failure.
//
// Parameters:
//   - path: the file path
//   - cause: the underlying OS error
//
// Returns:
//   - error: "write <path>: <cause>"
func WriteInitFile(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSteeringWriteInitFile), path, cause,
	)
}
