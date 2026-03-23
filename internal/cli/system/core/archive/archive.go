//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
	errBackup "github.com/ActiveMemory/ctx/internal/err/backup"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// finalizeArchive creates the archive, populates the result with size,
// and optionally copies to an SMB share.
//
// Parameters:
//   - w: writer for diagnostic output (typically stderr)
//   - archivePath: output file path for the archive
//   - archiveName: archive filename (for SMB destination path)
//   - scope: backup scope label (e.g., "project", "global")
//   - entries: directories and files to include
//   - smb: optional SMB configuration (nil to skip remote copy)
//
// Returns:
//   - BackupResult: archive path, size, and optional SMB destination
//   - error: non-nil on archive creation or SMB failure
func finalizeArchive(
	w io.Writer, archivePath, archiveName, scope string,
	entries []entity.ArchiveEntry, smb *SMBConfig,
) (entity.BackupResult, error) {
	if archiveErr := CreateArchive(archivePath, entries, w); archiveErr != nil {
		return entity.BackupResult{}, archiveErr
	}

	result := entity.BackupResult{Scope: scope, Archive: archivePath}
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

// addEntry adds a single ArchiveEntry (file or directory) to the tar writer.
// Optional entries that are not found emit a diagnostic message and are skipped.
//
// Parameters:
//   - tw: tar writer to add the entry to
//   - entry: archive entry describing the source and target
//   - w: writer for diagnostic output (typically stderr)
//
// Returns:
//   - error: non-nil on stat, walk, or tar write failure
func addEntry(tw *tar.Writer, entry entity.ArchiveEntry, w io.Writer) error {
	info, statErr := os.Stat(entry.SourcePath)
	if os.IsNotExist(statErr) {
		if entry.Optional {
			_, _ = fmt.Fprintf(
				w, desc.Text(text.DescKeyWriteBackupSkipEntry), entry.Prefix,
			)
			return nil
		}
		return errBackup.SourceNotFound(entry.SourcePath)
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

// addSingleFile writes a single non-directory file entry into the tar.
//
// Parameters:
//   - tw: tar writer
//   - path: absolute source file path
//   - name: name to use inside the archive
//   - info: file info for the tar header
//
// Returns:
//   - error: non-nil on header or content write failure
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
//
// Parameters:
//   - tw: tar writer
//   - path: absolute file path to read
//
// Returns:
//   - error: non-nil on open, read, or write failure
func copyFileToTar(tw *tar.Writer, path string) error {
	f, openErr := internalIo.SafeOpenUserFile(path)
	if openErr != nil {
		return openErr
	}
	defer func() { _ = f.Close() }()
	_, copyErr := io.Copy(tw, f)
	return copyErr
}
