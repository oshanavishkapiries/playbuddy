package cli

import (
	mainview "github.com/oshanavishkapiries/playbuddy/internal/cli/views/mainView"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "playbuddy",
	Short: "PlayBuddy Torrent Downloader",
	Run: func(cmd *cobra.Command, args []string) {
		mainview.ShowMainView()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
