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
	_ "embed"
	"os"

	"github.com/spf13/cobra"
)

var desclaimer = `
---------------------------------------------------------------------
ODSCI CLI - Copyright (C) 2026  Stratos Thivaios
This program comes with ABSOLUTELY NO WARRANTY; for details type 'show w'.
This is free software, and you are welcome to redistribute it under certain
conditions; type 'show c' for details.
---------------------------------------------------------------------
`

var rootCmd = &cobra.Command{
	Use:   "odsci",
	Short: "The ODSCI CLI",
	Long: `The ODSCI CLI is the official command line application
used to interface with ODSCI probes. Learn more about ODSCI
at https://github.com/sthivaios/odsci.
 
Begin by connecting your probe to your computer.

`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetUsageTemplate(rootCmd.UsageTemplate() + desclaimer)
}


