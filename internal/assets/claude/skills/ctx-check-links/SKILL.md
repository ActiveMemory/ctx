---
name: ctx-check-links
description: "Audit docs for dead links. Use before releases, after restructuring docs, or when /_ctx-audit runs."
allowed-tools: Bash(curl:*), Read, Grep, Glob
---

Scan all markdown files under `docs/` for broken links. Two passes:
internal (file targets) and external (HTTP URLs).

## When to Use

- Before releases or doc deployments
- After renaming, moving, or deleting doc pages
- After restructuring the `docs/` directory or nav
- When `/_ctx-audit` runs (audit check #12)
- When a user reports a 404 on the site

## When NOT to Use

- When editing a single doc (just eyeball links in that file)
- When offline and only external checks would matter

## Execution

### Pass 1: Internal Links

Scan every `.md` file in `docs/` for markdown links pointing to
other files: `[text](target.md)`, `[text](../path/file.md)`,
`[text](path/file.md#anchor)`.

For each link:

1. Resolve the target **relative to the source file's directory**
2. Strip any `#anchor` fragment before checking file existence
3. Skip external URLs (`http://`, `https://`, `mailto:`)
4. Skip bare anchors (`#section-name`) — these are intra-page
5. Verify the target file exists on disk

Collect all broken internal links as:

```
BROKEN: source-file.md:LINE → target.md (file not found)
```

### Pass 2: External Links

Scan every `.md` file in `docs/` for `http://` and `https://` URLs
in markdown link syntax.

For each URL:

1. Send an HTTP HEAD request with a 10-second timeout
2. If HEAD fails or returns 405, retry with GET
3. Record the HTTP status code

Report failures as:

```
WARN: source-file.md:LINE → https://example.com (HTTP 404)
WARN: source-file.md:LINE → https://example.com (timeout)
```

**Do not treat external failures as errors.** Network partitions,
rate limiting, and transient outages are common. Report them but
do not fail the check.

Exceptions — skip these URLs:
- `localhost` / `127.0.0.1` URLs (local dev servers)
- `example.com` / `example.org` (placeholder domains)

### Pass 3: Image References

Scan for image links: `![alt](path/to/image.png)` and
`![alt](images/file.jpg)`.

Verify the image file exists on disk. Same resolution rules as
internal links.

## Output Format

```
## Link Check Report

### Internal Links
- N broken links found (or "All clear")
- [list of broken links with file:line and target]

### External Links
- N warnings (or "All reachable")
- [list of failures with file:line, URL, and reason]

### Images
- N missing images (or "All present")
- [list of missing images with file:line and target]

### Summary
Internal: N broken / M total
External: N unreachable / M total
Images: N missing / M total
```

## Fixing

For broken internal links, offer specific fixes:

- If the target was renamed, suggest the new path
- If the target was deleted, suggest removing the link or
  pointing to an alternative
- If the target is a typo (close match exists), suggest the
  correction

For external links, just report. The user decides whether to
update, remove, or ignore.

## Integration with /_ctx-audit

When invoked as check #12 from `/_ctx-audit`:

- Run the full check (all 3 passes)
- Report findings in the same format as other consolidation checks
- Internal broken links count as findings to fix
- External failures count as warnings (informational)

## Quality Checklist

After running the check:
- [ ] All `.md` files under `docs/` were scanned
- [ ] Relative path resolution accounts for subdirectories
  (`recipes/`, `blog/`)
- [ ] Anchors stripped before file existence check
- [ ] External check used timeouts (not hanging on slow hosts)
- [ ] localhost/example URLs were skipped
- [ ] Report distinguishes errors (internal) from warnings
  (external)
