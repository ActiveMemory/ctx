//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package flagbind

import "github.com/spf13/cobra"

// BindStringFlagsP registers multiple string flags that each
// have a shorthand letter. All four slices must have the same
// length; each index i produces one [StringFlagP] call. Use
// this to replace repetitive sequences of individual
// [StringFlagP] registrations.
//
// Parameters:
//   - c: Cobra command to register on
//   - ptrs: Pointers to the string variables
//   - names: Flag name constants
//   - shorts: Shorthand letters
//   - descKeys: YAML DescKeys for the flag descriptions
func BindStringFlagsP(
	c *cobra.Command,
	ptrs []*string,
	names, shorts, descKeys []string,
) {
	for i, p := range ptrs {
		StringFlagP(c, p, names[i], shorts[i], descKeys[i])
	}
}

// BindStringFlags registers multiple string flags that have
// no shorthand letter. All three slices must have the same
// length; each index i produces one [StringFlag] call. Use
// this to replace repetitive sequences of individual
// [StringFlag] registrations.
//
// Parameters:
//   - c: Cobra command to register on
//   - ptrs: Pointers to the string variables
//   - names: Flag name constants
//   - descKeys: YAML DescKeys for the flag descriptions
func BindStringFlags(
	c *cobra.Command,
	ptrs []*string,
	names, descKeys []string,
) {
	for i, p := range ptrs {
		StringFlag(c, p, names[i], descKeys[i])
	}
}

// BindBoolFlags registers multiple boolean flags that have
// no shorthand letter, each defaulting to false. All three
// slices must have the same length; each index i produces one
// [BoolFlag] call. Use this to replace repetitive sequences
// of individual [BoolFlag] registrations.
//
// Parameters:
//   - c: Cobra command to register on
//   - ptrs: Pointers to the bool variables
//   - names: Flag name constants
//   - descKeys: YAML DescKeys for the flag descriptions
func BindBoolFlags(
	c *cobra.Command,
	ptrs []*bool,
	names, descKeys []string,
) {
	for i, p := range ptrs {
		BoolFlag(c, p, names[i], descKeys[i])
	}
}

// BindBoolFlagsP registers multiple boolean flags that each
// have a shorthand letter, each defaulting to false. All four
// slices must have the same length; each index i produces one
// [BoolFlagP] call. Use this to replace repetitive sequences
// of individual [BoolFlagP] registrations.
//
// Parameters:
//   - c: Cobra command to register on
//   - ptrs: Pointers to the bool variables
//   - names: Flag name constants
//   - shorts: Shorthand letters
//   - descKeys: YAML DescKeys for the flag descriptions
func BindBoolFlagsP(
	c *cobra.Command,
	ptrs []*bool,
	names, shorts, descKeys []string,
) {
	for i, p := range ptrs {
		BoolFlagP(c, p, names[i], shorts[i], descKeys[i])
	}
}

// BindStringFlagShorts registers multiple no-pointer string
// flags that each have a shorthand letter. All three slices
// must have the same length; each index i produces one
// [StringFlagShort] call. Use when the value is retrieved
// later via cmd.Flags().GetString().
//
// Parameters:
//   - c: Cobra command to register on
//   - names: Flag name constants
//   - shorts: Shorthand letters
//   - descKeys: YAML DescKeys for the flag descriptions
func BindStringFlagShorts(
	c *cobra.Command,
	names, shorts, descKeys []string,
) {
	for i, name := range names {
		StringFlagShort(c, name, shorts[i], descKeys[i])
	}
}

// BindStringFlagsPDefault registers multiple string flags
// that each have a shorthand letter and a non-empty default
// value. All five slices must have the same length; each
// index i produces one [StringFlagPDefault] call.
//
// Parameters:
//   - c: Cobra command to register on
//   - ptrs: Pointers to the string variables
//   - names: Flag name constants
//   - shorts: Shorthand letters
//   - defaults: Default values for each flag
//   - descKeys: YAML DescKeys for the flag descriptions
func BindStringFlagsPDefault(
	c *cobra.Command,
	ptrs []*string,
	names, shorts, defaults, descKeys []string,
) {
	for i, p := range ptrs {
		StringFlagPDefault(
			c, p, names[i], shorts[i],
			defaults[i], descKeys[i],
		)
	}
}
