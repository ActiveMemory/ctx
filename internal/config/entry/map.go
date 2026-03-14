//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"github.com/ActiveMemory/ctx/internal/config/ctx"
)

// ToCtxFile maps short names to actual file names.
var ToCtxFile = map[string]string{
	Decision:   ctx.Decision,
	Task:       ctx.Task,
	Learning:   ctx.Learning,
	Convention: ctx.Convention,
}
