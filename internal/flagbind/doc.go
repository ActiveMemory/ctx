//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package flagbind provides helpers for cobra flag
// registration that enforce the YAML-backed description
// pipeline.
//
// All cobra flag registration must go through this
// package. Direct calls to cobra's Flags().StringVar,
// Flags().BoolVar, and similar methods are prohibited
// outside flagbind. This ensures every flag description
// is routed through [desc.Flag] rather than hardcoded
// inline, keeping flag text localizable and consistent.
//
// # Single-Flag Helpers
//
// Each helper accepts a descKey that maps to a YAML
// entry in internal/assets/commands/flags.yaml. Flag
// name constants come from internal/config/flag.
//
//   - [BoolFlag], [BoolFlagP] register boolean flags
//     with optional shorthand, defaulting to false.
//   - [BoolFlagDefault] registers a boolean flag with
//     a non-false default value.
//   - [BoolFlagNoPtr], [BoolFlagShort] register flags
//     retrieved later via cmd.Flags().GetBool().
//   - [IntFlag], [IntFlagP] register integer flags.
//   - [DurationFlag] registers a time.Duration flag.
//   - [StringFlag], [StringFlagP] register string
//     flags with optional shorthand.
//   - [StringFlagDefault], [StringFlagPDefault]
//     register string flags with non-empty defaults.
//   - [StringFlagShort] registers a no-pointer string
//     flag with shorthand.
//   - [StringArrayFlagP] registers repeatable string
//     flags (--tag x --tag y).
//   - [PersistentBoolFlag] registers a persistent bool
//     flag inherited by children.
//   - [LastJSON] registers the --last/--json pair for
//     list-style commands.
//
// # Batch Helpers
//
// Batch functions register multiple flags of the same
// kind in a single call via parallel slices:
//
//   - [BindStringFlagsP], [BindStringFlags]
//   - [BindBoolFlags], [BindBoolFlagsP]
//   - [BindStringFlagShorts]
//   - [BindStringFlagsPDefault]
//
// All slice arguments must have matching lengths;
// each index produces one single-flag call.
package flagbind
