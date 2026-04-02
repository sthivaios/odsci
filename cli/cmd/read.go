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
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sthivaios/odsci/utils"
	"go.bug.st/serial"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read the temperature directly in the terminal",
	Long: `Unlike the 'capture' command which saves the output as a file,
the 'read' command prints the temperature readings directly into your terminal.

You can either print the temperature just once, or use the --watch
flag, to continuously update the reading in your terminal.

The command accepts other arguments too.`,
	Run: func(cmd *cobra.Command, args []string) {
		
		// extract flags/args
		serialPort, _ := cmd.Flags().GetString("port")
		watch, _ := cmd.Flags().GetBool("watch");
		noLog, _ := cmd.Flags().GetBool("no-log");
		interval, _ := cmd.Flags().GetInt("interval")
		unit, _ := cmd.Flags().GetString("unit")

		if (unit != "c" && unit != "f" && unit != "k") {
			print(color.MagentaString("The --unit (-u) flag only accepts: 'c', 'f', 'k',\nbut you entered \"%s\".\r\n\r\nReminder that the default is 'c'.\r\n", unit))
			return
		}

		// exit gracefully on ctrl+c
		var sigChan = make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// configure the serial port
		mode := &serial.Mode{
			BaudRate: 115200,
		}
		port, err := serial.Open(serialPort, mode)
		if err != nil {
			var errorString string = color.HiRedString("\r\nThere was an error while trying to connect to the ODSCI probe.\r\nThe serial port you entered may be incorrect.\r\nTo scan for serial ports on your computer, run ") + color.HiMagentaString("'odsci scan'") + color.HiRedString(".\r\n\r\nError details:\r\n\r\n")
			print(errorString)
			log.Fatal(err)
		}

		// handle ctrl+c
		go func() {
			<-sigChan
			color.HiRed("\r\n\r\nCancelled.")
			os.Exit(0)
		}()

		// new scanner for the serial port
		scanner := bufio.NewScanner(port)

		// clear serial buffer
		utils.ClearBuffer(port, scanner);

		// main logic
		if (watch) {
			fmt.Print(color.HiBlueString("Reading ODSCI probe on %s, at a %ds interval\r\n\r\n", serialPort, interval))
			if (noLog) {
				fmt.Print(color.HiYellowString("You are using the \"--no-log\" flag. If the CLI looks like it has frozen, it hasn't.\r\nThe temperature is just not updating.\r\n\r\n"))
			}
			_, raw_temp := utils.ReadTemperature(port, scanner)
			var temp_to_print string
			switch unit {
				case "c":
					temp_to_print = fmt.Sprintf("%0.2f", raw_temp)
				case "f":
					temp_to_print = fmt.Sprintf("%0.2f", utils.ConvertCelsiusToFahrenheit(raw_temp))
				case "k":
					temp_to_print = fmt.Sprintf("%0.2f", utils.ConvertCelsiusToKelvin(raw_temp))
			}
			for (true) {
				if (!noLog) {
					timestamp := time.Now().UTC().Format("15:04:05")
					fmt.Printf("[%s]: %s\r\n",timestamp, temp_to_print)
				} else {
					timestamp := time.Now().UTC().Format("15:04:05")
					fmt.Printf("\r[%s]: %-10s",timestamp, temp_to_print)
				}
				time.Sleep(time.Duration(interval) * time.Second)
			}
		} else {
			fmt.Print(color.HiBlueString("Reading ODSCI probe on %s\r\n\r\n", serialPort))
			fmt.Println(utils.ReadTemperature(port, scanner));
		}

	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().StringP("port", "p", "", "Serial port of the ODSCI probe")
	readCmd.MarkFlagRequired("port")
	readCmd.Flags().BoolP("watch", "w", false, `Continuously watch the temperature reading`)
	readCmd.Flags().IntP("interval", "i", 1, "Interval for watching, if using the --watch flag")
	readCmd.Flags().Bool("no-log", false, "No log: applies only if --watch is being used and does't log previous values")
	readCmd.Flags().StringP("unit", "u", "c", "Select a unit: 'c' for Celsius, 'f' for Fahrenheit or 'k' for Kelvin")
}
