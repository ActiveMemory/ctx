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
# 1. Verifies release notes exist (generate with /release-notes in Claude Code)
# 2. Builds binaries for all platforms
# 3. Creates and pushes a signed git tag
# 4. Updates the "latest" tag
#
# Usage: ./hack/release.sh
#
# =============================================================================
# RELEASE CHECKLIST - Before running this script:
# =============================================================================
#
# 1. UPDATE THE VERSION in the VERSION file at the repository root
#
# 2. GENERATE RELEASE NOTES using Claude Code:
#    /release-notes
#
# 3. UPDATE DOCUMENTATION with new version:
#    - docs/index.md: Update download URLs to new version
#
# 4. COMMIT all version-related changes
#
# 5. ENSURE working tree is clean:
#    git status (should show "nothing to commit")
#
# After running this script:
#
# 1. CREATE GitHub release at:
#    https://github.com/ActiveMemory/ctx/releases/new
#    - Select the pushed tag
#    - Copy release notes from dist/RELEASE_NOTES.md
#    - Upload all binaries and .sha256 files from dist/
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

VERSION="v$(tr -d '[:space:]' < "$ROOT_DIR/VERSION")"
# -----------------------------------------------------------------------------

# Derived values
TAG_NAME="${VERSION}"
RELEASE_NOTES="dist/RELEASE_NOTES.md"

echo "=============================================="
echo "  Context CLI Release: ${VERSION}"
echo "=============================================="
echo ""

# Check for release notes first
if [ ! -f "${RELEASE_NOTES}" ]; then
    echo "ERROR: ${RELEASE_NOTES} not found."
    echo ""
    echo "Generate release notes first using Claude Code:"
    echo "  /release-notes"
    echo ""
    exit 1
fi
echo "Found ${RELEASE_NOTES}"
echo ""

# Check for clean working tree (before we make changes)
if [ -n "$(git status --porcelain)" ]; then
    echo "ERROR: Working tree is not clean."
    echo "Please commit or stash your changes before releasing."
    echo ""
    git status --short
    exit 1
fi

# Update version references in documentation
echo "Updating version references in docs/index.md..."
VERSION_NUM="${VERSION#v}"  # Remove 'v' prefix
sed -i.bak -E "s|/v[0-9]+\.[0-9]+\.[0-9]+/ctx-[0-9]+\.[0-9]+\.[0-9]+|/v${VERSION_NUM}/ctx-${VERSION_NUM}|g" docs/index.md
sed -i.bak -E "s|ctx-[0-9]+\.[0-9]+\.[0-9]+-linux-amd64|ctx-${VERSION_NUM}-linux-amd64|g" docs/index.md
sed -i.bak -E "s|ctx-[0-9]+\.[0-9]+\.[0-9]+-linux-arm64|ctx-${VERSION_NUM}-linux-arm64|g" docs/index.md
sed -i.bak -E "s|ctx-[0-9]+\.[0-9]+\.[0-9]+-darwin-amd64|ctx-${VERSION_NUM}-darwin-amd64|g" docs/index.md
sed -i.bak -E "s|ctx-[0-9]+\.[0-9]+\.[0-9]+-darwin-arm64|ctx-${VERSION_NUM}-darwin-arm64|g" docs/index.md
sed -i.bak -E "s|ctx-[0-9]+\.[0-9]+\.[0-9]+-windows-amd64|ctx-${VERSION_NUM}-windows-amd64|g" docs/index.md
rm -f docs/index.md.bak

# Update versions.md with new release
echo "Updating docs/versions.md..."
RELEASE_DATE=$(date +%Y-%m-%d)
# Add new version row after the table header
sed -i.bak "/^| Version | Release Date | Documentation |$/,/^|/{
    /^|-/a\\
| ${VERSION} | ${RELEASE_DATE} | [View docs](https://github.com/ActiveMemory/ctx/tree/${VERSION}/docs) |
}" docs/versions.md
# Update the "latest stable" reference
sed -i.bak -E "s|see \[v[0-9]+\.[0-9]+\.[0-9]+\]|see [${VERSION}]|g" docs/versions.md
sed -i.bak -E "s|/tree/v[0-9]+\.[0-9]+\.[0-9]+/docs\)\.$|/tree/${VERSION}/docs).|g" docs/versions.md
rm -f docs/versions.md.bak

# Rebuild site with updated docs
echo "Rebuilding documentation site..."
make site

# Commit docs and site updates
echo "Committing documentation updates..."
git add docs/index.md docs/versions.md site/
git commit -m "docs: update download links and versions page for ${VERSION}"
echo ""

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

# Create signed tag
echo "Creating signed tag ${TAG_NAME}..."
git tag -s "${TAG_NAME}" -m "Release ${VERSION}

Context CLI ${VERSION}

See dist/RELEASE_NOTES.md for details."
echo ""

# Push the version tag
echo "Pushing tag ${TAG_NAME} to origin..."
git push origin "${TAG_NAME}"
echo ""

# Update the "latest" tag
echo "Updating 'latest' tag..."
git tag -d latest 2>/dev/null || true
git push origin :refs/tags/latest 2>/dev/null || true
git tag latest "${TAG_NAME}"
git push origin latest
echo ""

echo "=============================================="
echo "  Release ${VERSION} complete!"
echo "=============================================="
echo ""
echo "Created and pushed:"
echo "  - Tag: ${TAG_NAME}"
echo "  - Tag: latest -> ${TAG_NAME}"
echo ""
echo "Built artifacts in dist/:"
find dist -maxdepth 1 -name 'ctx-*' 2>/dev/null | sed 's/^/  /'
echo ""
echo "Next step:"
echo ""
echo "  Create GitHub release at:"
echo "  https://github.com/ActiveMemory/ctx/releases/new?tag=${TAG_NAME}"
echo ""
echo "  - Paste contents of dist/RELEASE_NOTES.md"
echo "  - Upload all files from dist/"
echo ""
