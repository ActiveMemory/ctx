//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// CtxGoRun matches go run ./cmd/ctx.
var CtxGoRun = regexp.MustCompile(`go run \./cmd/ctx`)

// CtxAbsoluteStart matches absolute paths to ctx at start of command.
var CtxAbsoluteStart = regexp.MustCompile(`^\s*(/home/|/tmp/|/var/)\S*/ctx(\s|$)`)

// AbsoluteSep matches absolute paths to ctx after command separator.
var AbsoluteSep = regexp.MustCompile(`(&&|;|\|\||\|)\s*(/home/|/tmp/|/var/)\S*/ctx(\s|$)`)

// CtxTestException matches /tmp/ctx-test for integration test exemption.
var CtxTestException = regexp.MustCompile(`/tmp/ctx-test`)

// CtxRelativeSep matches ./ctx or ./dist/ctx after command separator.
var CtxRelativeSep = regexp.MustCompile(`(&&|;|\|\||\|)\s*(\./ctx(\s|$)|\./dist/ctx)`)

// CtxRelativeStart matches ./ctx or ./dist/ctx at start of command.
var CtxRelativeStart = regexp.MustCompile(`^\s*(\./ctx(\s|$)|\./dist/ctx)`)
