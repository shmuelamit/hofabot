package parsers

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

func getRSSImage(event_url string, config RSSConfig) string {
	// Sometimes websites do stupid stuff and I handle it stupidly as well
	if _, image, hooked := GetDescriptionHook(event_url, config.Url, "", config.Image); hooked {
		return image
	}

	res := GetRequest(event_url)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	image := doc.Find(config.Image).First()

	return GetImageSource(image)
}

func GetRSSChannel(config RSSConfig) (chan Show, chan bool) {
	show := make(chan Show)
	stop := make(chan bool)
	ticker := time.NewTicker(config.Refresh)

	go func() {
		for {
			println("TICK", config.Url)
			select {
			case <-stop:
				return
			case tick := <-ticker.C:
				fp := gofeed.NewParser()
				feed, err := fp.ParseURL(config.Url)
				if err != nil {
					log.Fatal("gofeed failed to parse url")
				}

				for _, item := range feed.Items {
					if item.PublishedParsed.After(tick.Add(-config.Refresh)) {
						img := ""

						if config.Image != "" {
							img = getRSSImage(item.Link, config)
						}

						content, err := goquery.NewDocumentFromReader(strings.NewReader(item.Content))
						if err != nil {
							log.Fatal(err)
						}

						show <- Show{
							Name:  item.Title,
							Url:   item.Link,
							Desc:  HTMLToText(content.Selection),
							Image: img,
						}
					}
				}
			}
		}
	}()

	return show, stop
}
