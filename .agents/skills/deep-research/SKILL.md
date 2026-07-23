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

Use this when the user asks an open-ended question that needs synthesis
from multiple web sources rather than extraction from a known URL set.
The output is a cited markdown report on disk that the user can read,
share, and audit back to the original sources.

## When to use

Trigger this playbook on open-ended questions that need synthesis
across sources:

- "What are the current best practices for X?"
- "Compare A vs B vs C for use case Y"
- "What's the latest on `<topic>` in 2026?"
- "Find evidence for/against claim Z"

Do **not** use this for structured extraction from known URLs (use
`web-scraping`) or validating a known flow (use `webapp-testing`).

The distinguishing signal is **discovery**: if the user already knows
which pages to read, it's a scraping/testing job. If the agent has to
figure out *which sources matter*, weigh them, and synthesize - this
skill.

## Workflow: plan-then-execute (6 steps)

The most common failure modes in deep research are (a) jumping
straight to the first search result, (b) treating N restatements of
one primary source as N independent confirmations, and (c)
hallucinating a citation that nothing on disk actually supports. The
workflow is structured to prevent all three.

1. **Decompose** - break the question into 3-7 sub-questions covering
   distinct angles: technical, comparative, contrarian, recency. A
   good decomposition has at least one sub-question whose answer
   could *contradict* the others - if every sub-question points the
   same direction, you're building an echo chamber.
   - `write` the plan to `/tmp/research-<timestamp>/plan.md` so the
     run is auditable and resumable. Use a unix timestamp for
     `<timestamp>` so concurrent runs don't collide.

2. **Search** - one query per sub-question.
   - **Default endpoint: `https://html.duckduckgo.com/`** - it renders
     without JS and is friendlier to automated clients than Google.
     `navigate_to_url` with `wait_until: networkidle`.
   - `extract_data` the top results (title, URL, snippet) per query.
     The DuckDuckGo HTML result container is typically
     `.result` with `.result__a` (title+href), `.result__snippet`
     (snippet). Confirm via `execute_script` if the layout has
     changed.
   - **Swapping the search backend**: if DuckDuckGo HTML gets blocked
     or rate-limits, the simplest swap is to point at another
     no-JS-required HTML search frontend (e.g. a hosted SearXNG
     instance the user trusts, or `https://www.bing.com/search` with
     manual selector tweaks). Update step 2's URL and the result
     selectors; the rest of the workflow is search-engine-agnostic.

3. **Triage sources** - score each result by domain authority,
   recency, and topical relevance to its sub-question. Keep 2-4 per
   sub-question. Heuristics:
   - **Prefer primary sources** (official docs, vendor engineering
     blogs, RFCs, peer-reviewed papers, authoritative project READMEs)
     over aggregators ("Top 10 X in 2026" listicles, content farms).
   - **Date check**: when the question is time-sensitive ("latest",
     "current", "in 2026"), drop results older than ~12 months unless
     they're foundational references.
   - **Diversity check**: if the kept corpus all skews to one
     viewpoint (e.g. all from one vendor, all pro-X), add a
     deliberate contrarian search ("X limitations", "problems with X",
     "<competitor> on X") before moving on.

4. **Gather** - for each kept source, pick the right transport
   *before* committing to a Playwright session:
   - **Static text** (raw GitHub files, RFCs, plaintext docs, blog
     posts that render server-side): `fetch` returns the body
     directly. Much faster than rendering and leaves the body
     trivially greppable for citation. Verify it's not a hydration
     shell by checking the response body actually contains the
     article text rather than `<div id="root"></div>`.
   - **PDFs** (research papers, vendor whitepapers, standards
     documents): `fetch` with `save_path` into
     `/tmp/research-<timestamp>/sources/<n>.pdf`, then summarize
     from disk. The browser cannot extract PDF text reliably.
   - **Structured data** (JSON APIs, RSS/Atom feeds, OpenAPI specs):
     `fetch` the endpoint directly - it's the citable source, not
     the human-facing page that wraps it.
   - **JS-rendered articles** (most modern news/blog platforms,
     anything where the source HTML is a hydration shell): use the
     `navigate_to_url` path below.

   For the JS-rendered path:
   - `navigate_to_url` with `wait_until: networkidle` (JS-rendered
     articles often hydrate after the initial load).
   - `take_screenshot` with `full_page: true` for audit.
   - `extract_data` the main content - article body, paragraphs,
     headings. Common selectors: `article`, `main`,
     `[role="main"]`, `.post-content`. When the article is hidden
     behind a paywall fade or cookie wall, use `execute_script` to
     dump `document.querySelector('article').innerText` and check
     whether it's truncated.
   - Persist each source's extract to
     `/tmp/research-<timestamp>/sources/<n>.md` with the URL and
     publication date at the top. Numbered files keep citations
     stable - `[3]` in the report maps to `sources/3.md`. This means
     citations survive a crash and stay traceable end-to-end.

5. **Verify & cross-reference** - every load-bearing claim must be
   supported by **≥2 independent sources** or be marked
   `[single-source]` in the report.
   - "Independent" means upstream-distinct: if `sources/2.md` and
     `sources/5.md` both cite the same original paper, that's *one*
     citation chain, not two.
   - Note contradictions explicitly rather than silently picking one
     side. A report that says "sources disagree on X: [2] says A, [5]
     says B" is more useful than one that quietly picks A.
   - If a sub-question has only single-source coverage, run one more
     targeted search before falling back to the `[single-source]`
     label.

6. **Synthesize** - `write` `/tmp/research-<timestamp>/report.md`
   containing:
   - **Executive summary** (3-5 bullets) - what the user actually
     wanted to know, up top.
   - **One section per sub-question** with inline numeric citations
     (`[1]`, `[2]`, ...). Cite at the claim, not at the paragraph
     end - the reader needs to know *which* sentence came from
     *which* source.
   - **Sources** list mapping numbers → URL → local source-file path
     (`sources/<n>.md`) → publication date when known.
   - **Limitations** section listing gaps, single-source claims,
     contradictions, and any sub-questions that were deferred. This
     is what separates a research report from a confident-sounding
     summary.

## Authentication

If a source is paywalled or login-gated, **surface it; don't try to
bypass**. Add a note to the source's `sources/<n>.md` and the
report's *Limitations* section. Offer to use `handle_authentication`
if the user supplies credentials for that specific source.

## Pitfalls

- **Search-engine rate-limiting / CAPTCHA**: if DuckDuckGo HTML
  returns a challenge page, surface it; don't loop. Suggest the user
  rerun with longer `wait_for_condition` gaps between queries, or
  swap the search backend per step 2.
- **JS-rendered articles**: always `wait_until: networkidle` and
  inspect the *rendered* DOM via `execute_script`, not the raw page
  source. Modern article sites server-render a skeleton and hydrate
  the body client-side; the source HTML can be missing the content
  entirely.
- **Hallucinated citations**: every claim in the report must trace
  to a file in `/tmp/research-<timestamp>/sources/`. If you can't
  trace it, drop the claim or label it as inference, not citation.
- **Echo-chamber sources**: N articles citing the same primary
  source ≠ N independent confirmations. Treat the upstream as one
  citation and note the downstream sources as restatements.
- **Scope creep**: stop at the planned sub-question set. If
  intriguing tangents come up, list them under *Limitations →
  deferred* rather than silently expanding scope - it keeps the
  report bounded and gives the user a clean follow-up surface.
- **Stale recency**: when the question is time-sensitive, prefer
  sources from the last 12 months and **surface publication dates in
  the report**. An undated citation on a "latest in 2026" question is
  a red flag.
- **Single-vendor bias on comparison questions**: if the user asks
  "A vs B" and your corpus is 80% A's blog posts, the report will be
  pro-A no matter how carefully you write it. The contrarian search
  in step 3 exists to prevent this; don't skip it.

## Out of scope (future work)

A dedicated `web_search` tool (Brave / Tavily / SerpAPI-backed) would
replace step 2 and remove the HTML-scraping-a-search-engine
dependency entirely. Until that exists, the DuckDuckGo HTML endpoint
is the default and step 2 documents how to swap it.
