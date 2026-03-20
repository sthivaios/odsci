/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
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
