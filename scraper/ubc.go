package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	ubcURL = "https://www.ubc.ca"
)

type UBCAlertMessage struct {
	Category string
	Title    string
	Message  string
	Time     time.Time
}

func ScrapeUBCAlert() (rv *UBCAlertMessage, err error) {
	res, err := http.Get(ubcURL)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	defer res.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}
	if alertContent := doc.Find("div.alert-content").First(); alertContent != nil {
		rv = &UBCAlertMessage{}
		if timeNode := alertContent.Find("div.alert-date > em"); timeNode != nil {
			timeString := timeNode.Text()[len("Updated: "):]
			if strings.HasPrefix(timeString, "Updated: ") {
				timeString = timeString[len("Updated: "):]
			}
			if rv.Time, err = time.Parse("January 02, 2006 3:04 pm MST", strings.Replace(timeString, ".", "", 2)); err != nil {
				rv.Time = time.Now()
				err = nil
			}
		}
		if messageNode := alertContent.Find("div.alert-message"); messageNode != nil {
			if spanNode := messageNode.Find("span"); spanNode != nil {
				rv.Category = spanNode.Text()
				spanNode.Remove()
			}

			if strongNode := messageNode.Find("strong"); strongNode != nil {
				rv.Title = strongNode.Text()
				strongNode.Remove()
			}

			rv.Message = messageNode.Text()
			rv.Message = rv.Message[:strings.LastIndex(rv.Message, ".")]
		}
	}
	return
}
