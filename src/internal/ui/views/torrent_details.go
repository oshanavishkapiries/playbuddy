package views

import (
	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
	"github.com/oshanavishkapiries/playbuddy/src/internal/ui/sharedstyles"
)

func ViewTorrentDetails(t *models.Torrent) string {
	if t == nil {
		return "No torrent selected"
	}
	s := sharedstyles.TitleStyle.Render("Torrent Details") + "\n\n"

	s += sharedstyles.DetailStyle.Render(sharedstyles.DetailLabelStyle.Render("Name: ")+sharedstyles.DetailValueStyle.Render(t.Name)) + "\n"
	s += sharedstyles.DetailStyle.Render(sharedstyles.DetailLabelStyle.Render("Size: ")+sharedstyles.DetailValueStyle.Render(t.Size)) + "\n"
	s += sharedstyles.DetailStyle.Render(sharedstyles.DetailLabelStyle.Render("Category: ")+sharedstyles.DetailValueStyle.Render(t.Category)) + "\n"
	s += sharedstyles.DetailStyle.Render(sharedstyles.DetailLabelStyle.Render("Uploaded: ")+sharedstyles.DetailValueStyle.Render(t.DateUploaded)) + "\n"
	s += sharedstyles.DetailStyle.Render(sharedstyles.DetailLabelStyle.Render("Uploader: ")+sharedstyles.DetailValueStyle.Render(t.UploadedBy)) + "\n"
	s += sharedstyles.DetailStyle.Render(sharedstyles.DetailLabelStyle.Render("Seeders: ")+sharedstyles.DetailValueStyle.Render(t.Seeders)) + "\n"
	s += sharedstyles.DetailStyle.Render(sharedstyles.DetailLabelStyle.Render("Leechers: ")+sharedstyles.DetailValueStyle.Render(t.Leechers)) + "\n"
	s += sharedstyles.DetailStyle.Render(sharedstyles.DetailLabelStyle.Render("URL: ")+sharedstyles.DetailValueStyle.Render(t.Url)) + "\n"
	if t.Magnet != "" {
		s += sharedstyles.DetailStyle.Render(sharedstyles.DetailLabelStyle.Render("Magnet: ")+sharedstyles.DetailValueStyle.Render(t.Magnet)) + "\n"
	}
	return s
}
