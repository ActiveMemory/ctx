//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	fs2 "github.com/ActiveMemory/ctx/internal/config/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/backup"
	io2 "github.com/ActiveMemory/ctx/internal/io"
	"github.com/spf13/cobra"
)

// BackupProject creates a project-scoped backup archive.
//
// Parameters:
//   - cmd: Cobra command for diagnostic output
//   - home: user home directory
//   - timestamp: formatted timestamp for the archive filename
//   - smb: optional SMB configuration (nil to skip remote copy)
//
// Returns:
//   - BackupResult: archive path, size, and optional SMB destination
//   - error: non-nil on archive or SMB failure
func BackupProject(
	cmd *cobra.Command, home, timestamp string, smb *SMBConfig,
) (BackupResult, error) {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return BackupResult{}, cwdErr
	}

	archiveName := fmt.Sprintf(archive.BackupTplProjectArchive, timestamp)
	archivePath := filepath.Join(os.TempDir(), archiveName)

	entries := []ArchiveEntry{
		{SourcePath: filepath.Join(cwd, dir.Context), Prefix: dir.Context, ExcludeDir: dir.JournalSite},
		{SourcePath: filepath.Join(cwd, dir.Claude), Prefix: dir.Claude},
		{SourcePath: filepath.Join(cwd, dir.Ideas), Prefix: dir.Ideas, Optional: true},
		{SourcePath: filepath.Join(home, archive.Bashrc), Prefix: archive.Bashrc},
	}

	result, finalizeErr := finalizeArchive(
		cmd, archivePath, archiveName, archive.BackupScopeProject, entries, smb,
	)
	if finalizeErr != nil {
		return result, finalizeErr
	}

	// Touch marker file for check-backup-age hook.
	markerDir := filepath.Join(home, archive.BackupMarkerDir)
	_ = os.MkdirAll(markerDir, fs2.PermExec)
	markerPath := filepath.Join(markerDir, archive.BackupMarkerFile)
	TouchFile(markerPath)

	return result, nil
}

// BackupGlobal creates a global-scoped backup archive.
//
// Parameters:
//   - cmd: Cobra command for diagnostic output
//   - home: user home directory
//   - timestamp: formatted timestamp for the archive filename
//   - smb: optional SMB configuration (nil to skip remote copy)
//
// Returns:
//   - BackupResult: archive path, size, and optional SMB destination
//   - error: non-nil on archive or SMB failure
func BackupGlobal(
	cmd *cobra.Command, home, timestamp string, smb *SMBConfig,
) (BackupResult, error) {
	archiveName := fmt.Sprintf(archive.BackupTplGlobalArchive, timestamp)
	archivePath := filepath.Join(os.TempDir(), archiveName)

	entries := []ArchiveEntry{
		{SourcePath: filepath.Join(home, dir.Claude), Prefix: dir.Claude, ExcludeDir: archive.BackupExcludeTodos},
	}

	return finalizeArchive(
		cmd, archivePath, archiveName, archive.BackupScopeGlobal, entries, smb,
	)
}

// finalizeArchive creates the archive, populates the result with size,
// and optionally copies to an SMB share.
func finalizeArchive(
	cmd *cobra.Command, archivePath, archiveName, scope string,
	entries []ArchiveEntry, smb *SMBConfig,
) (BackupResult, error) {
	if archiveErr := CreateArchive(archivePath, entries, cmd); archiveErr != nil {
		return BackupResult{}, archiveErr
	}

	result := BackupResult{Scope: scope, Archive: archivePath}
	if info, statErr := os.Stat(archivePath); statErr == nil {
		result.Size = info.Size()
	}

	if smb != nil {
		if mountErr := EnsureSMBMount(smb); mountErr != nil {
			return result, mountErr
		}
		if copyErr := CopyToSMB(smb, archivePath); copyErr != nil {
			return result, copyErr
		}
		result.SMBDest = filepath.Join(smb.GVFSPath, smb.Subdir, archiveName)
	}

	return result, nil
}

// CreateArchive builds a tar.gz archive from the given entries.
//
// Parameters:
//   - archivePath: output file path for the archive
//   - entries: directories and files to include
//   - cmd: Cobra command for diagnostic output
//
// Returns:
//   - error: non-nil on file creation or tar writing failure
func CreateArchive(
	archivePath string, entries []ArchiveEntry, cmd *cobra.Command,
) error {
	outFile, createErr := os.Create(archivePath) //nolint:gosec // tmp path from our own code
	if createErr != nil {
		return ctxerr.CreateArchive(createErr)
	}
	defer func() { _ = outFile.Close() }()

	gzw := gzip.NewWriter(outFile)
	defer func() { _ = gzw.Close() }()

	tw := tar.NewWriter(gzw)
	defer func() { _ = tw.Close() }()

	for _, entry := range entries {
		if addErr := addEntry(tw, entry, cmd); addErr != nil {
			return addErr
		}
	}
	return nil
}

// addEntry adds a single ArchiveEntry (file or directory) to the tar writer.
func addEntry(tw *tar.Writer, entry ArchiveEntry, cmd *cobra.Command) error {
	info, statErr := os.Stat(entry.SourcePath)
	if os.IsNotExist(statErr) {
		if entry.Optional {
			cmd.PrintErrln(fmt.Sprintf("skipping %s (not found)", entry.Prefix))
			return nil
		}
		return ctxerr.SourceNotFound(entry.SourcePath)
	}
	if statErr != nil {
		return statErr
	}

	if !info.IsDir() {
		return addSingleFile(tw, entry.SourcePath, entry.Prefix, info)
	}

	return filepath.WalkDir(entry.SourcePath,
		func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() && entry.ExcludeDir != "" && d.Name() == entry.ExcludeDir {
				return filepath.SkipDir
			}
			if d.Type()&os.ModeSymlink != 0 {
				return nil
			}

			rel, relErr := filepath.Rel(entry.SourcePath, path)
			if relErr != nil {
				return relErr
			}

			name := filepath.ToSlash(filepath.Join(entry.Prefix, rel))

			fileInfo, infoErr := d.Info()
			if infoErr != nil {
				return infoErr
			}

			header, headerErr := tar.FileInfoHeader(fileInfo, "")
			if headerErr != nil {
				return headerErr
			}
			header.Name = name

			if writeErr := tw.WriteHeader(header); writeErr != nil {
				return writeErr
			}

			if d.IsDir() {
				return nil
			}
			return copyFileToTar(tw, path)
		})
}

// addSingleFile writes a single file entry into the tar.
func addSingleFile(
	tw *tar.Writer, path, name string, info fs.FileInfo,
) error {
	header, headerErr := tar.FileInfoHeader(info, "")
	if headerErr != nil {
		return headerErr
	}
	header.Name = name

	if writeErr := tw.WriteHeader(header); writeErr != nil {
		return writeErr
	}
	return copyFileToTar(tw, path)
}

// copyFileToTar reads a file and writes its contents to the tar writer.
func copyFileToTar(tw *tar.Writer, path string) error {
	f, openErr := io2.SafeOpenUserFile(path)
	if openErr != nil {
		return openErr
	}
	defer func() { _ = f.Close() }()
	_, copyErr := io.Copy(tw, f)
	return copyErr
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
			fmt.Sprintf(desc.TextDesc(text.DescKeyBackupSMBNotMounted), cfg.Host),
			desc.TextDesc(text.DescKeyBackupSMBUnavailable),
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
			desc.TextDesc(text.DescKeyBackupNoMarker),
			desc.TextDesc(text.DescKeyBackupRunHint),
		)
	}
	if statErr != nil {
		return warnings
	}

	ageDays := int(time.Since(info.ModTime()).Hours() / 24)
	if ageDays >= archive.BackupMaxAgeDays {
		return append(warnings,
			fmt.Sprintf(desc.TextDesc(text.DescKeyBackupStale), ageDays),
			desc.TextDesc(text.DescKeyBackupRunHint),
		)
	}

	return warnings
}
