package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/spf13/cobra"
)

// twitchCmd represents the twitch command
var twitchCmd = &cobra.Command{
	Use:   "twitch [channel]",
	Short: "Connects to a Twitch chat and displays messages in the terminal",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		channel := args[0]
		connectTwitchChat(channel)
	},
}

// Connect to Twitch chat
func connectTwitchChat(channel string) {
	client := twitch.NewAnonymousClient() // Connect anonymously (no login required)

	client.OnConnect(func() {
		fmt.Println(color.GreenString("✅ Connected to Twitch Chat!"))
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Printf(
			color.GreenString("%s: ")+color.WhiteString("%s\n"),
			message.User.DisplayName, message.Message,
		)
	})

	client.Join(channel)

	err := client.Connect()
	if err != nil {
		fmt.Println(color.RedString("❌ Error connecting to Twitch chat: %v"), err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(twitchCmd)
}
