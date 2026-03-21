//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

import (
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// Init loads all embedded YAML description maps. Call once from main()
// before building the command tree. Tests that need descriptions must
// call Init() in their setup.
func Init() {
	CommandsMap = loadYAML(asset.PathCommandsYAML)
	FlagsMap = loadYAML(asset.PathFlagsYAML)
	TextMap = loadYAMLDir(asset.DirCommandsText)
	ExamplesMap = loadYAML(asset.PathExamplesYAML)
	allowPerms = loadPermissions(asset.PathAllowTxt)
	denyPerms = loadPermissions(asset.PathDenyTxt)
	stopWordsMap = loadStopWords()
}
