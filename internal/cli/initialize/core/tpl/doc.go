//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tpl handles template deployment during the
// ctx init pipeline.
//
// # Overview
//
// This package provides a generic engine for copying
// embedded template files into the user's .context/
// directory. It is used by the entry and other init
// sub-packages to deploy context file templates.
//
// # Behavior
//
// [DeployTemplates] creates a target subdirectory under
// .context/ and writes each embedded template file into
// it, skipping files that already exist unless force mode
// is enabled.
//
// # Data Flow
//
// When [DeployTemplates] is called it:
//
//  1. Creates the target subdirectory under contextDir
//     (e.g. .context/templates/) with executable
//     permissions.
//  2. Calls the provided list function to enumerate
//     all embedded template names.
//  3. For each template, checks whether the target
//     file already exists. If it exists and force is
//     false, the file is skipped with a diagnostic.
//  4. Calls the provided read function to obtain the
//     template content.
//  5. Writes the content to the target path using
//     safe file I/O.
//  6. Reports each created or skipped file via the
//     initialize write layer.
//
// The list and read functions are injected as
// parameters, making DeployTemplates reusable across
// different template sets (entries, prompts, etc.).
package tpl
