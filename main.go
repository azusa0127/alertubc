package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/azusa0127/alertubc/feed"
	"github.com/azusa0127/alertubc/scraper"
)

var err error
var rss string

func scrapDaemon() {
	alert, err := scraper.ScrapeUBCAlert()
	if err == nil {
		rss, err = feed.GenerateRSS(alert)
	}

	for range time.Tick(time.Hour) {
		alert, err = scraper.ScrapeUBCAlert()
		if err == nil {
			rss, err = feed.GenerateRSS(alert)
		}
	}
}

func main() {
	go scrapDaemon()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("[500]", err.Error())
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(rss))
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
	fmt.Println("AlertUBC is now listening at port 8080")
}
