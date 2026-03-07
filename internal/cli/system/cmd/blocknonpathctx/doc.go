//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package blocknonpathctx implements the ctx system block-non-path-ctx
// subcommand.
//
// It blocks non-PATH ctx invocations such as ./ctx, go run ./cmd/ctx,
// and absolute-path ctx calls to enforce consistent binary usage.
package blocknonpathctx
