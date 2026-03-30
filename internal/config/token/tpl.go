//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package token

// Go template syntax constants.
const (
	// GoTplFieldOpen is the Go template field access opening delimiter.
	GoTplFieldOpen = "{{."
	// GoTplClose is the Go template closing delimiter.
	GoTplClose = "}}"
)

// TemplateMarkers are content substrings that indicate a file is a template.
var TemplateMarkers = []string{
	"YOUR_",
	"<your",
	"{{",
	"REPLACE_",
	"TODO",
	"CHANGEME",
	"FIXME",
}
