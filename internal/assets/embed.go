//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package assets provides embedded assets for ctx: .context/ templates
// stamped by "ctx init" and the Claude Code plugin (skills, hooks,
// manifest) served directly from claude/.
package assets

import (
	"embed"
	"encoding/json"
)

//go:embed context/*.md project/* claude/CLAUDE.md entry-templates/*.md claude/skills/*/SKILL.md claude/.claude-plugin/plugin.json ralph/*.md hooks/messages/*/*.txt why/*.md
var FS embed.FS

// Template reads a template file by name from the embedded filesystem.
//
// Parameters:
//   - name: Template filename (e.g., "TASKS.md")
//
// Returns:
//   - []byte: Template content
//   - error: Non-nil if the file is not found or read fails
func Template(name string) ([]byte, error) {
	return FS.ReadFile("context/" + name)
}

// List returns all available template file names.
//
// Returns:
//   - []string: List of template filenames in the root templates directory
//   - error: Non-nil if directory read fails
func List() ([]string, error) {
	entries, err := FS.ReadDir("context")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// ListEntry returns available entry template file names.
//
// Returns:
//   - []string: List of template filenames in entry-templates/
//   - error: Non-nil if directory read fails
func ListEntry() ([]string, error) {
	entries, err := FS.ReadDir("entry-templates")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// Entry reads an entry template by name.
//
// Parameters:
//   - name: Template filename (e.g., "decision.md")
//
// Returns:
//   - []byte: Template content from entry-templates/
//   - error: Non-nil if the file is not found or read fails
func Entry(name string) ([]byte, error) {
	return FS.ReadFile("entry-templates/" + name)
}

// ListSkills returns available skill directory names.
//
// Each skill is a directory containing a SKILL.md file following the
// Agent Skills specification (https://agentskills.io/specification).
//
// Returns:
//   - []string: List of skill directory names in claude/skills/
//   - error: Non-nil if directory read fails
func ListSkills() ([]string, error) {
	entries, err := FS.ReadDir("claude/skills")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// SkillContent reads a skill's SKILL.md file by skill name.
//
// Parameters:
//   - name: Skill directory name (e.g., "ctx-status")
//
// Returns:
//   - []byte: SKILL.md content from claude/skills/<name>/
//   - error: Non-nil if the file not found or read fails
func SkillContent(name string) ([]byte, error) {
	return FS.ReadFile("claude/skills/" + name + "/SKILL.md")
}

// MakefileCtx reads the ctx-owned Makefile include template.
//
// Returns:
//   - []byte: Makefile.ctx content
//   - error: Non-nil if the file is not found or read fails
func MakefileCtx() ([]byte, error) {
	return FS.ReadFile("project/Makefile.ctx")
}

// ProjectFile reads a project-root file by name from the embedded filesystem.
//
// These files are deployed to the project root (not .context/) by dedicated
// handlers during initialization.
//
// Parameters:
//   - name: Filename (e.g., "IMPLEMENTATION_PLAN.md")
//
// Returns:
//   - []byte: File content
//   - error: Non-nil if the file is not found or read fails
func ProjectFile(name string) ([]byte, error) {
	return FS.ReadFile("project/" + name)
}

// ClaudeMd reads the CLAUDE.md template from the embedded filesystem.
//
// CLAUDE.md is deployed to the project root by a dedicated handler
// during initialization, separate from the .context/ templates.
//
// Returns:
//   - []byte: CLAUDE.md content
//   - error: Non-nil if the file is not found or read fails
func ClaudeMd() ([]byte, error) {
	return FS.ReadFile("claude/CLAUDE.md")
}

// RalphTemplate reads a Ralph-mode template file by name.
//
// Ralph mode templates are designed for autonomous loop operation,
// with instructions for one-task-per-iteration, completion signals,
// and no clarifying questions.
//
// Parameters:
//   - name: Template filename (e.g., "PROMPT.md")
//
// Returns:
//   - []byte: Template content from ralph/
//   - error: Non-nil if the file is not found or read fails
func RalphTemplate(name string) ([]byte, error) {
	return FS.ReadFile("ralph/" + name)
}

// HookMessage reads a hook message template by hook name and filename.
//
// Parameters:
//   - hook: Hook directory name (e.g., "qa-reminder")
//   - filename: Template filename (e.g., "gate.txt")
//
// Returns:
//   - []byte: Template content from hooks/messages/<hook>/
//   - error: Non-nil if the file is not found or read fails
func HookMessage(hook, filename string) ([]byte, error) {
	return FS.ReadFile("hooks/messages/" + hook + "/" + filename)
}

// ListHookMessages returns available hook message directory names.
//
// Each hook is a directory under hooks/messages/ containing one or
// more variant .txt template files.
//
// Returns:
//   - []string: List of hook directory names
//   - error: Non-nil if directory read fails
func ListHookMessages() ([]string, error) {
	entries, readErr := FS.ReadDir("hooks/messages")
	if readErr != nil {
		return nil, readErr
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// ListHookVariants returns available variant filenames for a hook.
//
// Parameters:
//   - hook: Hook directory name (e.g., "qa-reminder")
//
// Returns:
//   - []string: List of variant filenames (e.g., "gate.txt")
//   - error: Non-nil if the hook directory is not found or read fails
func ListHookVariants(hook string) ([]string, error) {
	entries, readErr := FS.ReadDir("hooks/messages/" + hook)
	if readErr != nil {
		return nil, readErr
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// WhyDoc reads a "why" document by name from the embedded filesystem.
//
// Parameters:
//   - name: Document name (e.g., "manifesto", "about", "design-invariants")
//
// Returns:
//   - []byte: Document content from why/
//   - error: Non-nil if the file is not found or read fails
func WhyDoc(name string) ([]byte, error) {
	return FS.ReadFile("why/" + name + ".md")
}

// ListWhyDocs returns available "why" document names (without extension).
//
// Returns:
//   - []string: List of document names in why/
//   - error: Non-nil if directory read fails
func ListWhyDocs() ([]string, error) {
	entries, readErr := FS.ReadDir("why")
	if readErr != nil {
		return nil, readErr
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if len(name) > 3 && name[len(name)-3:] == ".md" {
				names = append(names, name[:len(name)-3])
			}
		}
	}
	return names, nil
}

// PluginVersion returns the version string from the embedded plugin.json.
func PluginVersion() (string, error) {
	data, err := FS.ReadFile("claude/.claude-plugin/plugin.json")
	if err != nil {
		return "", err
	}
	var manifest struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", err
	}
	return manifest.Version, nil
}
