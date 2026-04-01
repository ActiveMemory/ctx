//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
)

// HookScript is the prepare-commit-msg hook content installed by ctx.
const HookScript = `#!/bin/sh
# ctx: prepare-commit-msg hook for commit context tracing.
# Installed by: ctx trace hook enable
# Remove with:  ctx trace hook disable

COMMIT_MSG_FILE="$1"
COMMIT_SOURCE="$2"

# Only inject on normal commits (not merges, squashes, or amends)
case "$COMMIT_SOURCE" in
  merge|squash) exit 0 ;;
esac

# Ensure ctx is available
command -v ctx >/dev/null 2>&1 || exit 0

# Collect context refs
TRAILER=$(ctx trace collect 2>/dev/null)

if [ -n "$TRAILER" ]; then
  # Append trailer with a blank line separator
  echo "" >> "$COMMIT_MSG_FILE"
  echo "$TRAILER" >> "$COMMIT_MSG_FILE"
fi
`

// PostCommitScript is the post-commit hook content installed by ctx.
const PostCommitScript = `#!/bin/sh
# ctx: post-commit hook for recording commit context history.
# Installed by: ctx trace hook enable
# Remove with:  ctx trace hook disable

# Ensure ctx is available
command -v ctx >/dev/null 2>&1 || exit 0

COMMIT_HASH=$(git rev-parse HEAD)
ctx trace collect --record "$COMMIT_HASH" 2>/dev/null || true
`

// Enable installs both the prepare-commit-msg and post-commit hooks.
//
// Parameters:
//   - cmd: Cobra command for output stream
//
// Returns:
//   - error: non-nil on installation failure
func Enable(cmd *cobra.Command) error {
	prepPath, prepErr := HookFilePath("prepare-commit-msg")
	if prepErr != nil {
		return prepErr
	}
	if installErr := InstallHook(prepPath, HookScript, "prepare-commit-msg"); installErr != nil {
		return installErr
	}

	postPath, postErr := HookFilePath("post-commit")
	if postErr != nil {
		return postErr
	}
	if installErr := InstallHook(postPath, PostCommitScript, "post-commit"); installErr != nil {
		return installErr
	}

	cmd.Println("ctx trace hooks enabled (prepare-commit-msg, post-commit)")
	return nil
}

// Disable removes both the prepare-commit-msg and post-commit hooks if they
// were installed by ctx.
//
// Parameters:
//   - cmd: Cobra command for output stream
//
// Returns:
//   - error: non-nil on removal failure
func Disable(cmd *cobra.Command) error {
	prepPath, err := HookFilePath("prepare-commit-msg")
	if err != nil {
		return err
	}
	RemoveHook(prepPath)

	postPath, err := HookFilePath("post-commit")
	if err != nil {
		return err
	}
	RemoveHook(postPath)

	cmd.Println("ctx trace hooks disabled")
	return nil
}

// InstallHook writes the hook script to path, checking for existing non-ctx hooks.
//
// Parameters:
//   - path: absolute path to the hook file
//   - script: hook script content to write
//   - name: hook name for error messages
//
// Returns:
//   - error: non-nil if a non-ctx hook already exists or write fails
func InstallHook(path, script, name string) error {
	if _, err := os.Stat(path); err == nil {
		//nolint:gosec // path from HookFilePath which calls git rev-parse
		existing, readErr := os.ReadFile(path)
		if readErr == nil && !strings.Contains(string(existing), "ctx trace") {
			return fmt.Errorf("%s hook already exists at %s (not installed by ctx); remove it manually first", name, path)
		}
	}
	if err := os.WriteFile(path, []byte(script), cfgFs.PermExec); err != nil {
		return fmt.Errorf("write %s hook: %w", name, err)
	}
	return nil
}

// RemoveHook removes the hook at path if it was installed by ctx.
//
// Parameters:
//   - path: absolute path to the hook file
func RemoveHook(path string) {
	//nolint:gosec // path from HookFilePath which calls git rev-parse
	existing, err := os.ReadFile(path)
	if err != nil {
		return
	}
	if strings.Contains(string(existing), "ctx trace") {
		_ = os.Remove(path)
	}
}

// HookFilePath returns the absolute path to a git hook by name.
//
// Parameters:
//   - hookName: name of the git hook (e.g. "prepare-commit-msg")
//
// Returns:
//   - string: absolute path to the hook file
//   - error: non-nil if git rev-parse fails
func HookFilePath(hookName string) (string, error) {
	out, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse --git-dir: %w", err)
	}
	gitDir := strings.TrimSpace(string(out))
	return filepath.Join(gitDir, "hooks", hookName), nil
}
