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

	readJournal "github.com/ActiveMemory/ctx/internal/assets/read/journal"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/collapse"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/consolidate"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/format"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/generate"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/normalize"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/parse"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/reduce"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/section"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/turn"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/wrap"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/zensical"
	"github.com/ActiveMemory/ctx/internal/entity"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/err/journal"
	errSite "github.com/ActiveMemory/ctx/internal/err/site"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/err"
	writeJournal "github.com/ActiveMemory/ctx/internal/write/journal"
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
	_, lookErr := exec.LookPath(zensical.Bin)
	if lookErr != nil {
		return errSite.ZensicalNotFound()
	}

	// G204: binary is a constant, command is from the caller
	cmd := exec.Command(zensical.Bin, command) //nolint:gosec
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// Run handles the journal site command.
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
func Run(
	cmd *cobra.Command, output string, build, serve bool,
) error {
	journalDir := filepath.Join(rc.ContextDir(), dir.Journal)

	// Check if the journal directory exists
	if _, statErr := os.Stat(journalDir); os.IsNotExist(statErr) {
		return journal.NoDir(journalDir)
	}

	// Load journal state for per-file processing flags
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return journal.LoadStateErr(loadErr)
	}

	// Scan journal files
	entries, scanErr := parse.ScanJournalEntries(journalDir)
	if scanErr != nil {
		return journal.Scan(scanErr)
	}

	if len(entries) == 0 {
		return journal.NoEntries(journalDir)
	}

	// Create output directory structure
	docsDir := filepath.Join(output, dir.JournalDocs)
	if mkErr := os.MkdirAll(docsDir, fs.PermExec); mkErr != nil {
		return errFs.Mkdir(docsDir, mkErr)
	}

	// Write the stylesheet for <pre> overflow control
	stylesDir := filepath.Join(docsDir, zensical.Stylesheets)
	if mkErr := os.MkdirAll(stylesDir, fs.PermExec); mkErr != nil {
		return errFs.Mkdir(stylesDir, mkErr)
	}
	cssPath := filepath.Join(stylesDir, zensical.ExtraCSS)
	cssData, cssReadErr := readJournal.ExtraCSS()
	if cssReadErr != nil {
		return cssReadErr
	}
	if writeErr := os.WriteFile(
		cssPath, cssData, fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(cssPath, writeErr)
	}

	// Write README
	readmePath := filepath.Join(output, file.Readme)
	if writeErr := os.WriteFile(
		readmePath,
		[]byte(generate.GenerateSiteReadme(journalDir)), fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(readmePath, writeErr)
	}

	// Soft-wrap source journal files in-place, then copy to docs/
	for _, entry := range entries {
		src := entry.Path
		dst := filepath.Join(docsDir, entry.Filename)

		content, readErr := os.ReadFile(filepath.Clean(src))
		if readErr != nil {
			err.WarnFile(cmd, entry.Filename, readErr)
			continue
		}

		// Normalize the source file for readability
		normalized := collapse.CollapseToolOutputs(
			wrap.SoftWrapContent(
				turn.MergeConsecutiveTurns(
					consolidate.ConsolidateToolRuns(
						reduce.CleanToolOutputJSON(
							reduce.StripSystemReminders(string(content)),
						),
					),
				),
			),
		)
		if normalized != string(content) {
			if writeErr := os.WriteFile(
				src, []byte(normalized), fs.PermFile,
			); writeErr != nil {
				err.WarnFile(cmd, entry.Filename, writeErr)
			}
		}

		// Generate site copy with Markdown fixes
		fv := jstate.FencesVerified(entry.Filename)
		withLinks := generate.InjectedSourceLink(normalized, src)
		if entry.Summary != "" {
			withLinks = generate.InjectedSummary(withLinks, entry.Summary)
		}
		siteContent := normalize.NormalizeContent(withLinks, fv)
		if writeErr := os.WriteFile(
			dst, []byte(siteContent), fs.PermFile,
		); writeErr != nil {
			err.WarnFile(cmd, entry.Filename, writeErr)
			continue
		}
	}

	// Remove orphan site files: entries whose source was renamed or deleted.
	knownFiles := make(map[string]bool, len(entries)+1)
	knownFiles[file.Index] = true
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
				writeJournal.InfoOrphanRemoved(cmd, f.Name())
			}
		}
	}

	// Generate index.md
	indexContent := generate.GenerateIndex(entries)
	indexPath := filepath.Join(docsDir, file.Index)
	if writeErr := os.WriteFile(
		indexPath, []byte(indexContent), fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(indexPath, writeErr)
	}

	// Generate topic pages
	var topicEntries []entity.JournalEntry
	for _, e := range entries {
		if e.Suggestive || section.ContinuesMultipart(e.Filename) || len(e.Topics) == 0 {
			continue
		}
		topicEntries = append(topicEntries, e)
	}

	topics := section.BuildTopicIndex(topicEntries)

	if len(topics) > 0 {
		if writeErr := section.WriteSection(
			docsDir, dir.JournTopics,
			section.GenerateTopicsIndex(topics),
			func(dir string) {
				for _, t := range topics {
					if !t.Popular {
						continue
					}
					pagePath := filepath.Join(dir, t.Name+file.ExtMarkdown)
					if pageErr := os.WriteFile(
						pagePath, []byte(section.GenerateTopicPage(t)),
						fs.PermFile,
					); pageErr != nil {
						err.WarnFile(cmd, pagePath, pageErr)
					}
				}
			}); writeErr != nil {
			return writeErr
		}
	}

	// Generate key files pages
	var keyFileEntries []entity.JournalEntry
	for _, e := range entries {
		if e.Suggestive || section.ContinuesMultipart(
			e.Filename,
		) || len(e.KeyFiles) == 0 {
			continue
		}
		keyFileEntries = append(keyFileEntries, e)
	}

	keyFiles := section.BuildKeyFileIndex(keyFileEntries)

	if len(keyFiles) > 0 {
		if writeErr := section.WriteSection(
			docsDir, dir.JournalFiles,
			section.GenerateKeyFilesIndex(keyFiles),
			func(dir string) {
				for _, kf := range keyFiles {
					if !kf.Popular {
						continue
					}
					slug := format.KeyFileSlug(kf.Path)
					pagePath := filepath.Join(dir, slug+file.ExtMarkdown)
					if pageErr := os.WriteFile(
						pagePath, []byte(
							section.GenerateKeyFilePage(kf)),
						fs.PermFile,
					); pageErr != nil {
						err.WarnFile(cmd, pagePath, pageErr)
					}
				}
			}); writeErr != nil {
			return writeErr
		}
	}

	// Generate session type pages
	var typeEntries []entity.JournalEntry
	for _, e := range entries {
		if e.Suggestive || section.ContinuesMultipart(e.Filename) || e.Type == "" {
			continue
		}
		typeEntries = append(typeEntries, e)
	}

	sessionTypes := section.BuildTypeIndex(typeEntries)

	if len(sessionTypes) > 0 {
		if writeErr := section.WriteSection(
			docsDir,
			dir.JournalTypes,
			section.GenerateTypesIndex(sessionTypes),
			func(dir string) {
				for _, st := range sessionTypes {
					pagePath := filepath.Join(dir, st.Name+file.ExtMarkdown)
					if pageErr := os.WriteFile(
						pagePath,
						[]byte(section.GenerateTypePage(st)), fs.PermFile,
					); pageErr != nil {
						err.WarnFile(cmd, pagePath, pageErr)
					}
				}
			}); writeErr != nil {
			return writeErr
		}
	}

	// Generate zensical.toml
	tomlContent := generate.GenerateZensicalToml(
		entries, topics, keyFiles, sessionTypes,
	)
	tomlPath := filepath.Join(output, zensical.Toml)
	if writeErr := os.WriteFile(
		tomlPath,
		[]byte(tomlContent), fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(tomlPath, writeErr)
	}

	if serve {
		writeJournal.InfoSiteStarting(cmd)
		return runZensical(output, zensical.CmdServe)
	} else if build {
		writeJournal.InfoSiteBuilding(cmd)
		return runZensical(output, zensical.CmdBuild)
	}

	writeJournal.InfoSiteGenerated(cmd, len(entries), output, zensical.Bin)

	return nil
}
