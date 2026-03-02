//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	internalConfig "github.com/ActiveMemory/ctx/internal/config"
)

// Profile file names and identifiers.
const (
	fileCtxRC     = ".ctxrc"
	fileCtxRCBase = ".ctxrc.base"
	fileCtxRCDev  = ".ctxrc.dev"

	profileDev  = "dev"
	profileBase = "base"
)

func switchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "switch [dev|base]",
		Short: "Switch .ctxrc profile",
		Long: `Switch between .ctxrc configuration profiles.

With no argument, toggles between dev and base.
Accepts "prod" as an alias for "base".

Source files (.ctxrc.base, .ctxrc.dev) are committed to git.
The working copy (.ctxrc) is gitignored.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, rootErr := gitRoot()
			if rootErr != nil {
				return rootErr
			}
			return runSwitch(cmd, root, args)
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show active .ctxrc profile",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, rootErr := gitRoot()
			if rootErr != nil {
				return rootErr
			}
			return runStatus(cmd, root)
		},
	}
}

func runSwitch(cmd *cobra.Command, root string, args []string) error {
	var target string
	if len(args) > 0 {
		target = args[0]
	}

	// Normalize "prod" alias.
	if target == "prod" {
		target = profileBase
	}

	switch target {
	case profileDev:
		return switchTo(cmd, root, profileDev)
	case profileBase:
		return switchTo(cmd, root, profileBase)
	case "":
		// Toggle.
		current := detectProfile(root)
		if current == profileDev {
			return switchTo(cmd, root, profileBase)
		}
		return switchTo(cmd, root, profileDev)
	default:
		return fmt.Errorf(
			"unknown profile %q: must be dev, base, or prod", target)
	}
}

func switchTo(cmd *cobra.Command, root, profile string) error {
	current := detectProfile(root)
	if current == profile {
		cmd.Printf("already on %s profile\n", profile)
		return nil
	}

	var srcFile string
	if profile == profileDev {
		srcFile = fileCtxRCDev
	} else {
		srcFile = fileCtxRCBase
	}

	if copyErr := copyProfile(root, srcFile); copyErr != nil {
		return copyErr
	}

	if current == "" {
		cmd.Printf("created %s from %s profile\n", fileCtxRC, profile)
	} else {
		cmd.Printf("switched to %s profile\n", profile)
	}
	return nil
}

func runStatus(cmd *cobra.Command, root string) error {
	profile := detectProfile(root)
	switch profile {
	case profileDev:
		cmd.Println("active: dev (verbose logging enabled)")
	case profileBase:
		cmd.Println("active: base (defaults)")
	default:
		cmd.Printf("active: none (%s does not exist)\n", fileCtxRC)
	}
	return nil
}

// detectProfile reads .ctxrc and returns "dev" or "base" based on the
// presence of an uncommented "notify:" line. Returns "" if the file is missing.
func detectProfile(root string) string {
	data, readErr := os.ReadFile(filepath.Join(root, fileCtxRC)) //nolint:gosec // project-local config file
	if readErr != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "notify:") {
			return profileDev
		}
	}
	return profileBase
}

// copyProfile copies a source profile file to .ctxrc.
func copyProfile(root, srcFile string) error {
	src := filepath.Join(root, srcFile)
	data, readErr := os.ReadFile(src) //nolint:gosec // project-local file
	if readErr != nil {
		return fmt.Errorf("read %s: %w", srcFile, readErr)
	}

	dst := filepath.Join(root, fileCtxRC)
	return os.WriteFile(dst, data, internalConfig.PermFile)
}

// gitRoot returns the git repository root directory.
func gitRoot() (string, error) {
	out, execErr := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if execErr != nil {
		return "", fmt.Errorf("not in a git repository: %w", execErr)
	}
	return strings.TrimSpace(string(out)), nil
}
