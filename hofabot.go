package main

import (
	"fmt"
	"os"
	"os/signal"
	"parsers"
	"syscall"
	"time"
)

func main() {
	channel, _ := parsers.GetRSSChannel("https://lorem-rss.herokuapp.com/feed?unit=second", time.Second)

	for {
		item := <-channel
		fmt.Println("item", item)
	}

	return
	client := get_client()

	groups, _ := client.GetJoinedGroups()
	for _, group := range groups {
		if group.OwnerJID == client.Store.ID.ToNonAD() {

		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
