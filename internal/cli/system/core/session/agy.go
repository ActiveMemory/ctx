//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/agy"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

var (
	// agyMu guards per-invocation Antigravity session state.
	agyMu sync.Mutex
	// agyMode indicates whether the current session is Antigravity.
	agyMode bool
	// agySessionID holds the Antigravity conversation or session ID.
	agySessionID string
	// agyTranscriptPath holds the path to the transcript JSONL.
	agyTranscriptPath string
	// agyModelName holds the active Antigravity model name.
	agyModelName string
	// agySteps holds the queued system messages for Antigravity injectSteps.
	agySteps []string
	// agyOutputWritten tracks whether hook output was already emitted.
	agyOutputWritten bool
	// agyProvenance stores the buffered provenance line in Antigravity mode.
	agyProvenance string
)

// SetAgyMode updates the per-invocation Antigravity state.
//
// Parameters:
//   - mode: true to enable Antigravity mode
//   - sessionID: session identifier
//   - transcriptPath: path to transcript JSONL
//   - modelName: active model name
func SetAgyMode(mode bool, sessionID, transcriptPath, modelName string) {
	agyMu.Lock()
	defer agyMu.Unlock()
	agyMode = mode
	agySessionID = sessionID
	agyTranscriptPath = transcriptPath
	agyModelName = modelName
	agySteps = nil
	agyOutputWritten = false
	agyProvenance = ""
}

// Agy returns true if the current session is in Antigravity mode.
//
// Returns:
//   - bool: true if Antigravity mode is active
func Agy() bool {
	agyMu.Lock()
	defer agyMu.Unlock()
	return agyMode
}

// AgyTranscriptPath returns the active Antigravity transcript path.
//
// Returns:
//   - string: transcript JSONL path
func AgyTranscriptPath() string {
	agyMu.Lock()
	defer agyMu.Unlock()
	return agyTranscriptPath
}

// AgyModelName returns the active Antigravity model name.
//
// Returns:
//   - string: model name
func AgyModelName() string {
	agyMu.Lock()
	defer agyMu.Unlock()
	return agyModelName
}

// AgyOutputWritten returns true if Antigravity output has been emitted.
//
// Returns:
//   - bool: true if output has been written
func AgyOutputWritten() bool {
	agyMu.Lock()
	defer agyMu.Unlock()
	return agyOutputWritten
}

// SetAgyOutputWritten updates the output written flag.
//
// Parameters:
//   - written: true if output has been emitted
func SetAgyOutputWritten(written bool) {
	agyMu.Lock()
	defer agyMu.Unlock()
	agyOutputWritten = written
}

// SetAgyProvenance buffers a provenance line for inclusion in the
// next injected step in Antigravity mode.
//
// Parameters:
//   - prov: Formatted provenance string
func SetAgyProvenance(prov string) {
	agyMu.Lock()
	defer agyMu.Unlock()
	agyProvenance = prov
}

// AddAgyStep queues a message step for Antigravity injectSteps output.
// If provenance was previously buffered via [SetAgyProvenance], it is
// prepended to the step and cleared.
//
// Parameters:
//   - step: Formatted message string to add
func AddAgyStep(step string) {
	if step == "" {
		return
	}
	agyMu.Lock()
	defer agyMu.Unlock()
	if agyProvenance != "" {
		step = agyProvenance + token.NewlineLF + token.NewlineLF + step
		agyProvenance = ""
	}
	agySteps = append(agySteps, step)
}

// DrainAgySteps returns all queued steps and resets the slice.
//
// Returns:
//   - []string: slice of queued step messages
func DrainAgySteps() []string {
	agyMu.Lock()
	defer agyMu.Unlock()
	if len(agySteps) == 0 {
		return nil
	}
	out := agySteps
	agySteps = nil
	return out
}

// QueuePendingAgy appends a nudge to the session's pending queue file.
//
// Parameters:
//   - msg: Nudge text to queue
func QueuePendingAgy(msg string) {
	if msg == "" {
		return
	}
	agyMu.Lock()
	mode := agyMode
	sessionID := agySessionID
	agyMu.Unlock()

	if !mode {
		return
	}
	initialized, initErr := state.Initialized()
	if initErr != nil || !initialized {
		return
	}
	stateDir, dirErr := state.Dir()
	if dirErr != nil {
		return
	}
	if sessionID == "" {
		sessionID = cfgSession.IDUnknown
	}
	path := filepath.Join(
		stateDir,
		agy.PrefixPendingQueue+sessionID+file.ExtTxt,
	)
	payload := []byte(msg + agy.PendingDelimiter)
	_ = internalIo.AppendBytes(path, payload, fs.PermSecret)
}

// DrainPendingAgy reads and removes all pending nudges for the session.
//
// Returns:
//   - []string: List of queued nudge messages
func DrainPendingAgy() []string {
	agyMu.Lock()
	mode := agyMode
	sessionID := agySessionID
	agyMu.Unlock()

	if !mode {
		return nil
	}
	initialized, initErr := state.Initialized()
	if initErr != nil || !initialized {
		return nil
	}
	stateDir, dirErr := state.Dir()
	if dirErr != nil {
		return nil
	}
	if sessionID == "" {
		sessionID = cfgSession.IDUnknown
	}
	path := filepath.Join(
		stateDir,
		agy.PrefixPendingQueue+sessionID+file.ExtTxt,
	)
	data, readErr := internalIo.SafeReadUserFile(path)
	if readErr != nil || len(data) == 0 {
		return nil
	}
	_ = os.Remove(path)

	raw := string(data)
	parts := strings.Split(raw, agy.PendingDelimiter)
	var steps []string
	seen := make(map[string]struct{})
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		if _, dup := seen[trimmed]; dup {
			continue
		}
		seen[trimmed] = struct{}{}
		steps = append(steps, trimmed)
	}
	return steps
}

// FormatAgyInjectSteps builds an Antigravity injectSteps JSON response.
// Returns "{}" when steps is empty.
//
// Parameters:
//   - steps: Slice of system message strings to inject
//
// Returns:
//   - string: JSON-encoded injectSteps response or "{}"
func FormatAgyInjectSteps(steps []string) string {
	if len(steps) == 0 {
		return agy.EmptyJSON
	}
	items := make([]agyInjectStep, 0, len(steps))
	for _, s := range steps {
		if s == "" {
			continue
		}
		items = append(items, agyInjectStep{
			SystemMessage: agySystemMsg{
				SystemMessage: s,
			},
		})
	}
	if len(items) == 0 {
		return agy.EmptyJSON
	}
	resp := agyInjectStepsResponse{
		InjectSteps: items,
	}
	data, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		return agy.EmptyJSON
	}
	return string(data)
}

// FormatAgyDecision builds an Antigravity PreToolUse decision response.
//
// Parameters:
//   - decision: "allow" or "deny"
//   - reason: Explanation for denial (omitted if empty)
//
// Returns:
//   - string: JSON-encoded decision response
func FormatAgyDecision(decision, reason string) string {
	resp := agyDecisionResponse{
		Decision: decision,
		Reason:   reason,
	}
	data, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		return agy.EmptyJSON
	}
	return string(data)
}
