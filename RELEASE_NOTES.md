# Release Notes

## 2026-03-24

### Added
- `storeready publish` command for an end-to-end Apple release lane:
  - Runs local StoreReady preflight gates
  - Runs ASC metadata scan gates
  - Executes `asc release run` (dry-run by default, `--confirm` for real submission)
- `storeready release-checklist` command to output structured manual/hybrid App Store release gates with app-type profiles.

### Apple ASC Verification Coverage Expanded
- Added checks for:
  - content rights declaration
  - privacy policy URL coverage on app-info localizations
  - version copyright metadata presence

### Documentation Updates
- Updated `README.md` with Apple end-to-end lane usage.
- Updated `docs/commands.md` with `publish` and `release-checklist`.
- Updated `docs/apple-app-store.md` with corrected scan tier descriptions and release-lane docs.
