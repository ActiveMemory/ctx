//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package extract

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/backend"
	coreExtract "github.com/ActiveMemory/ctx/internal/cli/ai/core/extract"
	"github.com/ActiveMemory/ctx/internal/cli/ai/core/resolve"
	cfgExtract "github.com/ActiveMemory/ctx/internal/config/extract"
	errProposal "github.com/ActiveMemory/ctx/internal/err/proposal"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeAI "github.com/ActiveMemory/ctx/internal/write/ai"
	writeProposal "github.com/ActiveMemory/ctx/internal/write/proposal"
)

// Run executes `ctx ai extract`: reads input text from
// the command's stdin, dispatches a JSON-mode chat
// completion through the named or default backend, and
// writes the response to a new proposal file under
// `.context/proposals/`.
//
// Fails closed on unconfigured registry, empty input,
// transport / upstream failure, and proposal-write
// failure. `.context/*.md` files are never written.
//
// Parameters:
//   - cobraCmd: cobra command for stdin / stdout / ctx.
//   - backendName: --backend value; empty falls back to
//     registry default.
//
// Returns:
//   - error: typed sentinel on failure.
func Run(cobraCmd *cobra.Command, backendName string) error {
	body, readErr := io.ReadAll(cobraCmd.InOrStdin())
	if readErr != nil {
		return errProposal.ReadInput(readErr)
	}
	text := strings.TrimSpace(string(body))
	if text == "" {
		return errProposal.EmptyInput()
	}
	r, buildErr := resolve.Build()
	if buildErr != nil {
		return buildErr
	}
	b, pickErr := resolve.Pick(r, backendName)
	if pickErr != nil {
		return pickErr
	}
	resp, completeErr := b.Complete(cobraCmd.Context(), backend.Request{
		Messages: []backend.Message{
			{Role: cfgExtract.RoleSystem, Content: cfgExtract.SystemPrompt},
			{Role: cfgExtract.RoleUser, Content: text},
		},
		ResponseFormat: &backend.ResponseFormat{
			Type: cfgExtract.ResponseFormatType,
		},
		Temperature: cfgExtract.UnsetTemperature,
	})
	if completeErr != nil {
		return completeErr
	}
	ctxDir, dirErr := rc.ContextDir()
	if dirErr != nil {
		return dirErr
	}
	path, writeErr := writeProposal.Write(
		ctxDir, cfgExtract.ProposalSlug,
		coreExtract.Compose(b.Name(), resp.Model, resp.Content),
	)
	if writeErr != nil {
		return writeErr
	}
	writeAI.InfoExtractWritten(cobraCmd, filepath.ToSlash(path))
	return nil
}
