//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/entity"
	errBackup "github.com/ActiveMemory/ctx/internal/err/backup"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// Create builds a tar.gz archive from the given entries.
//
// Parameters:
//   - archivePath: output file path for the archive
//   - entries: directories and files to include
//   - w: writer for diagnostic output (typically stderr)
//
// Returns:
//   - error: non-nil on file creation or tar writing failure
func Create(
	archivePath string, entries []entity.ArchiveEntry, w io.Writer,
) error {
	outFile, createErr := internalIo.SafeCreateFile(archivePath, cfgFs.PermFile)
	if createErr != nil {
		return errBackup.CreateArchive(createErr)
	}
	defer func() { _ = outFile.Close() }()

	gzw := gzip.NewWriter(outFile)
	defer func() { _ = gzw.Close() }()

	tw := tar.NewWriter(gzw)
	defer func() { _ = tw.Close() }()

	for _, entry := range entries {
		if addErr := addEntry(tw, entry, w); addErr != nil {
			return addErr
		}
	}
	return nil
}

// BackupProject creates a project-scoped backup archive.
//
// Parameters:
//   - w: writer for diagnostic output (typically stderr)
//   - home: user home directory
//   - timestamp: formatted timestamp for the archive filename
//   - smb: optional SMB configuration (nil to skip remote copy)
//
// Returns:
//   - BackupResult: archive path, size, and optional SMB destination
//   - error: non-nil on archive or SMB failure
func BackupProject(
	w io.Writer, home, timestamp string, smb *SMBConfig,
) (entity.BackupResult, error) {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return entity.BackupResult{}, cwdErr
	}

	archiveName := fmt.Sprintf(archive.TplProjectArchive, timestamp)
	archivePath := filepath.Join(os.TempDir(), archiveName)

	entries := []entity.ArchiveEntry{
		{SourcePath: filepath.Join(cwd, dir.Context), Prefix: dir.Context, ExcludeDir: dir.JournalSite},
		{SourcePath: filepath.Join(cwd, dir.Claude), Prefix: dir.Claude},
		{SourcePath: filepath.Join(cwd, dir.Ideas), Prefix: dir.Ideas, Optional: true},
		{SourcePath: filepath.Join(home, archive.Bashrc), Prefix: archive.Bashrc},
	}

	result, finalizeErr := finalizeArchive(
		w, archivePath, archiveName, archive.BackupScopeProject, entries, smb,
	)
	if finalizeErr != nil {
		return result, finalizeErr
	}

	// Touch marker file for check-backup-age hook.
	markerDir := filepath.Join(home, archive.BackupMarkerDir)
	_ = os.MkdirAll(markerDir, cfgFs.PermExec)
	markerPath := filepath.Join(markerDir, archive.BackupMarkerFile)
	internalIo.TouchFile(markerPath)

	return result, nil
}

// BackupGlobal creates a global-scoped backup archive.
//
// Parameters:
//   - w: writer for diagnostic output (typically stderr)
//   - home: user home directory
//   - timestamp: formatted timestamp for the archive filename
//   - smb: optional SMB configuration (nil to skip remote copy)
//
// Returns:
//   - BackupResult: archive path, size, and optional SMB destination
//   - error: non-nil on archive or SMB failure
func BackupGlobal(
	w io.Writer, home, timestamp string, smb *SMBConfig,
) (entity.BackupResult, error) {
	archiveName := fmt.Sprintf(archive.TplGlobalArchive, timestamp)
	archivePath := filepath.Join(os.TempDir(), archiveName)

	entries := []entity.ArchiveEntry{
		{SourcePath: filepath.Join(home, dir.Claude), Prefix: dir.Claude, ExcludeDir: archive.BackupExcludeTodos},
	}

	return finalizeArchive(
		w, archivePath, archiveName, archive.BackupScopeGlobal, entries, smb,
	)
}

// CheckSMBMountWarnings checks whether the GVFS mount for the given SMB URL
// exists and appends warning strings if the share is not mounted.
//
// Parameters:
//   - smbURL: the SMB share URL from the environment
//   - warnings: existing warning slice to append to
//
// Returns:
//   - []string: the warnings slice, possibly with SMB mount warnings appended
func CheckSMBMountWarnings(smbURL string, warnings []string) []string {
	cfg, cfgErr := ParseSMBConfig(smbURL, "")
	if cfgErr != nil {
		return warnings
	}

	if _, statErr := os.Stat(cfg.GVFSPath); os.IsNotExist(statErr) {
		warnings = append(warnings,
			fmt.Sprintf(desc.Text(text.DescKeyBackupSMBNotMounted), cfg.Host),
			desc.Text(text.DescKeyBackupSMBUnavailable),
		)
	}

	return warnings
}

// CheckBackupMarker checks the backup marker file age and appends warnings
// when the marker is missing or older than config.BackupMaxAgeDays.
//
// Parameters:
//   - markerPath: absolute path to the backup marker file
//   - warnings: existing warning slice to append to
//
// Returns:
//   - []string: the warnings slice, possibly with staleness warnings appended
func CheckBackupMarker(markerPath string, warnings []string) []string {
	info, statErr := os.Stat(markerPath)
	if os.IsNotExist(statErr) {
		return append(warnings,
			desc.Text(text.DescKeyBackupNoMarker),
			desc.Text(text.DescKeyBackupRunHint),
		)
	}
	if statErr != nil {
		return warnings
	}

	ageDays := int(time.Since(info.ModTime()).Hours() / 24)
	if ageDays >= archive.BackupMaxAgeDays {
		return append(warnings,
			fmt.Sprintf(desc.Text(text.DescKeyBackupStale), ageDays),
			desc.Text(text.DescKeyBackupRunHint),
		)
	}

	return warnings
}
