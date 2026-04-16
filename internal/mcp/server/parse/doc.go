//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parse unmarshals raw JSON bytes into MCP
// request structs.
//
// # Request Parsing
//
// Request takes raw JSON bytes from stdin and returns
// a parsed proto.Request. It handles three cases:
//
//   - Valid request with ID: returns the parsed
//     request and nil error response.
//   - Valid notification (no ID): returns nil for
//     both, since notifications expect no response.
//   - Malformed JSON: returns nil request and a
//     parse-error response ready to send back.
//
// # Usage
//
//	req, errResp := parse.Request(data)
//	if errResp != nil {
//	    // send error response
//	}
//	if req == nil {
//	    // notification, skip
//	}
//
// # Error Codes
//
// Parse errors use the standard JSON-RPC parse error
// code from the schema config package. The error
// message comes from the embedded description text.
package parse
