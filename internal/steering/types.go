//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgSteering "github.com/ActiveMemory/ctx/internal/config/steering"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// SteeringFile describes a parsed steering file with YAML
// frontmatter and Markdown body content.
//
// Fields:
//   - Name: Unique identifier from frontmatter
//   - Description: Used for auto-inclusion matching
//   - Inclusion: Determines when the file is injected
//     (default: manual)
//   - Tools: Tool identifiers this file applies to
//     (nil means all tools)
//   - Priority: Injection order; lower values are injected
//     first (default: 50)
//   - Body: Markdown content after frontmatter
//   - Path: Filesystem path to the steering file
type SteeringFile struct {
	Name        string                    `yaml:"name"`
	Description string                    `yaml:"description,omitempty"`
	Inclusion   cfgSteering.InclusionMode `yaml:"inclusion"`
	Tools       []string                  `yaml:"tools,omitempty"`
	Priority    int                       `yaml:"priority"`
	Body        string                    `yaml:"-"`
	Path        string                    `yaml:"-"`
}

// SyncReport summarizes the result of syncing steering
// files to tool-native formats.
//
// Fields:
//   - Written: Files that were written or updated
//   - Skipped: Files that were skipped (unchanged or
//     excluded)
//   - Errors: Errors encountered during sync
type SyncReport struct {
	Written []string
	Skipped []string
	Errors  []error
}

// cursorFrontmatter holds YAML frontmatter for Cursor rule
// files.
type cursorFrontmatter struct {
	Description string `yaml:"description"`
	Globs       []any  `yaml:"globs"`
	AlwaysApply bool   `yaml:"alwaysApply"`
}

// kiroFrontmatter holds YAML frontmatter for Kiro steering
// files.
type kiroFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Mode        string `yaml:"mode"`
}

// FoundationFile describes a foundation steering file to
// generate during ctx steering init.
//
// Fields:
//   - Name: Unique steering file identifier
//   - Description: Human-readable purpose summary
//   - Body: Markdown content for the file
type FoundationFile struct {
	Name        string
	Description string
	Body        string
}

// FoundationFiles returns the set of files created by
// ctx steering init. Descriptions and bodies are loaded
// from YAML text assets at call time.
//
// Returns:
//   - []FoundationFile: Foundation steering files with
//     names, descriptions, and body content
func FoundationFiles() []FoundationFile {
	guidance := desc.Text(text.DescKeyWriteSteeringGuidance) +
		token.NewlineLF + token.NewlineLF
	return []FoundationFile{
		{
			Name:        cfgSteering.NameProduct,
			Description: desc.Text(text.DescKeyWriteSteeringDescProduct),
			Body:        guidance + desc.Text(text.DescKeyWriteSteeringBodyProduct),
		},
		{
			Name:        cfgSteering.NameTech,
			Description: desc.Text(text.DescKeyWriteSteeringDescTech),
			Body:        guidance + desc.Text(text.DescKeyWriteSteeringBodyTech),
		},
		{
			Name:        cfgSteering.NameStructure,
			Description: desc.Text(text.DescKeyWriteSteeringDescStructure),
			Body:        guidance + desc.Text(text.DescKeyWriteSteeringBodyStructure),
		},
		{
			Name:        cfgSteering.NameWorkflow,
			Description: desc.Text(text.DescKeyWriteSteeringDescWorkflow),
			Body:        guidance + desc.Text(text.DescKeyWriteSteeringBodyWorkflow),
		},
	}
}
