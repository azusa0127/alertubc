package feed

import (
	"time"

	"github.com/azusa0127/alertubc/scraper"
	"github.com/gorilla/feeds"
)

var lastRSS string
var lastUpdated = time.Now()

func GenerateRSS(a *scraper.UBCAlertMessage) (string, error) {
	var err error
	if a != nil {
		if a.Time == lastUpdated {
			return lastRSS, nil
		}
		channel := &feeds.Feed{
			Title:       "UBC Campus Notifications",
			Link:        &feeds.Link{Href: "https://www.ubc.ca/campus-notifications/"},
			Description: a.Category,
			Created:     a.Time,
			Items: []*feeds.Item{
				&feeds.Item{
					Title:       a.Title,
					Link:        &feeds.Link{Href: "https://www.ubc.ca/campus-notifications/"},
					Description: a.Message,
					Created:     a.Time},
			},
		}

		lastUpdated = a.Time
		lastRSS, err = channel.ToRss()
		return lastRSS, err
	}
	return (&feeds.Feed{
		Title:       "UBC Campus Notifications",
		Link:        &feeds.Link{Href: "https://www.ubc.ca/campus-notifications/"},
		Description: "No alert available at the moment",
		Created:     lastUpdated,
	}).ToRss()
}
