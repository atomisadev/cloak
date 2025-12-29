package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/atomisadev/cloak/pkg/injector"
	"github.com/atomisadev/cloak/pkg/store"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run -- [COMMAND]",
	Short: "Run a command with secrets injected",
	Long:  `Decrypts the secret store and injects environment variables into the specified command`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		masterKey := os.Getenv("CLOAK_MASTER_KEY")
		if masterKey == "" {
			color.Red("Error: CLOAK_MASTER_KEY is not set.")
			os.Exit(1)
		}

		secrets, err := store.Load("cloak.encrypted", masterKey)
		if err != nil {
			color.Red("Failed to load store: %v", err)
			os.Exit(1)
		}

		cyan := color.New(color.FgCyan, color.Bold)
		cyan.Printf("[CLOAK] Injecting %d secrets into %s\n", len(secrets), strings.Join(args, " "))

		if err := injector.RunCommand(args, secrets); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			color.Red("Command execution failed: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	runCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(runCmd)
}
