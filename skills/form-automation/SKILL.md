---
name: form-automation
description: Use this when the user asks to complete a multi-step form, optionally behind a login. Orchestrates handle_authentication, navigate_to_url, fill_form, click_element, wait_for_condition, and take_screenshot to capture the post-submit confirmation.
tags:
  - forms
  - automation
  - workflow
---

# form-automation

Use this when the user wants to submit a form on their behalf -
contact forms, signup flows, intake questionnaires, multi-step
wizards. The output is the post-submit confirmation: a screenshot
plus any reference number or URL the form returned.

## When to use

Trigger this playbook for requests that name a specific form and the
data to put in it:

- "Submit this contact form with my email and message X"
- "Sign me up at `<url>` with name=Y email=Z"
- "Fill out the support ticket form and attach this description"
- "Walk through the 4-step onboarding with these answers"

Do **not** use this for testing whether a form *works* (use
`webapp-testing`) or extracting data *from* a form (use
`web-scraping`).

## Workflow

1. **Optional: authenticate** - if the form requires login, run
   `handle_authentication` first. The session carries cookies across
   subsequent calls.

2. **Navigate and inspect**
   - `navigate_to_url` with `wait_until: networkidle`.
   - `take_screenshot` so you can verify you're on the right form
     before you start typing.
   - If the form layout is unclear from the screenshot, use
     `execute_script` to enumerate field selectors:
     `Array.from(document.querySelectorAll('input,select,textarea')).map(e => ({name: e.name, type: e.type, id: e.id}))`.

3. **Map user-supplied data to fields** - the user gave you values;
   you need to match each to a selector. Prefer in this order:
   `name` attribute, `id`, accessible label (`aria-label` /
   associated `<label for=...>`), placeholder text. Build the
   `fields[]` array for `fill_form`.

4. **Fill** - one `fill_form` call per logical group (per page of a
   multi-step wizard). Use the correct field `type`
   (`text`/`select`/`checkbox`/`radio`/`file`) - filling a checkbox
   as text silently fails.

5. **For each step in a multi-step form**:
   - `click_element` on the "Next" / "Continue" button.
   - `wait_for_condition` (`condition: selector`, `state: visible`)
     on a known element of the next step. Don't assume it's a
     re-render of the same URL - some wizards navigate, some swap
     the DOM in place.
   - Repeat fill → next.

6. **Submit and capture**
   - `click_element` on the submit button.
   - `wait_for_condition` for the confirmation - the success URL, a
     thank-you message, a confirmation number selector. **Do not
     screenshot before this** - you'll capture a half-rendered page.
   - `take_screenshot` of the confirmation.
   - `extract_data` to read back any confirmation number, ticket ID,
     or redirected URL - that's what the user actually wants.

7. **Report back** - include the screenshot path and the extracted
   confirmation values. If submission was rejected (validation
   error), include the error message and a screenshot of the
   failing state.

## Sensitive data

Treat anything the user passed in (passwords, payment details,
personal info) as one-shot: use it for the fill, do not log it, do
not include it in the response. If the user pasted credentials in
plaintext, complete the task but mention in the response that they
should rotate the credential.

## Pitfalls

- **Required-field validation**: many forms reject a submit with no
  visible error if a required field is blank or has the wrong
  format. After click_element on submit, always check whether you
  got to the confirmation OR an error state - don't assume success
  from absence of crash.
- **CAPTCHA / bot detection**: if the form is gated by a CAPTCHA,
  surface it - the agent cannot solve them. The stealth_mode
  browser config helps with passive detection but not with
  interactive CAPTCHAs.
- **File uploads**: `fill_form` with `type: file` expects a local
  path readable by the browser process. If the user passed a URL,
  download it first with the `fetch` tool - set `save_path` to a
  location under `/tmp/playwright/artifacts/` and pass that path as
  the field value. Don't reach for `execute_script` to invoke the
  browser-side `fetch()` API for this; the tool-level `fetch` is
  simpler and writes the bytes directly to disk.
- **Hidden steps**: some wizards inject extra confirmation steps
  ("Are you sure?") that aren't in the user's instructions. Treat
  them as part of the flow and proceed if the answer is obvious;
  ask the user otherwise.
