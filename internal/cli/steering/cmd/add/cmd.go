//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errSteering "github.com/ActiveMemory/ctx/internal/err/steering"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/steering"
	writeSteering "github.com/ActiveMemory/ctx/internal/write/steering"
)

// defaultPriority is the default priority for new steering files.
const defaultPriority = 50

// Cmd returns the "ctx steering add" subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySteeringAdd)

	return &cobra.Command{
		Use:   cmd.UseSteeringAdd,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return Run(c, args[0])
		},
	}
}

// Run creates a new steering file with default frontmatter.
//
// Parameters:
//   - c: The cobra command for output
//   - name: The steering file name (without .md extension)
func Run(c *cobra.Command, name string) error {
	contextDir := rc.ContextDir()

	// Check that .context/ directory exists.
	if _, statErr := ctxIo.SafeStat(contextDir); os.IsNotExist(statErr) {
		return errSteering.ContextDirMissing()
	}

	steeringDir := rc.SteeringDir()

	// Ensure the steering directory exists.
	if mkdirErr := ctxIo.SafeMkdirAll(
		steeringDir, fs.PermRestrictedDir,
	); mkdirErr != nil {
		return errSteering.CreateDir(mkdirErr)
	}

	filePath := filepath.Join(
		steeringDir, name+file.ExtMarkdown,
	)

	// Error if file already exists.
	if _, statErr := ctxIo.SafeStat(filePath); statErr == nil {
		return errSteering.FileExists(filePath)
	}

	sf := &steering.SteeringFile{
		Name:      name,
		Inclusion: steering.InclusionManual,
		Priority:  defaultPriority,
	}

	data := steering.Print(sf)
	if writeErr := ctxIo.SafeWriteFile(
		filePath, data, fs.PermFile,
	); writeErr != nil {
		return errSteering.Write(writeErr)
	}

	writeSteering.Created(c, filePath)
	return nil
}
