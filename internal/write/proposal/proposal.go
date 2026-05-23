//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package proposal

import (
	"path/filepath"
	"time"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgProposal "github.com/ActiveMemory/ctx/internal/config/proposal"
	errProposal "github.com/ActiveMemory/ctx/internal/err/proposal"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Write persists body under
// `<ctxDir>/proposals/<TS>-<slug>.md`. Creates the
// `proposals/` subdirectory on first call. Returns the
// absolute path of the file written so the caller can
// surface it to the user.
//
// Parameters:
//   - ctxDir: the project's `.context/` directory
//     (typically rc.ContextDir()).
//   - slug: free-text slug; normalised to kebab-case.
//     Empty falls back to [cfgProposal.DefaultSlug].
//   - body: the markdown body to write. No frontmatter
//     is injected by this package.
//
// Returns:
//   - string: absolute path of the proposal file.
//   - error: typed err/proposal sentinel on mkdir or
//     write failure.
func Write(ctxDir, slug, body string) (string, error) {
	dir := filepath.Join(ctxDir, cfgProposal.Subdir)
	mkErr := ctxIo.SafeMkdirAll(dir, cfgFs.PermExec)
	if mkErr != nil {
		return "", errProposal.MkdirProposals(dir, mkErr)
	}
	name := buildFilename(time.Now(), slug)
	full := filepath.Join(dir, name)
	wErr := ctxIo.SafeWriteFile(full, []byte(body), cfgFs.PermFile)
	if wErr != nil {
		return "", errProposal.Write(full, wErr)
	}
	return full, nil
}
