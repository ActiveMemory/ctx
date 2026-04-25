#!/bin/bash
# ctx session start hook for Copilot CLI
# Bootstraps context and loads the agent packet
set -euo pipefail

# Bootstrap ctx context
ctx system bootstrap 2>/dev/null || true

# Load AI-optimized context packet
ctx agent 2>/dev/null || true
