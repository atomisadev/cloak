package main

import (
	"fmt"
	"os"

	"github.com/atomisadev/cloak/pkg/keychain"
	"github.com/fatih/color"
)

func RequireKey() string {
	if envKey := os.Getenv("CLOAK_MASTER_KEY"); envKey != "" {
		return envKey
	}

	wd, err := os.Getwd()
	if err == nil {
		if key, err := keychain.Get(wd); err == nil && key != "" {
			return key
		}
	}

	color.Red("✖ Error: Master Key not found.")
	color.New(color.FgHiBlack).Println("  Cloak cannot decrypt your secrets without the key.")
	fmt.Println()
	color.Yellow("  Solution 1: Run 'cloak init' to generate and save a key.")
	color.Yellow("  Solution 2: Set 'export CLOAK_MASTER_KEY=...' manually.")

	os.Exit(1)
	return ""
}
