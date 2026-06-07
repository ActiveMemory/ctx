//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pass

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	dreamPaths "github.com/ActiveMemory/ctx/internal/cli/dream/core/paths"
	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	cfgToken "github.com/ActiveMemory/ctx/internal/config/token"
	engine "github.com/ActiveMemory/ctx/internal/dream"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
	execDream "github.com/ActiveMemory/ctx/internal/exec/dream"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeDream "github.com/ActiveMemory/ctx/internal/write/dream"
)

// minutesPerStep bounds the executor wall-clock budget per step so a
// runaway loop cannot hang a headless cron pass indefinitely.
const minutesPerStep = 2

// dreamDue is the auto-trigger gate honored unless the run is forced.
// The dream must be opt-in enabled; when a cron cadence is configured,
// a fresh pass is deferred until the quiet window has elapsed since the
// latest run, so back-to-back auto passes do not thrash. A forced
// manual run bypasses this entirely.
//
// Parameters:
//   - dreamsDir: the dreams/ notebook directory
//
// Returns:
//   - bool: true when an auto pass is due to proceed
func dreamDue(dreamsDir string) bool {
	if !rc.DreamEnabled() {
		return false
	}
	if rc.DreamCadence() == "" {
		return true
	}
	last, lastErr := engine.LatestRunDir(dreamsDir)
	if lastErr != nil || last == "" {
		return true
	}
	stamp, parseErr := time.Parse(
		cfgDream.RunTimeLayout, filepath.Base(last),
	)
	if parseErr != nil {
		return true
	}
	quiet := time.Duration(rc.DreamQuietMinutes()) * time.Minute
	return time.Since(stamp) >= quiet
}

// selectSources scans ideas/, computes the delta against saved state,
// and bounds the result to at most maxFiles for this pass.
//
// Parameters:
//   - loc: the resolved dream working locations
//   - maxFiles: ceiling on sources for this pass
//
// Returns:
//   - []string: the bounded, sorted delta paths to process
//   - error: a scan or state-read failure
func selectSources(
	loc dreamPaths.Resolved, maxFiles int,
) ([]string, error) {
	current, scanErr := engine.ScanIdeas(loc.Root, loc.Ideas)
	if scanErr != nil {
		return nil, scanErr
	}
	prior, stateErr := engine.LoadState(loc.Dreams)
	if stateErr != nil {
		return nil, stateErr
	}
	selected := engine.DeltaSelect(prior, current)
	if maxFiles > 0 && len(selected) > maxFiles {
		selected = selected[:maxFiles]
	}
	return selected, nil
}

// persist merges the processed sources into saved state, marking each
// active with its current hash so the discipline clock skips it until
// the content changes.
//
// Parameters:
//   - loc: the resolved dream working locations
//   - processed: the source paths handled this pass
//
// Returns:
//   - error: a scan, state-read, or state-write failure
func persist(loc dreamPaths.Resolved, processed []string) error {
	current, scanErr := engine.ScanIdeas(loc.Root, loc.Ideas)
	if scanErr != nil {
		return scanErr
	}
	prior, stateErr := engine.LoadState(loc.Dreams)
	if stateErr != nil {
		return stateErr
	}
	byPath := make(map[string]engine.SourceState, len(prior))
	for _, s := range prior {
		byPath[s.Path] = s
	}
	now := time.Now().UTC()
	for _, p := range processed {
		rec := byPath[p]
		rec.Path = p
		rec.Hash = current[p]
		rec.LastModified = now
		rec.LastSurfaced = now
		if rec.Status == "" {
			rec.Status = cfgDream.SourceActive
		}
		byPath[p] = rec
	}
	merged := make([]engine.SourceState, 0, len(byPath))
	for _, s := range byPath {
		merged = append(merged, s)
	}
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Path < merged[j].Path
	})
	return engine.SaveState(loc.Dreams, merged)
}

// invoke resolves and runs the configured executor for one bounded
// pass. A missing executor binary is a fail-loud error; the caller
// writes the failmark. The executor reads ideas/ and writes proposals
// into runDir.
//
// Parameters:
//   - cmd: cobra command for the executor's stdout/stderr wiring
//   - loc: the resolved dream working locations
//   - runDir: the per-run directory the executor writes proposals into
//   - opts: the resolved run parameters
//
// Returns:
//   - error: ExecutorNotFound or ExecutorRun on failure; nil on a
//     clean pass
func invoke(
	cmd *cobra.Command, loc dreamPaths.Resolved,
	runDir string, opts Opts,
) error {
	bin := cfgDream.ExecutorDefaultBin
	if override := rc.DreamExecutor(); override != "" {
		bin = override
	}
	resolved, lookErr := execDream.LookPath(bin)
	if lookErr != nil {
		return errDream.ExecutorNotFound(bin, lookErr)
	}

	budget := opts.Budget
	if budget <= 0 {
		budget = cfgDream.DefaultBudget
	}
	prompt := fmt.Sprintf(
		cfgDream.ExecutorPromptTemplate,
		opts.Mode, opts.Max, budget, loc.Ideas, runDir,
	)
	timeout := time.Duration(budget*minutesPerStep) * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	args := []string{
		cfgDream.ExecutorPromptFlag, prompt,
		cfgDream.ExecutorMaxTurnsFlag, strconv.Itoa(budget),
	}
	if model := rc.DreamModel(); model != "" {
		args = append(args, cfgDream.ExecutorModelFlag, model)
	}
	c := execDream.CommandContext(ctx, resolved, args...)
	c.Dir = loc.Root
	c.Stdout = cmd.OutOrStdout()
	c.Stderr = cmd.ErrOrStderr()
	if runErr := c.Run(); runErr != nil {
		return errDream.ExecutorRun(bin, runErr)
	}
	return nil
}

// validateRun reads the proposals the executor wrote into runDir and
// validates each against the proposal schema, returning the count of
// valid proposals.
//
// Parameters:
//   - runDir: the per-run directory the executor wrote into
//
// Returns:
//   - int: number of schema-valid proposals
//   - error: a read failure or the first invalid-proposal error
func validateRun(runDir string) (int, error) {
	proposals, readErr := engine.ReadProposals(runDir)
	if readErr != nil {
		return 0, readErr
	}
	for _, p := range proposals {
		if validErr := engine.ProposalValid(p); validErr != nil {
			return 0, validErr
		}
	}
	return len(proposals), nil
}

// writeFailmark records the fail-loud failmark under dreams/ and
// reports it, so a missing/failed executor never silently no-ops.
//
// Parameters:
//   - cmd: cobra command for output
//   - dreamsDir: the dreams/ notebook directory
func writeFailmark(cmd *cobra.Command, dreamsDir string) {
	path := filepath.Join(dreamsDir, cfgDream.FileFailed)
	stamp := time.Now().UTC().Format(cfgTime.RFC3339Compact) +
		cfgToken.NewlineLF
	if writeErr := ctxIo.SafeWriteFile(
		path, []byte(stamp), cfgFs.PermSecret,
	); writeErr == nil {
		writeDream.Failmark(cmd, path)
	}
}
