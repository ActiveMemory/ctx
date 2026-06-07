# Run the Dream

The **dream** is a scheduled, out-of-band pass that triages your
gitignored `ideas/` folder — classifying each idea against your codebase
and specs, and emitting gated **proposals** (archive / merge / promote /
mark-blog / keep) for you to review. It only ever proposes; it never
writes canonical memory and never acts on a proposal. You review the
proposals in a ~15-minute "garden walk" and accept / reject / amend.

The dream is **opt-in and off by default.** Nothing runs until you turn
it on. This recipe wires it up for Claude Code (the reference executor).
To run it under a different harness, see
[the executor contract](../reference/dream-executor-contract.md).

## Prerequisites

- A ctx project (a git working tree with `.context/`).
- An `ideas/` folder at the project root (gitignored).
- The `ctx-dream` and `ctx-serendipity` skills installed (shipped with
  `ctx setup`).
- A non-interactive Claude Code credential (cron has no interactive
  fallback).

## 1. Enable it in `.ctxrc`

Add a `dream:` section. `enabled: false` is the default — set it true:

```yaml
dream:
  enabled: true
  mode: discipline      # the only mode in v1
  max: 50               # max ideas processed per pass
  quiet_minutes: 60     # skip a pass if you were active within the window
  cadence: "30 2 * * *" # the cron schedule you'll install below
  budget: 40            # step/token ceiling per pass
  model: null           # null = the session default model
  executor: ""          # empty = the claude -p reference executor
```

## 2. Confirm `dreams/` is gitignored

The dream writes its notebook (proposals, per-source state, ledger,
backups) to a root-level `dreams/` directory. It inherits `ideas/`'s
privacy class, so it **must** stay gitignored — `ctx init` adds the
entry, and the don't-leak guard refuses any write that resolves to a
tracked path. Verify:

```bash
git check-ignore dreams && echo "ok: dreams/ is ignored"
```

## 3. Wire the guard hook

A headless pass runs with a PreToolUse guard so the agent can only write
under `dreams/`. Point a dream-specific settings file at the bundled
`guard.sh` (do **not** add it to your project's default settings — the
dream is opt-in):

```json
{
  "hooks": {
    "PreToolUse": [
      { "matcher": "Write|Edit|MultiEdit",
        "hooks": [{ "type": "command",
          "command": "<skills>/ctx-dream/guard.sh" }] },
      { "matcher": "Bash",
        "hooks": [{ "type": "command",
          "command": "<skills>/ctx-dream/guard.sh" }] }
    ]
  }
}
```

## 4. Install the cron entry

Run one pass nightly. `ctx dream` does the gate (skips when there's no
new idea delta or you were recently active), takes a lock, and invokes
the executor:

```cron
30 2 * * *  cd /path/to/project && PATH=/usr/local/bin:$PATH ctx dream >> ~/.ctx/dream.cron.log 2>&1
```

!!! warning "cron's PATH is minimal"
    cron will not see a node/nvm-managed `claude` or even `ctx` unless
    you set `PATH` in the entry (as above). If the executor binary is not
    found, `ctx dream` fails **loud** and writes `dreams/.failed` — it
    never silently no-ops.

## 5. Review what it found

The dream nags you (via `ctx remind`) when a round is waiting. Walk the
garden:

```
/ctx-serendipity
```

Each proposal shows its summary, evidence, and a one-line rationale.
Accept / reject / amend / skip — no pressure to clear the set.
Mechanical dispositions apply instantly; `merge`/`promote` are done from
the full source. Rejections are recorded so they don't re-surface.

You can also drive it directly:

```
ctx dream review
ctx dream accept <id>
ctx dream reject <id>
ctx dream amend <id> --action keep
```

## What it will never do

- Write the five canonical files (DECISIONS / LEARNINGS / CONVENTIONS /
  CONSTITUTION / TASKS). Ever.
- Act on a proposal without you. Every disposition into a tracked
  artifact passes through the human gate.
- Write anything outside `dreams/` during a pass (the guard enforces it),
  except your deliberate `promote` of an idea into `specs/`.
