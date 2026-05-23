//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package proposal

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// MkdirProposals wraps an os.MkdirAll failure on the
// `.context/proposals/` directory.
//
// Parameters:
//   - path: the proposals-dir path.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func MkdirProposals(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrProposalMkdir), path, cause,
	)
}

// Write wraps an os.WriteFile failure on a single
// proposal markdown file.
//
// Parameters:
//   - path: the proposal file path.
//   - cause: the underlying OS error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func Write(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrProposalWrite), path, cause,
	)
}

// ReadInput wraps an io.ReadAll failure on the input
// stream passed to the extract command.
//
// Parameters:
//   - cause: the underlying I/O error.
//
// Returns:
//   - error: wrapped for operator-friendly output.
func ReadInput(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrProposalReadInput), cause,
	)
}

// EmptyInput reports that the input passed to extract
// was empty after trimming whitespace; the caller has
// nothing to extract from.
//
// Returns:
//   - error: typed error string from text/errors.yaml.
func EmptyInput() error {
	return errors.New(desc.Text(text.DescKeyErrProposalEmptyInput))
}
