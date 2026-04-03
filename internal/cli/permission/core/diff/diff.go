//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package diff

// StringSlices computes the set difference between golden
// and local slices.
//
// Parameters:
//   - golden: Entries from the golden settings file
//   - local: Entries from the local settings file
//
// Returns:
//   - restored: Entries in golden but not in local
//   - dropped: Entries in local but not in golden
//
// Both output slices preserve the source ordering of their
// respective inputs.
func StringSlices(
	golden, local []string,
) (restored, dropped []string) {
	goldenSet := make(map[string]struct{}, len(golden))
	for _, s := range golden {
		goldenSet[s] = struct{}{}
	}

	localSet := make(map[string]struct{}, len(local))
	for _, s := range local {
		localSet[s] = struct{}{}
	}

	for _, s := range golden {
		if _, ok := localSet[s]; !ok {
			restored = append(restored, s)
		}
	}

	for _, s := range local {
		if _, ok := goldenSet[s]; !ok {
			dropped = append(dropped, s)
		}
	}

	return restored, dropped
}
