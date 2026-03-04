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
	"strings"
	"sync"
)

//go:embed claude/.claude-plugin/plugin.json claude/CLAUDE.md claude/skills/*/references/*.md claude/skills/*/SKILL.md context/*.md project/* entry-templates/*.md hooks/messages/*/*.txt hooks/messages/registry.yaml prompt-templates/*.md ralph/*.md schema/*.json why/*.md permissions/*.txt
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

// ListPromptTemplates returns available prompt template file names.
//
// Returns:
//   - []string: List of template filenames in prompt-templates/
//   - error: Non-nil if directory read fails
func ListPromptTemplates() ([]string, error) {
	entries, err := FS.ReadDir("prompt-templates")
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

// PromptTemplate reads a prompt template by name.
//
// Parameters:
//   - name: Template filename (e.g., "code-review.md")
//
// Returns:
//   - []byte: Template content from prompt-templates/
//   - error: Non-nil if the file is not found or read fails
func PromptTemplate(name string) ([]byte, error) {
	return FS.ReadFile("prompt-templates/" + name)
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

// SkillReference reads a reference file from a skill's references/ directory.
//
// Parameters:
//   - skill: Skill directory name (e.g., "ctx-skill-audit")
//   - filename: Reference filename (e.g., "anthropic-best-practices.md")
//
// Returns:
//   - []byte: Reference file content
//   - error: Non-nil if the file is not found or read fails
func SkillReference(skill, filename string) ([]byte, error) {
	return FS.ReadFile("claude/skills/" + skill + "/references/" + filename)
}

// ListSkillReferences returns available reference filenames for a skill.
//
// Parameters:
//   - skill: Skill directory name (e.g., "ctx-skill-audit")
//
// Returns:
//   - []string: List of reference filenames
//   - error: Non-nil if the references directory is not found or read fails
func ListSkillReferences(skill string) ([]string, error) {
	entries, readErr := FS.ReadDir("claude/skills/" + skill + "/references")
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

// ProjectReadme reads a project directory README template by directory name.
//
// Templates are stored as project/<dir>-README.md in the embedded filesystem.
//
// Parameters:
//   - dir: Directory name (e.g., "specs", "ideas")
//
// Returns:
//   - []byte: README.md content for the directory
//   - error: Non-nil if the file is not found or read fails
func ProjectReadme(dir string) ([]byte, error) {
	return FS.ReadFile("project/" + dir + "-README.md")
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

// HookMessageRegistry reads the embedded registry.yaml that describes
// all hook message templates.
//
// Returns:
//   - []byte: Raw YAML content
//   - error: Non-nil if the file is not found or read fails
func HookMessageRegistry() ([]byte, error) {
	return FS.ReadFile("hooks/messages/registry.yaml")
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

// Schema reads the embedded JSON Schema for .ctxrc.
//
// Returns:
//   - []byte: JSON Schema content
//   - error: Non-nil if the file is not found or read fails
func Schema() ([]byte, error) {
	return FS.ReadFile("schema/ctxrc.schema.json")
}

var (
	allowOnce  sync.Once
	allowPerms []string

	denyOnce  sync.Once
	denyPerms []string
)

// parsePermissions splits a text file into permission entries.
//
// Lines are trimmed; empty lines and lines starting with '#' are skipped.
func parsePermissions(data []byte) []string {
	var result []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		result = append(result, line)
	}
	return result
}

// DefaultAllowPermissions returns the default allow permissions for ctx
// commands and skills, parsed from the embedded permissions/allow.txt.
func DefaultAllowPermissions() []string {
	allowOnce.Do(func() {
		data, readErr := FS.ReadFile("permissions/allow.txt")
		if readErr != nil {
			return
		}
		allowPerms = parsePermissions(data)
	})
	return allowPerms
}

// DefaultDenyPermissions returns the default deny permissions that block
// dangerous operations, parsed from the embedded permissions/deny.txt.
func DefaultDenyPermissions() []string {
	denyOnce.Do(func() {
		data, readErr := FS.ReadFile("permissions/deny.txt")
		if readErr != nil {
			return
		}
		denyPerms = parsePermissions(data)
	})
	return denyPerms
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
