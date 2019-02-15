package feed

import (
	"time"

	"github.com/azusa0127/alertubc/scraper"
	"github.com/gorilla/feeds"
)

var lastUpdated = time.Now()

func GenerateRSS(a *scraper.UBCAlertMessage) (string, error) {
	channel := &feeds.Feed{
		Title:       "UBC Campus Notifications",
		Link:        &feeds.Link{Href: "https://www.ubc.ca/campus-notifications/"},
		Description: "No alert available at the moment",
		Created:     lastUpdated,
	}

	if a != nil {
		channel.Description = a.Category
		channel.Created = a.Time
		channel.Items = []*feeds.Item{
			&feeds.Item{
				Title:       a.Title,
				Link:        &feeds.Link{Href: "https://www.ubc.ca/campus-notifications/"},
				Description: a.Message,
				Created:     a.Time,
			},
		}

		lastUpdated = a.Time
	}
	return channel.ToRss()
}
