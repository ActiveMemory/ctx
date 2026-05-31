//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

import (
	"bytes"
	"text/template"
)

// ObsidianReadme renders the README for a generated Obsidian vault.
// Data: [ObsidianData]. Call sites render via [Render], passing this
// handle — never a name literal — so non-exempt caller packages stay
// clean under audit/magic_strings.
var ObsidianReadme *template.Template

// JournalSiteReadme renders the README for the journal-site directory.
// Data: [JournalSiteData].
var JournalSiteReadme *template.Template

// TriggerScript renders the scaffold bash script for `ctx trigger add`.
// Data: [TriggerData].
var TriggerScript *template.Template

// Learning renders a learning entry section. Data: [LearningData].
var Learning *template.Template

// Decision renders a decision (ADR) entry section. Data: [DecisionData].
var Decision *template.Template

// Render executes a parsed template handle against data.
//
// The handle is always non-nil for a registered template (a parse
// failure still yields a usable empty template, recorded for
// TestTemplatesParse), so this never panics on a nil handle. An
// execution error (e.g. a renamed data field) is returned, not
// panicked; golden tests gate template correctness.
//
// Parameters:
//   - t: a parsed template handle (e.g. [ObsidianReadme])
//   - data: the template's typed data struct
//
// Returns:
//   - string: the rendered output
//   - error: non-nil on an execution failure
func Render(t *template.Template, data any) (string, error) {
	var buf bytes.Buffer
	if execErr := t.Execute(&buf, data); execErr != nil {
		return "", execErr
	}
	return buf.String(), nil
}
