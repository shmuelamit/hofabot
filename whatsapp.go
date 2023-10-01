package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"parsers"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func GetClient() *whatsmeow.Client {
	dbLog := waLog.Stdout("Database", "INFO", true)

	container, err := sqlstore.New("sqlite3", "file:whatsapp.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				// fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	return client
}

func getImageMsg(client *whatsmeow.Client, image string, text string) (*waProto.Message, error) {
	if len(image) == 0 {
		log.Println("Image url empty")
		return nil, errors.New("Image url empty")
	}

	res, err := parsers.GetRequest(image)
	defer res.Body.Close()
	if err != nil {
		log.Println("GET request error", err)
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to upload file", err)
		return nil, err
	}

	uploaded, err := client.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		log.Println("Error uploading message", err)
		return nil, err
	}

	return &waProto.Message{ImageMessage: &waProto.ImageMessage{
		Caption:       &text,
		Url:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(data)),
		FileEncSha256: uploaded.FileEncSHA256,
		FileSha256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(data))),
	}}, nil
}

func SendShow(client *whatsmeow.Client, group *types.GroupInfo, show parsers.Show) {
	text := show.String()
	msg := &waProto.Message{Conversation: &text}

	if imagemsg, err := getImageMsg(client, show.Image, text); err == nil {
		msg = imagemsg
	} else {
		log.Println("Failed to create image message")
	}

	_, err := client.SendMessage(context.Background(), group.JID, msg)
	if err != nil {
		log.Println("Error sending message", err)
	}

	log.Println("SENT MESSAGE ------------------------- ")
}

func GetOwnedGroup(groups []*types.GroupInfo, client_jid types.JID, name string) *types.GroupInfo {
	for _, group := range groups {
		if group.OwnerJID == client_jid && group.GroupName.Name == name {
			return group
		}
	}

	return nil
}
