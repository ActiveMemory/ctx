//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handover

import (
	"path/filepath"
	"strings"
	"time"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/entity"
	errHandover "github.com/ActiveMemory/ctx/internal/err/handover"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/write/closeout"
)

// Write writes a new handover file under handoversDir. When
// entry.NoFold is false, it also reads the closeouts under
// closeoutsDir, finds those postdating the latest handover,
// folds their summaries into the new handover's body, and
// archives them under archiveDir.
//
// Parameters:
//   - handoversDir: absolute path to .context/handovers/.
//   - closeoutsDir: absolute path to .context/ingest/closeouts/.
//   - archiveDir: absolute path to .context/archive/closeouts/.
//   - projectRoot: absolute path to the project root; passed to
//     [github.com/ActiveMemory/ctx/internal/gitmeta.ResolveHead].
//   - entry: caller-supplied content + flags.
//
// Returns:
//   - Result: written file + folded + malformed metadata.
//   - error: non-nil on git-resolve / I/O / archive failure.
func Write(
	handoversDir, closeoutsDir, archiveDir, projectRoot string,
	entry Entry,
) (Result, error) {
	if strings.TrimSpace(entry.Title) == "" {
		return Result{}, errHandover.ErrTitleRequired
	}
	if strings.TrimSpace(entry.Summary) == "" {
		return Result{}, errHandover.ErrSummaryRequired
	}
	if strings.TrimSpace(entry.Next) == "" {
		return Result{}, errHandover.ErrNextRequired
	}

	sha, branch, provErr := resolveProvenance(projectRoot, entry.CommitOverride)
	if provErr != nil {
		return Result{}, provErr
	}

	now := time.Now().UTC().Truncate(time.Second)
	fm := Frontmatter{
		SHA:         sha,
		Branch:      branch,
		GeneratedAt: now,
		Title:       entry.Title,
	}

	var folded []entity.CloseoutFile
	var malformed []string
	if !entry.NoFold {
		latestAt, _, latestErr := Latest(handoversDir)
		if latestErr != nil {
			return Result{}, errHandover.Latest(latestErr)
		}
		all, bad, listErr := closeout.List(closeoutsDir)
		if listErr != nil {
			return Result{}, errHandover.ListCloseouts(listErr)
		}
		folded = closeout.PostdatedBy(all, latestAt)
		malformed = bad
	}

	body := renderBody(entry, folded)
	rendered, composeErr := composeMarkdown(fm, body)
	if composeErr != nil {
		return Result{}, composeErr
	}

	if mkErr := io.SafeMkdirAll(handoversDir, cfgFs.PermExec); mkErr != nil {
		return Result{}, errHandover.MkdirHandovers(mkErr)
	}
	name := buildFilename(now, entry.Title)
	path := filepath.Join(handoversDir, name)
	writeErr := io.SafeWriteFile(path, []byte(rendered), cfgFs.PermSecret)
	if writeErr != nil {
		return Result{}, errHandover.WriteFailed(writeErr)
	}

	if !entry.NoFold && len(folded) > 0 {
		if archErr := closeout.Archive(archiveDir, folded); archErr != nil {
			return Result{}, errHandover.ArchiveFoldedCloseouts(archErr)
		}
	}

	return Result{
		File: File{
			Path:        path,
			Frontmatter: fm,
			Body:        body,
		},
		FoldedCloseouts:    folded,
		MalformedCloseouts: malformed,
	}, nil
}
