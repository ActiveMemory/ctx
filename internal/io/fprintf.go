//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"fmt"
	"io"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/log/warn"
)

// SafeFprintf writes formatted output to w, logging to the warning
// sink on error.
//
// Parameters:
//   - w: destination writer
//   - format: Printf-style format string
//   - a: format arguments
func SafeFprintf(w io.Writer, format string, a ...any) {
	if _, err := fmt.Fprintf(w, format, a...); err != nil {
		warn.Warn(desc.Text(text.DescKeyErrFsWriteBuffer), err)
	}
}
