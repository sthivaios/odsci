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
	"log"

	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

// captureCmd represents the capture command
var captureCmd = &cobra.Command{
	Use: "capture",
	Short: "Used to capture a number of samples from the probe",
	Long: `The capture command is used to capture a given number
of samples from the probe, at a specific interval, and export
the captured data.`,

	Run: func(cmd *cobra.Command, args []string) {
		// port, _ := cmd.Flags().GetString("port")
        // samples, _ := cmd.Flags().GetInt("samples")
        // interval, _ := cmd.Flags().GetInt("interval")
        // output, _ := cmd.Flags().GetInt("output")

		mode := &serial.Mode{
			BaudRate: 115200,
		}
		_, err := serial.Open("/dev/tty.usbmodem3871397E34321", mode)
		if err != nil {
			log.Fatal(err)
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
