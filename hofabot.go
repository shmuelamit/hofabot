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

// func RegisterChat(config Config, client *whatsmeow.Client, groups)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Llongfile)

	var config Config
	_, err := toml.DecodeFile(CONFIG_FILE, &config)
	if err != nil {
		log.Fatal("Error decoding config file", err)
	}

	client := GetClient()
	groups, err := client.GetJoinedGroups()
	if err != nil {
		log.Fatal("Error getting cleint groups", err)
	}

	for _, config := range config.Rss {
		channel, _ := parsers.GetRSSChannel(config.RSSConfig)
		group := GetOwnedGroup(groups, client.Store.ID.ToNonAD(), config.Group)
		if group == nil {
			log.Fatal("no group named " + config.Group)
		}

		println("init", config.Url)

		go func() {
			for {
				show := <-channel
				log.Println("Got show", show.Name, show.Url)
				SendShow(client, group, show)
			}
		}()
	}

	for _, config := range config.Generic {
		channel, _ := parsers.GetGenericChannel(config.GenericConfig)
		group := GetOwnedGroup(groups, client.Store.ID.ToNonAD(), config.Group)
		if group == nil {
			log.Fatal("no group named " + config.Group)
		}

		go func() {
			for {
				show := <-channel
				log.Println("Got show", show.Name, show.Url)
				SendShow(client, group, show)
			}
		}()
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
