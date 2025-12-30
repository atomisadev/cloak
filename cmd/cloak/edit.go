package main

import (
	"fmt"
	"os"

	"github.com/atomisadev/cloak/internal/ui"
	"github.com/atomisadev/cloak/pkg/store"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit secrets in a TUI",
	Long:  `Launch the interactive spreadsheet editor for your secrets.`,
	Run: func(cmd *cobra.Command, args []string) {
		masterKey := RequireKey()

		secrets, err := store.Load("cloak.encrypted", masterKey)
		if err != nil {
			color.Red("Failed to load store: %v", err)
			os.Exit(1)
		}

		p := tea.NewProgram(ui.InitialModel(secrets), tea.WithAltScreen())

		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}

		m, ok := finalModel.(ui.Model)
		if !ok {
			return
		}

		if m.ToSave != nil {
			if err := store.Save("cloak.encrypted", m.ToSave, masterKey); err != nil {
				color.Red("Failed to save store: %v", err)
				os.Exit(1)
			}
			color.Green("✔ Vault updated securely.")
		} else {
			color.New(color.FgHiBlack).Println("No changes saved.")
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
