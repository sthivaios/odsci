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
	"time"

	"go.bug.st/serial"
)

func ClearBuffer(port serial.Port, scanner *bufio.Scanner) {
	port.SetReadTimeout(100 * time.Millisecond)
	buf := make([]byte, 256)
	for {
		n, _ := port.Read(buf)
		if n == 0 {
			break
		}
	}
	// reset to blocking
	port.SetReadTimeout(serial.NoTimeout)
}