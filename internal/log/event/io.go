//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// readLogFile reads and parses all events from a JSONL log file.
// Malformed lines are silently skipped.
//
// Parameters:
//   - path: absolute path to the JSONL event log
//
// Returns:
//   - []notify.Payload: parsed events in file order; nil when file
//     does not exist
//   - error: non-nil only when the file exists but cannot be opened
func readLogFile(path string) ([]notify.Payload, error) {
	f, openErr := io.SafeOpenUserFile(path)
	if openErr != nil {
		if os.IsNotExist(openErr) {
			return nil, nil
		}
		return nil, openErr
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			logWarn.Warn(warn.Close, path, closeErr)
		}
	}()

	var events []notify.Payload
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var p notify.Payload
		if unmarshalErr := json.Unmarshal(
			scanner.Bytes(), &p,
		); unmarshalErr != nil {
			continue // skip malformed lines
		}
		events = append(events, p)
	}

	return events, nil
}
