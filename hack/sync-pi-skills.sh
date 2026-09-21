#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# sync-pi-skills.sh — sync Pi skills from the generated OpenCode tree.
#
# specs/pi-cli-integration.md defines the Pi skill set as "the same
# bundled skill set as OpenCode", and TestPiSkillsMirrorOpenCode
# freezes the two trees byte-for-byte. This script makes that
# derivation structural: the Pi tree is a verbatim copy of the
# generated OpenCode tree (SKILL.md files and their references/
# directories), so the freeze can never be broken by a one-sided
# edit.
#
# Run order matters: hack/sync-opencode-skills.sh regenerates the
# OpenCode tree from the canonical ctx skills first; `make build`
# and `make audit` sequence them correctly.

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

OPENCODE_SKILLS="internal/assets/integrations/opencode/skills"
PI_SKILLS="internal/assets/integrations/pi/skills"

if [ ! -d "$OPENCODE_SKILLS" ]; then
  echo "FAIL: $OPENCODE_SKILLS not found — run sync-opencode-skills first" >&2
  exit 1
fi

rm -rf "$PI_SKILLS"
mkdir -p "$(dirname "$PI_SKILLS")"
cp -R "$OPENCODE_SKILLS" "$PI_SKILLS"

count=$(find "$PI_SKILLS" -name SKILL.md | wc -l | tr -d ' ')
echo "Pi skills synced: $count mirrored from the OpenCode tree."
