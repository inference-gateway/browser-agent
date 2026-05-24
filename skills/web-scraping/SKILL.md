---
name: web-scraping
description: Use this when the user asks to extract structured data from one or more pages. Drives extract_data across paginated URLs, normalizes results, and writes a JSON/CSV artifact via the write tool.
tags:
  - scraping
  - extraction
  - data
---

# web-scraping

Use this when the user asks to pull structured data off one or more
web pages and hand it back in a usable form (JSON, CSV, list of
records).

## When to use

Trigger this playbook for read-only data-collection requests:

- "Grab all product names and prices from this catalogue"
- "Get the titles and links from the first 5 pages of this blog"
- "Extract the table of contributors from this GitHub page as CSV"
- "Pull every job posting that matches `senior engineer` and save them"

Do **not** use this for tasks that mutate page state (use
`webapp-testing` or `form-automation`).

## Workflow

0. **Skip the browser when you can** - before opening a Playwright
   session, check if the data is reachable without one:
   - Does the site expose a JSON/XML API? Many SPAs render from a
     backend endpoint the page calls; if you can identify it (e.g.
     from the user's URL pattern, a `/api/` path, or a probe), `fetch`
     it directly and skip the rest of this workflow.
   - Is there a `/sitemap.xml`? `fetch` it to enumerate URLs instead
     of clicking through pagination.
   - Is the target a static page (RFC, raw GitHub README, plaintext
     docs)? `fetch` returns the body directly; no DOM rendering needed.
   - Are the records the user wants linked as downloadable files (CSV
     exports, PDFs)? `fetch` each with `save_path` and parse offline.

   Only fall through to the Playwright workflow below when the data is
   exclusive to the rendered DOM.

1. **Reconnaissance** - on the first page only:
   - `navigate_to_url` with `wait_until: networkidle`.
   - `take_screenshot` so you can verify the layout matches what the
     user described.
   - Use `execute_script` to inspect the structure of a single record
     (e.g. `document.querySelectorAll('.product')[0].outerHTML`).
     Pick the smallest stable selector that uniquely identifies each
     record and its inner fields.

2. **Define the extractor schema** - decide the field set up front and
   keep it consistent across all pages. Each `extract_data` call
   should use the same `extractors[]` shape with `multiple: true` so
   you get an array of records per page.

3. **Paginate**
   - Detect the pagination strategy from the screenshot/DOM:
     - **URL-based**: increment a `?page=N` parameter - loop with
       `navigate_to_url` + `extract_data`.
     - **Click-based**: `click_element` on the "next" button, then
       `wait_for_condition` with `networkidle` before the next
       `extract_data`.
     - **Infinite scroll**: `execute_script` to scroll
       (`window.scrollTo(0, document.body.scrollHeight)`), then
       `wait_for_condition` with `networkidle`.
   - **Stop conditions** to check on every page: empty results,
     duplicate first-record (looped), pagination button disabled.
   - **Respect rate limits**: insert a `wait_for_condition` with
     `condition: timeout` between pages (1-3 seconds). Don't hammer.

4. **Normalize**
   - Strip whitespace, decode HTML entities, coerce numeric fields.
     `execute_script` is fine for one-off cleanup that's hard to
     express in CSS selectors.
   - Deduplicate by a stable key (URL, ID) when paginating - the same
     record sometimes appears on adjacent pages.

5. **Persist** - write the result to disk so the user can download it:
   - JSON: `write` to `/tmp/scrape-<timestamp>.json` with the array of
     records.
   - CSV: prefer `format: csv` in `extract_data` when the schema is
     flat; otherwise serialize yourself and `write`.
   - For large scrapes, write incrementally (per page) so a crash
     mid-run doesn't lose everything.
   - **Binary assets**: if the scrape produces non-HTML artifacts
     (PDFs, images, CSV exports linked from the page), use `fetch`
     with `save_path` to download each one straight into the artifact
     directory rather than going through `execute_script`.

6. **Report back** - tell the user how many records, which file, and
   show a 3-5 row sample inline. Mention any pages that returned
   zero results or failed.

## Authentication

If the data is behind a login, run `handle_authentication` once
before step 1. Cookies persist across navigation within the session.

## Pitfalls

- **Scraping the wrong DOM**: if a page is server-rendered for SEO but
  hydrates on the client, the source HTML and the rendered DOM can
  differ. Always `wait_until: networkidle` and inspect the rendered
  state.
- **`extract_data` returning fewer results than expected**: usually a
  selector that only matches the first card because of accidental
  specificity. Use `execute_script` to count
  `document.querySelectorAll(...)` and confirm the selector matches
  every record.
- **Silent pagination loops**: a "next" button that's disabled at the
  end may still be clickable in the DOM. Check the button's
  `disabled` attribute or detect a duplicate first record.
- **robots.txt and terms of service**: scraping is the user's
  responsibility - surface concerns when the target is a site with
  obvious scraping policies, but don't refuse.
