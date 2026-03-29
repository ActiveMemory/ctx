//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// OversizeTokens matches "Injected:  NNNNN tokens" in the
// injection-oversize flag file.
//
// Groups:
//   - 1: token count digits
var OversizeTokens = regexp.MustCompile(`Injected:\s+(\d+)\s+tokens`)
