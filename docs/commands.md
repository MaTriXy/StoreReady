# Command Reference

Store-specific docs: [Google Play](google-play.md) &bull; [Apple App Store](apple-app-store.md)

---

## `codescan` — Cross-platform code pattern scan

```bash
storeready codescan .
```

Scans Swift, Objective-C, Kotlin, React Native, and Expo projects for store rejection and policy-risk patterns.

| Severity | Pattern | Guideline |
|----------|---------|-----------|
| **CRITICAL** | Private API usage | §2.5.1 |
| **CRITICAL** | Hardcoded secrets / API keys | §1.6 |
| **CRITICAL** | External payment for digital goods | §3.1.1 |
| **CRITICAL** | Dynamic code execution | §2.5.2 |
| **CRITICAL** | Cryptocurrency mining | §3.1.5 |
| WARN | Missing Sign in with Apple (social login present) | §4.8 |
| WARN | Missing Restore Purchases for IAP | §3.1.1 |
| WARN | Missing ATT for ad/tracking SDKs | §5.1.2 |
| WARN | Account creation without deletion option | §5.1.1 |
| WARN | Placeholder content in strings | §2.1 |
| WARN | References to competing platforms | §2.3 |
| INFO | Hardcoded IPv4 addresses | §2.5 |
| INFO | Insecure HTTP URLs | §1.6 |
| INFO | Vague Info.plist purpose strings | §5.1.1 |
| INFO | Expo config issues | §2.1 |

## Output Formats

All commands support:

```bash
--format terminal   # colored terminal output (default)
--format json       # JSON for CI/CD pipelines
--output file.json  # write to file instead of stdout
```

## Architecture

```
storeready
├── playstore-checkup     Google Play local checks
├── play-guidelines       Google Play policy matrix + checklists
│
├── appstore-checkup      Apple local preflight checks (alias: preflight)
│   ├── metadata          app.json / Info.plist analysis
│   ├── codescan          30+ rejection-risk code patterns
│   ├── privacy           Privacy manifest + Required Reason APIs
│   └── ipa               Binary inspection (optional)
│
├── codescan              Code-only scanning (both stores)
├── privacy               Privacy-only scanning
├── ipa                   Binary-only inspection
│
├── scan                  App Store Connect API checks (tiers 1–4)
├── auth                  App Store Connect authentication
├── guidelines            Apple guideline browser
├── play-guidelines       Google Play policy browser
└── version               Print version
```
