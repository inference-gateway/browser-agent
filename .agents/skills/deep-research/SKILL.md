---
name: deep-research
description: Use this when the user asks an open-ended question that needs synthesis from multiple web sources. Plans sub-questions, drives a search engine, visits and cross-references sources via navigate_to_url + extract_data, and writes a cited markdown report with write.
tags:
  - research
  - synthesis
  - citations
  - investigation
license: Apache-2.0
---

# deep-research

TODO: Describe when and how the agent should use this skill. Lead with an
action-oriented "Use this when…" sentence so the model can decide whether
to apply it. The full body of this file is prepended to the system prompt
at runtime.

## When to use

Describe the user intents or task shapes that should trigger this skill.
Be concrete - list the kinds of requests, signals, or context that map to
this playbook.

## Workflow

1. ...
2. ...
3. ...

## Tools

List the tools this skill expects to call (declared under `spec.tools` in
the ADL manifest), and the order in which they're typically invoked.

## Bundled assets

This skill lives in its own directory under `.agents/skills/deep-research/`
(also reachable as `.claude/skills/deep-research/` via the generated
`.claude/skills` -> `../.agents/skills` symlink). You can ship arbitrary scripts, templates, or
reference material alongside `SKILL.md` - the `.adl-ignore` file protects
the whole directory from being clobbered on regeneration. Suggested layout:

```
.agents/skills/deep-research/
├── SKILL.md          # this file
├── scripts/          # optional helper scripts (Python, shell, etc.)
├── templates/        # optional file templates the agent can fill in
└── resources/        # optional static reference material
```

Reference bundled files by relative path from `SKILL.md` (e.g.
`scripts/triage.py`, `templates/report.md`) so the agent can locate them
at runtime.
