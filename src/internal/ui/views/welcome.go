package views

import (
	"github.com/oshanavishkapiries/playbuddy/src/internal/ui/sharedstyles"
)

func ViewWelcome() string {
	s := sharedstyles.AsciiStyle.Render(sharedstyles.AsciiArt) + "\n"
	s += sharedstyles.VersionStyle.Render(sharedstyles.Version) + "\n\n"
	s += sharedstyles.HelpStyle.Render("Press Enter to continue or Q to quit")
	return s
}
