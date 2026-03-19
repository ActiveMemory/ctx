//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// NotInitialized returns an error indicating ctx has not been initialized.
//
// Returns:
//   - error: "ctx: not initialized: run \"ctx init\" first"
func NotInitialized() error {
	return errors.New(desc.TextDesc(text.TextDescKeyErrInitNotInitialized))
}

// ContextNotInitialized returns an error when no .context/ directory is found.
//
// Returns:
//   - error: "no .context/ directory found. Run 'ctx init' first"
func ContextNotInitialized() error {
	return errors.New(
		desc.TextDesc(text.TextDescKeyErrInitContextNotInitialized),
	)
}

// DetectReferenceTime wraps a failure to detect the reference time for changes.
//
// Parameters:
//   - cause: the underlying detection error
//
// Returns:
//   - error: "detecting reference time: <cause>"
func DetectReferenceTime(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.TextDescKeyErrInitDetectReferenceTime), cause,
	)
}

// HomeDir wraps a failure to determine the home directory.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "cannot determine home directory: <cause>"
func HomeDir(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.TextDescKeyErrInitHomeDir), cause,
	)
}

// ReadProjectReadme wraps a failure to read a project README template.
//
// Parameters:
//   - dir: directory name whose README failed to read
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to read <dir> README template: <cause>"
func ReadProjectReadme(dir string, cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.TextDescKeyErrInitReadProjectReadme), dir, cause,
	)
}

// ReadTemplate wraps a failure to read an init template file.
//
// Parameters:
//   - name: template filename that failed to read
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to read <name> template: <cause>"
func ReadTemplate(name string, cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.TextDescKeyErrInitReadInitTemplate), name, cause,
	)
}

// CreateMakefile wraps a failure to create a new Makefile.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create Makefile: <cause>"
func CreateMakefile(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.TextDescKeyErrInitCreateMakefile), cause,
	)
}

// CtxNotInPath returns an error indicating that ctx was not found in
// PATH.
//
// Returns:
//   - error: "ctx not found in PATH"
func CtxNotInPath() error {
	return errors.New(
		desc.TextDesc(text.TextDescKeyErrValidateCtxNotInPath),
	)
}
