//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package project provides access to project-root
// files and directory README templates from embedded
// assets.
//
// # README Templates
//
// Readme reads a directory-specific README template
// by directory name. Templates are stored as
// project/<dir>-README.md in the embedded filesystem.
//
//	content, err := project.Readme("specs")
//	content, err := project.Readme("ideas")
//
// # Deployment
//
// During ctx init, README templates are deployed into
// project subdirectories (specs/, ideas/, etc.) to
// provide guidance on how each directory should be
// used. The directory name is sanitized via path.Base
// before constructing the asset path.
package project
