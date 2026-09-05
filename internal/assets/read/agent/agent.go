//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"io/fs"
	"path"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// CopilotInstructions reads the embedded Copilot instructions template.
//
// Returns:
//   - []byte: Template content from integrations/copilot-instructions.md
//   - error: Non-nil if the file is not found or read fails
func CopilotInstructions() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathCopilotInstructions)
}

// CopilotCLIHooksJSON reads the embedded Copilot CLI hooks config.
//
// Returns:
//   - []byte: JSON content from integrations/copilot-cli/ctx-hooks.json
//   - error: Non-nil if the file is not found or read fails
func CopilotCLIHooksJSON() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathCopilotCLIHooksJSON)
}

// AgentsMd reads the embedded AGENTS.md template.
//
// Returns:
//   - []byte: Template content from integrations/agents.md
//   - error: Non-nil if the file is not found or read fails
func AgentsMd() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathAgentsMd)
}

// AgentsCtxMd reads the embedded .github/agents/ctx.md template.
//
// Returns:
//   - []byte: Template content from integrations/copilot-cli/agents-ctx.md
//   - error: Non-nil if the file is not found or read fails
func AgentsCtxMd() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathAgentsCtxMd)
}

// InstructionsCtxMd reads the embedded path-specific instructions.
//
// Returns:
//   - []byte: Template content from
//     integrations/copilot-cli/instructions-context.md
//   - error: Non-nil if the file is not found or read fails
func InstructionsCtxMd() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathInstructionsCtxMd)
}

// OpenCodePlugin reads all embedded OpenCode plugin files.
// Returns a map of filename to content for files in
// integrations/opencode/plugin/.
//
// Returns:
//   - map[string][]byte: Filename -> content for each plugin file
//   - error: Non-nil if the directory read fails
func OpenCodePlugin() (map[string][]byte, error) {
	files := make(map[string][]byte)
	entries, dirErr := fs.ReadDir(
		assets.FS, asset.DirIntegrationsOpenCodePlugin)
	if dirErr != nil {
		return nil, dirErr
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		p := path.Join(asset.DirIntegrationsOpenCodePlugin, name)
		content, readErr := assets.FS.ReadFile(p)
		if readErr != nil {
			return nil, readErr
		}
		files[name] = content
	}
	return files, nil
}

// OpenCodeSkills reads all embedded OpenCode skill templates.
// Returns a map of skill directory name to SKILL.md content for skills
// in integrations/opencode/skills/.
//
// Returns:
//   - map[string][]byte: Skill name -> SKILL.md content
//   - error: Non-nil if the directory read fails
func OpenCodeSkills() (map[string][]byte, error) {
	skills := make(map[string][]byte)
	entries, dirErr := fs.ReadDir(
		assets.FS, asset.DirIntegrationsOpenCodeSkill)
	if dirErr != nil {
		return nil, dirErr
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		skillPath := path.Join(
			asset.DirIntegrationsOpenCodeSkill,
			name, asset.FileSKILLMd)
		content, readErr := assets.FS.ReadFile(skillPath)
		if readErr != nil {
			return nil, readErr
		}
		skills[name] = content
	}
	return skills, nil
}

// CopilotCLISkills reads all embedded Copilot CLI skill templates.
// Returns a map of skill directory name to SKILL.md content for skills
// in integrations/copilot-cli/skills/.
//
// Returns:
//   - map[string][]byte: Skill name -> SKILL.md content
//   - error: Non-nil if the directory read fails
func CopilotCLISkills() (map[string][]byte, error) {
	skills := make(map[string][]byte)
	entries, dirErr := fs.ReadDir(assets.FS, asset.DirIntegrationsCopilotSkill)
	if dirErr != nil {
		return nil, dirErr
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		skillPath := path.Join(
			asset.DirIntegrationsCopilotSkill,
			name, asset.FileSKILLMd)
		content, readErr := assets.FS.ReadFile(skillPath)
		if readErr != nil {
			return nil, readErr
		}
		skills[name] = content
	}
	return skills, nil
}

// CodexHooksJSON reads the embedded Codex hooks manifest. The same
// content serves the plugin (`hooks/hooks.json`) and the
// project-local deploy (`.codex/hooks.json`).
//
// Returns:
//   - []byte: JSON content from codex/hooks/hooks.json
//   - error: Non-nil if the file is not found or read fails
func CodexHooksJSON() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathCodexHooksJSON)
}

// CodexSkills reads all embedded Codex skill templates.
// Returns a map of skill directory name to SKILL.md content for
// skills in codex/skills/.
//
// Returns:
//   - map[string][]byte: Skill name -> SKILL.md content
//   - error: Non-nil if the directory read fails
func CodexSkills() (map[string][]byte, error) {
	skills := make(map[string][]byte)
	entries, dirErr := fs.ReadDir(assets.FS, asset.DirCodexSkills)
	if dirErr != nil {
		return nil, dirErr
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		skillPath := path.Join(
			asset.DirCodexSkills, name, asset.FileSKILLMd)
		content, readErr := assets.FS.ReadFile(skillPath)
		if readErr != nil {
			return nil, readErr
		}
		skills[name] = content
	}
	return skills, nil
}

// CodexSkillReferences reads the embedded reference files of every
// Codex skill. Keys are skill names; values map a reference file
// name to its content. Skills without a references directory are
// absent from the map.
//
// Returns:
//   - map[string]map[string][]byte: Skill -> reference file -> content
//   - error: Non-nil if a read fails
func CodexSkillReferences() (map[string]map[string][]byte, error) {
	refs := make(map[string]map[string][]byte)
	entries, dirErr := fs.ReadDir(assets.FS, asset.DirCodexSkills)
	if dirErr != nil {
		return nil, dirErr
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		refDir := path.Join(
			asset.DirCodexSkills, name, asset.DirReferences)
		refEntries, refErr := fs.ReadDir(assets.FS, refDir)
		if refErr != nil {
			// No references directory for this skill.
			continue
		}
		for _, ref := range refEntries {
			if ref.IsDir() {
				continue
			}
			content, readErr := assets.FS.ReadFile(
				path.Join(refDir, ref.Name()))
			if readErr != nil {
				return nil, readErr
			}
			if refs[name] == nil {
				refs[name] = make(map[string][]byte)
			}
			refs[name][ref.Name()] = content
		}
	}
	return refs, nil
}
