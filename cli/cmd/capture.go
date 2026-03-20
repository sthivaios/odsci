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

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

type Sample struct {
	timestamp int64
	value float64
}

func timeString(seconds int64) string {
	if (seconds > 60) {
		if (seconds % 60 == 0) {
			return fmt.Sprintf("%dm", seconds/60)
		} else {
			return fmt.Sprintf("%dm %ds", seconds/60, seconds % 60)
		}
	} else {
		return fmt.Sprintf("%ds", seconds)
	}
}

// captureCmd represents the capture command
var captureCmd = &cobra.Command{
	Use: "capture",
	Short: "Used to capture a number of samples from the probe",
	Long: `The capture command is used to capture a given number
of samples from the probe, at a specific interval, and export
the captured data.`,

	Run: func(cmd *cobra.Command, args []string) {
		// port, _ := cmd.Flags().GetString("port")
        samples, _ := cmd.Flags().GetInt("samples")
        interval, _ := cmd.Flags().GetInt("interval")
        // output, _ := cmd.Flags().GetInt("output")

		mode := &serial.Mode{
			BaudRate: 115200,
		}
		port, err := serial.Open("/dev/tty.usbmodem3871397E34321", mode)
		if err != nil {
			log.Fatal(err)
		} else {
			port.Write([]byte("SET_CLED_ON\r"))
		}
		scanner := bufio.NewScanner(port)

		capturedSamples := make([]Sample, 0, samples)

		// time estimate
		var totalSeconds float64 = float64(samples-1) * float64(interval) + (float64(samples) * float64(0.85))
		print(fmt.Sprintf("Capturing %d samples, at a %s interval\r\nEstimated time until completetion: %s\r\nCapture should be completed at around %s\r\n\r\n", samples, timeString(int64(interval)), timeString(int64(totalSeconds)), time.Unix((time.Now().Unix() + int64(totalSeconds)), 0).Format("15:04:05")))

		bar := progressbar.New(samples)
		for i := range samples {
			var sample Sample
			sample.timestamp = time.Now().Unix()
			port.Write([]byte("GET_TEMPERATURE\r"))
			scanner.Scan()
			line := scanner.Text()
			value, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
			if err != nil {
				sample.value = -999
			}
			sample.value = value
			capturedSamples = append(capturedSamples, sample)
			bar.Add(1)
			
			// print(fmt.Sprintf("Capturing samples: %.0f%% [%d/%d]\r\n", (float64(i+1)/float64(samples))*100, i+1, samples))
			if (i+1 < samples) {
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}

		port.Write([]byte("SET_CLED_OFF\r"))

		for i, value := range capturedSamples {
			print(fmt.Sprintf("Slice index: %d --- Value: %f --- Timestamp: %d\r\n", i+1, value.value, value.timestamp))
		}
	},
}

func init() {
	rootCmd.AddCommand(captureCmd)
	captureCmd.Flags().StringP("port", "p", "", "Serial port")
	captureCmd.Flags().IntP("samples", "n", 100, "Number of samples")
	captureCmd.Flags().IntP("interval", "i", 10, "Interval between samples in seconds")
	captureCmd.Flags().StringP("output", "o", "", "Output path")
}
