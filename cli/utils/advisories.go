package utils

import "github.com/fatih/color"

func AdvisoryStringIWDG(boardInfo BoardInfo) string {
	return color.HiBlueString("\r\n========= ADVISORY ==========\r\n") + color.HiYellowString(`[SEVERITY: WARNING] IWDG RESET DETECTED

Uh oh, something isn't right... :(

The board was reset by the Independent Watchdog (IWDG). This usually means
the system became unresponsive and had to be forcibly restarted.

This should normally never happen and usually indicates a bug or unexpected state.

Please consider reporting this bug to aid with the development of the project.
You should clearly describe what you were doing when the reset occurred, as well
as the version of the firmware that the board is running. Any additional information is
also helpful. You can submit a new issue at: https://github.com/sthivaios/odsci/issues/new

If you are using a board given to you by the maintainers of the project, please include
the serial number in the bug report.

Your board's details:
- Firmware version: %s
- Serial Number: %s

For your future reference, the red error LED blinking 5 times in a row
during startup indicates an IWDG reset.

To stop this advisory from displaying again, reset the board
manually by unplugging it and plugging it back in.`, boardInfo.FirmwareVersion, boardInfo.SerialNumber) + color.HiGreenString("\r\n\r\nThe command you ran will not be affected by this.")+ color.HiBlueString("\r\n====== END OF ADVISORY ======\r\n\r\n")
}

func AdvisoryStringCRC(boardInfo BoardInfo) string {
	if (boardInfo.SerialNumber == "") {
		boardInfo.SerialNumber = "No serial number returned by the device.\nPerhaps you are running firmware older than v1.2.0?"
	}
	return color.HiBlueString("\r\n========= ADVISORY ==========\r\n") + color.HiYellowString(`[SEVERITY: WARNING] CRC8 DATA VALIDATION FAILED

Uh oh, something isn't right... :(

The board has reported a failed sensor CRC data validation. The ODSCI firmware running
onboard the microcontroller on your ODSCI probe, uses the CRC8 algorithm to ensure that
the data that it read from the sensor was correct. If that validation fails, the device
throws this error, to avoid returning potentially wrong or corrupted temperature readings.

This usually doesn't happen. It can potentially be caused by electromagnetic interference,
basically noise, on the sensor data line. If your sensor is connected to the ODSCI device
with a very long cable, that could potentially be affecting it.

Your board's details:
- Firmware version: %s
- Serial Number: %s

This advisory only displays once, when a CRC error is detected.`, boardInfo.FirmwareVersion, boardInfo.SerialNumber) + color.HiGreenString("\r\n\r\nThe CLI will still attempt to poll the sensor again. Don't cancel this command yet.\r\nThis type of error is usually intermittent.")+ color.HiBlueString("\r\n====== END OF ADVISORY ======\r\n\r\n")
}