# Security Architecture (Diagrams)

Parent: [ARCHITECTURE.md](ARCHITECTURE.md)

## Defense Layers

```
  ┌─────────────────────────────────────────────────────────┐
  │                    Layer 5: Plugin Hooks                  │
  │  block-non-path-ctx: reject ./ctx or /abs/path/ctx      │
  │  qa-reminder: gate before commit                         │
  ├─────────────────────────────────────────────────────────┤
  │                    Layer 4: Permission Deny List          │
  │  Bash(sudo *), Bash(rm -rf *), Bash(curl *),            │
  │  Bash(wget *), Bash(go install *), force push            │
  ├─────────────────────────────────────────────────────────┤
  │                    Layer 3: Boundary Validation           │
  │  ValidateBoundary(): resolved .context/ must be under    │
  │  project root (prevents path traversal)                  │
  ├─────────────────────────────────────────────────────────┤
  │                    Layer 2: Symlink Rejection             │
  │  CheckSymlinks(): .context/ dir and children must not    │
  │  be symlinks (M-2 defense against link attacks)          │
  ├─────────────────────────────────────────────────────────┤
  │                    Layer 1: File Permissions              │
  │  Keys: 0600 (owner rw)                                   │
  │  Executables: 0755                                       │
  │  Regular files: 0644                                     │
  ├─────────────────────────────────────────────────────────┤
  │                    Layer 0: Encryption                    │
  │  AES-256-GCM for scratchpad and webhook URLs             │
  │  12-byte random nonce + 16-byte auth tag                 │
  └─────────────────────────────────────────────────────────┘
```

## Secret Detection (Drift Check)

```
  drift.Detect() scans for files matching:
    .env, credentials*, *secret*, *api_key*, *password*

  Exceptions (not flagged):
    *.example, *.sample, files with template markers

  Constitution invariants:
    Never commit secrets, tokens, API keys, or credentials
    Never store customer/user data in context files
```
