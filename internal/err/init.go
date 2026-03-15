//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"errors"
	"fmt"
)

// NotInitialized returns an error indicating ctx has not been initialized.
//
// Returns:
//   - error: "ctx: not initialized — run \"ctx init\" first"
func NotInitialized() error {
	return fmt.Errorf("ctx: not initialized — run \"ctx init\" first")
}

// ContextNotInitialized returns an error when no .context/ directory is found.
//
// Returns:
//   - error: "no .context/ directory found. Run 'ctx init' first"
func ContextNotInitialized() error {
	return errors.New("no .context/ directory found. Run 'ctx init' first")
}

// DetectReferenceTime wraps a failure to detect the reference time for changes.
//
// Parameters:
//   - cause: the underlying detection error
//
// Returns:
//   - error: "detecting reference time: <cause>"
func DetectReferenceTime(cause error) error {
	return fmt.Errorf("detecting reference time: %w", cause)
}

// HomeDir wraps a failure to determine the home directory.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "cannot determine home directory: <cause>"
func HomeDir(cause error) error {
	return fmt.Errorf("cannot determine home directory: %w", cause)
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
	return fmt.Errorf("failed to read %s README template: %w", dir, cause)
}

// ReadInitTemplate wraps a failure to read an init template file.
//
// Parameters:
//   - name: template filename that failed to read
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to read <name> template: <cause>"
func ReadInitTemplate(name string, cause error) error {
	return fmt.Errorf("failed to read %s template: %w", name, cause)
}

// CreateMakefile wraps a failure to create a new Makefile.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create Makefile: <cause>"
func CreateMakefile(cause error) error {
	return fmt.Errorf("failed to create Makefile: %w", cause)
}
