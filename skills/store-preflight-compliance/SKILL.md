---
name: store-preflight-compliance
description: "Audit mobile app source code for Google Play and Apple App Store submission compliance. Checks AndroidManifest flags, Gradle release metadata, high-risk permissions, privacy manifests, purpose strings, hardcoded secrets, and common rejection patterns. Produces a severity-ranked findings report with fix recommendations and a READY/NOT READY verdict. Use when reviewing mobile projects for store rejection risks, submission readiness, privacy and policy compliance, or release checkups across Android and iOS."
---

# Store Preflight Compliance

Audit mobile app source for Google Play and Apple App Store compliance. Produce a severity-ranked report with fix recommendations and a READY/NOT READY verdict — no binary build required.

## Step 1: Detect Platform Scope

Identify which store targets are present:

```bash
# Android indicators
find . -name "AndroidManifest.xml" -o -name "build.gradle" -o -name "build.gradle.kts" | head -5

# Apple indicators
find . -name "Info.plist" -o -name "PrivacyInfo.xcprivacy" -o -name "*.xcodeproj" | head -5
grep -rl '"expo".*"ios"' app.json 2>/dev/null
```

If only one platform is detected, review that platform and note the other was not found.

## Step 2: Run Google Play Source Checks

Follow the detailed checklist in `references/play-checklist.md`. At minimum:

```bash
# Release-blocking flags
grep -rn 'android:debuggable="true"' . --include="AndroidManifest.xml"
grep -rn 'android:usesCleartextTraffic="true"' . --include="AndroidManifest.xml"

# High-risk permissions
grep -rn 'android.permission.\(READ_SMS\|READ_CALL_LOG\|MANAGE_EXTERNAL_STORAGE\|QUERY_ALL_PACKAGES\)' . --include="AndroidManifest.xml"

# Gradle release metadata
grep -rn 'applicationId\|targetSdk\|versionCode' . --include="*.gradle" --include="*.gradle.kts" | head -10
```

Also flag: data safety form accuracy, account deletion requirements, billing policy compliance for digital goods, and listing accuracy — these require manual Play Console review.

**Checkpoint:** Confirm at least one Android config file was found and parsed before proceeding.

## Step 3: Run Apple Source Checks

Follow the detailed checklist in `references/apple-checklist.md`. At minimum:

```bash
# Privacy manifest
find . -name "PrivacyInfo.xcprivacy" | head -3

# Hardcoded secrets / insecure URLs
grep -rn 'http://' . --include="*.swift" --include="*.m" --include="*.js" --include="*.ts" | grep -v node_modules | head -10
grep -rn 'sk_live_\|pk_live_\|AIza\|AKIA' . --include="*.swift" --include="*.js" --include="*.ts" | head -5

# Placeholder content
grep -rni 'lorem ipsum\|coming soon\|\bTBD\b' . --include="*.swift" --include="*.js" --include="*.tsx" | head -5
```

Also flag: missing privacy policy URL, social login without Sign in with Apple, account creation without deletion path, and competing platform references.

**Checkpoint:** Confirm at least one Apple config file was found and parsed before proceeding.

## Step 4: Produce Report

Structure the output as:

1. **Scope detected** — which platforms and key config files found
2. **Google Play findings** — sorted by severity
3. **Apple findings** — sorted by severity
4. **Manual console checklist** — items requiring human review in Play Console / App Store Connect
5. **Release recommendation** — `READY` (zero CRITICAL) or `NOT READY`

Each finding must include:

| Field | Example |
|-------|---------|
| Severity | `CRITICAL` |
| Title | Debuggable flag enabled in release manifest |
| Evidence | `AndroidManifest.xml:12 — android:debuggable="true"` |
| Fix | Set `android:debuggable="false"` or remove the attribute (defaults to false) |

## Optional: StoreReady CLI Fast Path

If `storeready` is available in PATH, use it to accelerate automated checks:

```bash
storeready playstore-checkup .
storeready appstore-checkup .
```

Still validate manual policy checklist items from reference files — not every store requirement is automatable.
