//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package test

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	cfgEvent "github.com/ActiveMemory/ctx/internal/config/event"
	"github.com/ActiveMemory/ctx/internal/config/project"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	errNotify "github.com/ActiveMemory/ctx/internal/err/notify"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Send loads the configured webhook, builds a test payload, and posts it.
// Returns a Result for the cmd layer to report.
//
// Returns:
//   - Result: Outcome of the test notification
//   - error: Non-nil on webhook load, marshal, or HTTP failure
func Send() (Result, error) {
	url, loadErr := notify.LoadWebhook()
	if loadErr != nil {
		return Result{}, errNotify.LoadWebhook(loadErr)
	}
	if url == "" {
		return Result{NoWebhook: true}, nil
	}

	projectName := project.FallbackName
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		projectName = filepath.Base(cwd)
	} else {
		ctxLog.Warn(warn.Getwd, cwdErr)
	}

	payload := entity.NotifyPayload{
		Event:     cfgEvent.TypeTest,
		Message:   cfgEvent.TestMessage,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Project:   projectName,
	}

	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return Result{}, errNotify.MarshalPayload(marshalErr)
	}

	filtered := !notify.EventAllowed(cfgEvent.TypeTest, rc.NotifyEvents())

	resp, postErr := notify.PostJSON(url, body)
	if postErr != nil {
		return Result{}, errNotify.SendNotification(postErr)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			ctxLog.Warn(
				warn.Close, warn.ResponseBody, closeErr,
			)
		}
	}()

	return Result{
		Filtered:   filtered,
		StatusCode: resp.StatusCode,
	}, nil
}

// OK reports whether the HTTP response indicates success.
//
// Parameters:
//   - r: Result from Send
//
// Returns:
//   - bool: True if status code is 2xx
func OK(r Result) bool {
	return r.StatusCode >= http.StatusOK &&
		r.StatusCode < http.StatusMultipleChoices
}
