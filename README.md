# StoreReady

**Know before you submit.** Pre-submission checkup for mobile apps on Apple App Store and Google Play.

StoreReady checks your app — source code, privacy manifests, Android manifests/Gradle config, IPA binaries, and store metadata — against mobile store submission requirements before you release.

## Install

```bash
# Homebrew (macOS)
brew install matrixy/tap/storeready

# Go
go install github.com/MaTriXy/StoreReady/cmd/storeready@latest

# Build from source
git clone https://github.com/MaTriXy/StoreReady.git
cd storeready && make build
# Binary at: build/storeready
```

## Quick Start

```bash
# Apple App Store checkup (offline)
storeready appstore-checkup /path/to/your/project

# Google Play checkup (offline)
storeready playstore-checkup /path/to/your/project

# Browse Play policy matrix and checklist
storeready play-guidelines list
storeready play-guidelines show GP-3

# Include IPA for binary analysis
storeready preflight . --ipa build.ipa
```

That's it. You get store-specific findings fast, with automated checks and manual policy checklist coverage.

## Store Coverage

- **Apple App Store**: local preflight checks (`preflight` / `appstore-checkup`), built-in Apple guideline browser (`guidelines`), and App Store Connect API checks (`scan`).
- **Apple App Store**: local preflight checks (`preflight` / `appstore-checkup`), built-in Apple guideline browser (`guidelines`), App Store Connect API checks (`scan`), and end-to-end ASC release lane (`publish`).
- **Apple App Store**: local preflight checks (`preflight` / `appstore-checkup`), built-in Apple guideline browser (`guidelines`), App Store Connect API checks (`scan`), structured manual release gates (`release-checklist`), and end-to-end ASC release lane (`publish`).
- **Google Play**: local project checks (`playstore-checkup`) and built-in Play policy matrix/checklists (`play-guidelines`).

## Commands

### `storeready preflight [path]` — Apple local preflight checks

Runs Apple-focused local scanners in parallel. No account needed. Entirely offline.

```bash
storeready preflight .                          # scan current directory
storeready preflight ./my-app --ipa build.ipa   # with binary inspection
storeready preflight . --format json            # JSON output for CI/CD
storeready preflight . --output report.json     # write to file
```

**Scanners included:**

| Scanner | Checks |
|---------|--------|
| **metadata** | app.json / Info.plist: name, version, bundle ID format, icon, privacy policy URL, purpose strings |
| **codescan** | 30+ code patterns: private APIs, secrets, payment violations, missing ATT, social login, placeholders |
| **privacy** | PrivacyInfo.xcprivacy completeness, Required Reason APIs, tracking SDKs vs ATT implementation |
| **ipa** | Binary: Info.plist keys, launch storyboard, app icons, app size, framework privacy manifests |

### `storeready appstore-checkup [path]` — Apple checkup alias

This is an alias of `storeready preflight [path]` with the same output and flags.

### `storeready playstore-checkup [path]` — Google Play local checks

```bash
storeready playstore-checkup .
storeready playstore-checkup ./android-app --format json --output play-report.json
```

Checks include:
- AndroidManifest presence and risky release flags (`debuggable`, `usesCleartextTraffic`, `allowBackup`)
- High-risk Play permissions that require strong policy justification
- Gradle metadata required for release readiness (`applicationId`, `targetSdk`, `versionCode`)
- Package ID sanity checks (placeholder IDs)
- Policy checklist coverage output (automated pass/fail + required manual review controls)

### `storeready play-guidelines` — Google Play policy matrix

```bash
storeready play-guidelines list                  # top-level matrix
storeready play-guidelines show GP-2             # exact section checks + sources
storeready play-guidelines search permissions    # full-text search
```

Sections are tagged as:
- `AUTOMATED`: verified by scanner rules
- `HYBRID`: partially automated, still requires console/policy review
- `MANUAL`: must be reviewed manually every release

### `storeready codescan [path]` — Code pattern scan (Apple + Play risk patterns)

```bash
storeready codescan /path/to/project
```

Scans Swift, Objective-C, React Native, and Expo projects for common store rejection and policy-risk patterns:
- Private API usage (§2.5.1) — **CRITICAL**
- Hardcoded secrets/API keys (§1.6) — **CRITICAL**
- External payment for digital goods (§3.1.1) — **CRITICAL**
- Dynamic code execution (§2.5.2) — **CRITICAL**
- Cryptocurrency mining (§3.1.5) — **CRITICAL**
- Missing Sign in with Apple when using social login (§4.8)
- Missing Restore Purchases for IAP (§3.1.1)
- Missing ATT for ad/tracking SDKs (§5.1.2)
- Account creation without deletion option (§5.1.1)
- Placeholder content in strings (§2.1)
- References to competing platforms (§2.3)
- Hardcoded IPv4 addresses (§2.5)
- Insecure HTTP URLs (§1.6)
- Vague Info.plist purpose strings (§5.1.1)
- Expo config issues (§2.1)

### `storeready privacy [path]` — Privacy manifest validator

```bash
storeready privacy /path/to/project
```

Deep privacy compliance scan:
- PrivacyInfo.xcprivacy exists and is properly configured
- Required Reason APIs detected in code vs declared in manifest
- Tracking SDKs detected vs ATT implementation
- Cross-references everything automatically

### `storeready ipa <path.ipa>` — Binary inspector

```bash
storeready ipa /path/to/build.ipa
```

Inspects a built IPA for:
- PrivacyInfo.xcprivacy presence
- Info.plist completeness and purpose string quality
- App Transport Security configuration
- App icon presence and sizes
- Launch storyboard presence
- App size vs 200MB cellular download limit
- Embedded framework privacy manifests

### `storeready scan --app-id <ID>` — Apple App Store Connect checks

```bash
storeready auth setup                    # one-time: configure API key
storeready auth login                    # or: sign in with Apple ID
storeready scan --app-id 6758967212     # run all tiers
```

API-based checks against your app in App Store Connect:
- Metadata completeness (descriptions, keywords, URLs)
- Screenshot verification for required device sizes
- Build processing status
- Age rating and encryption compliance
- Content rights declaration
- App-info privacy policy URL coverage by locale
- Version copyright metadata presence
- Content analysis (platform references, placeholders)

### `storeready release-checklist` — Manual ASC/App Review release gates

```bash
storeready release-checklist
storeready release-checklist --app-type subscription
storeready release-checklist --format json --output release-checklist.json
```

Outputs a structured list of **non-fully-automatable** checks to review before submit:
- Submission state, review notes, backend readiness, metadata truthfulness
- Policy-sensitive flows (account deletion, IAP/restore, SIWA parity, legal links)
- App-type specific gates (`subscription`, `social`, `kids`, `health`, `games`, `macos`, `ai`, `crypto`, `vpn`)

### `storeready publish` — End-to-end StoreReady gate + ASC release lane

```bash
# Safe default: runs local + ASC gates, then asc release dry-run
storeready publish \
  --app-id 6758967212 \
  --version 1.2.3 \
  --build 123456789

# Real run (submits through asc release flow after gates pass)
storeready publish \
  --app-id 6758967212 \
  --version 1.2.3 \
  --build 123456789 \
  --metadata-dir ./metadata/version/1.2.3 \
  --confirm
```

Behavior:
- Runs local StoreReady preflight checks (unless `--skip-local-checks`)
- Runs StoreReady ASC metadata gates via `scan` engine (unless `--skip-asc-scan`)
- Executes `asc release run` (dry-run by default; real run with `--confirm`)

### `storeready guidelines` — Browse Apple App Store guidelines

```bash
storeready guidelines list               # all sections
storeready guidelines show 2.1           # specific guideline
storeready guidelines search "privacy"   # full-text search
```

### Output formats

All scan commands support:

```bash
--format terminal   # colored terminal output (default)
--format json       # JSON for CI/CD pipelines
--output file.json  # write to file instead of stdout
```

## Claude Code Skill

StoreReady includes a Claude Code skill that can be installed via the `skills.sh` CLI.

### Setup

Install directly from this GitHub repo:

```bash
npx skills add https://github.com/MaTriXy/StoreReady --skill "Store Preflight Compliance"
```

Then in Claude Code, invoke:

```text
Use $store-preflight-compliance to review this repo for Apple App Store and Google Play submission readiness.
```

Claude will:
1. Detect Apple/Google Play scope from source files
2. Read every finding
3. Prioritize issues by severity (`CRITICAL`, `WARN`, `INFO`)
4. Include manual store-console checklist items not covered by static checks

## Codex Skill

StoreReady includes a Codex-native skill package at `codex-skill/`.

### Setup

```bash
mkdir -p ~/.codex/skills/store-preflight-compliance
cp -R codex-skill/* ~/.codex/skills/store-preflight-compliance/
```

Then in Codex, invoke:

```text
Use $store-preflight-compliance to run StoreReady Apple + Play checks and fix all findings until READY.
```

## Architecture

```
StoreReady
├── appstore-checkup   Alias for Apple preflight checks
├── playstore-checkup  Google Play local checks
├── play-guidelines    Google Play policy matrix + checklist references
│
├── preflight         Run ALL checks — one command
│   ├── metadata      app.json / Info.plist local analysis
│   ├── codescan      30+ rejection-risk code patterns
│   ├── privacy       Privacy manifest + Required Reason APIs
│   └── ipa           Binary inspection (optional)
│
├── codescan          Code-only scanning
├── privacy           Privacy-only scanning
├── ipa               Binary-only inspection
│
├── scan              Apple App Store Connect API checks (tiers 1-4)
│   ├── Tier 1        Metadata & completeness
│   ├── Tier 2        Content analysis
│   ├── Tier 3        Binary inspection
│   └── Tier 4        Historical pattern matching
│
├── release-checklist Structured manual ASC/App Review gate list
│
├── publish           End-to-end local+ASC gate and ASC CLI release lane
│
├── auth              App Store Connect authentication
│   ├── login         Apple ID + 2FA session auth
│   ├── setup         API key configuration
│   ├── status        Show current auth state
│   └── logout        Remove credentials
│
└── guidelines        Built-in Apple App Store Review Guidelines database
    ├── list          All 5 sections with subsections
    ├── show          Specific guideline details
    └── search        Full-text search
```

## CI/CD Integration

```yaml
# GitHub Actions
- name: Mobile store compliance check
  run: |
    storeready preflight . --format json --output storeready-report.json
    # Fail the pipeline if critical issues found
    if jq -e '.summary.critical > 0' storeready-report.json > /dev/null; then
      echo "CRITICAL issues found — fix before submission"
      exit 1
    fi
```

```yaml
# JUnit output for test reporting (scan command only)
storeready scan --app-id $APP_ID --format junit --output storeready.xml
```
