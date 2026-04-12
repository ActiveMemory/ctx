//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package flagbind provides helpers for cobra flag registration.
//
// All cobra flag registration must go through this package. Direct calls
// to cobra's Flags().StringVar, Flags().BoolVar, and similar methods are
// prohibited outside flagbind. This ensures every flag description is
// routed through the YAML-backed assets pipeline ([desc.Flag]) rather
// than hardcoded inline, keeping flag text localizable and consistent.
//
// Each helper accepts a descKey that maps to a YAML entry in
// internal/assets/commands/flags.yaml. Flag name constants come from
// internal/config/flag.
//
// Key exports: [BoolFlag], [BoolFlagP], [IntFlagP],
// [StringFlag], [StringFlagP], [StringFlagPDefault], [LastJSON].
//
// Batch helpers ([BindStringFlagsP], [BindStringFlags],
// [BindBoolFlags], [BindBoolFlagsP], [BindStringFlagShorts],
// [BindStringFlagsPDefault]) register multiple flags of the
// same kind in a single call via parallel slices, replacing
// repetitive one-at-a-time registrations.
package flagbind
