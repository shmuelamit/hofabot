package parsers

import (
	"fmt"
	"log"
	"net/url"
	"time"
)

type Show struct {
	Image string
	Name  string
	Url   string
	Desc  string
}

func (s Show) String() string {
	url, err := url.PathUnescape(s.Url)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("🎵 *%s* 🎵\n*%s*\n\n%s", s.Name, url, s.Desc)
}

type RSSConfig struct {
	Url     string
	Image   string
	Refresh time.Duration
}

type GenericConfig struct {
	Url     string
	Title   string
	Desc    string
	Image   string
	Refresh time.Duration
}
