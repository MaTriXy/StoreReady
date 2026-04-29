---
name: store-preflight-compliance
description: "Run StoreReady compliance checks against mobile app source and configs to catch Google Play and Apple App Store rejection risks. Audits Android Manifest flags, Gradle metadata, permissions, privacy manifests, hardcoded secrets, and common rejection patterns. Use when reviewing mobile projects for store rejection risks, submission readiness, privacy and policy compliance, or release checkups across Android and iOS."
---

# Store Preflight Compliance

Run StoreReady checks, fix findings, and repeat until the project reaches READY status.

## Workflow

1. Run `storeready playstore-checkup` and `storeready appstore-checkup` at the project root.
2. Triage findings by severity (`CRITICAL`, then `WARN`, then `INFO`).
3. Apply concrete code/configuration fixes.
4. Re-run and continue until no `CRITICAL` findings remain.
5. Complete manual store-console checklist items for policies that are not fully automatable.

## Step 1: Run Scan

```bash
storeready playstore-checkup .
storeready appstore-checkup .
```

If an IPA is available:

```bash
storeready appstore-checkup . --ipa /path/to/build.ipa
```

If `storeready` is missing, install it:

```bash
# Homebrew (macOS)
brew install matrixy/tap/storeready

# Go
go install github.com/MaTriXy/StoreReady/cmd/storeready@latest

# Build from source
git clone https://github.com/MaTriXy/StoreReady.git
cd StoreReady && make build
```

## Step 2: Fix Findings

Fix in order:

1. `CRITICAL`: must fix before submission.
2. `WARN`: high rejection risk, strongly recommended to fix.
3. `INFO`: best-practice improvements.

Store-specific fixes (apply standard code hygiene for generic issues):

- **Android release flags** → Remove `android:debuggable="true"` and `android:usesCleartextTraffic="true"` from `AndroidManifest.xml`. Verify `targetSdk` meets current Play requirements.
- **Digital goods payments** → Replace Stripe/PayPal with Play Billing (Android) and StoreKit/IAP (Apple) for in-app digital content.
- **Sign in with Apple** → Add Apple authentication when social login (Google/Facebook) exists — Apple requires this.
- **Account deletion** → Add "Delete Account" path when account creation exists (both stores require this).
- **Purpose strings** → Rewrite vague permission descriptions: not "Camera needed" but "PostureGuard uses your camera to analyze sitting posture in real-time."
- **Privacy policy** → Add URL in Play Console and App Store Connect if missing.

## Step 3: Re-Run Until READY

```bash
storeready playstore-checkup .
storeready appstore-checkup .
```

Continue until output reports READY (zero `CRITICAL` findings). Some fixes introduce new issues (e.g. adding a tracking SDK requires ATT) — re-run after each batch of changes.
