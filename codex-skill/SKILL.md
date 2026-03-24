---
name: store-preflight-compliance
description: Pre-submission compliance checkup for Google Play and Apple App Store apps. Use when reviewing mobile projects for store rejection risks, submission readiness, privacy/policy compliance, and release checkups across Android and iOS.
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

Common fixes:

- Resolve Android release/policy risks (`debuggable`, cleartext traffic, sensitive permissions, target SDK, versionCode).
- Move hardcoded secrets to environment variables.
- Replace external payment flows for digital goods with Play Billing on Android and StoreKit/IAP on Apple.
- Add Sign in with Apple when social login exists (Apple policy).
- Add account deletion when account creation exists (both stores).
- Remove references to competing platforms.
- Replace placeholder text (`Lorem ipsum`, `TBD`, `Coming soon`).
- Rewrite vague purpose strings with concrete app behavior.
- Replace hardcoded IPs with hostnames.
- Replace `http://` URLs with `https://`.
- Remove debug logs or gate them behind development flags.
- Add missing privacy policy URL and required store metadata.

## Step 3: Re-Run Until READY

```bash
storeready playstore-checkup .
storeready appstore-checkup .
```

Continue until output reports READY (zero `CRITICAL` findings).

## Useful Commands

```bash
storeready play-guidelines list
storeready codescan .
storeready privacy .
storeready ipa /path/to/build.ipa
storeready scan --app-id <ID>
storeready release-checklist --app-type all
storeready publish --app-id <ID> --version <X.Y.Z> [--build <BUILD_ID>] [--confirm]
storeready guidelines search "privacy"
```
