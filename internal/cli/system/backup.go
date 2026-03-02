//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
)

// backupScope enumerates the valid --scope values.
const (
	scopeProject = "project"
	scopeGlobal  = "global"
	scopeAll     = "all"
)

// archiveEntry describes a directory or file to include in a backup archive.
type archiveEntry struct {
	// SourcePath is the absolute path to the directory or file.
	SourcePath string
	// Prefix is the path prefix inside the tar archive.
	Prefix string
	// ExcludeDir is a directory name to skip (e.g. "journal-site").
	ExcludeDir string
	// Optional means a missing source is not an error.
	Optional bool
}

// backupResult holds the outcome of a single archive creation.
type backupResult struct {
	Scope   string `json:"scope"`
	Archive string `json:"archive"`
	Size    int64  `json:"size"`
	SMBDest string `json:"smb_dest,omitempty"`
}

func backupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup context and Claude data",
		Long: `Create timestamped tar.gz archives of project context and/or global
Claude Code data. Optionally copies archives to an SMB share.

Scopes:
  project  .context/, .claude/, ideas/, ~/.bashrc
  global   ~/.claude/ (excludes todos/)
  all      Both project and global (default)

Environment:
  CTX_BACKUP_SMB_URL    - SMB share URL (e.g. smb://host/share)
  CTX_BACKUP_SMB_SUBDIR - Subdirectory on share (default: ctx-sessions)`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runBackup(cmd)
		},
	}
	cmd.Flags().String("scope", scopeAll, "Backup scope: project, global, or all")
	cmd.Flags().Bool("json", false, "Output results as JSON")
	return cmd
}

func runBackup(cmd *cobra.Command) error {
	scope, _ := cmd.Flags().GetString("scope")
	jsonOut, _ := cmd.Flags().GetBool("json")

	switch scope {
	case scopeProject, scopeGlobal, scopeAll:
	default:
		return fmt.Errorf("invalid scope %q: must be project, global, or all", scope)
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return fmt.Errorf("determine home directory: %w", homeErr)
	}

	// Parse SMB config if configured.
	smbURL := os.Getenv(config.EnvBackupSMBURL)
	smbSubdir := os.Getenv(config.EnvBackupSMBSubdir)
	var smb *smbConfig
	if smbURL != "" {
		var smbErr error
		smb, smbErr = parseSMBConfig(smbURL, smbSubdir)
		if smbErr != nil {
			return fmt.Errorf("parse SMB config: %w", smbErr)
		}
	}

	timestamp := time.Now().Format("20060102-150405")
	var results []backupResult

	if scope == scopeProject || scope == scopeAll {
		result, projErr := backupProject(cmd, home, timestamp, smb)
		if projErr != nil {
			return fmt.Errorf("project backup: %w", projErr)
		}
		results = append(results, result)
	}

	if scope == scopeGlobal || scope == scopeAll {
		result, globalErr := backupGlobal(cmd, home, timestamp, smb)
		if globalErr != nil {
			return fmt.Errorf("global backup: %w", globalErr)
		}
		results = append(results, result)
	}

	if jsonOut {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	for _, r := range results {
		cmd.Printf("%s: %s (%s)",
			r.Scope, r.Archive, formatSize(r.Size))
		if r.SMBDest != "" {
			cmd.Printf(" → %s", r.SMBDest)
		}
		cmd.Println()
	}
	return nil
}

func backupProject(
	cmd *cobra.Command, home, timestamp string, smb *smbConfig,
) (backupResult, error) {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return backupResult{}, cwdErr
	}

	archiveName := fmt.Sprintf("ctx-backup-%s.tar.gz", timestamp)
	archivePath := filepath.Join(os.TempDir(), archiveName)

	entries := []archiveEntry{
		{SourcePath: filepath.Join(cwd, ".context"), Prefix: ".context", ExcludeDir: "journal-site"},
		{SourcePath: filepath.Join(cwd, ".claude"), Prefix: ".claude"},
		{SourcePath: filepath.Join(cwd, "ideas"), Prefix: "ideas", Optional: true},
		{SourcePath: filepath.Join(home, ".bashrc"), Prefix: ".bashrc"},
	}

	archiveErr := createArchive(archivePath, entries, cmd)
	if archiveErr != nil {
		return backupResult{}, archiveErr
	}

	result := backupResult{Scope: scopeProject, Archive: archivePath}
	info, statErr := os.Stat(archivePath)
	if statErr == nil {
		result.Size = info.Size()
	}

	if smb != nil {
		if mountErr := ensureSMBMount(smb); mountErr != nil {
			return result, mountErr
		}
		if copyErr := copyToSMB(smb, archivePath); copyErr != nil {
			return result, copyErr
		}
		result.SMBDest = filepath.Join(smb.GVFSPath, smb.Subdir, archiveName)
	}

	// Touch marker file for check-backup-age hook.
	markerDir := filepath.Join(home, ".local", "state")
	_ = os.MkdirAll(markerDir, config.PermExec)
	markerPath := filepath.Join(markerDir, config.BackupMarkerFile)
	touchFile(markerPath)

	return result, nil
}

func backupGlobal(
	cmd *cobra.Command, home, timestamp string, smb *smbConfig,
) (backupResult, error) {
	archiveName := fmt.Sprintf("claude-global-backup-%s.tar.gz", timestamp)
	archivePath := filepath.Join(os.TempDir(), archiveName)

	entries := []archiveEntry{
		{SourcePath: filepath.Join(home, ".claude"), Prefix: ".claude", ExcludeDir: "todos"},
	}

	archiveErr := createArchive(archivePath, entries, cmd)
	if archiveErr != nil {
		return backupResult{}, archiveErr
	}

	result := backupResult{Scope: scopeGlobal, Archive: archivePath}
	info, statErr := os.Stat(archivePath)
	if statErr == nil {
		result.Size = info.Size()
	}

	if smb != nil {
		if mountErr := ensureSMBMount(smb); mountErr != nil {
			return result, mountErr
		}
		if copyErr := copyToSMB(smb, archivePath); copyErr != nil {
			return result, copyErr
		}
		result.SMBDest = filepath.Join(smb.GVFSPath, smb.Subdir, archiveName)
	}

	return result, nil
}

// createArchive builds a tar.gz archive from the given entries.
func createArchive(
	archivePath string, entries []archiveEntry, cmd *cobra.Command,
) error {
	outFile, createErr := os.Create(archivePath) //nolint:gosec // tmp path
	if createErr != nil {
		return fmt.Errorf("create archive file: %w", createErr)
	}
	defer func() { _ = outFile.Close() }()

	gzw := gzip.NewWriter(outFile)
	defer func() { _ = gzw.Close() }()

	tw := tar.NewWriter(gzw)
	defer func() { _ = tw.Close() }()

	for _, entry := range entries {
		addErr := addEntry(tw, entry, cmd)
		if addErr != nil {
			return addErr
		}
	}
	return nil
}

// addEntry adds a single archiveEntry (file or directory) to the tar writer.
func addEntry(tw *tar.Writer, entry archiveEntry, cmd *cobra.Command) error {
	info, statErr := os.Stat(entry.SourcePath)
	if os.IsNotExist(statErr) {
		if entry.Optional {
			cmd.PrintErrf("skipping %s (not found)\n", entry.Prefix)
			return nil
		}
		return fmt.Errorf("source not found: %s", entry.SourcePath)
	}
	if statErr != nil {
		return statErr
	}

	// Single file (e.g. ~/.bashrc).
	if !info.IsDir() {
		return addSingleFile(tw, entry.SourcePath, entry.Prefix, info)
	}

	// Directory walk.
	return filepath.WalkDir(entry.SourcePath,
		func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			// Skip excluded directories.
			if d.IsDir() && entry.ExcludeDir != "" && d.Name() == entry.ExcludeDir {
				return filepath.SkipDir
			}

			// Skip symlinks.
			if d.Type()&os.ModeSymlink != 0 {
				return nil
			}

			rel, relErr := filepath.Rel(entry.SourcePath, path)
			if relErr != nil {
				return relErr
			}

			name := filepath.Join(entry.Prefix, rel)
			// Normalize to forward slashes in tar.
			name = filepath.ToSlash(name)

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
	f, openErr := os.Open(path) //nolint:gosec // paths are from our own entries
	if openErr != nil {
		return openErr
	}
	defer func() { _ = f.Close() }()
	_, copyErr := io.Copy(tw, f)
	return copyErr
}

// formatSize returns a human-readable file size string.
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMG"[exp])
}
