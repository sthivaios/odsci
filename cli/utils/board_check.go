/*
ODSCI CLI - An STM32-based USB interface for DS18B20 temperature sensors
# Copyright (C) 2026  Stratos Thivaios

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package utils

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"go.bug.st/serial"
)

func BoardCheck(port serial.Port, scanner *bufio.Scanner) (BoardInfo, error) {
	port.Write([]byte("GET_INFO\r"))
	scanner.Scan()
	line := scanner.Text()
	if strings.HasPrefix(line, "ERROR:") {
		var errorString string = fmt.Sprintf("ODSCI error: %s", line);
		return BoardInfo{}, errors.New(errorString)
	}
	parts := strings.Split(line, ",")

	var boardInfo BoardInfo

	for _, part := range parts {
		key, value, found := strings.Cut(part, "=")
		if !found {
			return BoardInfo{}, errors.New("error parsing response")
		}

		switch key {
			case "FIRMWARE_VERSION":
				boardInfo.FirmwareVersion = value

			case "CLED_IS_FOR_ERRORS_INSTEAD":
				boardInfo.CledIsUsedForErrors = (value == "1")

			case "LAST_RESET_DUE_TO_IWDG":
				boardInfo.LastResetWasIWDG = (value == "1")

			case "SERIAL_NUMBER":
				boardInfo.SerialNumber = value // add this field to struct
		}
	}

	// placeholder for now
	return boardInfo, nil
}