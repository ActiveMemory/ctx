//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

// BlockResponse is the JSON output for blocked commands.
type BlockResponse struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}
