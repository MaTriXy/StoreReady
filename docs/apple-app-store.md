# Apple App Store

## `appstore-checkup` — Local preflight checks

Runs all Apple scanners in parallel. No account needed. Entirely offline.

```bash
storeready appstore-checkup .
storeready appstore-checkup ./my-app --ipa build.ipa
storeready appstore-checkup . --format json --output report.json
```

| Scanner | What it checks |
|---------|---------------|
| **metadata** | app.json / Info.plist — name, version, bundle ID, icon, privacy policy URL, purpose strings |
| **codescan** | 30+ code patterns — private APIs, secrets, payments, ATT, social login, placeholders |
| **privacy** | PrivacyInfo.xcprivacy completeness, Required Reason APIs, tracking SDKs vs ATT |
| **ipa** | Binary — Info.plist keys, launch storyboard, icons, app size, framework privacy manifests |

## `privacy` — Privacy manifest validator

```bash
storeready privacy .
```

Checks PrivacyInfo.xcprivacy completeness, cross-references Required Reason APIs in code vs manifest, and detects tracking SDKs vs ATT implementation.

## `ipa` — Binary inspector

```bash
storeready ipa /path/to/build.ipa
```

| Check | Details |
|-------|---------|
| Privacy manifest | PrivacyInfo.xcprivacy presence in app + embedded frameworks |
| Info.plist | Completeness and purpose string quality |
| App Transport Security | ATS configuration |
| App icons | Presence and required sizes |
| Launch storyboard | Presence check |
| App size | vs 200MB cellular download limit |

## `scan` — App Store Connect API checks

Requires authentication. One-time setup:

```bash
storeready auth setup                    # configure API key
storeready auth login                    # or: sign in with Apple ID
```

Then run:

```bash
storeready scan --app-id 6758967212
```

| Tier | Checks |
|------|--------|
| 1 | Metadata/completeness — app access, version state, metadata limits, screenshots, build status, age rating, encryption, content rights, privacy policy URL, copyright, territory/pricing |
| 2 | Content analysis — platform references, placeholders, URL reachability, external TestFlight coverage |
| 3 | Reserved for future binary/API expansion |
| 4 | Reserved for future historical pattern matching |

## `release-checklist` — Manual/hybrid release gate list

```bash
storeready release-checklist
storeready release-checklist --app-type subscription
```

Use this before final submit to review non-fully-automatable checks:
- Review notes and tester credentials
- Account deletion, IAP/restore, SIWA parity
- Legal links and metadata truthfulness
- App-type specific risk gates (subscription, social, kids, health, games, macos, ai, crypto, vpn)

## `publish` — End-to-end StoreReady + ASC lane

```bash
# Safe default: dry-run release flow after all gates pass
storeready publish --app-id <ID> --version <X.Y.Z> --build <BUILD_ID>

# Real submit
storeready publish --app-id <ID> --version <X.Y.Z> --build <BUILD_ID> --confirm
```

Flow:
1. Local StoreReady preflight gates
2. ASC API scan gates
3. `asc release run` execution

## `guidelines` — Apple guideline browser

```bash
storeready guidelines list               # all sections
storeready guidelines show 2.1           # specific guideline
storeready guidelines search "privacy"   # full-text search
```

## `auth` — App Store Connect authentication

```bash
storeready auth setup      # API key configuration
storeready auth login      # Apple ID + 2FA session auth
storeready auth status     # show current auth state
storeready auth logout     # remove credentials
```

## CI/CD

```yaml
- name: Apple App Store compliance check
  run: |
    storeready appstore-checkup . --format json --output appstore-report.json
    if jq -e '.summary.critical > 0' appstore-report.json > /dev/null; then
      echo "CRITICAL issues found — fix before submission"
      exit 1
    fi
```

```yaml
# JUnit output for test reporting (scan command only)
storeready scan --app-id $APP_ID --format junit --output storeready.xml
```
