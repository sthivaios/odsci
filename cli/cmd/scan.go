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
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.bug.st/serial/enumerator"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Provides an easy way to list all USB serial ports",
	Long: `Scans the computer for any USB CDC serial ports,
to help the user identify the port which the ODSCI probe
is connected to.`,

	Run: func(cmd *cobra.Command, args []string) {
		ports, err := enumerator.GetDetailedPortsList()
		if err != nil {
			color.HiRed(`
There was an error while attempting to list serial ports...
Perhaps this is a permission error?

Error details:

`)
			log.Fatal(err)
		}
		if len(ports) == 0 {
			fmt.Println("No serial ports found!")
			return
		}
		for _, port := range ports {
			if port.IsUSB {
				fmt.Printf("\r\nFound port: ")
				color.Magenta(port.Name);
				fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
				fmt.Printf("   USB serial %s\n", port.SerialNumber)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
