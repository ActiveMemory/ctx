//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// NoConfig returns an error when the zensical config file is missing.
//
// Parameters:
//   - dir: directory where the config was expected.
//
// Returns:
//   - error: "no zensical.toml found in <dir>"
func NoConfig(dir string) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrSiteNoSiteConfig), dir)
}

// MarshalFeed wraps a failure to marshal an Atom feed.
//
// Parameters:
//   - cause: the underlying marshal error
//
// Returns:
//   - error: "cannot marshal feed: <cause>"
func MarshalFeed(cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSiteMarshalFeed), cause,
	)
}

// ZensicalNotFound returns an error when zensical is not installed.
//
// Returns:
//   - error: includes installation instructions
func ZensicalNotFound() error {
	return errors.New(desc.Text(text.DescKeyErrSiteZensicalNotFound))
}
