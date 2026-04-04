//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

// Shell script template for new hook scripts created by
// ctx hook add.
const (
	// TriggerScript is the bash hook script template.
	// Args: name, hookType.
	TriggerScript = `#!/usr/bin/env bash
# Hook: %s
# Type: %s
# Created by: ctx hook add

set -euo pipefail

INPUT=$(cat)

# Parse input fields
HOOK_TYPE=$(echo "$INPUT" | jq -r '.hookType')
TOOL=$(echo "$INPUT" | jq -r '.tool // empty')

# Your hook logic here

# Return output
echo '{"cancel": false, "context": "", "message": ""}'
`
)
