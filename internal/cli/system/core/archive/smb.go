//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/backup"
	fserr "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/io"
)

// SMBConfig holds parsed SMB share connection details.
type SMBConfig struct {
	Host      string
	Share     string
	Subdir    string
	GVFSPath  string
	SourceURL string
}

// ParseSMBConfig parses an SMB URL and subdirectory into a config struct
// with the derived GVFS mount path.
//
// Parameters:
//   - smbURL: SMB share URL (e.g. smb://host/share)
//   - subdir: Subdirectory on share (empty uses default)
//
// Returns:
//   - *SMBConfig: Parsed config
//   - error: Non-nil on invalid URL
func ParseSMBConfig(smbURL, subdir string) (*SMBConfig, error) {
	u, parseErr := url.Parse(smbURL)
	if parseErr != nil || u.Host == "" {
		return nil, ctxerr.InvalidSMBURL(smbURL)
	}

	host := u.Host
	share := u.Path
	if len(share) > 0 && share[0] == '/' {
		share = share[1:]
	}
	if share == "" {
		return nil, ctxerr.SMBMissingShare(smbURL)
	}

	if subdir == "" {
		subdir = archive.BackupDefaultSubdir
	}

	gvfsPath := fmt.Sprintf(desc.Text(text.DescKeyWriteFormatGVFSPath),
		os.Getuid(), host, share)

	return &SMBConfig{
		Host:      host,
		Share:     share,
		Subdir:    subdir,
		GVFSPath:  gvfsPath,
		SourceURL: smbURL,
	}, nil
}

// EnsureSMBMount checks if the GVFS mount exists and attempts gio mount if not.
//
// Parameters:
//   - cfg: SMB configuration
//
// Returns:
//   - error: Non-nil if mount fails
func EnsureSMBMount(cfg *SMBConfig) error {
	if _, statErr := os.Stat(cfg.GVFSPath); statErr == nil {
		return nil
	}

	//nolint:gosec // G204: smbURL is from user env, not untrusted input
	mountCmd := exec.Command("gio", "mount", cfg.SourceURL)
	if mountErr := mountCmd.Run(); mountErr != nil {
		return ctxerr.MountFailed(cfg.SourceURL, mountErr)
	}

	return nil
}

// CopyToSMB copies a local file to the SMB share destination directory.
//
// Parameters:
//   - cfg: SMB configuration
//   - localPath: Path to the local file to copy
//
// Returns:
//   - error: Non-nil on copy failure
func CopyToSMB(cfg *SMBConfig, localPath string) error {
	dest := filepath.Join(cfg.GVFSPath, cfg.Subdir)
	if mkdirErr := os.MkdirAll(dest, fs.PermExec); mkdirErr != nil {
		return fserr.CreateDir(dest, mkdirErr)
	}

	data, readErr := io.SafeReadUserFile(localPath)
	if readErr != nil {
		return fserr.ReadFile(readErr)
	}

	destFile := filepath.Join(dest, filepath.Base(localPath))
	if writeErr := os.WriteFile(destFile, data, fs.PermFile); writeErr != nil {
		return ctxerr.WriteSMB(writeErr)
	}

	return nil
}
