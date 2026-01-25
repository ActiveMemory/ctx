#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

#
# Release script for Context CLI
#
# This script prepares and creates a new release. It:
# 1. Builds binaries for all platforms
# 2. Generates release notes
# 3. Creates a signed git tag
#
# Usage: ./hack/release.sh
#
# =============================================================================
# RELEASE CHECKLIST - Update these for each release:
# =============================================================================
#
# Before running this script:
#
# 1. UPDATE THE VERSION in the VERSION file at the repository root
#
# 2. UPDATE DOCUMENTATION with new version:
#    - docs/index.md: Change download URLs from "latest" to "v0.1.0"
#      (lines with: releases/latest/download/ -> releases/download/v0.1.0/)
#
# 3. ENSURE all tests pass:
#    make test
#    make smoke
#
# 4. ENSURE working tree is clean:
#    git status (should show "nothing to commit")
#
# 5. COMMIT any version-related changes before running this script
#
# After running this script:
#
# 1. PUSH the tag:
#    git push origin v0.1.0
#
# 2. CREATE GitHub release:
#    - Go to https://github.com/ActiveMemory/ctx/releases/new
#    - Select the tag v0.1.0
#    - Copy release notes from dist/RELEASE_NOTES.md
#    - Upload all binaries and .sha256 files from dist/
#
# 3. UPDATE the "latest" tag (optional, for docs compatibility):
#    git tag -d latest 2>/dev/null || true
#    git push origin :refs/tags/latest 2>/dev/null || true
#    git tag latest v0.1.0
#    git push origin latest
#
# =============================================================================

set -e

# -----------------------------------------------------------------------------
# CONFIGURATION - Read from VERSION file
# -----------------------------------------------------------------------------
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

if [ ! -f "$ROOT_DIR/VERSION" ]; then
    echo "ERROR: VERSION file not found"
    exit 1
fi

VERSION="v$(cat "$ROOT_DIR/VERSION" | tr -d '[:space:]')"
# -----------------------------------------------------------------------------

# Derived values
TAG_NAME="${VERSION}"
RELEASE_NOTES="dist/RELEASE_NOTES.md"

echo "=============================================="
echo "  Context CLI Release: ${VERSION}"
echo "=============================================="
echo ""

# Check for clean working tree
if [ -n "$(git status --porcelain)" ]; then
    echo "ERROR: Working tree is not clean."
    echo "Please commit or stash your changes before releasing."
    echo ""
    git status --short
    exit 1
fi

# Check if tag already exists
if git rev-parse "${TAG_NAME}" >/dev/null 2>&1; then
    echo "ERROR: Tag ${TAG_NAME} already exists."
    echo "If you need to recreate it, delete it first:"
    echo "  git tag -d ${TAG_NAME}"
    echo "  git push origin :refs/tags/${TAG_NAME}"
    exit 1
fi

# Run tests
echo "Running tests..."
make test
echo ""

# Run smoke tests
echo "Running smoke tests..."
make smoke
echo ""

# Build binaries
echo "Building binaries for all platforms..."
./hack/build-all.sh "${VERSION#v}"  # Remove 'v' prefix for build script
echo ""

# Generate release notes
echo "Generating release notes..."
cat > "${RELEASE_NOTES}" << 'NOTES_HEADER'
# Context CLI v0.1.0

Initial release of the Context CLI (`ctx`) - a tool for persistent AI context management.

## What's New

This is the first stable release of `ctx`, providing:

- **Context Management**: Create and maintain `.context/` directories with structured markdown files
- **AI Integration**: Built-in support for Claude Code with hooks and slash commands
- **Session Persistence**: Automatic session saving and transcript management
- **Drift Detection**: Track staleness of context files
- **Multi-tool Support**: Integration guides for Claude Code, Cursor, Aider, Copilot, and Windsurf

## Installation

### Linux (x86_64)
```bash
curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.1.0/ctx-0.1.0-linux-amd64
chmod +x ctx-0.1.0-linux-amd64
sudo mv ctx-0.1.0-linux-amd64 /usr/local/bin/ctx
```

### Linux (ARM64)
```bash
curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.1.0/ctx-0.1.0-linux-arm64
chmod +x ctx-0.1.0-linux-arm64
sudo mv ctx-0.1.0-linux-arm64 /usr/local/bin/ctx
```

### macOS (Apple Silicon)
```bash
curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.1.0/ctx-0.1.0-darwin-arm64
chmod +x ctx-0.1.0-darwin-arm64
sudo mv ctx-0.1.0-darwin-arm64 /usr/local/bin/ctx
```

### macOS (Intel)
```bash
curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.1.0/ctx-0.1.0-darwin-amd64
chmod +x ctx-0.1.0-darwin-amd64
sudo mv ctx-0.1.0-darwin-amd64 /usr/local/bin/ctx
```

### Windows
Download `ctx-0.1.0-windows-amd64.exe` or `ctx-0.1.0-windows-arm64.exe` and add to your PATH.

### Verifying Checksums
Each binary has a `.sha256` file. Verify with:
```bash
curl -LO https://github.com/ActiveMemory/ctx/releases/download/v0.1.0/ctx-0.1.0-linux-amd64.sha256
sha256sum -c ctx-0.1.0-linux-amd64.sha256
```

## Quick Start

```bash
# Initialize context in your project
ctx init

# Check context status
ctx status

# Get AI-ready context packet
ctx agent --budget 4000
```

## Documentation

Full documentation available at [ctx.ist](https://ctx.ist)

## Checksums

Each binary has a corresponding `.sha256` file for verification.
NOTES_HEADER

echo "Release notes written to ${RELEASE_NOTES}"
echo ""

# Create signed tag
echo "Creating signed tag ${TAG_NAME}..."
git tag -s "${TAG_NAME}" -m "Release ${VERSION}

Context CLI ${VERSION} - Initial release

See RELEASE_NOTES.md for details."

echo ""
echo "=============================================="
echo "  Release preparation complete!"
echo "=============================================="
echo ""
echo "Created:"
echo "  - Binaries in dist/"
echo "  - Checksums (.sha256 per binary) in dist/"
echo "  - Release notes in dist/RELEASE_NOTES.md"
echo "  - Signed tag: ${TAG_NAME}"
echo ""
echo "Next steps:"
echo ""
echo "  1. Verify the tag:"
echo "     git show ${TAG_NAME}"
echo ""
echo "  2. Push the tag:"
echo "     git push origin ${TAG_NAME}"
echo ""
echo "  3. Create GitHub release at:"
echo "     https://github.com/ActiveMemory/ctx/releases/new"
echo ""
echo "  4. Upload these files to the release:"
ls -1 dist/ctx-* 2>/dev/null | sed 's/^/     /'
echo ""
echo "  5. (Optional) Update 'latest' tag:"
echo "     git tag -d latest 2>/dev/null || true"
echo "     git push origin :refs/tags/latest 2>/dev/null || true"
echo "     git tag latest ${TAG_NAME}"
echo "     git push origin latest"
echo ""
