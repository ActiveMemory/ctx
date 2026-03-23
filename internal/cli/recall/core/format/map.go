//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import "github.com/ActiveMemory/ctx/internal/config/session"

// toolDisplayKey maps tool names to the JSON input key that best
// describes each invocation.
var toolDisplayKey = map[string]string{
	session.ToolRead:      session.ToolInputFilePath,
	session.ToolWrite:     session.ToolInputFilePath,
	session.ToolEdit:      session.ToolInputFilePath,
	session.ToolBash:      session.ToolInputCommand,
	session.ToolGrep:      session.ToolInputPattern,
	session.ToolGlob:      session.ToolInputPattern,
	session.ToolWebFetch:  session.ToolInputURL,
	session.ToolWebSearch: session.ToolInputQuery,
	session.ToolTask:      session.ToolInputDescription,
}
