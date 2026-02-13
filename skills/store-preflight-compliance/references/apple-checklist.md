# Apple App Store Checklist (Source-Only)

Use this list when you cannot rely on compiled binary checks.

## Metadata and Config

- Verify app name and bundle identifier are clearly defined.
- Check for privacy policy URL references in project config (if missing, mark manual App Store Connect requirement).
- Flag placeholder metadata and obvious draft content.

## Privacy and Tracking

- Check for `PrivacyInfo.xcprivacy`.
- Look for tracking/ads SDK usage and ATT implementation consistency.
- Flag missing purpose-string quality for sensitive permissions where visible in config/plist files.

## Code and Content Risks

- Private API or risky dynamic execution indicators.
- Hardcoded secrets/tokens.
- Insecure `http://` URLs.
- Hardcoded IPs.
- Placeholder copy (`Lorem ipsum`, `TBD`, `Coming soon`).
- Mentions of competing platforms in user-facing App Store strings.

## Account/Identity/Payments

- Social login patterns should consider Apple Sign In requirements where relevant.
- If account creation exists, verify deletion-path evidence exists.
- For digital goods, flag external payment patterns that may violate App Store policy.

## Manual Console Items (Always Include)

- App Store Connect metadata completeness
- Screenshots/device size coverage
- Age rating and encryption declarations
- Pricing/availability validation
