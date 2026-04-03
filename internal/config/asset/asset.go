//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package asset

import "path"

// Embedded asset directory names.
const (
	DirClaude                   = "claude"
	DirClaudePlugin             = "claude/.claude-plugin"
	DirClaudeSkills             = "claude/skills"
	DirCommands                 = "commands"
	DirCommandsText             = "commands/text"
	DirContext                  = "context"
	DirEntryTemplates           = "entry-templates"
	DirIntegrations             = "integrations"
	DirIntegrationsCopilot      = "integrations/copilot"
	DirIntegrationsCopilotCLI   = "integrations/copilot-cli"
	DirIntegrationsCopilotScrp  = "integrations/copilot-cli/scripts"
	DirIntegrationsCopilotSkill = "integrations/copilot-cli/skills"
	DirHooksMessages            = "hooks/messages"
	DirHooksTrace               = "hooks/trace"
	DirJournal                  = "journal"
	DirPermissions              = "permissions"
	DirProject                  = "project"
	DirSchema                   = "schema"
	DirWhy                      = "why"
)

// JSON field keys used when parsing embedded asset files.
const (
	JSONKeyVersion = "version"
)

// Naming patterns used to construct embedded asset filenames.
const (
	// SuffixReadme is appended to a directory name to form the README
	// template filename (e.g., "specs" -> "specs-README.md").
	SuffixReadme = "-README.md"
)

// Embedded asset file names (leaf names only).
const (
	FileAllowTxt              = "allow.txt"
	FileCLAUDEMd              = "CLAUDE.md"
	FileCommandsYAML          = "commands.yaml"
	FileAgentsMd              = "agents.md"
	FileAgentsCtxMd           = "agents-ctx.md"
	FileCopilotCLIHooksJSON   = "ctx-hooks.json"
	FileCopilotInstructionsMd = "copilot-instructions.md"
	FileCtxrcSchemaJSON       = "ctxrc.schema.json"
	FileDenyTxt               = "deny.txt"
	FileExamplesYAML          = "examples.yaml"
	FileExtraCSS              = "extra.css"
	FileFlagsYAML             = "flags.yaml"
	FileMakefileCtx           = "Makefile.ctx"
	FilePluginJSON            = "plugin.json"
	FileRegistryYAML          = "registry.yaml"
	FileSKILLMd               = "SKILL.md"
)

// Subdirectory name within a skill directory.
const (
	DirReferences = "references"
)

// Full embedded paths for files accessed directly.
var (
	PathCLAUDEMd            = path.Join(DirClaude, FileCLAUDEMd)
	PathPluginJSON          = path.Join(DirClaudePlugin, FilePluginJSON)
	PathCommandsYAML        = path.Join(DirCommands, FileCommandsYAML)
	PathFlagsYAML           = path.Join(DirCommands, FileFlagsYAML)
	PathExamplesYAML        = path.Join(DirCommands, FileExamplesYAML)
	PathAgentsMd            = path.Join(DirIntegrations, FileAgentsMd)
	PathAgentsCtxMd         = path.Join(DirIntegrationsCopilotCLI, FileAgentsCtxMd)
	PathCopilotCLIHooksJSON = path.Join(
		DirIntegrationsCopilotCLI, FileCopilotCLIHooksJSON)
	PathCopilotInstructions = path.Join(
		DirIntegrationsCopilot, FileCopilotInstructionsMd)
	PathInstructionsCtxMd = path.Join(
		DirIntegrationsCopilotCLI, "instructions-context.md")
	PathMessageRegistry = path.Join(DirHooksMessages, FileRegistryYAML)
	PathExtraCSS        = path.Join(DirJournal, FileExtraCSS)
	PathMakefileCtx     = path.Join(DirProject, FileMakefileCtx)
	PathAllowTxt        = path.Join(DirPermissions, FileAllowTxt)
	PathDenyTxt         = path.Join(DirPermissions, FileDenyTxt)
	PathCtxrcSchema     = path.Join(DirSchema, FileCtxrcSchemaJSON)
)
