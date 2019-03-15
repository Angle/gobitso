package bitso

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const WEBSOCKET_ENDPOINT = "wss://ws.bitso.com"
const LOG_PREFIX = "BitsoWebsocket: "

const MAX_FEED_QUEUE_SIZE = 1000

type Websocket struct {
	conn 		*websocket.Conn

	feed 		chan FeedMessage

	quit 		chan bool
	quitOnce 	*sync.Once
}

// NewWebsocketListener returns a pointer to a new instance of the WebsocketListener
func NewWebsocketListener() *Websocket {
	return &Websocket{
		feed: make(chan FeedMessage, MAX_FEED_QUEUE_SIZE),
		quit: make(chan bool),
		quitOnce: new(sync.Once),
	}
}

// Connect establishes the initial connection to the websocket, must be called before subscribing to a channel
func (ws *Websocket) Connect() (<-chan FeedMessage, error) {
	log.Printf(LOG_PREFIX + "connecting to %s", WEBSOCKET_ENDPOINT)

	conn, _, err := websocket.DefaultDialer.Dial(WEBSOCKET_ENDPOINT, nil)
	if err != nil {
		return nil, NewWebSocketError(fmt.Sprintf("error on dial: %v", err))
	}

	ws.conn = conn

	log.Printf(LOG_PREFIX + "connected!")

	// Launch the reader process
	go ws.reader()

	// Launch the writer process
	go ws.writer()

	// Pass down the FeedMessage channel to the consumer
	return ws.feed, nil
}

// Subscribe subscribes to the specified channel for an specific book
func (ws *Websocket) Subscribe(book BookCode, channel Channel) error {
	if ws.conn == nil {
		return NewWebSocketError("websocket connection has not been initialized yet")
	}

	subscribePayload, err := json.Marshal(SubscribeRequestMessage{
		Action: ActionType_SUBSCRIBE,
		Book: book,
		Channel: channel,
	})

	if err != nil {
		return NewWebSocketError(fmt.Sprintf("subscribe message build failed: %v", err))
	}

	err = ws.conn.WriteMessage(websocket.TextMessage, subscribePayload)
	if err != nil {
		return NewWebSocketError(fmt.Sprintf("subscribe message send failed: %v", err))
	}

	return nil
}

// Disconnect closes the current Websocket connection as cleanly as possible
func (ws *Websocket) Disconnect() {
	ws.quitOnce.Do(func() {
		// Close the quit channel
		close(ws.quit)

		// when the 'quit' channel is closed, the writer should attempt a clean disconnect
		// we'll wait a little to allow that last message to be sent
		time.Sleep(time.Second)

		// Now we can close the connection.
		if err := ws.conn.Close(); err != nil {
			// Failed to properly close the connection
			// TODO: verbose error?
		}

		// send a last message to the Feed to notify upstream of the disconnection
		ws.sendFeedMessage(FeedMessage{Channel: Channel_DISCONNECTED})
	})
}

func (ws *Websocket) reader() {
	ReadLoop:
	for {
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			log.Println(LOG_PREFIX + "read:", err)
			break ReadLoop
		}
		//log.Printf("recv: %s", message)

		// Parse the incoming message
		incoming := IncomingMessage{}
		err = json.Unmarshal(message, &incoming)
		if err != nil {
			log.Println(LOG_PREFIX + "unknown incoming message format", err)
			break ReadLoop
		}

		//log.Println(incoming.Channel, incoming.Action)

		switch incoming.Channel {
		case Channel_KEEP_ALIVE:
			// received a server heartbeat, all is good, do nothing.
			continue ReadLoop

		case Channel_ORDERS:
			switch incoming.Action {
			case ActionType_SUBSCRIBE:
				log.Println(LOG_PREFIX + "ORDERS subscription ok!")
			case ActionType_NULL:
				// no action was specified, therefore it's a regular Orders Channel message
				ordersPayload := Orders{}
				err = json.Unmarshal(*incoming.Payload, &ordersPayload)

				// pass down the orders message
				ws.sendFeedMessage(FeedMessage{
					Channel: Channel_ORDERS,
					Book: incoming.Book,
					Payload: ordersPayload,
				})
			default:
				// woah, what happened? unknown action!
				break ReadLoop
			}

		case Channel_TRADES:
			switch incoming.Action {
			case ActionType_SUBSCRIBE:
				log.Println(LOG_PREFIX + "TRADES subscription ok!")
			case ActionType_NULL:
				// no action was specified, therefore it's a regular Trades Channel message
				tradesPayload := make([]Trade, 0)
				err = json.Unmarshal(*incoming.Payload, &tradesPayload)

				// pass down the orders message
				ws.sendFeedMessage(FeedMessage{
					Channel: Channel_TRADES,
					Book: incoming.Book,
					Payload: tradesPayload,
				})
			default:
				// woah, what happened? unknown action!
				break ReadLoop
			}

		default:
			// unknown channel, not yet implemented
			log.Printf(LOG_PREFIX + "unknown channel '%s'", string(incoming.Channel))
			break ReadLoop
		}
	}

	// if the loop breaks it means that a message failed to be read or parsed
	// and we'll consider that as a connection failure.
	// we'll throw a Disconnect() for good measure
	ws.Disconnect()
}

func (ws *Websocket) sendFeedMessage(m FeedMessage) {
	// attempt to send a FeedMessage upstream
	select {
	case ws.feed <- m:
	default:
		// the queue is full, we'll disconnect ourselves, the upstream is not responding.
		log.Println(LOG_PREFIX + "feed message queue is full")
		ws.Disconnect()
	}
}

func (ws *Websocket) writer() {
	ticker := time.NewTicker(time.Second) // keep-alive
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			err := ws.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println(LOG_PREFIX + "write:", err)
				return
			}
		case <-ws.quit:
			// this will only happen
			log.Println(LOG_PREFIX + "quit, attempting clean disconnect")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println(LOG_PREFIX + "write close:", err)
			}

			return
		}
	}
}