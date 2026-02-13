---
name: store-preflight-compliance
description: This skill should be used when the user asks to "run an app store and play store checkup", "review mobile app submission readiness", "check Apple App Store compliance", "check Google Play policy compliance", "audit store submission risks", or "prepare a mobile app for store submission".
version: 0.1.0
---

# Store Preflight Compliance

Run a full pre-submission compliance review for mobile apps targeting Apple App Store and Google Play.

## Purpose

Use this skill to produce a release-readiness report without requiring any binary build/install step. Perform static repository checks directly from source files and produce store-specific findings plus manual console review items.

## No-Build Workflow

Do not require `go build`, `make build`, or tool installation to run this skill.

If the `storeready` CLI is already available in PATH, it may be used as an optional accelerator. If it is not available, continue with source-only checks.

## Step 1: Detect Platform Scope

Identify which store targets are present in the repo:

- Apple indicators: `Info.plist`, `app.json` with `expo.ios`, `PrivacyInfo.xcprivacy`, iOS project files.
- Android indicators: `AndroidManifest.xml`, `build.gradle`, `build.gradle.kts`.

If only one platform is present, review that platform and still include a note that the other platform was not detected.

## Step 2: Run Apple Source Checks

Use the Apple checklist in `references/apple-checklist.md`.

At minimum, check:

- Metadata completeness risks (app name, bundle identifier, privacy policy references).
- Privacy manifest and required-reason API consistency.
- Common rejection patterns in code and copy (placeholder text, insecure URLs, platform-reference mistakes).
- Account and authentication policy pitfalls (for example social login patterns needing Apple Sign in support where applicable).

## Step 3: Run Google Play Source Checks

Use the Play checklist in `references/play-checklist.md`.

At minimum, check:

- Manifest release flags (`debuggable`, cleartext traffic, backup behavior).
- High-risk permissions requiring Play declarations.
- Gradle release metadata (`applicationId`, `targetSdk`, `versionCode`).
- Policy-sensitive areas requiring manual Play Console review (Data safety, account deletion, payments disclosures, listing accuracy).

## Step 4: Produce Report

Produce output in this structure:

1. Scope detected
2. Apple findings
3. Google Play findings
4. Manual console checklist items
5. Release recommendation (`READY` / `NOT READY`)

For each finding include:

- Severity (`CRITICAL`, `WARN`, `INFO`)
- Title
- Evidence (file path + short snippet/observation)
- Fix recommendation

## Optional Fast Path

If `storeready` is installed, optionally run:

```bash
storeready appstore-checkup .
storeready playstore-checkup .
```

Still validate manual policy checklist items from reference files, because not every store requirement is automatable.
