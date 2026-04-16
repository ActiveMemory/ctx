//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package zensical defines site configuration, TOML
// formatting, and MkDocs stripping constants for the
// zensical static site generator integration.
//
// ctx uses zensical (a MkDocs-compatible static site
// generator) to build both the main documentation
// site and per-project journal sites. This package
// provides the config filename, binary name, build
// commands, TOML nav-array formatting tokens, and
// MkDocs admonition/tab stripping constants.
//
// # Site Configuration
//
//   - [Toml] ("zensical.toml"): the site config
//     filename.
//   - [Bin] ("zensical"): the binary name.
//   - [CmdServe], [CmdBuild]: subcommand names
//     for serving and building.
//   - [Stylesheets], [ExtraCSS]: CSS asset paths
//     for journal site styling.
//
// # TOML Nav Formatting
//
//   - [TomlNavOpen]: opening bracket for
//     the nav array.
//   - [TomlNavSectionClose]: closing bracket for a
//     nav section group.
//   - [TomlNavClose]: closing bracket for
//     the top-level nav array.
//
// # MkDocs Stripping
//
// The ctx why command renders embedded docs that
// may contain MkDocs-specific syntax. These
// constants control the stripping:
//
//   - [MkDocsAdmonitionPrefix] ("!!!"): admonition
//     marker.
//   - [MkDocsTabPrefix] ("=== "): tab marker.
//   - [MkDocsIndent], [MkDocsIndentWidth]: body
//     indentation to remove.
//   - [MkDocsFrontmatterDelim] ("---"): YAML
//     frontmatter fence.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package zensical
