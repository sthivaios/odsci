# Changelog (for the ODSCI CLI only)

Note: Changelog started at v1.0.0. See commit history for changes prior to this version.

## Unreleased

### Added
- Added the `--silence-advisories` flag to... well... silence advisories
- The capture command now informs the user at the end about how many samples were actually valid and written to the file

## [v1.2.0] - 2026-04-07

### Added

- Added the ability to use ISO 8601 instead of UNIX timestmaps in the CSV output of the `capture` command
- Added the ability to select a unit of measurement in the `read` command with the `--unit` or `-u` flag
- The `capture` command now stores values in degrees Celsius, Fahrenheit and Kelvin in the CSV in separate columns
- The `read` command now shows the difference from the last read, next to the value
- All commands that communicate with the board, now print an advisory if the last reset reason was an IWDG timeout
- All commands that communicate with the board, now print an advisory if the sensor CRC8 data validation fails

### Fixed

- Fixed minor bug where the capture command would attempt to turn on the CLED without checking the flag first
- Fixed major bug where the `read` command wouldn't actually update because it was only reading the temperature once

### Changed

- Major improvements to error handling
- Updated from v1.0.0 of the schollz/progressbar library, to v3
- Updated the way that the GET_INFO command result is decoded

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
