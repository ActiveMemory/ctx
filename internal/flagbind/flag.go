//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package flagbind

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// BoolFlag registers a boolean flag with no shorthand, defaulting to false.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the bool variable
//   - name: Flag name constant
//   - descKey: YAML DescKey for the flag description
func BoolFlag(c *cobra.Command, p *bool, name, descKey string) {
	c.Flags().BoolVar(p, name, false, desc.Flag(descKey))
}

// BoolFlagP registers a boolean flag with a shorthand letter, defaulting
// to false.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the bool variable
//   - name: Flag name constant
//   - short: Shorthand letter
//   - descKey: YAML DescKey for the flag description
func BoolFlagP(c *cobra.Command, p *bool, name, short, descKey string) {
	c.Flags().BoolVarP(p, name, short, false, desc.Flag(descKey))
}

// IntFlagP registers an integer flag with a shorthand letter.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the int variable
//   - name: Flag name constant
//   - short: Shorthand letter
//   - defaultVal: Default value for the flag
//   - descKey: YAML DescKey for the flag description
func IntFlagP(
	c *cobra.Command, p *int, name, short string, defaultVal int, descKey string,
) {
	c.Flags().IntVarP(p, name, short, defaultVal, desc.Flag(descKey))
}

// StringFlag registers a string flag with no shorthand.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the string variable
//   - name: Flag name constant
//   - descKey: YAML DescKey for the flag description
func StringFlag(c *cobra.Command, p *string, name, descKey string) {
	c.Flags().StringVar(p, name, "", desc.Flag(descKey))
}

// StringFlagP registers a string flag with a shorthand letter.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the string variable
//   - name: Flag name constant
//   - short: Shorthand letter
//   - descKey: YAML DescKey for the flag description
func StringFlagP(c *cobra.Command, p *string, name, short, descKey string) {
	c.Flags().StringVarP(p, name, short, "", desc.Flag(descKey))
}

// StringFlagPDefault registers a string flag with a shorthand letter and
// a non-empty default value.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the string variable
//   - name: Flag name constant
//   - short: Shorthand letter
//   - defaultVal: Default value for the flag
//   - descKey: YAML DescKey for the flag description
func StringFlagPDefault(
	c *cobra.Command, p *string, name, short, defaultVal, descKey string,
) {
	c.Flags().StringVarP(p, name, short, defaultVal, desc.Flag(descKey))
}

// LastJSON registers the --last (int) and --json (bool) flag pair used by
// list-style commands.
//
// Parameters:
//   - c: Cobra command to register on
//   - lastDefault: Default value for --last
//   - lastDescKey: YAML DescKey for the --last flag description
//   - jsonDescKey: YAML DescKey for the --json flag description
func LastJSON(
	c *cobra.Command,
	lastDefault int,
	lastDescKey, jsonDescKey string,
) {
	c.Flags().IntP(
		cFlag.Last, cFlag.ShortLast,
		lastDefault, desc.Flag(lastDescKey),
	)
	c.Flags().BoolP(cFlag.JSON, cFlag.ShortJSON, false, desc.Flag(jsonDescKey))
}
