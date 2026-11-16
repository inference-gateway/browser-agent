package stealth

import _ "embed"

// StealthJS is the browser fingerprint evasion script injected into every page.
// Vendored from github.com/jonfriesen/playwright-go-stealth v0.0.3 (MIT, see LICENSE)
// to avoid pulling its legacy playwright-go dependency and the unpatched go-jose/v3.
//
//go:embed stealth.min.js
var StealthJS string
