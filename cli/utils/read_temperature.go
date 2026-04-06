package utils

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"go.bug.st/serial"
)

func ReadTemperature(port serial.Port, scanner *bufio.Scanner) (string, float64, error) {
	port.Write([]byte("GET_TEMPERATURE\r"))
	scanner.Scan()
	line := scanner.Text()
	if strings.HasPrefix(line, "ERROR:") {
		return color.HiRedString("[ERROR FROM DEVICE] %s", line), -4096, errors.New("CRC")
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
	if err != nil {
		return color.HiRedString("Error parsing the temperature reading"), -4096, errors.New("PARSE")
	} else {
		return fmt.Sprintf("%0.2f", value), value, nil
	}
}