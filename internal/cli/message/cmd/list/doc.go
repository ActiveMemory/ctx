//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package list implements the "ctx hook message list"
// command.
//
// # Overview
//
// The list command displays all registered hook message
// templates, showing each template's hook name, variant,
// category, and whether a local override exists. This
// gives the user a complete inventory of customizable
// messages.
//
// # Flags
//
//	--json    Output the template list as formatted
//	          JSON instead of the default table format.
//
// # Behavior
//
// [Cmd] builds the cobra.Command and registers the
// --json flag. [Run] iterates over the message registry,
// checks for local overrides, and outputs the results.
//
// In table mode, prints a header row followed by one
// row per template with columns for hook, variant,
// category, and override status.
//
// In JSON mode, encodes the full list as a JSON array
// with indentation, including template variables and
// descriptions.
//
// # Output
//
// Table mode (default):
//
//	HOOK         VARIANT    CATEGORY   OVERRIDE
//	PreToolUse   default    general    no
//	...
//
// JSON mode produces an array of objects with hook,
// variant, category, description, templateVars, and
// hasOverride fields.
package list
