//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
)

//go:embed extra.css
var extraCSS []byte

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
	_, err := exec.LookPath(config.BinZensical)
	if err != nil {
		return errZensicalNotFound()
	}

	cmd := exec.Command(config.BinZensical, command) //nolint:gosec // G204: binary is a constant, command is from caller
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
	if _, err := os.Stat(journalDir); os.IsNotExist(err) {
		return errNoJournalDir(journalDir)
	}

	// Load journal state for per-file processing flags
	jstate, err := state.Load(journalDir)
	if err != nil {
		return fmt.Errorf("load journal state: %w", err)
	}

	// Scan journal files
	entries, err := scanJournalEntries(journalDir)
	if err != nil {
		return errScanJournal(err)
	}

	if len(entries) == 0 {
		return errNoEntries(journalDir)
	}

	green := color.New(color.FgGreen).SprintFunc()

	// Create output directory structure
	docsDir := filepath.Join(output, config.JournalDirDocs)
	if err = os.MkdirAll(docsDir, config.PermExec); err != nil {
		return errMkdir(docsDir, err)
	}

	// Write stylesheet for <pre> overflow control
	stylesDir := filepath.Join(docsDir, "stylesheets")
	if err = os.MkdirAll(stylesDir, config.PermExec); err != nil {
		return errMkdir(stylesDir, err)
	}
	cssPath := filepath.Join(stylesDir, "extra.css")
	if err = os.WriteFile(
		cssPath, extraCSS, config.PermFile,
	); err != nil {
		return errFileWrite(cssPath, err)
	}

	// Write README
	readmePath := filepath.Join(output, config.FilenameReadme)
	if err = os.WriteFile(
		readmePath,
		[]byte(generateSiteReadme(journalDir)), config.PermFile,
	); err != nil {
		return errFileWrite(readmePath, err)
	}

	// Soft-wrap source journal files in-place, then copy to docs/
	for _, entry := range entries {
		src := entry.Path
		dst := filepath.Join(docsDir, entry.Filename)

		var content []byte
		content, err = os.ReadFile(filepath.Clean(src))
		if err != nil {
			warnFileErr(cmd, entry.Filename, err)
			continue
		}

		// Normalize the source file for readability
		normalized := collapseToolOutputs(
			softWrapContent(
				mergeConsecutiveTurns(
					consolidateToolRuns(
						cleanToolOutputJSON(
							stripSystemReminders(string(content)),
						),
					),
				),
			),
		)
		if normalized != string(content) {
			if err = os.WriteFile(
				src, []byte(normalized), config.PermFile,
			); err != nil {
				warnFileErr(cmd, entry.Filename, err)
			}
		}

		// Generate site copy with Markdown fixes
		fv := jstate.IsFencesVerified(entry.Filename)
		withLinks := injectSourceLink(normalized, src)
		if entry.Summary != "" {
			withLinks = injectSummary(withLinks, entry.Summary)
		}
		siteContent := normalizeContent(withLinks, fv)
		if err = os.WriteFile(
			dst, []byte(siteContent), config.PermFile,
		); err != nil {
			warnFileErr(cmd, entry.Filename, err)
			continue
		}
	}

	// Generate index.md
	indexContent := generateIndex(entries)
	indexPath := filepath.Join(docsDir, config.FilenameIndex)
	if err = os.WriteFile(
		indexPath, []byte(indexContent), config.PermFile,
	); err != nil {
		return errFileWrite(indexPath, err)
	}

	// Generate topic pages
	var topicEntries []journalEntry
	for _, e := range entries {
		if e.Suggestive || continuesMultipart(e.Filename) || len(e.Topics) == 0 {
			continue
		}
		topicEntries = append(topicEntries, e)
	}

	topics := buildTopicIndex(topicEntries)

	if len(topics) > 0 {
		if err = writeSection(
			docsDir, config.JournalDirTopics,
			generateTopicsIndex(topics),
			func(dir string) {
				for _, t := range topics {
					if !t.Popular {
						continue
					}
					pagePath := filepath.Join(dir, t.Name+config.ExtMarkdown)
					if writeErr := os.WriteFile(
						pagePath, []byte(generateTopicPage(t)),
						config.PermFile,
					); writeErr != nil {
						warnFileErr(cmd, pagePath, writeErr)
					}
				}
			}); err != nil {
			return err
		}
	}

	// Generate key files pages
	var keyFileEntries []journalEntry
	for _, e := range entries {
		if e.Suggestive || continuesMultipart(e.Filename) || len(e.KeyFiles) == 0 {
			continue
		}
		keyFileEntries = append(keyFileEntries, e)
	}

	keyFiles := buildKeyFileIndex(keyFileEntries)

	if len(keyFiles) > 0 {
		if err = writeSection(
			docsDir, config.JournalDirFiles,
			generateKeyFilesIndex(keyFiles),
			func(dir string) {
				for _, kf := range keyFiles {
					if !kf.Popular {
						continue
					}
					slug := keyFileSlug(kf.Path)
					pagePath := filepath.Join(dir, slug+config.ExtMarkdown)
					if writeErr := os.WriteFile(
						pagePath, []byte(
							generateKeyFilePage(kf)),
						config.PermFile,
					); writeErr != nil {
						warnFileErr(cmd, pagePath, writeErr)
					}
				}
			}); err != nil {
			return err
		}
	}

	// Generate session type pages
	var typeEntries []journalEntry
	for _, e := range entries {
		if e.Suggestive || continuesMultipart(e.Filename) || e.Type == "" {
			continue
		}
		typeEntries = append(typeEntries, e)
	}

	sessionTypes := buildTypeIndex(typeEntries)

	if len(sessionTypes) > 0 {
		if err = writeSection(
			docsDir,
			config.JournalDirTypes,
			generateTypesIndex(sessionTypes),
			func(dir string) {
				for _, st := range sessionTypes {
					pagePath := filepath.Join(dir, st.Name+config.ExtMarkdown)
					if writeErr := os.WriteFile(
						pagePath,
						[]byte(generateTypePage(st)), config.PermFile,
					); writeErr != nil {
						warnFileErr(cmd, pagePath, writeErr)
					}
				}
			}); err != nil {
			return err
		}
	}

	// Generate zensical.toml
	tomlContent := generateZensicalToml(
		entries, topics, keyFiles, sessionTypes,
	)
	tomlPath := filepath.Join(output, config.FileZensicalToml)
	if err = os.WriteFile(
		tomlPath,
		[]byte(tomlContent), config.PermFile,
	); err != nil {
		return errFileWrite(tomlPath, err)
	}

	cmd.Println(fmt.Sprintf(
		"%s Generated site with %d entries in %s",
		green("✓"), len(entries), output,
	))

	// Build or serve if requested
	if serve {
		cmd.Println()
		cmd.Println("Starting local server...")
		return runZensical(output, "serve")
	} else if build {
		cmd.Println()
		cmd.Println("Building site...")
		return runZensical(output, "build")
	}

	cmd.Println()
	cmd.Println("Next steps:")
	cmd.Println(fmt.Sprintf("  cd %s && %s serve", output, config.BinZensical))
	cmd.Println("  or")
	cmd.Println("  ctx journal site --serve")

	return nil
}
