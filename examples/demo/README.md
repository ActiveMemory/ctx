# Demo Project

This is a sample project demonstrating Active Memory context structure.

## Using Active Memory

1. View the context files in `.context/`:
   - `CONSTITUTION.md` - Inviolable rules
   - `TASKS.md` - Current work items
   - `CONVENTIONS.md` - Coding standards
   - `ARCHITECTURE.md` - System overview
   - `DECISIONS.md` - Technical decisions

2. Run Active Memory commands:
   ```bash
   # View context status
   amem status

   # Get AI-ready context packet
   amem agent

   # Add a new task
   amem add task "Implement feature X"

   # Mark a task complete
   amem complete "feature X"

   # Check for stale context
   amem drift
   ```

## Context Structure

The `.context/` directory contains markdown files that provide persistent
context for AI coding assistants. This helps AI tools understand:

- What rules must never be broken (CONSTITUTION)
- What work is currently in progress (TASKS)
- How code should be written (CONVENTIONS)
- How the system is organized (ARCHITECTURE)
- Why things are the way they are (DECISIONS)
