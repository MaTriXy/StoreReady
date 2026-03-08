# Google Play Checklist (Source-Only)

Use this list when you cannot rely on compiled/built artifact checks.

## Manifest and Release Flags

- Verify `AndroidManifest.xml` exists.
- Flag `android:debuggable="true"` as release-blocking.
- Flag `android:usesCleartextTraffic="true"` unless justified.
- Flag risky backup defaults for sensitive apps.

## Permissions and Policy Declarations

- Detect high-risk permissions (SMS/call-log/all-files/query-all-packages).
- Mark each high-risk permission for manual Play declaration verification.

## Gradle Release Metadata

- Verify presence of `applicationId`.
- Verify `targetSdk`/`targetSdkVersion` exists and is current enough for Play requirements.
- Verify `versionCode` exists and is incrementable for release flow.

## Store Policy Manual Checks (Always Include)

- Data safety form accuracy vs app behavior
- Account deletion requirements (if accounts exist)
- Billing/payments policy compliance for digital goods
- Listing accuracy (screenshots/claims/metadata)
- App quality checks (crash/ANR health, compatibility)

## Notes

Automated source checks reduce risk, but final Play approval also depends on Play Console declarations and policy forms. Always include these manual checks in the final report.
