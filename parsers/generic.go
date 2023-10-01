package parsers

import (
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getGenericDescription(event_url string, config GenericConfig) (string, string, error) {
	// Sometimes websites do stupid stuff and I handle it stupidly as well
	if desc, image, err := GetDescriptionHook(event_url, config.Url, config.Desc, config.Image); err == nil {
		return desc, image, nil
	}

	res, err := GetRequest(event_url)
	defer res.Body.Close()
	if err != nil {
		log.Println("GET request error", err)
		return "", "", err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("goquery parsing error", err)
		return "", "", err
	}

	image := doc.Find(config.Image).First()
	desc := HTMLToText(doc.Find(config.Desc).First())

	img_source, err := GetImageSource(image)
	if err != nil {
		return "", "", err
	}

	return desc, img_source, nil
}

func GetGenericShows(config GenericConfig) ([]Show, error) {
	res, err := GetRequest(config.Url)
	defer res.Body.Close()
	if err != nil {
		log.Println("GET request error", err)
		return nil, err
	}

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s\n", res.StatusCode, res.Status)
		return nil, errors.New("Status code error")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("failed to parse generic show body", err)
		return nil, err
	}

	shows := []Show{}

	// Can parallelize, won't
	doc.Find(config.Title).Each(func(i int, title *goquery.Selection) {
		if len(title.Text()) == 0 {
			log.Println(err)
			return
		}

		// Find relevant href
		link, exists := title.Parents().Has("a[href]").First().Find("a[href]").Attr("href")
		if !exists {
			log.Println(err)
			return
		}

		parsed_link, err := url.Parse(link)
		if err != nil {
			log.Println(err)
			return
		}

		parsed_url, err := url.Parse(config.Url)
		if err != nil {
			log.Println(err)
			return
		}

		resolved_link := parsed_url.ResolveReference(parsed_link).String()

		desc, image, err := getGenericDescription(resolved_link, config)
		if err != nil {
			log.Println(err)
			return
		}

		shows = append(shows, Show{
			Name:  title.Text(),
			Url:   resolved_link,
			Desc:  desc,
			Image: image,
		})
	})

	return shows, nil
}

func GetGenericChannel(config GenericConfig) (chan Show, chan bool) {

	show_channel := make(chan Show)
	stop := make(chan bool)
	ticker := time.NewTicker(config.Refresh)

	go func() {
		last_shows, err := GetGenericShows(config)
		if err != nil {
			log.Fatal("Failed to get initial shows", err, config)
		}

		for {
			println("TICK", config.Url)
			select {
			case <-stop:
				return
			case <-ticker.C:
				current_shows, err := GetGenericShows(config)
				if err != nil {
					log.Println("Failed to get shows, skipping", err, config)
					continue
				}

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
