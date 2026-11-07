---
title: 'Add a generic skill for shipping a Wails app to Google Play internal testing'
priority: p2
section: ralvaskills
topics:
    - skills
    - android
---
ralvafi proved the process (2026-11-07) and five more apps will repeat it. Capture it as a generic skill, for example
`wails-android-release`, with no project or private names (the registry publishes on push to `main`, so ask the owner before
adding anything under `skills/**`):

- Application id from a domain the owner controls, reversed (it is locked after the first upload), while the Java package and
  Gradle `namespace` stay `com.wails.app` because Wails v3 exports its JNI entry points under that name.
- `minSdk` 24 or higher (Play Protect), `targetSdk` per the current Play rule, 16 KB aligned native libraries.
- `versionName` and a monotonic `versionCode` derived from the git tag, with a script and a rule for release candidates.
- Passing the signing and version env to the Gradle task, a full JDK 21, the upload key kept outside the repo.
- A Wails Android web view downgrades API calls (bodiless `GET`, only 200 passes): use a fetch bridge, and restrict it to the
  API routes.
- Verify with `bundletool validate` and `dump manifest`, then on an emulator with `-gpu host`, before uploading.
- Optional: the Play Developer API and a CI job as a later step.

Done when the skill is written, within the description and body budgets (`8k19zz`), linted, and the owner approved
publishing it. Source of truth for the details: the pilot app's `docs/guides/google-play-internal-testing.md` and
`docs/guides/wails-shell.md`.
