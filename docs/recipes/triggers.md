---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Authoring Lifecycle Triggers
icon: lucide/zap
---

![ctx](../images/ctx-banner.png)

# Authoring Lifecycle Triggers

Triggers are **executable shell scripts** that fire at
specific events during an AI session. They're how you express
"when the AI saves a file, also do X" or "before the AI edits
this path, check Y first." This recipe walks through writing
your first trigger, testing it, and enabling it safely.

!!! danger "Triggers execute arbitrary code"
    A trigger is a shell script with the executable bit set.
    It runs with the same privileges as your AI tool and
    receives JSON input on stdin. **Treat triggers like
    pre-commit hooks**:

    - Only enable scripts you have read and understand.
    - Never enable a trigger you downloaded from the internet
      without reviewing every line.
    - Avoid shelling out to user-controlled values (`jq -r`
      output, `path` field, `tool` field) without quoting.
    - A malicious or buggy trigger can block tool calls,
      corrupt context files, or exfiltrate data.

    The generated trigger template starts **disabled** (no
    executable bit) so you cannot accidentally run an unreviewed
    script. Enable it explicitly with `ctx trigger enable`.

## Scenario

You want a `pre-tool-use` trigger that blocks the AI from
editing anything in `internal/crypto/` without explicit
confirmation. Cryptographic code is sensitive, and accidental
edits have caused outages before — you want a hard gate.

## Step 1 — scaffold the script

```bash
ctx trigger add pre-tool-use protect-crypto
```

That creates `.context/hooks/pre-tool-use/protect-crypto.sh`
with a template:

```bash
#!/usr/bin/env bash
set -euo pipefail

# Read the JSON event from stdin.
payload=$(cat)

# Parse fields with jq.
tool=$(echo "$payload" | jq -r '.tool // empty')
path=$(echo "$payload" | jq -r '.path // empty')

# Your logic here.

# Return a JSON result. action can be "allow", "block", or absent.
echo '{"action": "allow"}'
```

Note: the directory is `.context/hooks/pre-tool-use/` — the
on-disk layout still uses `hooks/` even though the command is
`ctx trigger`. If you `ls .context/hooks/`, that's where
your triggers live.

## Step 2 — write the logic

Open the file and replace the template body:

```bash
#!/usr/bin/env bash
set -euo pipefail

payload=$(cat)
tool=$(echo "$payload" | jq -r '.tool // empty')
path=$(echo "$payload" | jq -r '.path // empty')

# Only gate write-family tools.
case "$tool" in
  write_file|edit_file|apply_patch) ;;
  *)
    echo '{"action": "allow"}'
    exit 0
    ;;
esac

# Block any path under internal/crypto/.
case "$path" in
  internal/crypto/*|*/internal/crypto/*)
    jq -n --arg p "$path" '{
      action: "block",
      message: ("Edits to " + $p + " require manual review. " +
                "See CONVENTIONS.md for the crypto-change process.")
    }'
    exit 0
    ;;
esac

echo '{"action": "allow"}'
```

A few things to note:

- **`set -euo pipefail`** — any unhandled error aborts the
  script. Critical for a security-relevant trigger.
- **Quote everything from `jq`** — the `path` field comes from
  the AI tool; treat it as untrusted input.
- **Explicit `allow` case** — the default is allow. An
  empty or missing response is a risky default.
- **Use `jq -n --arg`** for output construction — safer than
  string concatenation when the message may contain special
  characters.

## Step 3 — test with a mock payload

Before enabling the trigger, test it with a realistic mock
input using `ctx trigger test`. This runs the script against
a synthetic JSON payload without actually firing any AI tool.

```bash
# Test the "should block" case
ctx trigger test pre-tool-use --tool write_file --path internal/crypto/aes.go
```

Expected: the trigger returns `{"action":"block", "message": "..."}`.

```bash
# Test the "should allow" case
ctx trigger test pre-tool-use --tool write_file --path internal/memory/mirror.go
```

Expected: the trigger returns `{"action":"allow"}`.

```bash
# Test that non-write tools pass through
ctx trigger test pre-tool-use --tool read_file --path internal/crypto/aes.go
```

Expected: `{"action":"allow"}` because the `case` statement
only gates write-family tools.

If any of these cases misbehave, **fix the trigger before
enabling it.** The trigger is disabled at this point, so
misbehavior doesn't affect real AI sessions.

## Step 4 — enable it

Once the test cases pass, enable the trigger:

```bash
ctx trigger enable protect-crypto
```

That sets the executable bit. Next time the AI starts a
`pre-tool-use` event, the trigger will fire.

Verify it's enabled:

```bash
ctx trigger list
```

Should show `protect-crypto` under `pre-tool-use` with an
enabled indicator.

## Step 5 — iterate safely

If you discover a bug after enabling, **disable first, fix
second**:

```bash
ctx trigger disable protect-crypto
# ...edit the script...
ctx trigger test pre-tool-use --tool write_file --path internal/crypto/aes.go
ctx trigger enable protect-crypto
```

Disabling simply clears the executable bit — the script stays
on disk, and `ctx trigger enable` re-enables it without
rewriting anything.

## Patterns worth copying

### Logging, not blocking

For auditing or analytics, return `{"action":"allow"}` always
and append to a log as a side effect:

```bash
#!/usr/bin/env bash
set -euo pipefail
payload=$(cat)
echo "$payload" >> .context/logs/tool-use.jsonl
echo '{"action":"allow"}'
```

### Context injection at session start

A `session-start` trigger can prepend text to the agent's
initial prompt by emitting `{"action":"inject", "content": "..."}`
— useful for injecting daily standup notes, open PRs, or
rotating TODOs without storing them in a steering file.

### Chaining triggers of the same type

Multiple scripts in the same type directory all run. If any
returns `action: block`, the block wins. Keep individual
triggers single-purpose and rely on composition.

## Common mistakes

**Forgetting the shebang.** Without `#!/usr/bin/env bash`,
the trigger won't execute even with the executable bit set.

**Not quoting `$path`.** If you use `$path` in a command
substitution or a `case` glob without quoting, a file name
with spaces or metacharacters will break the trigger in
surprising ways.

**Enabling before testing.** `ctx trigger enable` makes the
script live immediately. Always `ctx trigger test` first.

**Outputting non-JSON.** The trigger's stdout must be valid
JSON or ctx's trigger runner will log a parse error. Use
`jq -n` to construct output rather than hand-writing JSON
strings.

**Mixing `hook` and `trigger` vocabulary.** The command is
`ctx trigger` but the on-disk directory is `.context/hooks/`.
The feature was renamed; the directory name lags behind.
Don't let this confuse you — they refer to the same thing.

## See also

- [`ctx trigger` reference](../cli/trigger.md) — full
  command, flag, and event-type reference.
- [`ctx steering`](../cli/steering.md) — persistent rules,
  not scripts. Use steering when the thing you want is "tell
  the AI to always do X" rather than "run a script when Y
  happens."
- [Writing steering files](steering.md) — the rule-based
  equivalent of this recipe.
