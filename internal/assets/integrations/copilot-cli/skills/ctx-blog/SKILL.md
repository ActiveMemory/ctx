---
name: ctx-blog
description: "Generate blog post drafts from project activity. Use to communicate progress, decisions, or technical insights."
tools: [bash, read, write]
---

Generate blog post drafts from project activity.

## When to Use

- After completing a significant feature
- When a decision or learning is worth sharing publicly
- For project updates and changelogs
- When the user says "write a blog post about..."

## When NOT to Use

- For internal context (use learnings/decisions instead)
- When there's nothing noteworthy to share

## Process

### 1. Gather material

- Recent commits: `git log --oneline -20`
- Recent decisions from DECISIONS.md
- Recent learnings from LEARNINGS.md
- Completed tasks from TASKS.md

### 2. Identify the narrative

What's the story? Options:
- Feature announcement
- Technical deep-dive
- Lessons learned
- Project update / changelog

### 3. Draft the post

Structure:
- **Title**: clear and engaging
- **Introduction**: what and why (2-3 sentences)
- **Body**: the story with technical details
- **Conclusion**: what's next

### 4. Write to blog directory

Target: `site/blog/{date}-{slug}/index.html` or
`docs/blog/{date}-{slug}.md` per project convention.

## Quality Checklist

- [ ] Title is clear and engaging
- [ ] Technical accuracy verified
- [ ] No sensitive information exposed
- [ ] Proper frontmatter/metadata
- [ ] Links to relevant specs/docs where appropriate
