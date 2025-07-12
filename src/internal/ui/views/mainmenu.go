package views

import (
	"github.com/oshanavishkapiries/playbuddy/src/internal/ui/sharedstyles"
)

func ViewMainMenu(cursor int, errorMsg string) string {
	s := sharedstyles.AsciiStyle.Render(sharedstyles.AsciiArt) + "\n"
	s += sharedstyles.VersionStyle.Render(sharedstyles.Version) + "\n\n"

	menuItems := []string{"Search Torrent", "Download Hub", "Settings"}
	for i, item := range menuItems {
		if cursor == i {
			s += sharedstyles.SelectedMenuStyle.Render("> "+item) + "\n"
		} else {
			s += sharedstyles.MenuStyle.Render("  "+item) + "\n"
		}
	}

	if errorMsg != "" {
		s += "\n" + sharedstyles.ErrorStyle.Render("Error: "+errorMsg) + "\n"
	}

	s += "\n" + sharedstyles.HelpStyle.Render("Use arrow keys to navigate, Enter to select, Esc to go back")
	return s
}
