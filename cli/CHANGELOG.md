# Changelog (for the ODSCI CLI only)

Note: Changelog started at v1.0.0. See commit history for changes prior to this version.

## [Unreleased]

### Added

- Added the ability to use ISO 8601 instead of UNIX timestmaps in the CSV output of the capture command

### Fixed

- Fixed minor bug where the capture command would attempt to turn on the CLED without checking the flag first

## [v1.1.1] - 2026-03-31

### Fixed

- The `read` and `capture` commands, now require the `--port` and `--output` flags where applicable.

## [v1.1.0] - 2026-03-31

### Added

- Added the `read` command

### Changed

- Capture command checks firmware CLED flag before sending CLED commands; see PR #6

### Fixed

- Fixed a bug where the first temperature reading would fail because the serial buffer wasn't cleared
- Improved error handling in mulitple functions

## [v1.0.0] - 2026-03-29

### Added

- Added the `version` command to print the version string

### Breaking Changes

- Removed the `scan` command
