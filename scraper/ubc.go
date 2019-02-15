package scraper

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	ubcURL = "https://www.ubc.ca"

	timeStringFormat = "January 02, 2006 3:04 pm -0700"

	requestTimeout = 10 * time.Second
)

type UBCAlertMessage struct {
	Category string
	Title    string
	Message  string
	Time     time.Time
}

// Sample ts: Updated: February 15, 2019 6:00 a.m. PST
func preprocessTimestring(ts string) (string, error) {
	timeStrings := strings.Split(strings.TrimSpace(ts), " ")

	// Strip off prefix
	if timeStrings[0] == "Updated:" {
		timeStrings = timeStrings[1:]
	}

	// Validate timestring length
	if len(timeStrings) != 6 {
		// Invalid timeString, return for parseTimeString to handle
		return "", fmt.Errorf("Invalid timeString <%s> encounted @preprocessTimestring", ts)
	}

	// Format am/pm symbol
	switch timeStrings[4] {
	case "a.m.":
		timeStrings[4] = "am"
	case "p.m.":
		timeStrings[4] = "pm"
	}

	// Format PST
	switch timeStrings[5] {
	case "PST":
		timeStrings[5] = "-0800"
	case "PDT":
		timeStrings[5] = "-0700"
	}
	return strings.Join(timeStrings, " "), nil
}

func parseTimeString(timeString string) time.Time {
	t := time.Now()
	pts, err := preprocessTimestring(timeString)
	if err == nil {
		t, err = time.Parse(timeStringFormat, pts)
	}
	if err != nil {
		log.Println("[Warning]", err.Error())
	}
	return t
}

func processMessage(message string) string {
	rv := []string{}
	for _, sentence := range strings.Split(message[:strings.LastIndex(message, ".")], ".") {
		sentence = strings.TrimSpace(sentence)
		if !strings.HasPrefix(sentence, "Due to current weather conditions, members") &&
			!strings.HasPrefix(sentence, "Drive safely and") &&
			!strings.HasPrefix(sentence, "For information on transit") {
			rv = append(rv, sentence)
		}
	}
	return strings.Join(rv, ".")
}

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}
var netClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: netTransport,
}

func ScrapeUBCAlert() (rv *UBCAlertMessage, err error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), requestTimeout)
	defer cancelCtx()
	req, _ := http.NewRequest("GET", ubcURL, nil)
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)

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
		rv = &UBCAlertMessage{Time: time.Now()}
		if timeNode := alertContent.Find("div.alert-date > em"); timeNode != nil {
			rv.Time = parseTimeString(timeNode.Text())
		}
		if messageNode := alertContent.Find("div.alert-message"); messageNode != nil {
			if spanNode := messageNode.Find("span"); spanNode != nil {
				rv.Category = spanNode.Text()
			}
			if strongNode := messageNode.Find("strong"); strongNode != nil {
				rv.Title = strongNode.Text()
			}
			rv.Message = processMessage(messageNode.Text())
		}
	}
	return
}
