package main

import (
	"fmt"
	"os"

	"github.com/atomisadev/cloak/pkg/crypto"
	"github.com/atomisadev/cloak/pkg/keychain"
	"github.com/atomisadev/cloak/pkg/store"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a new encrypted secret store",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("cloak.encrypted"); err == nil {
			color.Red("Error: 'cloak.encrypted' already exists. Aborting to prevent overwrite.")
			os.Exit(1)
		}

		masterKey, err := crypto.GenerateKey()
		if err != nil {
			fmt.Printf("Failed to generate key: %v\n", err)
			os.Exit(1)
		}

		emptyStore := make(map[string]string)
		if err := store.Save("cloak.encrypted", emptyStore, masterKey); err != nil {
			fmt.Printf("Failed to write store: %v\n", err)
			os.Exit(1)
		}

		color.Green("✔ Store initialized successfully.")
		fmt.Println("Here is your MASTER KEY. Save it securely!")
		fmt.Println()

		keyStyle := color.New(color.FgGreen, color.Bold)
		keyStyle.Println(masterKey)
		fmt.Println()

		err = keychain.Save(masterKey)
		if err == nil {
			color.Cyan("Master key saved to System Keychain.")
			color.New(color.FgHiBlack).Println("(You don't need to set env vars manually)")
		} else {
			color.Yellow("⚠ Could not save to Keychain: %v", err)
			fmt.Println("Set it manually: export CLOAK_MASTER_KEY=...")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
