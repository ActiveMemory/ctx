#!/bin/bash
# ctx session end hook for Copilot CLI
# Checks for unsaved context and records heartbeat
set -euo pipefail

ctx system check-context-size 2>/dev/null || true
ctx system check-persistence 2>/dev/null || true
ctx system check-journal 2>/dev/null || true
ctx system check-version 2>/dev/null || true
ctx system heartbeat 2>/dev/null || true
