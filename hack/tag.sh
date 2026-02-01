#!/bin/bash

#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

#
# Tag script for Context CLI
#
# Creates a signed git tag and pushes it to origin.
# Reads version from VERSION file.
#
# Usage: ./hack/tag.sh
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Read version from VERSION file
if [ ! -f "$ROOT_DIR/VERSION" ]; then
    echo "ERROR: VERSION file not found"
    exit 1
fi

VERSION="v$(tr -d '[:space:]' < "$ROOT_DIR/VERSION")"

echo "Creating tag: $VERSION"

# Check if tag already exists locally
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo "ERROR: Tag $VERSION already exists locally."
    echo "To recreate it:"
    echo "  git tag -d $VERSION"
    exit 1
fi

# Check if tag exists on remote
if git ls-remote --tags origin | grep -q "refs/tags/$VERSION$"; then
    echo "ERROR: Tag $VERSION already exists on origin."
    echo "To recreate it:"
    echo "  git push origin :refs/tags/$VERSION"
    exit 1
fi

# Create signed tag
git tag -s "$VERSION" -m "$VERSION"

echo ""
echo "Tag $VERSION created locally."
echo ""
echo "To push:"
echo "  git push origin --tags"
echo ""
echo "Or to push just this tag:"
echo "  git push origin $VERSION"
