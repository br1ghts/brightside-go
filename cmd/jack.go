package cmd

import (
	"fmt"

	"brightside-go/ui"

	"github.com/spf13/cobra"
)

// Jack CLI Command
var jackCmd = &cobra.Command{
	Use:   "jack",
	Short: "Launches the Brightside Jack AI Terminal Dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🕶️ Brightside Jack booting up...")
		ui.StartDashboard()
	},
}

func init() {
	rootCmd.AddCommand(jackCmd)
}
