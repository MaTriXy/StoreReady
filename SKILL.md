---
name: storeready
description: "Run automated Google Play and Apple App Store compliance checks using the storeready CLI. Scans Android Manifest flags, Gradle metadata, permissions declarations, privacy manifests, purpose strings, hardcoded secrets, insecure URLs, and common rejection patterns across Kotlin, Swift, React Native, and Expo projects. Use when preparing a mobile app for submission, auditing store-policy compliance, or checking Play Store and App Store readiness."
---

# StoreReady — Mobile Store Pre-Submission Checkup

Run automated compliance checks with the `storeready` CLI, fix every issue by severity, and re-run until the app reaches READY status (zero CRITICAL findings).

## Step 1: Run the scan

Run both store checkups on the project root. The `storeready` CLI is already available in PATH:

```bash
storeready playstore-checkup .
storeready appstore-checkup .
```

If the user has a built IPA, include it:

```bash
storeready appstore-checkup . --ipa /path/to/build.ipa
```

If `storeready` is not found, install via one of:

```bash
brew install matrixy/tap/storeready                              # Homebrew (macOS)
go install github.com/MaTriXy/StoreReady/cmd/storeready@latest  # Go install
```

## Step 2: Read the output and fix every issue

Fix findings in severity order:

1. **CRITICAL** — Will be rejected. Must fix.
2. **WARN** — High rejection risk. Should fix.
3. **INFO** — Best practice. Consider fixing.

Common fix patterns:

- **Android release flags** → Set `debuggable` to false, disable cleartext traffic, review high-risk permissions, update target SDK and versionCode.
- **Hardcoded secrets** → Move to environment variables (`process.env.VAR_NAME` or Expo `Constants.expoConfig.extra`).
- **External payment for digital goods** → Replace Stripe/PayPal with Play Billing (Android) and StoreKit/IAP (Apple).
- **Social login without Sign in with Apple** → Add `expo-apple-authentication` alongside Google/Facebook login.
- **Account creation without deletion** → Add a "Delete Account" option in settings.
- **Platform references** → Remove mentions of competing platforms.
- **Placeholder content** → Replace "Lorem ipsum", "Coming soon", "TBD" with real content.
- **Vague purpose strings** → Rewrite to explain specifically why the app needs the permission (e.g. "PostureGuard uses your camera to analyze sitting posture in real-time").
- **Hardcoded IPv4 / HTTP URLs** → Replace IPs with hostnames, `http://` with `https://`.
- **Console logs** → Remove or gate behind `__DEV__` flag.
- **Missing privacy policy** → Set in Play Console and App Store Connect.

## Step 3: Re-run and repeat

```bash
storeready playstore-checkup .
storeready appstore-checkup .
```

**Keep looping until READY status (zero CRITICAL findings).** Some fixes introduce new issues (e.g. adding a tracking SDK requires ATT). The scan runs in under 1 second — re-run frequently.

## Other CLI Commands

```bash
storeready play-guidelines list            # Browse Google Play policy matrix
storeready codescan .                      # Code-only scan
storeready privacy .                       # Privacy manifest scan
storeready ipa /path/to/build.ipa          # Binary inspection
storeready scan --app-id <ID>              # App Store Connect checks (needs auth)
storeready release-checklist --app-type all
storeready publish --app-id <ID> --version <X.Y.Z> [--build <BUILD_ID>] [--confirm]
storeready guidelines search "privacy"     # Search Apple guidelines
```
