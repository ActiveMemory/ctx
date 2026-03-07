//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import "github.com/ActiveMemory/ctx/internal/cli/add/core"

// Params is the shared entry parameter type.
//
// Aliased from add/core where the struct is defined, to avoid import
// cycles (entry imports add/core for insert logic, so add/core owns
// the struct).
type Params = core.EntryParams
