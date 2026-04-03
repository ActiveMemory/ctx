//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package rust implements GraphBuilder for Rust projects.
// Uses cargo metadata to parse the workspace dependency
// graph. Internal graphs show workspace member dependencies;
// full graphs include all external crate dependencies.
package rust
