package parsers

// To handle weird cases.

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Barby is annoying, can also get YT url but useless in my opinion
func getBarbyDescription(event_url string, config GenericConfig) (string, string) {
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

	parsed_url, err := url.Parse(event_url)
	if err != nil {
		log.Fatal(err)
	}

	res = GetRequest("http://tickets.barby.co.il/api/shows/d-y/" + parsed_url.Query().Get("id2"))
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", GetImageSource(image)
	}

	var barby_desc map[string]string

	err = json.NewDecoder(res.Body).Decode(&barby_desc)
	if err != nil {
		log.Fatal(err)
	}

	doc, err = goquery.NewDocumentFromReader(strings.NewReader(barby_desc["description"]))
	if err != nil {
		log.Fatal(err)
	}
	desc := HTMLToText(doc.Selection)

	return desc, GetImageSource(image)
}

func GenericDescriptionHook(event_url string, config GenericConfig) (string, string, bool) {
	switch config.Url {
	case "https://www.barby.co.il":
		event_url, config := getBarbyDescription(event_url, config)
		return event_url, config, true
	default:
		return "", "", false
	}
}
