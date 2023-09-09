// https://github.com/ProjectOpenSea/stream-js/blob/main/src/types.ts
package sdk

import (
	"encoding/json"
	"sync"
)

// EventType
const (
	ITEM_METADATA_UPDATED string = "item_metadata_updated"
	ITEM_LISTED                  = "item_listed"
	ITEM_SOLD                    = "item_sold"
	ITEM_TRANSFERRED             = "item_transferred"
	ITEM_RECEIVED_OFFER          = "item_received_offer"
	ITEM_RECEIVED_BID            = "item_received_bid"
	ITEM_CANCELLED               = "item_cancelled"
	COLLECTION_OFFER             = "collection_offer"
	TRAIT_OFFER                  = "trait_offer"
	ITEM_FEACH_ALL               = "*"
)

const (
	MAIN_NET = iota
	TEST_NET
)

type ItemListedRes struct {
	EventTimestamp string `json:"event_timestamp"`
	BasePrice      string `json:"base_price"`
	Collection     struct {
		Slug string `json:"slug"`
	} `json:"collection"`
	ExpirationDate string `json:"expiration_date"`
	IsPrivate      bool   `json:"is_private"`
	ListingDate    string `json:"listing_date"`
	ListingType    string `json:"listing_type"`
	Maker          struct {
		Address string `json:"address"`
	} `json:"maker"`
	PaymentToken struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Item struct {
		Chain struct {
			Name string `json:"name"`
		} `json:"chain"`
		Metadata  ItemListRes `json:"metadata"`
		NFTId     string      `json:"nft_id"`
		Permalink string      `json:"permalink"`
	} `json:"item"`
	Quantity int    `json:"quantity"`
	Taker    string `json:"taker"`
}

type ItemSoldRes struct {
	EventTimestamp string `json:"event_timestamp"`
	ClosingDate    string `json:"closing_date"`
	IsPrivate      bool   `json:"is_private"`
	ListingDate    string `json:"listing_date"`
	ListingType    string `json:"listing_type"`
	Maker          struct {
		Address string `json:"address"`
	} `json:"maker"`
	PaymentToken struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Quantity int `json:"quantity"`
	Taker    struct {
		Address string `json:"address"`
	} `json:"taker"`
	Transaction struct {
		Timestamp string `json:"Timestamp"`
		Hash      string `json:"hash"`
	} `json:"transaction"`
}

type ItemTransferredRes struct {
	EventTimestamp string `json:"event_timestamp"`
	Transaction    struct {
		Timestamp string `json:"Timestamp"`
		Hash      string `json:"hash"`
	} `json:"transaction"`
	FromAccount struct {
		Address string `json:"address"`
	} `json:"from_account"`
	ToAccount struct {
		Address string `json:"address"`
	} `json:"to_account"`
	Quantity int `json:"quantity"`
}

type ItemMetadataUpdatedRes struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ImagePreviewUrl string `json:"image_preview_url"`
	AnimationUrl    string `json:"animation_url"`
	BackgroundColor string `json:"background_color"`
	MetadataUrl     string `json:"metadata_url"`
}

type ItemListRes struct {
	AnimationUrl string `json:"animation_url"`
	ImageUrl     string `json:"image_url"`
	MetadataUrl  string `json:"metadata_url"`
	Name         string `json:"name"`
}

type ItemCancelledRes struct {
	EventTimestamp string `json:"event_timestamp"`
	ListingType    string `json:"listing_type"`
	PaymentToken   struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Quantity    int `json:"quantity"`
	Transaction struct {
		Timestamp string `json:"Timestamp"`
		Hash      string `json:"hash"`
	} `json:"transaction"`
}

type ItemReceivedOfferRes struct {
	EventTimestamp string `json:"event_timestamp"`
	BasePrice      string `json:"base_price"`
	CreatedDate    string `json:"created_date"`
	ExpirationDate string `json:"expiration_date"`
	Item           struct {
		Chain struct {
			Name string `json:"name"`
		} `json:"chain"`
		Metadata  ItemMetadataUpdatedRes `json:"metadata"`
		NFTId     string                 `json:"nft_id"`
		Permalink string                 `json:"permalink"`
	} `json:"item"`
	Maker struct {
		Address string `json:"address"`
	} `json:"maker"`
	PaymentToken struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Quantity int `json:"quantity"`
	Taker    struct {
		Address string `json:"address"`
	} `json:"taker"`
}

type ItemReceivedBidRes struct {
	EventTimestamp string `json:"event_timestamp"`
	BasePrice      string `json:"base_price"`
	CreatedDate    string `json:"created_date"`
	ExpirationDate string `json:"expiration_date"`
	Maker          struct {
		Address string `json:"address"`
	} `json:"maker"`
	Item struct {
		Chain struct {
			Name string `json:"name"`
		} `json:"chain"`
		Metadata  ItemMetadataUpdatedRes `json:"metadata"`
		NFTId     string                 `json:"nft_id"`
		Permalink string                 `json:"permalink"`
	} `json:"item"`
	PaymentToken struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Quantity int `json:"quantity"`
	Taker    struct {
		Address string `json:"address"`
	} `json:"taker"`
}

type CollectionOfferRes struct {
	AssetContractCriteria struct {
		Address string `json:"address"`
	} `json:"asset_contract_criteria"`
	BasePrice  string `json:"base_price"`
	Collection struct {
		Slug string `json:"slug"`
	} `json:"collection"`
	CollectionCriteria struct {
		Slug string `json:"slug"`
	} `json:"collection_criteria"`
	CreatedDate    string `json:"created_date"`
	EventTimestamp string `json:"event_timestamp"`
	ExpirationDate string `json:"expiration_date"`
	Maker          struct {
		Address string `json:"address"`
	} `json:"maker"`
	OrderHash    string `json:"order_hash"`
	PaymentToken struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Quantity int `json:"quantity"`
}

type TraitOfferRes struct {
	AssetContractCriteria struct {
		Address string `json:"address"`
	} `json:"asset_contract_criteria"`
	BasePrice  string `json:"base_price"`
	Collection struct {
		Slug string `json:"slug"`
	} `json:"collection"`
	CollectionCriteria struct {
		Slug string `json:"slug"`
	} `json:"collection_criteria"`
	CreatedDate    string `json:"created_date"`
	EventTimestamp string `json:"event_timestamp"`
	ExpirationDate string `json:"expiration_date"`
	Maker          struct {
		Address string `json:"address"`
	} `json:"maker"`
	OrderHash    string `json:"order_hash"`
	PaymentToken struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Quantity      int `json:"quantity"`
	TraitCriteria struct {
		TraitName string `json:"trait_name"`
		TraitType string `json:"trait_type"`
	} `json:"trait_criteria"`
}

type PayloadJson struct {
	EventType string      `json:"event_type,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
	SentAt    string      `json:"sent_at,omitempty"`
}

type Message struct {
	Topic   string      `json:"topic"`
	Event   string      `json:"event"`
	Payload PayloadJson `json:"payload"`
	Ref     int         `json:"ref"`
}

func (m *Message) UnmarshalJSON(data []byte) error {

	type alice Message
	var temp alice

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	m.Event = temp.Event
	m.Topic = temp.Topic
	m.Ref = temp.Ref
	m.Payload = PayloadJson{}

	payload := PayloadJson{}
	switch temp.Event {
	case "phx_reply":
		return nil
	case "item_metadata_update":
		payload.Payload = &ItemMetadataUpdatedRes{}
	case "item_listed":
		payload.Payload = &ItemListedRes{}
	case "item_sold":
		payload.Payload = &ItemSoldRes{}
	case "item_transferred":
		payload.Payload = &ItemTransferredRes{}
	case "item_metadata_updated":
		payload.Payload = &ItemMetadataUpdatedRes{}
	case "item_cancelled":
		payload.Payload = &ItemCancelledRes{}
	case "item_received_offer":
		payload.Payload = &ItemReceivedOfferRes{}
	case "item_received_bid":
		payload.Payload = &ItemReceivedBidRes{}
	case "collection_offer":
		payload.Payload = &CollectionOfferRes{}
	case "trait_offer":
		payload.Payload = &TraitOfferRes{}
	}

	jsonObj, _ := json.Marshal(temp.Payload)
	json.Unmarshal(jsonObj, &payload)
	m.Payload = payload
	return nil
}

type SafeCounter struct {
	num int
	mux sync.Mutex
}

func (c *SafeCounter) Inc() {
	c.mux.Lock()
	c.num++
	c.mux.Unlock()
}

func (c *SafeCounter) Value() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.num
}
