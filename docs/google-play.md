# Google Play

## `playstore-checkup` — Local project checks

```bash
storeready playstore-checkup .
storeready playstore-checkup ./android-app --format json --output play-report.json
```

| Check | What it catches |
|-------|----------------|
| Manifest flags | `debuggable`, `usesCleartextTraffic`, `allowBackup` left enabled |
| Permissions | High-risk permissions requiring Play policy declarations |
| Gradle metadata | Missing `applicationId`, `targetSdk`, `versionCode` |
| Package ID | Placeholder or default package names |
| Policy checklist | Automated pass/fail + required manual review items |

## `play-guidelines` — Policy matrix browser

```bash
storeready play-guidelines list                  # top-level matrix
storeready play-guidelines show GP-2             # section details + sources
storeready play-guidelines search permissions    # full-text search
```

Sections are tagged `AUTOMATED`, `HYBRID`, or `MANUAL` so you know what the scanner covers vs. what needs manual review.

## CI/CD

```yaml
- name: Google Play compliance check
  run: |
    storeready playstore-checkup . --format json --output play-report.json
    if jq -e '.summary.critical > 0' play-report.json > /dev/null; then
      echo "CRITICAL issues found — fix before submission"
      exit 1
    fi
```
