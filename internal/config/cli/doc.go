//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cli defines constants that shape the ctx command-
// line interface: cobra annotations, XML attribute names,
// binary identification, and user confirmation values.
//
// The ctx CLI is built on cobra. Several cross-cutting
// concerns (initialization guards, context-update parsing,
// stdin handling) rely on shared string constants. This
// package collects them so that command implementations
// and middleware reference the same values.
//
// # Cobra Annotations
//
//   - [AnnotationSkipInit] is the annotation key that
//     exempts a command from the PersistentPreRunE
//     initialization guard. Commands like "completion"
//     and "version" set this to bypass context loading.
//   - [AnnotationTrue] is the canonical "true" value for
//     boolean annotations.
//
// # XML Attribute Names
//
// ctx supports structured context updates via XML tags
// in AI output. The attribute constants name the fields
// parsed from <context-update> tags:
//
//   - [AttrType], [AttrContext], [AttrLesson],
//     [AttrApplication], [AttrSection], [AttrRationale],
//     [AttrConsequence]
//
// These are consumed by the hook output parser to route
// updates to the correct context file.
//
// # Binary and Input
//
//   - [Binary] ("ctx") is the executable name used for
//     PATH validation during setup.
//   - [CmdCompletion] names the shell completion
//     subcommand for annotation checks.
//   - [StdinSentinel] ("-") is the conventional argument
//     meaning "read from stdin", used by commands that
//     accept piped input.
//
// # User Confirmation
//
//   - [ConfirmShort] ("y") and [ConfirmLong] ("yes") are
//     the accepted affirmative responses for interactive
//     y/N prompts.
//
// # Why Centralized
//
// Annotation keys and attribute names appear in both
// command registration and middleware. Centralizing them
// prevents silent mismatches that would cause commands
// to skip initialization or drop parsed attributes.
package cli
