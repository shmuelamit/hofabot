package main

import (
	"log"
	"os"
	"os/signal"
	"parsers"
	"syscall"

	"github.com/BurntSushi/toml"
)

const CONFIG_FILE = "config.toml"

type Config struct {
	Rss []struct {
		Group string
		parsers.RSSConfig
	}

	Generic []struct {
		Group string
		parsers.GenericConfig
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var config Config
	_, err := toml.DecodeFile(CONFIG_FILE, &config)
	if err != nil {
		log.Fatal(err)
	}

	// channel, _ := parsers.GetGenericChannel(config.Generic[0].GenericConfig)

	// for {
	// 	_ = <-channel
	// }

	// return
	client := GetClient()
	groups, _ := client.GetJoinedGroups()

	ozen := GetOwnedGroup(groups, client.Store.ID.ToNonAD(), "אוזן")
	println("yes", ozen.Name, parsers.GetFutureShows())

	// for _, config := range config.Rss {
	// 	// channel, _ := parsers.GetRSSChannel(config.RSSConfig)
	// 	group := GetOwnedGroup(groups, client.Store.ID.ToNonAD(), config.Group)
	// 	if group == nil {
	// 		log.Fatal("no group named " + config.Group)
	// 	}
	// }

	// for _, config := range config.Generic {
	// 	// channel, _ := parsers.GetRSSChannel(config.GenericConfig)
	// 	group := GetOwnedGroup(groups, client.Store.ID.ToNonAD(), config.Group)
	// 	if group == nil {
	// 		log.Fatal("no group named " + config.Group)
	// 	}
	// }

	println("AHAAAAAAAAAAA")
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
