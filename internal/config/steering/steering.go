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
