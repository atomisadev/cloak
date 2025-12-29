package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cloak",
	Short: "Secure environment injector",
	Long: `
	CLOAK // SECURE SECRET MANAGEMENT SYSTEM
	Injects encrypted secrets into child processes without writing to disk.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		printWelcome()
		cmd.Help()
	},
}

func printWelcome() {
	cyan := color.New(color.FgCyan, color.Bold)

	banner := `
   ______   __    ____    ___    __ __
  / ____/  / /   / __ \  /   |  / //_/
 / /      / /   / / / / / /| | / ,<
/ /___   / /___/ /_/ / / ___ |/ /| |
\____/  /_____/\____/ /_/  |_/_/ |_|

:: CLOAK PROTOCOL :: [ONLINE]
`

	cyan.Print(banner)
	fmt.Println()
}

func init() {

}
