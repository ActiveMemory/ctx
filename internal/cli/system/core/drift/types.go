//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

// MarketplaceManifest is the structure of .claude-plugin/marketplace.json.
type MarketplaceManifest struct {
	Plugins []struct {
		Version string `json:"version"`
	} `json:"plugins"`
}
