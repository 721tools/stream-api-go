// event schemas: https://docs.opensea.io/reference/stream-api-event-schemas
// api document: https://docs.opensea.io/reference/using-stream-api-without-sdk
package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
)

// ErrNotConnected is returned when the application read/writes
// a message and the connection is closed
var ErrNotConnected = errors.New("websocket: not connected")
var DefaultNet = MAIN_NET

type MsgHandlerFanc func(msg *Message) error
type UnSubscribeHandlerFanc func()

type payload struct {
	ref   int
	event string
	fn    func(msg *Message) error
}

type notifyService struct {
	mu                sync.Mutex
	endPoint          string
	ref               *SafeCounter
	topic_vendor      sync.Map
	connection        *websocket.Conn
	keepAliveTimeout  time.Duration
	keepAliveResponse keepAliveResponse
	afterConnect      chan bool
	q                 Quit
	// more...
}

func NewNotifyService(network int, token string) *notifyService {
	var endPoint string
	if network == MAIN_NET {
		endPoint = fmt.Sprintf("wss://stream.openseabeta.com/socket/websocket?token=%s", token)
	} else {
		endPoint = fmt.Sprintf("wss://testnets-stream.openseabeta.com/socket/websocket?token=%s", token)
	}

	ns := notifyService{}
	ns.q = *NewQuit()
	ns.q.WatchOsSignal()
	ns.endPoint = endPoint
	ns.ref = &SafeCounter{}
	ns.keepAliveTimeout = time.Second * 20
	ns.keepAliveResponse = keepAliveResponse{}
	ns.afterConnect = make(chan bool)

	return &ns
}

func (ns *notifyService) Subscribe(slug, event string, fn MsgHandlerFanc) (UnSubscribeHandlerFanc, error) {
	topic := collectionTopic(slug)
	if data, ok := ns.topic_vendor.Load(topic); !ok {
		event_chain := make(map[string]payload)
		event_chain[event] = payload{0, event, fn}
		ns.topic_vendor.Store(topic, event_chain)
	} else {
		event_chain := data.(map[string]payload)
		event_chain[event] = payload{0, event, fn}
		ns.topic_vendor.Store(topic, event_chain)
	}

	return func() {
		ns.unsubscribe(topic)
	}, nil
}

func (ns *notifyService) Unsubscribe(slug string) {
	topic := collectionTopic(slug)
	ns.unsubscribe(topic)
}

func (ns *notifyService) Start() error {
	// ws connect
	go ns.dial()
	<-ns.afterConnect

	// join
	go ns.join()
	// reading and processing callback
	go ns.repl()

	if ns.q.IsQuit() {
		return nil
	}

	return nil
}

func (ns *notifyService) Stop() error {
	// leave each topic
	ns.topic_vendor.Range(func(k, v interface{}) bool {
		topic := k.(string)
		ns.unsubscribe(topic)
		fmt.Println("Leave: " + topic)
		return true
	})

	// close ws
	ns.connection.Close()

	return nil
}

func (ns *notifyService) unsubscribe(topic string) {
	ns.leave(topic)
	ns.topic_vendor.Delete(topic)
}

func (ns *notifyService) writeTo(message Message) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	t, _ := json.Marshal(message)
	fmt.Printf("write msg %s \n", string(t))
	err := ns.connection.WriteJSON(message)
	ns.ref.Inc()

	return err
}

// Read, Evaluate, Process, Loop
func (ns *notifyService) repl() error {
	for {
		select {
		default:
			msg := Message{}
			err := ns.connection.ReadJSON(&msg)
			if err != nil {
				log.Println("quit on readFrom when read err:", err, "start reconnect....")
				ns.closeAndReconnect()
			} else {
				ns.keepAliveResponse.setLastResponse()
			}

			// t, _ := json.Marshal(msg)
			// fmt.Printf("recv msg %s \n", string(t))
			// processing msg
			if v, ok := ns.topic_vendor.Load(msg.Topic); ok {

				s := v.(map[string]payload)
				if load, ok := s[msg.Event]; ok {
					load.fn(&msg)
				}
			}
		}
	}
}

// closeAndReconnect will try to reconnect.
func (ns *notifyService) closeAndReconnect() {
	ns.connection.Close()
	go ns.dial()
	<-ns.afterConnect
	go ns.join()
}

func (ns *notifyService) getBackoff() *backoff.Backoff {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	return &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    time.Second * 2,
		Factor: 2,
		Jitter: true,
	}
}

func (ns *notifyService) dial() {
	b := ns.getBackoff()

	go func() {
		for {
			nextItvl := b.Duration()
			c, _, err := websocket.DefaultDialer.Dial(ns.endPoint, nil)
			if err != nil {
				log.Println("Dial: err with", err, " will try again in", nextItvl)
			} else {
				ns.connection = c

				log.Printf("Dial: connection was successfully established with %s\n", ns.endPoint)
				// init keep alive base timestamp
				ns.keepAliveResponse.setLastResponse()

				go ns.heartbeat()

				ns.afterConnect <- true
				break
			}

			time.Sleep(nextItvl)
		}
	}()

	if ns.q.IsQuit() {
		os.Exit(0)
	}
}

func (ns *notifyService) heartbeat() {
	ticker := time.NewTicker(ns.keepAliveTimeout)
	go func() {
		defer ticker.Stop()
		for {
			t := <-ticker.C
			fmt.Println("Tick at", t)
			message := Message{
				Topic:   "phoenix",
				Event:   "heartbeat",
				Payload: PayloadJson{},
				Ref:     ns.ref.Value(),
			}
			ns.writeTo(message)

			// reconnect when the former response timestamp is 2 requests before
			fmt.Println("ts from resp: ", ns.keepAliveResponse.getLastResponse(), "ts now: ", time.Now())
			if time.Since(ns.keepAliveResponse.getLastResponse()) > ns.keepAliveTimeout*2 {
				fmt.Println("heartbeat response timeout")
				ns.closeAndReconnect()
				break
			}
		}
		fmt.Println("heartbeat goroutine out")
	}()

	fmt.Println("heartbeat out")
}

// Unsubscribing from a collection
func (ns *notifyService) leave(topic string) {
	if v, ok := ns.topic_vendor.Load(topic); ok {
		var ref int
		payload := v.(map[string]payload)
		for _, v := range payload {
			ref = v.ref
			break
		}

		message := Message{
			Topic:   topic,
			Event:   "phx_leave",
			Payload: PayloadJson{},
			Ref:     ref,
		}
		ns.writeTo(message)
	}
}

// Subscribing to a collection
func (ns *notifyService) join() {
	if v, ok := ns.topic_vendor.Load("collection:*"); ok {
		ref := ns.ref.Value()
		payload := v.(map[string]payload)
		for _, v := range payload {
			v.ref = ref
		}

		message := Message{
			Topic:   "collection:*",
			Event:   "phx_join",
			Payload: PayloadJson{},
			Ref:     ref,
		}
		fmt.Println("Start: collection:*")
		ns.writeTo(message)
	} else {
		// or join each slug's topic
		ns.topic_vendor.Range(func(k, v interface{}) bool {
			topic := k.(string)
			ref := ns.ref.Value()

			payload := v.(map[string]payload)
			for _, v := range payload {
				v.ref = ref
			}
			message := Message{
				Topic:   topic,
				Event:   "phx_join",
				Payload: PayloadJson{},
				Ref:     ref,
			}
			ns.writeTo(message)

			fmt.Println("Start: " + topic)
			return true
		})
	}
}
