package cli

import (
	"github.com/oshanavishkapiries/playbuddy/internal/cli/views"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "playbuddy",
	Short: "PlayBuddy Torrent Downloader",
	Run: func(cmd *cobra.Command, args []string) {
		views.ShowHomeView()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
