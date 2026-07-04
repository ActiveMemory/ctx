# Spec: ctx is optional in deployed agent instructions

## Problem

The agent-instruction templates ctx deploys (CLAUDE.md, Copilot
INSTRUCTIONS.md) treated a missing `ctx` binary as a fatal error:
"CRITICAL, not optional", relay the error, STOP. That posture is
wrong for the audience those files actually meet. A contributor
who clones a ctx-using project without installing ctx gets an
agent that refuses to help with ordinary work, which punishes
exactly the newcomers the project wants.

## Decision

`ctx` is an optional companion in deployed instructions: helpful,
never required to build, test, or contribute.

- **Binary not found**: not an error. The agent mentions once how
  to enable persistent context, then proceeds with the user's
  task without ctx. The Claude variant recommends installing the
  plugin from a local clone (releases are infrequent, so
  marketplace versions may lag the repository); the Copilot
  variant points at the getting-started guide.
- **Installed but erroring**: unchanged posture. Relay the error
  verbatim, point at getting-started, STOP; recovery actions are
  the user's decision, not the agent's.

Both templates carry the same policy in their own idiom (the
Copilot kit mirrors the Claude template per the parity kit).

## Non-Goals

- Changing this repository's own CLAUDE.md (contributors to ctx
  itself are expected to run ctx)
- Auto-installing anything on the user's behalf

## Acceptance

- [ ] Both templates distinguish "not found" (mention once,
      proceed) from "installed but erroring" (relay, STOP)
- [ ] Neither template calls ctx "CRITICAL, not optional"
- [ ] Embedded-asset tests pass with the updated templates
