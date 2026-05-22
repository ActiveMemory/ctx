# Session Log Schema

Shape for `.context/ingest/SESSION_LOG.md`. One line per
phase boundary during a `/ctx-kb-ingest` pass. The session
log is mid-flight working memory: it records the agent's
movement through the pass so a partial closeout can name
exactly where the work stopped.

The session log is not the recall artifact. The handover
(`.context/handovers/<TS>-<slug>.md`) is the sole
authoritative source for "what happened in the last
session"; the session log is operational telemetry for the
pass currently in flight.

## Fields

The single-line format encodes every field positionally:

```
[YYYY-MM-DD HH:MM:SS sha=<short> branch=<name>] phase=<name> status=<state> note=<≤120 chars>
```

| Field      | Description                                                       |
|------------|-------------------------------------------------------------------|
| timestamp  | `YYYY-MM-DD HH:MM:SS` local time at the phase boundary.           |
| `sha`      | 7-char git short SHA of HEAD at the boundary.                     |
| `branch`   | Current symbolic ref, or `detached` if HEAD is not on a branch.   |
| `phase`    | One of `resolve`, `synthesise`, `reconcile`, `closeout`.          |
| `status`   | One of `done`, `partial`, `blocked`.                              |
| `note`     | Free text, at most 120 characters; no newlines.                   |

## Example

```
[2026-05-16 14:32:08 sha=88d5287 branch=main] phase=synthesise status=done note=widget topic page extended; cited EV-042..EV-046
```

## Rules

- Do not read `SESSION_LOG.md` on session start; it is
  mid-flight working memory, not a recall artifact. The
  handover is the recall surface.
- Append-only; never edit a past line. Corrections go on a
  new line with a `note=` explaining the prior entry.
- The four phase tokens are the closed set; introducing a
  new phase requires a `domain-decisions.md` entry and a
  spec update.
- Keep notes under the 120-character budget; longer
  observations belong in the closeout body, not here.
- Lines without `sha=` and `branch=` are malformed; doctor
  advisory flags them. Provenance is non-optional on this
  operational artifact.
