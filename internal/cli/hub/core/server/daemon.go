//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	cfgFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	errServe "github.com/ActiveMemory/ctx/internal/err/serve"
	execDaemon "github.com/ActiveMemory/ctx/internal/exec/daemon"
	"github.com/ActiveMemory/ctx/internal/io"
	writeServe "github.com/ActiveMemory/ctx/internal/write/serve"
)

// RunDaemon starts the hub server as a background process.
//
// Writes a PID file to <dataDir>/hub.pid for later
// stopping via Stop.
//
// Parameters:
//   - cmd: cobra command for output
//   - port: TCP port to listen on
//   - dataDir: hub data directory (empty = default)
//
// Returns:
//   - error: non-nil if fork or PID file write fails
func RunDaemon(
	cmd *cobra.Command, port int, dataDir string,
) error {
	if dataDir == "" {
		defaultDir, dirErr := defaultDataDir()
		if dirErr != nil {
			return dirErr
		}
		dataDir = defaultDir
	}

	binPath, lookErr := os.Executable()
	if lookErr != nil {
		return lookErr
	}

	args := []string{
		cfgHub.ArgHub, cfgHub.ArgStart,
		cfgHub.FmtFlagPrefix + cfgFlag.Port, strconv.Itoa(port),
		cfgHub.FmtFlagPrefix + cfgFlag.DataDir, dataDir,
	}

	pid, startErr := execDaemon.Start(binPath, args)
	if startErr != nil {
		return startErr
	}

	pidPath := filepath.Join(dataDir, cfgHub.FilePID)
	if writeErr := io.SafeWriteFile(
		pidPath,
		[]byte(strconv.Itoa(pid)),
		fs.PermFile,
	); writeErr != nil {
		return writeErr
	}

	writeServe.Daemonized(cmd, pid)
	return nil
}

// Stop kills a running hub daemon via PID file.
//
// Parameters:
//   - cmd: cobra command for output
//   - dataDir: hub data directory (empty = default)
//
// Returns:
//   - error: non-nil if PID file missing or kill fails
func Stop(cmd *cobra.Command, dataDir string) error {
	if dataDir == "" {
		defaultDir, dirErr := defaultDataDir()
		if dirErr != nil {
			return dirErr
		}
		dataDir = defaultDir
	}

	pidPath := filepath.Join(dataDir, cfgHub.FilePID)
	data, readErr := io.SafeReadUserFile(pidPath)
	if readErr != nil {
		return errServe.NoRunningHub(readErr)
	}

	pid, parseErr := strconv.Atoi(
		strings.TrimSpace(string(data)),
	)
	if parseErr != nil {
		return errServe.InvalidPID(parseErr)
	}

	proc, findErr := os.FindProcess(pid)
	if findErr != nil {
		return findErr
	}

	if killErr := proc.Signal(
		syscall.SIGTERM,
	); killErr != nil {
		return errServe.Kill(pid, killErr)
	}

	_ = os.Remove(pidPath)
	writeServe.Stopped(cmd, pid)
	return nil
}
