//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package node implements GraphBuilder for Node.js projects.
// Parses package.json for both single-package and monorepo
// workspace layouts. Internal graphs show workspace-to-workspace
// dependencies; full graphs include all external packages.
package node
