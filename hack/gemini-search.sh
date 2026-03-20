#!/usr/bin/env bash
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# Register the gemini-search MCP server with Claude Code.

set -euo pipefail

if [ -z "${GEMINI_API_KEY:-}" ]; then
    echo "GEMINI_API_KEY is not set."
    echo ""
    echo "Setup steps:"
    echo "  1. Get a Gemini API key from https://aistudio.google.com/apikey"
    echo "  2. Add to your shell profile (~/.bashrc or ~/.zshrc):"
    echo ""
    echo "     export GEMINI_API_KEY=\"your-key-here\""
    echo ""
    echo "  3. Reload your shell and run this again:"
    echo ""
    echo "     make gemini-search"
    exit 1
fi

# Check if already registered.
if claude mcp list 2>/dev/null | grep -q "gemini-search"; then
    echo "gemini-search MCP server is already registered."
    claude mcp list 2>/dev/null | grep "gemini-search"
    exit 0
fi

claude mcp add -s user \
    -e "GEMINI_API_KEY=\${GEMINI_API_KEY}" \
    -- gemini-search bash -lc 'npx -y gemini-grounding'

echo "gemini-search MCP server registered."
echo "Restart Claude Code to pick it up."
