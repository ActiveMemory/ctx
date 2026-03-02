//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config"
)

// smbConfig holds parsed SMB share connection details.
type smbConfig struct {
	Host      string
	Share     string
	Subdir    string
	GVFSPath  string
	SourceURL string
}

// parseSMBConfig parses an SMB URL and subdirectory into a config struct
// with the derived GVFS mount path.
func parseSMBConfig(smbURL, subdir string) (*smbConfig, error) {
	u, parseErr := url.Parse(smbURL)
	if parseErr != nil || u.Host == "" {
		return nil, fmt.Errorf("invalid SMB URL: %s", smbURL)
	}

	host := u.Host
	share := u.Path
	if len(share) > 0 && share[0] == '/' {
		share = share[1:]
	}
	if share == "" {
		return nil, fmt.Errorf("SMB URL missing share name: %s", smbURL)
	}

	if subdir == "" {
		subdir = config.BackupDefaultSubdir
	}

	gvfsPath := fmt.Sprintf("/run/user/%d/gvfs/smb-share:server=%s,share=%s",
		os.Getuid(), host, share)

	return &smbConfig{
		Host:      host,
		Share:     share,
		Subdir:    subdir,
		GVFSPath:  gvfsPath,
		SourceURL: smbURL,
	}, nil
}

// ensureSMBMount checks if the GVFS mount exists and attempts gio mount if not.
func ensureSMBMount(cfg *smbConfig) error {
	if _, statErr := os.Stat(cfg.GVFSPath); statErr == nil {
		return nil
	}

	//nolint:gosec // G204: smbURL is from user env, not untrusted input
	mountCmd := exec.Command("gio", "mount", cfg.SourceURL)
	if mountErr := mountCmd.Run(); mountErr != nil {
		return fmt.Errorf("failed to mount %s: %w", cfg.SourceURL, mountErr)
	}

	return nil
}

// copyToSMB copies a local file to the SMB share destination directory.
func copyToSMB(cfg *smbConfig, localPath string) error {
	dest := filepath.Join(cfg.GVFSPath, cfg.Subdir)
	if mkdirErr := os.MkdirAll(dest, config.PermExec); mkdirErr != nil {
		return fmt.Errorf("create destination dir: %w", mkdirErr)
	}

	data, readErr := os.ReadFile(localPath) //nolint:gosec // path from our own archive
	if readErr != nil {
		return fmt.Errorf("read archive: %w", readErr)
	}

	destFile := filepath.Join(dest, filepath.Base(localPath))
	if writeErr := os.WriteFile(destFile, data, config.PermFile); writeErr != nil {
		return fmt.Errorf("write to SMB: %w", writeErr)
	}

	return nil
}
