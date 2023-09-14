package parsers

// To handle weird cases.

import (
	"encoding/json"
	"log"
	"net/url"
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Barby is annoying, can also get YT url but useless in my opinion
func getBarbyDescription(event_url string, desc_sel string, image_sel string) (string, string, error) {
	res, err := GetRequest(event_url)
	defer res.Body.Close()
	if err != nil {
		log.Println("GET request error", err)
		return "", "", err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	image := doc.Find(image_sel).First()

	parsed_url, err := url.Parse(event_url)
	if err != nil {
		log.Fatal(err)
	}

	res, err = GetRequest("http://tickets.barby.co.il/api/shows/d-y/" + parsed_url.Query().Get("id2"))
	defer res.Body.Close()
	if err != nil {
		log.Println("GET request error", err)
		return "", "", err
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

	return desc, GetImageSource(image), nil
}

func GetDescriptionHook(event_url string, website string, desc_sel string, image_sel string) (string, string, error) {
	switch website {
	case "https://www.barby.co.il":
		event_url, config, err := getBarbyDescription(event_url, desc_sel, image_sel)
		if err != nil {
			log.Println("Error handling barby hook", err)
		}

		return event_url, config, err
	default:
		return "", "", errors.New("Hook doesn't exist")
	}
}
