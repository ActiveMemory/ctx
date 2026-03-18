//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/spf13/cobra"
)

// InfoNoProject reports that no supported project was detected.
//
// Parameters:
//   - cmd: Cobra command for output
//   - builderNames: Comma-separated list of supported project types
func InfoNoProject(cmd *cobra.Command, builderNames string) {
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteDepsNoProject))
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteDepsLookingFor))
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteDepsUseType), builderNames))
}

// InfoNoDeps reports that no dependencies were found.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoNoDeps(cmd *cobra.Command) {
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteDepsNoDeps))
}
