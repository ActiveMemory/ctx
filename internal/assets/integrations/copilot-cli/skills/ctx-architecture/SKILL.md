---
name: ctx-architecture
description: "Build and maintain ARCHITECTURE.md and DETAILED_DESIGN.md. Use when working on structure, adding packages, or tracing flow."
tools: [bash, read, write, glob, grep]
---

Build and maintain architecture documentation with incremental
coverage tracking.

## When to Use

- Working on system structure or adding packages
- Tracing data flow across the codebase
- Onboarding to understand the system
- After significant structural changes

## When NOT to Use

- For code-level changes within a single package
- When ARCHITECTURE.md is already up-to-date
- For documentation-only projects

## Process

### 1. Scan the codebase

```bash
ctx status
```

Read the existing ARCHITECTURE.md if it exists.
Scan the directory tree to identify:
- Top-level packages and their responsibilities
- Data flow between components
- External dependencies

### 2. Build or update ARCHITECTURE.md

Structure:
- **Overview**: 2-3 sentence system description
- **Package Map**: table of packages → responsibilities
- **Data Flow**: how data moves through the system
- **Key Interfaces**: important boundaries
- **Dependencies**: external deps and why they're used

### 3. Build DETAILED_DESIGN.md (optional)

Deeper dive into internals for complex packages:
- Function-level documentation
- State machines
- Error handling patterns
- Concurrency model

### 4. Coverage tracking

Track which packages have been documented:

```
Coverage: 18/24 packages documented (75%)
Missing: internal/hub, internal/crypto, ...
```

## Quality Checklist

- [ ] Every top-level package mentioned
- [ ] Data flow is traceable end-to-end
- [ ] External dependencies listed with rationale
- [ ] Coverage percentage reported
- [ ] No stale references to removed packages
