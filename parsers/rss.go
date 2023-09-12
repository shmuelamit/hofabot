package parsers

import (
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

func GetRSSChannel(config RSSConfig) (chan Show, chan bool) {
	show := make(chan Show)
	stop := make(chan bool)
	ticker := time.NewTicker(config.Refresh)

	go func() {
		for {
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
						show <- Show{
							Name: item.Title,
							Url:  item.Link,
							Desc: item.Content,
						}
					}
				}
			}
		}
	}()

	return show, stop
}
