//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

const (
	// NewlineCRLF is the Windows new line.
	//
	// We check NewlineCRLF first, then NewlineLF to handle both formats.
	NewlineCRLF = "\r\n"
	// NewlineLF is Unix new line.
	NewlineLF = "\n"
	// Whitespace is the set of inline whitespace characters (space and tab).
	Whitespace = " \t"
	// Space is a single space character.
	Space = " "
	// Tab is a horizontal tab character.
	Tab = "\t"
	// Colon is the colon character used as a key-value separator.
	Colon = ":"
	// Dash is a hyphen used as a timestamp segment separator.
	Dash = "-"
	// KeyValueSep is the equals sign used as a key-value separator in state files.
	KeyValueSep = "="
	// Separator is a Markdown horizontal rule used between sections.
	Separator = "---"
	// Ellipsis is a Markdown ellipsis.
	Ellipsis = "..."
	// HeadingLevelOneStart is the Markdown heading for the first section.
	HeadingLevelOneStart = "# "
	// HeadingLevelTwoStart is the Markdown heading for subsequent sections.
	HeadingLevelTwoStart = "## "
	// HeadingLevelThreeStart is the Markdown heading level three.
	HeadingLevelThreeStart = "### "
	// PrefixListDash is the prefix for a dash list item.
	PrefixListDash = "- "
	// PrefixListStar is the prefix for a star list item.
	PrefixListStar = "* "
	// MemoryMirrorPrefix is the filename prefix for archived mirror files.
	MemoryMirrorPrefix = "mirror-"
	// CodeFence is the standard Markdown code fence delimiter.
	CodeFence = "```"
	// Backtick is a single backtick character.
	Backtick = "`"
	// PipeSeparator is the inline separator used between navigation links.
	PipeSeparator = " | "
	// LinkPrefixParent is the relative link prefix to the parent directory.
	LinkPrefixParent = "../"
	// PrefixHeading is the Markdown heading character used for prefix checks.
	PrefixHeading = "#"
	// PrefixBracket is the opening bracket used for placeholder checks.
	PrefixBracket = "["
	// LoopComplete is the banner printed when the loop finishes.
	LoopComplete = "=== Loop Complete ==="
	// TomlNavOpen is the opening bracket for the TOML nav array.
	TomlNavOpen = "nav = ["
	// TomlNavSectionClose closes a nav section group.
	TomlNavSectionClose = "  ]}"
	// TomlNavClose closes the top-level nav array.
	TomlNavClose = "]"
	// NudgeBoxBottom is the bottom border of a nudge/notification box.
	NudgeBoxBottom = "└──────────────────────────────────────────────────"
)

// MCP constants.
const (
	// MCPResourceURIPrefix is the URI scheme prefix for MCP context resources.
	MCPResourceURIPrefix = "ctx://context/"
	// MimeMarkdown is the MIME type for Markdown content.
	MimeMarkdown = "text/markdown"
	// MCPScanMaxSize is the maximum scanner buffer size for MCP messages (1 MB).
	MCPScanMaxSize = 1 << 20
	// MCPMethodInitialize is the MCP initialize handshake method.
	MCPMethodInitialize = "initialize"
	// MCPMethodPing is the MCP ping method.
	MCPMethodPing = "ping"
	// MCPMethodResourcesList is the MCP method for listing resources.
	MCPMethodResourcesList = "resources/list"
	// MCPMethodResourcesRead is the MCP method for reading a resource.
	MCPMethodResourcesRead = "resources/read"
	// MCPMethodToolsList is the MCP method for listing tools.
	MCPMethodToolsList = "tools/list"
	// MCPMethodToolsCall is the MCP method for calling a tool.
	MCPMethodToolsCall = "tools/call"
	// JSONRPCVersion is the JSON-RPC protocol version string.
	JSONRPCVersion = "2.0"
	// MCPServerName is the server name reported during initialization.
	MCPServerName = "ctx"
	// MCPContentTypeText is the content type for text tool output.
	MCPContentTypeText = "text"
	// SchemaObject is the JSON Schema type for objects.
	SchemaObject = "object"
	// SchemaString is the JSON Schema type for strings.
	SchemaString = "string"
	// MCPToolStatus is the MCP tool name for context status.
	MCPToolStatus = "ctx_status"
	// MCPToolAdd is the MCP tool name for adding entries.
	MCPToolAdd = "ctx_add"
	// MCPToolComplete is the MCP tool name for completing tasks.
	MCPToolComplete = "ctx_complete"
	// MCPToolDrift is the MCP tool name for drift detection.
	MCPToolDrift = "ctx_drift"
)

// Content detection constants.
const (
	// ByteNewline is the newline character as a byte.
	ByteNewline = '\n'
	// ByteHeading is the heading character as a byte for content scanning.
	ByteHeading = '#'
	// ByteDash is the dash character as a byte for separator detection.
	ByteDash = '-'
	// MaxSeparatorLen is the maximum length of a line to be considered a
	// Markdown separator (e.g. "---" or "----").
	MaxSeparatorLen = 5
	// ParserBufInitSize is the initial scanner buffer size for session parsing (64 KB).
	ParserBufInitSize = 64 * 1024
	// ParserBufMaxSize is the maximum scanner buffer size for session parsing (1 MB).
	ParserBufMaxSize = 1024 * 1024
)

// MkDocs stripping constants (used by "ctx why" to clean embedded docs).
const (
	// MkDocsAdmonitionPrefix is the prefix for admonition lines in MkDocs.
	MkDocsAdmonitionPrefix = "!!!"
	// MkDocsTabPrefix is the prefix for tab marker lines in MkDocs.
	MkDocsTabPrefix = "=== "
	// MkDocsIndent is the 4-space indentation used in admonition/tab bodies.
	MkDocsIndent = "    "
	// MkDocsIndentWidth is the number of characters to dedent from body lines.
	MkDocsIndentWidth = 4
	// MkDocsFrontmatterDelim is the YAML frontmatter delimiter.
	MkDocsFrontmatterDelim = "---"
)

// SecretPatterns are filename substrings that indicate potential secret files.
var SecretPatterns = []string{
	".env",
	"credentials",
	"secret",
	"api_key",
	"apikey",
	"password",
}

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

// XML attribute name constants for context-update tag parsing.
const (
	// AttrType is the "type" attribute on a context-update tag.
	AttrType = "type"
	// AttrContext is the "context" attribute on a context-update tag.
	AttrContext = "context"
	// AttrLesson is the "lesson" attribute on a context-update tag.
	AttrLesson = "lesson"
	// AttrApplication is the "application" attribute on a context-update tag.
	AttrApplication = "application"
	// AttrRationale is the "rationale" attribute on a context-update tag.
	AttrRationale = "rationale"
	// AttrConsequences is the "consequences" attribute on a context-update tag.
	AttrConsequences = "consequences"
)

// Runtime configuration defaults (overridable via .ctxrc).
const (
	// DefaultRcTokenBudget is the default token budget for context assembly.
	DefaultRcTokenBudget = 8000
	// DefaultRcArchiveAfterDays is the default days before archiving completed tasks.
	DefaultRcArchiveAfterDays = 7
	// DefaultRcEntryCountLearnings is the entry count threshold for LEARNINGS.md.
	DefaultRcEntryCountLearnings = 30
	// DefaultRcEntryCountDecisions is the entry count threshold for DECISIONS.md.
	DefaultRcEntryCountDecisions = 20
	// DefaultRcConventionLineCount is the line count threshold for CONVENTIONS.md.
	DefaultRcConventionLineCount = 200
	// DefaultRcInjectionTokenWarn is the token threshold for oversize injection warning.
	DefaultRcInjectionTokenWarn = 15000
	// DefaultRcContextWindow is the default context window size in tokens.
	DefaultRcContextWindow = 200000
	// DefaultRcTaskNudgeInterval is the Edit/Write calls between task completion nudges.
	DefaultRcTaskNudgeInterval = 5
	// DefaultRcKeyRotationDays is the days before encryption key rotation nudge.
	DefaultRcKeyRotationDays = 90
)
