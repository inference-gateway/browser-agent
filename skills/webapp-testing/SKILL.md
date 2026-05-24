---
name: webapp-testing
description: Use this when the user asks to verify, validate, or test a webapp end-to-end. Performs reconnaissance-then-action: navigate, screenshot the rendered DOM, identify selectors, then exercise the flow using navigate_to_url, click_element, fill_form, wait_for_condition, and take_screenshot.
tags:
  - testing
  - qa
  - e2e
  - playwright
---

# webapp-testing

Use this when the user asks to verify, validate, or test a webapp's
behavior end-to-end - login flows, checkout funnels, form submissions,
multi-step wizards, post-deploy smoke tests.

## When to use

Trigger this playbook for requests that contain action verbs against a
live URL paired with an expected outcome. Examples:

- "Make sure the login on `https://app.example.com` still works"
- "Verify the signup form rejects passwords under 8 characters"
- "Smoke-test the checkout flow and screenshot every step"
- "After the deploy, confirm the dashboard loads and shows charts"

Do **not** use this for read-only data extraction (use `web-scraping`)
or for filling a single form (use `form-automation`).

## Workflow: reconnaissance-then-action

The most common failure mode in webapp testing is acting on stale
assumptions about the DOM. Always inspect the rendered page before
choosing selectors.

1. **Reconnaissance**
   - **Pre-flight probe** (post-deploy smoke tests only): before
     opening the browser, `fetch GET /health` (or whatever the app's
     status endpoint is) and confirm a 2xx response. If the backend
     is already down, surface the failure and stop - exercising the
     UI tells you nothing the probe didn't already.
   - `navigate_to_url` with `wait_until: networkidle` so client-side
     rendering completes before you look at the DOM.
   - `take_screenshot` with `full_page: true` and save the path - you
     will reference it when reporting back to the user.
   - If the page is JS-heavy and selectors are unclear, use
     `execute_script` to dump the relevant section of the DOM
     (`document.querySelector('main').outerHTML`) and pick selectors
     from the actual rendered output, not the page source.

2. **Identify selectors** from the screenshot + DOM dump. Prefer in
   this order: `data-testid`/`data-test` attributes, ARIA role +
   accessible name, stable IDs, text content. Avoid nth-child
   selectors and deeply nested CSS - they break under the lightest
   refactor.

3. **Action sequence** - for each step in the user's flow:
   - `wait_for_condition` (`condition: selector`, `state: visible`)
     **before** every interaction. Skipping this is the #1 source of
     flaky failures.
   - `click_element` / `fill_form` to perform the action.
   - `wait_for_condition` (`condition: networkidle` or a result
     selector) to confirm the action completed.

4. **Validation**
   - `take_screenshot` of the expected end state.
   - `extract_data` to read back values that prove the flow worked
     (confirmation message, order ID, redirected URL, etc.).
   - Use `execute_script` to assert on client-side state when no
     visible artifact exists (e.g. `return localStorage.getItem('token') !== null`).
   - **API-side assertion**: when a UI action is supposed to produce
     server-side state (an order, a record, a job), `fetch` the
     corresponding read endpoint and confirm the resource exists
     with the expected fields. A "success" toast does not prove the
     row was written - the API does.

5. **Report back** - include the screenshot paths and the extracted
   confirmation values. If any step failed, include the screenshot
   taken just before the failing assertion.

## Authentication

If the flow requires login, run `handle_authentication` once at the
start, then proceed with the rest of the workflow. Don't re-login
between steps - the browser session persists across tool calls within
the same task.

## Pitfalls

- **Acting before networkidle**: clicking a button that hasn't been
  hydrated yet by the frontend framework silently no-ops. Always wait
  for `networkidle` after the initial navigate.
- **Asserting on DOM that hasn't updated**: after a click that triggers
  a re-render, wait for the new selector before extracting, not
  immediately.
- **Hardcoding timeouts above 30s**: if a page takes longer than that,
  the app has a real problem - report it rather than papering over it.
- **Iframes**: most selectors do not cross iframe boundaries. If the
  target is inside one, surface that to the user; the underlying
  Playwright service has iframe-aware lookups but they require
  explicit handling.
