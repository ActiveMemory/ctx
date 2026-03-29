//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package register provides helpers for cobra flag registration.
package flagbind

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

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
