#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#

# Register the gitnexus MCP server with Claude Code.

set -euo pipefail

# Check if gitnexus is installed.
if ! command -v gitnexus &>/dev/null; then
    echo "gitnexus is not installed."
    echo ""
    echo "Setup steps:"
    echo "  1. Install gitnexus:"
    echo ""
    echo "     npm install -g gitnexus"
    echo ""
    echo "  2. Run this again:"
    echo ""
    echo "     make register-mcp"
    exit 1
fi

# Check if already registered.
if claude mcp list 2>/dev/null | grep -q "gitnexus"; then
    echo "gitnexus MCP server is already registered."
    claude mcp list 2>/dev/null | grep "gitnexus"
    exit 0
fi

claude mcp add -s user \
    -- gitnexus bash -lc 'gitnexus mcp'

echo "gitnexus MCP server registered."
echo "Restart Claude Code to pick it up."
