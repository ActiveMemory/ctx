//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package closeout

import (
	"path/filepath"
	"time"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	errCloseout "github.com/ActiveMemory/ctx/internal/err/closeout"
	"github.com/ActiveMemory/ctx/internal/git_meta"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Write assembles a closeout file under the supplied closeouts
// directory. The caller supplies mode (one of cfgKB.CloseoutMode*),
// optional pass-mode + life-stage (set only for ingest passes),
// and the rendered body sections.
//
// SHA and branch are read from [gitmeta.ResolveHead] against the
// project root. The generated-at timestamp is captured at write
// time in UTC.
//
// Parameters:
//   - closeoutsDir: absolute path to the closeouts directory
//     (typically `.context/ingest/closeouts/`); created if absent.
//   - projectRoot: absolute path to the project root (parent of
//     `.context/`); passed to [gitmeta.ResolveHead].
//   - mode: one of `cfgKB.CloseoutMode*` (ingest, ask, etc.).
//   - passMode: one of "topic-page" | "triage" | "evidence-only"
//     for ingest passes; empty string for other modes.
//   - lifeStage: one of "bootstrap" | "maintenance" for ingest
//     passes; empty string for other modes.
//   - body: rendered closeout body (everything after the closing
//     `---`).
//
// Returns:
//   - File: the written file with parsed frontmatter and the
//     body echoed back.
//   - error: non-nil on stat / git-resolve / write failure.
func Write(
	closeoutsDir, projectRoot, mode, passMode, lifeStage, body string,
) (File, error) {
	if mode == "" {
		return File{}, errCloseout.ErrModeRequired
	}
	ref, headErr := gitmeta.ResolveHead(projectRoot)
	if headErr != nil {
		return File{}, errCloseout.ResolveHead(headErr)
	}
	fm := Frontmatter{
		SHA:         ref.SHA,
		Branch:      ref.Branch,
		Mode:        mode,
		PassMode:    passMode,
		LifeStage:   lifeStage,
		GeneratedAt: time.Now().UTC().Truncate(time.Second),
	}

	if mkErr := io.SafeMkdirAll(closeoutsDir, cfgFs.PermExec); mkErr != nil {
		return File{}, errCloseout.MkdirCloseouts(mkErr)
	}

	rendered, renderErr := renderMarkdown(fm, body)
	if renderErr != nil {
		return File{}, renderErr
	}

	name := buildFilename(fm)
	path := filepath.Join(closeoutsDir, name)
	writeErr := io.SafeWriteFile(path, []byte(rendered), cfgFs.PermSecret)
	if writeErr != nil {
		return File{}, errCloseout.WriteFailed(writeErr)
	}

	return File{
		Path:        path,
		Frontmatter: fm,
		Body:        body,
	}, nil
}
