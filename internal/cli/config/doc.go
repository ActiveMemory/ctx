//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package config implements the ctx config command group for
// managing runtime configuration profiles.
//
// Subcommands: schema (output JSON Schema), status (show active
// config), switch (change profile). Profiles are .ctxrc.<name>
// files that can be swapped via ctx config switch.
package config
