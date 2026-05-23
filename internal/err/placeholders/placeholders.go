//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package placeholders

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// ReadLocale wraps a failure to read the embedded
// placeholder locale YAML.
//
// Parameters:
//   - locale: the locale identifier (e.g. "en").
//   - cause: the underlying read error.
//
// Returns:
//   - error: "read placeholder locale <locale>: <cause>".
func ReadLocale(locale string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPlaceholdersReadLocale), locale, cause,
	)
}

// ParseLocale wraps a failure to YAML-parse the embedded
// placeholder locale file.
//
// Parameters:
//   - locale: the locale identifier (e.g. "en").
//   - cause: the underlying YAML parse error.
//
// Returns:
//   - error: "parse placeholder locale <locale>: <cause>".
func ParseLocale(locale string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrPlaceholdersParseLocale), locale, cause,
	)
}
