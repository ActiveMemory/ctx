//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import "fmt"

// PromptExists returns an error when a prompt template already exists.
//
// Parameters:
//   - name: the prompt name that already exists.
//
// Returns:
//   - error: "prompt <name> already exists"
func PromptExists(name string) error {
	return fmt.Errorf("prompt %q already exists", name)
}

// PromptNotFound returns an error when a prompt template does not exist.
//
// Parameters:
//   - name: the prompt name that was not found.
//
// Returns:
//   - error: "prompt <name> not found"
func PromptNotFound(name string) error {
	return fmt.Errorf("prompt %q not found", name)
}

// RemovePrompt wraps a failure to remove a prompt template.
//
// Parameters:
//   - cause: the underlying OS error.
//
// Returns:
//   - error: "remove prompt: <cause>"
func RemovePrompt(cause error) error {
	return fmt.Errorf("remove prompt: %w", cause)
}

// NoPromptTemplate returns an error when no embedded template exists.
//
// Parameters:
//   - name: the template name that was not found.
//
// Returns:
//   - error: advises the user to use --stdin
func NoPromptTemplate(name string) error {
	return fmt.Errorf(
		"no embedded template %q — use --stdin to provide content", name,
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
	return fmt.Errorf("failed to list prompt templates: %w", cause)
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
	return fmt.Errorf("failed to read prompt template %s: %w", name, cause)
}

// ListEntryTemplates wraps a failure to list entry templates.
//
// Parameters:
//   - cause: the underlying error
//
// Returns:
//   - error: "failed to list entry templates: <cause>"
func ListEntryTemplates(cause error) error {
	return fmt.Errorf("failed to list entry templates: %w", cause)
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
	return fmt.Errorf("failed to read entry template %s: %w", name, cause)
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
	return fmt.Errorf("no template available for %s: %w", filename, cause)
}

// ListTemplates wraps a failure to list embedded templates.
//
// Parameters:
//   - cause: the underlying error from the list operation
//
// Returns:
//   - error: "failed to list templates: <cause>"
func ListTemplates(cause error) error {
	return fmt.Errorf("failed to list templates: %w", cause)
}

// ReadTemplate wraps a failure to read an embedded template.
//
// Parameters:
//   - name: template name that failed to read
//   - cause: the underlying error from the read operation
//
// Returns:
//   - error: "failed to read template <name>: <cause>"
func ReadTemplate(name string, cause error) error {
	return fmt.Errorf("failed to read template %s: %w", name, cause)
}

// TemplateMissingMarkers returns an error when a template lacks markers.
//
// Parameters:
//   - kind: marker kind (e.g. "ctx", "plan", "prompt")
//
// Returns:
//   - error: "template missing <kind> markers"
func TemplateMissingMarkers(kind string) error {
	return fmt.Errorf("template missing %s markers", kind)
}

// MarkerNotFound returns an error when a section marker is missing.
//
// Parameters:
//   - kind: marker kind (e.g. "ctx", "plan", "prompt")
//
// Returns:
//   - error: "<kind> start marker not found"
func MarkerNotFound(kind string) error {
	return fmt.Errorf("%s start marker not found", kind)
}
