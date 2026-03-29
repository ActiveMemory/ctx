//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package moc

import "github.com/ActiveMemory/ctx/internal/entity"

// SplitPopular partitions a slice into popular and long-tail groups.
//
// Parameters:
//   - items: Slice of items implementing PopularSplittable
//
// Returns:
//   - popular: Items where IsPopular() is true
//   - longtail: Items where IsPopular() is false
func SplitPopular[T entity.PopularSplittable](
	items []T,
) (popular, longtail []T) {
	for _, item := range items {
		if item.IsPopular() {
			popular = append(popular, item)
		} else {
			longtail = append(longtail, item)
		}
	}
	return popular, longtail
}
