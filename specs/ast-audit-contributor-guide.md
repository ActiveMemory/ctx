---
title: AST audit contributor guide
date: 2026-04-03
status: ready
---

# AST Audit Contributor Guide

## Problem

Contributors (human and AI) routinely introduce code that violates
project conventions enforced by `internal/audit/` AST tests. The
violations are caught at `go test` time, but the fix patterns are
not documented — contributors must reverse-engineer the convention
from the test name and error message.

## Solution

A contributor-facing document in `docs/reference/` that catalogs
every common violation pattern, shows a before/after code example,
and explains the rationale. Organized by convention category, not
by test name (contributors think "I have a magic string" not
"TestNoMagicStrings is failing").

## Scope

- Document only patterns enforced by `internal/audit/` tests
- Before/after examples drawn from actual commits
- Link each pattern to its test and the CONVENTIONS.md entry
- No changes to tests or source code
