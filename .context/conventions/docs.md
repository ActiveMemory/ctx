# docs

## Typography and Document Shape

See [`typography.md`](typography.md) for the full guide: Title Case
headings, monotype `` `ctx` ``, no em-dashes / smart quotes / quad
backticks, doc frontmatter / banner conventions, recipe arc, admonition
variants. Linters in `hack/` enforce the hard rules.

## Documentation

- **Godoc format**: Use canonical sections
  ```go
  // FunctionName does X.
  //
  // Longer description if needed.
  //
  // Parameters:
  //   - param1: Description
  //   - param2: Description
  //
  // Returns:
  //   - Type: Description of return value
  func FunctionName(param1, param2 string) error
  ```
- **Struct field documentation**: Exported structs with 2+ fields
  must document every field. Two accepted forms:
  ```go
  // Option A: Fields section in docblock (preferred for 4+ fields)
  // TypeName describes X.
  //
  // Fields:
  //   - FieldA: Description
  //   - FieldB: Description
  type TypeName struct {

  // Option B: Inline comments (acceptable for 2-3 fields)
  // TypeName describes X.
  type TypeName struct {
      // FieldA is the description.
      FieldA string
      FieldB string // Description
  }
  ```
- **Package doc in doc.go**: Each package gets a `doc.go` with package-level
  documentation describing behavior, not structure. Do NOT include
  `# File Organization` sections listing files — they drift when files are
  added, renamed, or removed, and the filesystem is self-documenting
- **Copyright headers**: All source files get the project copyright header

## Blog Publishing

- **Checklist for ideas/ → docs/blog/ promotion**:
  1. Update date in frontmatter to publish date
  2. Fix relative paths (from `../docs/blog/` to peer references)
  3. Add cross-links to/from companion posts ("See also" sections)
  4. Add "The Arc" section connecting to the series narrative
  5. Update `docs/blog/index.md` with entry (newest first)
  6. Verify all link targets exist
  7. Build and test before commit
- **Arc section**: Every post includes "The Arc" near the end, framing
  where the post sits in the broader blog narrative
- **See also links**: Use italic `*See also: [Title](file) -- one-line
  description connecting the two posts.*` format at the end of posts
- **Frontmatter**: Include copyright header, title, date, author, topics list
- **Blog index order**: Newest post first, with topic tags and 3-4 line summary

- **Update admonitions for historical blog content**: Use MkDocs admonitions
  (`!!! note "Update"`) at the top of blog post sections where features have
  been superseded or installation has changed. Link to current documentation.
  Keep original content intact below for historical context.
- **New CLI subcommand documentation checklist**: Update docs in at least
  three places: (1) Feature page — commands table, usage section, skill/NL
  table. (2) CLI reference — full reference entry with args, flags, examples.
  (3) Relevant recipes. (4) zensical.toml — only if adding a new page.
- **Rename/refactor documentation checklist**: Scope ALL documentation impact
  before implementation. Three anchors plus one tangential: (1) Docstrings.
  (2) User-facing docs (`docs/`). (3) Recipes (`docs/recipes/`). (4) Blog
  posts and release notes. Also check: skills, hook messages, YAML text
  files, `.context/` files, and specs.
- **Stage site/ with docs/ changes**: The generated HTML is tracked in git
  with no CI build step

