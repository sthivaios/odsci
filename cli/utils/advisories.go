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
during startup, indicates an IWDG reset.

To stop this advisory from displaying again, reset the board
manually by unplugging it and plugging it back in.`, boardInfo.FirmwareVersion, boardInfo.SerialNumber) + color.HiGreenString("\r\n\r\nThe command you ran will not be affected by this.")+ color.HiBlueString("\r\n====== END OF ADVISORY ======\r\n")
}