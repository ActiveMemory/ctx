# PROMPT.md — Demo Project

## CORE PRINCIPLE

You have NO conversational memory. Your memory IS the file system.
Your goal: advance the project by exactly ONE task, update context, and exit.

---

## PROJECT CONTEXT

**Project**: Demo API Server
**Language**: Go 1.22+
**Current Focus**: Phase 2 — Authentication

---

## PHASE 0: ORIENT

1. Read `.context/TASKS.md` — Current work items
2. Read `.context/CONSTITUTION.md` — Rules to never violate
3. Read `.context/CONVENTIONS.md` — How to write code
4. Read relevant spec in `specs/` for the current task

---

## PHASE 1: SELECT TASK

1. Read `.context/TASKS.md`
2. Find the **first unchecked item** (line starting with `- [ ]`)
3. That is your ONE task for this iteration

**IF NO UNCHECKED ITEMS:**
1. Run validation: `go build ./...`, `go test ./...`
2. If all pass, output `<promise>PHASE_COMPLETE</promise>`
3. If any fail, add fix task and continue

---

## PHASE 2: EXECUTE

1. **Read the spec** — Check `specs/` for detailed requirements
2. **Search first** — Don't assume code doesn't exist
3. **Implement ONE task** — Complete it fully. No placeholders.
4. **Follow conventions** — Check `.context/CONVENTIONS.md`

---

## PHASE 3: VALIDATE

After implementing, run:

```bash
go build ./...           # Must compile
go test ./...            # Tests must pass
go vet ./...             # No vet errors
```

---

## PHASE 4: UPDATE CONTEXT

1. Mark completed task `[x]` in `.context/TASKS.md`
2. Add `#done:YYYY-MM-DD-HHMMSS` timestamp
3. If you made an architectural decision → add to `.context/DECISIONS.md`
4. If you learned a gotcha → add to `.context/LEARNINGS.md`

**EXIT.** Do not continue to next task. The loop will restart you.

---

## CRITICAL CONSTRAINTS

### ONE TASK ONLY
Complete ONE task, then stop. The loop handles continuation.

### NO CHAT
Never ask questions. If blocked:
1. Add reason to task in `.context/TASKS.md`
2. Move to next task

### MEMORY IS THE FILESYSTEM
You will not remember this conversation. Write everything important to files.

---

## REFERENCE: SPECS

| Spec | Description |
|------|-------------|
| `specs/oauth2.md` | OAuth2 authentication implementation |

---

Now read `.context/TASKS.md` and begin.
