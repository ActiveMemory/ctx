//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

// Shell script template used by `ctx trigger add` to
// scaffold new lifecycle triggers.
const (
	// TriggerScript is the bash template written to
	// .context/hooks/<type>/<name>.sh by ctx trigger add.
	//
	// Args (in order):
	//   - name: trigger script base name (without .sh)
	//   - type: trigger type (e.g. pre-tool-use, session-start)
	//
	// The generated script has no executable bit; users
	// must run `ctx trigger enable <name>` after review, so
	// unreviewed code never fires on real events.
	TriggerScript = `#!/usr/bin/env bash
# Trigger: %s
# Type:    %s
# Created by: ctx trigger add
#
# Enable with: ctx trigger enable %[1]s
# Test with:   ctx trigger test %[2]s

set -euo pipefail

# Read the JSON event payload from stdin.
INPUT=$(cat)

# Parse the fields you need from the payload.
TRIGGER_TYPE=$(echo "$INPUT" | jq -r '.hookType // empty')
TOOL=$(echo "$INPUT" | jq -r '.tool // empty')
PATH_ARG=$(echo "$INPUT" | jq -r '.path // empty')

# Your trigger logic here.

# Return a JSON response on stdout. "cancel": true blocks
# the tool call (pre-tool-use only); "context" injects
# additional context; "message" is shown to the user.
echo '{"cancel": false, "context": "", "message": ""}'
`
)
