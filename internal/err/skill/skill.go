//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// CreateDest wraps a failure to create the skill destination directory.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "skill: create destination: <cause>"
func CreateDest(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillCreateDest), cause,
	)
}

// Install wraps a skill installation copy failure.
//
// Parameters:
//   - name: skill name
//   - cause: the underlying error
//
// Returns:
//   - error: "skill: install <name>: <cause>"
func Install(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillInstall), name, cause,
	)
}

// InvalidManifest wraps an invalid skill manifest parse failure.
//
// Parameters:
//   - manifest: manifest filename
//   - cause: the underlying error
//
// Returns:
//   - error: "skill: source has invalid <manifest>: <cause>"
func InvalidManifest(manifest string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillInvalidManifest), manifest, cause,
	)
}

// InvalidYAML wraps a YAML parse failure in a skill manifest.
//
// Parameters:
//   - name: skill name
//   - cause: the underlying YAML error
//
// Returns:
//   - error: "skill: <name>: invalid YAML frontmatter: <cause>"
func InvalidYAML(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillInvalidYAML), name, cause,
	)
}

// List wraps a failure to list embedded skill directories.
//
// Parameters:
//   - cause: the underlying error from the list operation
//
// Returns:
//   - error: "failed to list skills: <cause>"
func List(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrSkillList), cause)
}

// Load wraps a failure to load a skill by name.
//
// Parameters:
//   - name: skill name
//   - cause: the underlying error
//
// Returns:
//   - error: "skill: <name>: <cause>"
func Load(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillLoad), name, cause,
	)
}

// MissingName returns an error for a skill manifest missing the
// required name field.
//
// Parameters:
//   - manifest: the manifest filename
//
// Returns:
//   - error: "skill: <manifest> is missing required 'name' field"
func MissingName(manifest string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillMissingName), manifest,
	)
}

// NotFound returns an error when a skill cannot be found by name.
//
// Parameters:
//   - name: the skill name
//
// Returns:
//   - error: "skill <name> not found"
func NotFound(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillNotFound), name,
	)
}

// NotValidDir returns an error when a skill path is not a directory.
//
// Parameters:
//   - name: the skill name
//
// Returns:
//   - error: "skill: <name> is not a valid skill directory"
func NotValidDir(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillNotValidDir), name,
	)
}

// NotValidSource wraps an error when the source is not a valid skill.
//
// Parameters:
//   - cause: the underlying read error
//
// Returns:
//   - error: "skill: source is not a valid skill: <cause>"
func NotValidSource(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillNotValidSource), cause,
	)
}

// Read wraps a failure to read a skill's content.
//
// Parameters:
//   - name: Skill directory name that failed to read
//   - cause: the underlying error from the read operation
//
// Returns:
//   - error: "failed to read skill <name>: <cause>"
func Read(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillRead), name, cause,
	)
}

// ReadDir wraps a failure to read the skills directory.
//
// Parameters:
//   - dir: the directory path
//   - cause: the underlying error
//
// Returns:
//   - error: "skill: read directory <dir>: <cause>"
func ReadDir(dir string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillReadDir), dir, cause,
	)
}

// Remove wraps a skill removal failure.
//
// Parameters:
//   - name: skill name
//   - cause: the underlying error
//
// Returns:
//   - error: "skill: remove <name>: <cause>"
func Remove(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillRemove), name, cause,
	)
}

// LoadQuoted wraps a skill load failure with the skill name
// quoted.
//
// Parameters:
//   - name: skill name
//   - cause: the underlying error
//
// Returns:
//   - error: "skill <name>: <cause>"
func LoadQuoted(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillSkillLoad), name, cause,
	)
}
