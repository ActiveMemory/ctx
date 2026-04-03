---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: "Code Structure as an Agent Interface:
  What 19 AST Tests Taught Us About Agent-Readable Code"
date: 2026-04-02
author: Jose Alekhinne
reviewed_and_finalized: true
topics:
  - ast
  - code quality
  - agent readability
  - conventions
  - field notes
---

# Code Structure as an Agent Interface

## What 19 AST Tests Taught Us About Agent-Readable Code

![ctx](../images/ctx-banner.png)

**When an agent sees `token.Slash` instead of `"/"`, it cannot pattern-match 
against the millions of `strings.Split(s, "/")` calls in its training data
and coast on statistical inference. It has to actually look up what 
`token.Slash` is.**

*Jose Alekhinne / April 2, 2026*

We set out to replace a shell script with Go tests.

We ended up discovering that "**code quality**" and
"**agent readability**" are the **same thing**.

This is not about linting. This is about controlling
how an agent **perceives** your system.

One term will recur throughout this post, so let me pin it down:

!!! tip "Agent Readability"
    **Agent Readability** is the degree to which a codebase
    can be understood through structured traversal, not 
    statistical pattern matching.

This is the story of **19 AST-based audit tests**, a
single-day session that touched **300+ files**, and what
happens when you treat your codebase's structure as
an interface for the machines that read it.

---

## The Shell Script Problem

`ctx` had a file called `hack/lint-drift.sh`. It ran
five checks using `grep` and `awk`: literal `"\n"`
strings, `cmd.Printf` calls outside the write package,
magic directory strings in `filepath.Join`, hardcoded
`.md` extensions, and DescKey-to-YAML linkage.

It worked. **Until it didn't**.

The script had three structural weaknesses that kept biting us:

1. **No type awareness.** It could not distinguish a
   `Use*` constant from a `DescKey*` constant, causing
   71 false positives in one run.
2. **Fragile exclusions.** When a constant moved from
   `token.go` to `whitespace.go`, the exclusion glob
   broke silently.
3. **Ceiling on detection.** Checks that require
   understanding call sites, import graphs, or type
   relationships are impossible in shell.

We wrote a *spec* to replace all five checks
with Go tests using `go/ast` and `go/packages`. The
tests would run as part of `go test ./...`: no separate
script, no separate CI step.

What we did not expect was where the work would lead.

## The AST Migration

The pattern for each test is identical:

```go
func TestNoLiteralWhitespace(t *testing.T) {
    pkgs := loadPackages(t)
    var violations []string
    for _, pkg := range pkgs {
        for _, file := range pkg.Syntax {
            ast.Inspect(file, func(n ast.Node) bool {
                // check node, append to violations
                return true
            })
        }
    }
    for _, v := range violations {
        t.Error(v)
    }
}
```

Load packages once via `sync.Once`, walk every syntax
tree, collect violations, report. The shared helpers
(`loadPackages`, `isTestFile`, `posString`) live in
`helpers_test.go`. Each test is a `_test.go` file in
`internal/audit/`, producing no binary output and not
importable by production code.

In a single session, we built 13 new tests on top of
6 that already existed, bringing the total to 19:

| Test                          | What it catches                                                |
|-------------------------------|----------------------------------------------------------------|
| `TestNoLiteralWhitespace`     | `"\n"`, `"\t"`, `'\r'` outside `config/token/`                 |
| `TestNoNakedErrors`           | `fmt.Errorf`/`errors.New` outside `internal/err/`              |
| `TestNoStrayErrFiles`         | `err.go` files outside `internal/err/`                         |
| `TestNoRawLogging`            | `fmt.Fprint*(os.Stderr)`, `log.Print*` outside `internal/log/` |
| `TestNoInlineSeparators`      | `strings.Join` with literal separator arg                      |
| `TestNoStringConcatPaths`     | Path-like variables built with `+`                             |
| `TestNoStutteryFunctions`     | `write.WriteJournal` repeats package name                      |
| `TestDocComments`             | Missing doc comments on any declaration                        |
| `TestNoMagicValues`           | Numeric literals outside const definitions                     |
| `TestNoMagicStrings`          | String literals outside const definitions                      |
| `TestLineLength`              | Lines exceeding 80 characters                                  |
| `TestNoRegexpOutsideRegexPkg` | `regexp.MustCompile` outside `config/regex/`                   |

Plus the six that preceded the session:
`TestNoErrorsAs`, `TestNoCmdPrintOutsideWrite`,
`TestNoExecOutsideExecPkg`, `TestNoInlineRegexpCompile`,
`TestNoRawFileIO`, `TestNoRawPermissions`.

The migration touched **300+ files across 25 commits**.

**Not** because the tests were hard to write, **but** because
every test we wrote revealed violations that needed fixing.

## The Tightening Loop

The most instructive part was not writing the tests.
It was the iterative tightening.

The following process was repeated for every test:

1. Write the test with reasonable exemptions
2. Run it, see violations
3. Fix the violations (*migrate to config constants*)
4. The human reviews the result
5. The human spots something the test missed
6. Fix the test first, verify it catches the issue
7. Fix the newly caught violations
8. Repeat from step 4

This loop drove the tests from "*basically correct*" to
"**actually useful**". 

Three examples:

### Example 1: The Local Const Loophole

`TestNoMagicValues` initially exempted local constants
inside function bodies. This let code like this pass:

```go
const descMaxWidth = 70
desc := truncateDescription(
    meta.Description, descMaxWidth,
)
```

The test saw a `const` definition and moved on. But
`const descMaxWidth = 70` on the line before its only
use is just renaming a magic number. The `70` should
live in `config/format/TruncateDescription` where it is
discoverable, reusable, and auditable.

We removed the local const exemption. The test caught
it. The value moved to config.

### Example 2: The Single-Character Dodge

`TestNoMagicStrings` initially exempted all single-character strings as 
"*structural punctuation*". 

This  let `"/"`, `"-"`, and `"."` pass everywhere.

But `"/"` is a **directory separator**. It is **OS-specific**
and a **security surface**. 

`"-"` used in `strings.Repeat("-", width)` is creating visual output,
not acting as a delimiter. 

`"."` in `strings.SplitN(ver, ".", 3)` is a version separator.

None of these are "*just punctuation*": 
They are **domain values with specific meanings**.

We removed the blanket exemption: 
30 violations surfaced. 

Every one was a real magic value that should have been 
`token.Slash`, `token.Dash`, or `token.Dot`.

### Example 3: The Replacer versus Regex

After migrating magic strings, we had this:

```go
func MermaidID(pkg string) string {
    r := strings.NewReplacer(
        token.Slash, token.Underscore,
        token.Dot, token.Underscore,
        token.Dash, token.Underscore,
    )
    return r.Replace(pkg)
}
```

Six token references and a `NewReplacer` allocation.
The magic values were gone, but we had replaced them
with token soup: **structure without abstraction.** 

The correct tool was a regex:

```go
// In config/regex/file.go:
var MermaidUnsafe = regexp.MustCompile(`[/.\-]`)

// In the caller:
func MermaidID(pkg string) string {
    return regex.MermaidUnsafe.ReplaceAllString(
        pkg, token.Underscore,
    )
}
```

One config regex, one call. The regex lives in
`config/regex/file.go` where every other compiled
pattern lives. An agent reading the code sees
`regex.MermaidUnsafe` and immediately knows: this is a
sanitization pattern, it lives in the regex registry,
and it has a name that explains its purpose.

**Clean is better than clever**.

---

## A Before-and-After

To make the agent-readability claim concrete, consider
one function through the full transformation.

**Before** (the code we started with):

```go
func MermaidID(pkg string) string {
    r := strings.NewReplacer(
        "/", "_", ".", "_", "-", "_",
    )
    return r.Replace(pkg)
}
```

An agent reading this sees six string literals. To
understand what the function does, it must: (1) parse
the `NewReplacer` pair semantics, (2) infer that `/`,
`.`, `-` are being replaced, (3) guess why, (4) hope
the guess is right.

There is nothing to follow. No import to trace. No
name to search. **The meaning is locked inside the
function body**.

**After** (*the code we ended with*):

```go
func MermaidID(pkg string) string {
    return regex.MermaidUnsafe.ReplaceAllString(
        pkg, token.Underscore,
    )
}
```

An agent reading this sees two named references:
`regex.MermaidUnsafe` and `token.Underscore`. 

To understand the function, it can: (1) look up
`MermaidUnsafe` in `config/regex/file.go` and see the
pattern `[/.\-]` with a doc comment explaining it
matches invalid Mermaid characters, (2) look up
`Underscore` in `config/token/delim.go` and see it is
the replacement character.

The agent now has: a named pattern, a named
replacement, a package location, documentation, and
neighboring context (*other regex patterns, other
delimiters*). 

It got all of this **for free** by following just two references.

The indirection is not an overhead. It is the **retrieval query**.

---

## The Principles

You are not just improving code quality. You are
shaping the **input space** that determines how an LLM
can reason about your system.

Every structural constraint we enforce converts
implicit semantics into explicit structure. 

**LLMs struggle when meaning is implicit and patterns are
statistical**. 

They thrive when **meaning is explicit** and **structure is navigable**.

Here is what we learned, organized into three
categories.

### Cognitive Constraints

These force agents (*and humans*) to **think harder**.

**Indirection acts as a built-in retrieval mechanism**:

Moving magic values to config forces the agent to
follow the reference. `errMemory.WriteFile(cause)`
tells the agent "there is a memory error package, go
look." `fmt.Errorf("writing MEMORY.md: %w", cause)`
inlines everything and makes the call graph invisible.
The indirection IS the retrieval query.

**Unfamiliar patterns force reasoning**:

When an agent sees `token.Slash` instead of `"/"`, it
cannot coast on corpus frequency. It has to actually
look up what `token.Slash` is, which forces it through
the dependency graph, which means it encounters
documentation and neighboring constants, which gives
it richer context. You are exploiting the agent's
weakness (over-reliance on training data) to make it
behave more carefully.

**Documentation helps everyone**:

Extensive documentation helps humans reading the code,
agents reasoning about it, and RAG systems indexing it.

Our `TestDocComments` check added 308 doc comments in
one commit. Every function, every type, every constant
block now has a doc comment. 

This is not busywork: it is the content that agents and embeddings consume.

### Structural Constraints

These shape the codebase into a navigable graph.

**Shorter files save tokens**:

Forcing private helper functions out of main files
makes the main file shorter. An agent loading a file
spends fewer tokens on boilerplate and more on the
logic that matters.

**Fixed-width constraints force decomposition**:

A function that cannot be expressed in 80 columns is
either too deeply nested (*extract a helper*), has too
many parameters (*introduce a struct)*, or has a variable
name that is too long (*rethink the abstraction*). 

The constraint forces structural improvements that happen
to also make the code more parseable.

**Chunk-friendly structure helps RAG**

Code intelligence tools chunk files for embedding and
retrieval. Short, well-documented, single-responsibility
files produce better chunks than monolithic files with
mixed concerns. 

The structural constraints create files
that RAG systems can index effectively.

**Centralization creates debuggable seams**:

All error handling in `internal/err/`, all logging in
`internal/log/`, all file operations in `internal/io/`.
One place to debug, one place to test, one place to see
patterns. An agent analyzing "how does this project
handle errors" gets one answer from one package, not
200 scattered `fmt.Errorf` calls.

**Private functions become public patterns**:

When you extract a private function to satisfy a
constraint, it often ends up as a semi-public function
in a `core/` package. Then you realize it is generic
enough to be factored into a purpose-specific module.

The constraint drives discovery of reusable
abstractions hiding inside monolithic functions.

### Operational Benefits

These pay dividends in daily development.

**Single-edit renames**:

Renaming a flag is one edit to a config constant
instead of find-and-replace across 30,000 lines with
possible misses. `grep token.Slash` gives you every
place that uses a forward slash semantically.

`grep "/"` gives you noise.

**Blast radius containment**:

When every magic value is a config constant, a search
is one result. This matters for impact analysis,
security audits, and agents trying to understand "*what
uses this*".

**Compile-time contract enforcement**:

When `err/memory.WriteFile` exists, the compiler
guarantees the error message exists and the call
signature is correct. An inline `fmt.Errorf` can have
a typo in the format string and nothing catches it
until runtime. Centralization turns runtime failures
into compile errors.

**Semantic `git blame`**:

When `token.Slash` is used everywhere and someone
changes its value, `git blame` on the config file
shows exactly when and why. 

With inline `"/"` scattered  across 30 files, the history is invisible.

**Test surface reduction**:

Centralizing into `internal/err/`, `internal/io/`,
`internal/config/` means you test behavior once at
the boundary and trust the callers. 

You do not need 30 tests for 30 `fmt.Errorf` calls. You need 1 test for
`errMemory.WriteFile` and 30 trivial call-site audits,
which is exactly what these AST tests provide.

## The Numbers

One session. 25 commits. The raw stats:

| Metric                      | Count |
|-----------------------------|-------|
| New audit tests             | 13    |
| Total audit tests           | 19    |
| Files touched               | 300+  |
| Magic values migrated       | 90+   |
| Functions renamed           | 17    |
| Doc comments added          | 323   |
| Lines rewrapped to 80 chars | 190   |
| Config constants created    | 40+   |
| Config regexes created      | 3     |

Every number represents a violation that existed
before the test caught it. The tests did not create
work: they revealed work that was already needed.

## The Uncomfortable Implication

None of this is Go-specific.

If an AI agent interacts with your codebase, your
codebase *already is* an interface. You just have not
designed it as one.

If your error messages are scattered across 200 files,
an agent cannot reason about error handling as a
concept. If your magic values are inlined, an agent
cannot distinguish "this is a path separator" from
"this is a division operator." If your functions are
named `write.WriteJournal`, the agent wastes tokens
on redundant information.

What we discovered, through the unglamorous work of
writing lint tests and migrating string literals, is
that the structural constraints software engineering
has valued for decades are exactly the constraints
that make code readable to machines.

**This is not a coincidence**: These constraints exist
because they **reduce the cognitive load of
understanding code**. 

**Agents have cognitive load too**: It is called **the context window**.

You are not converting code to a new paradigm.

You are **making the latent graph visible**.

You are converting implicit semantics into explicit
structure that both humans and machines can traverse.

## What's Next

The spec lists 8 more tests we have not built yet,
including `TestDescKeyYAMLLinkage` (verifying that
every DescKey constant has a corresponding YAML entry),
`TestCLICmdStructure` (enforcing the `cmd.go` /
`run.go` / `doc.go` file convention), and
`TestNoFlagBindOutsideFlagbind` (which requires
migrating ~50 flag registration sites first).

The broader question: should these principles be
codified as a reusable linting framework? The patterns
(`loadPackages` + `ast.Inspect` + violation collection)
are generic. 

The specific checks are project-specific.
But the *categories* of checks (*centralization
enforcement, magic value detection, naming conventions,
documentation requirements*) are universal.

For now, 19 tests in `internal/audit/` is enough.
They run in 2 seconds as part of `go test ./...`. They
catch real issues. 

And they encode a theory of code quality that serves
both humans and the agents that work alongside them.

---

Agents are not going away. They are reading your code
right now, forming representations of your system in
context windows that forget everything between sessions.

The codebases that structure themselves for that
reality will compound. The ones that do not will slowly
become illegible to the tools they depend on.

Structure is no longer just for maintainability. It is
for **reasonability**.
