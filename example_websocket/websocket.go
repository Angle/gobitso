package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/angle/gobitso"
)

func main() {
	bitsoClient := bitso.NewClient()

	// Pull books & limits from the API
	log.Println("Pulling books..")
	books, err := bitsoClient.AvailableBooks()
	if err != nil {
		log.Fatalf("Error pulling books from Bitso: %v", err)
	}

	log.Printf("Loaded %d books", len(books))


	///////////////////////////////
	// INITIALIZE WEBSOCKET BITSO CLIENT
	bitsoWs := bitso.NewWebsocketListener()

	feed, err := bitsoWs.Connect()

	if err != nil {
		log.Fatal("Error connecting to Bitso's websocket")
	}

	// subscribe to _every_ book
	for bookCode := range books {
		err = bitsoWs.Subscribe(bookCode, bitso.Channel_ORDERS)
		if err != nil {
			log.Fatalf("Error subscribing to '%s' Orders channel", bookCode)
		}

		err = bitsoWs.Subscribe(bookCode, bitso.Channel_TRADES)
		if err != nil {
			log.Fatalf("Error subscribing to '%s' Trades channel", bookCode)
		}
	}

	// Catch interrupts (Ctrl-c)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start main websocket listener loop
	MainLoop:
	for {
		select {
		case msg := <-feed:
			// Incoming message from the websocket!
			switch msg.Channel {
			case bitso.Channel_DISCONNECTED:
				log.Println("Bitso websocket was disconnected!")
				break MainLoop

			case bitso.Channel_ORDERS:
				orders, ok := msg.Payload.(bitso.Orders)

				if !ok {
					// invalid payload in orders
					log.Println("invalid orders payload")
					continue
				}

				// Do something with the incoming orders payload for the book
				out := fmt.Sprintf("[%s] ", msg.Book)


				// Example: print the Best Bid / Ask
				if len(orders.Bids) > 0 {
					out += fmt.Sprintf("BID: %s @ $%s", orders.Bids[0].Amount.StringFixed(8), orders.Bids[0].Rate.StringFixed(8))
				} else {
					out += "BID: --"
				}

				out += " | "

				if len(orders.Asks) > 0 {
					out += fmt.Sprintf("ASK: %s @ $%s", orders.Asks[0].Amount.StringFixed(8), orders.Asks[0].Rate.StringFixed(8))
				} else {
					out += "ASK: --"
				}

				log.Println(out)



			case bitso.Channel_TRADES:
				trades, ok := msg.Payload.([]bitso.Trade)

				if !ok {
					// invalid payload in trades
					log.Println("invalid trades payload")
					continue
				}

				// Do something with the incoming trades payload for the book
				// Example: print the Trade
				for _, trade := range trades {
					log.Printf("[%s] TRADE: [%s] %s @ $%s = $%s", msg.Book, trade.Side.String(), trade.Amount.StringFixed(8), trade.Rate.StringFixed(8), trade.Value.StringFixed( 2))
				}

			}
		case <-interrupt:
			log.Println("interrupt!")
			bitsoWs.Disconnect()
			break MainLoop
		}
	}

	log.Println("Program ended.")
}