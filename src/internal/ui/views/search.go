package views

import (
	"fmt"

	"github.com/oshanavishkapiries/playbuddy/src/internal/ui/sharedstyles"
)

func ViewSearch(input string, loading bool, loaderFrame int, loaderFrames []string, errorMsg string) string {
	s := sharedstyles.AsciiStyle.Render(sharedstyles.AsciiArt) + "\n"
	s += sharedstyles.VersionStyle.Render(sharedstyles.Version) + "\n\n"

	s += sharedstyles.InputStyle.Render("Search: "+input+"█") + "\n\n"

	if loading {
		s += sharedstyles.LoadingStyle.Render(fmt.Sprintf(" Searching... %s", loaderFrames[loaderFrame])) + "\n\n"
	}

	if errorMsg != "" {
		s += sharedstyles.ErrorStyle.Render("Error: "+errorMsg) + "\n"
	}

	s += "\n" + sharedstyles.NavigationStyle.Render("Type your search and press Enter, \u2190 Back")
	return s
}
