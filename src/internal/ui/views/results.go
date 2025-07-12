package views

import (
	"fmt"

	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
	"github.com/oshanavishkapiries/playbuddy/src/internal/ui/sharedstyles"
)

func ViewResults(results []models.Torrent, cursor int, truncateString func(string, int) string) string {
	s := ""
	if len(results) == 0 {
		s += sharedstyles.ResultStyle.Render("No results found") + "\n"
	} else {
		for i, torrent := range results {
			icon := "[🧲]"
			if cursor == i {
				icon = "[➡️]"
			}
			row := fmt.Sprintf("%s %s (%s)", icon, truncateString(torrent.Name, 45), torrent.Size)
			if cursor == i {
				s += sharedstyles.SelectedResultStyle.Render(row) + "\n"
			} else {
				s += sharedstyles.ResultStyle.Render(row) + "\n"
			}
		}
	}
	return s
}
