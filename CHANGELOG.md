# Changelog

All notable changes to this project will be documented in this file.

## [0.7.3](https://github.com/inference-gateway/browser-agent/compare/v0.7.2...v0.7.3) (2026-07-23)

### ♻️ Improvements

* un-fork main.go and .gitattributes from .adl-ignore ([#135](https://github.com/inference-gateway/browser-agent/issues/135)) ([fa94901](https://github.com/inference-gateway/browser-agent/commit/fa9490123d503ad9aa8e30580bca5c52993222a0)), references [#133](https://github.com/inference-gateway/browser-agent/issues/133) [#128](https://github.com/inference-gateway/browser-agent/issues/128) [#134](https://github.com/inference-gateway/browser-agent/issues/134)

### 🐛 Bug Fixes

* migrate to mxschmitt/playwright-go fork to fix retired driver CDN ([#134](https://github.com/inference-gateway/browser-agent/issues/134)) ([8250db7](https://github.com/inference-gateway/browser-agent/commit/8250db7dca8281c438e51840c87b899de1b9aec5))

### 🔧 Miscellaneous

* **flox:** add lockfile ([c5d6523](https://github.com/inference-gateway/browser-agent/commit/c5d652379a7cdc614df1ce72e545080f3f1c8cd2))

## [0.7.2](https://github.com/inference-gateway/browser-agent/compare/v0.7.1...v0.7.2) (2026-07-23)

### 🐛 Bug Fixes

* copy and load skills from .agents/skills ([#133](https://github.com/inference-gateway/browser-agent/issues/133)) ([aec6515](https://github.com/inference-gateway/browser-agent/commit/aec6515fafbf8d886c74e0e81d8283d070805b21)), references [#308](https://github.com/inference-gateway/browser-agent/issues/308)

## [0.7.1](https://github.com/inference-gateway/browser-agent/compare/v0.7.0...v0.7.1) (2026-07-23)

### 🔧 Miscellaneous

* **adl:** refresh agent.yaml defaults from ADL CLI v0.54.0 ([#131](https://github.com/inference-gateway/browser-agent/issues/131)) ([31d4b96](https://github.com/inference-gateway/browser-agent/commit/31d4b964f634c836f53ce51ac9dea840363b55ba))
* **deps:** bump ADL CLI to v0.54.0 ([#132](https://github.com/inference-gateway/browser-agent/issues/132)) ([0d6f2f9](https://github.com/inference-gateway/browser-agent/commit/0d6f2f9889469d17aeeb4d71ce64d94a54ebdd19))

## [0.7.0](https://github.com/inference-gateway/browser-agent/compare/v0.6.4...v0.7.0) (2026-07-17)

### ✨ Features

* **telemetry:** add OpenTelemetry support via agent.yaml manifest ([#128](https://github.com/inference-gateway/browser-agent/issues/128)) ([f8b3e99](https://github.com/inference-gateway/browser-agent/commit/f8b3e9984c21f7e826fa9b6fa06ee14cc5fbd917)), closes [#105](https://github.com/inference-gateway/browser-agent/issues/105)

### 🐛 Bug Fixes

* **playwright:** add return after t.Fatal to fix SA5011 nil pointer dereference lint ([#108](https://github.com/inference-gateway/browser-agent/issues/108)) ([dc1369f](https://github.com/inference-gateway/browser-agent/commit/dc1369fb3a8f4081bcf3d2dc7a1553b349bd665d))

### 👷 CI

* **claude:** change effort to max ([a3d814f](https://github.com/inference-gateway/browser-agent/commit/a3d814fb643ef254ccc382e5732e5b07fba75e96))
* **claude:** remove system prompt - use default community maintained prompt ([54e9e0f](https://github.com/inference-gateway/browser-agent/commit/54e9e0ffde28bbaa7d7c960e5bf8cf25c047946a))
* **claude:** standardize workflow + task-based branch prefix ([a30166d](https://github.com/inference-gateway/browser-agent/commit/a30166ddc69264e4d04640b61c414d57c4e89ea0))
* **deps:** bump actions/checkout from 6.0.3 to 7.0.0 in the github-actions group ([#92](https://github.com/inference-gateway/browser-agent/issues/92)) ([408414c](https://github.com/inference-gateway/browser-agent/commit/408414c7135cdb662acbbead61da2cc3d5546d21))
* **deps:** bump anthropics/claude-code-action from 1.0.135 to 1.0.144 in the github-actions group ([#88](https://github.com/inference-gateway/browser-agent/issues/88)) ([12974f0](https://github.com/inference-gateway/browser-agent/commit/12974f0941200f4773c4b1cbc0eff958da73201b))
* **deps:** bump anthropics/claude-code-action from 1.0.144 to 1.0.151 in the github-actions group ([#90](https://github.com/inference-gateway/browser-agent/issues/90)) ([bde6b2d](https://github.com/inference-gateway/browser-agent/commit/bde6b2d265e77a9243e134ad04362053b6d1f834))
* **deps:** bump anthropics/claude-code-action from 1.0.173 to 1.0.174 in the github-actions group ([#110](https://github.com/inference-gateway/browser-agent/issues/110)) ([c48820e](https://github.com/inference-gateway/browser-agent/commit/c48820e7336b4ac955c24532bd08dce6c1fd65a6))
* **deps:** bump anthropics/claude-code-action from 1.0.173 to 1.0.174 in the github-actions group ([#112](https://github.com/inference-gateway/browser-agent/issues/112)) ([8ce7111](https://github.com/inference-gateway/browser-agent/commit/8ce71119955db5b497bd6eecdd1efc24e362d006))
* **deps:** bump anthropics/claude-code-action from 1.0.174 to 1.0.175 in the github-actions group ([#123](https://github.com/inference-gateway/browser-agent/issues/123)) ([32c9fde](https://github.com/inference-gateway/browser-agent/commit/32c9fde8997d832fa351b0e5c475cc8a70ec1e57))
* **deps:** bump github.com/inference-gateway/adk from 0.19.0 to 0.20.0 in the gomod group ([#96](https://github.com/inference-gateway/browser-agent/issues/96)) ([c9d6c3c](https://github.com/inference-gateway/browser-agent/commit/c9d6c3c860acb3ca2ecba752ae3f8befcab55e09))
* **deps:** bump github.com/inference-gateway/adk from 0.23.0 to 0.23.2 in the gomod group ([#114](https://github.com/inference-gateway/browser-agent/issues/114)) ([51b540b](https://github.com/inference-gateway/browser-agent/commit/51b540bbb6105eb3b85fe2c998cd7916783875f6))
* **deps:** bump github.com/inference-gateway/adk from 0.23.2 to 0.23.3 in the gomod group ([#124](https://github.com/inference-gateway/browser-agent/issues/124)) ([f99c629](https://github.com/inference-gateway/browser-agent/commit/f99c6290de0074b6915d1d27269d17a6a076268f))
* **deps:** bump github.com/quic-go/quic-go from 0.59.0 to 0.59.1 ([#89](https://github.com/inference-gateway/browser-agent/issues/89)) ([5fbed6c](https://github.com/inference-gateway/browser-agent/commit/5fbed6cf9daf8ac85f513966642d24cfdbdd5754))
* **deps:** bump github.com/sethvargo/go-envconfig from 1.3.1 to 1.4.0 in the gomod group across 1 directory ([#130](https://github.com/inference-gateway/browser-agent/issues/130)) ([e25d875](https://github.com/inference-gateway/browser-agent/commit/e25d8751f451f063883c97b6473b8c222a8fde22))
* **deps:** bump Go version from 1.26.2 to 1.26.4 and update package entries ([9c3ab60](https://github.com/inference-gateway/browser-agent/commit/9c3ab609970494b226c1155e51b85537e3b9ca64))
* **deps:** bump Go version from 1.26.2 to 1.26.4 in agent configuration ([59dea40](https://github.com/inference-gateway/browser-agent/commit/59dea40ebf453e3bc56653a642fffc9585b1aa4d))
* **deps:** bump golang from 1.26.2-alpine to 1.26.4-alpine in the docker group ([#104](https://github.com/inference-gateway/browser-agent/issues/104)) ([b0c83e6](https://github.com/inference-gateway/browser-agent/commit/b0c83e6941d835d4d1fb6890c5136d34a406f647))
* **deps:** bump golang.org/x/crypto from 0.51.0 to 0.52.0 ([#98](https://github.com/inference-gateway/browser-agent/issues/98)) ([3319f33](https://github.com/inference-gateway/browser-agent/commit/3319f333c1c8c48f6e0e7586e35ddfe2d680e4e2))
* **deps:** bump golang.org/x/net from 0.52.0 to 0.55.0 ([#95](https://github.com/inference-gateway/browser-agent/issues/95)) ([f0fc48a](https://github.com/inference-gateway/browser-agent/commit/f0fc48a7932421f20bc03b32dab6531cbb14e5d4))
* **deps:** bump inference-gateway/infer-action from 0.32.1 to 0.32.2 in the github-actions group ([#117](https://github.com/inference-gateway/browser-agent/issues/117)) ([7cd5a17](https://github.com/inference-gateway/browser-agent/commit/7cd5a175051f5ec4b945cbdb247bee56cb47044b))
* **deps:** bump inference-gateway/infer-action from 0.32.2 to 0.34.1 in the github-actions group ([#127](https://github.com/inference-gateway/browser-agent/issues/127)) ([1df1161](https://github.com/inference-gateway/browser-agent/commit/1df1161caac0ea91ce29ef071755dab72b4b6561))
* **deps:** bump the github-actions group with 2 updates ([#125](https://github.com/inference-gateway/browser-agent/issues/125)) ([c404e78](https://github.com/inference-gateway/browser-agent/commit/c404e781856ccf595a1a50f16d3a3cd536b89952))
* **deps:** bump the github-actions group with 2 updates ([#93](https://github.com/inference-gateway/browser-agent/issues/93)) ([374b5e3](https://github.com/inference-gateway/browser-agent/commit/374b5e3fa453d946e9b0f2e6c0571962c4120da0))
* **deps:** bump the github-actions group with 2 updates ([#97](https://github.com/inference-gateway/browser-agent/issues/97)) ([29855af](https://github.com/inference-gateway/browser-agent/commit/29855afceea9ec09e95a67e855d59d564a2729e2))
* **deps:** bump the github-actions group with 6 updates ([#94](https://github.com/inference-gateway/browser-agent/issues/94)) ([180aea3](https://github.com/inference-gateway/browser-agent/commit/180aea33c7013a7e938964620f184103610ce3e0))
* **deps:** downgrade task version from 3.51.1 to 3.48.0 in workflows and manifest ([be04608](https://github.com/inference-gateway/browser-agent/commit/be04608062589ea7e48b7776ab783e27418e19a9))
* **release:** update semantic release and plugins to latest versions with local installation ([cbccbe4](https://github.com/inference-gateway/browser-agent/commit/cbccbe401160ee84c12edd7b34b9901de0595e7b))

### 📚 Documentation

* author spec.documentation and spec.examples in agent.yaml ([#119](https://github.com/inference-gateway/browser-agent/issues/119)) ([2b48375](https://github.com/inference-gateway/browser-agent/commit/2b4837506054589ee31d94d6f23360b6106bab78)), closes [#118](https://github.com/inference-gateway/browser-agent/issues/118)

### 🔧 Miscellaneous

* **adl:** refresh agent.yaml defaults from ADL CLI v0.50.2 ([#120](https://github.com/inference-gateway/browser-agent/issues/120)) ([8bcd04d](https://github.com/inference-gateway/browser-agent/commit/8bcd04d6a3a551746c3815c7c27eb6d40fdba2c1))
* **deps:** bump ADL CLI to v0.40.0 ([#85](https://github.com/inference-gateway/browser-agent/issues/85)) ([7b334a5](https://github.com/inference-gateway/browser-agent/commit/7b334a5aad19347ae874b062f24b9fd3340ef887))
* **deps:** bump ADL CLI to v0.43.2 ([#87](https://github.com/inference-gateway/browser-agent/issues/87)) ([f34c9e2](https://github.com/inference-gateway/browser-agent/commit/f34c9e2d362813d3818372a6339a64b3041350b6))
* **deps:** bump ADL CLI to v0.44.0 ([#91](https://github.com/inference-gateway/browser-agent/issues/91)) ([f172652](https://github.com/inference-gateway/browser-agent/commit/f1726526a02cb184bce9c6ef7c7ce7cca22e825f))
* **deps:** bump ADL CLI to v0.46.0 ([#99](https://github.com/inference-gateway/browser-agent/issues/99)) ([b93f453](https://github.com/inference-gateway/browser-agent/commit/b93f453d3babcf2b4e1754d5c09e603356445e8d))
* **deps:** bump ADL CLI to v0.46.5 ([#103](https://github.com/inference-gateway/browser-agent/issues/103)) ([ca3537d](https://github.com/inference-gateway/browser-agent/commit/ca3537d442d364f75ae5e514edd8755f8028a033))
* **deps:** bump ADL CLI to v0.47.1 ([#107](https://github.com/inference-gateway/browser-agent/issues/107)) ([c9e470d](https://github.com/inference-gateway/browser-agent/commit/c9e470d7cecec61c8cbcade62982de717af1ebf1))
* **deps:** bump ADL CLI to v0.48.0 ([#109](https://github.com/inference-gateway/browser-agent/issues/109)) ([c776682](https://github.com/inference-gateway/browser-agent/commit/c7766823fffd0a0624c38e47bcebe59843250c25))
* **deps:** bump ADL CLI to v0.48.1 ([#111](https://github.com/inference-gateway/browser-agent/issues/111)) ([c8c3e16](https://github.com/inference-gateway/browser-agent/commit/c8c3e160460ee4d100f51c6f63145f2a65fa4ae8))
* **deps:** bump ADL CLI to v0.48.4 ([#113](https://github.com/inference-gateway/browser-agent/issues/113)) ([7b0c04a](https://github.com/inference-gateway/browser-agent/commit/7b0c04a149d8516dacd6a72807741d0a8d87ed62))
* **deps:** bump ADL CLI to v0.48.5 ([#115](https://github.com/inference-gateway/browser-agent/issues/115)) ([c57ce47](https://github.com/inference-gateway/browser-agent/commit/c57ce47519be212d36b4a066f84778a1d7fb1dbe))
* **deps:** bump ADL CLI to v0.49.0 ([#116](https://github.com/inference-gateway/browser-agent/issues/116)) ([f3fff87](https://github.com/inference-gateway/browser-agent/commit/f3fff871a132ad9e2f059c2f8349ffd0f0a5262d))
* **deps:** bump ADL CLI to v0.50.2 ([#121](https://github.com/inference-gateway/browser-agent/issues/121)) ([dafba5b](https://github.com/inference-gateway/browser-agent/commit/dafba5b569d832730062f5f893a6905b3fb47490))
* **deps:** bump ADL CLI to v0.51.0 ([#122](https://github.com/inference-gateway/browser-agent/issues/122)) ([b91902c](https://github.com/inference-gateway/browser-agent/commit/b91902c66f9b5a4f0b4da74026b7fb8efaa16a87))
* **deps:** bump ADL CLI to v0.51.4 ([#126](https://github.com/inference-gateway/browser-agent/issues/126)) ([70b5e63](https://github.com/inference-gateway/browser-agent/commit/70b5e63bc2d75d471443a806623d08ba33088030))
* **deps:** bump ADL CLI to v0.52.0 ([#129](https://github.com/inference-gateway/browser-agent/issues/129)) ([4b06b25](https://github.com/inference-gateway/browser-agent/commit/4b06b2554518cac7652f0e377d0e9076c500f740))
* **deps:** bump docker/setup-qemu-action version v4.0.0 -> v4.1.0 ([704487e](https://github.com/inference-gateway/browser-agent/commit/704487ea69c81f13c82075a527ca85746a08bbb2))
* **flox:** downgrade deps ([d47a169](https://github.com/inference-gateway/browser-agent/commit/d47a16925aa66eda4cc0d7516ce78aca94ba0a8f))
* remove obsolete configuration and shortcut files ([f4f03e7](https://github.com/inference-gateway/browser-agent/commit/f4f03e77f9ec90617038b382afc84588d0c50ac0))
* **schema:** update adl schema to latest ([d3f21bc](https://github.com/inference-gateway/browser-agent/commit/d3f21bc63d830bbbbe47b759f4cc6ba5fdc22b97))
* small fix ([fda71ce](https://github.com/inference-gateway/browser-agent/commit/fda71ce66f1534bed009455da9182ec82612f29f))

## [0.6.4](https://github.com/inference-gateway/browser-agent/compare/v0.6.3...v0.6.4) (2026-05-26)

### 🔧 Miscellaneous

* **deps:** Bump ADL CLI to v0.39.3 ([#84](https://github.com/inference-gateway/browser-agent/issues/84)) ([3ec3901](https://github.com/inference-gateway/browser-agent/commit/3ec39018def3c99a3b512ae6edb7475b48c1b9d9))

## [0.6.3](https://github.com/inference-gateway/browser-agent/compare/v0.6.2...v0.6.3) (2026-05-24)

### ♻️ Improvements

* **tools:** Correctness and quality hardening pass ([#83](https://github.com/inference-gateway/browser-agent/issues/83)) ([43c2e95](https://github.com/inference-gateway/browser-agent/commit/43c2e954e71e1397078677376635e987de7423d8))

## [0.6.2](https://github.com/inference-gateway/browser-agent/compare/v0.6.1...v0.6.2) (2026-05-24)

### 🐛 Bug Fixes

* **execute_script:** Actionable rejection reasons & fix function-expression false positive ([#80](https://github.com/inference-gateway/browser-agent/issues/80)) ([105a855](https://github.com/inference-gateway/browser-agent/commit/105a85520908b52d1c4a62600f89d7d2b887ce54))

### 👷 CI

* **deps:** Bump github.com/go-jose/go-jose/v3 from 3.0.4 to 3.0.5 ([#81](https://github.com/inference-gateway/browser-agent/issues/81)) ([64db80d](https://github.com/inference-gateway/browser-agent/commit/64db80d86918e271e236785b28ec2a4cd5432451))

### 🔧 Miscellaneous

* **deps:** Bump ADL CLI to v0.39.2 ([#82](https://github.com/inference-gateway/browser-agent/issues/82)) ([c059510](https://github.com/inference-gateway/browser-agent/commit/c0595100562fc59fef52c2cb5e978309f5de753d))
* Replace em dash with regular dash ([6683e28](https://github.com/inference-gateway/browser-agent/commit/6683e2816bd6c5a066af67cbde654e6f9d29d5be))
* Run task generate and add apache-2.0 license to skills ([9692bbe](https://github.com/inference-gateway/browser-agent/commit/9692bbe9d707fcc51bfe79581b4d94b0a4bf95e1))

## [0.6.1](https://github.com/inference-gateway/browser-agent/compare/v0.6.0...v0.6.1) (2026-05-24)

### 🔧 Miscellaneous

* **deps:** Bump ADL CLI to v0.39.0 ([#77](https://github.com/inference-gateway/browser-agent/issues/77)) ([0ee1bb6](https://github.com/inference-gateway/browser-agent/commit/0ee1bb6b4ed3b057e665dabd1bbd95a892303dc2))
* **deps:** Bump ADL CLI to v0.39.1 ([#78](https://github.com/inference-gateway/browser-agent/issues/78)) ([9568e8b](https://github.com/inference-gateway/browser-agent/commit/9568e8b8dd807553d8b57e4e6d9eba862ac52431))

## [0.6.0](https://github.com/inference-gateway/browser-agent/compare/v0.5.1...v0.6.0) (2026-05-24)

### ✨ Features

* Add artifact storage configuration ([1ad4511](https://github.com/inference-gateway/browser-agent/commit/1ad4511172452ba7f015fe0dcb5ee94b2e5c751c))
* **tools:** Enable fetch built-in with fetch-vs-browser guidance ([8466bec](https://github.com/inference-gateway/browser-agent/commit/8466bec61a0351f91ed98c0a64a0cfe43e32123a))

### 🐛 Bug Fixes

* Remove flags, use declarative approach ([0dc1b1a](https://github.com/inference-gateway/browser-agent/commit/0dc1b1a67cb3486191393605e51664a71c2fc6e6))

### 👷 CI

* **claude:** Simplify conditions for triggering Claude Code actions ([f76ceea](https://github.com/inference-gateway/browser-agent/commit/f76ceeafbd7e4c2f1f92d0e75803ddb59ef1e12c))
* **deps:** Bump the github-actions group with 3 updates ([#66](https://github.com/inference-gateway/browser-agent/issues/66)) ([e9281c9](https://github.com/inference-gateway/browser-agent/commit/e9281c9289ce20048264d17f2ddf42eb4540acb9))
* **deps:** Update claude-code-action to version 1.0.130 ([2413f3e](https://github.com/inference-gateway/browser-agent/commit/2413f3ebe517ab9f73d28ae9d95dd1fcfe1de96f))

### 🔧 Miscellaneous

* **adl:** Refresh agent.yaml defaults from ADL CLI v0.33.1 ([#63](https://github.com/inference-gateway/browser-agent/issues/63)) ([6341b52](https://github.com/inference-gateway/browser-agent/commit/6341b5222ec1cd067f5c26052d7990b3d3ce0cc6))
* **adl:** Refresh agent.yaml defaults from ADL CLI v0.36.0 ([#69](https://github.com/inference-gateway/browser-agent/issues/69)) ([a957176](https://github.com/inference-gateway/browser-agent/commit/a957176f774fc4b0e9762211e1d406f28e881de2))
* **adl:** Refresh agent.yaml defaults from ADL CLI v0.38.1 ([#76](https://github.com/inference-gateway/browser-agent/issues/76)) ([6efba31](https://github.com/inference-gateway/browser-agent/commit/6efba3127a74ec607b788142498d5ed604b837dd))
* Allow the CLI to fetch artifacts and enable artifacts server ([be5701f](https://github.com/inference-gateway/browser-agent/commit/be5701fedf4efddcc6ebd2bf56dcde9f2dbc236a))
* **dependabot:** Update golang and ubuntu version ignore rules in dependabot configuration ([fc92512](https://github.com/inference-gateway/browser-agent/commit/fc92512bba0b81d3beb9e043dc4408c1b7a4fd46))
* **deps:** Bump ADL CLI to v0.31.0 ([#62](https://github.com/inference-gateway/browser-agent/issues/62)) ([15b33e2](https://github.com/inference-gateway/browser-agent/commit/15b33e26506f7c2026e57578dc0d24949d9be610))
* **deps:** Bump ADL CLI to v0.34.0 ([#65](https://github.com/inference-gateway/browser-agent/issues/65)) ([5d75e1e](https://github.com/inference-gateway/browser-agent/commit/5d75e1e0497221f1d61cbcb442bdaaabf84df3be))
* **deps:** Bump ADL CLI to v0.34.1 ([#67](https://github.com/inference-gateway/browser-agent/issues/67)) ([ae29031](https://github.com/inference-gateway/browser-agent/commit/ae2903116f092437a03d6a9f520b8fa7cad40878))
* **deps:** Bump ADL CLI to v0.34.2 ([#68](https://github.com/inference-gateway/browser-agent/issues/68)) ([811796b](https://github.com/inference-gateway/browser-agent/commit/811796b6bb626ccf56ba59fa4dec1cdd928d9d99))
* **deps:** Bump ADL CLI to v0.36.1 ([#71](https://github.com/inference-gateway/browser-agent/issues/71)) ([a7abd51](https://github.com/inference-gateway/browser-agent/commit/a7abd5155c6d80e29ceeee0efcec5427ce8d09b2))
* **deps:** Bump ADL CLI to v0.36.2 ([#73](https://github.com/inference-gateway/browser-agent/issues/73)) ([9697df9](https://github.com/inference-gateway/browser-agent/commit/9697df9e18f54fffb7b30f21dfd915a46ddd92de))
* **deps:** Bump ADL CLI to v0.36.4 ([#74](https://github.com/inference-gateway/browser-agent/issues/74)) ([c3e5690](https://github.com/inference-gateway/browser-agent/commit/c3e56903b20371f7e4cf46497f03271b52aff4e8))
* **deps:** Bump ADL CLI to v0.38.1 ([#75](https://github.com/inference-gateway/browser-agent/issues/75)) ([9c0170d](https://github.com/inference-gateway/browser-agent/commit/9c0170d1265dfb99206ca0a1c13f3be00bd13d14))
* Enable a few tools for artifacts ([cf18eb5](https://github.com/inference-gateway/browser-agent/commit/cf18eb53cc98aac4ac6fef688ded6c30f94fc311))
* Enable docker compose and generate ([0b5a7af](https://github.com/inference-gateway/browser-agent/commit/0b5a7af9c613bbcd7f9019207f04a4a2eabeb796))
* **flox:** Add manifest lock file ([7ca21e9](https://github.com/inference-gateway/browser-agent/commit/7ca21e90cfd3902044e88f0d0487fbf279edaaed))
* **flox:** Generate manifest lock file ([63aa810](https://github.com/inference-gateway/browser-agent/commit/63aa810a6c7665948fbb9d0e3e63d35249e42a4f))
* Generate missing parts ([1f24afe](https://github.com/inference-gateway/browser-agent/commit/1f24afeee411296a8cc766cca7ab2d0387c5624e))
* **license:** Update license to Apache 2.0 ([9642254](https://github.com/inference-gateway/browser-agent/commit/96422542154e3cb56f15aa32f4a5c111248070fb))

### ✅ Miscellaneous

* Re-generate the tests ([df97e9e](https://github.com/inference-gateway/browser-agent/commit/df97e9e51edc95b9df573674de7d8d68ef75e7d4))

## [0.5.1](https://github.com/inference-gateway/browser-agent/compare/v0.5.0...v0.5.1) (2026-05-21)

### 🐛 Bug Fixes

* **container:** Match BROWSER_ENGINE in second playwright install ([#61](https://github.com/inference-gateway/browser-agent/issues/61)) ([5681879](https://github.com/inference-gateway/browser-agent/commit/56818793e7bc563aa9613334662de7beb124f58d)), closes [#58](https://github.com/inference-gateway/browser-agent/issues/58)

### 👷 CI

* **deps:** Bump anthropics/claude-code-action from 1.0.128 to 1.0.129 in the github-actions group ([#60](https://github.com/inference-gateway/browser-agent/issues/60)) ([8e58890](https://github.com/inference-gateway/browser-agent/commit/8e588905472bff860512f901d3096c8c65c84324))
* **deps:** Bump github.com/inference-gateway/adk from 0.18.3 to 0.18.4 in the gomod group ([#57](https://github.com/inference-gateway/browser-agent/issues/57)) ([d344637](https://github.com/inference-gateway/browser-agent/commit/d344637d8f6c659401a8387e5d7f83631f766bc9))

### 🔧 Miscellaneous

* **deps:** Add ignore rule for golang dependency versions >=1.26.3 in dependabot configuration ([ddc20db](https://github.com/inference-gateway/browser-agent/commit/ddc20dbfbc35806590d53d7ebb5b4dbc471516c8))

## [0.5.0](https://github.com/inference-gateway/browser-agent/compare/v0.4.19...v0.5.0) (2026-05-21)

### ✨ Features

* **skills:** Add deep-research bare skill for multi-source investigative research ([#53](https://github.com/inference-gateway/browser-agent/issues/53)) ([f8a9750](https://github.com/inference-gateway/browser-agent/commit/f8a9750f8e60b7ca6e4baa83c336e10c9916c9bc)), closes [#52](https://github.com/inference-gateway/browser-agent/issues/52)

### 👷 CI

* **deps:** Bump anthropics/claude-code-action from 1.0.122 to 1.0.128 in the github-actions group ([#56](https://github.com/inference-gateway/browser-agent/issues/56)) ([308683f](https://github.com/inference-gateway/browser-agent/commit/308683f194ad1e2fc32445654cd5b4cf5e82a729))

### 🔧 Miscellaneous

* **deps:** Bump github.com/inference-gateway/adk from 0.18.2 to 0.18.3 in the gomod group ([#55](https://github.com/inference-gateway/browser-agent/issues/55)) ([383b078](https://github.com/inference-gateway/browser-agent/commit/383b0783763b0006154429da21b9ff215415e2f6))

## [0.4.19](https://github.com/inference-gateway/browser-agent/compare/v0.4.18...v0.4.19) (2026-05-20)

### 👷 CI

* **deps:** Update installation methods for golangci-lint and task in workflows ([9e2eb7a](https://github.com/inference-gateway/browser-agent/commit/9e2eb7a3d9c4ca60318154830de2daa3613d690f))

### 🔧 Miscellaneous

* **deps:** Bump ADL CLI to version 0.30.10 ([58a9164](https://github.com/inference-gateway/browser-agent/commit/58a916463d1983730d62c6a18d90cc6af51f9209))
* **deps:** Bump ADL CLI to version 0.30.9 ([9aa5eb2](https://github.com/inference-gateway/browser-agent/commit/9aa5eb2b9019d32867546b678fa884913c3f72ed))
* Run task generate ([9e34f0d](https://github.com/inference-gateway/browser-agent/commit/9e34f0d492618dc4e5aa09cd733f2c3c857d93fd))

## [0.4.18](https://github.com/inference-gateway/browser-agent/compare/v0.4.17...v0.4.18) (2026-05-19)

### ♻️ Improvements

* Separate function-call tools from markdown skill playbooks ([#51](https://github.com/inference-gateway/browser-agent/issues/51)) ([775ce5c](https://github.com/inference-gateway/browser-agent/commit/775ce5cfda6d3b47ea91cc6573d69cd68230b18d))

### 👷 CI

* **dependabot:** Add dependabot to help with dependecies upgrades ([88a5a3f](https://github.com/inference-gateway/browser-agent/commit/88a5a3f0a58f9ef785c898e1bcbebcfe4de4f6ed))
* **deps:** Bump the github-actions group with 5 updates ([#49](https://github.com/inference-gateway/browser-agent/issues/49)) ([e4ca5c8](https://github.com/inference-gateway/browser-agent/commit/e4ca5c8715b88f4a409be6467a24644967fd336e))
* Enable display report for Claude Code action ([40650b2](https://github.com/inference-gateway/browser-agent/commit/40650b22d68cfc43bfeca8d7db179fe29f28d767))
* Update create-github-app-token action to v3.2.0 ([a9a450a](https://github.com/inference-gateway/browser-agent/commit/a9a450a4d5a9d89f121752bdefd782b95f225448))

### 🔧 Miscellaneous

* Add CODEOWNERS file for repository ownership ([95d1a0c](https://github.com/inference-gateway/browser-agent/commit/95d1a0cfb7e0679e4bc33d974cf4de00683a76e9))
* **deps:** Bump the docker group with 2 updates ([#48](https://github.com/inference-gateway/browser-agent/issues/48)) ([d67c8c9](https://github.com/inference-gateway/browser-agent/commit/d67c8c9437769a8438ddab7982f73ca71bf7f22e))
* **deps:** Bump the gomod group with 3 updates ([#50](https://github.com/inference-gateway/browser-agent/issues/50)) ([c8e1d04](https://github.com/inference-gateway/browser-agent/commit/c8e1d04cd0018cb39a27e4864ef6475d60b85bcf))
* Remove outdated issue templates for bug reports, feature requests, and refactor requests ([3cc7fc7](https://github.com/inference-gateway/browser-agent/commit/3cc7fc77c865fde1bd476f1b4807e9b001d59d94))

## [0.4.17](https://github.com/inference-gateway/browser-agent/compare/v0.4.16...v0.4.17) (2026-05-07)

### ♻️ Improvements

* Rename all instances of deepseek-chat to deepseek-v4-flash ([71b223e](https://github.com/inference-gateway/browser-agent/commit/71b223e36ce53cc53adeb8caaf08162b792269f3))

### 👷 CI

* **deps:** Bump golangci-lint to latest ([5479089](https://github.com/inference-gateway/browser-agent/commit/5479089c4d959d95562ec567eb8bac9b814c27d2))
* Update golangci-lint installation script and version in CI workflows ([dd4e1c0](https://github.com/inference-gateway/browser-agent/commit/dd4e1c099e5769dfd22f39486b9795df315d09a7))

### 🔧 Miscellaneous

* **deps:** Bump claude code action ([3f3476c](https://github.com/inference-gateway/browser-agent/commit/3f3476c79c4c4dbc191bc252431d277e3c9723e1))
* Update task installation method in CI and CD workflows ([494df52](https://github.com/inference-gateway/browser-agent/commit/494df52eeac218223fb68b29a6a32d3e83e57054))

## [0.4.16](https://github.com/inference-gateway/browser-agent/compare/v0.4.15...v0.4.16) (2026-04-17)

### 🔧 Miscellaneous

* **deps:** Bump ADL CLI and re-generate ([6ed0d25](https://github.com/inference-gateway/browser-agent/commit/6ed0d25d2000ed18ba981abeec9b63651c7c4d5c))
* **deps:** Bump ADL CLI to 0.27.6 and re-generate ([f16ee2c](https://github.com/inference-gateway/browser-agent/commit/f16ee2c840bad747d29427cb4e0f67070111eed0))

## [0.4.15](https://github.com/inference-gateway/browser-agent/compare/v0.4.14...v0.4.15) (2026-04-10)

### ♻️ Improvements

* Use fmt.Fprintf for strings.Builder writes ([054c42d](https://github.com/inference-gateway/browser-agent/commit/054c42d1b263e66583b1457a905f588102c96f13))

## [0.4.14](https://github.com/inference-gateway/browser-agent/compare/v0.4.13...v0.4.14) (2026-04-10)

### 🔧 Miscellaneous

* **deps:** Re-generate files - bump the versions ([1f37063](https://github.com/inference-gateway/browser-agent/commit/1f37063acb5a6bc2bc835c30502f8861d3703747))

## [0.4.13](https://github.com/inference-gateway/browser-agent/compare/v0.4.12...v0.4.13) (2026-01-27)

### 🐛 Bug Fixes

* Update artifact parts access for ADK 0.17.0 type changes ([cbf1b01](https://github.com/inference-gateway/browser-agent/commit/cbf1b01153d982f6bd25b615ca0c2d90bcd00a2b))

### 🔧 Miscellaneous

* **ci:** Update Claude Code workflow configuration ([b38f763](https://github.com/inference-gateway/browser-agent/commit/b38f7630635cf4e356510aae279e5f16a7f688e1))
* **deps:** Update dependencies and regenerate with ADL CLI v0.27.1 ([6e308bb](https://github.com/inference-gateway/browser-agent/commit/6e308bb61a6dda8423ab53c36aa3a896b1330d0f))
* Update Infer CLI configuration ([7561647](https://github.com/inference-gateway/browser-agent/commit/7561647255e03c92e3f215c3f12627114dda7bee))

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
