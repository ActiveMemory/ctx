//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// NoTemplate wraps a failure to find an embedded template.
//
// Parameters:
//   - filename: Name of the file without a template
//   - cause: the underlying read error
//
// Returns:
//   - error: "no template available for <filename>: <cause>"
func NoTemplate(filename string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptNoTemplate), filename, cause,
	)
}

// ListTemplates wraps a failure to list embedded templates.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to list templates: <cause>"
func ListTemplates(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptListTemplates), cause,
	)
}

// ReadTemplate wraps a failure to read an embedded template.
//
// Parameters:
//   - name: template name that failed to read
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to read template <name>: <cause>"
func ReadTemplate(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptReadTemplate), name, cause,
	)
}

// TemplateMissingMarkers returns an error when a template lacks markers.
//
// Parameters:
//   - kind: marker kind (e.g. "ctx", "prompt")
//
// Returns:
//   - error: "template missing <kind> markers"
func TemplateMissingMarkers(kind string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptTemplateMissingMarkers), kind,
	)
}

// MarkerNotFound returns an error when a section marker is missing.
//
// Parameters:
//   - kind: marker kind (e.g. "ctx", "prompt")
//
// Returns:
//   - error: "<kind> start marker not found"
func MarkerNotFound(kind string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptMarkerNotFound), kind,
	)
}
