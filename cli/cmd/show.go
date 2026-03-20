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

	"github.com/spf13/cobra"
)

//go:embed LICENSE
var licenseText string

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows information about GPLv3",
	Long: `Shows information about the The GNU General Public License v3.0`,
}

var showWCmd = &cobra.Command{
	Use:   "w",
	Short: "Shows the warranty desclaimer of the GPLv3",
	Long: `Shows the warranty desclaimer of
the The GNU General Public License v3.0`,
	
	Run: func(cmd *cobra.Command, args []string) {
		print(`From section "15. Disclaimer of Warranty" of
The GNU General Public License v3.0 terms:

THERE IS NO WARRANTY FOR THE PROGRAM, TO THE EXTENT PERMITTED BY
APPLICABLE LAW. EXCEPT WHEN OTHERWISE STATED IN WRITING THE COPYRIGHT
HOLDERS AND/OR OTHER PARTIES PROVIDE THE PROGRAM "AS IS" WITHOUT WARRANTY
OF ANY KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING, BUT NOT LIMITED TO,
THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
PURPOSE. THE ENTIRE RISK AS TO THE QUALITY AND PERFORMANCE OF THE PROGRAM
IS WITH YOU. SHOULD THE PROGRAM PROVE DEFECTIVE, YOU ASSUME THE COST OF
ALL NECESSARY SERVICING, REPAIR OR CORRECTION.
`)
	},
}

var showCCmd = &cobra.Command{
	Use:   "c",
	Short: "Shows the license terms of the GPLv3",
	Long: `Shows the full license terms of The GNU General Public License v3.0`,
	
	Run: func(cmd *cobra.Command, args []string) {
		print(licenseText);
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.AddCommand(showWCmd)
	showCmd.AddCommand(showCCmd)
}
