package parsers

import (
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const TEDER_URL = "https://www.teder.fm/"

func getTederDescription(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return HTMLToText(doc.Find(".content-body").First())
}

func GetFutureShows() []Show {
	res, err := http.Get(TEDER_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	shows := []Show{}

	doc.Find(".item-future").Each(func(i int, s *goquery.Selection) {
		title, exists := s.Find("a").Attr("title")
		if !exists {
			log.Fatal(err)
		}

		link, exists := s.Find("a").Attr("href")
		if !exists {
			log.Fatal(err)
		}

		shows = append(shows, Show{
			Name: title,
			Url:  link,
			Desc: getTederDescription(link),
		})
	})

	return shows
}

func GetTederChannel(refresh time.Duration) (chan Show, chan bool) {

	show_channel := make(chan Show)
	stop := make(chan bool)
	ticker := time.NewTicker(refresh)

	go func() {
		last_shows := GetFutureShows()

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				current_shows := GetFutureShows()
				new_shows := ShowsDifference(current_shows, last_shows)

				for _, show := range new_shows {
					show_channel <- show
				}

				last_shows = current_shows
			}
		}
	}()

	return show_channel, stop
}
