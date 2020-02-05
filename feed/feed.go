package feed

import (
	"time"

	"github.com/azusa0127/alertubc/scraper"
	"github.com/gorilla/feeds"
)

var lastRSS string
var lastUpdated = time.Now()

func GenerateRSS(alerts []*scraper.UBCAlertMessage) (string, error) {
	if alerts == nil || len(alerts) == 0 {
		return (&feeds.Feed{
			Title:       "UBC Campus Notifications",
			Link:        &feeds.Link{Href: "https://www.ubc.ca/campus-notifications/"},
			Description: "No alert available at the moment",
			Created:     lastUpdated,
		}).ToRss()
	}

	if alerts[0].Time == lastUpdated {
		return lastRSS, nil
	}

	var err error
	channel := &feeds.Feed{
		Title:       "UBC Campus Notifications",
		Link:        &feeds.Link{Href: "https://www.ubc.ca/campus-notifications/"},
		Description: "Feed for UBC Campus Notifications",
		Items:       []*feeds.Item{},
	}

	for _, a := range alerts {
		channel.Items = append(channel.Items, &feeds.Item{
			Title:       a.Title,
			Link:        &feeds.Link{Href: "https://www.ubc.ca/campus-notifications/"},
			Description: a.Message,
			Created:     a.Time})

		if lastUpdated.Before(a.Time) {
			lastUpdated = a.Time
		}
		if channel.Created.After(a.Time) {
			channel.Created = a.Time
		}
	}

	channel.Updated = lastUpdated
	lastRSS, err = channel.ToRss()
	return lastRSS, err
}
