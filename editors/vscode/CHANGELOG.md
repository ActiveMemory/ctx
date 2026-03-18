# Changelog

All notable changes to the **ctx — Persistent Context for AI** extension
will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [0.9.0] — 2026-03-19

### Added

- **@ctx chat participant** with 45 slash commands covering context
  lifecycle, task management, session recall, and discovery
- **Natural language routing** — type plain English after `@ctx` and
  the extension maps it to the correct handler
- **Auto-bootstrap** — downloads the ctx CLI binary if not found on PATH
- **Detection ring** — terminal command watcher and file edit watcher
  record governance violations for the MCP engine
- **Status bar reminders** — `$(bell) ctx` indicator for pending reminders
- **Automatic hooks** — file save, git commit, dependency change, and
  context file change handlers
- **Follow-up suggestions** — context-aware buttons after each command
- **`/diag` command** — diagnose extension issues with step-by-step timing

### Configuration

- `ctx.executablePath` — path to the ctx CLI binary (default: `ctx`)

## [Unreleased]

- Marketplace publication
