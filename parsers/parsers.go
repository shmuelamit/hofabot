package parsers

import (
	"fmt"
	"time"
)

type Show struct {
	Image string
	Name  string
	Url   string
	Desc  string
}

func (s Show) String() string {
	return fmt.Sprintf("%s\n--------\n%s\n\n%s", s.Name, s.Url, s.Desc)
}

type RSSConfig struct {
	Url     string
	Refresh time.Duration
}

type GenericConfig struct {
	Url     string
	Title   string
	Desc    string
	Image   string
	Refresh time.Duration
}
