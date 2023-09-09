package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/721tools/stream-api-go/sdk"

	flag "github.com/spf13/pflag"
)

var mainnet bool
var token string

func init() {
	flag.BoolVar(&mainnet, "net", true, "Please Select the Network to Connect")
	flag.StringVar(&token, "key", "", "Please provide an API key for authentication")
}

func main() {
	flag.Parse()

	if token == "" {
		log.Fatal("An API key is required for connection. If you do not have an API key, please request one from OpenSea." + token)
	}

	net := sdk.MAIN_NET
	if !mainnet {
		// todo log here
		net = sdk.TEST_NET
	}

	log.SetFlags(0)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ns := sdk.NewNotifyService(net, token)

	ns.Subscribe("*", sdk.ITEM_LISTED, func(msg *sdk.Message) error {
		pl := msg.Payload.Payload.(*sdk.ItemListedRes)
		fmt.Printf("ITEM_LISTED: msg.Payload.Payload: %v\n", pl.Item.Metadata.ImageUrl)
		return nil
	})

	ns.Subscribe("*", sdk.COLLECTION_OFFER, func(msg *sdk.Message) error {
		pl := msg.Payload.Payload.(*sdk.CollectionOfferRes)
		fmt.Printf("COLLECTION_OFFER: msg.Payload.Payload: %v\n", pl.Collection.Slug)
		return nil
	})

	ns.Start()
}
