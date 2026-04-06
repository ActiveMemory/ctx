//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package flagbind

import (
	"time"

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

// StringFlagDefault registers a string flag with no shorthand
// and a non-empty default value.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the string variable
//   - name: Flag name constant
//   - defaultVal: Default value for the flag
//   - descKey: YAML DescKey for the flag description
func StringFlagDefault(
	c *cobra.Command, p *string,
	name, defaultVal, descKey string,
) {
	c.Flags().StringVar(
		p, name, defaultVal, desc.Flag(descKey),
	)
}

// BoolFlagDefault registers a boolean flag with no
// shorthand and a non-false default value.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the bool variable
//   - name: Flag name constant
//   - defaultVal: Default value for the flag
//   - descKey: YAML DescKey for the flag description
func BoolFlagDefault(
	c *cobra.Command, p *bool,
	name string, defaultVal bool, descKey string,
) {
	c.Flags().BoolVar(
		p, name, defaultVal, desc.Flag(descKey),
	)
}

// IntFlag registers an integer flag with no shorthand.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the int variable
//   - name: Flag name constant
//   - defaultVal: Default value for the flag
//   - descKey: YAML DescKey for the flag description
func IntFlag(
	c *cobra.Command, p *int,
	name string, defaultVal int, descKey string,
) {
	c.Flags().IntVar(
		p, name, defaultVal, desc.Flag(descKey),
	)
}

// DurationFlag registers a duration flag with no shorthand.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the time.Duration variable
//   - name: Flag name constant
//   - defaultVal: Default value for the flag
//   - descKey: YAML DescKey for the flag description
func DurationFlag(
	c *cobra.Command, p *time.Duration,
	name string, defaultVal time.Duration,
	descKey string,
) {
	c.Flags().DurationVar(
		p, name, defaultVal, desc.Flag(descKey),
	)
}

// PersistentBoolFlag registers a persistent boolean
// flag with no shorthand, defaulting to false.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the bool variable
//   - name: Flag name constant
//   - descKey: YAML DescKey for the flag description
func PersistentBoolFlag(
	c *cobra.Command, p *bool, name, descKey string,
) {
	c.PersistentFlags().BoolVar(
		p, name, false, desc.Flag(descKey),
	)
}

// PersistentStringFlag registers a persistent string
// flag with no shorthand, defaulting to empty string.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the string variable
//   - name: Flag name constant
//   - descKey: YAML DescKey for the flag description
func PersistentStringFlag(
	c *cobra.Command, p *string, name, descKey string,
) {
	c.PersistentFlags().StringVar(
		p, name, "", desc.Flag(descKey),
	)
}

// StringArrayFlagP registers a string array flag with a shorthand
// letter. The flag can be repeated: --tag x --tag y.
//
// Parameters:
//   - c: Cobra command to register on
//   - p: Pointer to the string slice variable
//   - name: Flag name constant
//   - short: Shorthand letter
//   - descKey: YAML DescKey for the flag description
func StringArrayFlagP(
	c *cobra.Command, p *[]string, name, short, descKey string,
) {
	c.Flags().StringArrayVarP(p, name, short, nil, desc.Flag(descKey))
}

// BoolFlagNoPtr registers a boolean flag with no
// shorthand and no pointer, defaulting to false.
// Use when the value is retrieved via
// cmd.Flags().GetBool().
//
// Parameters:
//   - c: Cobra command to register on
//   - name: Flag name constant
//   - descKey: YAML DescKey for the flag description
func BoolFlagNoPtr(
	c *cobra.Command, name, descKey string,
) {
	c.Flags().Bool(
		name, false, desc.Flag(descKey),
	)
}

// BoolFlagShort registers a boolean flag with a
// shorthand, returning a pointer. Use when the value
// is retrieved via cmd.Flags().GetBool().
//
// Parameters:
//   - c: Cobra command to register on
//   - name: Flag name constant
//   - short: Shorthand letter
//   - descKey: YAML DescKey for the flag description
func BoolFlagShort(
	c *cobra.Command, name, short, descKey string,
) {
	c.Flags().BoolP(
		name, short, false, desc.Flag(descKey),
	)
}

// StringFlagShort registers a string flag with a
// shorthand, returning a pointer. Use when the value
// is retrieved via cmd.Flags().GetString().
//
// Parameters:
//   - c: Cobra command to register on
//   - name: Flag name constant
//   - short: Shorthand letter
//   - descKey: YAML DescKey for the flag description
func StringFlagShort(
	c *cobra.Command, name, short, descKey string,
) {
	c.Flags().StringP(
		name, short, "", desc.Flag(descKey),
	)
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
