# Changelog

All notable changes to this project will be documented in this file.

## [0.4.12](https://github.com/inference-gateway/browser-agent/compare/v0.4.11...v0.4.12) (2025-12-14)

### ♻️ Improvements

* Improve VNC container X11 connection logic ([ed106eb](https://github.com/inference-gateway/browser-agent/commit/ed106ebafc9bdc053b49b31fa4bc0ea0e5089e38))

### 🐛 Bug Fixes

* Xvfb not listening over tcp ([11b8827](https://github.com/inference-gateway/browser-agent/commit/11b882723dbaf6780d605abe66b1fe773798177d))

### 🔧 Miscellaneous

* Update config and add shortcut files ([57afcd1](https://github.com/inference-gateway/browser-agent/commit/57afcd1991d68e5d57037ca1e7da5c053c9063df))

## [0.4.11](https://github.com/inference-gateway/browser-agent/compare/v0.4.10...v0.4.11) (2025-12-13)

### 🔧 Miscellaneous

* **deps:** Bump ADK to version 0.16.2 ([4fe2d19](https://github.com/inference-gateway/browser-agent/commit/4fe2d19fb7089caa07a4e1d8e84dfdf0e694e6de))

## [0.4.10](https://github.com/inference-gateway/browser-agent/compare/v0.4.9...v0.4.10) (2025-12-13)

### 🔧 Miscellaneous

* **deps:** Bump ADK to 0.16.1 ([c9f6046](https://github.com/inference-gateway/browser-agent/commit/c9f6046a13f96a6c7df09b6f56f4ed31ab10d23f))

## [0.4.9](https://github.com/inference-gateway/browser-agent/compare/v0.4.8...v0.4.9) (2025-12-13)

### 🐛 Bug Fixes

* Sync scripts should be wrapped ([0699131](https://github.com/inference-gateway/browser-agent/commit/0699131ddb535b2c24392140d4f4d25514aeffad))

## [0.4.8](https://github.com/inference-gateway/browser-agent/compare/v0.4.7...v0.4.8) (2025-12-12)

### 🔧 Miscellaneous

* **deps:** Bump adl-cli to v0.26.2 ([5817930](https://github.com/inference-gateway/browser-agent/commit/58179308feb7faedbded18ddc18a7d95a7d8ad38))
* **deps:** Run flox activate ([5dd0ede](https://github.com/inference-gateway/browser-agent/commit/5dd0ede607bd03a62b4232339fb0fa9a4d7ff9ed))

## [0.4.7](https://github.com/inference-gateway/browser-agent/compare/v0.4.6...v0.4.7) (2025-11-25)

### 🐛 Bug Fixes

* Improve Playwright installation order in Dockerfile ([#47](https://github.com/inference-gateway/browser-agent/issues/47)) ([088ac1e](https://github.com/inference-gateway/browser-agent/commit/088ac1e0a867ab4a32526477116653470c81b81e))

## [0.4.6](https://github.com/inference-gateway/browser-agent/compare/v0.4.5...v0.4.6) (2025-11-25)

### 🐛 Bug Fixes

* Ensure proper cache directory ownership in Dockerfile ([#46](https://github.com/inference-gateway/browser-agent/issues/46)) ([976c15c](https://github.com/inference-gateway/browser-agent/commit/976c15c0273ea785546def99a75fe3e65ee251d9))

## [0.4.5](https://github.com/inference-gateway/browser-agent/compare/v0.4.4...v0.4.5) (2025-11-25)

### 🐛 Bug Fixes

* Move browser cache to user directory ([#44](https://github.com/inference-gateway/browser-agent/issues/44)) ([acb379b](https://github.com/inference-gateway/browser-agent/commit/acb379bac6cbba9be6bb31d0787e7765bb544d56))

### 🔧 Miscellaneous

* Add agents config for testing agents without docker-compose ([#45](https://github.com/inference-gateway/browser-agent/issues/45)) ([927cccc](https://github.com/inference-gateway/browser-agent/commit/927cccc5282282dd9f4b1f5eefaca82356fd5701))

## [0.4.4](https://github.com/inference-gateway/browser-agent/compare/v0.4.3...v0.4.4) (2025-11-24)

### 🔧 Miscellaneous

* **deps:** Update to ADL CLI v0.26.0 and dependency versions ([#43](https://github.com/inference-gateway/browser-agent/issues/43)) ([cfd9437](https://github.com/inference-gateway/browser-agent/commit/cfd943750b5e087b92b415b10d9e35373fe358fc))

## [0.4.3](https://github.com/inference-gateway/browser-agent/compare/v0.4.2...v0.4.3) (2025-10-20)

### ♻️ Improvements

* **tests:** Remove redundant comments in session isolation tests ([1c54151](https://github.com/inference-gateway/browser-agent/commit/1c54151b4e279e08d3018d2e14180ad6cafeec2a))

### 🔧 Miscellaneous

* **deps:** Update ADL CLI version to 0.23.11 in generated files ([c7dbb1c](https://github.com/inference-gateway/browser-agent/commit/c7dbb1cc85f56ec4b41f403f191f5821112f9b02))

## [0.4.2](https://github.com/inference-gateway/browser-agent/compare/v0.4.1...v0.4.2) (2025-10-19)

### ♻️ Improvements

* Implement multi-tenant browser session isolation ([#41](https://github.com/inference-gateway/browser-agent/issues/41)) ([5661bde](https://github.com/inference-gateway/browser-agent/commit/5661bdebbb1b9e069d97de3cc4d51780ee2fe56a)), closes [#40](https://github.com/inference-gateway/browser-agent/issues/40)
* Improve the configurations ([#38](https://github.com/inference-gateway/browser-agent/issues/38)) ([9579694](https://github.com/inference-gateway/browser-agent/commit/9579694e5b877fad9ce394765ff95ab84c5af3f7))

## [0.4.1](https://github.com/inference-gateway/browser-agent/compare/v0.4.0...v0.4.1) (2025-10-18)

### ♻️ Improvements

* **logs:** Add browser configuration logging at service initialization ([7faea50](https://github.com/inference-gateway/browser-agent/commit/7faea50267018038e26ee6f2a0e8af7a645af889))
* **manifest:** Remove duplicate Go package entries and consolidate versions ([f95fd38](https://github.com/inference-gateway/browser-agent/commit/f95fd38f2829caddebd4c819b1a73260e36836e4))

## [0.4.0](https://github.com/inference-gateway/browser-agent/compare/v0.3.4...v0.4.0) (2025-10-18)

### ✨ Features

* Add headless configuration ([#37](https://github.com/inference-gateway/browser-agent/issues/37)) ([c08b4fb](https://github.com/inference-gateway/browser-agent/commit/c08b4fb4d2a7a42c22000037811f3e0dcf6b1509))

## [0.3.4](https://github.com/inference-gateway/browser-agent/compare/v0.3.3...v0.3.4) (2025-10-17)

### 🔧 Miscellaneous

* **deps:** Update ADL CLI version to 0.23.7 in generated files ([a4e520e](https://github.com/inference-gateway/browser-agent/commit/a4e520ef4823e58114d8c994a0061535c1a26177))

## [0.3.3](https://github.com/inference-gateway/browser-agent/compare/v0.3.2...v0.3.3) (2025-10-14)

### ♻️ Improvements

* Remove write_to_csv skill and related tests ([#36](https://github.com/inference-gateway/browser-agent/issues/36)) ([144d901](https://github.com/inference-gateway/browser-agent/commit/144d901cf409200da76b4e90baae074b0c1cbf57))

## [0.3.2](https://github.com/inference-gateway/browser-agent/compare/v0.3.1...v0.3.2) (2025-10-12)

### ♻️ Improvements

* Remove duplicate Go package entries in manifest.lock ([e6aa1cd](https://github.com/inference-gateway/browser-agent/commit/e6aa1cdeba6f1d7aa446f3f2548b2a7f366d39b7))

### 📚 Documentation

* Update README to include instructions for collecting prices and writing to CSV ([132ba44](https://github.com/inference-gateway/browser-agent/commit/132ba442568a4ccbe47e85e189158b5ed85a6054))

### 🔧 Miscellaneous

* **deps:** Update ADL CLI version to 0.23.1 in generated files ([f72060f](https://github.com/inference-gateway/browser-agent/commit/f72060f1204f27d80fffc16cfa5fb6c6e0219bb0))
* Update dependencies and generated files to ADL CLI v0.23.2 ([9ae2bd2](https://github.com/inference-gateway/browser-agent/commit/9ae2bd207e87f4dd3b5e7b14bd6a93d580fa5891))

## [0.3.1](https://github.com/inference-gateway/browser-agent/compare/v0.3.0...v0.3.1) (2025-10-06)

### ♻️ Improvements

* Improve screenshot and CSV writing skills with artifact integration ([#35](https://github.com/inference-gateway/browser-agent/issues/35)) ([c95f964](https://github.com/inference-gateway/browser-agent/commit/c95f964906ff5ab87e2c6b38958bdabf9043a566))

## [0.3.0](https://github.com/inference-gateway/browser-agent/compare/v0.2.1...v0.3.0) (2025-10-01)

### ✨ Features

* Update ADL CLI version to 0.22.1 and add artifacts configuration ([#34](https://github.com/inference-gateway/browser-agent/issues/34)) ([19c7ff0](https://github.com/inference-gateway/browser-agent/commit/19c7ff06ec28371f16d8e0db2eada1570efd4a7e))

## [0.2.1](https://github.com/inference-gateway/browser-agent/compare/v0.2.0...v0.2.1) (2025-09-26)

### ♻️ Improvements

* Bump ADL-CLI to 0.21.7 and ADK version to 0.11.1 ([81ee9bf](https://github.com/inference-gateway/browser-agent/commit/81ee9bf69c5aaa7f2c3a917cf0e9ebad722ed75b))

## [0.2.0](https://github.com/inference-gateway/browser-agent/compare/v0.1.3...v0.2.0) (2025-09-26)

### ✨ Features

* **skills:** Add write_to_csv skill for data export workflows ([#25](https://github.com/inference-gateway/browser-agent/issues/25)) ([5b7509f](https://github.com/inference-gateway/browser-agent/commit/5b7509f3bf96d3f5e6f17c54497e35f4c88aebec)), closes [#24](https://github.com/inference-gateway/browser-agent/issues/24)

### ♻️ Improvements

* Update agent metadata to use agent-card.json and increment ADL CLI version to 0.21.6 ([#33](https://github.com/inference-gateway/browser-agent/issues/33)) ([7d91dde](https://github.com/inference-gateway/browser-agent/commit/7d91dde1d272a11c48db71d77d430489926f45af))

## [0.1.3](https://github.com/inference-gateway/browser-agent/compare/v0.1.2...v0.1.3) (2025-09-22)

### ♻️ Improvements

* Rename playwright-agent to browser-agent across the project ([d05fb16](https://github.com/inference-gateway/browser-agent/commit/d05fb1686bf90ad5d6b0c13f4154849034b57a17))

## [0.1.2](https://github.com/inference-gateway/playwright-agent/compare/v0.1.1...v0.1.2) (2025-09-22)

### 🐛 Bug Fixes

* **container:** Use Dockerfile.playwright for releases ([#27](https://github.com/inference-gateway/playwright-agent/issues/27)) ([cbfca28](https://github.com/inference-gateway/playwright-agent/commit/cbfca28cb5dec69c4357938a2428f0fa926216b9)), closes [#26](https://github.com/inference-gateway/playwright-agent/issues/26)

### 🔧 Miscellaneous

* Update ADL CLI version to 0.21.4 in generated files ([d4ea0a3](https://github.com/inference-gateway/playwright-agent/commit/d4ea0a3b18f2a31f2c1f422defd3e334b956d365))

## [0.1.1](https://github.com/inference-gateway/playwright-agent/compare/v0.1.0...v0.1.1) (2025-09-21)

### ♻️ Improvements

* **docker-compose:** Disable various infer tools in configuration ([e4a5c17](https://github.com/inference-gateway/playwright-agent/commit/e4a5c179b212c11d2291218dc229c98411c45411))
