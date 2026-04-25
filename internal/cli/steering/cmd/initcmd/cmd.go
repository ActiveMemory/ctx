//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initcmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgSteering "github.com/ActiveMemory/ctx/internal/config/steering"
	errSteering "github.com/ActiveMemory/ctx/internal/err/steering"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/steering"
	writeSteering "github.com/ActiveMemory/ctx/internal/write/steering"
)

// Cmd returns the "ctx steering init" subcommand.
//
// Returns:
//   - *cobra.Command: Configured init subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySteeringInit)

	return &cobra.Command{
		Use:     cmd.UseSteeringInit,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySteeringInit),
		Args:    cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			return Run(c)
		},
	}
}

// Run generates foundation steering files in the steering directory.
// Existing files are skipped and reported.
//
// Parameters:
//   - c: The cobra command for output
//
// Returns:
//   - error: nil on success, or if the context directory is missing
func Run(c *cobra.Command) error {
	contextDir, err := rc.RequireContextDir()
	if err != nil {
		c.SilenceUsage = true
		return err
	}
	return RunWithDir(c, contextDir)
}

// RunWithDir is the implementation of Run that accepts an explicit
// context directory. Used by `ctx init`, which has just created the
// directory and needs to scaffold foundation steering files without
// requiring the user to have declared CTX_DIR first.
//
// Parameters:
//   - c: The cobra command for output
//   - contextDir: absolute path to the .context/ directory
//
// Returns:
//   - error: nil on success, or a file creation error
func RunWithDir(c *cobra.Command, contextDir string) error {
	// Check that .context/ directory exists.
	if _, statErr := ctxIo.SafeStat(
		contextDir,
	); os.IsNotExist(statErr) {
		return errSteering.ContextDirMissing()
	}

	steeringDir := filepath.Join(contextDir, dir.Steering)

	// Ensure the steering directory exists.
	if mkdirErr := ctxIo.SafeMkdirAll(
		steeringDir, fs.PermRestrictedDir,
	); mkdirErr != nil {
		return errSteering.CreateDir(mkdirErr)
	}

	var created, skipped int

	for _, ff := range steering.FoundationFiles() {
		filePath := filepath.Join(
			steeringDir, ff.Name+file.ExtMarkdown,
		)

		if _, statErr := ctxIo.SafeStat(
			filePath,
		); statErr == nil {
			writeSteering.Skipped(c, filePath)
			skipped++
			continue
		}

		sf := &steering.SteeringFile{
			Name:        ff.Name,
			Description: ff.Description,
			Inclusion:   cfgSteering.InclusionAlways,
			Priority:    10,
			Body:        ff.Body,
		}

		data := steering.Print(sf)
		if writeErr := ctxIo.SafeWriteFile(
			filePath, data, fs.PermFile,
		); writeErr != nil {
			return errSteering.WriteInitFile(
				filePath, writeErr,
			)
		}

		writeSteering.Created(c, filePath)
		created++
	}

	writeSteering.InitSummary(c, created, skipped)
	return nil
}
