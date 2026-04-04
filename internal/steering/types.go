//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

// InclusionMode determines when a steering file is injected into
// an AI prompt.
type InclusionMode string

const (
	// InclusionAlways includes the steering file in every context packet.
	InclusionAlways InclusionMode = "always"
	// InclusionAuto includes the steering file when the prompt matches
	// the file's description.
	InclusionAuto InclusionMode = "auto"
	// InclusionManual includes the steering file only when explicitly
	// referenced by name.
	InclusionManual InclusionMode = "manual"
)

// SteeringFile represents a parsed steering file with YAML frontmatter
// and markdown body content.
//
// Fields:
//   - Name: Unique identifier from frontmatter
//   - Description: Used for auto inclusion matching
//   - Inclusion: Determines when the file is injected (default: manual)
//   - Tools: Tool identifiers this file applies to (nil means all tools)
//   - Priority: Injection order; lower values are injected first (default: 50)
//   - Body: Markdown content after frontmatter
//   - Path: Filesystem path to the steering file
type SteeringFile struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description,omitempty"`
	Inclusion   InclusionMode `yaml:"inclusion"`
	Tools       []string      `yaml:"tools,omitempty"`
	Priority    int           `yaml:"priority"`
	Body        string        `yaml:"-"`
	Path        string        `yaml:"-"`
}

// SyncReport summarizes the result of syncing steering files to
// tool-native formats.
//
// Fields:
//   - Written: Files that were written or updated
//   - Skipped: Files that were skipped (unchanged or excluded)
//   - Errors: Errors encountered during sync
type SyncReport struct {
	Written []string
	Skipped []string
	Errors  []error
}

// cursorFrontmatter is the YAML frontmatter for Cursor rule files.
type cursorFrontmatter struct {
	Description string `yaml:"description"`
	Globs       []any  `yaml:"globs"`
	AlwaysApply bool   `yaml:"alwaysApply"`
}

// kiroFrontmatter is the YAML frontmatter for Kiro steering files.
type kiroFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Mode        string `yaml:"mode"`
}

// FoundationFile describes a foundation steering file to generate.
type FoundationFile struct {
	Name        string
	Description string
	Body        string
}

// FoundationFiles defines the set of files created by ctx steering init.
var FoundationFiles = []FoundationFile{
	{
		Name:        "product",
		Description: "Product context, goals, and target users",
		Body: "# Product Context\n\n" +
			"Describe the product, its goals, and target users.\n",
	},
	{
		Name:        "tech",
		Description: "Technology stack, constraints, and dependencies",
		Body: "# Technology Stack\n\n" +
			"Describe the technology stack, " +
			"constraints, and key dependencies.\n",
	},
	{
		Name:        "structure",
		Description: "Project structure and directory conventions",
		Body: "# Project Structure\n\n" +
			"Describe the project layout " +
			"and directory conventions.\n",
	},
	{
		Name:        "workflow",
		Description: "Development workflow and process rules",
		Body: "# Development Workflow\n\n" +
			"Describe the development workflow, " +
			"branching strategy, and process rules.\n",
	},
}
