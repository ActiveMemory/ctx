//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// SplitPopular partitions a slice into popular and long-tail groups.
//
// Parameters:
//   - items: Slice of items implementing PopularSplittable
//
// Returns:
//   - popular: Items where IsPopular() is true
//   - longtail: Items where IsPopular() is false
func SplitPopular[T PopularSplittable](items []T) (popular, longtail []T) {
	for _, item := range items {
		if item.IsPopular() {
			popular = append(popular, item)
		} else {
			longtail = append(longtail, item)
		}
	}
	return popular, longtail
}

// writeSection writes a headed list section to sb if items is non-empty.
//
// Parameters:
//   - sb: String builder to write to
//   - headingKey: YAML DescKey for the section heading
//   - items: Slice of items to render
//   - formatFn: Function that renders one item as a string
func writeSection[T any](sb *strings.Builder, headingKey string, items []T, formatFn func(T) string) {
	if len(items) == 0 {
		return
	}
	nl := token.NewlineLF
	sb.WriteString(desc.TextDesc(headingKey) + nl + nl)
	for _, item := range items {
		sb.WriteString(formatFn(item))
	}
	sb.WriteString(nl)
}
