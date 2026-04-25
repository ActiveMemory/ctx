---
name: ctx-sanitize-permissions
description: "Audit tool permissions for dangerous or overly broad entries. Use to ensure safe agent configuration."
tools: [bash, read, write]
---

Audit agent permission configurations for dangerous patterns.

## When to Use

- After initial project setup
- When reviewing security posture
- When permissions seem overly broad
- Before sharing a project configuration

## When NOT to Use

- No permission config exists
- Already audited recently

## Categories to Check

### 1. Hook bypass permissions
Permissions that disable safety hooks entirely.

### 2. Destructive command permissions
Allow patterns that cover `rm -rf`, `git push --force`,
`git reset --hard`, etc.

### 3. Injection vectors
Overly broad shell permissions that could allow arbitrary
command execution.

### 4. Overly broad wildcards
Permissions like `Bash(*)` or `Write(*)` that grant
unrestricted access.

## Process

1. Read the permission configuration file
2. Check each entry against the four categories
3. Flag dangerous entries with severity level
4. Propose safer alternatives
5. Apply fixes with user approval

## Output Format

```
## Permission Audit Results

### 🔴 Critical (N)
1. `Bash(*)`: unrestricted shell access
   → Suggest: scope to specific commands

### 🟡 Warning (N)
1. `Write(/etc/*)`: write access to system dirs
   → Suggest: remove or scope to project

### ✅ Clean (N entries passed)
```

## Quality Checklist

- [ ] All permission entries reviewed
- [ ] Critical items flagged
- [ ] Safer alternatives proposed
- [ ] No changes made without user approval
