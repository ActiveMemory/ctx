//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

// Tool-native directory and extension constants used by
// steering sync to write files in each tool's format.
const (
	// DirCursorDot is the Cursor configuration directory.
	DirCursorDot = ".cursor"
	// DirRules is the Cursor rules subdirectory.
	DirRules = "rules"
	// ExtMDC is the Cursor MDC rule file extension.
	ExtMDC = ".mdc"
	// DirClinerules is the Cline rules directory.
	DirClinerules = ".clinerules"
	// DirKiroDot is the Kiro configuration directory.
	DirKiroDot = ".kiro"
	// DirSteering is the Kiro steering subdirectory.
	DirSteering = "steering"
)

// LabelAllTools is the display label when a steering
// or trigger item applies to all tools.
const LabelAllTools = "all"

// DefaultPriority is the default injection priority for
// steering files when omitted from frontmatter.
const DefaultPriority = 50

// Foundation steering file names used by ctx steering init
// and ctx init to scaffold the starter set.
const (
	// NameProduct is the file name for the product context file.
	NameProduct = "product"
	// NameTech is the file name for the technology stack file.
	NameTech = "tech"
	// NameStructure is the file name for the project structure file.
	NameStructure = "structure"
	// NameWorkflow is the file name for the development workflow file.
	NameWorkflow = "workflow"
)

// Tombstone is the literal marker line embedded in scaffolded
// foundation steering file bodies by ctx steering init. Its
// presence in a steering file's body signals that the file has
// not been customized yet and that the body content is
// unmodified placeholder text. Files containing this marker
// are excluded from the agent context packet, MCP
// ctx_steering_get results, and native-tool sync (Cursor /
// Cline / Kiro). Removing the line activates the file. The
// marker is HTML-comment-shaped so it survives Markdown
// rendering invisibly, and the match is a literal string
// comparison so non-English body customizations still remove
// it reliably.
const Tombstone = "<!-- remove this after you edit the steering file !-->"
