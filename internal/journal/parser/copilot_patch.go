//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"encoding/json"
	"strconv"

	cfgCopilot "github.com/ActiveMemory/ctx/internal/config/copilot"
)

// applyScalarPatch applies a kind=1 scalar patch to the session.
// These update individual properties like result, modelState, followups.
//
// Parameters:
//   - session: the session to patch
//   - keys: raw JSON key path from the JSONL line
//   - value: raw JSON value to apply
func (p *Copilot) applyScalarPatch(
	session *copilotRawSession, keys []json.RawMessage, value json.RawMessage,
) {
	path := p.parseKeyPath(keys)
	if len(path) < 2 {
		return
	}

	// Handle requests.<N>.result patches — these contain token counts
	if path[0] == cfgCopilot.KeyRequests &&
		len(path) == 3 && path[2] == cfgCopilot.KeyResult {
		idx, parseErr := strconv.Atoi(path[1])
		if parseErr != nil || idx < 0 || idx >= len(session.Requests) {
			return
		}
		var result copilotRawResult
		if unmarshalErr := json.Unmarshal(
			value, &result,
		); unmarshalErr == nil {
			session.Requests[idx].Result = &result
		}
	}
}

// applyPatch applies a kind=2 array/object patch to the session.
//
// Parameters:
//   - session: the session to patch
//   - keys: raw JSON key path from the JSONL line
//   - value: raw JSON value to apply
func (p *Copilot) applyPatch(
	session *copilotRawSession, keys []json.RawMessage, value json.RawMessage,
) {
	path := p.parseKeyPath(keys)
	if len(path) == 0 {
		return
	}

	switch {
	case len(path) == 1 && path[0] == cfgCopilot.KeyRequests:
		// New request(s) appended
		var requests []copilotRawRequest
		if unmarshalReqErr := json.Unmarshal(
			value, &requests,
		); unmarshalReqErr == nil {
			session.Requests = append(
				session.Requests, requests...,
			)
		}

	case len(path) == 3 &&
		path[0] == cfgCopilot.KeyRequests &&
		path[2] == cfgCopilot.KeyResponse:
		// Response update for a specific request
		idx, parseErr := strconv.Atoi(path[1])
		if parseErr != nil || idx < 0 || idx >= len(session.Requests) {
			return
		}
		var items []copilotRawRespItem
		if unmarshalItemsErr := json.Unmarshal(
			value, &items,
		); unmarshalItemsErr == nil {
			session.Requests[idx].Response = items
		}
	}
}

// parseKeyPath converts the K array from JSONL into string path segments.
//
// Parameters:
//   - keys: raw JSON key elements to decode
//
// Returns:
//   - []string: decoded path segments as strings
func (p *Copilot) parseKeyPath(keys []json.RawMessage) []string {
	path := make([]string, 0, len(keys))
	for _, k := range keys {
		var s string
		if unmarshalStrErr := json.Unmarshal(
			k, &s,
		); unmarshalStrErr == nil {
			path = append(path, s)
			continue
		}
		var n int
		if unmarshalIntErr := json.Unmarshal(
			k, &n,
		); unmarshalIntErr == nil {
			path = append(path, strconv.Itoa(n))
			continue
		}
	}
	return path
}
