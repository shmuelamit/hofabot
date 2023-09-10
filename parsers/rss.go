package parsers

import (
	"time"

	"github.com/mmcdole/gofeed"
)

func GetRSSChannel(url string, refresh time.Duration) (chan Show, chan bool) {
	show := make(chan Show)
	stop := make(chan bool)
	ticker := time.NewTicker(refresh)

	go func() {
		for {
			select {
			case <-stop:
				return
			case tick := <-ticker.C:
				fp := gofeed.NewParser()
				feed, err := fp.ParseURL(url)
				if err != nil {
					panic("gofeed failed to parse url")
				}

				for _, item := range feed.Items {
					if item.PublishedParsed.After(tick.Add(-refresh)) {
						show <- Show{
							name: item.Title,
							url:  item.Link,
							desc: item.Content,
						}
					}
				}
			}
		}
	}()

	return show, stop
}
