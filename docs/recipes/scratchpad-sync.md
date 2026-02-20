---
title: "Syncing Scratchpad Notes Across Machines"
icon: lucide/key-round
---

![ctx](../images/ctx-banner.png)

## Problem

You work from multiple machines — a desktop and a laptop, or a local
machine and a remote dev server.

The scratchpad entries are encrypted. The ciphertext (`.context/scratchpad.enc`)
travels with git, but the encryption key (`.context/.scratchpad.key`) is
gitignored. Without the key on each machine, you cannot read or write entries.

**How do you distribute the key and keep the scratchpad in sync?**

!!! tip "TL;DR"
    ```bash
    ctx init                                                 # 1. generates .scratchpad.key
    scp .context/.scratchpad.key user@machine-b:project/.context/  # 2. copy key
    chmod 600 project/.context/.scratchpad.key                # 3. secure it
    # Normal git push/pull syncs the encrypted scratchpad.enc
    # On conflict: ctx pad resolve → rebuild → git add + commit
    ```

## Commands and Skills Used

| Tool | Type | Purpose |
|------|------|---------|
| `ctx init` | CLI command | Initialize context (generates key automatically) |
| `ctx pad add` | CLI command | Add a scratchpad entry |
| `ctx pad rm` | CLI command | Remove a scratchpad entry |
| `ctx pad edit` | CLI command | Edit a scratchpad entry |
| `ctx pad resolve` | CLI command | Show both sides of a merge conflict |
| `ctx pad merge` | CLI command | Merge entries from other scratchpad files |
| `ctx pad import` | CLI command | Bulk-import lines from a file |
| `ctx pad export` | CLI command | Export blob entries to a directory |
| `scp` | Shell | Copy the key file between machines |
| `git push` / `git pull` | Shell | Sync the encrypted file via git |
| `/ctx-pad` | Skill | Natural language interface to pad commands |

## The Workflow

### Step 1: Initialize on Machine A

Run `ctx init` on your first machine. The key is created automatically:

```bash
ctx init
# ...
# Created .context/.scratchpad.key (0600)
# Created .context/scratchpad.enc
```

The key is gitignored. The `.enc` file is tracked.

### Step 2: Copy the Key to Machine B

Use any secure transfer method:

```bash
# scp
scp .context/.scratchpad.key user@machine-b:project/.context/

# Or use a password manager, USB drive, etc.
```

Set permissions on Machine B:

```bash
chmod 600 project/.context/.scratchpad.key
```

!!! warning "Secure the Transfer"
    The key is a raw 256-bit AES key. Anyone with the key can decrypt
    the scratchpad. Use an encrypted channel (SSH, password manager
    vault) — never paste it in plaintext over email or chat.

### Step 3: Normal Push/Pull Workflow

The encrypted file is committed, so standard git sync works:

```bash
# Machine A: add entries and push
ctx pad add "staging API key: sk-test-abc123"
git add .context/scratchpad.enc
git commit -m "Update scratchpad"
git push

# Machine B: pull and read
git pull
ctx pad
#   1. staging API key: sk-test-abc123
```

Both machines have the same key, so both can decrypt the same `.enc` file.

### Step 4: Read and Write from Either Machine

Once the key is distributed, all `ctx pad` commands work identically on
both machines. Entries added on Machine A are visible on Machine B after
a `git pull`, and vice versa.

### Step 5: Handle Merge Conflicts

If both machines add entries between syncs, pulling will create a merge
conflict on `.context/scratchpad.enc`. Git cannot merge binary (encrypted)
content automatically.

The fastest approach is `ctx pad merge` — it reads both conflict sides,
deduplicates, and writes the union:

```bash
# Extract theirs to a temp file, then merge it in
git show :3:.context/scratchpad.enc > /tmp/theirs.enc
git checkout --ours .context/scratchpad.enc
ctx pad merge /tmp/theirs.enc

# Done — commit the resolved scratchpad
git add .context/scratchpad.enc
git commit -m "Resolve scratchpad merge conflict"
```

Alternatively, use `ctx pad resolve` to inspect both sides manually:

```bash
ctx pad resolve
# === Ours (this machine) ===
#   1. staging API key: sk-test-abc123
#   2. check DNS after deploy
#
# === Theirs (incoming) ===
#   1. staging API key: sk-test-abc123
#   2. new endpoint: api.example.com/v2
```

Then reconstruct the merged scratchpad:

```bash
# Start fresh with all entries from both sides
ctx pad add "staging API key: sk-test-abc123"
ctx pad add "check DNS after deploy"
ctx pad add "new endpoint: api.example.com/v2"

# Mark the conflict resolved
git add .context/scratchpad.enc
git commit -m "Resolve scratchpad merge conflict"
```

## Merge Conflict Walkthrough

Here's a full scenario showing how conflicts arise and how to resolve them:

**1. Both machines start in sync** (1 entry):

```
Machine A: 1. staging API key: sk-test-abc123
Machine B: 1. staging API key: sk-test-abc123
```

**2. Both add entries independently**:

```
Machine A adds: "check DNS after deploy"
Machine B adds: "new endpoint: api.example.com/v2"
```

**3. Machine A pushes first. Machine B pulls and gets a conflict**:

```bash
git pull
# CONFLICT (content): Merge conflict in .context/scratchpad.enc
```

**4. Machine B runs `ctx pad resolve`**:

```bash
ctx pad resolve
# === Ours ===
#   1. staging API key: sk-test-abc123
#   2. new endpoint: api.example.com/v2
#
# === Theirs ===
#   1. staging API key: sk-test-abc123
#   2. check DNS after deploy
```

**5. Rebuild with entries from both sides and commit**:

```bash
# Clear and rebuild (or use the skill to guide you)
ctx pad add "staging API key: sk-test-abc123"
ctx pad add "check DNS after deploy"
ctx pad add "new endpoint: api.example.com/v2"

git add .context/scratchpad.enc
git commit -m "Merge scratchpad: keep entries from both machines"
```

### Conversational Approach

When working with an AI assistant, you can resolve conflicts naturally:

```text
You: "I have a scratchpad merge conflict. Can you resolve it?"

Agent: "Let me extract theirs and merge it in."
       [runs git show :3:.context/scratchpad.enc > /tmp/theirs.enc]
       [runs git checkout --ours .context/scratchpad.enc]
       [runs ctx pad merge /tmp/theirs.enc]
       "Merged 2 new entries (1 duplicate skipped). Want me to
       commit the resolution?"
```

## Tips

* **Back up the key**. If you lose it, you lose access to all encrypted
  entries. Store a copy in your password manager.
* **One key per project**. Each `ctx init` generates a unique key.
  Don't reuse keys across projects.
* **Plaintext fallback for non-sensitive projects**. If encryption adds
  friction and you have nothing sensitive, set `scratchpad_encrypt: false`
  in `.contextrc`. Merge conflicts become trivial text merges.
* **Never commit the key**. It's gitignored by default. Don't override
  this.

## Next Up

**[Parallel Agent Development with Git Worktrees →](parallel-worktrees.md)**:
Run multiple agents on independent task tracks using git worktrees.

## See Also

* [Scratchpad](../scratchpad.md): feature overview, all commands, when
  to use scratchpad vs context files
* [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md):
  for structured knowledge that outlives the scratchpad
