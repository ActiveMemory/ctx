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

// Exists returns an error when a prompt template already exists.
//
// Parameters:
//   - name: the prompt name that already exists.
//
// Returns:
//   - error: "prompt <name> already exists"
func Exists(name string) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrPromptExists), name)
}

// NotFound returns an error when a prompt template does not exist.
//
// Parameters:
//   - name: the prompt name that was not found.
//
// Returns:
//   - error: "prompt <name> not found"
func NotFound(name string) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrPromptNotFound), name)
}

// Remove wraps a failure to remove a prompt template.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "remove prompt: <cause>"
func Remove(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptRemovePrompt), cause,
	)
}

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

// NoPromptTemplate returns an error when no embedded template exists.
//
// Parameters:
//   - name: the template name that was not found
//
// Returns:
//   - error: advises the user to use --stdin
func NoPromptTemplate(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptNoPromptTemplate), name,
	)
}

// ListPromptTemplates wraps a failure to list prompt templates.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to list prompt templates: <cause>"
func ListPromptTemplates(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptListPromptTemplates), cause,
	)
}

// ReadPromptTemplate wraps a failure to read a prompt template.
//
// Parameters:
//   - name: template name that failed to read
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to read prompt template <name>: <cause>"
func ReadPromptTemplate(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptReadPromptTemplate), name, cause,
	)
}

// ListEntryTemplates wraps a failure to list entry templates.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to list entry templates: <cause>"
func ListEntryTemplates(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptListEntryTemplates), cause,
	)
}

// ReadEntryTemplate wraps a failure to read an entry template.
//
// Parameters:
//   - name: template name that failed to read
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to read entry template <name>: <cause>"
func ReadEntryTemplate(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptReadEntryTemplate), name, cause,
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
//   - kind: marker kind (e.g. "ctx", "plan", "prompt")
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
//   - kind: marker kind (e.g. "ctx", "plan", "prompt")
//
// Returns:
//   - error: "<kind> start marker not found"
func MarkerNotFound(kind string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPromptMarkerNotFound), kind,
	)
}
