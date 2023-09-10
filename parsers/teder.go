package parsers

import (
	"time"

	"github.com/mmcdole/gofeed"
)

func GetTederChannel(url string, refresh time.Duration) (chan Show, chan bool) {
	show := make(chan Show)
	stop := make(chan bool)
	ticker := time.NewTicker(refresh)

	go func() {
		for {
			select {
			case <-stop:
				return
			case tick := <-ticker.C:

			}
		}
	}()

	return show, stop
}
