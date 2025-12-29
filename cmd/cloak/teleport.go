package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/psanford/wormhole-william/wormhole"
	"github.com/spf13/cobra"
)

var teleportCmd = &cobra.Command{
	Use:   "teleport [CODE]",
	Short: "Securely send/receive Master Key via Magic Wormhole",
	Long: `Transmit your Master Key directly to another device using PAKE.

SENDER (No args):
  Generates a code. Keeps the connection open until a receiver connects.

RECEIVER (With code):
  Connects using the code and receives the key.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if len(args) > 0 {
			code := args[0]
			receiveKey(ctx, code)
			return
		}

		sendKey(ctx)
	},
}

func sendKey(ctx context.Context) {
	masterKey := os.Getenv("CLOAK_MASTER_KEY")
	if masterKey == "" {
		color.Red("Error: CLOAK_MASTER_KEY is not set. Nothing to teleport.")
		os.Exit(1)
	}

	var c wormhole.Client

	code, status, err := c.SendText(ctx, masterKey)
	if err != nil {
		color.Red("Failed to initialize wormhole: %v", err)
		os.Exit(1)
	}

	fmt.Println("Wormhole Open. Share this code with the receiver:")
	fmt.Println()

	codeStyle := color.New(color.FgGreen, color.Bold)
	codeStyle.Printf("   %s\n", code)

	fmt.Println()
	fmt.Println("Waiting for receiver... (Ctrl+C to cancel)")

	result := <-status
	if result.Error != nil {
		color.Red("\nTransfer failed: %v", result.Error)
		os.Exit(1)
	} else if result.OK {
		color.Green("\n✔ Teleport successful.")
	}
}

func receiveKey(ctx context.Context, code string) {
	var c wormhole.Client

	msg, err := c.Receive(ctx, code)
	if err != nil {
		color.Red("Failed to receive from wormhole: %v", err)
		os.Exit(1)
	}

	data, err := io.ReadAll(msg)
	if err != nil {
		color.Red("Failed to read message: %v", err)
		os.Exit(1)
	}

	fmt.Print(string(data))
}

func init() {
	rootCmd.AddCommand(teleportCmd)
}
