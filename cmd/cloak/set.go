package main

import (
	"os"

	"github.com/atomisadev/cloak/pkg/store"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [KEY] [VALUE]",
	Short: "Add or update a secret",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		masterKey := RequireKey()

		secrets, err := store.Load("cloak.encrypted", masterKey)
		if err != nil {
			color.Red("Failed to load store: %v", err)
			os.Exit(1)
		}

		secrets[key] = value

		if err := store.Save("cloak.encrypted", secrets, masterKey); err != nil {
			color.Red("Failed to save store: %v", err)
			os.Exit(1)
		}

		color.Cyan("✔ Set %s", key)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
