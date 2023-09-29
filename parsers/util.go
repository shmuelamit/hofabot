package parsers

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Gets show difference by title
func ShowsDifference(a, b []Show) []Show {
	mb := make(map[string]Show, len(b))

	for _, show := range b {
		mb[show.Name] = show
	}

	var diff []Show
	for _, x := range a {
		if _, found := mb[x.Name]; !found {
			diff = append(diff, x)
		}
	}

	return diff
}

func HTMLToText(sel *goquery.Selection) string {
	return strings.Join(sel.Children().Map(func(i int, s *goquery.Selection) string {
		return strings.Join(s.Contents().Map(func(i int, s *goquery.Selection) string {
			if !s.Is("br") {
				return strings.TrimSpace(s.Text())
			} else {
				return "\n"
			}
		}), "")

	}), "\n\n")
}

func GetRequest(web_url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", web_url, nil)
	if err != nil {
		log.Println("Failed to init GET request", err)
		return nil, err
	}

	parsed_url, err := url.Parse(web_url)
	if err != nil {
		log.Println("Failed to parse url", err)
		return nil, err
	}

	req.Header.Add("Host", parsed_url.Host)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "go")

	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send GET request", err)
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("Status Code Error")
	}

	return res, nil
}

func GetImageSource(img *goquery.Selection) (string, error) {
	if img.Is("div") {
		css, exists := img.Attr("style")
		if !exists {
			log.Println("no css attribute")
			return "", errors.New("No css attribute")
		}

		css_noprefix := css[strings.Index(css, "background-image: url(")+len("background-image: url("):]
		image_url := css_noprefix[:strings.Index(css_noprefix, ")")]

		return image_url, nil
	} else {
		return img.AttrOr("data-original", img.AttrOr("src", "")), nil
	}
}
