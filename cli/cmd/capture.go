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
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/sthivaios/odsci/utils"
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

		// extract flags/args
		serialPort, _ := cmd.Flags().GetString("port")
        samples, _ := cmd.Flags().GetInt("samples")
        interval, _ := cmd.Flags().GetInt("interval")
        output, _ := cmd.Flags().GetString("output")
		iso8601, _ := cmd.Flags().GetBool("iso-8601")

		// exit gracefully on ctrl+c
		var sigChan = make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// configure the serial port
		mode := &serial.Mode{
			BaudRate: 115200,
		}
		port, err := serial.Open(serialPort, mode)
		if err != nil {
			var errorString string = color.HiRedString("\r\nThere was an error while trying to connect to the ODSCI probe.\r\nThe serial port you entered may be incorrect." + color.HiRedString("\r\n\r\nError details:\r\n\r\n"))
			print(errorString)
			log.Fatal(err)
		}

		// handle ctrl+c
		go func() {
			<-sigChan
			port.Write([]byte("SET_CLED_OFF\r"))
			color.HiRed("\r\n\r\nCancelled by user via Ctrl+C")
			os.Exit(0)
		}()

		// new scanner for the serial port
		scanner := bufio.NewScanner(port)

		// clear serial buffer
		utils.ClearBuffer(port, scanner);

		// get board info, turn on CLED and check for iwdg reset
		boardInfo, _ := utils.BoardCheck(port, scanner);
		if (boardInfo.CledIsUsedForErrors == true) {
			port.Write([]byte("SET_CLED_ON\r"))
		}
		if (boardInfo.LastResetWasIWDG) {
			print(utils.AdvisoryStringIWDG(boardInfo))
		}

		capturedSamples := make([]utils.Sample, 0, samples)

		// time estimate
		var totalSeconds float64 = float64(samples-1) * float64(interval) + (float64(samples) * float64(0.85))
		print(fmt.Sprintf("\r\nCapturing %d samples, at a %s interval\r\nEstimated time until completion: %s\r\nCapture should be completed at around %s\r\n\r\n", samples, utils.TimeString(int64(interval)), utils.TimeString(int64(totalSeconds)), time.Unix((time.Now().Unix() + int64(totalSeconds)), 0).Format("15:04:05")))

		// new progress bar
		bar := progressbar.New(samples)

		var crcAdvisoryDisplayed bool = false;
		
		for i := range samples {
			// read and put values in the struct
			var sample utils.Sample

			// chose timestamp type depending on the iso8601 flag
			if (!iso8601) {
				sample.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
			} else {
				sample.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")
			}

			_, value, readError := utils.ReadTemperature(port, scanner);

			var errorString string;
			timestamp := time.Now().UTC().Format("15:04:05")
			if (readError != nil) {
				if (readError.Error() == "CRC") {
					errorString = color.HiRedString("CRC error, perhaps your sensor line is noisy?")
				} else if (readError.Error() == "PARSE") {
					errorString = color.HiRedString("Error parsing the temperature...")
				}
				if (!crcAdvisoryDisplayed) {
					fmt.Print(utils.AdvisoryStringCRC(boardInfo));
					crcAdvisoryDisplayed = true;
				}
			}

			if (readError != nil) {
				bar.Describe(fmt.Sprintf("[%s]: %s",timestamp, errorString))
			}

			sample.Value = value

			sample.ValueInFahrenheit = utils.ConvertCelsiusToFahrenheit(value)
			sample.ValueInKelvin = utils.ConvertCelsiusToKelvin(value)

			// append the struct to the samples
			capturedSamples = append(capturedSamples, sample)

			// advance the progress bar
			bar.Add(1)
			
			// print(fmt.Sprintf("Capturing samples: %.0f%% [%d/%d]\r\n", (float64(i+1)/float64(samples))*100, i+1, samples))
			if (i+1 < samples) {
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}

		// write csv
		f, err := os.Create(output)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		writer := csv.NewWriter(f)
		defer writer.Flush()

		// header row
		writer.Write([]string{"timestamp", "temperature_c", "temperature_f", "temperature_k"})

		for _, sample := range capturedSamples {
			writer.Write([]string{
				sample.Timestamp,
				strconv.FormatFloat(sample.Value, 'f', 2, 64),
				strconv.FormatFloat(sample.ValueInFahrenheit, 'f', 2, 64),
				strconv.FormatFloat(sample.ValueInKelvin, 'f', 2, 64),
			})
		}

		// turn off CLED
		if (boardInfo.CledIsUsedForErrors == true) {
			port.Write([]byte("SET_CLED_OFF\r"))
		}
	},
}


func init() {
	rootCmd.AddCommand(captureCmd)
	captureCmd.Flags().StringP("port", "p", "", "Serial port")
	captureCmd.MarkFlagRequired("port")
	captureCmd.Flags().IntP("samples", "n", 100, "Number of samples")
	captureCmd.Flags().IntP("interval", "i", 10, "Interval between samples in seconds")
	captureCmd.Flags().StringP("output", "o", "", "Output path")
	captureCmd.MarkFlagRequired("output")
	captureCmd.Flags().Bool("iso-8601", false, "Uses ISO 8601 timestamps in the CSV instead")
}
