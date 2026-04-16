//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package serve implements the "ctx serve" command for
// serving static sites locally.
//
// The serve command starts a local HTTP server that hosts
// the journal site or any specified directory. It uses
// zensical as the static site engine and opens the
// default browser to the served URL on startup.
//
// # Use Cases
//
//   - Preview journal site output from ctx journal site
//     before deploying
//   - Serve any directory of static HTML for local
//     testing
//   - Quick local web server for documentation review
//
// # Subpackages
//
//   - cmd/root: cobra command definition, HTTP server
//     setup, and browser launch
package serve
