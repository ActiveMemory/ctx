//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

import (
	"embed"
	"text/template"
)

// templatesFS holds the multi-line template bodies migrated out of the
// fmt.Sprintf format-string constants. The embed is local to tpl: tpl
// is a leaf package, and reaching into the parent assets.FS would
// couple it there and invite the import cycle the embed_test split
// fought.
//
//go:embed templates/*.tmpl
var templatesFS embed.FS

// parseErrs accumulates init-time template parse failures. It is empty
// in any correct build; TestTemplatesParse asserts so, turning a
// malformed embedded template into a CI failure rather than a runtime
// panic (the project forbids panic, and there is no template.Must
// precedent here).
var parseErrs []error

// init parses every embedded template into its exported handle.
func init() {
	ObsidianReadme = parseTemplate("templates/obsidian-readme.md.tmpl")
	JournalSiteReadme = parseTemplate("templates/journal-site-readme.md.tmpl")
	TriggerScript = parseTemplate("templates/trigger-script.sh.tmpl")
	Learning = parseTemplate("templates/learning.md.tmpl")
	Decision = parseTemplate("templates/decision.md.tmpl")
}

// parseTemplate reads and parses one embedded template. On failure it
// records the cause in parseErrs and returns the non-nil (empty)
// template, so Render never receives a nil handle: the failure path
// stays panic-free while TestTemplatesParse flags it.
//
// Parameters:
//   - path: embedded template path under templatesFS
//
// Returns:
//   - *template.Template: the parsed template (never nil)
func parseTemplate(path string) *template.Template {
	t := template.New(path)
	body, readErr := templatesFS.ReadFile(path)
	if readErr != nil {
		parseErrs = append(parseErrs, readErr)
		return t
	}
	if _, parseErr := t.Parse(string(body)); parseErr != nil {
		parseErrs = append(parseErrs, parseErr)
	}
	return t
}
