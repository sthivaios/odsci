# Changelog (for the ODSCI firmware only)
 
Note: Changelog started at v1.1.0. See commit history for changes prior to this version.

## [v1.2.0] - 2026-04-07

### Added
- Configured the IWDG to timeout after 2000ms

### Changed
- The GET_INFO command now also returns the serial number and whether the last reset reason was an IWDG timeout
- The serial number of the device as well as the USB descriptor serial is derived from a new function in command_handler.c
- Updated the USB device descriptors (product string, manufacturer, etc.)

### Fixed
- Fixed a minor bug in the main firmware loop, that could potentially cause a race condition
- Improved error handling when getting the temperature
- Improved error handling when the command buffer overflows

## [v1.1.0] - 2026-03-29

### Added
- CRC validation on DS18B20 scratchpad readings

### Breaking Changes
- CAPTURE_LED now used as error indicator, SET_CLED_ON and SET_CLED_OFF removed in this mode