package main

import (
	"log"
	"net/http"
	"time"

	"github.com/azusa0127/alertubc/feed"
	"github.com/azusa0127/alertubc/scraper"
)

var err error
var rss string

func scrapDaemon() {
	alerts, err := scraper.ScrapeUBCAlert()
	if err == nil {
		rss, err = feed.GenerateRSS(alerts)
	}

	for range time.Tick(time.Hour) {
		alerts, err = scraper.ScrapeUBCAlert()
		if err == nil {
			rss, err = feed.GenerateRSS(alerts)
		} else {
			log.Println("[Error]", err.Error(), "@ScrapeUBCAlert()")
		}
	}
}

func main() {
	go scrapDaemon()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(rss))
		}
	})
	log.Println("[Info] AlertUBC is now listening at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
