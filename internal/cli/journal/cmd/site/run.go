//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core"
	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write"
)

// runZensical executes zensical build or serve in the output directory.
//
// Parameters:
//   - dir: Directory containing the generated site
//   - command: "build" or "serve"
//
// Returns:
//   - error: Non-nil if zensical is not found or fails
func runZensical(dir, command string) error {
	// Check if zensical is available
	_, lookErr := exec.LookPath(config.BinZensical)
	if lookErr != nil {
		return ctxerr.ZensicalNotFound()
	}

	// G204: binary is a constant, command is from the caller
	cmd := exec.Command(config.BinZensical, command) //nolint:gosec
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// runJournalSite handles the journal site command.
//
// Scans .context/journal/ for Markdown files, generates a zensical project
// structure, and optionally builds or serves the site.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - output: Output directory for the generated site
//   - build: If true, run zensical build after generating
//   - serve: If true, run zensical serve after generating
//
// Returns:
//   - error: Non-nil if generation fails
func runJournalSite(
	cmd *cobra.Command, output string, build, serve bool,
) error {
	journalDir := filepath.Join(rc.ContextDir(), config.DirJournal)

	// Check if the journal directory exists
	if _, statErr := os.Stat(journalDir); os.IsNotExist(statErr) {
		return ctxerr.NoJournalDir(journalDir)
	}

	// Load journal state for per-file processing flags
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return ctxerr.LoadJournalStateErr(loadErr)
	}

	// Scan journal files
	entries, scanErr := core.ScanJournalEntries(journalDir)
	if scanErr != nil {
		return ctxerr.ScanJournal(scanErr)
	}

	if len(entries) == 0 {
		return ctxerr.NoJournalEntries(journalDir)
	}

	// Create output directory structure
	docsDir := filepath.Join(output, config.JournalDirDocs)
	if mkErr := os.MkdirAll(docsDir, config.PermExec); mkErr != nil {
		return ctxerr.Mkdir(docsDir, mkErr)
	}

	// Write the stylesheet for <pre> overflow control
	stylesDir := filepath.Join(docsDir, config.DirStylesheets)
	if mkErr := os.MkdirAll(stylesDir, config.PermExec); mkErr != nil {
		return ctxerr.Mkdir(stylesDir, mkErr)
	}
	cssPath := filepath.Join(stylesDir, config.FileExtraCSS)
	cssData, cssReadErr := assets.JournalExtraCSS()
	if cssReadErr != nil {
		return cssReadErr
	}
	if writeErr := os.WriteFile(
		cssPath, cssData, config.PermFile,
	); writeErr != nil {
		return ctxerr.FileWrite(cssPath, writeErr)
	}

	// Write README
	readmePath := filepath.Join(output, config.FilenameReadme)
	if writeErr := os.WriteFile(
		readmePath,
		[]byte(core.GenerateSiteReadme(journalDir)), config.PermFile,
	); writeErr != nil {
		return ctxerr.FileWrite(readmePath, writeErr)
	}

	// Soft-wrap source journal files in-place, then copy to docs/
	for _, entry := range entries {
		src := entry.Path
		dst := filepath.Join(docsDir, entry.Filename)

		content, readErr := os.ReadFile(filepath.Clean(src))
		if readErr != nil {
			write.WarnFileErr(cmd, entry.Filename, readErr)
			continue
		}

		// Normalize the source file for readability
		normalized := core.CollapseToolOutputs(
			core.SoftWrapContent(
				core.MergeConsecutiveTurns(
					core.ConsolidateToolRuns(
						core.CleanToolOutputJSON(
							core.StripSystemReminders(string(content)),
						),
					),
				),
			),
		)
		if normalized != string(content) {
			if writeErr := os.WriteFile(
				src, []byte(normalized), config.PermFile,
			); writeErr != nil {
				write.WarnFileErr(cmd, entry.Filename, writeErr)
			}
		}

		// Generate site copy with Markdown fixes
		fv := jstate.IsFencesVerified(entry.Filename)
		withLinks := core.InjectSourceLink(normalized, src)
		if entry.Summary != "" {
			withLinks = core.InjectSummary(withLinks, entry.Summary)
		}
		siteContent := core.NormalizeContent(withLinks, fv)
		if writeErr := os.WriteFile(
			dst, []byte(siteContent), config.PermFile,
		); writeErr != nil {
			write.WarnFileErr(cmd, entry.Filename, writeErr)
			continue
		}
	}

	// Remove orphan site files — entries whose source was renamed or deleted.
	knownFiles := make(map[string]bool, len(entries)+1)
	knownFiles[config.FilenameIndex] = true
	for _, e := range entries {
		knownFiles[e.Filename] = true
	}
	if siteFiles, readErr := os.ReadDir(docsDir); readErr == nil {
		for _, f := range siteFiles {
			if f.IsDir() || knownFiles[f.Name()] {
				continue
			}
			orphanPath := filepath.Join(docsDir, f.Name())
			if rmErr := os.Remove(orphanPath); rmErr == nil {
				write.InfoJournalOrphanRemoved(cmd, f.Name())
			}
		}
	}

	// Generate index.md
	indexContent := core.GenerateIndex(entries)
	indexPath := filepath.Join(docsDir, config.FilenameIndex)
	if writeErr := os.WriteFile(
		indexPath, []byte(indexContent), config.PermFile,
	); writeErr != nil {
		return ctxerr.FileWrite(indexPath, writeErr)
	}

	// Generate topic pages
	var topicEntries []core.JournalEntry
	for _, e := range entries {
		if e.Suggestive || core.ContinuesMultipart(e.Filename) || len(e.Topics) == 0 {
			continue
		}
		topicEntries = append(topicEntries, e)
	}

	topics := core.BuildTopicIndex(topicEntries)

	if len(topics) > 0 {
		if writeErr := core.WriteSection(
			docsDir, config.JournalDirTopics,
			core.GenerateTopicsIndex(topics),
			func(dir string) {
				for _, t := range topics {
					if !t.Popular {
						continue
					}
					pagePath := filepath.Join(dir, t.Name+config.ExtMarkdown)
					if pageErr := os.WriteFile(
						pagePath, []byte(core.GenerateTopicPage(t)),
						config.PermFile,
					); pageErr != nil {
						write.WarnFileErr(cmd, pagePath, pageErr)
					}
				}
			}); writeErr != nil {
			return writeErr
		}
	}

	// Generate key files pages
	var keyFileEntries []core.JournalEntry
	for _, e := range entries {
		if e.Suggestive || core.ContinuesMultipart(e.Filename) || len(e.KeyFiles) == 0 {
			continue
		}
		keyFileEntries = append(keyFileEntries, e)
	}

	keyFiles := core.BuildKeyFileIndex(keyFileEntries)

	if len(keyFiles) > 0 {
		if writeErr := core.WriteSection(
			docsDir, config.JournalDirFiles,
			core.GenerateKeyFilesIndex(keyFiles),
			func(dir string) {
				for _, kf := range keyFiles {
					if !kf.Popular {
						continue
					}
					slug := core.KeyFileSlug(kf.Path)
					pagePath := filepath.Join(dir, slug+config.ExtMarkdown)
					if pageErr := os.WriteFile(
						pagePath, []byte(
							core.GenerateKeyFilePage(kf)),
						config.PermFile,
					); pageErr != nil {
						write.WarnFileErr(cmd, pagePath, pageErr)
					}
				}
			}); writeErr != nil {
			return writeErr
		}
	}

	// Generate session type pages
	var typeEntries []core.JournalEntry
	for _, e := range entries {
		if e.Suggestive || core.ContinuesMultipart(e.Filename) || e.Type == "" {
			continue
		}
		typeEntries = append(typeEntries, e)
	}

	sessionTypes := core.BuildTypeIndex(typeEntries)

	if len(sessionTypes) > 0 {
		if writeErr := core.WriteSection(
			docsDir,
			config.JournalDirTypes,
			core.GenerateTypesIndex(sessionTypes),
			func(dir string) {
				for _, st := range sessionTypes {
					pagePath := filepath.Join(dir, st.Name+config.ExtMarkdown)
					if pageErr := os.WriteFile(
						pagePath,
						[]byte(core.GenerateTypePage(st)), config.PermFile,
					); pageErr != nil {
						write.WarnFileErr(cmd, pagePath, pageErr)
					}
				}
			}); writeErr != nil {
			return writeErr
		}
	}

	// Generate zensical.toml
	tomlContent := core.GenerateZensicalToml(
		entries, topics, keyFiles, sessionTypes,
	)
	tomlPath := filepath.Join(output, config.FileZensicalToml)
	if writeErr := os.WriteFile(
		tomlPath,
		[]byte(tomlContent), config.PermFile,
	); writeErr != nil {
		return ctxerr.FileWrite(tomlPath, writeErr)
	}

	if serve {
		write.InfoJournalSiteStarting(cmd)
		return runZensical(output, "serve")
	} else if build {
		write.InfoJournalSiteBuilding(cmd)
		return runZensical(output, "build")
	}

	write.InfoJournalSiteGenerated(cmd, len(entries), output, config.BinZensical)

	return nil
}
