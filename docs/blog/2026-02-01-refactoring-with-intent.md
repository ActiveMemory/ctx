---
title: "Refactoring with Intent: Human-Guided Sessions in AI Development"
date: 2026-02-01
author: Jose Alekhinne
reviewed_and_finalized: true
topics:
  - refactoring
  - code quality
  - documentation standards
  - module decomposition
  - YOLO versus intentional development
---

# Refactoring with Intent

![ctx](../images/ctx-banner.png)

## Human-Guided Sessions in AI Development

*Jose Alekhinne / 2026-02-01*

!!! question "What Happens When You Slow Down?"
    YOLO mode shipped 14 commands in a week. 

    But technical debt doesn't send invoices: It just waits.

This is the story of what happened when I stopped auto-accepting everything
and started **guiding the AI with intent**. 

The result: 27 commits across 4 days, a major version release, and 
lessons that apply far beyond `ctx`.

!!! info "The Refactoring Window"
    **January 28 - February 1, 2026**

    From commit `bb1cd20` to the v0.2.0 release merge.  
    *(this window matters more than the individual commits:
    it's where intent replaced velocity.)*

## The Velocity Trap

In the [previous post][first-post], I documented the "*YOLO mode*" that birthed
`ctx`: auto-accept everything, let the AI run free, ship features fast.

It worked: **until it didn't**.

[first-post]: 2026-01-27-building-ctx-using-ctx.md "Building ctx Using ctx"

The codebase had accumulated patterns I didn't notice during the sprint:

| YOLO Pattern            | Where Found            | Why It Hurts                       |
|-------------------------|------------------------|------------------------------------|
| `"TASKS.md"` as literal | 10+ files              | One typo = silent failure          |
| `dir + "/" + file`      | Path construction      | Breaks on Windows                  |
| Monolithic `embed.go`   | 150+ lines, 5 concerns | Untestable, hard to extend         |
| Inconsistent docstrings | Everywhere             | AI can't learn project conventions |

I didn't see these during "*YOLO mode*" because, honestly, **I wasn't looking**.

**Auto-accept means auto-ignore**.

In YOLO mode, every file you open looks fine until you try to change it.  

In contrast, refactoring mode is **when** you start paying attention to that 
**hidden friction**.

## The Shift: From **Velocity** to **Intent**

On January 28th, I changed the workflow:

1. **Read every diff before accepting**.
2. **Ask "why this way?" before committing**.
3. **Document patterns, not just features**.

The first commit of this era was telling:

```text
feat: add structured attributes to context. update XML format
```

Not a new feature: A *refinement*:

The `XML` format for context updates needed `type` and `timestamp` attributes. 

YOLO mode would have shipped something that worked. Intentional mode asked:  

**"What does well-structured look like?"**

## The Decomposition: `embed.go`

The most satisfying refactor was splitting `internal/claude/embed.go`.

**Before**: One 153-line file doing five things:

* Command registration
* Hook generation
* Permission handling
* Script templates
* Type definitions

... your "*de facto*" God object.

**After**: Five focused modules:

| File        | Lines | Responsibility       |
|-------------|-------|----------------------|
| `cmd.go`    | 46    | Command registration |
| `hook.go`   | 64    | Hook configuration   |
| `perm.go`   | 25    | Permission handling  |
| `script.go` | 47    | Script templates     |
| `types.go`  | 7     | Type definitions     |

The refactor also renamed functions to follow Go conventions:

```go
// Before: unnecessary prefixes
GetAutoSaveScript()
GetBlockNonPathCtxScript()
ListCommands()
CreateDefaultHooks()

// After: idiomatic Go
AutoSaveScript()
BlockNonPathCtxScript()
Commands()
DefaultHooks()
```

This wasn't about character count. It was about **teaching the AI**
what good Go looks like in **this** project.

!!! note "Project Conventions"
    What I wanted from AI was to **understand** and **follow** the project's 
    conventions, and **trust** the author.

The next time it generates code, it has **better** examples to learn from.

## The Documentation Debt

YOLO mode created features. It didn't create documentation standards.

The January 29th sessions focused on standardization.

### Terminology Fixes

* "context-update" → "entry" (*what users actually call them*)
* Consistent naming across CLI, docs, and code comments

### Go Docstrings

```go
// Before: inconsistent or missing
func Parse(s string) Entry { ... }

// After: standardized sections

// Parse extracts an entry from a markdown string.
//
// Parameters:
//   - s: The markdown string to parse
//
// Returns:
//   - Entry with populated fields, or zero value if parsing fails
func Parse(s string) Entry { ... }
```

This is intentionally more structured than typical GoDoc:

It serves as documentation **and** doubles as training data for future 
AI-generated code.

### CLI Output Convention

```markdown
All CLI output follows: [emoji] [Title]: [message]

Examples:
  ✓ Decision added: Use symbolic types for entry categories
  ⚠ Warning: No tasks found
  ✗ Error: File not found
```

A **consistent** output shape makes both **human scanning** and 
**AI reasoning** more **reliable**.

These aren't exciting commits. But they are **force multipliers**:

Every future AI session now has **better** examples to follow.

## The Journal System

**If you only read one section, read this one**:

This is where **v0.2.0** becomes more than a refactor.

The biggest feature of this change window wasn't a refactor; it was
the **journal system**.

!!! note "45 files changed, 1680 insertions"
    This commit added the infrastructure for synthesizing AI session
    history into human-readable content.

The **journal system** includes:

| Component                | Purpose                                            |
|--------------------------|----------------------------------------------------|
| `ctx recall export`      | Export sessions to markdown in `.context/journal/` |
| `ctx journal site`       | Generate static site from journal entries          |
| `ctx serve`              | Convenience wrapper for the static site server     |
| `/ctx-journal-enrich`    | Slash command to add frontmatter and tags          |
| `/ctx-blog`              | Generate blog posts from recent activity           |
| `/ctx-blog-changelog`    | Generate changelog-style blog posts                |

...and the meta continues: **this blog post was generated using `/ctx-blog`**.

The session history from January 28-31 was

* **exported**, 
* **enriched**,
* and **synthesized**.

into the narrative you are reading.

## The Constants Consolidation

The final refactoring session addressed the remaining *magic strings*:

```go
const (
    // Comment markers
    CommentOpen  = "<!--"
    CommentClose = "-->"

    // Index markers
    MarkerIndexStart = "<!-- INDEX:START -->"
    MarkerIndexEnd   = "<!-- INDEX:END -->"

    // Newlines
    NewlineLF   = "\n"
    NewlineCRLF = "\r\n"
)
```

The work also introduced **thread safety** in the **recall parser** and
centralized shared validation logic; removing duplication that had
quietly spread during YOLO mode.

## I (Re)learned My Lessons

Similar to what I've learned in 
[the former human-assisted refactoring post][first-post], this
journey also made me realize that "*AI-only code generation*"
isn't sustainable in the long term.

### 1. Velocity and Quality Aren't Opposites

YOLO mode has its place: for *prototyping*, *exploration*, and *discovery*.

**BUT** (*and it's a huge "but"*), it needs to be followed by 
[**consolidation sessions**][ratio].

The ratio that worked for me: **[3:1][ratio]**.

* Three YOLO sessions create enough surface area to reveal patterns;
* the fourth session turns those patterns into structure.

[ratio]: 2026-02-17-the-3-1-ratio.md "The 3:1 Ratio: the evidence and the practice"

### 2. Documentation **IS** Code

When I standardized docstrings, I wasn't just writing docs.
I was **training future AI sessions**.

Every example of good code becomes a template for generated code.

### 3. Decomposition > Deletion

When `embed.go` became unwieldy, the temptation was to remove functionality.

The right answer was decomposition:

* Same functionality;
* **Better** organization;
* **Easier** to test;
* **Easier** to extend.

The result: more lines overall, but dramatically better structure.

!!! tip "The AI Benefit"
    Smaller, focused files also help AI assistants. 

    When a file fits comfortably in the context window, the AI can 
    reason about it *completely* instead of working from truncated snippets,
    preserving token budget for the actual task.

### 4. Meta-Tools Pay Dividends

The **journal system** took almost a full day to implement.

Yet it paid for itself immediately:

* This blog post was generated from session history;
* Future posts will be easier;
* The archaeological record is now **browsable**, not just `grep`-able.

## The Release: v0.2.0

The refactoring window culminated in the v0.2.0 release.

What's in v0.2.0:

| Category      | Changes                                                      |
|---------------|--------------------------------------------------------------|
| **Features**  | Journal system, quick reference indexes, global flags        |
| **Refactors** | Module decomposition, constants consolidation, CRLF handling |
| **Docs**      | Standardized terminology, Go docstrings, CLI conventions     |
| **Quality**   | Thread safety, shared validation, linter fixes               |

The version bump was symbolic.

The real change was how the codebase felt.

Opening files no longer triggered the familiar
*"ugh, I need to clean this up"* reaction.

## The Meta Continues

This post was written using the tools built during this refactoring window:

1. Session history exported via `ctx recall export`;
2. Journal entries enriched via `/ctx-journal-enrich`;
3. Blog draft generated via `/ctx-blog`;
4. Final editing done (*by yours truly*), with full project context loaded.

!!! info "The Context Is Massive"
    The `ctx` session files now contain 50+ development snapshots: each one
    capturing **decisions**, **learnings**, and **intent**.

!!! quote "The Moral of the Story"
    * **YOLO mode** builds **the prototype**.
    * **Intentional mode** builds **the product**.


Schedule both, or you'll only get one, **if** you're lucky.

---

*This blog post was generated with the help of `ctx`, using session
history, decision logs, learning logs, and git history from the
refactoring window. The meta continues.*
