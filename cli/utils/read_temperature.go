package utils

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"go.bug.st/serial"
)

func ReadTemperature(port serial.Port, scanner *bufio.Scanner) (string, float64) {
	port.Write([]byte("GET_TEMPERATURE\r"))
	scanner.Scan()
	line := scanner.Text()
	if strings.HasPrefix(line, "ERROR:") {
		return color.MagentaString("Sensor error: %s", line), -999
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
	if err != nil {
		return color.MagentaString("Error parsing the temperature reading"), -999
	} else {
		return fmt.Sprintf("%0.2f", value), value
	}
}