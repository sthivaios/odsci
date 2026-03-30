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
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
		
		serialPort, _ := cmd.Flags().GetString("port")
		watch, _ := cmd.Flags().GetBool("watch");
		interval, _ := cmd.Flags().GetInt("interval")
		
		// make sure the user did enter a serial port
		if (serialPort == "") {
			fmt.Println(color.HiRedString("No serial port specified! Use --port to specify a port."))
		}

		var sigChan = make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		mode := &serial.Mode{
			BaudRate: 115200,
		}
		port, err := serial.Open(serialPort, mode)
		if err != nil {
			var errorString string = color.HiRedString("\r\nThere was an error while trying to connect to the ODSCI probe.\r\nThe serial port you entered may be incorrect.\r\nTo scan for serial ports on your computer, run ") + color.HiMagentaString("'odsci scan'") + color.HiRedString(".\r\n\r\nError details:\r\n\r\n")
			print(errorString)
			log.Fatal(err)
		} else {
			port.Write([]byte("SET_CLED_ON\r"))
		}

		go func() {
			<-sigChan
			port.Write([]byte("SET_CLED_OFF\r"))
			color.HiRed("\r\n\r\nCancelled.")
			os.Exit(0)
		}()

		scanner := bufio.NewScanner(port)

		if (watch) {
			fmt.Print(color.HiBlueString("Reading ODSCI probe on %s, at a %ds interval\r\n\r\n", serialPort, interval))
			for (true) {
				port.Write([]byte("GET_TEMPERATURE\r"))
				scanner.Scan()
				line := scanner.Text()
				value, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
				if err != nil {
					fmt.Println(color.MagentaString("Error polling sensor"))
					continue
				} else {
					fmt.Println(value);
				}
				time.Sleep(time.Duration(interval) * time.Second)
			}

		} else {
			port.Write([]byte("GET_TEMPERATURE\r"))
			scanner.Scan()
			line := scanner.Text()
			value, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
			if err != nil {
				fmt.Println(color.MagentaString("Error polling sensor"))
			} else {
				fmt.Println(value);
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().StringP("port", "p", "", "Serial port of the ODSCI probe")
	readCmd.Flags().BoolP("watch", "w", false, `Continuously watch the temperature reading`)
	readCmd.Flags().IntP("interval", "i", 1, "Interval for watching, if using the --watch flag")
}
