//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package token

// Constants for prefix tokens.
const (
	// PrefixHeading is the Markdown heading character used for prefix checks.
	PrefixHeading = "#"
	// PrefixDot is the dot prefix for hidden files and directories.
	PrefixDot = "."
	// PrefixListDash is the prefix for a dash list item.
	PrefixListDash = "- "
	// PrefixListStar is the prefix for a star list item.
	PrefixListStar = "* "
	// PrefixComment is the inline comment marker in YAML.
	PrefixComment = "#"
	// PrefixBang is the Markdown image/admonition prefix.
	PrefixBang = "!"
	// PrefixStar is the Markdown emphasis/bold prefix.
	PrefixStar = "*"
)

// Home directory prefix constants.
const (
	// PrefixHomeDir is the Unix home directory shorthand.
	PrefixHomeDir = "~/"
)

// URL and link prefix constants.
const (
	// LinkPrefixParent is the relative link prefix to the parent directory.
	LinkPrefixParent = "../"
	// PrefixHTTP is the scheme prefix for HTTP/HTTPS URLs.
	PrefixHTTP = "http"
	// PrefixProtocolRelative is the protocol-relative URL prefix.
	PrefixProtocolRelative = "//"
)

// Template and glob indicator characters.
const (
	// GlobStar is the wildcard character in glob patterns.
	GlobStar = "*"
	// TemplateBrace is the opening brace in template/placeholder patterns.
	TemplateBrace = "{"
)
