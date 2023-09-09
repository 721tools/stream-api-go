## OpenSea stream API unofficial implement in Go

> https://docs.opensea.io/reference/using-stream-api-without-sdk

## How to use

```
package main

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/721tools/stream-api-go/sdk"
)

func main() {
	api_token := "YOUR-API-TOKEN-HERE"
	ns := sdk.NewNotifyService(sdk.MAIN_NET, api_token)

	ns.Subscribe("*", sdk.ITEM_LISTED, func(msg *sdk.Message) error {
		t, _ := json.Marshal(msg)
		fmt.Printf("recv msg %s \n", string(t))
		return nil
	})

	go ns.Start()

	time.Sleep(time.Second * 100)
}


// ns.Stop()

```
