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
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errBackup "github.com/ActiveMemory/ctx/internal/err/backup"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	execGio "github.com/ActiveMemory/ctx/internal/exec/gio"
	"github.com/ActiveMemory/ctx/internal/io"
)

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
		return nil, errBackup.InvalidSMBURL(smbURL)
	}

	host := u.Host
	share := u.Path
	if len(share) > 0 && share[0] == '/' {
		share = share[1:]
	}
	if share == "" {
		return nil, errBackup.SMBMissingShare(smbURL)
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

	if mountErr := execGio.Mount(cfg.SourceURL); mountErr != nil {
		return errBackup.MountFailed(cfg.SourceURL, mountErr)
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
		return errFs.CreateDir(dest, mkdirErr)
	}

	data, readErr := io.SafeReadUserFile(localPath)
	if readErr != nil {
		return errFs.ReadFile(readErr)
	}

	destFile := filepath.Join(dest, filepath.Base(localPath))
	if writeErr := os.WriteFile(destFile, data, fs.PermFile); writeErr != nil {
		return errBackup.WriteSMB(writeErr)
	}

	return nil
}
